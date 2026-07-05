<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted } from 'vue'
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

const threads = ref<Thread[]>([])
const loading = ref(true)
const error = ref('')

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
    } catch (e: any) {
        error.value = e.message || t('salarie.forum.errorDelete')
    }
}

async function bannir(auteur: string) {
    if (!confirm(t('salarie.forum.confirmBan', { name: auteur }))) return
    try {
        const res = await fetch(`${API_BASE}/salarie/forum/bannir`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({ username: auteur }),
        })
        if (!res.ok) throw new Error(t('salarie.forum.errorBan'))
    } catch (e: any) {
        error.value = e.message || t('salarie.forum.errorBan')
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
                            <div class="thread-titre">
                                <span v-if="thread.epingle" class="pin-icon" :title="t('salarie.forum.pinned')">📌</span>
                                {{ thread.titre }}
                            </div>
                        </td>
                        <td class="td-muted">{{ thread.auteur }}</td>
                        <td class="td-muted">{{ thread.date }}</td>
                        <td>{{ thread.replies }}</td>
                        <td class="td-actions">
                            <button class="btn-action" :title="thread.epingle ? t('salarie.forum.unpin') : t('salarie.forum.pin')" @click="epingler(thread)">
                                {{ thread.epingle ? t('salarie.forum.unpin') : t('salarie.forum.pin') }}
                            </button>
                            <button class="btn-action btn-action--warn" :title="t('salarie.forum.ban')" @click="bannir(thread.auteur)">
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
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.error-banner { background: rgba(229, 62, 62, 0.08); border: 1px solid rgba(229, 62, 62, 0.25); border-radius: 8px; padding: 12px 16px; font-size: 0.85rem; color: #c53030; margin-bottom: 16px; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.row-pinned { background: rgba(215,236,225,0.2); }
.thread-titre { display: flex; align-items: center; gap: 8px; font-weight: 600; }
.pin-icon { font-size: 0.85rem; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-actions { display: flex; gap: 8px; flex-wrap: wrap; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.btn-action { padding: 5px 11px; border-radius: 6px; font-size: 0.78rem; font-weight: 600; cursor: pointer; border: 1.5px solid rgba(53,53,53,0.15); background: transparent; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; white-space: nowrap; }
.btn-action:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-action--warn { border-color: rgba(234,179,8,0.4); color: #854d0e; }
.btn-action--warn:hover { border-color: #854d0e; }
.btn-action--danger { border-color: rgba(220,38,38,0.3); color: #dc2626; }
.btn-action--danger:hover { border-color: #dc2626; background: #fee2e2; }
</style>
