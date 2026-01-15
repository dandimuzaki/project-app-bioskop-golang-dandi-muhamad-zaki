package adaptor

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type ScreeningHandler struct {
	Usecase usecase.Usecase
	Logger  *zap.Logger
	Config  utils.Configuration
}

func NewScreeningHandler(usecase usecase.Usecase, log *zap.Logger, config utils.Configuration) ScreeningHandler {
	return ScreeningHandler{
		Usecase: usecase,
		Logger:  log,
		Config:  config,
	}
}

func (h *ScreeningHandler) GetByCinema(w http.ResponseWriter, r *http.Request) {
	// Retrieve query
	q, err := utils.GetPaginationQuery(r, h.Logger, h.Config)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
		return
	}

	var query dto.ScreeningQuery
	query.Page = q.Page
	query.Limit = q.Limit
	query.All = q.All
	query.Date = r.URL.Query().Get("date")

	cinemaIDStr := r.URL.Query().Get("cinemaId")
	if cinemaIDStr != "" {
		query.CinemaID, err = strconv.Atoi(cinemaIDStr)
		if err != nil {
			utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
			return
		}
	}

	movieIDStr := r.URL.Query().Get("movieId")
	if movieIDStr != "" {
		query.MovieID, err = strconv.Atoi(movieIDStr)
		if err != nil {
			utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
			return
		}
	}

	// Execute get screenings
	result, pagination, err := h.Usecase.ScreeningUsecase.GetByCinema(query)
	if err != nil {
		h.Logger.Error("Error handling get screenings: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get screenings failed", err.Error())
		return
	}

	utils.ResponseWithPagination(w, http.StatusOK, "get screenings success", result, pagination)
}

func (h *ScreeningHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute get screenings
	result, err := h.Usecase.ScreeningUsecase.GetByID(id)
	if err != nil && err.Error() == utils.ErrNotFound("screening").Error() {
		h.Logger.Error("Error screening not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "screening not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling get screening by id: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get screening failed", err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "get screening success", result)
}

func (h *ScreeningHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ScreeningRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto screening request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create screening
	err = h.Usecase.ScreeningUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create screening: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create screening failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "create screening success", nil)
}

func (h *ScreeningHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	var req dto.UpdateScreeningRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto screening request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute update screening
	err = h.Usecase.ScreeningUsecase.Update(id, req)
	if err != nil && err.Error() == utils.ErrNotFound("screening").Error() {
		h.Logger.Error("Error screening not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "screening not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling update screening: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "update screening failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "update screening success", nil)
}

func (h *ScreeningHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute delete screening
	err = h.Usecase.ScreeningUsecase.Delete(id)
	if err != nil && err.Error() == utils.ErrNotFound("screening").Error() {
		h.Logger.Error("Error screening not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "screening not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling delete screening: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "delete screening failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "delete screening success", nil)
}