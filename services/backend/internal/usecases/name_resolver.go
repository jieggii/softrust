package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// NameResolver is a service that resolves product name from a free-form query.
type NameResolver struct {
	client *genai.Client
}

func NewNameResolver(client *genai.Client) *NameResolver {
	return &NameResolver{
		client: client,
	}
}

var ErrResolutionFailed = errors.New("product name resolution failed")

// ResolveName either resolves product name or returns an error if it is not possible.
func (r *NameResolver) ResolveName(query string) (string, error) {
	ctx := context.Background()

	resp, err := r.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",        // choose your model
		genai.Text(prompt(query)), // use the prompt
		nil,                       // default options
	)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrResolutionFailed, err)
	}

	// Trim whitespace from the output
	name := strings.TrimSpace(resp.Text())
	if name == "" {
		return "", ErrResolutionFailed
	}

	return name, nil
}

// prompt generates a prompt for the language model to resolve product name.
func prompt(query string) string {
	return fmt.Sprintf(`Hi! Here is a query for determining a software product: "%s".
Please output only the canonical product name as a short string.
Examples:
- "https://www.microsoft.com/en-us/windows" -> Windows
- "https://github.com/ip7z/7zip" → "7Zip"
- "steam" → "Steam"
Output format: single string only, no explanations.`, query)
}
