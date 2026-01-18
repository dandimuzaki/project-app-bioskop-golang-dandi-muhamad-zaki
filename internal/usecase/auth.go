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
	Register(data dto.RegisterRequest) (*dto.OTPResponse, error)
	VerifyOTP(dto.OTPRequest) (*dto.AuthResponse, error)
	ResendOTP(email string) (*string, error)
	ValidateToken(token string) (*int, error)
	Logout(token string) error
}

type authUsecase struct {
	Repo *repository.Repository
	Logger *zap.Logger
	emailJobs chan <- utils.EmailJob
	Config utils.Configuration
}

func NewAuthUsecase(repo *repository.Repository, log *zap.Logger, emailJobs chan <- utils.EmailJob, config utils.Configuration) AuthUsecase {
	return &authUsecase{
		Repo: repo,
		Logger: log,
		emailJobs: emailJobs,
		Config: config,
	}
}

func (u *authUsecase) Login(email, password string) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := u.Repo.UserRepo.FindByEmail(email)
	if err == utils.ErrNotFound("user") {
		u.Logger.Error("User not found: ", zap.Error(err))
		return nil, err
	}

	if err != nil {
		u.Logger.Error("Error find user by email: ", zap.Error(err))
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(password, *user.Password) {
		u.Logger.Error("Incorrect password: ", zap.Error(err))
		return nil, errors.New("incorrect password")
	}

	// Record session
	token, err := u.Repo.SessionRepo.Create(user.ID)
	if err != nil {
		u.Logger.Error("Error create token: ", zap.Error(err))
		return nil, errors.New("token error")
	}
	
	res := dto.AuthResponse{
		Name: user.Name,
		Email: user.Email,
		Token: token,
	}

	return &res, nil
}

func (u *authUsecase) Register(data dto.RegisterRequest) (*dto.OTPResponse, error) {
	// Check if email is registered
	user, err := u.Repo.UserRepo.FindByEmail(data.Email)
	if user != nil {
		return nil, errors.New("user already registered")
	}

	// Hash password
	passwordHashed := utils.HashPassword(data.Password)
	newUser := entity.User{
		Name: data.Name,
		Email: data.Email,
		Password: &passwordHashed,
		Role: "customer",
	}
	
	// Execute create user
	user, err = u.Repo.UserRepo.Create(&newUser)
	if err != nil {
		u.Logger.Error("Error create user usecase: ", zap.Error(err))
		return nil, err
	}

	// Generate OTP
	otpStr, err := utils.GenerateOTP(6)
	if err != nil {
		u.Logger.Error("Error generate OTP: ", zap.Error(err))
		return nil, err
	}

	// Hash OTP
	otpHash := utils.HashPassword(otpStr)

	otp := entity.OTP{
		Email: user.Email,
		OTPHash: otpHash,
	}

	// Insert into database
	err = u.Repo.OTPRepo.Create(otp)
	if err != nil {
		u.Logger.Error("Error insert OTP: ", zap.Error(err))
		return nil, err
	}

	res := dto.OTPResponse{
		Name: user.Name,
		Email: user.Email,
		OTP: otpStr,
	}

	// Format email content
	body := utils.SendOTP(res)
	to := user.Email
	subject := "Email Verification"
	content := dto.EmailRequest{
		To: to,
		Subject: subject,
		Body: body,
	}

	// Send OTP
	u.emailJobs <- utils.EmailJob{
		EmailContent: content,
		Config: u.Config,
		Log: u.Logger,
	}

	return &res, nil
}

func (u *authUsecase) VerifyOTP(data dto.OTPRequest) (*dto.AuthResponse, error) {
	// Find OTP by email
	otpHash, err := u.Repo.OTPRepo.Find(data.Email)
	if err != nil {
		u.Logger.Error("Error find OTP: ", zap.Error(err))
		return nil, err
	}

	// Compare OTP
	if !utils.CheckPassword(*data.OTP, *otpHash) {
		u.Logger.Error("Invalid OTP: ", zap.Error(err))
		return nil, errors.New("invalid OTP")
	}

	// Update OTP record
	err = u.Repo.OTPRepo.Update(*otpHash)
	if err != nil {
		u.Logger.Error("Invalid OTP: ", zap.Error(err))
		return nil, err
	}

	// Get user
	user, err := u.Repo.UserRepo.FindByEmail(data.Email)
	if err == utils.ErrNotFound("user") {
		u.Logger.Error("User not found: ", zap.Error(err))
		return nil, err
	}

	// Record session
	token, err := u.Repo.SessionRepo.Create(user.ID)
	if err != nil {
		u.Logger.Error("Error create token: ", zap.Error(err))
		return nil, errors.New("token error")
	}
	
	res := dto.AuthResponse{
		Name: user.Name,
		Email: user.Email,
		Token: token,
	}

	return &res, nil
}

func (u *authUsecase) ResendOTP(email string) (*string, error) {
	// Generate OTP
	otpStr, err := utils.GenerateOTP(6)
	if err != nil {
		u.Logger.Error("Error generate OTP: ", zap.Error(err))
		return nil, err
	}

	// Hash OTP
	otpHash := utils.HashPassword(otpStr)

	otp := entity.OTP{
		Email: email,
		OTPHash: otpHash,
	}

	// Insert into database
	err = u.Repo.OTPRepo.Create(otp)
	if err != nil {
		u.Logger.Error("Error insert OTP: ", zap.Error(err))
		return nil, err
	}

	user, err := u.Repo.UserRepo.FindByEmail(email)
	if err == utils.ErrNotFound("user") {
		u.Logger.Error("User not found: ", zap.Error(err))
		return nil, err
	}

	res := dto.OTPResponse{
		Name: user.Name,
		Email: user.Email,
		OTP: otpStr,
	}

	// Format email content
	body := utils.SendOTP(res)
	to := user.Email
	subject := "Email Verification"
	content := dto.EmailRequest{
		To: to,
		Subject: subject,
		Body: body,
	}

	// Send OTP
	u.emailJobs <- utils.EmailJob{
		EmailContent: content,
		Config: u.Config,
		Log: u.Logger,
	}

	return &otpStr, nil
}

func (u *authUsecase) ValidateToken(token string) (*int, error) {
	// Validate token to authorize user
	userID, err := u.Repo.SessionRepo.ValidateToken(token)
	if err != nil {
		u.Logger.Error("Error validate token usecase: ", zap.Error(err))
		return nil, err
	}

	return userID, nil
}

func (u *authUsecase) Logout(token string) error {
	err := u.Repo.SessionRepo.Revoke(token)
	if err != nil {
		u.Logger.Error("Error logout usecase: ", zap.Error(err))
		return err
	}
	return nil
}