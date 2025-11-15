package usecases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jieggii/softrust/services/backend/internal/adapters/mongoadapter"
	"github.com/jieggii/softrust/services/backend/internal/domain"
)

const reportCacheTTL = 1 * time.Hour

type Service struct {
	repo *mongoadapter.Repo

	reportGenerator *ReportContentGenerator

	reportCacheTTL time.Duration
	log            *log.Logger
}

func NewService(repo *mongoadapter.Repo, reportGenerator *ReportContentGenerator, log *log.Logger) *Service {
	return &Service{
		repo:            repo,
		reportGenerator: reportGenerator,
		reportCacheTTL:  reportCacheTTL,
		log:             log,
	}
}

// GenerateReport starts generating product report and returns report ID.
func (s *Service) GenerateReport(ctx context.Context, query string) (uuid.UUID, error) {
	// normalize query for consistency:
	normalizedQuery := domain.NormalizeQuery(query)

	// get existing report by query (in case it was generated recently):
	report, err := s.repo.GetReportByQuery(ctx, normalizedQuery)
	if err == nil {
		if time.Now().UTC().Sub(report.CreatedAt) <= s.reportCacheTTL {
			// return existing report ID if it's still fresh:
			return report.ID, nil
		}
	} else if !errors.Is(err, domain.ErrReportNotFound) {
		return uuid.Nil, err
	}

	// generate new report, because no fresh report exists:
	reportID := uuid.New()

	// create report in the repository:
	if err := s.repo.CreateReport(ctx, reportID, normalizedQuery); err != nil {
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
