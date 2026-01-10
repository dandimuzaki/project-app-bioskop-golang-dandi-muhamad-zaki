package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type CinemaUsecase interface{
	GetAll(q dto.PaginationQuery) ([]entity.Cinema, *dto.Pagination, error)
	GetByID(id int) (*entity.Cinema, error)
	Create(data dto.CinemaRequest) (*entity.Cinema, error)
	Update(id int, data dto.CinemaRequest) error
	Delete(id int) error
}

type cinemaUsecase struct {
	Repo   *repository.Repository
	Logger *zap.Logger
}

func NewCinemaUsecase(repo *repository.Repository, log *zap.Logger) CinemaUsecase {
	return &cinemaUsecase{
		Repo:   repo,
		Logger: log,
	}
}

func (s *cinemaUsecase) GetAll(q dto.PaginationQuery) ([]entity.Cinema, *dto.Pagination, error) {
	// Execute repo to get all cinemas
	cinemas, total, err := s.Repo.CinemaRepo.GetAll(q)
	if err != nil {
		s.Logger.Error("Error get all cinemas Usecase: ", zap.Error(err))
		return nil, nil, err
	}

	// Calculate total pages
	var totalPages int
	totalPages = utils.TotalPage(q.Limit, total)

	// Create pagination
	var pagination dto.Pagination

	if q.All {
		pagination = dto.Pagination{
			TotalRecords: total,
		}
	} else {
		pagination = dto.Pagination{
			CurrentPage:  &q.Page,
			Limit:        &q.Limit,
			TotalPages:   &totalPages,
			TotalRecords: total,
		}
	}

	return cinemas, &pagination, nil
}

func (s *cinemaUsecase) GetByID(id int) (*entity.Cinema, error) {
	cinema, err := s.Repo.CinemaRepo.GetByID(id)
	if err != nil {
		s.Logger.Error("Error get cinema by id Usecase: ", zap.Error(err))
		return nil, err
	}
	return cinema, err
}

func (s *cinemaUsecase) Create(data dto.CinemaRequest) (*entity.Cinema, error) {
	cinema := entity.Cinema{
		Name: data.Name,
	}
	newcinema, err := s.Repo.CinemaRepo.Create(cinema)
	if err != nil {
		s.Logger.Error("Error create cinema Usecase: ", zap.Error(err))
		return nil, err
	}
	return newcinema, err
}

func (s *cinemaUsecase) Update(id int, data dto.CinemaRequest) error {
	cinema := entity.Cinema{
		Name: data.Name,
	}
	err := s.Repo.CinemaRepo.Update(id, &cinema)
	if err != nil {
		s.Logger.Error("Error update cinema Usecase: ", zap.Error(err))
		return err
	}
	return nil
}

func (s *cinemaUsecase) Delete(id int) error {
	err := s.Repo.CinemaRepo.Delete(id)
	if err != nil {
		s.Logger.Error("Error delete cinema Usecase: ", zap.Error(err))
		return err
	}
	return nil
}