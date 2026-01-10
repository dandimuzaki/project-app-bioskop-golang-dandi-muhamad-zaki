package wire

import (
	"sync"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/project-app-bioskop-golang/internal/adaptor"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	mCustom "github.com/project-app-bioskop-golang/internal/middleware"
	"github.com/project-app-bioskop-golang/internal/usecase"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type App struct {
	Route *chi.Mux
	Stop chan struct{}
	WG *sync.WaitGroup
}

func Wiring(repo *repository.Repository, log *zap.Logger, config utils.Configuration) *App {
	r := chi.NewRouter()

	usecase := usecase.NewUsecase(repo, log)
	handler := adaptor.NewHandler(usecase, log, config)
	mw := mCustom.NewMiddlewareCustom(usecase, log)
	r.Mount("/api/v1", ApiV1(&handler, mw))

	// emailJobs := make(chan utils.EmailJob, 10) // BUFFER
	stop := make(chan struct{})
	// metrics := &utils.Metrics{}
	wg := &sync.WaitGroup{}

	// utils.StartEmailWorkers(3, emailJobs, stop, metrics, wg)

	return &App{
		Route: r,
		Stop: stop,
		WG: wg,
	}
}

func ApiV1(handler *adaptor.Handler, mw mCustom.MiddlewareCustom) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mw.Logging)

	r.Route("/", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Use(mw.RequirePermission("super admin", "admin", "staff"))
			
		})

		// Authentication
		r.Post("/login", handler.AuthHandler.Login)
		r.Post("/register", handler.AuthHandler.Register)
		r.Post("/logout", handler.AuthHandler.Logout)
	})

	r.Route("/cinemas", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Use(mw.RequirePermission("admin"))
			// CRUD cinemas
			r.Post("/", handler.CinemaHandler.Create)
			r.Put("/id", handler.CinemaHandler.Update)
			r.Delete("/id", handler.CinemaHandler.Delete)
		})

		r.Get("/", handler.CinemaHandler.GetAll)
		r.Get("/id", handler.CinemaHandler.GetByID)
	})
	
	return r
}