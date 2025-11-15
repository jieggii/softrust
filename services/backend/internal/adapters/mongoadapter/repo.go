package mongoadapter

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	// Hash of the query string used to generate the report.
	QueryHash string `bson:"query_hash"`

	// Date and time when the report was created.
	CreatedAt time.Time `bson:"created_at"`
}

type Repo struct {
}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) CreateReport(id uuid.UUID, queryHash string) error {
	return nil // todo implement
}

func (r *Repo) GetReportByID(id uuid.UUID) (*Report, error) {
	return nil, nil // todo implement
}
