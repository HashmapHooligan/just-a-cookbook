# Context
This project is called "Just a Cookbook", since it shall be exactly that: just a cookbook for my network 
- no unnecessary features
- no ads, 
- no walls of text explaining me the history of a recipe

# Planned features
- Recipe overview with fulltext search and option to delete recipe
- Recipe display with option to edit recipe
- Adding recipes manually via form
- Adding recipes by importing images and parsing them via LLM
- UI language toggle (English / German) — recipes themselves are not translated

# Architecture
As a backend language, Go is used - since I like the language for its simplicity, type safety and speed. Persistence is done with an SQLite database. LLMs are used from an OpenAI compatible API.

The application has a Vue.js frontend with Quasar and Typescript - since I know these and don't find them _to_ bad.

All cooking recipes are split into the following elements, denoted here as a JSON example:
```json
{
  "title": "Title of the recipe",
  "ingredients": [
    {
      "name": "Name of the ingredient",
      "amountNumber": 2.3,
      "amountUnit": "Unit of the amount",
      "emoji": "An emoji representing the ingredient"
    }
  ],
  "steps": [
    {
       "description": "Description of the step"
    }
  ],
  "tags": [
    {
      "name": "Name of the tag"
    }
  ],
  "source": "Source of the recipe"
}
```

# Testing
No frontend-tests - those are unnecessary for a small application like this.

For backend, test the APIs provided for the frontend, but test those thoroughly, including edge cases! No unit tests or stuff like that.

# Files
/frontend contains all frontend code. This is a Vue.js Quasar application, setup by Quasar-CLI
/backend contains all backend code
/data contains the SQLite database and recipes as JSON as a backup
/design contains rough designs created with stitch. There is funcionality displayed there, that is not necessary; stick to the "Planned features" section!
/design/the_design_system/DESIGN.md contains a description of the design system

# Rules
Only add external libraries when it is necessary. If you want to add a new library, **FIRST ASK!**

Make sure to use a _consistent_ architecture and code structure. No need to be perfect, but be _consistent_.

When is makes sense, extract utility functions and/or frontend modules in separate files.

Claude is used in sandbox mode. Thus, if you need to build or test something, let me know. I will execute the commands for you.