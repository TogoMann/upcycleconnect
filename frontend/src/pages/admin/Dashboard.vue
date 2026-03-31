<template>
    <div class="dashboard">
        <h1 class="page-title">Tableau de bord.</h1>

        <div class="kpi-grid">
            <div class="kpi-card">
                <div class="kpi-icon">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                        <circle cx="9" cy="7" r="4" />
                        <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                        <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                    </svg>
                </div>
                <div class="kpi-value">{{ adminStore.users.length }}</div>
                <div class="kpi-label">Utilisateurs</div>
            </div>

            <div class="kpi-card">
                <div class="kpi-icon">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                        <polyline points="14 2 14 8 20 8" />
                    </svg>
                </div>
                <div class="kpi-value">{{ adminStore.courses.length }}</div>
                <div class="kpi-label">Annonces</div>
            </div>

            <div class="kpi-card">
                <div class="kpi-icon">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                        <line x1="16" y1="2" x2="16" y2="6" />
                        <line x1="8" y1="2" x2="8" y2="6" />
                        <line x1="3" y1="10" x2="21" y2="10" />
                    </svg>
                </div>
                <div class="kpi-value">{{ adminStore.events.length }}</div>
                <div class="kpi-label">Événements</div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

onMounted(() => {
    adminStore.fetchUsers()
    adminStore.fetchCourses()
    adminStore.fetchEvents()
})
</script>

<style scoped>
.dashboard {
    --green-dark: #086a35;
    --green-mid: #34895b;
    --green-light: #8bbd94;
    --green-pale: #d7ece1;
    --charcoal: #353535;
    --white: #ffffff;

    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.page-title {
    font-size: clamp(2.2rem, 4vw, 3.2rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 32px;
    line-height: 1.08;
}

.kpi-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
}

.kpi-card {
    background: var(--green-pale);
    border-radius: 14px;
    padding: 24px 22px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.kpi-icon {
    width: 40px;
    height: 40px;
    background: var(--green-mid);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--white);
    margin-bottom: 4px;
}

.kpi-value {
    font-size: 2.2rem;
    font-weight: 800;
    color: var(--green-dark);
    line-height: 1;
    letter-spacing: -0.03em;
}

.kpi-label {
    font-size: 0.85rem;
    color: var(--green-dark);
    font-weight: 500;
    opacity: 0.75;
}

@media (max-width: 600px) {
    .kpi-grid {
        grid-template-columns: 1fr;
    }
}
</style>
