package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"go.uber.org/zap"
)

type UserUsecase interface {
	GetByID(id int) (entity.User, error)
}

type userUsecase struct {
	Repo *repository.Repository
	Logger *zap.Logger
}

func NewUserUsecase(repo *repository.Repository, log *zap.Logger) UserUsecase {
	return &userUsecase{
		Repo: repo,
		Logger: log,
	}
}

func (s *userUsecase) GetByID(id int) (entity.User, error) {
	user, err := s.Repo.UserRepo.GetByID(id)
	if err != nil {
		s.Logger.Error("Error get user by id Usecase: ", zap.Error(err))
		return user, err
	}
	return user, nil
}