package repository

import (
	"context"

	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type SeatRepository interface{
	GetAvailableSeats(screeningID int) ([]dto.SeatResponse, error)
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

func (r *seatRepository) GetAvailableSeats(screeningID int) ([]dto.SeatResponse, error) {
	query := `SELECT s.id, s.seat_code
	FROM seats s
	JOIN studios st ON st.id = s.studio_id
	JOIN screenings sc ON sc.studio_id = st.id
	WHERE sc.id = $1
	AND NOT EXISTS (
			SELECT 1
			FROM booking_seats bs
			JOIN bookings b ON b.id = bs.booking_id
			WHERE bs.seat_id = s.id
				AND bs.screening_id = sc.id
				AND b.status IN ('pending', 'paid')
				AND b.expired_at < NOW()
	);
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
		err := rows.Scan(&s.ID, &s.SeatCode)
		if err != nil {
			r.Logger.Error("Error query get available seats: ", zap.Error(err))
			return nil, err
		}
		seats = append(seats, s)
	}

	return seats, nil
}