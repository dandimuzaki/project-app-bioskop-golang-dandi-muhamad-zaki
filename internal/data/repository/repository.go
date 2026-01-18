package repository

import (
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type Repository struct {
	UserRepo UserRepository
	SessionRepo SessionRepository
	CinemaRepo CinemaRepository
	StudioRepo StudioRepository
	GenreRepo GenreRepository
	MovieRepo MovieRepository
	ScreeningRepo ScreeningRepository
	SeatRepo SeatRepository
	BookingRepo BookingRepository
	PaymentRepo PaymentRepository
	OTPRepo OTPRepository
	TicketRepo TicketRepository
}

func NewRepository(db database.PgxIface, log *zap.Logger) Repository {
	return Repository{
		UserRepo: NewUserRepository(db, log),
		SessionRepo: NewSessionRepository(db, log),
		CinemaRepo: NewCinemaRepository(db, log),
		StudioRepo: NewStudioRepository(db, log),
		GenreRepo: NewGenreRepository(db, log),
		MovieRepo: NewMovieRepository(db, log),
		ScreeningRepo: NewScreeningRepository(db, log),
		SeatRepo: NewSeatRepository(db, log),
		BookingRepo: NewBookingRepository(db, log),
		PaymentRepo: NewPaymentRepository(db, log),
		OTPRepo: NewOTPRepository(db, log),
		TicketRepo: NewTicketRepository(db, log),
	}
}