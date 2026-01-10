package adaptor

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Usecase usecase.Usecase
	Logger *zap.Logger
	Config utils.Configuration
}

func NewAuthHandler(usecase usecase.Usecase, log *zap.Logger, config utils.Configuration) AuthHandler {
	return AuthHandler{
		Usecase: usecase,
		Logger: log,
		Config: config,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto login request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute login
	result, err := h.Usecase.AuthUsecase.Login(req.Email, req.Password)
	if err != nil {
		h.Logger.Error("Error handling login user: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusUnauthorized, "username or password failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "login success", result)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto register request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute register
	result, err := h.Usecase.AuthUsecase.Register(req)
	if err != nil {
		h.Logger.Error("Error handling register user: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusUnauthorized, "register user failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "register success", result)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	token := strings.TrimSpace(strings.Replace(auth, "Bearer", "", 1))

	// Execute logout
	err := h.Usecase.AuthUsecase.Logout(token)
	if err != nil {
		h.Logger.Error("Error handling logout user: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "logout failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "logout success", nil)
}