package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jieggii/softrust/services/backend/internal/domain"
	"google.golang.org/genai"
)

type Issue struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ProductSecurityAssessment struct {
	Summary                    string  `json:"summary"`
	KeyIssues                  []Issue `json:"key_issues"`
	Verdict                    string  `json:"verdict"`
	SecurityScore              string  `json:"security_score"`
	SecurityScoreJustification string  `json:"security_score_justification"`
}

type ProductSecurityAssessor struct {
	Client *genai.Client
}

func (a *ProductSecurityAssessor) AssessSecurity(ctx context.Context, productName string, productVendor string, productClassification string) (domain.ProductSecurityAssessment, error) {
	tools := []*genai.Tool{
		// This instructs the model to use the built-in Google Search tool for grounding.
		{GoogleSearch: &genai.GoogleSearch{}},
	}

	resp, err := a.Client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(securityAssessorPrompt(productName, productVendor, productClassification)), // use the prompt
		&genai.GenerateContentConfig{
			Tools: tools,
		},
	)
	if err != nil {
		return domain.ProductSecurityAssessment{}, fmt.Errorf("%w: %v", ErrProductMetaResolutionFailed, err)
	}

	fmt.Println(resp.Text())
	assessment, err := parseProductSecurityAssessment(resp.Text())
	if err != nil {
		return domain.ProductSecurityAssessment{}, fmt.Errorf("parse AI response: %w", err)
	}

	// return the assessment in domain format:
	keyIssues := make([]domain.Issue, len(assessment.KeyIssues))
	for i, issue := range assessment.KeyIssues {
		keyIssues[i] = domain.Issue{
			Title:       issue.Title,
			Description: issue.Description,
		}
	}
	return domain.ProductSecurityAssessment{
		Summary:                    assessment.Summary,
		KeyIssues:                  keyIssues,
		Verdict:                    assessment.Verdict,
		SecurityScore:              assessment.SecurityScore,
		SecurityScoreJustification: assessment.SecurityScoreJustification,
	}, nil
}

func parseProductSecurityAssessment(response string) (*ProductSecurityAssessment, error) {
	text := strings.TrimSpace(response)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var assessment ProductSecurityAssessment
	if err := json.Unmarshal([]byte(text), &assessment); err != nil {
		return nil, fmt.Errorf("%w: failed to parse JSON: %v", ErrProductMetaResolutionFailed, err)
	}

	return &assessment, nil
}

func securityAssessorPrompt(productName string, productVendor string, productClassification string) string {
	return fmt.Sprintf(`You are a cybersecurity expert. Your task is to analyze the security of the provided software.

CONTEXT DATA (Static):
- Software Name: %s
- Vendor: %s
- Classification: %s 

INSTRUCTIONS:
1.  **Search:** Use **Google Search** (Grounding) to find the **most recent (2024–2025)** critical **CVEs, exploits, and official security advisories** for the %s.
2.  **Analyze:** Synthesize the static context with the search results. Base the score on vulnerability severity, patch timeliness, and support status.
3.  **Output:** Provide the result **strictly as a JSON object**.

**REQUIRED JSON FORMAT:**

{
	"summary": "[Brief summary of security status (max 3 sentences)]",
	"key_issues": [
		{
			"title": "[for example most critical recent CVE/Exploit]",
			"description": "[brief description and impact]"
		},
		{
			"title": "[Secondary security issue/weakness]",
			"description": "[Brief description.]"
		}
	],
	"verdict": "[key recommendation (e.g: you can use freely/never use)]",
	"security_score": "[IMPORTANT: json data type must be STRING; string must contain ONLY ONE INTEGER from 0 to 100 representing the security score. e.g.: 85]",
	"security_score_justification": "Brief reason for the assigned score (max 2 sentences).",
}
`, productName, productVendor, productClassification, productName)
}
