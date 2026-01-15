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

type CinemaHandler struct {
	Usecase usecase.Usecase
	Logger  *zap.Logger
	Config  utils.Configuration
}

func NewCinemaHandler(usecase usecase.Usecase, log *zap.Logger, config utils.Configuration) CinemaHandler {
	return CinemaHandler{
		Usecase: usecase,
		Logger:  log,
		Config:  config,
	}
}

func (h *CinemaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Retrieve query
	q, err := utils.GetPaginationQuery(r, h.Logger, h.Config)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
		return
	}

	// Execute get cinemas
	result, pagination, err := h.Usecase.CinemaUsecase.GetAll(q)
	if err != nil {
		h.Logger.Error("Error handling get cinemas: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get cinemas failed", err.Error())
		return
	}

	utils.ResponseWithPagination(w, http.StatusOK, "get cinemas success", result, pagination)
}

func (h *CinemaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute get cinemas
	result, err := h.Usecase.CinemaUsecase.GetByID(id)
	if err != nil && err.Error() == utils.ErrNotFound("cinema").Error() {
		h.Logger.Error("Error cinema not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "cinema not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling get cinema by id: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get cinema failed", err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "get cinema success", result)
}

func (h *CinemaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CinemaRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto cinema request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create cinema
	result, err := h.Usecase.CinemaUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create cinema: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create cinema failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "create cinema success", result)
}

func (h *CinemaHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	var req dto.CinemaRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto cinema request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute update cinema
	err = h.Usecase.CinemaUsecase.Update(id, req)
	if err != nil && err.Error() == utils.ErrNotFound("cinema").Error() {
		h.Logger.Error("Error cinema not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "cinema not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling update cinema: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "update cinema failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "update cinema success", nil)
}

func (h *CinemaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute delete Cinema
	err = h.Usecase.CinemaUsecase.Delete(id)
	if err != nil && err.Error() == utils.ErrNotFound("cinema").Error() {
		h.Logger.Error("Error cinema not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "cinema not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling delete cinema: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "delete cinema failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "delete cinema success", nil)
}