package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrReportNotFound = errors.New("report not found")

type ReportStatus string

const (
	ReportStatusUnknown ReportStatus = "unknown"
	ReportStatusPending ReportStatus = "pending"
	ReportStatusReady   ReportStatus = "ready"
	ReportStatusFailed  ReportStatus = "failed"
)

func ParseReportStatus(status string) ReportStatus {
	switch status {
	case "pending":
		return ReportStatusPending
	case "ready":
		return ReportStatusReady
	case "failed":
		return ReportStatusFailed
	default:
		return ReportStatusUnknown
	}
}

type Report struct {
	// Unique identifier of the report.
	ID uuid.UUID

	// Query used to generate the report.
	Query string

	// Current status of the report.
	Status ReportStatus

	CreatedAt time.Time

	// Report content is available only when Status is ReportStatusReady.
	Content *ReportContent
}

type ProductMeta struct {
	// Product name.
	Name string

	// Product vendor.
	Vendor string

	// Classification of the product, e.g. "antivirus", "vpn", etc.
	Classification string

	// A short description of the product.
	ShortDesc string

	// List of alternative products.
	Alternatives []string
}

type ReportContent struct {
	// ProductMeta information about the product.
	ProductMeta ProductMeta

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
