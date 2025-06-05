import { createRouter, createWebHistory } from 'vue-router';
import { compile } from 'vue/dist/vue.cjs.prod';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/Home/Index.vue'),
    },
    {
      path: '/room',
      component: () => import('@/views/Room/Index.vue'),
    },
  ],
});

export default router;
