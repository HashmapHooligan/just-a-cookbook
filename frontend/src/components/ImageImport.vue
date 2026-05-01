<template>
  <div
    class="image-import"
    :class="{ 'image-import--dragging': dragging }"
    @dragover.prevent="dragging = true"
    @dragleave="dragging = false"
    @drop.prevent="onDrop"
    @click="fileInput?.click()"
  >
    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      class="image-import__hidden-input"
      @change="onFileChange"
    />

    <div v-if="!loading" class="image-import__content">
      <span class="material-symbol" style="font-size: 48px; color: var(--color-outline)">image</span>
      <p class="font-body-md" style="color: var(--color-on-surface-variant); margin: 8px 0 0">
        {{ t('form.importHint') }}
      </p>
    </div>

    <div v-else class="image-import__content">
      <q-spinner-dots size="48px" color="secondary" />
      <p class="font-body-md" style="color: var(--color-on-surface-variant); margin: 8px 0 0">
        {{ t('form.importing') }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const emit = defineEmits<{ import: [file: File] }>();

const { t } = useI18n();
const fileInput = ref<HTMLInputElement | null>(null);
const dragging = ref(false);
const loading = defineModel<boolean>('loading', { default: false });

function onDrop(e: DragEvent) {
  dragging.value = false;
  const file = e.dataTransfer?.files[0];
  if (file && file.type.startsWith('image/')) {
    emit('import', file);
  }
}

function onFileChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (file) emit('import', file);
}
</script>

<style scoped>
.image-import {
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

.image-import--dragging {
  border-color: var(--color-secondary-container);
  background-color: var(--color-surface-container);
}

.image-import__content {
  display: flex;
  flex-direction: column;
  align-items: center;
  pointer-events: none;
}

.image-import__hidden-input {
  display: none;
}
</style>
