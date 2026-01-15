package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type StudioUsecase interface{
	GetAll(q dto.PaginationQuery) ([]entity.Studio, *dto.Pagination, error)
	GetByID(id int) (*entity.Studio, error)
	Create(data dto.StudioRequest) (*entity.Studio, error)
	Update(id int, data dto.StudioRequest) error
	Delete(id int) error
}

type studioUsecase struct {
	Repo   *repository.Repository
	Logger *zap.Logger
}

func NewStudioUsecase(repo *repository.Repository, log *zap.Logger) StudioUsecase {
	return &studioUsecase{
		Repo:   repo,
		Logger: log,
	}
}

func (s *studioUsecase) GetAll(q dto.PaginationQuery) ([]entity.Studio, *dto.Pagination, error) {
	// Execute repo to get all studios
	studios, total, err := s.Repo.StudioRepo.GetAll(q)
	if err != nil {
		s.Logger.Error("Error get all studios Usecase: ", zap.Error(err))
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

	return studios, &pagination, nil
}

func (s *studioUsecase) GetByID(id int) (*entity.Studio, error) {
	studio, err := s.Repo.StudioRepo.GetByID(id)
	if err != nil {
		s.Logger.Error("Error get studio by id Usecase: ", zap.Error(err))
		return nil, err
	}
	return studio, err
}

func (s *studioUsecase) Create(data dto.StudioRequest) (*entity.Studio, error) {
	newstudio, err := s.Repo.StudioRepo.Create(data)
	if err != nil {
		s.Logger.Error("Error create studio Usecase: ", zap.Error(err))
		return nil, err
	}
	return newstudio, err
}

func (s *studioUsecase) Update(id int, data dto.StudioRequest) error {
	err := s.Repo.StudioRepo.Update(id, data)
	if err != nil {
		s.Logger.Error("Error update studio Usecase: ", zap.Error(err))
		return err
	}
	return nil
}

func (s *studioUsecase) Delete(id int) error {
	err := s.Repo.StudioRepo.Delete(id)
	if err != nil {
		s.Logger.Error("Error delete studio Usecase: ", zap.Error(err))
		return err
	}
	return nil
}