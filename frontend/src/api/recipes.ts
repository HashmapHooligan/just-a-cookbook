import type { Recipe, RecipeSummary } from 'src/models/recipe';

const BASE = '/api/recipes';

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const resp = await fetch(url, options);
  if (!resp.ok) {
    const body = await resp.json().catch(() => ({ error: resp.statusText }));
    throw new Error(body.error ?? resp.statusText);
  }
  if (resp.status === 204) return undefined as T;
  return resp.json();
}

export function fetchRecipes(query?: string): Promise<RecipeSummary[]> {
  const url = query ? `${BASE}?q=${encodeURIComponent(query)}` : BASE;
  return request<RecipeSummary[]>(url);
}

export function fetchRecipe(id: number): Promise<Recipe> {
  return request<Recipe>(`${BASE}/${id}`);
}

export function createRecipe(recipe: Recipe): Promise<Recipe> {
  return request<Recipe>(BASE, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(recipe),
  });
}

export function updateRecipe(id: number, recipe: Recipe): Promise<Recipe> {
  return request<Recipe>(`${BASE}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(recipe),
  });
}

export function deleteRecipe(id: number): Promise<void> {
  return request<void>(`${BASE}/${id}`, { method: 'DELETE' });
}

export function importRecipeFromImage(file: File): Promise<Recipe> {
  const form = new FormData();
  form.append('image', file);
  return request<Recipe>(`${BASE}/import`, { method: 'POST', body: form });
}
