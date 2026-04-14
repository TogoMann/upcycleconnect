<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
    nom: '',
    description: '',
    type: '',
    date_debut: '',
    budget: '',
})
const error = ref('')
const loading = ref(false)

async function submit() {
    if (!form.value.nom || !form.value.type) {
        error.value = 'Le nom et le type sont requis.'
        return
    }
    loading.value = true
    error.value = ''
    try {
        const res = await fetch('http://localhost:8081/pro/projets', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({ ...form.value, budget: parseFloat(form.value.budget) || 0 }),
        })
        if (res.ok) {
            router.push('/pro/projets')
        } else {
            const d = await res.json()
            error.value = d.message ?? 'Erreur lors de la création.'
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
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="15 18 9 12 15 6" />
                </svg>
                Retour
            </router-link>
            <h1 class="page-title">Nouveau projet.</h1>
            <p class="page-subtitle">Renseignez les informations de votre projet.</p>
        </div>

        <form class="form-card" @submit.prevent="submit">
            <div v-if="error" class="alert alert--error">{{ error }}</div>

            <div class="form-group">
                <label class="form-label">Nom du projet *</label>
                <input v-model="form.nom" type="text" class="form-input" placeholder="Ex: Upcycling meubles vintage" />
            </div>

            <div class="form-group">
                <label class="form-label">Description</label>
                <textarea v-model="form.description" class="form-input form-textarea" placeholder="Décrivez votre projet…" rows="4"></textarea>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">Type *</label>
                    <select v-model="form.type" class="form-input">
                        <option value="">Choisir un type</option>
                        <option value="mobilier">Mobilier</option>
                        <option value="textile">Textile</option>
                        <option value="electronique">Électronique</option>
                        <option value="autre">Autre</option>
                    </select>
                </div>
                <div class="form-group">
                    <label class="form-label">Date de début</label>
                    <input v-model="form.date_debut" type="date" class="form-input" />
                </div>
            </div>

            <div class="form-group">
                <label class="form-label">Budget estimé (€)</label>
                <input v-model="form.budget" type="number" min="0" step="0.01" class="form-input" placeholder="0.00" />
            </div>

            <div class="form-actions">
                <router-link to="/pro/projets" class="btn-secondary">Annuler</router-link>
                <button type="submit" class="btn-primary" :disabled="loading">
                    {{ loading ? 'Création…' : 'Créer le projet' }}
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
.form-card { max-width: 600px; background: var(--white); border-radius: 16px; border: 1.5px solid rgba(53,53,53,0.1); padding: 32px; display: flex; flex-direction: column; gap: 20px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; min-height: 100px; }
.form-actions { display: flex; gap: 12px; justify-content: flex-end; padding-top: 8px; }
.btn-primary { padding: 12px 24px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.btn-secondary { padding: 12px 24px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; text-decoration: none; transition: border-color 0.2s; display: inline-flex; align-items: center; }
.btn-secondary:hover { border-color: var(--charcoal); }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; background: #fee2e2; color: #991b1b; }
@media (max-width: 560px) { .form-row { grid-template-columns: 1fr; } }
</style>
