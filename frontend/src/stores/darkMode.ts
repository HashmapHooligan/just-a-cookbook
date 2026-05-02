import { defineStore } from 'pinia';
import { ref } from 'vue';
import { Dark } from 'quasar';

export const useDarkModeStore = defineStore('darkMode', () => {
  const getInitial = (): boolean => {
    const stored = localStorage.getItem('darkMode');
    if (stored !== null) return stored === 'true';
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  };

  const isDark = ref(getInitial());

  function apply() {
    Dark.set(isDark.value);
  }

  function toggle() {
    isDark.value = !isDark.value;
    localStorage.setItem('darkMode', String(isDark.value));
    apply();
  }

  apply();

  return { isDark, toggle };
});
