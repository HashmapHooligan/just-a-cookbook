<template>
  <q-page class="form-page">
    <div class="form-page__inner">
      <div class="form-page__header">
        <q-btn flat round icon="close" @click="router.back()" />
        <h1 class="font-headline-md form-page__header-title">{{ t('bulkAdd.title') }}</h1>
        <q-btn
          unelevated
          :label="t('bulkAdd.importButton')"
          :loading="running"
          :disable="files.length === 0"
          class="font-label-lg"
          style="background-color: var(--color-primary-container); color: var(--color-on-primary)"
          @click="runImport"
        />
      </div>

      <div
        class="bulk-drop"
        :class="{ 'bulk-drop--dragging': dragging }"
        @dragover.prevent="dragging = true"
        @dragleave="dragging = false"
        @drop.prevent="onDrop"
        @click="fileInput?.click()"
      >
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          multiple
          class="bulk-drop__hidden-input"
          @change="onFileChange"
        />
        <span class="material-symbol" style="font-size: 48px; color: var(--color-outline)">
          photo_library
        </span>
        <p class="font-body-md" style="color: var(--color-on-surface-variant); margin: 8px 0 0">
          {{ t('bulkAdd.hint') }}
        </p>
      </div>

      <ul v-if="files.length > 0" class="bulk-file-list">
        <li v-for="(file, i) in files" :key="i" class="font-body-md bulk-file-list__item">
          <span class="material-symbol bulk-file-list__icon" style="font-size: 18px">image</span>
          {{ file.name }}
        </li>
      </ul>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { importRecipeFromImage, createRecipe } from 'src/api/recipes';

const router = useRouter();
const $q = useQuasar();
const { t } = useI18n();

const fileInput = ref<HTMLInputElement | null>(null);
const dragging = ref(false);
const files = ref<File[]>([]);
const running = ref(false);

function onDrop(e: DragEvent) {
  dragging.value = false;
  const dropped = Array.from(e.dataTransfer?.files ?? []).filter((f) =>
    f.type.startsWith('image/'),
  );
  if (dropped.length > 0) files.value = dropped;
}

function onFileChange(e: Event) {
  const selected = Array.from((e.target as HTMLInputElement).files ?? []);
  if (selected.length > 0) files.value = selected;
}

async function runLimited(tasks: (() => Promise<void>)[], limit: number): Promise<void> {
  const queue = [...tasks];
  async function worker() {
    while (queue.length > 0) {
      await queue.shift()!();
    }
  }
  await Promise.all(Array.from({ length: Math.min(limit, tasks.length) }, worker));
}

async function runImport() {
  running.value = true;
  const total = files.value.length;
  let succeeded = 0;
  let failed = 0;

  const tasks = files.value.map((file) => async () => {
    try {
      const recipe = await importRecipeFromImage(file);
      await createRecipe(recipe);
      succeeded++;
    } catch {
      failed++;
    }
  });

  await runLimited(tasks, 5);
  running.value = false;

  let message: string;
  if (failed === 0) {
    message = t('bulkAdd.success', { n: succeeded });
  } else if (succeeded === 0) {
    message = t('bulkAdd.allFailed');
  } else {
    message = t('bulkAdd.partial', { success: succeeded, total, failed });
  }

  $q.notify({ message, type: failed === 0 ? 'positive' : 'warning' });
  void router.push('/');
}
</script>

<style scoped>
.bulk-drop {
  border: 2px dashed var(--color-outline-variant);
  border-radius: 12px;
  padding: 40px 24px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s, background-color 0.2s;
  background-color: var(--color-surface-container-low);

  &:hover {
    border-color: var(--color-secondary-container);
    background-color: var(--color-surface-container);
  }
}

.bulk-drop--dragging {
  border-color: var(--color-secondary-container);
  background-color: var(--color-surface-container);
}

.bulk-drop__hidden-input {
  display: none;
}

.bulk-file-list {
  list-style: none;
  margin: 16px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.bulk-file-list__item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-on-surface-variant);
}

.bulk-file-list__icon {
  color: var(--color-outline);
}
</style>
