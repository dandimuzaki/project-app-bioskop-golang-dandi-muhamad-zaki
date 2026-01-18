package usecase

import (
	"context"

	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type BookingUsecase interface {
	Create(b dto.BookingRequest) (*dto.BookingResponse, error)
	GetBookingHistory(ctx context.Context, q dto.PaginationQuery) ([]dto.BookingHistory, *dto.Pagination, error)
	GetByID(id int) (*dto.BookingResponse, error)
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

	response, err := u.GetByID(booking.ID)
	if err != nil {
		u.Logger.Error("Error get booking by id usecase: ", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *bookingUsecase) GetByID(id int) (*dto.BookingResponse, error) {
	b, err := u.repo.BookingRepo.GetByID(id)
	if err != nil {
		u.Logger.Error("Error get booking by id usecase: ", zap.Error(err))
		return nil, err
	}

	screening, err := u.repo.ScreeningRepo.GetByID(b.ScreeningID)
	movie, err := u.repo.MovieRepo.GetByID(screening.MovieID)
	studio, err := u.repo.StudioRepo.GetByID(screening.StudioID)
	cinema, err := u.repo.CinemaRepo.GetByID(studio.CinemaID)
	seats, err := u.repo.SeatRepo.GetSeatsByBookingID(b.ID)
	
	response := dto.BookingResponse{
		BookingID: b.ID,
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
		Status: b.Status,
		TotalAmount: studio.Price * float64(len(seats)),
		ExpiredAt: b.ExpiredAt,
	}
	return &response, err
}

func (u *bookingUsecase) GetBookingHistory(ctx context.Context, q dto.PaginationQuery) ([]dto.BookingHistory, *dto.Pagination, error) {
	bookings, total, err := u.repo.BookingRepo.GetBookingHistory(ctx, q)
	if err != nil {
		u.Logger.Error("Error get booking history usecase: ", zap.Error(err))
		return nil, nil, err
	}

	// Calculate total pages
	var totalPages int
	totalPages = utils.TotalPage(q.Limit, total)

	// Create pagination
	var pagination dto.Pagination

	if q.All {
		pagination = dto.Pagination{
			TotalRecords: total,
		}
	} else {
		pagination = dto.Pagination{
			CurrentPage:  &q.Page,
			Limit:        &q.Limit,
			TotalPages:   &totalPages,
			TotalRecords: total,
		}
	}
	
	return bookings, &pagination, nil
}