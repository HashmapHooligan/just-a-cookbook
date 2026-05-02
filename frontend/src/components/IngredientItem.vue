<template>
  <div class="ingredient-item">
    <span class="ingredient-item__emoji" v-if="ingredient.emoji">{{ ingredient.emoji }}</span>
    <span class="ingredient-item__name font-body-md">
      {{ ingredient.name }}
    </span>
    <span class="ingredient-item__amount font-body-sm" v-if="ingredient.amountNumber || ingredient.amountUnit">
      {{ ingredient.amountNumber ? formatAmount(ingredient.amountNumber) : '' }}
      {{ ingredient.amountUnit ?? '' }}
    </span>
  </div>
</template>

<script setup lang="ts">
import type { Ingredient } from 'src/models/recipe';

defineProps<{ ingredient: Ingredient }>();

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
}

.ingredient-item__emoji {
  font-size: 20px;
  flex-shrink: 0;
}

.ingredient-item__name {
  flex: 1;
  color: var(--color-on-surface);
}

.ingredient-item__amount {
  color: var(--color-outline);
  white-space: nowrap;
}
</style>
