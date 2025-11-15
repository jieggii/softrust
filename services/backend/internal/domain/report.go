package domain

import "github.com/google/uuid"

type ReportStatus string

const (
	ReportStatusPending ReportStatus = "pending"
	ReportStatusReady   ReportStatus = "ready"
	ReportStatusFailed  ReportStatus = "failed"
)

type Report struct {
	// Unique identifier of the report.
	ID uuid.UUID

	// Query used to generate the report.
	Query string

	// Current status of the report.
	Status ReportStatus

	// Report content is available only when Status is ReportStatusReady.
	Content *ReportContent
}

type ReportContent struct {
	// Product title.
	Title string

	// Product vendor.
	Vendor string

	// Checks that were run on report.
	Checks []Check

	// Verdict that summarizes all analyzes.
	Verdict string

	// Score from 0 to 100 that is based on the verdict.
	Score int
}

// Check represents a single check that has been run on a product.
type Check interface {
	// Name returns analyze name.
	Name() string

	Verdict() string
}
