package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"justacookbook/models"
)

type ImportHandler struct {
	apiURL string
	apiKey string
	model  string
}

func NewImportHandler(apiURL, apiKey, model string) *ImportHandler {
	return &ImportHandler{apiURL: apiURL, apiKey: apiKey, model: model}
}

func (h *ImportHandler) Import(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "image field required")
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read image")
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	recipe, err := h.parseRecipeFromImage(r, imgBytes, mimeType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("LLM error: %v", err))
		return
	}
	writeJSON(w, http.StatusOK, recipe)
}

func (h *ImportHandler) parseRecipeFromImage(r *http.Request, imgBytes []byte, mimeType string) (*models.Recipe, error) {
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
Return ONLY the JSON object, no markdown, no explanation.`

	reqBody := map[string]any{
		"model": h.model,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{"type": "text", "text": prompt},
					{"type": "image_url", "image_url": map[string]string{
						"url": fmt.Sprintf("data:%s;base64,%s", mimeType, encoded),
					}},
				},
			},
		},
		"max_tokens": 16000,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost,
		h.apiURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("LLM API returned %d: %s", resp.StatusCode, string(b))
	}

	var llmResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return nil, err
	}
	if len(llmResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in LLM response")
	}

	var recipe models.Recipe
	if err := json.Unmarshal([]byte(llmResp.Choices[0].Message.Content), &recipe); err != nil {
		return nil, fmt.Errorf("parse recipe JSON: %w", err)
	}
	return &recipe, nil
}
