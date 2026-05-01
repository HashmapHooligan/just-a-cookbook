package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"justacookbook/db"
	"justacookbook/env"
	"justacookbook/handlers"
	"justacookbook/middleware"
)

func main() {
	if err := env.Load(".env"); err != nil {
		log.Printf("warning: .env load failed: %v", err)
	}

	dbPath := getenv("DB_PATH", "../data/cookbook.db")
	llmURL := getenv("LLM_API_URL", "https://api.openai.com/v1")
	llmKey := getenv("LLM_API_KEY", "")
	llmModel := getenv("LLM_MODEL", "gpt-4o")
	addr := getenv("ADDR", ":8080")
	originsRaw := getenv("ALLOWED_ORIGINS", "http://localhost:9000")

	allowedOrigins := splitTrim(originsRaw)

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	llmClient := handlers.NewLLMClient(llmURL, llmKey, llmModel)
	recipeHandler := handlers.NewRecipeHandler(database, llmClient)
	importHandler := handlers.NewImportHandler(llmClient)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS(allowedOrigins))

	r.Route("/api/recipes", func(r chi.Router) {
		r.Get("/", recipeHandler.List)
		r.Post("/", recipeHandler.Create)
		r.Post("/import", importHandler.Import)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", recipeHandler.Get)
			r.Put("/", recipeHandler.Update)
			r.Delete("/", recipeHandler.Delete)
		})
	})

	log.Printf("listening on %s (allowed origins: %v)", addr, allowedOrigins)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("listen: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
