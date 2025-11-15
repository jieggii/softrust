package usecases

import (
	"log"

	"github.com/google/uuid"
	"github.com/jieggii/softrust/services/backend/internal/adapters/mongoadapter"
	"github.com/jieggii/softrust/services/backend/internal/domain"
)

type Service struct {
	repo *mongoadapter.Repo

	nameResolver *NameResolver

	log *log.Logger
}

func NewService(repo *mongoadapter.Repo, nameResolver *NameResolver, log *log.Logger) *Service {
	return &Service{
		repo:         repo,
		nameResolver: nameResolver,
		log:          log,
	}
}

// GenerateReport starts generating product report and returns report ID.
func (s *Service) GenerateReport(query string) (uuid.UUID, error) {
	reportID := uuid.New()
	queryHash := domain.QueryHash(query)

	// create report in the repository:
	if err := s.repo.CreateReport(reportID, queryHash); err != nil {
		return uuid.Nil, err
	}

	// start generating report in a separate goroutine:
	go func() {
		if err := s.generateReport(reportID, query); err != nil {
			s.log.Printf("error generating report: %v", err)
		}
	}()

	return reportID, nil
}

func (s *Service) GetReportByID(id uuid.UUID) (mongoadapter.Report, error) {
	return s.repo.GetReportByID(id)
}

func (s *Service) generateReport(id uuid.UUID, query string) error {
	panic("implement me")
}
