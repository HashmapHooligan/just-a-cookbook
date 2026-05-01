<template>
  <q-page class="form-page">
    <div class="form-page__inner">
      <!-- Header -->
      <div class="form-page__header">
        <q-btn flat round icon="close" @click="router.back()" />
        <h1 class="font-headline-md form-page__header-title">
          {{ isEdit ? t('form.editTitle') : t('form.addTitle') }}
        </h1>
        <q-btn
          unelevated
          :label="t('form.save')"
          :loading="store.loading"
          class="font-label-lg"
          style="background-color: var(--color-primary-container); color: var(--color-on-primary)"
          @click="submit"
        />
      </div>

      <!-- Image Import (only on add mode) -->
      <section v-if="!isEdit" class="form-page__section">
        <h2 class="font-headline-md form-page__section-title">{{ t('form.importTitle') }}</h2>
        <ImageImport v-model:loading="importing" @import="onImport" />
        <p v-if="importError" class="font-body-sm form-page__import-error">
          {{ t('form.importError') }}
        </p>
      </section>

      <q-separator v-if="!isEdit" class="q-my-lg" style="background-color: var(--color-outline-variant)" />

      <!-- Title -->
      <div class="form-page__field">
        <label class="form-page__label font-label-lg">
          {{ t('form.title') }} <span style="color: var(--color-secondary-container)">*</span>
        </label>
        <input
          v-model="form.title"
          class="form-page__input font-body-md"
          :placeholder="t('form.title')"
          :class="{ 'form-page__input--error': titleError }"
        />
        <p v-if="titleError" class="form-page__field-error font-body-sm">
          {{ t('form.titleRequired') }}
        </p>
      </div>

      <!-- Source -->
      <div class="form-page__field">
        <label class="form-page__label font-label-lg">{{ t('form.source') }}</label>
        <input
          v-model="form.source"
          class="form-page__input font-body-md"
          :placeholder="t('form.source')"
        />
      </div>

      <!-- Tags -->
      <div class="form-page__field">
        <label class="form-page__label font-label-lg">{{ t('form.tags') }}</label>
        <TagInput v-model="form.tags" :placeholder="t('form.tagsHint')" />
      </div>

      <q-separator class="q-my-lg" style="background-color: var(--color-outline-variant)" />

      <!-- Ingredients -->
      <section class="form-page__section">
        <h2 class="font-headline-md form-page__section-title">{{ t('form.ingredients') }}</h2>
        <div
          v-for="(ing, i) in form.ingredients"
          :key="i"
          class="form-page__ingredient-row"
        >
          <input
            v-model="ing.name"
            class="form-page__input form-page__input--name font-body-md"
            :placeholder="t('form.ingredientName')"
            @keydown.enter.prevent="addIngredient"
          />
          <input
            v-model.number="ing.amountNumber"
            type="number"
            min="0"
            step="any"
            class="form-page__input form-page__input--amount font-body-md"
            :placeholder="t('form.ingredientAmount')"
            @keydown.enter.prevent="addIngredient"
          />
          <input
            v-model="ing.amountUnit"
            class="form-page__input form-page__input--unit font-body-md"
            :placeholder="t('form.ingredientUnit')"
            @keydown.enter.prevent="addIngredient"
          />
          <q-btn
            flat
            round
            dense
            icon="delete"
            style="color: var(--color-outline)"
            @click="removeIngredient(i)"
          />
        </div>
        <q-btn
          flat
          no-caps
          :label="t('form.addIngredient')"
          icon="add"
          class="font-label-lg q-mt-sm"
          style="color: var(--color-secondary-container)"
          @click="addIngredient"
        />
      </section>

      <q-separator class="q-my-lg" style="background-color: var(--color-outline-variant)" />

      <!-- Steps -->
      <section class="form-page__section">
        <h2 class="font-headline-md form-page__section-title">{{ t('form.steps') }}</h2>
        <div v-for="(step, i) in form.steps" :key="i" class="form-page__step-row">
          <div class="step-number">{{ i + 1 }}</div>
          <textarea
            v-model="step.description"
            class="form-page__textarea font-body-md"
            :placeholder="t('form.stepDescription')"
            rows="2"
            @keydown.enter.exact.prevent="addStep"
          />
          <q-btn
            flat
            round
            dense
            icon="delete"
            style="color: var(--color-outline)"
            @click="removeStep(i)"
          />
        </div>
        <q-btn
          flat
          no-caps
          :label="t('form.addStep')"
          icon="add"
          class="font-label-lg q-mt-sm"
          style="color: var(--color-secondary-container)"
          @click="addStep"
        />
      </section>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import { useRecipeStore } from 'src/stores/recipes';
import type { Recipe, Ingredient, Step } from 'src/models/recipe';
import TagInput from 'src/components/TagInput.vue';
import ImageImport from 'src/components/ImageImport.vue';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const store = useRecipeStore();
const $q = useQuasar();

const isEdit = computed(() => route.path.endsWith('/edit'));
const titleError = ref(false);
const importing = ref(false);
const importError = ref(false);

function emptyRecipe(): Recipe {
  return { title: '', source: '', ingredients: [], steps: [], tags: [] };
}

const form = ref<Recipe>(emptyRecipe());

onMounted(async () => {
  if (isEdit.value) {
    const id = Number(route.params.id);
    await store.loadRecipe(id);
    if (store.currentRecipe) {
      form.value = JSON.parse(JSON.stringify(store.currentRecipe)) as Recipe;
    }
  }
});

async function onImport(file: File) {
  importing.value = true;
  importError.value = false;
  try {
    const imported = await store.importFromImage(file);
    form.value = {
      ...imported,
      ingredients: imported.ingredients.map((ing, i) => ({ ...ing, position: i })),
      steps: imported.steps.map((s, i) => ({ ...s, position: i })),
      tags: imported.tags ?? [],
    };
    $q.notify({ type: 'positive', message: t('form.importSuccess') });
  } catch {
    importError.value = true;
    $q.notify({ type: 'negative', message: t('form.importError') });
  } finally {
    importing.value = false;
  }
}

function addIngredient() {
  form.value.ingredients.push({ name: '', position: form.value.ingredients.length });
}

function removeIngredient(i: number) {
  form.value.ingredients.splice(i, 1);
}

function addStep() {
  form.value.steps.push({ description: '', position: form.value.steps.length });
}

function removeStep(i: number) {
  form.value.steps.splice(i, 1);
}

async function submit() {
  if (!form.value.title.trim()) {
    titleError.value = true;
    return;
  }
  titleError.value = false;

  const payload: Recipe = {
    ...form.value,
    ingredients: form.value.ingredients
      .filter((ing) => ing.name.trim())
      .map((ing, i): Ingredient => ({
        name: ing.name,
        position: i,
        ...(ing.id !== undefined && { id: ing.id }),
        ...(ing.amountNumber && { amountNumber: ing.amountNumber }),
        ...(ing.amountUnit && { amountUnit: ing.amountUnit }),
        ...(ing.emoji && { emoji: ing.emoji }),
      })),
    steps: form.value.steps
      .filter((s) => s.description.trim())
      .map((s, i): Step => ({ ...s, position: i })),
  };

  try {
    const saved = await store.saveRecipe(payload);
    $q.notify({ type: 'positive', message: 'Recipe saved!' });
    void router.push(`/recipes/${saved.id}`);
  } catch {
    $q.notify({ type: 'negative', message: t('errors.saveFailed') });
  }
}
</script>

<style scoped>
.form-page {
  padding: var(--spacing-md);
  padding-bottom: var(--spacing-xl);
}

.form-page__inner {
  max-width: 800px;
  margin: 0 auto;
}

.form-page__header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: var(--spacing-md);
  position: sticky;
  top: 50px;
  z-index: 10;
  background-color: var(--color-background);
  padding: 12px 0;
  border-bottom: 1px solid var(--color-outline-variant);
}

.form-page__header-title {
  flex: 1;
  margin: 0;
  color: var(--color-primary-container);
}

.form-page__section {
  margin-bottom: var(--spacing-md);
}

.form-page__section-title {
  margin: 0 0 16px;
  color: var(--color-primary-container);
}

.form-page__field {
  margin-bottom: var(--spacing-sm);
}

.form-page__label {
  display: block;
  color: var(--color-on-surface-variant);
  margin-bottom: 6px;
}

.form-page__input {
  width: 100%;
  box-sizing: border-box;
  border: none;
  border-bottom: 2px solid var(--color-outline-variant);
  background-color: var(--color-surface-container-low);
  border-radius: 8px 8px 0 0;
  padding: 12px 16px;
  outline: none;
  color: var(--color-on-surface);
  transition: border-bottom-color 0.2s;

  &:focus {
    border-bottom-color: var(--color-secondary-container);
  }

  &::placeholder {
    color: var(--color-outline);
  }
}

.form-page__input--error {
  border-bottom-color: var(--color-negative);
}

.form-page__field-error,
.form-page__import-error {
  color: var(--color-negative);
  margin: 4px 0 0;
}

.form-page__input--name {
  flex: 1;
}

.form-page__input--amount {
  width: 80px;
}

.form-page__input--unit {
  width: 80px;
}

.form-page__ingredient-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.form-page__step-row {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  margin-bottom: 12px;
}

.form-page__textarea {
  flex: 1;
  border: none;
  border-bottom: 2px solid var(--color-outline-variant);
  background-color: var(--color-surface-container-low);
  border-radius: 8px 8px 0 0;
  padding: 12px 16px;
  outline: none;
  color: var(--color-on-surface);
  resize: vertical;
  font-family: 'Be Vietnam Pro', sans-serif;
  transition: border-bottom-color 0.2s;

  &:focus {
    border-bottom-color: var(--color-secondary-container);
  }

  &::placeholder {
    color: var(--color-outline);
  }
}
</style>
