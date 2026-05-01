<template>
  <q-page class="detail-page" v-if="store.currentRecipe">
    <div class="detail-page__inner">
      <!-- Hero -->
      <div class="detail-page__hero">
        <div class="detail-page__tags">
          <TagChip v-for="tag in store.currentRecipe.tags" :key="tag.name" :tag="tag" />
        </div>
        <h1 class="font-headline-xl detail-page__title">{{ store.currentRecipe.title }}</h1>
        <p v-if="store.currentRecipe.source" class="font-body-md detail-page__source">
          {{ t('detail.source') }}: {{ store.currentRecipe.source }}
        </p>
      </div>

      <!-- Actions -->
      <div class="detail-page__actions">
        <q-btn
          unelevated
          :label="t('detail.edit')"
          icon="edit"
          :to="`/recipes/${store.currentRecipe.id}/edit`"
          class="font-label-lg"
          style="background-color: var(--color-primary-container); color: var(--color-on-primary)"
        />
        <q-btn
          flat
          :label="t('detail.delete')"
          icon="delete"
          class="font-label-lg"
          style="color: var(--color-negative)"
          @click="deleteDialog = true"
        />
      </div>

      <q-separator class="q-my-lg" style="background-color: var(--color-outline-variant)" />

      <!-- Ingredients -->
      <section class="detail-page__section">
        <h2 class="font-headline-lg detail-page__section-title">
          {{ t('detail.ingredients') }}
        </h2>
        <div class="detail-page__ingredients">
          <IngredientItem
            v-for="ing in store.currentRecipe.ingredients"
            :key="ing.id ?? ing.name"
            :ingredient="ing"
          />
        </div>
      </section>

      <q-separator class="q-my-lg" style="background-color: var(--color-outline-variant)" />

      <!-- Steps -->
      <section class="detail-page__section">
        <h2 class="font-headline-lg detail-page__section-title">
          {{ t('detail.steps') }}
        </h2>
        <div class="detail-page__steps">
          <StepItem
            v-for="(step, i) in store.currentRecipe.steps"
            :key="step.id ?? i"
            :step="step"
            :number="i + 1"
          />
        </div>
      </section>
    </div>

    <!-- Delete dialog -->
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

  <q-page v-else-if="store.loading" class="detail-page detail-page--loading">
    <q-spinner-dots size="64px" color="secondary" />
  </q-page>

  <q-page v-else class="detail-page detail-page--loading">
    <p class="font-body-lg" style="color: var(--color-on-surface-variant)">
      {{ t('errors.notFound') }}
    </p>
    <q-btn flat :label="t('nav.recipes')" to="/" class="q-mt-md font-label-lg" />
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import { useRecipeStore } from 'src/stores/recipes';
import TagChip from 'src/components/TagChip.vue';
import IngredientItem from 'src/components/IngredientItem.vue';
import StepItem from 'src/components/StepItem.vue';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const store = useRecipeStore();
const $q = useQuasar();

const deleteDialog = ref(false);

onMounted(() => {
  const id = Number(route.params.id);
  void store.loadRecipe(id);
});

async function doDelete() {
  if (!store.currentRecipe?.id) return;
  deleteDialog.value = false;
  try {
    await store.removeRecipe(store.currentRecipe.id);
    $q.notify({ type: 'positive', message: 'Recipe deleted' });
    void router.push('/');
  } catch {
    $q.notify({ type: 'negative', message: t('errors.deleteFailed') });
  }
}
</script>

<style scoped>
.detail-page {
  padding: var(--spacing-md);
}

.detail-page__inner {
  max-width: 800px;
  margin: 0 auto;
}

.detail-page__hero {
  padding: var(--spacing-md) 0 var(--spacing-sm);
}

.detail-page__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.detail-page__title {
  margin: 0 0 8px;
  color: var(--color-primary-container);
}

.detail-page__source {
  margin: 0;
  color: var(--color-outline);
}

.detail-page__actions {
  display: flex;
  gap: 12px;
  margin-top: var(--spacing-sm);
}

.detail-page__section-title {
  margin: 0 0 16px;
  color: var(--color-primary-container);
}

.detail-page__ingredients {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-page--loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 50vh;
}
</style>
