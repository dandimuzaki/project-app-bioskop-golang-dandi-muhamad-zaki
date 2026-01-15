package usecase

import (
	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type MovieUsecase interface{
	GetAll(q dto.MovieQuery) ([]dto.MovieResponse, *dto.Pagination, error)
	GetByID(id int) (*dto.MovieResponse, error)
	Create(data dto.MovieRequest) (*dto.MovieResponse, error)
	Update(id int, data dto.MovieRequest) error
	Delete(id int) error
}

type movieUsecase struct {
	Repo   *repository.Repository
	Logger *zap.Logger
}

func NewMovieUsecase(repo *repository.Repository, log *zap.Logger) MovieUsecase {
	return &movieUsecase{
		Repo:   repo,
		Logger: log,
	}
}

func (s *movieUsecase) GetAll(q dto.MovieQuery) ([]dto.MovieResponse, *dto.Pagination, error) {
	// Execute repo to get all movies
	movies, total, err := s.Repo.MovieRepo.GetAll(q)
	if err != nil {
		s.Logger.Error("Error get all movies usecase: ", zap.Error(err))
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

	var response []dto.MovieResponse
	for _, m := range movies {
		response = append(response, dto.MovieResponse{
			MovieID: m.ID,
			Title: m.Title,
			Synopsis: m.Synopsis,
			Genres: m.Genres,
			PosterURL: m.PosterURL,
			TrailerURL: m.TrailerURL,
			Duration: m.Duration,
			ReleaseDate: m.ReleaseDate,
			Language: m.Language,
			RatingAge: m.RatingAge,
		})
	}

	return response, &pagination, nil
}

func (s *movieUsecase) GetByID(id int) (*dto.MovieResponse, error) {
	m, err := s.Repo.MovieRepo.GetByID(id)
	if err != nil {
		s.Logger.Error("Error get movie by id usecase: ", zap.Error(err))
		return nil, err
	}
	return m, err
}

func (s *movieUsecase) Create(data dto.MovieRequest) (*dto.MovieResponse, error) {
	newMovieID, err := s.Repo.MovieRepo.Create(data)
	if err != nil {
		s.Logger.Error("Error create movie usecase: ", zap.Error(err))
		return nil, err
	}
	m, err := s.Repo.MovieRepo.GetByID(newMovieID)
	if err != nil {
		s.Logger.Error("Error create movie usecase: ", zap.Error(err))
		return nil, err
	}
	return m, err
}

func (s *movieUsecase) Update(id int, data dto.MovieRequest) error {
	err := s.Repo.MovieRepo.Update(id, data)
	if err != nil {
		s.Logger.Error("Error update movie usecase: ", zap.Error(err))
		return err
	}
	return nil
}

func (s *movieUsecase) Delete(id int) error {
	err := s.Repo.MovieRepo.Delete(id)
	if err != nil {
		s.Logger.Error("Error delete movie usecase: ", zap.Error(err))
		return err
	}
	return nil
}