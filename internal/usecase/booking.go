package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"go.uber.org/zap"
)

type BookingUsecase interface {
	Create(b dto.BookingRequest) (*dto.BookingResponse, error)
}

type bookingUsecase struct {
	repo *repository.Repository
	Logger *zap.Logger
}

func NewBookingUsecase(repo *repository.Repository, log *zap.Logger) BookingUsecase {
	return &bookingUsecase{
		repo: repo,
		Logger: log,
	}
}

func (u *bookingUsecase) Create(b dto.BookingRequest) (*dto.BookingResponse, error) {
	booking, err := u.repo.BookingRepo.Create(b)
	if err != nil {
		u.Logger.Error("Error create booking usecase: ", zap.Error(err))
		return nil, err
	}

	screening, err := u.repo.ScreeningRepo.GetByID(booking.ScreeningID)
	movie, err := u.repo.MovieRepo.GetByID(screening.MovieID)
	studio, err := u.repo.StudioRepo.GetByID(screening.StudioID)
	cinema, err := u.repo.CinemaRepo.GetByID(studio.CinemaID)

	var seats []dto.SeatResponse
	for _, s := range booking.Seats {
		seats = append(seats, dto.SeatResponse{
			ID: s.ID,
			SeatCode: s.SeatCode,
		})
	}
	
	response := dto.BookingResponse{
		BookingID: booking.ID,
		BookingDate: screening.StartTime.Format("02-01-2006"),
		Movie: *movie,
		Studio: dto.StudioResponse{
			StudioID: studio.ID,
			Name: studio.Name,
			Type: studio.Type,
			Price: studio.Price,
		},
		Cinema: *cinema,
		Seats: seats,
		Status: booking.Status,
		TotalAmount: studio.Price * float64(len(b.Seats)),
		ExpiredAt: booking.ExpiredAt,
	}
	return &response, nil
}