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

type StudioHandler struct {
	Usecase usecase.Usecase
	Logger  *zap.Logger
	Config  utils.Configuration
}

func NewStudioHandler(usecase usecase.Usecase, log *zap.Logger, config utils.Configuration) StudioHandler {
	return StudioHandler{
		Usecase: usecase,
		Logger:  log,
		Config:  config,
	}
}

func (h *StudioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Retrieve query
	q, err := utils.GetPaginationQuery(r, h.Logger, h.Config)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
		return
	}

	// Execute get studios
	result, pagination, err := h.Usecase.StudioUsecase.GetAll(q)
	if err != nil {
		h.Logger.Error("Error handling get studios: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get studios failed", err.Error())
		return
	}

	utils.ResponseWithPagination(w, http.StatusOK, "get studios success", result, pagination)
}

func (h *StudioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute get studios
	result, err := h.Usecase.StudioUsecase.GetByID(id)
	if err != nil && err.Error() == utils.ErrNotFound("studio").Error() {
		h.Logger.Error("Error studio not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "studio not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling get studio by id: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get studio failed", err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "get studio success", result)
}

func (h *StudioHandler) CreateStudioType(w http.ResponseWriter, r *http.Request) {
	var req dto.StudioType

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto studio type request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create studio type
	result, err := h.Usecase.StudioUsecase.CreateStudioType(req)
	if err != nil {
		h.Logger.Error("Error handling create studio: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create studio type failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "create studio type success", result)
}

func (h *StudioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.StudioRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto studio request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create studio
	result, err := h.Usecase.StudioUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create studio: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create studio failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "create studio success", result)
}

func (h *StudioHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	var req dto.StudioRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto studio request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute update studio
	err = h.Usecase.StudioUsecase.Update(id, req)
	if err != nil && err.Error() == utils.ErrNotFound("studio").Error() {
		h.Logger.Error("Error studio not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "studio not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling update studio: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "update studio failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "update studio success", nil)
}

func (h *StudioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute delete studio
	err = h.Usecase.StudioUsecase.Delete(id)
	if err != nil && err.Error() == utils.ErrNotFound("studio").Error() {
		h.Logger.Error("Error studio not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "studio not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling delete studio: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "delete studio failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "delete studio success", nil)
}