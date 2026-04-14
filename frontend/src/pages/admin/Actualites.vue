<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface News {
    id: number
    created_by: number | null
    title: string
    content: string
    upvotes: number
    downvotes: number
    created_at: string | null
}

const news = ref<News[]>([])
const showForm = ref(false)
const form = ref({ title: '', content: '' })
const saving = ref(false)

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/news', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) news.value = await res.json()
    } catch {}
})

async function creer() {
    if (!form.value.title.trim() || !form.value.content.trim()) return
    saving.value = true
    try {
        const res = await fetch('http://localhost:8081/news/', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({ title: form.value.title, content: form.value.content }),
        })
        if (res.ok) {
            const created = await res.json()
            news.value.unshift(created)
            form.value = { title: '', content: '' }
            showForm.value = false
        }
    } catch {}
    saving.value = false
}

async function supprimer(id: number) {
    if (!confirm('Supprimer cette actualité ?')) return
    try {
        const res = await fetch(`http://localhost:8081/news/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) news.value = news.value.filter(n => n.id !== id)
    } catch {}
}

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
}
</script>

<template>
    <div class="actualites">
        <div class="page-header">
            <h1 class="page-title">Actualités.</h1>
            <p class="page-subtitle">Publiez et gérez les actualités de la plateforme.</p>
        </div>

        <div class="toolbar">
            <button class="btn-create" @click="showForm = !showForm">
                {{ showForm ? 'Annuler' : '+ Nouvelle actualité' }}
            </button>
        </div>

        <div v-if="showForm" class="form-card">
            <h3 class="form-title">Nouvelle actualité</h3>
            <form @submit.prevent="creer" class="form-fields">
                <div class="field">
                    <label class="field-label">Titre</label>
                    <input v-model="form.title" type="text" class="field-input" placeholder="Titre de l'actualité" required />
                </div>
                <div class="field">
                    <label class="field-label">Contenu</label>
                    <textarea v-model="form.content" class="field-textarea" rows="4" placeholder="Contenu…" required />
                </div>
                <button type="submit" class="btn-save" :disabled="saving">
                    {{ saving ? 'Publication…' : 'Publier' }}
                </button>
            </form>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Date</th>
                        <th>Votes</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="news.length === 0">
                        <td colspan="4" class="empty">Aucune actualité.</td>
                    </tr>
                    <tr v-for="n in news" :key="n.id">
                        <td class="td-bold">{{ n.title }}</td>
                        <td class="td-muted">{{ fmtDate(n.created_at) }}</td>
                        <td>▲ {{ n.upvotes }}</td>
                        <td class="td-actions">
                            <button class="btn-sm btn-sm--danger" @click="supprimer(n.id)">Supprimer</button>
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
.toolbar { margin-bottom: 16px; }
.btn-create { padding: 10px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; font-family: inherit; transition: background 0.2s; }
.btn-create:hover { background: var(--green-mid); }
.form-card { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); padding: 24px; margin-bottom: 20px; }
.form-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 20px; }
.form-fields { display: flex; flex-direction: column; gap: 16px; }
.field { display: flex; flex-direction: column; gap: 6px; }
.field-label { font-size: 0.78rem; font-weight: 600; color: var(--charcoal); opacity: 0.6; text-transform: uppercase; letter-spacing: 0.04em; }
.field-input, .field-textarea { padding: 9px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--white); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.field-textarea { resize: vertical; }
.field-input:focus, .field-textarea:focus { border-color: var(--green-mid); }
.btn-save { align-self: flex-start; padding: 9px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-save:disabled { opacity: 0.5; }
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
