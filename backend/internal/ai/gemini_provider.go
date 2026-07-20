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
	"strings"
	"time"
)

//go:embed prompts/extract_profile.txt
var extractProfilePrompt string

type GeminiProvider struct {
	apiKey string
}

func NewGeminiProvider() (*GeminiProvider, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}
	return &GeminiProvider{apiKey: apiKey}, nil
}

func (p *GeminiProvider) ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error) {
	currentYear := time.Now().Year()
	prompt := fmt.Sprintf(extractProfilePrompt, currentYear, "```json", chatHistory)

	jsonText, err := p.generateContent(ctx, prompt)
	if err != nil {
		return nil, err
	}

	jsonText = strings.TrimSpace(jsonText)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	var result ExtractionResult
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON: %v, raw text: %s", err, jsonText)
	}

	return &result, nil
}

func (p *GeminiProvider) generateContent(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", p.apiKey)

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
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini API error: status %d, body: %s", resp.StatusCode, string(respBody))
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

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
