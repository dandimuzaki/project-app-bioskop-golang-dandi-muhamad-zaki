package repository

import (
	"github.com/project-app-bioskop-golang/pkg/database"
	"go.uber.org/zap"
)

type Repository struct {
	UserRepo UserRepository
	SessionRepo SessionRepository
	CinemaRepo CinemaRepository
}

func NewRepository(db database.PgxIface, log *zap.Logger) Repository {
	return Repository{
		UserRepo: NewUserRepository(db, log),
		SessionRepo: NewSessionRepository(db, log),
		CinemaRepo: NewCinemaRepository(db, log),
	}
}