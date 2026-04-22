<template>
    <div class="admin-page">
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

                        <td>
                            <div class="user-name-wrap">
                                <div class="user-avatar">
                                    {{ user.first_name?.[0] ?? '' }}{{ user.last_name?.[0] ?? '' }}
                                </div>
                                <span>{{ user.first_name }} {{ user.last_name }}</span>
                            </div>
                        </td>

                        <td class="cell-email">{{ user.email }}</td>

                        <td>
                            <span
                                class="badge"
                                :class="`badge--${user.role}`"
                            >
                                {{ user.role }}
                            </span>
                        </td>

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

                        <td class="cell-actions">
                            <button class="btn-history" @click="viewHistory(user)">
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
                                    <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                                Historique
                            </button>
                            <button class="btn-edit" @click="openEdit(user)">
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
                                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                </svg>
                                Modifier
                            </button>
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

                    <tr v-if="adminStore.users.length === 0">
                        <td colspan="6" class="cell-empty">Aucun utilisateur trouvé.</td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Modal Historique Points -->
        <div v-if="showHistoryModal" class="modal-overlay" @click.self="closeHistory">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">Historique des points - {{ historyUser?.username }}</h3>
                    <button class="btn-close" @click="closeHistory">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <line x1="18" y1="6" x2="6" y2="18" />
                            <line x1="6" y1="6" x2="18" y2="18" />
                        </svg>
                    </button>
                </div>
                
                <div class="modal-body history-modal-body">
                    <div v-if="isLoadingHistory" class="history-loading">Chargement...</div>
                    <div v-else-if="scoreHistory.length === 0" class="history-empty">Aucun historique trouvé.</div>
                    <div v-else class="history-list">
                        <div v-for="item in scoreHistory" :key="item.id" class="history-item">
                            <div class="history-item-main">
                                <span class="history-points" :class="item.points >= 0 ? 'pos' : 'neg'">
                                    {{ item.points >= 0 ? '+' : '' }}{{ item.points }}
                                </span>
                                <span class="history-desc">{{ item.description }}</span>
                            </div>
                            <span class="history-date">{{ new Date(item.created_at).toLocaleDateString() }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Modal Edition User -->
        <div v-if="showEditModal" class="modal-overlay" @click.self="closeEdit">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">Modifier l'utilisateur</h3>
                    <button class="btn-close" @click="closeEdit">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <line x1="18" y1="6" x2="6" y2="18" />
                            <line x1="6" y1="6" x2="18" y2="18" />
                        </svg>
                    </button>
                </div>
                
                <form @submit.prevent="saveUser" class="modal-body">
                    <div class="form-row">
                        <div class="form-group">
                            <label>Prénom</label>
                            <input type="text" v-model="editForm.first_name" required class="form-input" />
                        </div>
                        <div class="form-group">
                            <label>Nom</label>
                            <input type="text" v-model="editForm.last_name" required class="form-input" />
                        </div>
                    </div>

                    <div class="form-group">
                        <label>Nom d'utilisateur</label>
                        <input type="text" v-model="editForm.username" required class="form-input" />
                    </div>

                    <div class="form-group">
                        <label>Email</label>
                        <input type="email" v-model="editForm.email" required class="form-input" />
                    </div>
                    
                    <div class="form-group">
                        <label>Rôle</label>
                        <select v-model="editForm.role" class="form-select">
                            <option value="client">Client</option>
                            <option value="pro">Professionnel</option>
                            <option value="interne">Salarié</option>
                            <option value="admin">Administrateur</option>
                        </select>
                    </div>

                    <div class="form-actions">
                        <button type="button" class="btn-cancel" @click="closeEdit">Annuler</button>
                        <button type="submit" class="btn-save" :disabled="isSaving">
                            {{ isSaving ? 'Enregistrement...' : 'Enregistrer' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

const showEditModal = ref(false)
const showHistoryModal = ref(false)
const isLoadingHistory = ref(false)
const scoreHistory = ref<any[]>([])
const historyUser = ref<any>(null)
const isSaving = ref(false)
const editingUser = ref<any>(null)
const editForm = reactive({
    first_name: '',
    last_name: '',
    email: '',
    role: '',
    username: '',
    language_preference: 'fr'
})

const viewHistory = async (user: any) => {
    historyUser.value = user
    showHistoryModal.value = true
    isLoadingHistory.value = true
    try {
        scoreHistory.value = await adminStore.getScoreHistory(user.id)
    } catch (err) {
        console.error('Failed to fetch history', err)
    } finally {
        isLoadingHistory.value = false
    }
}

const closeHistory = () => {
    showHistoryModal.value = false
    historyUser.value = null
    scoreHistory.value = []
}

const openEdit = (user: any) => {
    editingUser.value = user
    editForm.first_name = user.first_name || ''
    editForm.last_name = user.last_name || ''
    editForm.email = user.email
    editForm.role = user.role
    editForm.username = user.username || ''
    editForm.language_preference = user.language_preference || 'fr'
    showEditModal.value = true
}

const closeEdit = () => {
    showEditModal.value = false
    editingUser.value = null
}

const saveUser = async () => {
    console.log('saveUser called', editingUser.value)
    if (!editingUser.value) return
    
    isSaving.value = true
    try {
        console.log('updating user', editingUser.value.id, editForm)
        await adminStore.updateUser(editingUser.value.id, {
            first_name: editForm.first_name,
            last_name: editForm.last_name,
            email: editForm.email,
            role: editForm.role,
            username: editForm.username,
            language_preference: editForm.language_preference
        })
        closeEdit()
    } catch (err) {
        console.error('Update failed', err)
        // Error handled by store
    } finally {
        isSaving.value = false
    }
}

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
/* ... existing styles ... */
.cell-actions {
    display: flex;
    gap: 8px;
    white-space: nowrap;
    width: 1%;
}

.btn-edit {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    color: var(--green-mid);
    border: 1.5px solid rgba(52, 137, 91, 0.25);
    padding: 7px 14px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s, border-color 0.2s;
}
.btn-edit:hover {
    background: var(--green-pale);
    border-color: var(--green-mid);
}

.btn-history {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    color: var(--charcoal);
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    padding: 7px 14px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s, border-color 0.2s;
}
.btn-history:hover {
    background: rgba(53, 53, 53, 0.05);
    border-color: rgba(53, 53, 53, 0.3);
}

.history-modal-body {
    max-height: 400px;
    overflow-y: auto;
}
.history-loading, .history-empty {
    padding: 40px;
    text-align: center;
    color: rgba(53, 53, 53, 0.5);
    font-size: 0.9rem;
}
.history-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.history-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px;
    background: var(--cream);
    border-radius: 8px;
}
.history-item-main {
    display: flex;
    align-items: center;
    gap: 12px;
}
.history-points {
    font-weight: 800;
    font-size: 0.95rem;
    min-width: 40px;
}
.history-points.pos { color: var(--green-dark); }
.history-points.neg { color: #c53030; }
.history-desc {
    font-size: 0.88rem;
    font-weight: 600;
}
.history-date {
    font-size: 0.78rem;
    color: rgba(53, 53, 53, 0.5);
    font-weight: 500;
}

/* Modal Styles */
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
}
.modal-content {
    background: var(--white);
    width: 100%;
    max-width: 450px;
    border-radius: 16px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}
.modal-header {
    padding: 20px 24px;
    border-bottom: 1px solid rgba(53, 53, 53, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.modal-title {
    margin: 0;
    font-size: 1.2rem;
    font-weight: 700;
}
.btn-close {
    background: none;
    border: none;
    cursor: pointer;
    color: rgba(53, 53, 53, 0.4);
}
.modal-body {
    padding: 24px;
}
.form-row {
    display: flex;
    gap: 16px;
}
.form-row .form-group {
    flex: 1;
}
.form-group {
    margin-bottom: 20px;
}
.form-group label {
    display: block;
    font-size: 0.85rem;
    font-weight: 600;
    margin-bottom: 8px;
    color: var(--charcoal);
}
.form-input, .form-select {
    width: 100%;
    padding: 10px 14px;
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    border-radius: 8px;
    font-family: inherit;
    font-size: 0.95rem;
}
.form-input:focus, .form-select:focus {
    outline: none;
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px var(--green-pale);
}
.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 32px;
}
.btn-cancel {
    background: transparent;
    border: none;
    font-weight: 600;
    cursor: pointer;
    color: rgba(53, 53, 53, 0.6);
}
.btn-save {
    background: var(--green-mid);
    color: white;
    border: none;
    padding: 10px 24px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
}
.btn-save:hover {
    background: var(--green-dark);
}
.btn-save:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

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
        color 0.2s,
        border-color 0.2s;
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
    overflow-x: auto;
}
.data-table {
    width: 100%;
    min-width: 900px;
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

.cell-id {
    font-weight: 600;
    color: rgba(53, 53, 53, 0.45);
    font-size: 0.82rem;
    width: 60px;
}

.cell-email {
    font-size: 0.84rem;
    color: rgba(53, 53, 53, 0.7);
}

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
.badge--client {
    background: rgba(139, 189, 148, 0.25);
    color: var(--green-mid);
}
.badge--pro {
    background: rgba(52, 137, 91, 0.12);
    color: #2c3e50;
}
.badge--interne {
    background: rgba(52, 152, 219, 0.12);
    color: #2980b9;
}

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

.cell-empty {
    text-align: center;
    padding: 48px !important;
    color: rgba(53, 53, 53, 0.45);
    font-size: 0.9rem;
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
