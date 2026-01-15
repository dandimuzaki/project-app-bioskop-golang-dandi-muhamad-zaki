package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type PaymentUsecase interface {
	Create(b dto.PaymentRequest) (*dto.PaymentResponse, error)
	Update(b dto.UpdatePayment) error
}

type paymentUsecase struct {
	repo *repository.Repository
	Logger *zap.Logger
}

func NewPaymentUsecase(repo *repository.Repository, log *zap.Logger) PaymentUsecase {
	return &paymentUsecase{
		repo: repo,
		Logger: log,
	}
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

func (u *paymentUsecase) Update(b dto.UpdatePayment) error {
	err := u.repo.PaymentRepo.Update(b)
	if err != nil {
		u.Logger.Error("Error update payment usecase: ", zap.Error(err))
		return err
	}
	
	return nil
}