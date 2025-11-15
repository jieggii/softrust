package server

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jieggii/softrust/services/backend/internal/adapters/oapiadapter"
	"github.com/jieggii/softrust/services/backend/internal/domain"
	"github.com/jieggii/softrust/services/backend/internal/oapi"
	"github.com/jieggii/softrust/services/backend/internal/usecases"
)

var _ oapi.StrictServerInterface = (*HTTP)(nil)

type HTTP struct {
	svc *usecases.Service

	log *log.Logger
}

func NewHTTP(svc *usecases.Service, log *log.Logger) *HTTP {
	return &HTTP{
		svc: svc,
		log: log,
	}
}

func (s *HTTP) CreateReport(ctx context.Context, request oapi.CreateReportRequestObject) (oapi.CreateReportResponseObject, error) {
	reportID, err := s.svc.GenerateReport(ctx, request.Body.Query)
	if err != nil {
		s.log.Printf("start generating report: %v", err)
		errMsg := err.Error()

		return oapi.CreateReport500JSONResponse{
			Error: &errMsg,
		}, nil
	}

	return oapi.CreateReport200JSONResponse{
		Id: reportID.String(),
	}, nil
}

func (s *HTTP) GetReportByID(ctx context.Context, request oapi.GetReportByIDRequestObject) (oapi.GetReportByIDResponseObject, error) {
	reportID, err := uuid.Parse(request.Id)
	if err != nil {
		errMsg := "invalid report ID"
		return oapi.GetReportByID400JSONResponse{
			Error: &errMsg,
		}, nil
	}

	report, err := s.svc.GetReportByID(ctx, reportID)
	if err != nil {
		if errors.Is(err, domain.ErrReportNotFound) {
			errMsg := "report not found"
			return oapi.GetReportByID404JSONResponse{
				Error: &errMsg,
			}, nil
		}

		s.log.Printf("get report by ID: %v", err)
		errMsg := err.Error()
		return oapi.GetReportByID500JSONResponse{
			Error: &errMsg,
		}, nil
	}

	return oapiadapter.GetReportByID200Response(report), nil
}
