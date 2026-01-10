package adaptor

import (
	"encoding/json"
	"net/http"

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
	q := utils.GetPaginationQuery(w, r, h.Logger, h.Config)

	// Execute get Cinemas
	result, pagination, err := h.Usecase.CinemaUsecase.GetAll(q)
	if err != nil {
		h.Logger.Error("Error handling get Cinemas: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get Cinemas failed", err.Error())
		return
	}

	utils.ResponseWithPagination(w, http.StatusOK, "get Cinemas success", result, pagination)
}

func (h *CinemaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	id := utils.GetIDParam(w, r, h.Logger)

	// Execute get Cinemas
	result, err := h.Usecase.CinemaUsecase.GetByID(id)
	if err != nil && err.Error() == utils.ErrNotFound("Cinema").Error() {
		h.Logger.Error("Error Cinema not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "Cinema not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling get Cinema by id: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get Cinema failed", err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "get Cinema success", result)
}

func (h *CinemaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CinemaRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto Cinema request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create Cinema
	result, err := h.Usecase.CinemaUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create Cinema: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create Cinema failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "create Cinema success", result)
}

func (h *CinemaHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	id := utils.GetIDParam(w, r, h.Logger)

	var req dto.CinemaRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto Cinema request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute update Cinema
	err = h.Usecase.CinemaUsecase.Update(id, req)
	if err != nil && err.Error() == utils.ErrNotFound("Cinema").Error() {
		h.Logger.Error("Error Cinema not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "Cinema not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling update Cinema: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "update Cinema failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "update Cinema success", nil)
}

func (h *CinemaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	id := utils.GetIDParam(w, r, h.Logger)

	// Execute delete Cinema
	err := h.Usecase.CinemaUsecase.Delete(id)
	if err != nil && err.Error() == utils.ErrNotFound("Cinema").Error() {
		h.Logger.Error("Error Cinema not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "Cinema not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling delete Cinema: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "delete Cinema failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "delete Cinema success", nil)
}