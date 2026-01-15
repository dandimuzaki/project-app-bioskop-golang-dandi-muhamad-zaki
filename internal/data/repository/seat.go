package repository

import (
	"context"
	"errors"

	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type SeatRepository interface{
	GetSeats(screeningID int) ([]dto.SeatResponse, error)
}

type seatRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewSeatRepository(db database.PgxIface, log *zap.Logger) SeatRepository {
	return &seatRepository{
		db: db,
		Logger: log,
	}
}

func (r *seatRepository) GetSeats(screeningID int) ([]dto.SeatResponse, error) {
	query := `SELECT s.id, s.seat_code,
	CASE
	WHEN NOT EXISTS (
			SELECT 1
			FROM booking_seats bs
			JOIN bookings b ON b.id = bs.booking_id
			WHERE bs.seat_id = s.id
				AND bs.screening_id = sc.id
				AND b.status IN ('pending', 'paid')
				AND b.expired_at > NOW()
	) THEN 'available'
	ELSE 'booked'
	END AS status
	FROM seats s
	JOIN studios st ON st.id = s.studio_id
	JOIN screenings sc ON sc.studio_id = st.id
	JOIN movies m ON sc.movie_id = m.id
	WHERE sc.id = $1 AND NOW() < sc.start_time + (m.duration_minute * INTERVAL '1 minute')
	`

	rows, err := r.db.Query(context.Background(), query, screeningID)
	if err != nil {
		r.Logger.Error("Error query get available seats: ", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	
	var seats []dto.SeatResponse
	for rows.Next() {
		var s dto.SeatResponse
		err := rows.Scan(&s.ID, &s.SeatCode, &s.Status)
		if err != nil {
			r.Logger.Error("Error query get available seats: ", zap.Error(err))
			return nil, err
		}
		seats = append(seats, s)
	}

	if len(seats) == 0 {
		return nil, errors.New("seats for this screening is already closed")
	}

	return seats, nil
}