package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/datnguyen/life_capital/backend/internal/model"
)

type GeminiProvider struct {
	apiKey    string
	model     string
	maxTokens int
}

func NewGeminiProvider() (*GeminiProvider, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-3.5-flash"
	}

	return &GeminiProvider{
		apiKey:    apiKey,
		model:     model,
		maxTokens: GetMaxOutputTokens(),
	}, nil
}

func (p *GeminiProvider) Name() string {
	return "Gemini"
}

func (p *GeminiProvider) Model() string {
	return p.model
}

func (p *GeminiProvider) ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error) {
	return extractProfileHelper(ctx, chatHistory, p.generateContent)
}

func (p *GeminiProvider) GenerateIPS(ctx context.Context, profile *model.InvestorProfile, assets []model.Asset, incomes []*model.IncomeStream, dependents []*model.Dependent, preferredAssets []string) (*IPSExtractionResult, error) {
	return generateIPSHelper(ctx, profile, assets, incomes, dependents, preferredAssets, p.generateContent)
}

func (p *GeminiProvider) GenerateThesis(ctx context.Context, req *model.ThesisGenerationRequest) (*model.InvestmentThesis, error) {
	return generateThesisHelper(ctx, req, p.generateContent)
}

func (p *GeminiProvider) generateContent(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)
	log.Printf("=== [Gemini] Sending Prompt to AI ===\n%s\n===================================\n", prompt)

	reqBody := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]interface{}{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"maxOutputTokens":  p.maxTokens,
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)

	// Retry up to 3 times for transient errors (503, 429)
	maxRetries := 3
	var statusCode int
	var respBody []byte
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		statusCode, respBody, err = doHTTPRequest(ctx, "POST", url, nil, bodyBytes)
		if err != nil {
			return "", err
		}

		if statusCode == http.StatusOK {
			break
		}

		// Retry on transient errors
		if (statusCode == 503 || statusCode == 429) && attempt < maxRetries-1 {
			waitSec := (attempt + 1) * 3 // 3s, 6s
			log.Printf("⏳ Gemini returned %d, retrying in %ds (attempt %d/%d)...", statusCode, waitSec, attempt+1, maxRetries)
			select {
			case <-time.After(time.Duration(waitSec) * time.Second):
			case <-ctx.Done():
				return "", ctx.Err()
			}
			continue
		}

		return "", fmt.Errorf("gemini API error: status %d, body: %s", statusCode, string(respBody))
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API error after %d retries: status %d, body: %s", maxRetries, statusCode, string(respBody))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

func (p *GeminiProvider) AnalyzeLifeEvent(ctx context.Context, promptContext string) (*LifeEventAnalysisResult, error) {
	return analyzeLifeEventHelper(ctx, promptContext, p.generateContent)
}

func (p *GeminiProvider) GenerateMonthlyReview(ctx context.Context, replacements map[string]string) (*model.MonthlyReviewRecommendationResponse, error) {
	return generateMonthlyReviewHelper(ctx, replacements, p.generateContent)
}
