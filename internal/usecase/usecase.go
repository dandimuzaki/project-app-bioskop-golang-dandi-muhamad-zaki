package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"go.uber.org/zap"
)

type Usecase struct {
	AuthUsecase AuthUsecase
	UserUsecase UserUsecase
	CinemaUsecase CinemaUsecase
}

func NewUsecase(repo *repository.Repository, log *zap.Logger) Usecase {
	return Usecase{
		AuthUsecase: NewAuthUsecase(repo, log),
		UserUsecase: NewUserUsecase(repo, log),
		CinemaUsecase: NewCinemaUsecase(repo, log),
	}
}