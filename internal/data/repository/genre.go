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

type GenreRepository interface{
	GetAll(q dto.PaginationQuery) ([]entity.Genre, int, error)
	GetByID(id int) (*entity.Genre, error)
	Create(genre entity.Genre) (*entity.Genre, error)
	Delete(id int) error
}

type genreRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewGenreRepository(db database.PgxIface, log *zap.Logger) GenreRepository {
	return &genreRepository{
		db:     db,
		Logger: log,
	}
}

func (r *genreRepository) GetAll(q dto.PaginationQuery) ([]entity.Genre, int, error) {
	var offset int
	offset = (q.Page - 1) * q.Limit
	
	// Get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM genres WHERE deleted_at IS NULL`
	err := r.db.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("Error query count genres: ", zap.Error(err))
		return nil, 0, err
	}

	// Initiate rows
	var rows pgx.Rows
	
	// Conditional query based on page, limit, and all param
	query := `SELECT id, name, created_at, updated_at FROM genres WHERE deleted_at IS NULL`

	if !q.All && q.Limit > 0 {
		query += ` LIMIT $1 OFFSET $2`
		rows, err = r.db.Query(context.Background(), query, q.Limit, offset)
	} else {
		rows, err = r.db.Query(context.Background(), query)
	}
	
	if err != nil {
		r.Logger.Error("Error query get all genres: ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var genres []entity.Genre
	for rows.Next() {
		var g entity.Genre
		err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			r.Logger.Error("Error scan genre: ", zap.Error(err))
			return nil, 0, err
		}
		genres = append(genres, g)
	}
	return genres, total, nil
}

func (r *genreRepository) GetByID(id int) (*entity.Genre, error) {
	var genre entity.Genre
	query := "SELECT id, name, created_at, updated_at FROM genres WHERE id = $1 AND deleted_at IS NULL"

	err := r.db.QueryRow(context.Background(), query, id).Scan(&genre.ID, &genre.Name, &genre.CreatedAt, &genre.UpdatedAt)

	if err == pgx.ErrNoRows {
		r.Logger.Error("Error not found genre: ", zap.Error(err))
		return nil, utils.ErrNotFound("genre")
	}

	if err != nil {
		r.Logger.Error("Error query get genre by id: ", zap.Error(err))
		return nil, err
	}

	return &genre, nil
}

func (r *genreRepository) Create(genre entity.Genre) (*entity.Genre, error) {
	query := `
		INSERT INTO genres (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, genre.Name).Scan(&genre.ID)
	if err != nil {
		r.Logger.Error("Error query create genre: ", zap.Error(err))
		return nil, err
	}

	genre.CreatedAt = time.Now()
	genre.UpdatedAt = time.Now()
	return &genre, nil
}

func (r *genreRepository) Delete(id int) error {
	query := `
		UPDATE genres
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query, id)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found genre: ", zap.Error(err))
		return utils.ErrNotFound("genre")
	}

	if err != nil {
		r.Logger.Error("Error query delete genre: ", zap.Error(err))
		return err
	}

	return nil
}