import { defineStore } from 'pinia';
import { ref } from 'vue';
import { i18n, type MessageLanguages } from 'src/boot/i18n';

const LOCALES: { value: MessageLanguages; label: string }[] = [
  { value: 'en-US', label: 'EN' },
  { value: 'de', label: 'DE' },
];

export const useLocaleStore = defineStore('locale', () => {
  const current = ref<MessageLanguages>(
    (localStorage.getItem('locale') as MessageLanguages) ?? 'en-US',
  );

  function toggle() {
    const next = current.value === 'en-US' ? 'de' : 'en-US';
    set(next);
  }

  function set(locale: MessageLanguages) {
    current.value = locale;
    localStorage.setItem('locale', locale);
    (i18n.global.locale as unknown as { value: string }).value = locale;
  }

  return { current, toggle, set, locales: LOCALES };
});
