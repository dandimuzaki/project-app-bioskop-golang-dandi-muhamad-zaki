package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type Usecase struct {
	AuthUsecase AuthUsecase
	UserUsecase UserUsecase
	CinemaUsecase CinemaUsecase
	StudioUsecase StudioUsecase
	GenreUsecase GenreUsecase
	MovieUsecase MovieUsecase
	ScreeningUsecase ScreeningUsecase
	SeatUsecase SeatUsecase
	BookingUsecase BookingUsecase
	PaymentUsecase PaymentUsecase
}

func NewUsecase(repo *repository.Repository, log *zap.Logger, emailJobs chan <- utils.EmailJob, ticketJobs chan <- utils.TicketJob, config utils.Configuration) Usecase {
	return Usecase{
		AuthUsecase: NewAuthUsecase(repo, log, emailJobs, config),
		UserUsecase: NewUserUsecase(repo, log),
		CinemaUsecase: NewCinemaUsecase(repo, log),
		StudioUsecase: NewStudioUsecase(repo, log),
		GenreUsecase: NewGenreUsecase(repo, log),
		MovieUsecase: NewMovieUsecase(repo, log),
		ScreeningUsecase: NewScreeningUsecase(repo, log),
		SeatUsecase: NewSeatUsecase(repo, log),
		BookingUsecase: NewBookingUsecase(repo, log),
		PaymentUsecase: NewPaymentUsecase(repo, log, ticketJobs, config),
	}
}