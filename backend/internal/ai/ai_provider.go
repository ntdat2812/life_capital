package ai

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

//go:embed prompts/extract_profile.txt
var extractProfilePrompt string

type ExtractionResult struct {
	DateOfBirth                 *string `json:"date_of_birth"`
	RiskScore                   int     `json:"risk_score"`
	RiskTolerance               string  `json:"risk_tolerance"`
	MaritalStatus               string  `json:"marital_status"`
	TotalMonthlyIncome          float64 `json:"total_monthly_income"`
	EssentialMonthlyExpense     float64 `json:"essential_monthly_expense"`
	DiscretionaryMonthlyExpense float64 `json:"discretionary_monthly_expense"`
	FITargetAmount              float64 `json:"fi_target_amount"`
}

type AIProvider interface {
	Name() string
	Model() string
	ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error)
}

// ExecuteWithFallback is a generic helper to run any AIProvider method with a fallback mechanism.
func ExecuteWithFallback[T any](providers []AIProvider, action func(AIProvider) (T, error)) (T, error) {
	var zero T
	if len(providers) == 0 {
		return zero, fmt.Errorf("no AI providers configured")
	}

	var lastErr error
	for _, provider := range providers {
		result, err := action(provider)
		if err == nil {
			return result, nil
		}
		fmt.Printf("Warning: AI provider [%s] with model [%s] failed, trying next. Error: %v\n", provider.Name(), provider.Model(), err)
		lastErr = err
	}

	return zero, fmt.Errorf("all AI providers failed, last error: %v", lastErr)
}

// ExtractProfileHelper handles common JSON extraction logic for all providers
func extractProfileHelper(ctx context.Context, chatHistory string, generateContent func(context.Context, string) (string, error)) (*ExtractionResult, error) {
	currentYear := time.Now().Year()
	cleanTemplate := strings.TrimSpace(extractProfilePrompt)
	prompt := fmt.Sprintf(cleanTemplate, currentYear, chatHistory)

	jsonText, err := generateContent(ctx, prompt)
	if err != nil {
		return nil, err
	}

	jsonText = extractJSON(jsonText)

	var result ExtractionResult
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON: %v, raw text: %s", err, jsonText)
	}

	return &result, nil
}

// doHTTPRequest is a common helper for sending HTTP requests and reading the response body.
func doHTTPRequest(ctx context.Context, method, url string, headers map[string]string, body []byte) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, respBody, nil
}

// extractJSON safely extracts the JSON object from a string that might contain markdown wrappers or random text.
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start != -1 && end != -1 && end > start {
		return s[start : end+1]
	}
	return s
}
