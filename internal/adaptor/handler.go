package adaptor

import (
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type Handler struct{
	AuthHandler AuthHandler
	CinemaHandler CinemaHandler
	StudioHandler StudioHandler
	GenreHandler GenreHandler
	MovieHandler MovieHandler
	ScreeningHandler ScreeningHandler
	SeatHandler SeatHandler
	BookingHandler BookingHandler
	PaymentHandler PaymentHandler
}

func NewHandler(uc usecase.Usecase, log *zap.Logger, config utils.Configuration) Handler {
	return Handler{
		AuthHandler: NewAuthHandler(uc, log, config),
		CinemaHandler: NewCinemaHandler(uc, log, config),
		StudioHandler: NewStudioHandler(uc, log, config),
		GenreHandler: NewGenreHandler(uc, log, config),
		MovieHandler: NewMovieHandler(uc, log, config),
		ScreeningHandler: NewScreeningHandler(uc, log, config),
		SeatHandler: NewSeatHandler(uc, log, config),
		BookingHandler: NewBookingHandler(uc, log, config),
		PaymentHandler: NewPaymentHandler(uc, log, config),
	}
}