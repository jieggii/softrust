package oapiadapter

import (
	"github.com/jieggii/softrust/services/backend/internal/domain"
	"github.com/jieggii/softrust/services/backend/internal/oapi"
)

func GetReportByID200Response(r *domain.Report) oapi.GetReportByID200JSONResponse {
	var content *oapi.ReportContent
	if r.Content != nil {
		content = reportContent(r.Content)
	}

	return oapi.GetReportByID200JSONResponse{
		Id:        r.ID.String(),
		Query:     r.Query,
		CreatedAt: r.CreatedAt.UTC(),
		Status:    reportStatus(r.Status),
		Content:   content,
	}
}

func reportContent(c *domain.ReportContent) *oapi.ReportContent {
	return &oapi.ReportContent{
		Meta:               productMeta(c.ProductMeta),
		SecurityAssessment: productSecurityAssessment(c.SecurityAssessment),
	}
}

func productMeta(m domain.ProductMeta) oapi.ProductMeta {
	return oapi.ProductMeta{
		Name:             m.Name,
		Vendor:           m.Vendor,
		Classification:   m.Classification,
		ShortDescription: m.ShortDesc,
		Alternatives:     m.Alternatives,
	}
}

func productSecurityAssessment(a domain.ProductSecurityAssessment) oapi.ProductSecurityAssessment {
	issues := make([]oapi.Issue, len(a.KeyIssues))
	for i, v := range a.KeyIssues {
		issues[i] = issue(v)
	}

	return oapi.ProductSecurityAssessment{
		KeyIssues:                  issues,
		SecurityScore:              a.SecurityScore,
		SecurityScoreJustification: a.SecurityScoreJustification,
		Summary:                    a.Summary,
		Verdict:                    a.Verdict,
	}
}

func issue(i domain.Issue) oapi.Issue {
	return oapi.Issue{
		Title:       i.Title,
		Description: i.Description,
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
