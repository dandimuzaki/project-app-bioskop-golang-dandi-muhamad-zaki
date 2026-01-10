package utils

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/project-app-bioskop-golang/internal/dto"
	"go.uber.org/zap"
)

func TotalPage(limit int, totalData int) int {
	if totalData <= 0 {
		return 0
	}

	flimit := float64(limit)
	fdata := float64(totalData)

	res := math.Ceil(fdata / flimit)

	return int(res)
}

func ErrNotFound(field string) error {
	return fmt.Errorf("%s not found", field)
}

func StrToBool (str string) bool {
	if str == "true" {
		return true
	}
	return false
}

// Get query string for pagination
func GetPaginationQuery(w http.ResponseWriter, r *http.Request, log *zap.Logger, config Configuration) dto.PaginationQuery {
	// Retrieve query
	var page int = 1
	var err error
	var limit int
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Error("Error convert string to int: ", zap.Error(err))
			ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
			return dto.PaginationQuery{}
		}
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Error("Error convert string to int: ", zap.Error(err))
			ResponseFailed(w, http.StatusBadRequest, "invalid query param", err.Error())
			return dto.PaginationQuery{}
		}
	} else {
		limit = config.Limit
	}
	var all bool
	allStr := r.URL.Query().Get("all")
	if allStr != "" {
		all = StrToBool(allStr)
	}
	
	return dto.PaginationQuery{
		Page: page,
		Limit: limit,
		All: all,
	}
}

func GetIDParam(w http.ResponseWriter, r *http.Request, log *zap.Logger) int {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("Error convert string to int: ", zap.Error(err))
		ResponseFailed(w, http.StatusBadRequest, "error data", err.Error())
		return 0
	}
	return id
}