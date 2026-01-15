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

type MovieRepository interface{
	GetAll(q dto.MovieQuery) ([]entity.Movie, int, error)
	GetByID(id int) (*dto.MovieResponse, error)
	Create(data dto.MovieRequest) (int, error)
	Update(id int, data dto.MovieRequest) error
	Delete(id int) error
}

type movieRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewMovieRepository(db database.PgxIface, log *zap.Logger) MovieRepository {
	return &movieRepository{
		db:     db,
		Logger: log,
	}
}

func (r *movieRepository) GetAll(q dto.MovieQuery) ([]entity.Movie, int, error) {
	var offset int
	offset = (q.Page - 1) * q.Limit
	
	// Get total data for pagination
	var total int
	countQuery := `
	SELECT COUNT(DISTINCT m.id)
	FROM movies m
	WHERE m.deleted_at IS NULL
	AND (
		$1 = '' OR
		EXISTS (
			SELECT 1
			FROM genre_movies gm
			JOIN genres g ON g.id = gm.genre_id
			WHERE gm.movie_id = m.id
				AND g.name = $1
		)
	)`
	err := r.db.QueryRow(context.Background(), countQuery, q.Genre).Scan(&total)
	if err != nil {
		r.Logger.Error("Error query count movies: ", zap.Error(err))
		return nil, 0, err
	}

	// Initiate rows
	var rows pgx.Rows
	
	// Conditional query based on page, limit, and all param
	query := `SELECT m.id, title, synopsis, poster_url, 
	trailer_url, duration_minute, release_date, language, 
	rating_age, ARRAY_AGG(g.name) AS genres, m.created_at, m.updated_at
	FROM movies m 
	LEFT JOIN genre_movies gm ON gm.movie_id = m.id
	LEFT JOIN genres g ON g.id = gm.genre_id 
	WHERE m.deleted_at IS NULL
	AND (
		$1 = '' OR
		EXISTS (
			SELECT 1
			FROM genre_movies gm2
			JOIN genres g2 ON g2.id = gm2.genre_id
			WHERE gm2.movie_id = m.id
				AND g2.name = $1
		)
	)
	GROUP BY m.id
	ORDER BY m.created_at DESC
	`

	if !q.All && q.Limit > 0 {
		query += ` LIMIT $2 OFFSET $3`
		rows, err = r.db.Query(context.Background(), query, q.Genre, q.Limit, offset)
	} else {
		rows, err = r.db.Query(context.Background(), query, q.Genre)
	}
	
	if err != nil {
		r.Logger.Error("Error query get all movies: ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var movies []entity.Movie
	for rows.Next() {
		var m entity.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Synopsis, &m.PosterURL, &m.TrailerURL,
			&m.Duration, &m.ReleaseDate, &m.Language, &m.RatingAge,
			&m.Genres, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			r.Logger.Error("Error scan movie: ", zap.Error(err))
			return nil, 0, err
		}
		movies = append(movies, m)
	}
	return movies, total, nil
}

func (r *movieRepository) GetByID(id int) (*dto.MovieResponse, error) {
	var m dto.MovieResponse
	query := `SELECT m.id, title, synopsis, poster_url, 
	trailer_url, duration_minute, release_date, language, 
	rating_age, ARRAY_AGG(g.name) AS genres
	FROM movies m 
	LEFT JOIN genre_movies gm ON gm.movie_id = m.id
	LEFT JOIN genres g ON g.id = gm.genre_id 
	WHERE m.deleted_at IS NULL
	AND m.id = $1
	GROUP BY m.id`

	err := r.db.QueryRow(context.Background(), query, id).Scan(&m.MovieID, 
		&m.Title, &m.Synopsis, &m.PosterURL, &m.TrailerURL,
		&m.Duration, &m.ReleaseDate, &m.Language,
		&m.RatingAge, &m.Genres)

	if err == pgx.ErrNoRows {
		r.Logger.Error("Error not found movie: ", zap.Error(err))
		return nil, utils.ErrNotFound("movie")
	}

	if err != nil {
		r.Logger.Error("Error query get movie by id: ", zap.Error(err))
		return nil, err
	}

	return &m, nil
}

func (r *movieRepository) Create(m dto.MovieRequest) (int, error) {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background()); 
	if err != nil {
		r.Logger.Error("Error start db transaction: ", zap.Error(err))
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	var movie entity.Movie
	// Convert string to datetime
	releaseDate, err := time.Parse("2006-01-02", m.ReleaseDate)
	
	// Create movie
	query := `
		INSERT INTO movies (title, synopsis, poster_url, 
		trailer_url, duration_minute, release_date, language, 
		rating_age, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id
	`
	err = tx.QueryRow(context.Background(), query,
		&m.Title, &m.Synopsis, &m.PosterURL, &m.TrailerURL,
		&m.Duration, &releaseDate, &m.Language,
		&m.RatingAge).Scan(&movie.ID)
	if err != nil {
		r.Logger.Error("Error query create movie: ", zap.Error(err))
		return 0, err
	}

	for _, genreID := range m.Genres {
		fmt.Println(genreID, movie.ID)
		query = `INSERT INTO genre_movies (genre_id, movie_id) VALUES ($1, $2)`
		_, err = tx.Exec(context.Background(), query, genreID, movie.ID)
		if err != nil {
			r.Logger.Error("Error query create genre_movies: ", zap.Error(err))
			return 0, err
		}
	}

	return movie.ID, nil
}

func (r *movieRepository) Update(id int, m dto.MovieRequest) error {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background()); 
	if err != nil {
		r.Logger.Error("Error start db transaction: ", zap.Error(err))
		return err
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

	// Convert string to datetime
	releaseDate, err := time.Parse("2006-01-02", m.ReleaseDate)

	// Update movies
	query := `
		UPDATE movies
		SET title = $1, 
		synopsis = $2, 
		poster_url = $3, 
		trailer_url = $4, 
		duration_minute = $5, 
		release_date = $6,  
		language = $7, 
		rating_age = $8,
		updated_at = NOW()
		WHERE id = $9 AND deleted_at IS NULL
	`

	result, err := tx.Exec(context.Background(), query,
		&m.Title, &m.Synopsis, &m.PosterURL, &m.TrailerURL,
		&m.Duration, &releaseDate, &m.Language,
		&m.RatingAge, id,
	)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found movie: ", zap.Error(err))
		return utils.ErrNotFound("movie")
	}

	if err != nil {
		r.Logger.Error("Error query update movie: ", zap.Error(err))
		return err
	}

	// Delete genre_movies
	query = `DELETE FROM genre_movies WHERE movie_id = $1`
	result, err = tx.Exec(context.Background(), query, id)
	rowsAffected = result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found movie: ", zap.Error(err))
		return utils.ErrNotFound("movie")
	}

	if err != nil {
		r.Logger.Error("Error query delete genre_movies: ", zap.Error(err))
		return err
	}
	
	query = `INSERT INTO genre_movies (movie_id, genre_id) SELECT $1, unnest($2::int[]);`
	_, err = tx.Exec(context.Background(), query, id, m.Genres)

	if err != nil {
		r.Logger.Error("Error query create genre_movies: ", zap.Error(err))
		return err
	}

	return nil
}

func (r *movieRepository) Delete(id int) error {
	query := `
		UPDATE movies
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query, id)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found movie: ", zap.Error(err))
		return utils.ErrNotFound("movie")
	}

	if err != nil {
		r.Logger.Error("Error query delete movie: ", zap.Error(err))
		return err
	}

	return nil
}