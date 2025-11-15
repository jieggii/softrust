package oapiadapter

import (
	"github.com/jieggii/softrust/services/backend/internal/domain"
	"github.com/jieggii/softrust/services/backend/internal/oapi"
)

func GetReportByID200Response(r *domain.Report) oapi.GetReportByID200JSONResponse {
	return oapi.GetReportByID200JSONResponse{
		Id:        r.ID.String(),
		Query:     r.Query,
		CreatedAt: r.CreatedAt.UTC(),
		Status:    reportStatus(r.Status),
	}
}

func reportStatus(s domain.ReportStatus) oapi.ReportGetResponseStatus {
	switch s {
	case domain.ReportStatusPending:
		return oapi.Pending
	case domain.ReportStatusReady:
		return oapi.Ready
	case domain.ReportStatusFailed:
		return oapi.Failed
	default:
		return oapi.Unknown
	}
}
