package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type CinemaRepository interface{
	GetAll(q dto.PaginationQuery) ([]entity.Cinema, int, error)
	GetByID(id int) (*entity.Cinema, error)
	Create(cinema entity.Cinema) (*entity.Cinema, error)
	Update(id int, w *entity.Cinema) error
	Delete(id int) error
}

type cinemaRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewCinemaRepository(db database.PgxIface, log *zap.Logger) CinemaRepository {
	return &cinemaRepository{
		db:     db,
		Logger: log,
	}
}

func (r *cinemaRepository) GetAll(q dto.PaginationQuery) ([]entity.Cinema, int, error) {
	var offset int
	offset = (q.Page - 1) * q.Limit
	
	// Get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM cinemas WHERE deleted_at IS NULL`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("Error query count cinemas: ", zap.Error(err))
		return nil, 0, err
	}

	// Initiate rows
	var rows pgx.Rows
	
	// Conditional query based on page, limit, and all param
	query := `SELECT id, name, location, created_at, updated_at FROM cinemas WHERE deleted_at IS NULL`

	if !q.All && q.Limit > 0 {
		query += ` LIMIT $1 OFFSET $2`
		rows, err = r.db.Query(context.Background(), query, q.Limit, offset)
	} else {
		rows, err = r.db.Query(context.Background(), query)
	}
	
	if err != nil {
		r.Logger.Error("Error query get all cinemas: ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var cinemas []entity.Cinema
	for rows.Next() {
		var w entity.Cinema
		err := rows.Scan(&w.ID, &w.Name, &w.Location, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			r.Logger.Error("Error scan cinema: ", zap.Error(err))
			return nil, 0, err
		}
		cinemas = append(cinemas, w)
	}
	return cinemas, total, nil
}

func (r *cinemaRepository) GetByID(id int) (*entity.Cinema, error) {
	var cinema entity.Cinema
	query := "SELECT id, name, location, created_at, updated_at FROM cinemas WHERE id = $1 AND deleted_at IS NULL"

	err := r.db.QueryRow(context.Background(), query, id).Scan(&cinema.ID, &cinema.Name, &cinema.Location, &cinema.CreatedAt, &cinema.UpdatedAt)

	if err == pgx.ErrNoRows {
		r.Logger.Error("Error not found cinema: ", zap.Error(err))
		return nil, utils.ErrNotFound("cinema")
	}

	if err != nil {
		r.Logger.Error("Error query get cinema by id: ", zap.Error(err))
		return nil, err
	}

	return &cinema, nil
}

func (r *cinemaRepository) Create(cinema entity.Cinema) (*entity.Cinema, error) {
	query := `
		INSERT INTO cinemas (name, location, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, cinema.Name, cinema.Location).Scan(&cinema.ID)
	if err != nil {
		r.Logger.Error("Error query create cinema: ", zap.Error(err))
		return nil, err
	}

	cinema.CreatedAt = time.Now()
	cinema.UpdatedAt = time.Now()
	return &cinema, nil
}

func (r *cinemaRepository) Update(id int, w *entity.Cinema) error {
	query := `
		UPDATE cinemas
		SET name = COALESCE($1, name),
		updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query,
		&w.Name, id,
	)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found cinema: ", zap.Error(err))
		return utils.ErrNotFound("cinema")
	}

	if err != nil {
		r.Logger.Error("Error query update cinema: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *cinemaRepository) Delete(id int) error {
	query := `
		UPDATE cinemas
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query, id)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found cinema: ", zap.Error(err))
		return utils.ErrNotFound("cinema")
	}

	if err != nil {
		r.Logger.Error("Error query delete cinema: ", zap.Error(err))
		return err
	}

	return nil
}