package mongoadapter

import (
	"time"

	"github.com/google/uuid"
	"github.com/jieggii/softrust/services/backend/internal/domain"
)

type ReportDocument struct {
	ReportID string `bson:"report_id"`

	Query string `bson:"query"`

	CreatedAt time.Time `bson:"created_at"`

	Status string `bson:"status"`

	Content *ReportContent `bson:"content,omitempty"`
}

func (r *ReportDocument) Domain() (*domain.Report, error) {
	var content *domain.ReportContent
	if r.Content != nil {
		content = r.Content.Domain()
	}

	reportID, err := uuid.Parse(r.ReportID)
	if err != nil {
		return nil, err
	}

	return &domain.Report{
		ID:        reportID,
		Query:     r.Query,
		CreatedAt: r.CreatedAt,
		Status:    domain.ParseReportStatus(r.Status),
		Content:   content,
	}, nil
}

type ProductMeta struct {
	Name string `bson:"name"`

	// Product vendor.
	Vendor string `bson:"vendor"`

	// Classification of the product, e.g. "antivirus", "vpn", etc.
	Classification string `bson:"classification"`

	// A short description of the product.
	ShortDesc string `bson:"shortDesc"`

	// List of alternative products.
	Alternatives []string `bson:"alternatives"`
}

func NewProductMeta(m domain.ProductMeta) ProductMeta {
	return ProductMeta{
		Name:           m.Name,
		Vendor:         m.Vendor,
		Classification: m.Classification,
		ShortDesc:      m.ShortDesc,
		Alternatives:   m.Alternatives,
	}
}

func (m *ProductMeta) Domain() domain.ProductMeta {
	return domain.ProductMeta{
		Name:           m.Name,
		Vendor:         m.Vendor,
		Classification: m.Classification,
		ShortDesc:      m.ShortDesc,
		Alternatives:   m.Alternatives,
	}
}

type ReportContent struct {
	ProductMeta ProductMeta `bson:"product_meta"`
}

func NewReportContent(c domain.ReportContent) ReportContent {
	return ReportContent{
		ProductMeta: NewProductMeta(c.ProductMeta),
	}
}

func (r *ReportContent) Domain() *domain.ReportContent {
	return &domain.ReportContent{
		ProductMeta: r.ProductMeta.Domain(),
	}
}
