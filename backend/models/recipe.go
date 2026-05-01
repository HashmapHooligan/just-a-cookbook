package models

type Ingredient struct {
	ID           int64    `json:"id,omitempty"`
	RecipeID     int64    `json:"-"`
	Name         string   `json:"name"`
	AmountNumber *float64 `json:"amountNumber,omitempty"`
	AmountUnit   string   `json:"amountUnit,omitempty"`
	Emoji        string   `json:"emoji,omitempty"`
	Position     int      `json:"position"`
}

type Step struct {
	ID          int64  `json:"id,omitempty"`
	RecipeID    int64  `json:"-"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}

type Tag struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
}

type Recipe struct {
	ID          int64        `json:"id,omitempty"`
	Title       string       `json:"title"`
	Source      string       `json:"source,omitempty"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []Step       `json:"steps"`
	Tags        []Tag        `json:"tags"`
}

type RecipeSummary struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Tags  []Tag  `json:"tags"`
}
