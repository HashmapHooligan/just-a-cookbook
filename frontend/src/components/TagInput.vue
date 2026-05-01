<template>
  <div class="tag-input">
    <div class="tag-input__chips">
      <TagChip
        v-for="tag in modelValue"
        :key="tag.name"
        :tag="tag"
        :removable="true"
        @remove="removeTag(tag.name)"
      />
      <input
        v-model="inputVal"
        class="tag-input__field font-body-md"
        :placeholder="placeholder"
        @keydown.enter.prevent="addTag"
        @keydown.backspace="onBackspace"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Tag } from 'src/models/recipe';
import TagChip from 'src/components/TagChip.vue';

const props = defineProps<{ modelValue: Tag[]; placeholder?: string }>();
const emit = defineEmits<{ 'update:modelValue': [tags: Tag[]] }>();

const inputVal = ref('');

function addTag() {
  const name = inputVal.value.trim();
  if (!name || props.modelValue.some((t) => t.name === name)) return;
  emit('update:modelValue', [...props.modelValue, { name }]);
  inputVal.value = '';
}

function removeTag(name: string) {
  emit('update:modelValue', props.modelValue.filter((t) => t.name !== name));
}

function onBackspace() {
  if (inputVal.value === '' && props.modelValue.length > 0) {
    removeTag(props.modelValue[props.modelValue.length - 1]!.name);
  }
}
</script>

<style scoped>
.tag-input {
  border-bottom: 2px solid var(--color-outline-variant);
  background-color: var(--color-surface-container-low);
  border-radius: 8px 8px 0 0;
  padding: 12px 16px;
  min-height: 52px;
  transition: border-color 0.2s;

  &:focus-within {
    border-bottom-color: var(--color-secondary-container);
  }
}

.tag-input__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

.tag-input__field {
  border: none;
  outline: none;
  background: transparent;
  flex: 1;
  min-width: 120px;
  color: var(--color-on-surface);

  &::placeholder {
    color: var(--color-outline);
  }
}
</style>
