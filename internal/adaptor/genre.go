package adaptor

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type GenreHandler struct {
	Usecase usecase.Usecase
	Logger  *zap.Logger
	Config  utils.Configuration
}

func NewGenreHandler(usecase usecase.Usecase, log *zap.Logger, config utils.Configuration) GenreHandler {
	return GenreHandler{
		Usecase: usecase,
		Logger:  log,
		Config:  config,
	}
}

func (h *GenreHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Retrieve query
	q, err := utils.GetPaginationQuery(r, h.Logger, h.Config)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
		return
	}

	// Execute get genres
	result, pagination, err := h.Usecase.GenreUsecase.GetAll(q)
	if err != nil {
		h.Logger.Error("Error handling get genres: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get genres failed", err.Error())
		return
	}

	utils.ResponseWithPagination(w, http.StatusOK, "get genres success", result, pagination)
}

func (h *GenreHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute get genres
	result, err := h.Usecase.GenreUsecase.GetByID(id)
	if err != nil && err.Error() == utils.ErrNotFound("genre").Error() {
		h.Logger.Error("Error genre not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "genre not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling get genre by id: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "get genre failed", err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "get genre success", result)
}

func (h *GenreHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.Genre

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Error decode request body to dto genre request: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Execute create genre
	result, err := h.Usecase.GenreUsecase.Create(req)
	if err != nil {
		h.Logger.Error("Error handling create genre: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "create genre failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "create genre success", result)
}

func (h *GenreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Retrieve id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("Error convert string to int: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return
	}

	// Execute delete genre
	err = h.Usecase.GenreUsecase.Delete(id)
	if err != nil && err.Error() == utils.ErrNotFound("genre").Error() {
		h.Logger.Error("Error genre not found: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusNotFound, "genre not found", err.Error())
		return
	}
	if err != nil {
		h.Logger.Error("Error handling delete genre: ", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "delete genre failed", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "delete genre success", nil)
}