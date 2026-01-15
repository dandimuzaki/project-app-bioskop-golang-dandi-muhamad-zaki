package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type CinemaUsecase interface{
	GetAll(q dto.PaginationQuery) ([]dto.CinemaResponse, *dto.Pagination, error)
	GetByID(id int) (*dto.CinemaResponse, error)
	Create(data dto.CinemaRequest) (*dto.CinemaResponse, error)
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

func (s *cinemaUsecase) GetAll(q dto.PaginationQuery) ([]dto.CinemaResponse, *dto.Pagination, error) {
	// Execute repo to get all cinemas
	cinemas, total, err := s.Repo.CinemaRepo.GetAll(q)
	if err != nil {
		s.Logger.Error("Error get all cinemas usecase: ", zap.Error(err))
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

	var response []dto.CinemaResponse
	for _, c := range cinemas {
		response = append(response, dto.CinemaResponse{
			CinemaID: c.ID,
			Name: c.Name,
			Location: c.Location,
		})
	}

	return response, &pagination, nil
}

func (s *cinemaUsecase) GetByID(id int) (*dto.CinemaResponse, error) {
	c, err := s.Repo.CinemaRepo.GetByID(id)
	if err != nil {
		s.Logger.Error("Error get cinema by id usecase: ", zap.Error(err))
		return nil, err
	}
	return c, err
}

func (s *cinemaUsecase) Create(data dto.CinemaRequest) (*dto.CinemaResponse, error) {
	cinema := entity.Cinema{
		Name: data.Name,
		Location: data.Location,
	}
	newCinema, err := s.Repo.CinemaRepo.Create(cinema)
	if err != nil {
		s.Logger.Error("Error create cinema usecase: ", zap.Error(err))
		return nil, err
	}
	response := dto.CinemaResponse{
		CinemaID: newCinema.ID,
		Name: newCinema.Name,
		Location: newCinema.Location,
	}
	return &response, err
}

func (s *cinemaUsecase) Update(id int, data dto.CinemaRequest) error {
	cinema := entity.Cinema{
		Name: data.Name,
		Location: data.Location,
	}
	err := s.Repo.CinemaRepo.Update(id, &cinema)
	if err != nil {
		s.Logger.Error("Error update cinema usecase: ", zap.Error(err))
		return err
	}
	return nil
}

func (s *cinemaUsecase) Delete(id int) error {
	err := s.Repo.CinemaRepo.Delete(id)
	if err != nil {
		s.Logger.Error("Error delete cinema usecase: ", zap.Error(err))
		return err
	}
	return nil
}