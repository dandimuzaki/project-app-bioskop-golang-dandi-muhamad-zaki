package usecase

import (
	"context"
	"fmt"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type PaymentUsecase interface {
	GetPaymentMethod() ([]entity.PaymentMethod, error)
	Create(b dto.PaymentRequest) (*dto.PaymentResponse, error)
	Update(ctx context.Context, b dto.UpdatePayment) error
	SendTicket(ctx context.Context, id int) error
}

type paymentUsecase struct {
	repo *repository.Repository
	Logger *zap.Logger
	ticketJobs chan <- utils.TicketJob
	config utils.Configuration
}

func NewPaymentUsecase(repo *repository.Repository, log *zap.Logger, ticketJobs chan <- utils.TicketJob, config utils.Configuration) PaymentUsecase {
	return &paymentUsecase{
		repo: repo,
		Logger: log,
		ticketJobs: ticketJobs,
		config: config,
	}
}

func (u *paymentUsecase) GetPaymentMethod() ([]entity.PaymentMethod, error) {
	pm, err := u.repo.PaymentRepo.GetPaymentMethod()
	if err != nil {
		u.Logger.Error("Error get payment method usecase: ", zap.Error(err))
		return nil, err
	}
	return pm, nil
}

func (u *paymentUsecase) Create(b dto.PaymentRequest) (*dto.PaymentResponse, error) {
	paymentID, err := u.repo.PaymentRepo.Create(b)
	if err != nil {
		u.Logger.Error("Error create payment usecase: ", zap.Error(err))
		return nil, err
	}

	transactionID, err := utils.GenerateTransactionID(10)
	if err != nil {
		u.Logger.Error("Error generate transaction ID: ", zap.Error(err))
		return nil, err
	}

	response := dto.PaymentResponse{
		PaymentID: *paymentID,
		TransactionID: &transactionID,
	}
	
	return &response, nil
}

func (u *paymentUsecase) Update(ctx context.Context, b dto.UpdatePayment) error {
	bookingID, err := u.repo.PaymentRepo.Update(b)
	if err != nil {
		u.Logger.Error("Error update payment usecase: ", zap.Error(err))
		return err
	}

	if b.Status == "success" {
		u.SendTicket(ctx, *bookingID)
		return nil
	}
	
	return nil
}

func (u *paymentUsecase) SendTicket(ctx context.Context, id int) error {
	booking, err := u.repo.BookingRepo.GetByID(id)
	if err != nil {
		return err
	}

	screening, err := u.repo.ScreeningRepo.GetByID(booking.ScreeningID)
	movie, err := u.repo.MovieRepo.GetByID(screening.MovieID)
	studio, err := u.repo.StudioRepo.GetByID(screening.StudioID)
	cinema, err := u.repo.CinemaRepo.GetByID(studio.CinemaID)
	tickets, err := u.repo.TicketRepo.GetByBookingID(id)
	user := ctx.Value("user").(entity.User)

	res := dto.TicketEmail{
		Profile: dto.ProfileResponse{
			Name: user.Name,
			Email: user.Email,
		},
		BookingID: booking.ID,
		BookingDate: screening.StartTime.Format("02-01-2006"),
		Movie: *movie,
		Cinema: *cinema,
		Studio: dto.StudioResponse{
			StudioID: studio.ID,
			Name: studio.Name,
			Type: studio.Type,
			Price: studio.Price,
		},
		Screening: dto.ScreeningResponse{
			ScreeningID: screening.ID,
			StartTime: screening.StartTime.Format("15.04"),
		},
		Tickets: tickets,
	}

	fmt.Println(tickets)

	// Send ticket
	u.ticketJobs <- utils.TicketJob{
		Config: u.config,
		Log: u.Logger,
		Data: res,
	}

	return nil
}