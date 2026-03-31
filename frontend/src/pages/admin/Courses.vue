<template>
    <div class="admin-page">
        <div class="section-header">
            <h2 class="section-title">Gestion des annonces</h2>
            <button class="btn-refresh" @click="adminStore.fetchCourses()">
                <svg
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <polyline points="23 4 23 10 17 10" />
                    <polyline points="1 20 1 14 7 14" />
                    <path
                        d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"
                    />
                </svg>
                Actualiser
            </button>
        </div>

        <div v-if="adminStore.isLoading" class="state-loading">
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                class="spin"
            >
                <path d="M21 12a9 9 0 1 1-6.219-8.56" />
            </svg>
            Chargement des annonces...
        </div>

        <div v-if="adminStore.error" class="state-error">
            <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
            {{ adminStore.error }}
        </div>

        <div class="table-wrap" v-if="!adminStore.isLoading">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Créé par</th>
                        <th>Prix</th>
                        <th>Statut</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="course in adminStore.courses" :key="course.id" class="table-row">
                        <td class="cell-id">{{ course.id }}</td>
                        <td>Utilisateur {{ course.created_by }}</td>
                        <td class="cell-price">{{ course.price }} €</td>
                        <td>
                            <span v-if="course.approved" class="badge badge--valide">Validé</span>
                            <span v-else class="badge badge--attente">En attente</span>
                        </td>
                        <td>
                            <button class="btn-delete" @click="deleteCourse(course.id)">
                                <svg
                                    width="14"
                                    height="14"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                >
                                    <polyline points="3 6 5 6 21 6" />
                                    <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
                                    <path d="M10 11v6M14 11v6" />
                                    <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
                                </svg>
                                Supprimer
                            </button>
                        </td>
                    </tr>

                    <tr v-if="adminStore.courses.length === 0">
                        <td colspan="5" class="cell-empty">Aucune annonce trouvée.</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

const deleteCourse = (id: number) => {
    if (confirm('Êtes-vous sûr de vouloir supprimer cette annonce ?')) {
        adminStore.deleteCourse(id)
    }
}

onMounted(() => {
    adminStore.fetchCourses()
})
</script>

<style scoped>
.admin-page {
    --cream: #f8f5ee;
    --green-dark: #086a35;
    --green-mid: #34895b;
    --green-light: #8bbd94;
    --green-pale: #d7ece1;
    --charcoal: #353535;
    --white: #ffffff;

    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 28px;
}
.section-title {
    font-size: 1.6rem;
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.02em;
    margin: 0;
}

.btn-refresh {
    display: flex;
    align-items: center;
    gap: 8px;
    background: transparent;
    color: var(--green-mid);
    border: 1.5px solid var(--green-mid);
    padding: 9px 18px;
    border-radius: 7px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        background 0.2s,
        color 0.2s;
}
.btn-refresh:hover {
    background: var(--green-pale);
    color: var(--green-dark);
    border-color: var(--green-dark);
}

.state-loading {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 20px 0;
    font-size: 0.9rem;
    color: rgba(53, 53, 53, 0.6);
}
.state-error {
    display: flex;
    align-items: center;
    gap: 10px;
    background: #fdecea;
    color: #c0392b;
    border: 1px solid #f5c6cb;
    padding: 14px 18px;
    border-radius: 8px;
    font-size: 0.875rem;
    margin-bottom: 20px;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}
.spin {
    animation: spin 1s linear infinite;
}

.table-wrap {
    background: var(--green-pale);
    border-radius: 14px;
    overflow: hidden;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.875rem;
}

.data-table thead tr {
    border-bottom: 1px solid rgba(53, 53, 53, 0.1);
}
.data-table th {
    padding: 16px 20px;
    text-align: left;
    font-size: 0.78rem;
    font-weight: 700;
    color: var(--green-mid);
    letter-spacing: 0.04em;
    text-transform: uppercase;
    background: transparent;
}

.table-row {
    border-bottom: 1px solid rgba(53, 53, 53, 0.07);
    transition: background 0.15s;
}
.table-row:last-child {
    border-bottom: none;
}
.table-row:hover {
    background: rgba(255, 255, 255, 0.5);
}

.data-table td {
    padding: 16px 20px;
    color: var(--charcoal);
    vertical-align: middle;
}

.cell-id {
    font-weight: 600;
    color: rgba(53, 53, 53, 0.5);
    font-size: 0.82rem;
    width: 60px;
}
.cell-price {
    font-weight: 700;
    color: var(--green-dark);
}
.cell-empty {
    text-align: center;
    padding: 48px !important;
    color: rgba(53, 53, 53, 0.45);
    font-size: 0.9rem;
}

.badge {
    display: inline-block;
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.02em;
}
.badge--valide {
    background: rgba(52, 137, 91, 0.15);
    color: var(--green-dark);
}
.badge--attente {
    background: rgba(255, 193, 7, 0.15);
    color: #856404;
}

.btn-delete {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    color: #c0392b;
    border: 1.5px solid rgba(192, 57, 43, 0.25);
    padding: 7px 14px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        background 0.2s,
        border-color 0.2s;
}
.btn-delete:hover {
    background: rgba(192, 57, 43, 0.08);
    border-color: #c0392b;
}
</style>
