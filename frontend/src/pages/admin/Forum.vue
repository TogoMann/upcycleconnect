<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Thread {
    id: number
    created_by: number | null
    title: string
    content: string
    upvotes: number
    downvotes: number
    created_at: string | null
    last_post_at: string | null
}

const threads = ref<Thread[]>([])

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/thread', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) threads.value = await res.json()
    } catch {}
})

async function supprimerThread(id: number) {
    if (!confirm('Supprimer cette discussion ?')) return
    try {
        const res = await fetch(`http://localhost:8081/thread/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) threads.value = threads.value.filter(t => t.id !== id)
    } catch {}
}

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
}
</script>

<template>
    <div class="forum">
        <div class="page-header">
            <h1 class="page-title">Forum.</h1>
            <p class="page-subtitle">Modération des discussions.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Auteur</th>
                        <th>Votes</th>
                        <th>Date</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="threads.length === 0">
                        <td colspan="5" class="empty">Aucune discussion.</td>
                    </tr>
                    <tr v-for="t in threads" :key="t.id">
                        <td class="td-bold">{{ t.title }}</td>
                        <td class="td-muted">Utilisateur #{{ t.created_by ?? '—' }}</td>
                        <td>▲ {{ t.upvotes }} / ▼ {{ t.downvotes }}</td>
                        <td class="td-muted">{{ fmtDate(t.created_at) }}</td>
                        <td class="td-actions">
                            <button class="btn-sm btn-sm--danger" @click="supprimerThread(t.id)">Supprimer</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-actions { display: flex; gap: 8px; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.btn-sm { padding: 5px 12px; border-radius: 6px; font-size: 0.78rem; font-weight: 600; cursor: pointer; border: 1.5px solid rgba(53,53,53,0.15); background: transparent; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; white-space: nowrap; }
.btn-sm--danger { border-color: rgba(220,38,38,0.3); color: #dc2626; }
.btn-sm--danger:hover { border-color: #dc2626; background: #fee2e2; }
</style>
