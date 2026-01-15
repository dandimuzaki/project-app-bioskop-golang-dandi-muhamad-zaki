package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/repository"
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
}

func NewUsecase(repo *repository.Repository, log *zap.Logger) Usecase {
	return Usecase{
		AuthUsecase: NewAuthUsecase(repo, log),
		UserUsecase: NewUserUsecase(repo, log),
		CinemaUsecase: NewCinemaUsecase(repo, log),
		StudioUsecase: NewStudioUsecase(repo, log),
		GenreUsecase: NewGenreUsecase(repo, log),
		MovieUsecase: NewMovieUsecase(repo, log),
		ScreeningUsecase: NewScreeningUsecase(repo, log),
		SeatUsecase: NewSeatUsecase(repo, log),
		BookingUsecase: NewBookingUsecase(repo, log),
	}
}