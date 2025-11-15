package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jieggii/softrust/services/backend/internal/domain"
)

type ReportContentGenerator struct {
	MetaResolver *ProductMetaResolver
	Log          *log.Logger
}

func (r *ReportContentGenerator) GenerateReportContent(ctx context.Context, reportID uuid.UUID, query string) (domain.ReportContent, error) {
	// find out the product meta:
	meta, err := r.MetaResolver.ResolveMeta(ctx, query)
	if err != nil {
		return domain.ReportContent{}, fmt.Errorf("resolve name: %w", err)
	}

	return domain.ReportContent{ProductMeta: meta}, nil
}
