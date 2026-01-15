package usecase

import (
	"time"

	"github.com/project-app-bioskop-golang/internal/data/repository"
	"github.com/project-app-bioskop-golang/internal/dto"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

type ScreeningUsecase interface{
	GetByCinema(q dto.ScreeningQuery) ([]dto.MovieScreening, *dto.Pagination, error)
	GetByID(id int) (*dto.MovieScreeningRow, error)
	Create(s dto.ScreeningRequest) error
	Update(id int, data dto.UpdateScreeningRequest) error
	Delete(id int) error
}

type screeningUsecase struct {
	Repo   *repository.Repository
	Logger *zap.Logger
}

func NewScreeningUsecase(repo *repository.Repository, log *zap.Logger) ScreeningUsecase {
	return &screeningUsecase{
		Repo:   repo,
		Logger: log,
	}
}

func (s *screeningUsecase) GetByCinema(q dto.ScreeningQuery) ([]dto.MovieScreening, *dto.Pagination, error) {
	// Execute repo to get all screenings
	screenings, total, err := s.Repo.ScreeningRepo.GetByCinema(q)
	if err != nil {
		s.Logger.Error("Error get all screenings Usecase: ", zap.Error(err))
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

	dateStr := q.Date
	if dateStr == "" {
		dateStr = time.Now().Format("02-01-2006")
	}

	result := map[int]*dto.MovieScreening{}
	for _, sc := range screenings {
    if _, ok := result[sc.MovieID]; !ok {
			m, _ := s.Repo.MovieRepo.GetByID(sc.MovieID)
			st, _ := s.Repo.StudioRepo.GetByID(sc.StudioID)
			result[sc.MovieID] = &dto.MovieScreening{
				Date: dateStr,
				Movie: *m,
				Studio: dto.StudioResponse{
					StudioID: st.ID,
					Name: st.Name,
				},
			}
    }

		startTime := sc.StartTime.Format("15.04")
		endTime := sc.EndTime.Format("15.04")

    result[sc.MovieID].Screenings = append(
			result[sc.MovieID].Screenings,
			dto.ScreeningResponse{
				ScreeningID: sc.ID,
				StartTime: startTime,
				EndTime: endTime,
			},
    )
	}

	var response []dto.MovieScreening
	for _, screening := range result {
		response = append(response, *screening)
	}

	return response, &pagination, nil
}

func (s *screeningUsecase) GetByID(id int) (*dto.MovieScreeningRow, error) {
	sc, err := s.Repo.ScreeningRepo.GetByID(id)
	if err != nil {
		s.Logger.Error("Error get screening by id Usecase: ", zap.Error(err))
		return nil, err
	}
	m, _ := s.Repo.MovieRepo.GetByID(sc.MovieID)
	st, _ := s.Repo.StudioRepo.GetByID(sc.StudioID)

	// Format time into string
	startTime := sc.StartTime.Format("15.04")
	endTime := sc.EndTime.Format("15.04")
	date := sc.StartTime.Format("02-01-2006")

	response := &dto.MovieScreeningRow{
		Screening: dto.ScreeningResponse{
			ScreeningID: sc.ID,
			StartTime: startTime,
			EndTime: endTime,
		},
		Date: date,
		Movie: *m,
		Studio: dto.StudioResponse{
			StudioID: st.ID,
			Name: st.Name,
		},
	}
	return response, err
}

func (s *screeningUsecase) Create(data dto.ScreeningRequest) error {
	err := s.Repo.ScreeningRepo.Create(data)
	if err != nil {
		s.Logger.Error("Error create screening Usecase: ", zap.Error(err))
		return err
	}
	return err
}

func (s *screeningUsecase) Update(id int, data dto.UpdateScreeningRequest) error {
	err := s.Repo.ScreeningRepo.Update(id, data)
	if err != nil {
		s.Logger.Error("Error update screening Usecase: ", zap.Error(err))
		return err
	}
	return nil
}

func (s *screeningUsecase) Delete(id int) error {
	err := s.Repo.ScreeningRepo.Delete(id)
	if err != nil {
		s.Logger.Error("Error delete screening Usecase: ", zap.Error(err))
		return err
	}
	return nil
}