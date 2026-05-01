<template>
  <div class="recipe-card" :style="cardStyle" @click="emit('click')">
    <div class="recipe-card__body">
      <div class="recipe-card__tags q-mb-sm">
        <TagChip v-for="tag in recipe.tags.slice(0, 3)" :key="tag.name" :tag="tag" />
      </div>
      <div class="font-headline-md recipe-card__title">{{ recipe.title }}</div>
    </div>
    <div class="recipe-card__actions">
      <q-btn
        flat
        round
        dense
        icon="delete"
        :style="{ color: deleteColor }"
        @click.stop="emit('delete')"
        :aria-label="`Delete ${recipe.title}`"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { RecipeSummary } from 'src/models/recipe';
import TagChip from 'src/components/TagChip.vue';

const CARD_THEMES = [
  { bg: 'var(--color-surface-container-low)', color: 'var(--color-on-surface)', deleteFg: 'var(--color-secondary-container)' },
  { bg: 'var(--color-primary-container)', color: 'var(--color-on-primary)', deleteFg: 'var(--color-primary-fixed-dim)' },
  { bg: 'var(--color-surface-container)', color: 'var(--color-on-surface)', deleteFg: 'var(--color-secondary-container)' },
  { bg: 'var(--color-surface-container-highest)', color: 'var(--color-on-surface)', deleteFg: 'var(--color-secondary-container)' },
];

const props = defineProps<{ recipe: RecipeSummary; index: number }>();
const emit = defineEmits<{ click: []; delete: [] }>();

const theme = computed(() => CARD_THEMES[props.index % CARD_THEMES.length]!);
const cardStyle = computed(() => ({
  backgroundColor: theme.value.bg,
  color: theme.value.color,
}));
const deleteColor = computed(() => theme.value.deleteFg);
</script>

<style scoped>
.recipe-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 200px;
  border-radius: 16px;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

.recipe-card__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.recipe-card__title {
  margin-top: 8px;
}

.recipe-card__actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
</style>
