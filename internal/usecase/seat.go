package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"go.uber.org/zap"
)

type SeatUsecase interface {
	GetAvailableSeats(screeningID int) ([]dto.SeatResponse, error)
}

type seatUsecase struct {
	repo *repository.Repository
	Logger *zap.Logger
}

func NewSeatUsecase(repo *repository.Repository, log *zap.Logger) SeatUsecase {
	return &seatUsecase{
		repo: repo,
		Logger: log,
	}
}

func (u *seatUsecase) GetAvailableSeats(screeningID int) ([]dto.SeatResponse, error) {
	seats, err := u.repo.SeatRepo.GetAvailableSeats(screeningID)
	if err != nil {
		u.Logger.Error("Error get available seats usecase: ", zap.Error(err))
		return nil, err
	}
	return seats, nil
}