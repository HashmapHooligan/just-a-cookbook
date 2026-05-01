<template>
  <div class="recipe-card" :style="cardStyle" @click="emit('click')">
    <div class="recipe-card__body">
      <div class="recipe-card__tags q-mb-sm">
        <TagChip v-for="tag in recipe.tags" :key="tag.name" :tag="tag" />
      </div>
      <div class="font-headline-md recipe-card__title">{{ recipe.title }}</div>
      <div v-if="recipe.emojis?.length" class="recipe-card__emojis">
        {{ recipe.emojis.join(' ') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { RecipeSummary } from 'src/models/recipe';
import TagChip from 'src/components/TagChip.vue';

const CARD_THEMES = [
  { bg: 'var(--color-surface-container-low)', color: 'var(--color-on-surface)' },
  { bg: 'var(--color-surface-container)', color: 'var(--color-on-surface)' },
  { bg: 'var(--color-surface-container-highest)', color: 'var(--color-on-surface)' },
];

const props = defineProps<{ recipe: RecipeSummary; index: number }>();
const emit = defineEmits<{ click: [] }>();

const theme = computed(() => CARD_THEMES[props.index % CARD_THEMES.length]!);
const cardStyle = computed(() => ({
  backgroundColor: theme.value.bg,
  color: theme.value.color,
}));
</script>

<style scoped>
.recipe-card {
  height: 230px;
  border-radius: 16px;
  padding: 24px;
  overflow: hidden;
  cursor: pointer;
}

.recipe-card__body {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.recipe-card__tags {
  display: flex;
  flex-wrap: nowrap;
  gap: 6px;
  overflow: hidden;
  height: 24px;
  flex-shrink: 0;
  mask-image: linear-gradient(to right, black calc(100% - 20px), transparent 100%);
  -webkit-mask-image: linear-gradient(to right, black calc(100% - 20px), transparent 100%);
}

.recipe-card__title {
  margin-top: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.recipe-card__emojis {
  margin-top: auto;
  font-size: 1.5rem;
  line-height: 1.4;
  letter-spacing: 2px;
  white-space: nowrap;
  overflow: hidden;
  flex-shrink: 0;
}
</style>
