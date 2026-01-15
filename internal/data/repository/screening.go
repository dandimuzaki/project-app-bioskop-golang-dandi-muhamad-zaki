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

type ScreeningRepository interface{
	GetByCinema(q dto.ScreeningQuery) ([]entity.Screening, int, error)
	GetByID(id int) (*entity.Screening, error)
	Create(s dto.ScreeningRequest) error
	Update(id int, data dto.UpdateScreeningRequest) error
	Delete(id int) error
}

type screeningRepository struct {
	db     database.PgxIface
	Logger *zap.Logger
}

func NewScreeningRepository(db database.PgxIface, log *zap.Logger) ScreeningRepository {
	return &screeningRepository{
		db:     db,
		Logger: log,
	}
}

func (r *screeningRepository) GetByCinema(q dto.ScreeningQuery) ([]entity.Screening, int, error) {
	dateStr := q.Date
	if dateStr == "" {
		dateStr = time.Now().Format("02-01-2006")
	}

	selectedDate, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		r.Logger.Error("Invalid date format: ", zap.Error(err))
		return nil, 0, err
	}
	
	// Get total data for pagination
	var total int
	countQuery := `
	SELECT COUNT(DISTINCT m.id)
	FROM movies m
	LEFT JOIN screenings s ON s.movie_id = m.id
	LEFT JOIN studios st ON s.studio_id = st.id
	WHERE m.deleted_at IS NULL AND st.cinema_id = $1
		AND s.start_time >= $2
		AND s.start_time < $3
		AND s.deleted_at IS NULL
	GROUP BY s.movie_id, m.title, m.poster_url,
		m.duration_minute,
		m.rating_age, st.type, st.price, s.id
	`
	err = r.db.QueryRow(context.Background(), countQuery, q.CinemaID, selectedDate, selectedDate.AddDate(0,0,1)).Scan(&total)
	if err != nil {
		r.Logger.Error("Error query count screenings: ", zap.Error(err))
		return nil, 0, err
	}

	// Initiate rows
	var rows pgx.Rows

	// Conditional query based on page, limit, and all param
	query := `
	SELECT 
		s.id,
		s.studio_id,
		s.movie_id,
		s.start_time,
		s.start_time + (m.duration_minute * INTERVAL '1 minute') AS end_time
	FROM screenings s
	JOIN movies m ON m.id = s.movie_id
	LEFT JOIN studios st ON s.studio_id = st.id
	LEFT JOIN cinemas c ON st.cinema_id = c.id
	WHERE s.start_time >= $1
		AND s.start_time < $2
		AND s.deleted_at IS NULL
		AND c.id = $3
	ORDER BY s.start_time ASC
	`
	rows, err = r.db.Query(context.Background(), query, selectedDate, selectedDate.AddDate(0,0,1), q.CinemaID)
	if err != nil {
		r.Logger.Error("Error query get all screenings: ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var screenings []entity.Screening
	for rows.Next() {
    var s entity.Screening
		var startTime time.Time
		var endTime time.Time

    rows.Scan(&s.ID, &s.StudioID, &s.MovieID, &startTime, &endTime)

		// Convert to WIB
		loc, _ := time.LoadLocation("Asia/Jakarta")
		s.StartTime = startTime.In(loc)
		s.EndTime = endTime.In(loc)

    screenings = append(screenings, s)
	}

	return screenings, total, nil
}

func (r *screeningRepository) GetByID(id int) (*entity.Screening, error) {
	var s entity.Screening
	var startTime time.Time
	var endTime time.Time
	query := `SELECT 
		s.id,
		s.studio_id,
		s.movie_id,
		s.start_time,
		s.start_time + (m.duration_minute * INTERVAL '1 minute') AS end_time
	FROM screenings s
	JOIN movies m ON m.id = s.movie_id
	LEFT JOIN studios st ON s.studio_id = st.id
	LEFT JOIN cinemas c ON st.cinema_id = c.id
	WHERE s.id = $1`

	err := r.db.QueryRow(context.Background(), query, id).Scan(&s.ID, &s.StudioID, &s.MovieID, &startTime, &endTime)

	// Convert to WIB
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s.StartTime = startTime.In(loc)
	s.EndTime = endTime.In(loc)

	if err == pgx.ErrNoRows {
		r.Logger.Error("Error not found screening: ", zap.Error(err))
		return nil, utils.ErrNotFound("screening")
	}

	if err != nil {
		r.Logger.Error("Error query get screening by id: ", zap.Error(err))
		return nil, err
	}

	return &s, nil
}

func (r *screeningRepository) Create(s dto.ScreeningRequest) error {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background()); 
	if err != nil {
		r.Logger.Error("Error start db transaction: ", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	var screening entity.Screening

	// Parse date
	dateLayout := "02-01-2006"
	startDate, err := time.Parse(dateLayout, s.StartDate)
	endDate, err := time.Parse(dateLayout, s.EndDate)
	if err != nil {
		r.Logger.Error("Error parsing datetime: ", zap.Error(err))
		return err
	}
	
	// Create screening
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		for _, h := range s.StartHours {
			timeStr := fmt.Sprintf("%v %s", d.Format(dateLayout), h)
			// Convert to utc in db
			loc, _ := time.LoadLocation("Asia/Jakarta")

			wibTime, err := time.ParseInLocation(
				"02-01-2006 15.04",
				timeStr,
				loc,
			)

			// Convert to UTC for storage
			utcTime := wibTime.UTC()

			if err != nil {
				r.Logger.Error("Error parsing datetime: ", zap.Error(err))
				return err
			}
			query := `
				INSERT INTO screenings (studio_id, movie_id, start_time, created_at, updated_at)
				VALUES ($1, $2, $3, NOW(), NOW())
				RETURNING id
			`
			err = tx.QueryRow(context.Background(), query,
				&s.StudioID, &s.MovieID, &utcTime).Scan(&screening.ID)
			if err != nil {
				r.Logger.Error("Error query create screening: ", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

func (r *screeningRepository) Update(id int, s dto.UpdateScreeningRequest) error {
	// Handle db transaction
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()

	// Convert string to datetime
	startTime, err := time.Parse("2006-01-02 15:04", s.StartTime)

	// Update screenings
	query := `
		UPDATE screenings
		SET studio_id = $1,
		movie_id = $2,
		start_time = $3,
		updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := tx.Exec(context.Background(), query,
		&s.StudioID, &s.MovieID, &startTime, id)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found screening: ", zap.Error(err))
		return utils.ErrNotFound("screening")
	}

	if err != nil {
		r.Logger.Error("Error query update screening: ", zap.Error(err))
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *screeningRepository) Delete(id int) error {
	query := `
		UPDATE screenings
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query, id)

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.Logger.Error("Error not found screening: ", zap.Error(err))
		return utils.ErrNotFound("screening")
	}

	if err != nil {
		r.Logger.Error("Error query delete screening: ", zap.Error(err))
		return err
	}

	return nil
}