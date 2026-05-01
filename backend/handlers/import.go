package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type ImportHandler struct {
	llm *LLMClient
}

func NewImportHandler(llm *LLMClient) *ImportHandler {
	return &ImportHandler{llm: llm}
}

func (h *ImportHandler) Import(w http.ResponseWriter, r *http.Request) {
	log.Printf("import: parsing multipart form")
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		log.Printf("import: parse multipart form failed: %v", err)
		writeError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("import: get image field failed: %v", err)
		writeError(w, http.StatusBadRequest, "image field required")
		return
	}
	defer file.Close()

	log.Printf("import: received image %q (%d bytes), content-type=%s",
		header.Filename, header.Size, header.Header.Get("Content-Type"))

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("import: read image bytes failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to read image")
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
		log.Printf("import: content-type missing, defaulting to image/jpeg")
	}

	log.Printf("import: sending to LLM (model=%s, url=%s)", h.llm.model, h.llm.apiURL)
	recipe, err := h.llm.parseRecipeFromImage(r.Context(), imgBytes, mimeType)
	if err != nil {
		log.Printf("import: LLM parsing failed: %v", err)
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("LLM error: %v", err))
		return
	}

	log.Printf("import: parsed recipe %q (%d ingredients, %d steps, %d tags)",
		recipe.Title, len(recipe.Ingredients), len(recipe.Steps), len(recipe.Tags))
	writeJSON(w, http.StatusOK, recipe)
}
