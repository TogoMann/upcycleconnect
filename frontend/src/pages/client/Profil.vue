<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useClientStore } from '@/stores/client'

const authStore = useAuthStore()
const clientStore = useClientStore()

const editing = ref(false)
const saved = ref(false)
const saveError = ref('')
const showHistory = ref(false)

const form = reactive({
    first_name: authStore.user?.first_name ?? '',
    last_name: authStore.user?.last_name ?? '',
    email: authStore.user?.email ?? '',
})

watch(showHistory, (val) => {
    if (val) clientStore.fetchScoreHistory()
})

watch(
    () => authStore.user,
    user => {
        if (user) {
            form.first_name = user.first_name
            form.last_name = user.last_name
            form.email = user.email
        }
    },
)

const roleLabels: Record<string, string> = {
    client: 'Particulier',
    pro: 'Professionnel',
    interne: 'Interne',
    admin: 'Administrateur',
}

function startEdit() {
    editing.value = true
    saved.value = false
}

function cancelEdit() {
    if (authStore.user) {
        form.first_name = authStore.user.first_name
        form.last_name = authStore.user.last_name
        form.email = authStore.user.email
    }
    editing.value = false
}

async function handleSave() {
    saveError.value = ''
    try {
        await clientStore.updateProfile({
            first_name: form.first_name,
            last_name: form.last_name,
            email: form.email,
        })
        editing.value = false
        saved.value = true
        setTimeout(() => (saved.value = false), 3000)
    } catch (e: any) {
        saveError.value = e.message
    }
}
</script>

<template>
    <div class="page">
        <h1 class="page-title">Mon Profil.</h1>

        <div class="profile-layout">
            <div class="profile-card">
                <div class="avatar-section">
                    <div class="avatar">
                        <span class="avatar-initials">
                            {{ (authStore.user?.first_name?.[0] ?? '') + (authStore.user?.last_name?.[0] ?? authStore.user?.username?.[0] ?? '?') }}
                        </span>
                    </div>
                    <div class="avatar-info">
                        <p class="avatar-name">
                            {{ authStore.user?.first_name || '' }} {{ authStore.user?.last_name || authStore.user?.username || '' }}
                        </p>
                        <span class="role-badge">
                            {{ roleLabels[authStore.user?.role ?? ''] ?? authStore.user?.role ?? 'Particulier' }}
                        </span>
                    </div>
                </div>

                <div v-if="saved" class="saved-banner">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                        <polyline points="20 6 9 17 4 12" />
                    </svg>
                    Profil mis à jour
                </div>

                <template v-if="!editing">
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="info-label">Prénom</span>
                            <span class="info-value">{{ authStore.user?.first_name || '—' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Nom</span>
                            <span class="info-value">{{ authStore.user?.last_name || '—' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Nom d'utilisateur</span>
                            <span class="info-value">{{ authStore.user?.username || '—' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Email</span>
                            <span class="info-value">{{ authStore.user?.email || '—' }}</span>
                        </div>
                    </div>

                    <button class="btn-edit" @click="startEdit">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                        </svg>
                        Modifier le profil
                    </button>
                </template>

                <template v-else>
                    <form class="edit-form" @submit.prevent="handleSave">
                        <div class="form-row">
                            <div class="form-group">
                                <label class="form-label">Prénom</label>
                                <input v-model="form.first_name" type="text" class="form-input" placeholder="Prénom" />
                            </div>
                            <div class="form-group">
                                <label class="form-label">Nom</label>
                                <input v-model="form.last_name" type="text" class="form-input" placeholder="Nom" />
                            </div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">Email</label>
                            <input v-model="form.email" type="email" class="form-input" placeholder="email@exemple.com" />
                        </div>
                        <div v-if="saveError" class="save-error">{{ saveError }}</div>
                        <div class="form-actions">
                            <button type="button" class="btn-cancel" @click="cancelEdit">Annuler</button>
                            <button type="submit" class="btn-save">Enregistrer</button>
                        </div>
                    </form>
                </template>
            </div>

            <div class="side-cards">
                <div class="side-card">
                    <div class="side-card-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
                        </svg>
                    </div>
                    <div>
                        <p class="side-card-label">Score upcycling</p>
                        <p class="side-card-value">{{ clientStore.score }} pts</p>
                    </div>
                    <div class="side-card-actions">
                        <router-link to="/particulier/score" class="side-card-link">Voir →</router-link>
                        <button class="btn-history-small" @click="showHistory = true">Historique</button>
                    </div>
                </div>

                <div class="side-card">
                    <div class="side-card-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                            <polyline points="14 2 14 8 20 8" />
                        </svg>
                    </div>
                    <div>
                        <p class="side-card-label">Annonces</p>
                        <p class="side-card-value">{{ clientStore.annonces.length }}</p>
                    </div>
                    <router-link to="/particulier/annonces" class="side-card-link">Voir →</router-link>
                </div>

                <div class="side-card">
                    <div class="side-card-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                        </svg>
                    </div>
                    <div>
                        <p class="side-card-label">Dépôts</p>
                        <p class="side-card-value">{{ clientStore.depots.length }}</p>
                    </div>
                    <router-link to="/particulier/conteneurs" class="side-card-link">Voir →</router-link>
                </div>
            </div>
        </div>

        <!-- Modal Historique Points -->
        <div v-if="showHistory" class="modal-overlay" @click.self="showHistory = false">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title">Mon historique de points</h3>
                    <button class="btn-close" @click="showHistory = false">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <line x1="18" y1="6" x2="6" y2="18" />
                            <line x1="6" y1="6" x2="18" y2="18" />
                        </svg>
                    </button>
                </div>
                
                <div class="modal-body history-modal-body">
                    <div v-if="clientStore.isLoading" class="history-loading">Chargement...</div>
                    <div v-else-if="clientStore.scoreHistory.length === 0" class="history-empty">Aucun historique trouvé.</div>
                    <div v-else class="history-list">
                        <div v-for="item in clientStore.scoreHistory" :key="item.id" class="history-item">
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
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 32px;
    line-height: 1.08;
}

.profile-layout {
    display: flex;
    gap: 24px;
    align-items: flex-start;
}

.profile-card {
    flex: 1;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    padding: 28px;
    min-width: 0;
}

.avatar-section {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 28px;
    padding-bottom: 24px;
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
}
.avatar {
    width: 64px;
    height: 64px;
    background: var(--green-pale);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    border: 3px solid var(--green-light);
}
.avatar-initials {
    font-size: 1.4rem;
    font-weight: 800;
    color: var(--green-dark);
    text-transform: uppercase;
    letter-spacing: -0.02em;
}
.avatar-name {
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 6px;
}
.role-badge {
    display: inline-block;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 0.72rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
    background: var(--green-pale);
    color: var(--green-dark);
}

.saved-banner {
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--green-pale);
    color: var(--green-dark);
    padding: 10px 16px;
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    margin-bottom: 20px;
}
.saved-banner svg {
    width: 16px;
    height: 16px;
}

.info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    margin-bottom: 24px;
}
.info-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}
.info-label {
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.5;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
}
.info-value {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--charcoal);
}

.btn-edit {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 18px;
    background: transparent;
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--charcoal);
    cursor: pointer;
    font-family: inherit;
    transition: border-color 0.2s, background 0.2s;
}
.btn-edit:hover {
    border-color: var(--green-mid);
    background: var(--green-pale);
    color: var(--green-dark);
}
.btn-edit svg {
    width: 16px;
    height: 16px;
}

.edit-form {
    display: flex;
    flex-direction: column;
    gap: 18px;
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
}
.form-group {
    display: flex;
    flex-direction: column;
    gap: 7px;
}
.form-label {
    font-size: 0.82rem;
    font-weight: 600;
    color: var(--charcoal);
}
.form-input {
    width: 100%;
    padding: 12px 14px;
    font-size: 0.9rem;
    font-family: inherit;
    color: var(--charcoal);
    background: var(--cream);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
}
.form-input:focus {
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.1);
}
.form-actions {
    display: flex;
    gap: 10px;
    justify-content: flex-end;
    padding-top: 8px;
    border-top: 1px solid rgba(53, 53, 53, 0.08);
}
.btn-cancel {
    padding: 10px 18px;
    background: transparent;
    border: 1.5px solid rgba(53, 53, 53, 0.18);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    color: rgba(53, 53, 53, 0.6);
    cursor: pointer;
    font-family: inherit;
    transition: border-color 0.2s, color 0.2s;
}
.btn-cancel:hover {
    border-color: rgba(53, 53, 53, 0.35);
    color: var(--charcoal);
}
.save-error {
    background: rgba(229, 62, 62, 0.08);
    border: 1px solid rgba(229, 62, 62, 0.25);
    border-radius: 8px;
    padding: 10px 14px;
    font-size: 0.82rem;
    color: #c53030;
}
.btn-save {
    padding: 10px 20px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
}
.btn-save:hover {
    background: var(--green-mid);
}

.side-cards {
    width: 240px;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.side-card {
    background: var(--green-pale);
    border-radius: 12px;
    padding: 16px;
    display: flex;
    align-items: center;
    gap: 12px;
}
.side-card-icon {
    width: 36px;
    height: 36px;
    background: var(--white);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--green-mid);
    flex-shrink: 0;
}
.side-card-icon svg {
    width: 18px;
    height: 18px;
}
.side-card-label {
    font-size: 0.72rem;
    color: var(--green-dark);
    opacity: 0.65;
    margin: 0 0 2px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.03em;
}
.side-card-value {
    font-size: 1.1rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
    margin: 0;
}
.side-card-link {
    margin-left: auto;
    font-size: 0.78rem;
    color: var(--green-mid);
    text-decoration: none;
    font-weight: 600;
    flex-shrink: 0;
    transition: color 0.2s;
}
.side-card-link:hover {
    color: var(--green-dark);
}

.side-card-actions {
    margin-left: auto;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
}
.btn-history-small {
    font-size: 0.7rem;
    color: var(--green-mid);
    background: transparent;
    border: none;
    cursor: pointer;
    font-weight: 600;
    padding: 0;
    font-family: inherit;
    text-decoration: underline;
}
.btn-history-small:hover {
    color: var(--green-dark);
}

/* Modal Historique */
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
    font-size: 1.1rem;
    font-weight: 700;
}
.btn-close {
    background: none;
    border: none;
    cursor: pointer;
    color: rgba(53, 53, 53, 0.4);
}
.history-modal-body {
    padding: 24px;
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

@media (max-width: 800px) {
    .profile-layout {
        flex-direction: column;
    }
    .side-cards {
        width: 100%;
        flex-direction: row;
        flex-wrap: wrap;
    }
    .side-card {
        flex: 1;
        min-width: 160px;
    }
    .info-grid {
        grid-template-columns: 1fr;
    }
    .form-row {
        grid-template-columns: 1fr;
    }
}
</style>
