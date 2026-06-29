<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

interface Listing {
    id: { Int64: number; Valid: boolean }
    name: string
}

const listings = ref<Listing[]>([])
const form = ref({
    title: '',
    description: '',
    listing_id: 0,
})
const error = ref('')
const loading = ref(false)

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/listing`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            const data = await res.json()
            listings.value = Array.isArray(data) ? data : []
        }
    } catch {}
})

async function submit() {
    if (!form.value.title.trim()) {
        error.value = 'Le titre du projet est requis.'
        return
    }
    if (!form.value.description.trim()) {
        error.value = 'La description est requise.'
        return
    }
    loading.value = true
    error.value = ''
    try {
        const body: any = {
            title: form.value.title.trim(),
            description: form.value.description.trim(),
        }
        if (form.value.listing_id > 0) {
            body.listing_id = { Int64: form.value.listing_id, Valid: true }
        }
        const res = await fetch(`${API_BASE}/project`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify(body),
        })
        if (res.ok) {
            const data = await res.json()
            const projectId = data.id?.Int64 || data.id
            router.push(`/pro/projets/${projectId}`)
        } else {
            const d = await res.text()
            error.value = d || 'Erreur lors de la création.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
    loading.value = false
}
</script>

<template>
    <div class="nouveau-projet">
        <div class="page-header">
            <router-link to="/pro/projets" class="back-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
                Retour aux projets
            </router-link>
            <h1 class="page-title">Nouveau projet.</h1>
            <p class="page-subtitle">Créez un projet d'upcycling et documentez chaque étape de transformation.</p>
        </div>

        <form class="form-card" @submit.prevent="submit">
            <div v-if="error" class="alert alert--error">{{ error }}</div>

            <div class="form-group">
                <label class="form-label">Titre du projet *</label>
                <input v-model="form.title" type="text" class="form-input" placeholder="Ex: Transformation d'une palette en étagère murale" maxlength="128" />
                <p class="form-hint">{{ form.title.length }}/128 caractères</p>
            </div>

            <div class="form-group">
                <label class="form-label">Description *</label>
                <textarea v-model="form.description" class="form-input form-textarea" placeholder="Décrivez votre projet d'upcycling, les matériaux utilisés, l'objectif final..." rows="5"></textarea>
            </div>

            <div class="form-group">
                <label class="form-label">Annonce associée</label>
                <select v-model.number="form.listing_id" class="form-input">
                    <option :value="0">Aucune annonce liée</option>
                    <option v-for="l in listings" :key="l.id.Int64" :value="l.id.Int64">{{ l.name }}</option>
                </select>
                <p class="form-hint">Optionnel — Liez ce projet à une annonce existante sur la plateforme</p>
            </div>

            <div class="form-info">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
                <p>Après la création, vous pourrez ajouter les étapes détaillées de votre projet, modifier le statut, et le soumettre pour une mise en avant sur la plateforme.</p>
            </div>

            <div class="form-actions">
                <router-link to="/pro/projets" class="btn-secondary">Annuler</router-link>
                <button type="submit" class="btn-primary" :disabled="loading">
                    {{ loading ? 'Création...' : 'Créer le projet' }}
                </button>
            </div>
        </form>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.back-link { display: inline-flex; align-items: center; gap: 6px; font-size: 0.85rem; color: var(--green-mid); text-decoration: none; margin-bottom: 16px; transition: color 0.2s; }
.back-link:hover { color: var(--green-dark); }
.back-link svg { width: 16px; height: 16px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }

.form-card { max-width: 640px; background: var(--white); border-radius: 16px; border: 1.5px solid rgba(53,53,53,0.1); padding: 32px; display: flex; flex-direction: column; gap: 22px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; min-height: 120px; line-height: 1.6; }
.form-hint { font-size: 0.78rem; color: var(--charcoal); opacity: 0.4; margin: 0; }

.form-info { display: flex; gap: 12px; padding: 16px; background: #eff6ff; border-radius: 10px; align-items: flex-start; }
.form-info svg { width: 20px; height: 20px; flex-shrink: 0; color: #3b82f6; margin-top: 1px; }
.form-info p { font-size: 0.84rem; color: #1e40af; line-height: 1.5; margin: 0; }

.form-actions { display: flex; gap: 12px; justify-content: flex-end; padding-top: 8px; }
.btn-primary { padding: 12px 24px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.btn-secondary { padding: 12px 24px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; text-decoration: none; transition: border-color 0.2s; display: inline-flex; align-items: center; }
.btn-secondary:hover { border-color: var(--charcoal); }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.86rem; font-weight: 500; }
.alert--error { background: #fee2e2; color: #991b1b; }

@media (max-width: 560px) { .form-card { padding: 24px 20px; } }
</style>
