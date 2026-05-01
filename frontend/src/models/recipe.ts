export interface Ingredient {
  id?: number;
  name: string;
  amountNumber?: number;
  amountUnit?: string;
  emoji?: string;
  position: number;
}

export interface Step {
  id?: number;
  description: string;
  position: number;
}

export interface Tag {
  id?: number;
  name: string;
}

export interface Recipe {
  id?: number;
  title: string;
  source?: string;
  ingredients: Ingredient[];
  steps: Step[];
  tags: Tag[];
}

export interface RecipeSummary {
  id: number;
  title: string;
  tags: Tag[];
}
