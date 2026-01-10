package adaptor

import (
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type Handler struct{
	AuthHandler AuthHandler
	CinemaHandler CinemaHandler
}

func NewHandler(uc usecase.Usecase, log *zap.Logger, config utils.Configuration) Handler {
	return Handler{
		AuthHandler: NewAuthHandler(uc, log, config),
		CinemaHandler: NewCinemaHandler(uc, log, config),
	}
}