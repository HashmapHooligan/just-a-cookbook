package handlers_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/go-chi/chi/v5"

	"justacookbook/handlers"
	"justacookbook/models"
)

func setupImportServer(t *testing.T, llmURL string) *httptest.Server {
	t.Helper()
	llm := handlers.NewLLMClient(llmURL, "test-key", "test-model")
	h := handlers.NewImportHandler(llm)

	r := chi.NewRouter()
	r.Post("/kochbuch/api/recipes/import", h.Import)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

// fakeImageBytes returns a minimal JPEG header — enough for Content-Type detection.
func fakeImageBytes() []byte {
	return []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00}
}

func multipartImage(t *testing.T, data []byte, filename, mimeType string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, filename))
	h.Set("Content-Type", mimeType)
	fw, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("create multipart part: %v", err)
	}
	fw.Write(data)
	w.Close()
	return &buf, w.FormDataContentType()
}

func TestImport_ValidImage(t *testing.T) {
	recipeJSON := `{"title":"Pasta","source":"","ingredients":[{"name":"Spaghetti","amountNumber":200,"amountUnit":"g","emoji":"🍝"}],"steps":[{"description":"Cook pasta."}],"tags":[{"name":"Pasta"}]}`
	llmSrv := mockLLMSuccess(t, recipeJSON)
	server := setupImportServer(t, llmSrv.URL)

	body, ct := multipartImage(t, fakeImageBytes(), "recipe.jpg", "image/jpeg")
	resp, err := http.Post(server.URL+"/kochbuch/api/recipes/import", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	if got.Title != "Pasta" {
		t.Fatalf("expected title Pasta, got %q", got.Title)
	}
	if len(got.Ingredients) != 1 {
		t.Fatalf("expected 1 ingredient, got %d", len(got.Ingredients))
	}
}

func TestImport_LLMMarkdownFenced(t *testing.T) {
	fenced := "```json\n{\"title\":\"Soup\",\"source\":\"\",\"ingredients\":[],\"steps\":[],\"tags\":[]}\n```"
	llmSrv := mockLLMSuccess(t, fenced)
	server := setupImportServer(t, llmSrv.URL)

	body, ct := multipartImage(t, fakeImageBytes(), "recipe.jpg", "image/jpeg")
	resp, err := http.Post(server.URL+"/kochbuch/api/recipes/import", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for markdown-fenced JSON, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	if got.Title != "Soup" {
		t.Fatalf("expected title Soup, got %q", got.Title)
	}
}

func TestImport_MissingImageField(t *testing.T) {
	llmSrv := mockLLMSuccess(t, `{}`)
	server := setupImportServer(t, llmSrv.URL)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("other", "value")
	w.Close()

	resp, err := http.Post(server.URL+"/kochbuch/api/recipes/import", w.FormDataContentType(), &buf)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestImport_LLMReturnsInvalidJSON(t *testing.T) {
	llmSrv := mockLLMSuccess(t, "this is not json")
	server := setupImportServer(t, llmSrv.URL)

	body, ct := multipartImage(t, fakeImageBytes(), "recipe.jpg", "image/jpeg")
	resp, err := http.Post(server.URL+"/kochbuch/api/recipes/import", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

func TestImport_LLMApiError(t *testing.T) {
	llmSrv := mockLLMError(t, http.StatusServiceUnavailable)
	server := setupImportServer(t, llmSrv.URL)

	body, ct := multipartImage(t, fakeImageBytes(), "recipe.jpg", "image/jpeg")
	resp, err := http.Post(server.URL+"/kochbuch/api/recipes/import", ct, body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}
