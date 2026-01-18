package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type BookingRepository interface{
	Create(b dto.BookingRequest) (*entity.Booking, error)
	GetByID(id int) (*entity.Booking, error)
	GetBookingHistory(ctx context.Context, q dto.PaginationQuery) ([]dto.BookingHistory, int, error)
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
		query := `INSERT INTO booking_seats (booking_id, screening_id, seat_id, booking_status, created_at)
		VALUES ($1, $2, $3, 'pending', NOW())`

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

func (r *bookingRepository) GetByID(id int) (*entity.Booking, error) {
	var b entity.Booking
	query := `SELECT id, user_id, screening_id, status, expired_at, created_at, updated_at
	FROM bookings WHERE id = $1
	`
	err := r.db.QueryRow(context.Background(), query, id).Scan(&b.ID, &b.UserID, &b.ScreeningID, &b.Status, &b.ExpiredAt, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		r.Logger.Error("Error query get booking by id: ", zap.Error(err))
		return nil, err
	}

	return &b, nil
}

func (r *bookingRepository) GetBookingHistory(ctx context.Context, q dto.PaginationQuery) ([]dto.BookingHistory, int, error) {
	var offset int
	offset = (q.Page - 1) * q.Limit
	
	// Get total data for pagination
	user := ctx.Value("user").(entity.User)
	var total int
	countQuery := `
	SELECT id FROM bookings WHERE user_id = $1`
	err := r.db.QueryRow(context.Background(), countQuery, user.ID).Scan(&total)
	if err != nil {
		r.Logger.Error("Error query count booking history: ", zap.Error(err))
		return nil, 0, err
	}

	// Initiate rows
	var rows pgx.Rows
	
	// Conditional query based on page, limit, and all param
	query := `SELECT b.id, m.title, c.name, m.poster_url, 
	ARRAY_AGG(s.seat_code) AS seats, b.status, sc.start_time
	FROM bookings b
	LEFT JOIN screenings sc ON sc.id = b.screening_id
	LEFT JOIN studios st ON st.id = sc.studio_id
	LEFT JOIN cinemas c ON c.id = st.cinema_id
	LEFT JOIN movies m ON m.id = sc.movie_id
	LEFT JOIN booking_seats bs ON bs.booking_id = b.id
	LEFT JOIN seats s ON s.id = bs.seat_id
	WHERE user_id = $1
	GROUP BY b.id, m.title, c.name, m.poster_url, b.status, sc.start_time
	ORDER BY sc.start_time DESC
	`
	
	if !q.All && q.Limit > 0 {
		query += ` LIMIT $2 OFFSET $3`
		rows, err = r.db.Query(context.Background(), query, user.ID, q.Limit, offset)
	} else {
		rows, err = r.db.Query(context.Background(), query, user.ID)
	}

	if err != nil {
		r.Logger.Error("Error query get booking history: ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()
	
	var bookings []dto.BookingHistory
	for rows.Next() {
		var b dto.BookingHistory
		var startTime time.Time
		err = rows.Scan(&b.ID, &b.MovieTitle, &b.Cinema, &b.PosterURL, &b.Seats, &b.Status, &startTime)
		if err != nil {
			r.Logger.Error("Error query get booking by id: ", zap.Error(err))
			return nil, 0, err
		}
		// Convert to WIB
		loc, _ := time.LoadLocation("Asia/Jakarta")
		b.Date = startTime.In(loc).Format("02 January 2006")
		b.StartTime = startTime.In(loc).Format("15.04")
		bookings = append(bookings, b)
	}

	return bookings, total, nil
}