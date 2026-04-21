<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

interface Annonce {
    id: any
    name: string
    description: string
    price: any
    category: string
    created_by_name: string
    created_at: any
    approved: boolean
    status: string
}

const annonce = ref<Annonce | null>(null)
const loading = ref(true)

function formatPrice(p: any): string {
    if (!p) return '0.00'
    const val = typeof p === 'object' ? p.Float64 ?? p.Int64 : p
    return Number(val).toFixed(2)
}

function formatDate(ts: any): string {
    if (!ts) return '—'
    const date = new Date(ts.Time ?? ts)
    return date.toLocaleDateString('fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
}

onMounted(async () => {
    try {
        const res = await fetch(`http://localhost:8081/admin/listings/${route.params.id}`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            annonce.value = await res.json()
        } else {
            router.push('/admin/annonces')
        }
    } catch (e) {
        console.error('Fetch error:', e)
    } finally {
        loading.value = false
    }
})

async function valider() {
    try {
        const res = await fetch(`http://localhost:8081/listing/${route.params.id}/approve`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok && annonce.value) {
            annonce.value.approved = true
        }
    } catch (e) {
        console.error('Approve error:', e)
    }
}

async function refuser() {
    if (!confirm('Désapprouver cette annonce ?')) return
    try {
        const res = await fetch(`http://localhost:8081/listing/${route.params.id}/disapprove`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok && annonce.value) {
            annonce.value.approved = false
        }
    } catch (e) {
        console.error('Disapprove error:', e)
    }
}
</script>

<template>
    <div class="annonce-detail">
        <router-link to="/admin/annonces" class="back-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15 18 9 12 15 6" />
            </svg>
            Retour
        </router-link>

        <div v-if="loading" class="loading">Chargement…</div>

        <template v-else-if="annonce">
            <div class="page-header">
                <div class="header-row">
                    <h1 class="page-title">{{ annonce.name }}.</h1>
                    <span
                        class="badge"
                        :class="annonce.approved ? 'badge--active' : 'badge--pending'"
                    >
                        {{ annonce.approved ? 'Approuvée' : 'En attente' }}
                    </span>
                </div>
            </div>

            <div class="info-grid">
                <div class="info-card">
                    <div class="info-label">Auteur</div>
                    <div class="info-value">{{ annonce.created_by_name || 'Inconnu' }}</div>
                </div>
                <div class="info-card">
                    <div class="info-label">Catégorie</div>
                    <div class="info-value">{{ annonce.category }}</div>
                </div>
                <div class="info-card">
                    <div class="info-label">Prix</div>
                    <div class="info-value">{{ formatPrice(annonce.price) }} €</div>
                </div>
                <div class="info-card">
                    <div class="info-label">Date</div>
                    <div class="info-value">{{ formatDate(annonce.created_at) }}</div>
                </div>
            </div>

            <div class="desc-section">
                <h3 class="section-title">Description</h3>
                <p class="desc-text">{{ annonce.description }}</p>
            </div>

            <div class="actions-row">
                <button v-if="!annonce.approved" class="btn-validate" @click="valider">Approuver l'annonce</button>
                <button v-else class="btn-refuse" @click="refuser">Désapprouver</button>
            </div>
        </template>
    </div>
</template>

<style scoped>
.back-link { display: inline-flex; align-items: center; gap: 6px; font-size: 0.85rem; color: var(--green-mid); text-decoration: none; margin-bottom: 24px; transition: color 0.2s; }
.back-link:hover { color: var(--green-dark); }
.back-link svg { width: 16px; height: 16px; }
.page-header { margin-bottom: 28px; }
.header-row { display: flex; align-items: flex-start; gap: 16px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0; line-height: 1.08; flex: 1; }
.badge { display: inline-block; padding: 5px 12px; border-radius: 20px; font-size: 0.78rem; font-weight: 600; white-space: nowrap; align-self: center; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--danger { background: #fee2e2; color: #991b1b; }
.badge--pending { background: #fef3c7; color: #92400e; }
.info-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-bottom: 32px; }
.info-card { background: var(--green-pale); border-radius: 12px; padding: 18px; }
.info-label { font-size: 0.75rem; font-weight: 600; color: var(--green-dark); opacity: 0.65; text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 6px; }
.info-value { font-size: 1rem; font-weight: 700; color: var(--charcoal); }
.section-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 12px; }
.desc-section { background: var(--white); border: 1.5px solid rgba(53,53,53,0.08); border-radius: 12px; padding: 24px; margin-bottom: 32px; }
.desc-text { font-size: 0.9rem; color: var(--charcoal); line-height: 1.6; margin: 0; opacity: 0.8; }
.actions-row { display: flex; gap: 12px; }
.btn-validate { padding: 12px 28px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-validate:hover { background: var(--green-mid); }
.btn-refuse { padding: 12px 28px; background: transparent; color: #dc2626; border: 1.5px solid #dc2626; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-refuse:hover { background: #fee2e2; }
.loading { opacity: 0.5; font-size: 0.9rem; padding: 40px 0; }
@media (max-width: 700px) { .info-grid { grid-template-columns: repeat(2, 1fr); } }
</style>
