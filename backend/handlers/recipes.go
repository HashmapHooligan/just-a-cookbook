package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"justacookbook/models"
)

type RecipeHandler struct {
	db  *sql.DB
	llm *LLMClient
}

func NewRecipeHandler(db *sql.DB, llm *LLMClient) *RecipeHandler {
	return &RecipeHandler{db: db, llm: llm}
}

func (h *RecipeHandler) fillEmojis(r *http.Request, recipe *models.Recipe) {
	if h.llm == nil {
		return
	}
	var indices []int
	var names []string
	for i, ing := range recipe.Ingredients {
		if ing.Emoji == "" && ing.Name != "" {
			indices = append(indices, i)
			names = append(names, ing.Name)
		}
	}
	if len(names) == 0 {
		return
	}
	emojis, err := h.llm.inferEmojis(r.Context(), names)
	if err != nil {
		log.Printf("fillEmojis: failed: %v", err)
		return
	}
	for j, i := range indices {
		recipe.Ingredients[i].Emoji = emojis[j]
	}
}

func (h *RecipeHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var rows *sql.Rows
	var err error

	if query != "" {
		rows, err = h.db.QueryContext(r.Context(), `
			SELECT r.id, r.title
			FROM recipes r
			JOIN recipes_fts fts ON fts.rowid = r.id
			WHERE recipes_fts MATCH ?
			ORDER BY rank
		`, ftsPrefix(query))
	} else {
		rows, err = h.db.QueryContext(r.Context(), `
			SELECT id, title FROM recipes ORDER BY created_at DESC
		`)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}
	defer rows.Close()

	summaries := make([]models.RecipeSummary, 0)
	for rows.Next() {
		var s models.RecipeSummary
		if err := rows.Scan(&s.ID, &s.Title); err != nil {
			writeError(w, http.StatusInternalServerError, "scan failed")
			return
		}
		tags, err := h.loadTags(r, s.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "load tags failed")
			return
		}
		s.Tags = tags
		emojis, err := h.loadIngredientEmojis(r, s.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "load emojis failed")
			return
		}
		s.Emojis = emojis
		summaries = append(summaries, s)
	}
	writeJSON(w, http.StatusOK, summaries)
}

func (h *RecipeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	recipe, err := h.loadRecipe(r, id)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "recipe not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load failed")
		return
	}
	writeJSON(w, http.StatusOK, recipe)
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if recipe.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	h.fillEmojis(r, &recipe)

	id, err := h.insertRecipe(r, recipe)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "insert failed")
		return
	}

	created, err := h.loadRecipe(r, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load failed")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *RecipeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if recipe.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	h.fillEmojis(r, &recipe)

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "tx failed")
		return
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(r.Context(),
		`UPDATE recipes SET title=?, source=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		recipe.Title, recipe.Source, id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(w, http.StatusNotFound, "recipe not found")
		return
	}

	if _, err := tx.ExecContext(r.Context(), `DELETE FROM ingredients WHERE recipe_id=?`, id); err != nil {
		writeError(w, http.StatusInternalServerError, "delete ingredients failed")
		return
	}
	if _, err := tx.ExecContext(r.Context(), `DELETE FROM steps WHERE recipe_id=?`, id); err != nil {
		writeError(w, http.StatusInternalServerError, "delete steps failed")
		return
	}
	if _, err := tx.ExecContext(r.Context(), `DELETE FROM recipe_tags WHERE recipe_id=?`, id); err != nil {
		writeError(w, http.StatusInternalServerError, "delete recipe_tags failed")
		return
	}

	if err := insertRelations(r.Context(), tx, id, recipe); err != nil {
		writeError(w, http.StatusInternalServerError, "insert relations failed")
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "commit failed")
		return
	}

	updated, err := h.loadRecipe(r, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load failed")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *RecipeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := h.db.ExecContext(r.Context(), `DELETE FROM recipes WHERE id=?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "delete failed")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(w, http.StatusNotFound, "recipe not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) loadRecipe(r *http.Request, id int64) (*models.Recipe, error) {
	recipe := &models.Recipe{}
	err := h.db.QueryRowContext(r.Context(),
		`SELECT id, title, source FROM recipes WHERE id=?`, id,
	).Scan(&recipe.ID, &recipe.Title, &recipe.Source)
	if err != nil {
		return nil, err
	}

	rows, err := h.db.QueryContext(r.Context(),
		`SELECT id, name, amount_number, amount_unit, emoji, position
		 FROM ingredients WHERE recipe_id=? ORDER BY position`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	recipe.Ingredients = make([]models.Ingredient, 0)
	for rows.Next() {
		var ing models.Ingredient
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.AmountNumber, &ing.AmountUnit, &ing.Emoji, &ing.Position); err != nil {
			return nil, err
		}
		recipe.Ingredients = append(recipe.Ingredients, ing)
	}

	stepRows, err := h.db.QueryContext(r.Context(),
		`SELECT id, description, position FROM steps WHERE recipe_id=? ORDER BY position`, id,
	)
	if err != nil {
		return nil, err
	}
	defer stepRows.Close()
	recipe.Steps = make([]models.Step, 0)
	for stepRows.Next() {
		var s models.Step
		if err := stepRows.Scan(&s.ID, &s.Description, &s.Position); err != nil {
			return nil, err
		}
		recipe.Steps = append(recipe.Steps, s)
	}

	tags, err := h.loadTags(r, id)
	if err != nil {
		return nil, err
	}
	recipe.Tags = tags
	return recipe, nil
}

func (h *RecipeHandler) loadIngredientEmojis(r *http.Request, recipeID int64) ([]string, error) {
	rows, err := h.db.QueryContext(r.Context(),
		`SELECT emoji FROM ingredients WHERE recipe_id=? AND emoji != '' ORDER BY position`, recipeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	emojis := make([]string, 0)
	for rows.Next() {
		var emoji string
		if err := rows.Scan(&emoji); err != nil {
			return nil, err
		}
		emojis = append(emojis, emoji)
	}
	return emojis, nil
}

func (h *RecipeHandler) loadTags(r *http.Request, recipeID int64) ([]models.Tag, error) {
	rows, err := h.db.QueryContext(r.Context(),
		`SELECT t.id, t.name FROM tags t
		 JOIN recipe_tags rt ON rt.tag_id=t.id
		 WHERE rt.recipe_id=? ORDER BY t.name`, recipeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]models.Tag, 0)
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (h *RecipeHandler) insertRecipe(r *http.Request, recipe models.Recipe) (int64, error) {
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(r.Context(),
		`INSERT INTO recipes (title, source) VALUES (?, ?)`,
		recipe.Title, recipe.Source,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := insertRelations(r.Context(), tx, id, recipe); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// ftsPrefix wraps q in FTS5 double-quote syntax and appends * for prefix matching.
// Internal double-quotes are doubled ("" is the FTS5 escape sequence) so user input
// cannot inject FTS5 operators like OR, AND, NEAR, or column filters.
func ftsPrefix(q string) string {
	escaped := strings.ReplaceAll(q, `"`, `""`)
	return `"` + escaped + `"*`
}
