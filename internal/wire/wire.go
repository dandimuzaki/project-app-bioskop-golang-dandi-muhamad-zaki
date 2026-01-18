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

	emailJobs := make(chan utils.EmailJob, 10) // BUFFER
	ticketJobs := make(chan utils.TicketJob)
	stop := make(chan struct{})
	metrics := &utils.Metrics{}
	wg := &sync.WaitGroup{}

	utils.StartEmailWorkers(3, emailJobs, stop, metrics, wg)
	utils.StartTicketWorkers(3, ticketJobs, stop, metrics, wg)

	usecase := usecase.NewUsecase(repo, log, emailJobs, ticketJobs, config)
	handler := adaptor.NewHandler(usecase, log, config)
	mw := mCustom.NewMiddlewareCustom(usecase, log)
	r.Mount("/api/v1", ApiV1(&handler, mw))

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
			r.Use(mw.RequirePermission("admin"))
			r.Post("/studio-type", handler.StudioHandler.CreateStudioType)			
		})				
	})

	// Authentication
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.AuthHandler.Login)
		r.Post("/register", handler.AuthHandler.Register)
		r.Post("/logout", handler.AuthHandler.Logout)
		r.Post("/resend", handler.AuthHandler.ResendOTP)
		r.Post("/verify", handler.AuthHandler.VerifyOTP)
	})

	// Seats
	r.Route("/seats", func(r chi.Router) {
		r.Use(mw.AuthMiddleware())
		r.Get("/", handler.SeatHandler.GetSeatsByScreening)
	})

	r.Route("/cinemas", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Use(mw.RequirePermission("admin"))
			// CRUD cinemas
			r.Post("/", handler.CinemaHandler.Create)
			r.Put("/{id}", handler.CinemaHandler.Update)
			r.Delete("/{id}", handler.CinemaHandler.Delete)
		})

		r.Get("/", handler.CinemaHandler.GetAll)
		r.Get("/{id}", handler.CinemaHandler.GetByID)
	})

	r.Route("/studios", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Use(mw.RequirePermission("admin"))
			// CRUD studios
			r.Get("/", handler.StudioHandler.GetAll)
			r.Get("/{id}", handler.StudioHandler.GetByID)
			r.Post("/", handler.StudioHandler.Create)
			r.Put("/{id}", handler.StudioHandler.Update)
			r.Delete("/{id}", handler.StudioHandler.Delete)
		})
	})

	r.Route("/genres", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Use(mw.RequirePermission("admin"))
			// CRUD genres
			r.Get("/", handler.GenreHandler.GetAll)
			r.Get("/{id}", handler.GenreHandler.GetByID)
			r.Post("/", handler.GenreHandler.Create)
			r.Delete("/{id}", handler.GenreHandler.Delete)
		})
	})

	r.Route("/movies", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Use(mw.RequirePermission("admin"))
			// CRUD movies
			r.Post("/", handler.MovieHandler.Create)
			r.Put("/{id}", handler.MovieHandler.Update)
			r.Delete("/{id}", handler.MovieHandler.Delete)
		})

		r.Get("/", handler.MovieHandler.GetAll)
		r.Get("/{id}", handler.MovieHandler.GetByID)
	})

	r.Route("/screenings", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Use(mw.RequirePermission("admin"))
			// CRUD screenings
			r.Post("/", handler.ScreeningHandler.Create)
			r.Put("/{id}", handler.ScreeningHandler.Update)
			r.Delete("/{id}", handler.ScreeningHandler.Delete)
		})

		r.Get("/", handler.ScreeningHandler.GetByCinema)
		r.Get("/{id}", handler.ScreeningHandler.GetByID)
	})

	r.Route("/bookings", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Post("/", handler.BookingHandler.Create)
			r.Get("/", handler.BookingHandler.GetBookingHistory)
			r.Get("/{id}", handler.BookingHandler.GetByID)
		})
	})

	r.Route("/payments", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(mw.AuthMiddleware())
			r.Post("/", handler.PaymentHandler.Create)
			r.Post("/callback", handler.PaymentHandler.Callback)
			r.Get("/method", handler.PaymentHandler.GetPaymentMethod)
		})
	})
	
	return r
}