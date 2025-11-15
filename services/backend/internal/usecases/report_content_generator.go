package usecases

import (
	"log"

	"github.com/google/uuid"
)

type ReportContentGenerator struct {
	NameResolver *NameResolver
	Log          *log.Logger
}

func (r *ReportContentGenerator) GenerateReportContent(reportID uuid.UUID, query string) error {
	// find out the product name:
	productName, err := r.NameResolver.ResolveName(query)
	if err != nil {
		return err
	}

	r.Log.Printf("Generating report for product: %s", productName)
	return nil
}
