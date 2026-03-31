import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'

import HomePage from '@/pages/home/HomePage.vue'
import PrestationPage from '@/pages/home/PrestationPage.vue'
import EventPage from '@/pages/home/EventPage.vue'
import AnnouncePage from '@/pages/home/AnnoucePage.vue'
import RepairPage from '@/pages/home/RepairPage.vue'

import Login from '@/pages/auth/Login.vue'

import AdminDashboard from '@/pages/admin/Dashboard.vue'
import AdminUsers from '@/pages/admin/Users.vue'
import AdminCourses from '@/pages/admin/Courses.vue'
import AdminEvents from '@/pages/admin/Events.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            component: HomePage,
        },
        {
            path: '/prestations',
            component: PrestationPage,
        },
        {
            path: '/evenements',
            component: EventPage,
        },
        {
            path: '/annonces',
            component: AnnouncePage,
        },
        {
            path: '/reparer',
            component: RepairPage,
        },
        {
            path: '/auth/login',
            component: Login,
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
