import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/RecipeOverviewPage.vue') },
      { path: 'recipes/new', component: () => import('pages/RecipeFormPage.vue') },
      { path: 'recipes/:id', component: () => import('pages/RecipeDetailPage.vue') },
      { path: 'recipes/:id/edit', component: () => import('pages/RecipeFormPage.vue') },
    ],
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
