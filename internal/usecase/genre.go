package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type GenreUsecase interface{
	GetAll(q dto.PaginationQuery) ([]entity.Genre, *dto.Pagination, error)
	GetByID(id int) (*entity.Genre, error)
	Create(data entity.Genre) (*entity.Genre, error)
	Delete(id int) error
}

type genreUsecase struct {
	Repo   *repository.Repository
	Logger *zap.Logger
}

func NewGenreUsecase(repo *repository.Repository, log *zap.Logger) GenreUsecase {
	return &genreUsecase{
		Repo:   repo,
		Logger: log,
	}
}

func (s *genreUsecase) GetAll(q dto.PaginationQuery) ([]entity.Genre, *dto.Pagination, error) {
	// Execute repo to get all genres
	genres, total, err := s.Repo.GenreRepo.GetAll(q)
	if err != nil {
		s.Logger.Error("Error get all genres usecase: ", zap.Error(err))
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

	return genres, &pagination, nil
}

func (s *genreUsecase) GetByID(id int) (*entity.Genre, error) {
	genre, err := s.Repo.GenreRepo.GetByID(id)
	if err != nil {
		s.Logger.Error("Error get genre by id usecase: ", zap.Error(err))
		return nil, err
	}
	return genre, err
}

func (s *genreUsecase) Create(data entity.Genre) (*entity.Genre, error) {
	newgenre, err := s.Repo.GenreRepo.Create(data)
	if err != nil {
		s.Logger.Error("Error create genre usecase: ", zap.Error(err))
		return nil, err
	}
	return newgenre, err
}

func (s *genreUsecase) Delete(id int) error {
	err := s.Repo.GenreRepo.Delete(id)
	if err != nil {
		s.Logger.Error("Error delete genre usecase: ", zap.Error(err))
		return err
	}
	return nil
}