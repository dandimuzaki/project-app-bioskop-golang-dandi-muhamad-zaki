package usecase

import (
	"errors"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type AuthUsecase interface {
	Login(email, password string) (*dto.AuthResponse, error)
	Register(data dto.RegisterRequest) (*dto.AuthResponse, error)
	ValidateToken(token string) (*int, error)
	Logout(token string) error
}

type authUsecase struct {
	Repo *repository.Repository
	Logger *zap.Logger
}

func NewAuthUsecase(repo *repository.Repository, log *zap.Logger) AuthUsecase {
	return &authUsecase{
		Repo: repo,
		Logger: log,
	}
}

func (s *authUsecase) Login(email, password string) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.Repo.UserRepo.FindByEmail(email)
	if err == utils.ErrNotFound("user") {
		s.Logger.Error("User not found: ", zap.Error(err))
		return nil, err
	}

	if err != nil {
		s.Logger.Error("Error find user by email: ", zap.Error(err))
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(password, *user.Password) {
		s.Logger.Error("Incorrect password: ", zap.Error(err))
		return nil, errors.New("incorrect password")
	}

	// Record session
	token, err := s.Repo.SessionRepo.Create(user.ID)
	if err != nil {
		s.Logger.Error("Error create token: ", zap.Error(err))
		return nil, errors.New("token error")
	}
	
	res := dto.AuthResponse{
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		Token: token,
	}

	return &res, nil
}

func (s *authUsecase) Register(data dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email is registered
	user, err := s.Repo.UserRepo.FindByEmail(data.Email)
	if user != nil {
		return nil, errors.New("user already registered")
	}

	// Hash password
	passwordHashed := utils.HashPassword(data.Password)
	newUser := entity.User{
		Name: data.Name,
		Email: data.Email,
		Password: &passwordHashed,
		Role: "staff",
	}
	
	// Execute create user
	user, err = s.Repo.UserRepo.Create(&newUser)
	if err != nil {
		s.Logger.Error("Error create user Usecase: ", zap.Error(err))
		return nil, err
	}

	// Record session
	token, err := s.Repo.SessionRepo.Create(user.ID)
	res := dto.AuthResponse{
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		Token: token,
	}

	return &res, nil
}

func (s *authUsecase) ValidateToken(token string) (*int, error) {
	// Validate token to authorize user
	userID, err := s.Repo.SessionRepo.ValidateToken(token)
	if err != nil {
		s.Logger.Error("Error validate token Usecase: ", zap.Error(err))
		return nil, err
	}

	return userID, nil
}

func (s *authUsecase) Logout(token string) error {
	err := s.Repo.SessionRepo.Revoke(token)
	if err != nil {
		s.Logger.Error("Error logout Usecase: ", zap.Error(err))
		return err
	}
	return nil
}