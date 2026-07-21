package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/datnguyen/life_capital/backend/internal/model"
)

type GroqProvider struct {
	apiKey    string
	model     string
	maxTokens int
}

func NewGroqProvider() (*GroqProvider, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is required")
	}

	model := os.Getenv("GROQ_MODEL")
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}

	return &GroqProvider{
		apiKey:    apiKey,
		model:     model,
		maxTokens: GetMaxOutputTokens(),
	}, nil
}

func (p *GroqProvider) Name() string {
	return "Groq"
}

func (p *GroqProvider) Model() string {
	return p.model
}

func (p *GroqProvider) ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error) {
	return extractProfileHelper(ctx, chatHistory, p.generateContent)
}

func (p *GroqProvider) GenerateIPS(ctx context.Context, profile *model.InvestorProfile, assets []model.Asset, incomes []*model.IncomeStream, dependents []*model.Dependent, preferredAssets []string) (*IPSExtractionResult, error) {
	return generateIPSHelper(ctx, profile, assets, incomes, dependents, preferredAssets, p.generateContent)
}

func (p *GroqProvider) GenerateThesis(ctx context.Context, req *model.ThesisGenerationRequest) (*model.InvestmentThesis, error) {
	return generateThesisHelper(ctx, req, p.generateContent)
}

func (p *GroqProvider) generateContent(ctx context.Context, prompt string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	log.Printf("=== [Groq] Sending Prompt to AI ===\n%s\n===================================\n", prompt)

	reqBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"response_format": map[string]interface{}{
			"type": "json_object",
		},
		"max_tokens": p.maxTokens,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	headers := map[string]string{
		"Authorization": "Bearer " + p.apiKey,
	}

	statusCode, respBody, err := doHTTPRequest(ctx, "POST", url, headers, bodyBytes)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("groq API error: status %d, body: %s", statusCode, string(respBody))
	}

	var gptResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &gptResp); err != nil {
		return "", err
	}

	if len(gptResp.Choices) == 0 {
		return "", fmt.Errorf("groq returned empty response")
	}

	return gptResp.Choices[0].Message.Content, nil
}

func (p *GroqProvider) AnalyzeLifeEvent(ctx context.Context, promptContext string) (*LifeEventAnalysisResult, error) {
	return analyzeLifeEventHelper(ctx, promptContext, p.generateContent)
}
