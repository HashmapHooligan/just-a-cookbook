# Just a Cookbook

A personal cookbook for my network. No ads. No fluff. No walls of text about the history of pasta.

## Features

- Recipe overview with fulltext search
- Recipe display and editing
- Add recipes manually via form
- Import recipes from images using LLM parsing
- UI language toggle (EN / DE)

## Stack

| Layer    | Tech                         |
|----------|------------------------------|
| Backend  | Go + SQLite                  |
| Frontend | Vue.js + Quasar + TypeScript |
| LLM      | OpenAI-compatible API        |

## Structure

```
/frontend   — Vue.js Quasar app
/backend    — Go API server
/design     — Design references and design system docs
/data       — SQLite database + JSON recipe backups
```

## Running

Start backend:
```bash
cd backend && go run .
```

Start frontend:
```bash
cd frontend && quasar dev
```
