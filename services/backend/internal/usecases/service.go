package usecases

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jieggii/softrust/services/backend/internal/adapters/mongoadapter"
	"github.com/jieggii/softrust/services/backend/internal/domain"
)

type Service struct {
	repo *mongoadapter.Repo

	reportGenerator *ReportContentGenerator

	log *log.Logger
}

func NewService(repo *mongoadapter.Repo, reportGenerator *ReportContentGenerator, log *log.Logger) *Service {
	return &Service{
		repo:            repo,
		reportGenerator: reportGenerator,
		log:             log,
	}
}

// GenerateReport starts generating product report and returns report ID.
func (s *Service) GenerateReport(ctx context.Context, query string) (uuid.UUID, error) {
	reportID := uuid.New()
	queryHash := domain.NormalizeQuery(query)

	// create report in the repository:
	if err := s.repo.CreateReport(ctx, reportID, queryHash); err != nil {
		return uuid.Nil, err
	}

	// start generating report in a separate goroutine:
	go func() {
		ctx := context.Background()

		content, err := s.reportGenerator.GenerateReportContent(ctx, reportID, query)
		if err != nil {
			s.log.Printf("error generating report: %v", err)
			return
		}

		if err := s.repo.UpdateReportContent(ctx, reportID, content); err != nil {
			s.log.Printf("error updating report content: %v", err)
		}

		if err := s.repo.UpdateReportStatus(ctx, reportID, domain.ReportStatusReady); err != nil {
			s.log.Printf("error updating report status: %v", err)
		}
	}()

	return reportID, nil
}

func (s *Service) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.Report, error) {
	report, err := s.repo.GetReportByReportID(ctx, id)
	if err != nil {
		return nil, err
	}
	return report, nil
}
