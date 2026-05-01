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

      <div v-if="store.loading" class="overview-page__loading">
        <q-spinner-dots size="48px" color="secondary" />
      </div>

      <div v-else-if="store.recipes.length === 0" class="overview-page__empty">
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

      <div v-else class="overview-page__grid">
        <RecipeCard
          v-for="(recipe, i) in store.recipes"
          :key="recipe.id"
          :recipe="recipe"
          :index="i"
          @click="router.push(`/recipes/${recipe.id}`)"
          @delete="confirmDelete(recipe)"
        />
      </div>
    </div>

    <q-dialog v-model="deleteDialog">
      <q-card style="min-width: 320px; border-radius: 16px;">
        <q-card-section>
          <div class="font-headline-md">{{ t('detail.confirmDelete') }}</div>
          <p class="font-body-md q-mt-sm" style="color: var(--color-on-surface-variant)">
            {{ t('detail.confirmDeleteMessage') }}
          </p>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('confirm.no')" v-close-popup class="font-label-lg" />
          <q-btn
            unelevated
            :label="t('confirm.yes')"
            class="font-label-lg"
            style="background-color: var(--color-negative); color: white"
            @click="doDelete"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import { useRecipeStore } from 'src/stores/recipes';
import type { RecipeSummary } from 'src/models/recipe';
import RecipeCard from 'src/components/RecipeCard.vue';

const { t } = useI18n();
const router = useRouter();
const store = useRecipeStore();
const $q = useQuasar();

const searchQuery = ref('');
const deleteDialog = ref(false);
const pendingDelete = ref<RecipeSummary | null>(null);

let searchTimer: ReturnType<typeof setTimeout>;

onMounted(() => store.loadRecipes());

function onSearch(val: string | number | null) {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    void store.loadRecipes(val ? String(val) : undefined);
  }, 300);
}

function confirmDelete(recipe: RecipeSummary) {
  pendingDelete.value = recipe;
  deleteDialog.value = true;
}

async function doDelete() {
  if (!pendingDelete.value) return;
  deleteDialog.value = false;
  try {
    await store.removeRecipe(pendingDelete.value.id);
    $q.notify({ type: 'positive', message: `"${pendingDelete.value.title}" deleted` });
  } catch {
    $q.notify({ type: 'negative', message: t('errors.deleteFailed') });
  }
  pendingDelete.value = null;
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
