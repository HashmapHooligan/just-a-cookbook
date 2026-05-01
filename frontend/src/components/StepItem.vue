<template>
  <div class="step-item">
    <div class="step-number">{{ number }}</div>
    <p class="font-instruction-step step-item__text">{{ step.description }}</p>
  </div>
</template>

<script setup lang="ts">
import type { Step } from 'src/models/recipe';

defineProps<{ step: Step; number: number }>();
</script>

<style scoped>
.step-item {
  display: flex;
  gap: 20px;
  align-items: flex-start;
  padding: 16px 0;
}

.step-item__text {
  flex: 1;
  margin: 0;
  padding-top: 10px;
  color: var(--color-on-surface);
}
</style>
