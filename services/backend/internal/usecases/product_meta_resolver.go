package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jieggii/softrust/services/backend/internal/domain"
	"google.golang.org/genai"
)

type ProductMeta struct {
	Name           string   `json:"name"`
	Vendor         string   `json:"vendor"`
	Classification string   `json:"classification"`
	ShortDesc      string   `json:"short_description"`
	Alternatives   []string `json:"alternatives"`
}

// ProductMetaResolver is a service that resolves product name from a free-form query.
type ProductMetaResolver struct {
	Client *genai.Client
}

var ErrProductMetaResolutionFailed = errors.New("product meta resolution failed")

// ResolveMeta either resolves product name or returns an error if it is not possible.
func (r *ProductMetaResolver) ResolveMeta(ctx context.Context, query string) (domain.ProductMeta, error) {
	resp, err := r.Client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",                    // choose your model
		genai.Text(metaResolverPrompt(query)), // use the prompt
		nil,                                   // default options
	)
	if err != nil {
		return domain.ProductMeta{}, fmt.Errorf("%w: %v", ErrProductMetaResolutionFailed, err)
	}

	meta, err := parseProductMeta(resp.Text())
	if err != nil {
		return domain.ProductMeta{}, fmt.Errorf("parse AI response: %w", err)
	}

	return domain.ProductMeta{
		Name:           meta.Name,
		Vendor:         meta.Vendor,
		Classification: meta.Classification,
		ShortDesc:      meta.ShortDesc,
		Alternatives:   meta.Alternatives,
	}, nil

	//return domain.ProductMeta{
	//	Name:           "stub name",
	//	Vendor:         "staub vendor",
	//	Classification: "stub classification",
	//	ShortDesc:      "stub short description",
	//	Alternatives:   []string{"alt1", "alt2", "alt3"},
	//}, nil
}

func parseProductMeta(response string) (*ProductMeta, error) {
	// Parse JSON response
	text := strings.TrimSpace(response)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var meta ProductMeta
	if err := json.Unmarshal([]byte(text), &meta); err != nil {
		return nil, fmt.Errorf("%w: failed to parse JSON: %v", ErrProductMetaResolutionFailed, err)
	}

	return &meta, nil
}

// prompt generates a prompt for the language model to resolve product name.
func metaResolverPrompt(query string) string {
	return fmt.Sprintf(`Analyze this product query: "%s"

Return ONLY valid JSON with this exact structure (no markdown, no explanation):
{
  "name": "product name",
  "vendor": "manufacturer or vendor name",
  "classification": "product category",
  "short_description": "2-4 sentences about the product and its key features",
  "alternatives": ["alternative1", "alternative2", "alternative3"]
}

Examples:
- "Windows 11" → vendor: "Microsoft", classification: "Operating System", etc
- "7zip" → vendor: "Igor Pavlov", classification: "File Archiver", etc
- "steam" → vendor: "Valve Corporation", classification: "Digital Distribution Platform", etc`, query)
}
