<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const form = ref({ titre: '', categorie: '', description: '', duree: '', statut: 'brouillon' })
const error = ref('')
const loading = ref(false)

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch(`http://localhost:8081/salarie/formations/${route.params.id}`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) {
            const data = await res.json()
            form.value = { titre: data.titre, categorie: data.categorie, description: data.description, duree: data.duree, statut: data.statut }
        } else {
            router.push('/salarie/formations')
        }
    } catch {}
})

async function submit() {
    loading.value = true
    error.value = ''
    try {
        const res = await fetch(`http://localhost:8081/salarie/formations/${route.params.id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify(form.value),
        })
        if (res.ok) router.push('/salarie/formations')
        else {
            const d = await res.json()
            error.value = d.message ?? 'Erreur lors de la mise à jour.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
    loading.value = false
}
</script>

<template>
    <div class="edit-formation">
        <div class="page-header">
            <router-link to="/salarie/formations" class="back-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="15 18 9 12 15 6" />
                </svg>
                Retour
            </router-link>
            <h1 class="page-title">Modifier la formation.</h1>
            <p class="page-subtitle">Mettez à jour les informations du module.</p>
        </div>

        <form class="form-card" @submit.prevent="submit">
            <div v-if="error" class="alert alert--error">{{ error }}</div>

            <div class="form-group">
                <label class="form-label">Titre</label>
                <input v-model="form.titre" type="text" class="form-input" />
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">Catégorie</label>
                    <select v-model="form.categorie" class="form-input">
                        <option value="">Choisir</option>
                        <option value="textile">Textile</option>
                        <option value="mobilier">Mobilier</option>
                        <option value="electronique">Électronique</option>
                        <option value="general">Général</option>
                    </select>
                </div>
                <div class="form-group">
                    <label class="form-label">Durée</label>
                    <input v-model="form.duree" type="text" class="form-input" />
                </div>
            </div>

            <div class="form-group">
                <label class="form-label">Description</label>
                <textarea v-model="form.description" class="form-input form-textarea" rows="5"></textarea>
            </div>

            <div class="form-group">
                <label class="form-label">Statut</label>
                <select v-model="form.statut" class="form-input">
                    <option value="brouillon">Brouillon</option>
                    <option value="publiee">Publiée</option>
                </select>
            </div>

            <div class="form-actions">
                <router-link to="/salarie/formations" class="btn-secondary">Annuler</router-link>
                <button type="submit" class="btn-primary" :disabled="loading">
                    {{ loading ? 'Enregistrement…' : 'Enregistrer' }}
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
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; min-height: 120px; }
.form-actions { display: flex; gap: 12px; justify-content: flex-end; padding-top: 8px; }
.btn-primary { padding: 12px 24px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.btn-secondary { padding: 12px 24px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; }
.btn-secondary:hover { border-color: var(--charcoal); }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; background: #fee2e2; color: #991b1b; }
@media (max-width: 560px) { .form-row { grid-template-columns: 1fr; } }
</style>
