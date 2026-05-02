<template>
  <q-layout view="hHh lpR fFf">
    <q-header style="background-color: var(--color-header)">
      <q-toolbar style="color: var(--color-on-header)">
        <q-toolbar-title>
          <router-link to="/" class="no-underline" style="display: inline-flex; align-items: center; text-decoration: none">
            <JustALogo :dark="darkModeStore.isDark" />
          </router-link>
        </q-toolbar-title>

        <q-btn-dropdown
          flat
          no-caps
          split
          :label="t('nav.newRecipe')"
          icon="add"
          to="/recipes/new"
          style="color: var(--color-on-header)"
          class="font-label-lg q-mr-sm"
        >
          <q-list>
            <q-item clickable v-close-popup to="/recipes/bulk-add">
              <q-item-section avatar>
                <q-icon name="photo_library" />
              </q-item-section>
              <q-item-section>{{ t('nav.bulkAdd') }}</q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>

        <q-btn
          flat
          round
          :label="localeStore.current === 'en-US' ? 'DE' : 'EN'"
          class="font-label-lg q-mr-xs"
          style="color: var(--color-on-header)"
          @click="localeStore.toggle()"
          title="Switch language"
        />

        <q-btn
          flat
          round
          :icon="darkModeStore.isDark ? 'light_mode' : 'dark_mode'"
          style="color: var(--color-on-header)"
          :title="darkModeStore.isDark ? 'Switch to light mode' : 'Switch to dark mode'"
          @click="darkModeStore.toggle()"
        />
      </q-toolbar>
    </q-header>

    <q-page-container style="background-color: var(--color-background)">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { useLocaleStore } from 'src/stores/locale';
import { useDarkModeStore } from 'src/stores/darkMode';
import JustALogo from 'src/components/JustALogo.vue';

const { t } = useI18n();
const localeStore = useLocaleStore();
const darkModeStore = useDarkModeStore();
</script>
