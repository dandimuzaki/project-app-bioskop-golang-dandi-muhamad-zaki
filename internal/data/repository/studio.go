package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type StudioRepository interface{
	CreateStudioType(req dto.StudioType) (*entity.StudioType, error)
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

func (r *studioRepository) CreateStudioType(req dto.StudioType) (*entity.StudioType, error) {
	var t entity.StudioType
	query := `INSERT INTO studio_types (name, row, column, price)
	VALUES ($1, $2, $3, $4) RETURNING id, name, row, column, price`
	err := r.db.QueryRow(context.Background(), query, req.Name, req.Row, req.Column, req.Price).Scan(&t.ID, &t.Name, &t.Row, &t.Column, &t.Price)
	if err != nil {
		r.Logger.Error("Error query create studio type: ", zap.Error(err))
		return nil, err
	}
	return &t, nil
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
	query := `SELECT s.id, cinema_id, s.name, t.name AS type, t.price, created_at, updated_at
	FROM studios s
	LEFT JOIN studio_types t ON s.type = t.id
	WHERE deleted_at IS NULL`

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
	query := `SELECT s.id, cinema_id, s.name, t.name AS type, t.price, created_at, updated_at
	FROM studios s
	LEFT JOIN studio_types t ON s.type = t.id
	WHERE s.id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(context.Background(), query, id).Scan(&studio.ID, &studio.CinemaID, 
		&studio.Name, &studio.Type, &studio.Price, &studio.CreatedAt, &studio.UpdatedAt)

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
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
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

	// Get studio
	query = `SELECT s.id, cinema_id, s.name, t.name AS type, t.row, 
	t.column, t.price, created_at, updated_at
	FROM studios s
	LEFT JOIN studio_types t ON s.type = t.id
	WHERE id = $1 AND deleted_at IS NULL`
	err = tx.QueryRow(context.Background(), query, studio.ID).Scan(&studio.ID, &studio.CinemaID, 
		&studio.Name, &studio.Type, &studio.Row, &studio.Column, &studio.Price, &studio.CreatedAt, &studio.UpdatedAt)
	if err != nil {
		r.Logger.Error("Error query create studio: ", zap.Error(err))
		return nil, err
	}

	// Create seats
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P"}
	for i:=0; i<*studio.Row; i++ {
		for j:=1; j<=*studio.Column; j++ {
			seatCode := fmt.Sprintf("%s%d", letters[i], j)
			query = `INSERT INTO seats (studio_id, seat_code, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())`
			if _, err = tx.Exec(context.Background(), query, studio.ID, seatCode); err != nil {
				r.Logger.Error("Error query create seat: ", zap.Error(err))
				return nil, err
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

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