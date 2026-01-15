package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type StudioRepository interface{
	GetAll(q dto.PaginationQuery) ([]entity.Studio, int, error)
	GetByID(id int) (*entity.Studio, error)
	Create(data dto.StudioRequest) (*entity.Studio, error)
	Update(id int, data dto.StudioRequest) error
	Delete(id int) error
}

type studioRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewStudioRepository(db database.PgxIface, log *zap.Logger) StudioRepository {
	return &studioRepository{
		db:     db,
		Logger: log,
	}
}

func (r *studioRepository) GetAll(q dto.PaginationQuery) ([]entity.Studio, int, error) {
	var offset int
	offset = (q.Page - 1) * q.Limit
	
	// Get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM studios WHERE deleted_at IS NULL`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("Error query count studios: ", zap.Error(err))
		return nil, 0, err
	}

	// Initiate rows
	var rows pgx.Rows
	
	// Conditional query based on page, limit, and all param
	query := `SELECT id, cinema_id, name, type, created_at, updated_at FROM studios WHERE deleted_at IS NULL`

	if !q.All && q.Limit > 0 {
		query += ` LIMIT $1 OFFSET $2`
		rows, err = r.db.Query(context.Background(), query, q.Limit, offset)
	} else {
		rows, err = r.db.Query(context.Background(), query)
	}
	
	if err != nil {
		r.Logger.Error("Error query get all studios: ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var studios []entity.Studio
	for rows.Next() {
		var s entity.Studio
		err := rows.Scan(&s.ID, &s.CinemaID, &s.Name, &s.Type, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			r.Logger.Error("Error scan studio: ", zap.Error(err))
			return nil, 0, err
		}
		studios = append(studios, s)
	}
	return studios, total, nil
}

func (r *studioRepository) GetByID(id int) (*entity.Studio, error) {
	var studio entity.Studio
	query := "SELECT id, cinema_id, name, type, created_at, updated_at FROM studios WHERE id = $1 AND deleted_at IS NULL"

	err := r.db.QueryRow(context.Background(), query, id).Scan(&studio.ID, &studio.CinemaID, 
		&studio.Name, &studio.Type, &studio.CreatedAt, &studio.UpdatedAt)

	if err == pgx.ErrNoRows {
		r.Logger.Error("Error not found studio: ", zap.Error(err))
		return nil, utils.ErrNotFound("studio")
	}

	if err != nil {
		r.Logger.Error("Error query get studio by id: ", zap.Error(err))
		return nil, err
	}

	return &studio, nil
}

func (r *studioRepository) Create(data dto.StudioRequest) (*entity.Studio, error) {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background()); 
	if err != nil {
		r.Logger.Error("Error start db transaction: ", zap.Error(err))
		return nil, err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit(context.Background())
			r.Logger.Error("Error commit db transaction: ", zap.Error(err))
		default:
			_ = tx.Rollback(context.Background())
		}
	}()
	
	// Create studio
	query := `
		INSERT INTO studios (cinema_id, name, type, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	var studio entity.Studio
	err = tx.QueryRow(context.Background(), query, data.CinemaID, data.Name, data.Type).Scan(&studio.ID)
	if err != nil {
		r.Logger.Error("Error query create studio: ", zap.Error(err))
		return nil, err
	}

	// Create seats
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P"}
	studioType := make(map[string][]int)
	studioType["Regular"] = []int{14, 10}
	studioType["Large"] = []int{18, 14}
	studioType["Premium"] = []int{8, 6}
	seats, ok := studioType[data.Type]
	if !ok {
		r.Logger.Error("Invalid studio type: ", zap.Error(err))
		return nil, err
	}
	for i:=0; i<seats[1]; i++ {
		for j:=1; j<=seats[0]; j++ {
			seatCode := fmt.Sprintf("%s%d", letters[i], j)
			query = `INSERT INTO seats (studio_id, seat_code, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())`
			if _, err = tx.Exec(context.Background(), query, studio.ID, seatCode); err != nil {
				r.Logger.Error("Error query create seat: ", zap.Error(err))
				return nil, err
			}
		}
	}

	studio.CinemaID = data.CinemaID
	studio.Name = data.Type
	studio.Type = data.Type
	studio.CreatedAt = time.Now()
	studio.UpdatedAt = time.Now()
	return &studio, nil
}

func (r *studioRepository) Update(id int, s dto.StudioRequest) error {
	query := `
		UPDATE studios
		SET cinema_id = COALESCE($1, cinema_id)
		name = COALESCE($2, name)
		updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query,
		&s.CinemaID, &s.Name, id,
	)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found studio: ", zap.Error(err))
		return utils.ErrNotFound("studio")
	}

	if err != nil {
		r.Logger.Error("Error query update studio: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *studioRepository) Delete(id int) error {
	query := `
		UPDATE studios
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query, id)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found studio: ", zap.Error(err))
		return utils.ErrNotFound("studio")
	}

	if err != nil {
		r.Logger.Error("Error query delete studio: ", zap.Error(err))
		return err
	}

	return nil
}