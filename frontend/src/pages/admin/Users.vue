<template>
    <div class="admin-page">
        <!-- En-tête de section -->
        <div class="section-header">
            <h2 class="section-title">Gestion des utilisateurs</h2>
            <button class="btn-refresh" @click="adminStore.fetchUsers()">
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

        <!-- État : chargement -->
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
            Chargement des utilisateurs...
        </div>

        <!-- État : erreur -->
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

        <!-- Tableau -->
        <div class="table-wrap" v-if="!adminStore.isLoading">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Nom</th>
                        <th>Email</th>
                        <th>Rôle</th>
                        <th>Score</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="user in adminStore.users" :key="user.id" class="table-row">
                        <td class="cell-id">{{ user.id }}</td>

                        <!-- Nom avec avatar initiales -->
                        <td>
                            <div class="user-name-wrap">
                                <div class="user-avatar">
                                    {{ user.first_name?.[0] ?? '' }}{{ user.last_name?.[0] ?? '' }}
                                </div>
                                <span>{{ user.first_name }} {{ user.last_name }}</span>
                            </div>
                        </td>

                        <td class="cell-email">{{ user.email }}</td>

                        <!-- Badge rôle -->
                        <td>
                            <span
                                class="badge"
                                :class="user.role === 'admin' ? 'badge--admin' : 'badge--user'"
                            >
                                {{ user.role }}
                            </span>
                        </td>

                        <!-- Score avec barre visuelle -->
                        <td>
                            <div class="score-wrap">
                                <span class="score-value">{{ user.score }}</span>
                                <div class="score-bar-bg">
                                    <div
                                        class="score-bar-fill"
                                        :style="{
                                            width: Math.min((user.score / 100) * 100, 100) + '%',
                                        }"
                                    />
                                </div>
                            </div>
                        </td>

                        <td>
                            <button class="btn-delete" @click="deleteUser(user.id)">
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

                    <!-- État vide -->
                    <tr v-if="adminStore.users.length === 0">
                        <td colspan="6" class="cell-empty">Aucun utilisateur trouvé.</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

// ── Logique existante conservée intacte ──
const adminStore = useAdminStore()

const deleteUser = (id: number) => {
    if (confirm('Êtes-vous sûr de vouloir supprimer cet utilisateur ?')) {
        adminStore.deleteUser(id)
    }
}

onMounted(() => {
    adminStore.fetchUsers()
})
</script>

<style scoped>
/* ══ Charte graphique UCC ══ */
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

/* ── En-tête ── */
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
        color 0.2s,
        border-color 0.2s;
}
.btn-refresh:hover {
    background: var(--green-pale);
    color: var(--green-dark);
    border-color: var(--green-dark);
}

/* ── États ── */
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

/* ── Tableau ── */
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
    padding: 14px 20px;
    color: var(--charcoal);
    vertical-align: middle;
}

/* ID */
.cell-id {
    font-weight: 600;
    color: rgba(53, 53, 53, 0.45);
    font-size: 0.82rem;
    width: 60px;
}

/* Email */
.cell-email {
    font-size: 0.84rem;
    color: rgba(53, 53, 53, 0.7);
}

/* Nom avec avatar initiales */
.user-name-wrap {
    display: flex;
    align-items: center;
    gap: 10px;
}
.user-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: var(--green-mid);
    color: var(--white);
    font-size: 0.72rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    text-transform: uppercase;
    letter-spacing: 0.02em;
}

/* Badges rôle */
.badge {
    display: inline-block;
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.03em;
    text-transform: capitalize;
}
.badge--admin {
    background: rgba(8, 106, 53, 0.12);
    color: var(--green-dark);
}
.badge--user {
    background: rgba(139, 189, 148, 0.25);
    color: var(--green-mid);
}

/* Score avec barre */
.score-wrap {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 120px;
}
.score-value {
    font-weight: 700;
    color: var(--green-dark);
    font-size: 0.88rem;
    min-width: 28px;
    text-align: right;
}
.score-bar-bg {
    flex: 1;
    height: 6px;
    background: rgba(53, 53, 53, 0.12);
    border-radius: 99px;
    overflow: hidden;
}
.score-bar-fill {
    height: 100%;
    background: var(--green-mid);
    border-radius: 99px;
    transition: width 0.4s ease;
}

/* État vide */
.cell-empty {
    text-align: center;
    padding: 48px !important;
    color: rgba(53, 53, 53, 0.45);
    font-size: 0.9rem;
}

/* Bouton Supprimer */
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
