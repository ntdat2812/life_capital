package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type GroqProvider struct {
	apiKey string
}

func NewGroqProvider() (*GroqProvider, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is required")
	}
	return &GroqProvider{apiKey: apiKey}, nil
}

func (p *GroqProvider) ExtractProfile(ctx context.Context, chatHistory string) (*ExtractionResult, error) {
	currentYear := time.Now().Year()
	cleanTemplate := strings.TrimSpace(extractProfilePrompt)
	prompt := fmt.Sprintf(cleanTemplate, currentYear, chatHistory)

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
		return nil, fmt.Errorf("failed to parse AI JSON (Groq): %v, raw text: %s", err, jsonText)
	}

	return &result, nil
}

func (p *GroqProvider) generateContent(ctx context.Context, prompt string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	reqBody := map[string]interface{}{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"response_format": map[string]interface{}{
			"type": "json_object",
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	
	client := &http.Client{Timeout: 30 * time.Second}
	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.apiKey)

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("groq API error: status %d, body: %s", resp.StatusCode, string(respBody))
			
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				time.Sleep(time.Duration(attempt) * 2 * time.Second)
				continue
			}
			return "", lastErr
		}

		var gptResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&gptResp); err != nil {
			resp.Body.Close()
			return "", err
		}
		resp.Body.Close()

		if len(gptResp.Choices) == 0 {
			return "", fmt.Errorf("groq returned empty response")
		}

		return gptResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("groq failed after 3 attempts: %v", lastErr)
}
