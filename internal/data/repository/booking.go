package repository

import (
	"context"
	"errors"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type BookingRepository interface{
	Create(b dto.BookingRequest) (*entity.Booking, error)
}

type bookingRepository struct {
	db database.PgxIface
	Logger *zap.Logger
}

func NewBookingRepository(db database.PgxIface, log *zap.Logger) BookingRepository {
	return &bookingRepository{
		db: db,
		Logger: log,
	}
}

func (r *bookingRepository) Create(b dto.BookingRequest) (*entity.Booking, error) {
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

	// Validate screening
	var screeningID int
	query := `SELECT s.id
FROM screenings s
JOIN movies m ON m.id = s.movie_id
WHERE s.id = $1
  AND s.deleted_at IS NULL
  AND NOW() < s.start_time + (m.duration_minute * INTERVAL '1 minute')
`
	err = tx.QueryRow(context.Background(), query, b.ScreeningID).Scan(&screeningID)
	if err != nil {
		r.Logger.Error("Booking is closed for this screening: ", zap.Error(err))
		return nil, errors.New("booking is closed for this screening")
	}

	// Create booking
	var booking entity.Booking
	query = `INSERT INTO bookings (user_id, screening_id, status, expired_at, created_at, updated_at)
	VALUES ($1, $2, 'pending', NOW() + interval '10 minute', NOW(), NOW()) RETURNING id, user_id, screening_id, status, expired_at, created_at, updated_at`
	err = tx.QueryRow(context.Background(), query, b.UserID, b.ScreeningID).Scan(&booking.ID, &booking.UserID, &booking.ScreeningID, &booking.Status, &booking.ExpiredAt, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		r.Logger.Error("Error query create bookings: ", zap.Error(err))
		return nil, err
	}

	// Create booking_seats
	for _, seatID := range b.Seats {
		query := `INSERT INTO booking_seats (booking_id, screening_id, seat_id, created_at)
		VALUES ($1, $2, $3, NOW())`

		_, err := tx.Exec(context.Background(), query, booking.ID, b.ScreeningID, seatID)
		if err != nil {
			r.Logger.Error("Error query create booking_seats: ", zap.Error(err))
			return nil, errors.New("one or more seats already booked")
		}

		var seat entity.Seat
		query = `SELECT id, seat_code FROM seats WHERE id = $1`

		err = r.db.QueryRow(context.Background(), query, seatID).Scan(&seat.ID, &seat.SeatCode)
		if err != nil {
			r.Logger.Error("Error query get seat: ", zap.Error(err))
			return nil, err
		}

		booking.Seats = append(booking.Seats, seat)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &booking, nil
}