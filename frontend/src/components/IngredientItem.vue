<template>
  <label class="ingredient-item" :class="{ 'ingredient-item--checked': checked }">
    <input type="checkbox" v-model="checked" class="ingredient-item__checkbox" />
    <span class="ingredient-item__emoji" v-if="ingredient.emoji">{{ ingredient.emoji }}</span>
    <span class="ingredient-item__name font-body-md" :class="{ 'ingredient-item__name--done': checked }">
      {{ ingredient.name }}
    </span>
    <span class="ingredient-item__amount font-body-sm" v-if="ingredient.amountNumber || ingredient.amountUnit">
      {{ ingredient.amountNumber ? formatAmount(ingredient.amountNumber) : '' }}
      {{ ingredient.amountUnit ?? '' }}
    </span>
  </label>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Ingredient } from 'src/models/recipe';

defineProps<{ ingredient: Ingredient }>();

const checked = ref(false);

function formatAmount(n: number): string {
  return n % 1 === 0 ? String(n) : n.toFixed(1).replace('.0', '');
}
</script>

<style scoped>
.ingredient-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.15s;

  &:hover {
    background-color: var(--color-surface-variant);
    opacity: 0.6;
  }
}

.ingredient-item__checkbox {
  width: 20px;
  height: 20px;
  accent-color: var(--color-secondary-container);
  cursor: pointer;
  flex-shrink: 0;
}

.ingredient-item__emoji {
  font-size: 20px;
  flex-shrink: 0;
}

.ingredient-item__name {
  flex: 1;
  color: var(--color-on-surface);
  transition: opacity 0.2s;
}

.ingredient-item__name--done {
  text-decoration: line-through;
  opacity: 0.5;
}

.ingredient-item__amount {
  color: var(--color-outline);
  white-space: nowrap;
}
</style>
