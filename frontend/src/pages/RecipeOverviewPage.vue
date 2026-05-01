<template>
  <q-page class="overview-page">
    <div class="overview-page__inner">
      <div class="overview-page__header">
        <h1 class="font-headline-xl overview-page__title">{{ t('overview.title') }}</h1>
        <div class="overview-page__search">
          <q-input
            v-model="searchQuery"
            :placeholder="t('overview.searchPlaceholder')"
            outlined
            dense
            clearable
            class="overview-page__search-input"
            @update:model-value="onSearch"
          >
            <template #prepend>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
      </div>

      <div v-if="store.recipes.length === 0 && !store.loading" class="overview-page__empty">
        <p class="font-body-lg" style="color: var(--color-on-surface-variant)">
          {{ t('overview.noResults') }}
        </p>
        <p class="font-body-md" style="color: var(--color-outline)">
          {{ t('overview.noResultsHint') }}
        </p>
        <q-btn
          unelevated
          :label="t('nav.newRecipe')"
          icon="add"
          to="/recipes/new"
          style="background-color: var(--color-secondary-container); color: var(--color-on-secondary)"
          class="q-mt-md font-label-lg"
        />
      </div>

      <div v-if="store.recipes.length > 0" class="overview-page__grid">
        <RecipeCard
          v-for="(recipe, i) in store.recipes"
          :key="recipe.id"
          :recipe="recipe"
          :index="i"
          @click="router.push(`/recipes/${recipe.id}`)"
        />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useRecipeStore } from 'src/stores/recipes';
import RecipeCard from 'src/components/RecipeCard.vue';

const { t } = useI18n();
const router = useRouter();
const store = useRecipeStore();

const searchQuery = ref('');

onMounted(() => store.loadRecipes());

function onSearch(val: string | number | null) {
  void store.loadRecipes(val ? String(val) : undefined);
}
</script>

<style scoped>
.overview-page {
  padding: var(--spacing-md);
}

.overview-page__inner {
  max-width: 1200px;
  margin: 0 auto;
}

.overview-page__header {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
  flex-wrap: wrap;
}

.overview-page__title {
  flex: 1;
  margin: 0;
  color: var(--color-primary-container);
}

.overview-page__search {
  flex: 1;
  min-width: 240px;
  max-width: 400px;
}

.overview-page__search-input {
  :deep(.q-field__control) {
    background-color: var(--color-surface-container-low);
    border-radius: 9999px;
  }
}

.overview-page__loading,
.overview-page__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-lg) 0;
  text-align: center;
}

.overview-page__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-gutter);
}
</style>
