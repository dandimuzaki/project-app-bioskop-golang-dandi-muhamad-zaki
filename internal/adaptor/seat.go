package adaptor

import (
	"net/http"
	"strconv"

	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type SeatHandler struct {
	Usecase usecase.Usecase
	Logger *zap.Logger
	Config utils.Configuration
}

func NewSeatHandler(uc usecase.Usecase, log *zap.Logger, config utils.Configuration) SeatHandler {
	return SeatHandler{
		Usecase: uc,
		Logger: log,
		Config: config,
	}
}

func (h *SeatHandler) GetAvailableSeat(w http.ResponseWriter, r *http.Request) {
	// Retrieve screening id
	screeningIDStr := r.URL.Query().Get("screeningId")
	var screeningID int
	var err error
	if screeningIDStr != "" {
		screeningID, err = strconv.Atoi(screeningIDStr)
		if err != nil {
			h.Logger.Error("Error retrieve screening ID query param: ", zap.Error(err))
			utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
			return
		}
	} 

	// Execute get available seats
	result, err := h.Usecase.SeatUsecase.GetSeats(screeningID)
	if err != nil {
		h.Logger.Error("Error handling get available seats: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get available seats failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "get available seats success", result)
}