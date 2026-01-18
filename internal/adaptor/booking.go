package adaptor

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *BookingHandler) GetBookingHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Retrieve query
	q, err := utils.GetPaginationQuery(r, h.Logger, h.Config)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
		return
	}

	// Execute get booking history
	result, pagination, err := h.Usecase.BookingUsecase.GetBookingHistory(ctx, q)
	if err != nil {
		h.Logger.Error("Error handling get booking history: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get booking history failed", err.Error())
		return
	}

	utils.ResponseWithPagination(w, http.StatusOK, "get booking history success", result, pagination)
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute get booking
	result, err := h.Usecase.BookingUsecase.GetByID(id)
	if err != nil && err.Error() == utils.ErrNotFound("booking").Error() {
		h.Logger.Error("Error booking not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "booking not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling get booking by id: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get booking failed", err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "get booking success", result)
}