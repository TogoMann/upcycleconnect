import { createRouter, createWebHistory } from 'vue-router'
import { Login, Home } from '@/pages'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/auth/login',
      component: Login,
    },
    {
      path: '/',
      component: Home,
    },
  ],
})

export default router
