package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS recipes (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			title      TEXT    NOT NULL,
			source     TEXT    NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS ingredients (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id     INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
			name          TEXT    NOT NULL,
			amount_number REAL,
			amount_unit   TEXT    NOT NULL DEFAULT '',
			emoji         TEXT    NOT NULL DEFAULT '',
			position      INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS steps (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id   INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
			description TEXT    NOT NULL,
			position    INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS tags (
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT    NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS recipe_tags (
			recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
			tag_id    INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (recipe_id, tag_id)
		);

		CREATE VIRTUAL TABLE IF NOT EXISTS recipes_fts USING fts5(
			title,
			content=recipes,
			content_rowid=id
		);

		CREATE TRIGGER IF NOT EXISTS recipes_ai AFTER INSERT ON recipes BEGIN
			INSERT INTO recipes_fts(rowid, title) VALUES (new.id, new.title);
		END;

		CREATE TRIGGER IF NOT EXISTS recipes_ad AFTER DELETE ON recipes BEGIN
			INSERT INTO recipes_fts(recipes_fts, rowid, title) VALUES ('delete', old.id, old.title);
		END;

		CREATE TRIGGER IF NOT EXISTS recipes_au AFTER UPDATE ON recipes BEGIN
			INSERT INTO recipes_fts(recipes_fts, rowid, title) VALUES ('delete', old.id, old.title);
			INSERT INTO recipes_fts(rowid, title) VALUES (new.id, new.title);
		END;
	`)
	return err
}
