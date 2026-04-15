import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
    interface RouteMeta {
        requiresAuth?: boolean
        role?: string
    }
}

import PublicLayout from '@/layouts/PublicLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'
import ClientLayout from '@/layouts/ClientLayout.vue'
import ProLayout from '@/layouts/ProLayout.vue'
import SalarieLayout from '@/layouts/SalarieLayout.vue'

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
import AdminAnnonceDetail from '@/pages/admin/AnnonceDetail.vue'
import AdminConteneurs from '@/pages/admin/Conteneurs.vue'
import AdminDepots from '@/pages/admin/Depots.vue'
import AdminAbonnements from '@/pages/admin/Abonnements.vue'
import AdminFinancier from '@/pages/admin/Financier.vue'
import AdminCommissions from '@/pages/admin/Commissions.vue'
import AdminPublicites from '@/pages/admin/Publicites.vue'
import AdminCatalogue from '@/pages/admin/Catalogue.vue'
import AdminNotifications from '@/pages/admin/Notifications.vue'
import AdminDocuments from '@/pages/admin/Documents.vue'
import AdminPlannings from '@/pages/admin/Plannings.vue'
import AdminProjets from '@/pages/admin/Projets.vue'
import AdminEvenements from '@/pages/admin/Evenements.vue'
import AdminActualites from '@/pages/admin/Actualites.vue'
import AdminForum from '@/pages/admin/Forum.vue'
import AdminLangues from '@/pages/admin/Langues.vue'
import AdminLogs from '@/pages/admin/Logs.vue'
import AdminParametres from '@/pages/admin/Parametres.vue'

import ClientDashboard from '@/pages/client/Dashboard.vue'
import MesAnnonces from '@/pages/client/MesAnnonces.vue'
import NouvelleAnnonce from '@/pages/client/NouvelleAnnonce.vue'
import MesDepots from '@/pages/client/MesDepots.vue'
import DeposerObjet from '@/pages/client/DeposerObjet.vue'
import UpcyclingScore from '@/pages/client/UpcyclingScore.vue'
import Planning from '@/pages/client/Planning.vue'
import Catalogue from '@/pages/client/Catalogue.vue'
import Paiement from '@/pages/client/Paiement.vue'
import ConfirmationPaiement from '@/pages/client/ConfirmationPaiement.vue'
import ClientProfil from '@/pages/client/Profil.vue'

import ProDashboard from '@/pages/pro/Dashboard.vue'
import ProAbonnements from '@/pages/pro/Abonnements.vue'
import ProPublicites from '@/pages/pro/Publicites.vue'
import ProFacturation from '@/pages/pro/Facturation.vue'
import ProAnnonces from '@/pages/pro/Annonces.vue'
import ProRecupererObjet from '@/pages/pro/RecupererObjet.vue'
import ProMesProjets from '@/pages/pro/MesProjets.vue'
import ProNouveauProjet from '@/pages/pro/NouveauProjet.vue'
import ProProjetDetail from '@/pages/pro/ProjetDetail.vue'
import ProTableauBordAvance from '@/pages/pro/TableauBordAvance.vue'
import ProProfil from '@/pages/pro/Profil.vue'

import SalarieDashboard from '@/pages/salarie/Dashboard.vue'
import SalarieMesFormations from '@/pages/salarie/MesFormations.vue'
import SalarieNouvelleFormation from '@/pages/salarie/NouvelleFormation.vue'
import SalarieEditFormation from '@/pages/salarie/EditFormation.vue'
import SalariePlanning from '@/pages/salarie/Planning.vue'
import SalarieConseils from '@/pages/salarie/Conseils.vue'
import SalarieForum from '@/pages/salarie/Forum.vue'
import SalarieProfil from '@/pages/salarie/Profil.vue'

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
            path: '/particulier',
            component: ClientLayout,
            meta: { requiresAuth: true, role: 'client' },
            children: [
                { path: '', component: ClientDashboard },
                { path: 'annonces', component: MesAnnonces },
                { path: 'annonces/nouvelle', component: NouvelleAnnonce },
                { path: 'conteneurs', component: MesDepots },
                { path: 'conteneurs/deposer', component: DeposerObjet },
                { path: 'score', component: UpcyclingScore },
                { path: 'planning', component: Planning },
                { path: 'catalogue', component: Catalogue },
                { path: 'paiement', component: Paiement },
                { path: 'paiement/confirmation', component: ConfirmationPaiement },
                { path: 'profil', component: ClientProfil },
            ],
        },
        {
            path: '/pro',
            component: ProLayout,
            meta: { requiresAuth: true, role: 'pro' },
            children: [
                { path: '', component: ProDashboard },
                { path: 'abonnements', component: ProAbonnements },
                { path: 'publicites', component: ProPublicites },
                { path: 'facturation', component: ProFacturation },
                { path: 'annonces', component: ProAnnonces },
                { path: 'conteneurs/recuperer', component: ProRecupererObjet },
                { path: 'projets', component: ProMesProjets },
                { path: 'projets/nouveau', component: ProNouveauProjet },
                { path: 'projets/:id', component: ProProjetDetail },
                { path: 'dashboard-avance', component: ProTableauBordAvance },
                { path: 'profil', component: ProProfil },
            ],
        },
        {
            path: '/salarie',
            component: SalarieLayout,
            meta: { requiresAuth: true, role: 'interne' },
            children: [
                { path: '', component: SalarieDashboard },
                { path: 'formations', component: SalarieMesFormations },
                { path: 'formations/nouvelle', component: SalarieNouvelleFormation },
                { path: 'formations/:id/edit', component: SalarieEditFormation },
                { path: 'planning', component: SalariePlanning },
                { path: 'conseils', component: SalarieConseils },
                { path: 'forum', component: SalarieForum },
                { path: 'profil', component: SalarieProfil },
            ],
        },
        {
            path: '/admin',
            component: AdminLayout,
            meta: { requiresAuth: true, role: 'admin' },
            children: [
                { path: '', component: AdminDashboard },
                { path: 'annonces/:id', component: AdminAnnonceDetail },
                { path: 'conteneurs', component: AdminConteneurs },
                { path: 'depots', component: AdminDepots },
                { path: 'abonnements', component: AdminAbonnements },
                { path: 'financier', component: AdminFinancier },
                { path: 'commissions', component: AdminCommissions },
                { path: 'publicites', component: AdminPublicites },
                { path: 'catalogue', component: AdminCatalogue },
                { path: 'notifications', component: AdminNotifications },
                { path: 'documents', component: AdminDocuments },
                { path: 'plannings', component: AdminPlannings },
                { path: 'projets', component: AdminProjets },
                { path: 'evenements', component: AdminEvenements },
                { path: 'actualites', component: AdminActualites },
                { path: 'forum', component: AdminForum },
                { path: 'langues', component: AdminLangues },
                { path: 'logs', component: AdminLogs },
                { path: 'parametres', component: AdminParametres },
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

router.beforeEach((to) => {
    const auth = useAuthStore()

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        return '/auth/login'
    }

    if (to.meta.requiresAuth && to.meta.role && auth.userRole !== to.meta.role) {
        return '/403'
    }
})

export default router
