<template>
  <div class="just-a-logo">
    <q-chip
      :style="{ backgroundColor: chipFill, color: chipText }"
      class="just-a-chip"
      :ripple="false"
    >JUST A</q-chip>
    <span class="just-a-wordmark" :style="{ color: wordmarkFill }">Cookbook</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ dark: boolean }>();

const chipFill = computed(() => (props.dark ? '#FF653F' : '#1E104E'));
const chipText = computed(() => (props.dark ? '#1d100d' : '#FFC85C'));
const wordmarkFill = computed(() => (props.dark ? '#ffffff' : '#1E104E'));
</script>

<style scoped>
.just-a-logo {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.just-a-chip {
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 14px;
  font-weight: 800;
}

.just-a-wordmark {
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 26px;
  font-weight: 900;
  letter-spacing: -0.5px;
}
</style>
