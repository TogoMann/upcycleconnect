<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const form = ref({
    titre: '',
    categorie: '',
    description: '',
    statut: 'brouillon',
    date: '',
    start_time: '',
    end_time: '',
    prix: '',
    max_capacity: '',
})

const dureeCalculee = computed(() => {
    if (!form.value.start_time || !form.value.end_time) return ''
    const [h1, m1] = form.value.start_time.split(':').map(Number)
    const [h2, m2] = form.value.end_time.split(':').map(Number)
    const minutes = (h2 * 60 + m2) - (h1 * 60 + m1)
    if (minutes <= 0) return ''
    const h = Math.floor(minutes / 60)
    const m = minutes % 60
    return m === 0 ? `${h}h` : `${h}h${String(m).padStart(2, '0')}`
})
const error = ref('')
const loading = ref(false)

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch(`${API_BASE}/salarie/formations/${route.params.id}`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) {
            const data = await res.json()
            form.value = {
                titre: data.name || data.nom || '',
                categorie: data.categorie || 'general',
                description: data.description || '',
                statut: data.approved ? 'publiee' : 'brouillon',
                date: data.date ? data.date.substring(0, 10) : '',
                start_time: data.start_time ? data.start_time.substring(0, 5) : '',
                end_time: data.end_time ? data.end_time.substring(0, 5) : '',
                prix: data.price || '',
                max_capacity: data.max_capacity || '',
            }
        } else {
            router.push('/salarie/formations')
        }
    } catch {}
})

async function submit() {
    loading.value = true
    error.value = ''
    try {
        const payload = {
            ...form.value,
            duree: dureeCalculee.value,
            prix: parseFloat(form.value.prix) || 0,
            max_capacity: parseInt(form.value.max_capacity) || null,
        }
        const res = await fetch(`${API_BASE}/salarie/formations/${route.params.id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify(payload),
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
                    <input :value="dureeCalculee || '—'" type="text" class="form-input" readonly />
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">Date</label>
                    <input v-model="form.date" type="date" class="form-input" />
                </div>
                <div class="form-row" style="display: grid; grid-template-columns: 1fr 1fr; gap: 8px;">
                    <div class="form-group">
                        <label class="form-label">Début</label>
                        <input v-model="form.start_time" type="time" class="form-input" />
                    </div>
                    <div class="form-group">
                        <label class="form-label">Fin</label>
                        <input v-model="form.end_time" type="time" class="form-input" />
                    </div>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">Prix (€)</label>
                    <input v-model="form.prix" type="number" step="0.01" class="form-input" placeholder="0.00" />
                </div>
                <div class="form-group">
                    <label class="form-label">Capacité maximale</label>
                    <input v-model="form.max_capacity" type="number" class="form-input" placeholder="Ex: 15" />
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
