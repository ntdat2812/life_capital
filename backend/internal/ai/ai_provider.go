package ai

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/datnguyen/life_capital/backend/internal/model"
)

//go:embed prompts/extract_profile.txt
var extractProfilePrompt string

//go:embed prompts/analyze_life_event.txt
var analyzeLifeEventPrompt string

//go:embed prompts/generate_ips.txt
var generateIPSPrompt string

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

type IncomeStreamDraft struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	IsPassive bool    `json:"is_passive"`
}

type DependentDraft struct {
	Name         string  `json:"name"`
	Relationship string  `json:"relationship"`
	MonthlyCost  float64 `json:"monthly_cost"`
}

type LifeEventAnalysisResult struct {
	Category             string              `json:"category"`
	AIImpactAnalysis     string              `json:"ai_impact_analysis"`
	IncomeImpact         float64             `json:"income_impact"`
	ExpenseImpact        float64             `json:"expense_impact"`
	NewMaritalStatus     *string             `json:"new_marital_status"`
	NewRiskScore         *int                `json:"new_risk_score"`
	NewRiskTolerance     *string             `json:"new_risk_tolerance"`
	IncomeStreamsToAdd   []IncomeStreamDraft `json:"income_streams_to_add"`
	DependentsToAdd      []DependentDraft    `json:"dependents_to_add"`
}

type IPSExtractionResult struct {
	TargetAllocation map[string]float64 `json:"target_allocation"`
	DetailedStrategy string             `json:"detailed_strategy"`
}

type AIProvider interface {
	Name() string
	Model() string
	ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error)
	AnalyzeLifeEvent(ctx context.Context, promptContext string) (*LifeEventAnalysisResult, error)
	GenerateIPS(ctx context.Context, profile *model.InvestorProfile, assets []model.Asset, preferredAssets []string) (*IPSExtractionResult, error)
}

// GetMaxOutputTokens retrieves the maximum output tokens from the environment, defaulting to 8192
func GetMaxOutputTokens() int {
	maxTokens := 8192
	if tokensStr := os.Getenv("AI_MAX_OUTPUT_TOKENS"); tokensStr != "" {
		if val, err := strconv.Atoi(tokensStr); err == nil {
			maxTokens = val
		}
	}
	return maxTokens
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

// AnalyzeLifeEventHelper handles common JSON extraction logic for life events
func analyzeLifeEventHelper(ctx context.Context, promptContext string, generateContent func(context.Context, string) (string, error)) (*LifeEventAnalysisResult, error) {
	currentYear := time.Now().Year()
	cleanTemplate := strings.TrimSpace(analyzeLifeEventPrompt)
	
	prompt := fmt.Sprintf(cleanTemplate, currentYear, promptContext)

	jsonText, err := generateContent(ctx, prompt)
	if err != nil {
		return nil, err
	}

	jsonText = extractJSON(jsonText)

	var result LifeEventAnalysisResult
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON for event: %v, raw text: %s", err, jsonText)
	}

	return &result, nil
}

func generateIPSHelper(ctx context.Context, profile *model.InvestorProfile, assets []model.Asset, preferredAssets []string, generate func(context.Context, string) (string, error)) (*IPSExtractionResult, error) {
	cleanTemplate := strings.TrimSpace(generateIPSPrompt)
	
	age := 0
	if profile != nil && !profile.DateOfBirth.IsZero() {
		age = time.Now().Year() - profile.DateOfBirth.Year()
	}

	maritalStatus := ""
	essentialExpense := 0.0
	discretionaryExpense := 0.0
	riskScore := 0
	riskTolerance := ""
	fiTargetAmount := 0.0

	if profile != nil {
		maritalStatus = profile.MaritalStatus
		essentialExpense = profile.EssentialMonthlyExpense
		discretionaryExpense = profile.DiscretionaryMonthlyExpense
		riskScore = profile.RiskScore
		riskTolerance = profile.RiskTolerance
		fiTargetAmount = profile.FITargetAmount
	}

	var assetListStr string
	if len(assets) == 0 {
		assetListStr = "Chưa có tài sản nào được khai báo."
	} else {
		for _, a := range assets {
			qty := 0.0
			if a.Quantity != nil {
				qty = *a.Quantity
			}
			assetListStr += fmt.Sprintf("- %s (%s): %f (Giá trị: %f)\n", a.Name, a.Category, qty, a.CurrentValue)
		}
	}

	preferredAssetsStr := ""
	if len(preferredAssets) > 0 {
		preferredAssetsStr = strings.Join(preferredAssets, ", ")
	} else {
		preferredAssetsStr = "Khách hàng không chỉ định. Vui lòng tự động tư vấn phân bổ thông minh dựa trên khẩu vị rủi ro và các mục tiêu."
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.Local
	}
	currentTimeStr := time.Now().In(loc).Format("15:04 02/01/2006")

	prompt := fmt.Sprintf(cleanTemplate,
		age,
		maritalStatus,
		0.0, // passive income omitted for now
		essentialExpense,
		discretionaryExpense,
		riskScore,
		riskTolerance,
		fiTargetAmount,
		preferredAssetsStr,
		currentTimeStr, // %s in Tình trạng tài sản
		assetListStr,
		currentTimeStr, // %s in 1. Bức tranh Dòng tiền
	)

	respText, err := generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	respText = extractJSON(respText)

	var result IPSExtractionResult
	if err := json.Unmarshal([]byte(respText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON for IPS: %v, raw text: %s", err, respText)
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
