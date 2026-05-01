package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	r.Get("/api/recipes", h.List)
	r.Post("/api/recipes", h.Create)
	r.Get("/api/recipes/{id}", h.Get)
	r.Put("/api/recipes/{id}", h.Update)
	r.Delete("/api/recipes/{id}", h.Delete)

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

	resp, err := http.Get(server.URL + "/api/recipes")
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

	postJSON(t, server, "/api/recipes", sampleRecipe())
	postJSON(t, server, "/api/recipes", models.Recipe{Title: "Pizza", Tags: []models.Tag{}})

	resp, err := http.Get(server.URL + "/api/recipes")
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

	postJSON(t, server, "/api/recipes", sampleRecipe())
	postJSON(t, server, "/api/recipes", models.Recipe{Title: "Pizza Margherita"})

	resp, err := http.Get(server.URL + "/api/recipes?q=pasta")
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

func TestList_Search_NoMatch(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	postJSON(t, server, "/api/recipes", sampleRecipe())

	resp, err := http.Get(server.URL + "/api/recipes?q=sushi")
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

	resp := postJSON(t, server, "/api/recipes", sampleRecipe())
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

	resp := postJSON(t, server, "/api/recipes", models.Recipe{Source: "test"})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/recipes", "application/json", bytes.NewBufferString("not json"))
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

	postJSON(t, server, "/api/recipes", models.Recipe{Title: "A", Tags: []models.Tag{{Name: "Italian"}}})
	resp := postJSON(t, server, "/api/recipes", models.Recipe{Title: "B", Tags: []models.Tag{{Name: "Italian"}}})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

// --- Get ---

func TestGet_Found(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	created := decode[models.Recipe](t, postJSON(t, server, "/api/recipes", sampleRecipe()))

	resp, err := http.Get(server.URL + "/api/recipes/" + itoa(created.ID))
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

	resp, err := http.Get(server.URL + "/api/recipes/99999")
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

	resp, err := http.Get(server.URL + "/api/recipes/notanid")
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

	created := decode[models.Recipe](t, postJSON(t, server, "/api/recipes", sampleRecipe()))

	updated := created
	updated.Title = "Pasta Amatriciana"
	updated.Ingredients = []models.Ingredient{{Name: "Guanciale", AmountNumber: ptr(150.0), AmountUnit: "g"}}
	updated.Steps = []models.Step{{Description: "Fry guanciale."}}
	updated.Tags = []models.Tag{{Name: "Roman"}}

	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/api/recipes/"+itoa(created.ID),
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

func TestUpdate_NotFound(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPut,
		server.URL+"/api/recipes/99999",
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

	created := decode[models.Recipe](t, postJSON(t, server, "/api/recipes", sampleRecipe()))

	req, _ := http.NewRequest(http.MethodDelete,
		server.URL+"/api/recipes/"+itoa(created.ID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	getResp, _ := http.Get(server.URL + "/api/recipes/" + itoa(created.ID))
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getResp.StatusCode)
	}
}

func TestDelete_NotFound(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/recipes/99999", nil)
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

	created := decode[models.Recipe](t, postJSON(t, server, "/api/recipes", sampleRecipe()))

	req, _ := http.NewRequest(http.MethodDelete,
		server.URL+"/api/recipes/"+itoa(created.ID), nil)
	http.DefaultClient.Do(req)

	// After delete, list should be empty
	resp, _ := http.Get(server.URL + "/api/recipes")
	results := decode[[]models.RecipeSummary](t, resp)
	if len(results) != 0 {
		t.Fatalf("expected empty list after delete, got %d", len(results))
	}
}

// helpers

func itoa(id int64) string {
	return fmt.Sprintf("%d", id)
}

func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
