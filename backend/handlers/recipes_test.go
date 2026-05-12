package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"

	"justacookbook/db"
	"justacookbook/handlers"
	"justacookbook/models"
)

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	h := handlers.NewRecipeHandler(database)

	r := chi.NewRouter()
	r.Get("/kochbuch/api/recipes", h.List)
	r.Post("/kochbuch/api/recipes", h.Create)
	r.Get("/kochbuch/api/recipes/{id}", h.Get)
	r.Put("/kochbuch/api/recipes/{id}", h.Update)
	r.Delete("/kochbuch/api/recipes/{id}", h.Delete)

	return httptest.NewServer(r)
}

func postJSON(t *testing.T, server *httptest.Server, path string, body any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	resp, err := http.Post(server.URL+path, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func decode[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	resp.Body.Close()
	return v
}

func sampleRecipe() models.Recipe {
	return models.Recipe{
		Title:  "Pasta Carbonara",
		Source: "Nonna",
		Ingredients: []models.Ingredient{
			{Name: "Spaghetti", AmountNumber: ptr(200.0), AmountUnit: "g", Emoji: "🍝"},
			{Name: "Eggs", AmountNumber: ptr(3.0), Emoji: "🥚"},
		},
		Steps: []models.Step{
			{Description: "Cook pasta al dente."},
			{Description: "Mix eggs with cheese."},
		},
		Tags: []models.Tag{{Name: "Italian"}, {Name: "Quick"}},
	}
}

func ptr(f float64) *float64 { return &f }

// --- List ---

func TestList_Empty(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 0 {
		t.Fatalf("expected empty list, got %d items", len(results))
	}
}

func TestList_WithData(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza", Tags: []models.Tag{}})

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes")
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestList_Search(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza Margherita"})

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=pasta")
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Title != "Pasta Carbonara" {
		t.Fatalf("expected Pasta Carbonara, got %s", results[0].Title)
	}
}

func TestList_Search_FTSOperators(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())

	// Special characters and boolean-looking words must not crash — must not return 500.
	operators := []string{"OR", "AND", "NEAR", "NOT", "*", "pasta OR pizza", "pasta AND pizza"}
	for _, q := range operators {
		resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=" + url.QueryEscape(q))
		if err != nil {
			t.Fatalf("GET q=%q: %v", q, err)
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusInternalServerError {
			t.Fatalf("q=%q returned 500", q)
		}
	}
}

func TestList_Search_QuoteInQuery(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: `Chef's "Special" Pasta`})

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=" + url.QueryEscape(`"Special"`))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		t.Fatalf("search with double-quote returned 500")
	}
}

func TestList_Search_EmptyQueryParam(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza"})

	// ?q= (empty value) must return all recipes, not trigger FTS path.
	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=")
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 2 {
		t.Fatalf("expected 2 with empty q, got %d", len(results))
	}
}

func TestList_Search_WhitespaceQuery(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza"})

	// Whitespace-only query has no words after splitting — must return all recipes.
	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=" + url.QueryEscape("   "))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 2 {
		t.Fatalf("expected 2 with whitespace-only q, got %d", len(results))
	}
}

func TestList_Search_MultiWord(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()) // "Pasta Carbonara"
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza Margherita"})

	// Both words must match the same recipe.
	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=" + url.QueryEscape("pasta carbonara"))
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Title != "Pasta Carbonara" {
		t.Fatalf("expected Pasta Carbonara, got %s", results[0].Title)
	}
}

func TestList_Search_MultiWord_NoPartialMatch(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()) // "Pasta Carbonara"
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza Margherita"})

	// Words from different recipes — must return nothing.
	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=" + url.QueryEscape("pasta pizza"))
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 0 {
		t.Fatalf("expected 0, got %d", len(results))
	}
}

func TestList_Search_MultiWord_TitleAndTag(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// sampleRecipe() has title "Pasta Carbonara" and tag "Italian"
	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza Margherita"})

	// One word matches title, other matches tag — recipe must appear.
	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=" + url.QueryEscape("pasta italian"))
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Title != "Pasta Carbonara" {
		t.Fatalf("expected Pasta Carbonara, got %s", results[0].Title)
	}
}

func TestList_Search_MultiWord_ExtraSpaces(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()) // "Pasta Carbonara"
	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Pizza Margherita"})

	// Multiple spaces between words must behave identically to single space.
	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=" + url.QueryEscape("pasta   carbonara"))
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
}

func TestList_Search_NoMatch(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes?q=sushi")
	if err != nil {
		t.Fatal(err)
	}
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 0 {
		t.Fatalf("expected 0, got %d", len(results))
	}
}

// --- Create ---

func TestCreate_Valid(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp := postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe())
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	created := decode[models.Recipe](t, resp)
	if created.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
	if created.Title != "Pasta Carbonara" {
		t.Fatalf("expected Pasta Carbonara, got %s", created.Title)
	}
	if len(created.Ingredients) != 2 {
		t.Fatalf("expected 2 ingredients, got %d", len(created.Ingredients))
	}
	if len(created.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(created.Steps))
	}
	if len(created.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(created.Tags))
	}
}

func TestCreate_MissingTitle(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp := postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Source: "test"})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Post(server.URL+"/kochbuch/api/recipes", "application/json", bytes.NewBufferString("not json"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreate_SharedTags(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "A", Tags: []models.Tag{{Name: "Italian"}}})
	resp := postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "B", Tags: []models.Tag{{Name: "Italian"}}})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestCreate_EmptyNestedArrays(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp := postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Minimal"})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	if len(got.Ingredients) != 0 {
		t.Fatalf("expected 0 ingredients, got %d", len(got.Ingredients))
	}
	if len(got.Steps) != 0 {
		t.Fatalf("expected 0 steps, got %d", len(got.Steps))
	}
	if len(got.Tags) != 0 {
		t.Fatalf("expected 0 tags, got %d", len(got.Tags))
	}
}

// --- Get ---

func TestGet_Found(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes/" + itoa(created.ID))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	if got.Title != "Pasta Carbonara" {
		t.Fatalf("expected Pasta Carbonara, got %s", got.Title)
	}
}

func TestGet_NotFound(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes/99999")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGet_InvalidID(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/kochbuch/api/recipes/notanid")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// --- Update ---

func TestUpdate_Valid(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	updated := created
	updated.Title = "Pasta Amatriciana"
	updated.Ingredients = []models.Ingredient{{Name: "Guanciale", AmountNumber: ptr(150.0), AmountUnit: "g"}}
	updated.Steps = []models.Step{{Description: "Fry guanciale."}}
	updated.Tags = []models.Tag{{Name: "Roman"}}

	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/kochbuch/api/recipes/"+itoa(created.ID),
		jsonBody(updated))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	if got.Title != "Pasta Amatriciana" {
		t.Fatalf("expected Pasta Amatriciana, got %s", got.Title)
	}
	if len(got.Ingredients) != 1 {
		t.Fatalf("expected 1 ingredient, got %d", len(got.Ingredients))
	}
}

func TestUpdate_ClearsAllRelations(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	bare := models.Recipe{
		Title:       "Stripped",
		Ingredients: []models.Ingredient{},
		Steps:       []models.Step{},
		Tags:        []models.Tag{},
	}
	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/kochbuch/api/recipes/"+itoa(created.ID),
		jsonBody(bare))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	if len(got.Ingredients) != 0 {
		t.Fatalf("expected 0 ingredients after clear, got %d", len(got.Ingredients))
	}
	if len(got.Steps) != 0 {
		t.Fatalf("expected 0 steps after clear, got %d", len(got.Steps))
	}
	if len(got.Tags) != 0 {
		t.Fatalf("expected 0 tags after clear, got %d", len(got.Tags))
	}
}

func TestUpdate_AddsRelationsToEmpty(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{Title: "Minimal"}))

	full := sampleRecipe()
	full.Title = "Minimal"
	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/kochbuch/api/recipes/"+itoa(created.ID),
		jsonBody(full))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	if len(got.Ingredients) != 2 {
		t.Fatalf("expected 2 ingredients, got %d", len(got.Ingredients))
	}
	if len(got.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(got.Steps))
	}
	if len(got.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(got.Tags))
	}
}

func TestUpdate_InvalidJSON(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/kochbuch/api/recipes/"+itoa(created.ID),
		bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdate_MissingTitle(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/kochbuch/api/recipes/"+itoa(created.ID),
		jsonBody(models.Recipe{Source: "test"}))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdate_InvalidID(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/kochbuch/api/recipes/notanid",
		jsonBody(models.Recipe{Title: "Test"}))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/kochbuch/api/recipes/99999",
		jsonBody(models.Recipe{Title: "Test"}))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

// --- Delete ---

func TestDelete_Found(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	req, _ := http.NewRequest(http.MethodDelete,
		server.URL+"/kochbuch/api/recipes/"+itoa(created.ID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	getResp, _ := http.Get(server.URL + "/kochbuch/api/recipes/" + itoa(created.ID))
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestDelete_AlreadyDeleted(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	del := func() *http.Response {
		req, _ := http.NewRequest(http.MethodDelete,
			server.URL+"/kochbuch/api/recipes/"+itoa(created.ID), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	if resp := del(); resp.StatusCode != http.StatusNoContent {
		t.Fatalf("first delete: expected 204, got %d", resp.StatusCode)
	}
	if resp := del(); resp.StatusCode != http.StatusNotFound {
		t.Fatalf("second delete: expected 404, got %d", resp.StatusCode)
	}
}

func TestDelete_NotFound(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/kochbuch/api/recipes/99999", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDelete_CascadesRelations(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/kochbuch/api/recipes", sampleRecipe()))

	req, _ := http.NewRequest(http.MethodDelete,
		server.URL+"/kochbuch/api/recipes/"+itoa(created.ID), nil)
	http.DefaultClient.Do(req)

	// After delete, list should be empty
	resp, _ := http.Get(server.URL + "/kochbuch/api/recipes")
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 0 {
		t.Fatalf("expected empty list after delete, got %d", len(results))
	}
}

// --- LLM / emoji ---

func TestCreate_FillsEmojisViaLLM(t *testing.T) {
	llmSrv := mockLLMSuccess(t, `["🍝","🥚"]`)
	server := setupTestServerWithLLM(t, llmSrv.URL)
	defer server.Close()

	recipe := models.Recipe{
		Title: "Pasta",
		Ingredients: []models.Ingredient{
			{Name: "Spaghetti"},
			{Name: "Eggs"},
		},
	}
	resp := postJSON(t, server, "/kochbuch/api/recipes", recipe)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	got := decode[models.Recipe](t, resp)
	for _, ing := range got.Ingredients {
		if ing.Emoji == "" {
			t.Fatalf("ingredient %q missing emoji after LLM inference", ing.Name)
		}
	}
}

func TestCreate_LLMEmojiFailureIsGraceful(t *testing.T) {
	llmSrv := mockLLMError(t, http.StatusServiceUnavailable)
	server := setupTestServerWithLLM(t, llmSrv.URL)
	defer server.Close()

	// LLM failure must not abort the create — recipe saved without emojis.
	resp := postJSON(t, server, "/kochbuch/api/recipes", models.Recipe{
		Title:       "Pasta",
		Ingredients: []models.Ingredient{{Name: "Spaghetti"}},
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 even when LLM fails, got %d", resp.StatusCode)
	}
}

// --- helpers ---

func mockLLMSuccess(t *testing.T, content string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"content": content}},
			},
		})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func mockLLMError(t *testing.T, status int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "LLM unavailable", status)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func setupTestServerWithLLM(t *testing.T, llmURL string) *httptest.Server {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	llm := handlers.NewLLMClient(llmURL, "test-key", "test-model")
	h := handlers.NewRecipeHandler(database, llm)

	r := chi.NewRouter()
	r.Get("/kochbuch/api/recipes", h.List)
	r.Post("/kochbuch/api/recipes", h.Create)
	r.Get("/kochbuch/api/recipes/{id}", h.Get)
	r.Put("/kochbuch/api/recipes/{id}", h.Update)
	r.Delete("/kochbuch/api/recipes/{id}", h.Delete)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func itoa(id int64) string {
	return fmt.Sprintf("%d", id)
}

func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
