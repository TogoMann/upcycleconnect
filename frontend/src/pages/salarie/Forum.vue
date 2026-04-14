<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

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

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/salarie/forum', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) threads.value = await res.json()
    } catch {}
})

async function epingler(t: Thread) {
    await fetch(`http://localhost:8081/salarie/forum/${t.id}/epingler`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    t.epingle = !t.epingle
}

async function supprimer(id: number) {
    if (!confirm('Supprimer cette discussion ?')) return
    await fetch(`http://localhost:8081/salarie/forum/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    threads.value = threads.value.filter(t => t.id !== id)
}

async function bannir(auteur: string) {
    if (!confirm(`Bannir l'utilisateur "${auteur}" ?`)) return
    await fetch(`http://localhost:8081/salarie/forum/bannir`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
        body: JSON.stringify({ username: auteur }),
    })
}
</script>

<template>
    <div class="forum">
        <div class="page-header">
            <h1 class="page-title">Forum.</h1>
            <p class="page-subtitle">Modération : supprimer, épingler, bannir des utilisateurs.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Discussion</th>
                        <th>Auteur</th>
                        <th>Date</th>
                        <th>Réponses</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="threads.length === 0">
                        <td colspan="5" class="empty">Aucune discussion.</td>
                    </tr>
                    <tr v-for="t in threads" :key="t.id" :class="{ 'row-pinned': t.epingle }">
                        <td>
                            <div class="thread-titre">
                                <span v-if="t.epingle" class="pin-icon" title="Épinglé">📌</span>
                                {{ t.titre }}
                            </div>
                        </td>
                        <td class="td-muted">{{ t.auteur }}</td>
                        <td class="td-muted">{{ t.date }}</td>
                        <td>{{ t.replies }}</td>
                        <td class="td-actions">
                            <button class="btn-action" :title="t.epingle ? 'Désépingler' : 'Épingler'" @click="epingler(t)">
                                {{ t.epingle ? 'Désépingler' : 'Épingler' }}
                            </button>
                            <button class="btn-action btn-action--warn" title="Bannir" @click="bannir(t.auteur)">
                                Bannir
                            </button>
                            <button class="btn-action btn-action--danger" title="Supprimer" @click="supprimer(t.id)">
                                Supprimer
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
