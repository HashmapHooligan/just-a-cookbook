import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Recipe, RecipeSummary } from 'src/models/recipe';
import * as api from 'src/api/recipes';

export const useRecipeStore = defineStore('recipes', () => {
  const recipes = ref<RecipeSummary[]>([]);
  const currentRecipe = ref<Recipe | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function loadRecipes(query?: string) {
    loading.value = true;
    error.value = null;
    try {
      recipes.value = await api.fetchRecipes(query);
    } catch (e) {
      error.value = (e as Error).message;
    } finally {
      loading.value = false;
    }
  }

  async function loadRecipe(id: number) {
    loading.value = true;
    error.value = null;
    try {
      currentRecipe.value = await api.fetchRecipe(id);
    } catch (e) {
      error.value = (e as Error).message;
      currentRecipe.value = null;
    } finally {
      loading.value = false;
    }
  }

  async function saveRecipe(recipe: Recipe): Promise<Recipe> {
    loading.value = true;
    error.value = null;
    try {
      const saved = recipe.id
        ? await api.updateRecipe(recipe.id, recipe)
        : await api.createRecipe(recipe);
      return saved;
    } catch (e) {
      error.value = (e as Error).message;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  async function removeRecipe(id: number) {
    loading.value = true;
    error.value = null;
    try {
      await api.deleteRecipe(id);
      recipes.value = recipes.value.filter((r) => r.id !== id);
    } catch (e) {
      error.value = (e as Error).message;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  async function importFromImage(file: File): Promise<Recipe> {
    loading.value = true;
    error.value = null;
    try {
      return await api.importRecipeFromImage(file);
    } catch (e) {
      error.value = (e as Error).message;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  return {
    recipes,
    currentRecipe,
    loading,
    error,
    loadRecipes,
    loadRecipe,
    saveRecipe,
    removeRecipe,
    importFromImage,
  };
});
