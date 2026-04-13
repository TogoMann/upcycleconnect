import { createRouter, createWebHistory } from 'vue-router'

import PublicLayout from '@/layouts/PublicLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

import HomePage from '@/pages/home/HomePage.vue'
import PrestationPage from '@/pages/home/PrestationPage.vue'
import EventPage from '@/pages/home/EventPage.vue'
import AnnouncePage from '@/pages/home/AnnoucePage.vue'
import RepairPage from '@/pages/home/RepairPage.vue'
import ForumPage from '@/pages/home/ForumPage.vue'
import ForumThreadPage from '@/pages/home/ForumThreadPage.vue'
import AboutPage from '@/pages/home/AboutPage.vue'
import ConseilsPage from '@/pages/home/ConseilsPage.vue'
import AnnonceDetailPage from '@/pages/home/AnnonceDetailPage.vue'

import Login from '@/pages/auth/Login.vue'
import Register from '@/pages/auth/Register.vue'
import ForgotPassword from '@/pages/auth/ForgotPassword.vue'
import ResetPassword from '@/pages/auth/ResetPassword.vue'

import NotFound from '@/pages/errors/NotFound.vue'
import Forbidden from '@/pages/errors/Forbidden.vue'
import ServerError from '@/pages/errors/ServerError.vue'

import AdminDashboard from '@/pages/admin/Dashboard.vue'
import AdminUsers from '@/pages/admin/Users.vue'
import AdminCourses from '@/pages/admin/Courses.vue'
import AdminEvents from '@/pages/admin/Events.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            component: PublicLayout,
            children: [
                { path: '', component: HomePage },
                { path: 'prestations', component: PrestationPage },
                { path: 'evenements', component: EventPage },
                { path: 'annonces', component: AnnouncePage },
                { path: 'annonces/:id', component: AnnonceDetailPage },
                { path: 'reparer', component: RepairPage },
                { path: 'forum', component: ForumPage },
                { path: 'forum/:id', component: ForumThreadPage },
                { path: 'a-propos', component: AboutPage },
                { path: 'conseils', component: ConseilsPage },
                { path: 'auth/login', component: Login },
                { path: 'auth/register', component: Register },
                { path: 'auth/forgot-password', component: ForgotPassword },
                { path: 'auth/reset-password', component: ResetPassword },
                { path: '403', component: Forbidden },
                { path: '500', component: ServerError },
            ],
        },
        {
            path: '/admin',
            component: AdminLayout,
            children: [
                { path: '', component: AdminDashboard },
                { path: 'users', component: AdminUsers },
                { path: 'courses', component: AdminCourses },
                { path: 'events', component: AdminEvents },
            ],
        },
        {
            path: '/:pathMatch(.*)*',
            component: PublicLayout,
            children: [
                { path: '', component: NotFound },
            ],
        },
    ],
})

export default router
