<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()

interface Thread {
    id: number
    titre: string
    auteur: string
    date: string
    epingle: boolean
    statut: string
    replies: number
}

interface Post {
    id: number
    thread_id: number
    created_by: number
    username: string
    content: string
    created_at: string
}

const threads = ref<Thread[]>([])
const loading = ref(true)
const error = ref('')
const success = ref('')

const selectedThread = ref<Thread | null>(null)
const posts = ref<Post[]>([])
const loadingPosts = ref(false)
const showThreadModal = ref(false)

const showBanModal = ref(false)
const banForm = reactive({
    username: '',
    isPermanent: true,
    durationHours: 24,
})

async function fetchThreads() {
    loading.value = true
    error.value = ''
    const token = authStore.token
    if (!token) {
        loading.value = false
        return
    }
    try {
        const res = await fetch(`${API_BASE}/salarie/forum`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (!res.ok) throw new Error(t('salarie.forum.errorLoad'))
        threads.value = await res.json()
    } catch (e: any) {
        error.value = e.message || t('salarie.forum.errorLoad')
    } finally {
        loading.value = false
    }
}

onMounted(fetchThreads)

async function epingler(thread: Thread) {
    try {
        const res = await fetch(`${API_BASE}/salarie/forum/${thread.id}/epingler`, {
            method: 'POST',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (!res.ok) throw new Error(t('salarie.forum.errorAction'))
        thread.epingle = !thread.epingle
    } catch (e: any) {
        error.value = e.message || t('salarie.forum.errorAction')
    }
}

async function supprimer(id: number) {
    if (!confirm(t('salarie.forum.confirmDelete'))) return
    try {
        const res = await fetch(`${API_BASE}/salarie/forum/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (!res.ok) throw new Error(t('salarie.forum.errorDelete'))
        threads.value = threads.value.filter(item => item.id !== id)
        if (showThreadModal.value && selectedThread.value?.id === id) {
            showThreadModal.value = false
        }
    } catch (e: any) {
        error.value = e.message || t('salarie.forum.errorDelete')
    }
}

function openBan(username: string) {
    banForm.username = username
    banForm.isPermanent = true
    banForm.durationHours = 24
    showBanModal.value = true
}

async function handleBan() {
    error.value = ''
    success.value = ''
    try {
        const res = await fetch(`${API_BASE}/salarie/forum/bannir`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({
                username: banForm.username,
                is_permanent: banForm.isPermanent,
                duration_hours: banForm.isPermanent ? 0 : Number(banForm.durationHours),
            }),
        })
        if (res.ok) {
            success.value = `L'utilisateur ${banForm.username} a été banni avec succès.`
            showBanModal.value = false
            setTimeout(() => { success.value = '' }, 4000)
        } else {
            error.value = "Erreur lors du bannissement de l'utilisateur."
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
}

async function openThread(thread: Thread) {
    selectedThread.value = thread
    posts.value = []
    loadingPosts.value = true
    showThreadModal.value = true
    try {
        const res = await fetch(`${API_BASE}/thread/${thread.id}/posts`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            posts.value = await res.json()
        } else {
            error.value = 'Impossible de charger les messages.'
        }
    } catch {
        error.value = 'Erreur réseau lors de la récupération des messages.'
    } finally {
        loadingPosts.value = false
    }
}

async function deletePost(postId: number) {
    if (!confirm('Voulez-vous vraiment supprimer ce message du forum ?')) return
    error.value = ''
    try {
        const res = await fetch(`${API_BASE}/post/${postId}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            posts.value = posts.value.filter(p => p.id !== postId)
            if (selectedThread.value) {
                selectedThread.value.replies = Math.max(0, selectedThread.value.replies - 1)
            }
        } else {
            error.value = 'Erreur lors de la suppression du message.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
}
</script>

<template>
    <div class="forum">
        <div class="page-header">
            <h1 class="page-title">{{ t('salarie.forum.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('salarie.forum.subtitle') }}</p>
        </div>

        <div v-if="error" class="error-banner">{{ error }}</div>
        <div v-if="success" class="success-banner">{{ success }}</div>

        <div v-if="loading" class="loading-state">{{ t('salarie.forum.loading') }}</div>

        <div v-else class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('salarie.forum.colDiscussion') }}</th>
                        <th>{{ t('salarie.forum.colAuthor') }}</th>
                        <th>{{ t('salarie.forum.colDate') }}</th>
                        <th>{{ t('salarie.forum.colReplies') }}</th>
                        <th>{{ t('salarie.forum.colActions') }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="threads.length === 0">
                        <td colspan="5" class="empty">{{ t('salarie.forum.empty') }}</td>
                    </tr>
                    <tr v-for="thread in threads" :key="thread.id" :class="{ 'row-pinned': thread.epingle }">
                        <td>
                            <div class="thread-titre-wrap">
                                <span v-if="thread.epingle" class="pin-icon" :title="t('salarie.forum.pinned')">📌</span>
                                <button class="btn-link-title" @click="openThread(thread)">{{ thread.titre }}</button>
                            </div>
                        </td>
                        <td class="td-muted">{{ thread.auteur }}</td>
                        <td class="td-muted">{{ thread.date }}</td>
                        <td>{{ thread.replies }}</td>
                        <td class="td-actions">
                            <button class="btn-action" @click="openThread(thread)">
                                Voir discussion
                            </button>
                            <button class="btn-action" :title="thread.epingle ? t('salarie.forum.unpin') : t('salarie.forum.pin')" @click="epingler(thread)">
                                {{ thread.epingle ? t('salarie.forum.unpin') : t('salarie.forum.pin') }}
                            </button>
                            <button class="btn-action btn-action--warn" :title="t('salarie.forum.ban')" @click="openBan(thread.auteur)">
                                {{ t('salarie.forum.ban') }}
                            </button>
                            <button class="btn-action btn-action--danger" :title="t('salarie.forum.delete')" @click="supprimer(thread.id)">
                                {{ t('salarie.forum.delete') }}
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div v-if="showThreadModal && selectedThread" class="modal-overlay" @click.self="showThreadModal = false">
            <div class="modal-card">
                <div class="modal-header">
                    <div>
                        <div class="modal-label">Modération discussion</div>
                        <h2 class="modal-title">{{ selectedThread.titre }}</h2>
                    </div>
                    <button class="btn-close" @click="showThreadModal = false">&times;</button>
                </div>
                <div class="modal-body">
                    <div v-if="loadingPosts" class="modal-loading">Chargement des messages...</div>
                    <div v-else-if="posts.length === 0" class="modal-empty">Aucun message dans cette discussion.</div>
                    <div v-else class="posts-list">
                        <div v-for="post in posts" :key="post.id" class="post-item">
                            <div class="post-meta">
                                <span class="post-author">{{ post.username || 'Utilisateur' }}</span>
                                <span class="post-date">{{ post.created_at }}</span>
                            </div>
                            <div class="post-content">{{ post.content }}</div>
                            <div class="post-actions">
                                <button class="btn-action-small btn-action-small--warn" @click="openBan(post.username || '')">
                                    Bannir l'auteur
                                </button>
                                <button class="btn-action-small btn-action-small--danger" @click="deletePost(post.id)">
                                    Supprimer le message
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button class="btn-cancel" @click="showThreadModal = false">Fermer</button>
                    <button class="btn-danger" @click="supprimer(selectedThread.id)">Supprimer la discussion entière</button>
                </div>
            </div>
        </div>

        <div v-if="showBanModal" class="modal-overlay" @click.self="showBanModal = false">
            <div class="modal-card modal-card--small">
                <div class="modal-header">
                    <h2 class="modal-title">Bannir l'utilisateur</h2>
                    <button class="btn-close" @click="showBanModal = false">&times;</button>
                </div>
                <form @submit.prevent="handleBan">
                    <div class="modal-body">
                        <p class="ban-target-desc">Vous êtes sur le point de bannir l'utilisateur <strong>{{ banForm.username }}</strong>.</p>

                        <div class="form-group">
                            <label class="form-label">Type de bannissement</label>
                            <div class="radio-group">
                                <label class="radio-label">
                                    <input type="radio" v-model="banForm.isPermanent" :value="true" />
                                    Bannissement définitif (permanent)
                                </label>
                                <label class="radio-label">
                                    <input type="radio" v-model="banForm.isPermanent" :value="false" />
                                    Bannissement temporaire
                                </label>
                            </div>
                        </div>

                        <div v-if="!banForm.isPermanent" class="form-group">
                            <label class="form-label">Durée du bannissement (en heures)</label>
                            <input
                                v-model.number="banForm.durationHours"
                                type="number"
                                class="form-input"
                                min="1"
                                placeholder="Ex: 24, 48, 168 (1 semaine)"
                                required
                            />
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn-cancel" @click="showBanModal = false">Annuler</button>
                        <button type="submit" class="btn-submit-ban">Confirmer le bannissement</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }

.error-banner { background: rgba(229, 62, 62, 0.08); border: 1.5px solid rgba(229, 62, 62, 0.25); border-radius: 8px; padding: 12px 16px; font-size: 0.85rem; color: #c53030; margin-bottom: 16px; font-weight: 600; }
.success-banner { background: var(--green-pale); border: 1.5px solid rgba(8, 106, 53, 0.25); border-radius: 8px; padding: 12px 16px; font-size: 0.85rem; color: var(--green-dark); margin-bottom: 16px; font-weight: 600; }

.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.row-pinned { background: rgba(215,236,225,0.2); }

.thread-titre-wrap { display: flex; align-items: center; gap: 8px; }
.btn-link-title { background: none; border: none; padding: 0; font-weight: 600; font-family: inherit; color: var(--charcoal); text-align: left; font-size: 0.9rem; cursor: pointer; transition: color 0.2s; }
.btn-link-title:hover { color: var(--green-dark); text-decoration: underline; }

.pin-icon { font-size: 0.85rem; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-actions { display: flex; gap: 8px; flex-wrap: wrap; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }

.btn-action { padding: 5px 11px; border-radius: 6px; font-size: 0.78rem; font-weight: 600; cursor: pointer; border: 1.5px solid rgba(53,53,53,0.15); background: transparent; color: var(--charcoal); transition: border-color 0.2s, color 0.2s, background 0.2s; white-space: nowrap; }
.btn-action:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-action--warn { border-color: rgba(234,179,8,0.4); color: #854d0e; }
.btn-action--warn:hover { border-color: #854d0e; background: rgba(234,179,8,0.05); }
.btn-action--danger { border-color: rgba(220,38,38,0.3); color: #dc2626; }
.btn-action--danger:hover { border-color: #dc2626; background: #fee2e2; }

.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.4); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: 1100; padding: 20px; }
.modal-card { background: var(--white); border-radius: 16px; width: 100%; max-width: 750px; max-height: 90vh; display: flex; flex-direction: column; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1), 0 10px 10px -5px rgba(0,0,0,0.04); }
.modal-card--small { max-width: 480px; }

.modal-header { padding: 20px 24px; border-bottom: 1px solid rgba(53,53,53,0.08); display: flex; justify-content: space-between; align-items: flex-start; }
.modal-label { font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.08em; font-weight: 700; color: var(--green-dark); margin-bottom: 4px; }
.modal-title { font-size: 1.25rem; font-weight: 800; color: var(--charcoal); margin: 0; letter-spacing: -0.02em; }
.btn-close { background: none; border: none; font-size: 1.5rem; cursor: pointer; color: rgba(53,53,53,0.4); }
.btn-close:hover { color: var(--charcoal); }

.modal-body { padding: 24px; overflow-y: auto; flex: 1; }
.modal-loading, .modal-empty { text-align: center; padding: 40px 0; opacity: 0.5; font-size: 0.9rem; }

.posts-list { display: flex; flex-direction: column; gap: 16px; }
.post-item { border: 1.5px solid rgba(53,53,53,0.08); border-radius: 10px; padding: 16px; background: rgba(53,53,53,0.01); }
.post-meta { display: flex; justify-content: space-between; font-size: 0.78rem; margin-bottom: 8px; }
.post-author { font-weight: 700; color: var(--green-dark); }
.post-date { opacity: 0.5; }
.post-content { font-size: 0.9rem; color: var(--charcoal); line-height: 1.5; margin-bottom: 12px; white-space: pre-wrap; }
.post-actions { display: flex; gap: 8px; justify-content: flex-end; }

.btn-action-small { padding: 4px 8px; border-radius: 4px; font-size: 0.72rem; font-weight: 600; cursor: pointer; border: 1px solid rgba(53,53,53,0.15); background: var(--white); }
.btn-action-small--warn { color: #854d0e; border-color: rgba(234,179,8,0.3); }
.btn-action-small--warn:hover { background: rgba(234,179,8,0.05); }
.btn-action-small--danger { color: #dc2626; border-color: rgba(220,38,38,0.2); }
.btn-action-small--danger:hover { background: #fee2e2; }

.modal-footer { padding: 20px 24px; border-top: 1px solid rgba(53,53,53,0.08); display: flex; justify-content: flex-end; gap: 12px; }
.btn-cancel { padding: 9px 16px; border-radius: 8px; border: 1.5px solid rgba(53,53,53,0.15); background: transparent; font-weight: 600; font-family: inherit; font-size: 0.85rem; cursor: pointer; }
.btn-cancel:hover { background: rgba(53,53,53,0.03); }
.btn-danger { padding: 9px 16px; border-radius: 8px; border: none; background: #dc2626; color: var(--white); font-weight: 700; font-family: inherit; font-size: 0.85rem; cursor: pointer; }
.btn-danger:hover { background: #b91c1c; }

.ban-target-desc { font-size: 0.92rem; margin-bottom: 20px; line-height: 1.5; }
.form-group { display: flex; flex-direction: column; gap: 8px; margin-bottom: 20px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); }
.radio-group { display: flex; flex-direction: column; gap: 10px; }
.radio-label { display: flex; align-items: center; gap: 8px; font-size: 0.88rem; cursor: pointer; }
.form-input { padding: 10px 12px; border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-family: inherit; font-size: 0.9rem; outline: none; }
.form-input:focus { border-color: var(--green-mid); }
.btn-submit-ban { padding: 10px 20px; border-radius: 8px; border: none; background: #dc2626; color: var(--white); font-weight: 700; font-family: inherit; font-size: 0.88rem; cursor: pointer; }
.btn-submit-ban:hover { background: #b91c1c; }
</style>
