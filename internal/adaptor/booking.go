package adaptor

import (
	"encoding/json"
	"net/http"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type BookingHandler struct {
	Usecase usecase.Usecase
	Logger *zap.Logger
	Config utils.Configuration
}

func NewBookingHandler(uc usecase.Usecase, log *zap.Logger, config utils.Configuration) BookingHandler {
	return BookingHandler{
		Usecase: uc,
		Logger: log,
		Config: config,
	}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.BookingRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto booking request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Retrieve user info
	user := r.Context().Value("user").(entity.User)
	req.UserID = user.ID

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create booking
	result, err := h.Usecase.BookingUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create booking: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create booking failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "create booking success", result)
}