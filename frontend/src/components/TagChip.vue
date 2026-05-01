<template>
  <span class="tag-chip font-label-md">
    {{ tag.name }}
    <button
      v-if="removable"
      class="tag-chip__remove"
      @click.stop="emit('remove')"
      :aria-label="`Remove ${tag.name}`"
    >
      ✕
    </button>
  </span>
</template>

<script setup lang="ts">
import type { Tag } from 'src/models/recipe';

defineProps<{ tag: Tag; removable?: boolean }>();
const emit = defineEmits<{ remove: [] }>();
</script>

<style scoped>
.tag-chip__remove {
  background: none;
  border: none;
  cursor: pointer;
  color: inherit;
  opacity: 0.7;
  padding: 0 0 0 2px;
  font-size: 10px;
  line-height: 1;

  &:hover {
    opacity: 1;
  }
}
</style>
