import { createRouter, createWebHistory } from 'vue-router'
import { Login, Home } from '@/pages'
import AdminLayout from '@/layouts/AdminLayout.vue'
import AdminDashboard from '@/pages/admin/Dashboard.vue'
import AdminUsers from '@/pages/admin/Users.vue'
import AdminCourses from '@/pages/admin/Courses.vue'
import AdminEvents from '@/pages/admin/Events.vue'

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
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        {
          path: '',
          component: AdminDashboard,
        },
        {
          path: 'users',
          component: AdminUsers,
        },
        {
          path: 'courses',
          component: AdminCourses,
        },
        {
          path: 'events',
          component: AdminEvents,
        },
      ],
    },
  ],
})

export default router
