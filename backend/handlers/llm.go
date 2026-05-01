package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"justacookbook/models"
)

type LLMClient struct {
	apiURL string
	apiKey string
	model  string
}

func NewLLMClient(apiURL, apiKey, model string) *LLMClient {
	return &LLMClient{apiURL: apiURL, apiKey: apiKey, model: model}
}

func (c *LLMClient) chatComplete(ctx context.Context, messages []map[string]any, maxTokens int) (string, error) {
	reqBody := map[string]any{
		"model":      c.model,
		"messages":   messages,
		"max_tokens": maxTokens,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request to %s failed: %w", c.apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API returned %d: %s", resp.StatusCode, string(b))
	}

	var llmResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return "", fmt.Errorf("decode LLM response: %w", err)
	}
	if len(llmResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in LLM response")
	}

	return llmResp.Choices[0].Message.Content, nil
}

func stripMarkdownFence(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) >= 2 {
			end := len(lines) - 1
			if strings.TrimSpace(lines[end]) == "```" {
				end--
			}
			s = strings.TrimSpace(strings.Join(lines[1:end+1], "\n"))
		}
	}
	return s
}

func (c *LLMClient) inferEmojis(ctx context.Context, names []string) ([]string, error) {
	if len(names) == 0 {
		return nil, nil
	}

	namesJSON, err := json.Marshal(names)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(
		"For each ingredient name in this JSON array, return a single fitting emoji. Return ONLY a JSON array of emoji strings in the same order, no markdown, no explanation.\n%s",
		string(namesJSON),
	)

	log.Printf("llm: inferring emojis for %d ingredients", len(names))
	content, err := c.chatComplete(ctx, []map[string]any{
		{"role": "user", "content": prompt},
	}, 300)
	if err != nil {
		return nil, err
	}

	clean := stripMarkdownFence(content)
	log.Printf("llm: emoji response: %s", clean)

	var emojis []string
	if err := json.Unmarshal([]byte(clean), &emojis); err != nil {
		return nil, fmt.Errorf("parse emoji JSON %q: %w", clean, err)
	}
	if len(emojis) != len(names) {
		return nil, fmt.Errorf("expected %d emojis, got %d", len(names), len(emojis))
	}

	return emojis, nil
}

func (c *LLMClient) parseRecipeFromImage(ctx context.Context, imgBytes []byte, mimeType string) (*models.Recipe, error) {
	encoded := base64.StdEncoding.EncodeToString(imgBytes)

	prompt := `Extract the recipe from this image and return ONLY valid JSON with this exact structure:
{
  "title": "Recipe title",
  "source": "Source if visible, otherwise empty string",
  "ingredients": [
    {"name": "ingredient name", "amountNumber": 1.5, "amountUnit": "cup", "emoji": "🥕"}
  ],
  "steps": [
    {"description": "Step description"}
  ],
  "tags": [
    {"name": "tag name"}
  ]
}
Use 1-3 tags maximum (aim for 2). Choose broad category tags only (e.g. "vegetarian", "pasta", "dessert").
Return ONLY the JSON object, no markdown, no explanation.`

	content, err := c.chatComplete(ctx, []map[string]any{
		{
			"role": "user",
			"content": []map[string]any{
				{"type": "text", "text": prompt},
				{"type": "image_url", "image_url": map[string]string{
					"url": fmt.Sprintf("data:%s;base64,%s", mimeType, encoded),
				}},
			},
		},
	}, 16000)
	if err != nil {
		return nil, err
	}

	log.Printf("import: LLM response length=%d chars", len(content))

	clean := stripMarkdownFence(content)

	var recipe models.Recipe
	if err := json.Unmarshal([]byte(clean), &recipe); err != nil {
		log.Printf("import: failed to parse recipe JSON: %v\nraw response:\n%s", err, clean)
		return nil, fmt.Errorf("parse recipe JSON: %w", err)
	}
	return &recipe, nil
}
