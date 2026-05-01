<template>
  <q-layout view="hHh lpR fFf">
    <q-header style="background-color: var(--color-header)">
      <q-toolbar style="color: var(--color-primary-container)">
        <q-toolbar-title>
          <router-link
            to="/"
            class="font-headline-md no-underline"
            style="color: var(--color-primary-container)"
          >
            Just a Cookbook
          </router-link>
        </q-toolbar-title>

        <q-btn
          flat
          no-caps
          :label="t('nav.newRecipe')"
          icon="add"
          to="/recipes/new"
          style="color: var(--color-primary-container)"
          class="font-label-lg q-mr-sm"
        />

        <q-btn
          flat
          round
          :label="localeStore.current === 'en-US' ? 'DE' : 'EN'"
          class="font-label-lg"
          style="color: var(--color-primary-container)"
          @click="localeStore.toggle()"
          title="Switch language"
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

const { t } = useI18n();
const localeStore = useLocaleStore();
</script>
