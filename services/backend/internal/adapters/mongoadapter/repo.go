package mongoadapter

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jieggii/softrust/services/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	reportIDField = "report_id"
	contentField  = "content"
	statusField   = "status"
	queryField    = "query"
)

const reportsCollection = "reports"

type Repo struct {
	DB *mongo.Database
}

func (r *Repo) CreateReport(ctx context.Context, id uuid.UUID, queryHash string) error {
	doc := &ReportDocument{
		ReportID:  id.String(),
		Query:     queryHash,
		CreatedAt: time.Now().UTC(),
		Status:    string(domain.ReportStatusPending),
	}

	if _, err := r.DB.Collection(reportsCollection).InsertOne(ctx, doc); err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetReportByReportID(ctx context.Context, id uuid.UUID) (*domain.Report, error) {
	var doc ReportDocument
	filter := bson.M{reportIDField: id.String()}

	if err := r.DB.Collection(reportsCollection).FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrReportNotFound
		}
		return nil, err
	}

	return doc.Domain()
}

func (r *Repo) UpdateReportStatus(ctx context.Context, id uuid.UUID, status domain.ReportStatus) error {
	filter := bson.M{reportIDField: id.String()}
	query := bson.M{
		"$set": bson.M{
			statusField: string(status),
		},
	}

	res, err := r.DB.Collection(reportsCollection).UpdateOne(ctx, filter, query)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return domain.ErrReportNotFound
	}

	return nil
}

func (r *Repo) GetReportByQuery(ctx context.Context, query string) (*domain.Report, error) {
	var doc ReportDocument
	filter := bson.M{queryField: query}

	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}) // newest first

	if err := r.DB.Collection(reportsCollection).FindOne(ctx, filter, opts).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrReportNotFound
		}
		return nil, err
	}

	return doc.Domain()
}

func (r *Repo) UpdateReportContent(ctx context.Context, id uuid.UUID, content domain.ReportContent) error {

	filter := bson.M{reportIDField: id.String()}
	query := bson.M{
		"$set": bson.M{
			contentField: NewReportContent(content),
		},
	}

	res, err := r.DB.Collection(reportsCollection).UpdateOne(ctx, filter, query)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return domain.ErrReportNotFound
	}

	return nil
}
