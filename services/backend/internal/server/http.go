package server

import (
	"context"
	"log"

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
	reportID, err := s.svc.GenerateReport(request.Body.Query)
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
	//TODO implement me
	panic("implement me")
}
