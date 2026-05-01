package handlers

import (
	"context"
	"database/sql"

	"justacookbook/models"
)

func insertRelations(ctx context.Context, tx *sql.Tx, recipeID int64, recipe models.Recipe) error {
	for i, ing := range recipe.Ingredients {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO ingredients (recipe_id, name, amount_number, amount_unit, emoji, position)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			recipeID, ing.Name, ing.AmountNumber, ing.AmountUnit, ing.Emoji, i,
		)
		if err != nil {
			return err
		}
	}

	for i, step := range recipe.Steps {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO steps (recipe_id, description, position) VALUES (?, ?, ?)`,
			recipeID, step.Description, i,
		)
		if err != nil {
			return err
		}
	}

	for _, tag := range recipe.Tags {
		tagID, err := upsertTag(ctx, tx, tag.Name)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO recipe_tags (recipe_id, tag_id) VALUES (?, ?)`,
			recipeID, tagID,
		); err != nil {
			return err
		}
	}
	return nil
}

func upsertTag(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	if _, err := tx.ExecContext(ctx,
		`INSERT OR IGNORE INTO tags (name) VALUES (?)`, name,
	); err != nil {
		return 0, err
	}
	var id int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM tags WHERE name=?`, name).Scan(&id)
	return id, err
}
