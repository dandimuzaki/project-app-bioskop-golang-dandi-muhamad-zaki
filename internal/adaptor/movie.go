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

type MovieHandler struct {
	Usecase usecase.Usecase
	Logger  *zap.Logger
	Config  utils.Configuration
}

func NewMovieHandler(usecase usecase.Usecase, log *zap.Logger, config utils.Configuration) MovieHandler {
	return MovieHandler{
		Usecase: usecase,
		Logger:  log,
		Config:  config,
	}
}

func (h *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Retrieve query
	q, err := utils.GetPaginationQuery(r, h.Logger, h.Config)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
		return
	}

	genre := r.URL.Query().Get("genre")
	var query dto.MovieQuery
	query.Page = q.Page
	query.Limit = q.Limit
	query.All = q.All
	query.Genre = genre

	// Execute get movies
	result, pagination, err := h.Usecase.MovieUsecase.GetAll(query)
	if err != nil {
		h.Logger.Error("Error handling get movies: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get movies failed", err.Error())
		return
	}

	utils.ResponseWithPagination(w, http.StatusOK, "get movies success", result, pagination)
}

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute get movies
	result, err := h.Usecase.MovieUsecase.GetByID(id)
	if err != nil && err.Error() == utils.ErrNotFound("movie").Error() {
		h.Logger.Error("Error movie not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "movie not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling get movie by id: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get movie failed", err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "get movie success", result)
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.MovieRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto movie request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create movie
	result, err := h.Usecase.MovieUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create movie: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create movie failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "create movie success", result)
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	var req dto.MovieRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto movie request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute update movie
	err = h.Usecase.MovieUsecase.Update(id, req)
	if err != nil && err.Error() == utils.ErrNotFound("movie").Error() {
		h.Logger.Error("Error movie not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "movie not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling update movie: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "update movie failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "update movie success", nil)
}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute delete movie
	err = h.Usecase.MovieUsecase.Delete(id)
	if err != nil && err.Error() == utils.ErrNotFound("movie").Error() {
		h.Logger.Error("Error movie not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "movie not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling delete movie: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "delete movie failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "delete movie success", nil)
}