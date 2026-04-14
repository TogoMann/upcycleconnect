<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const form = ref({
    first_name: '',
    last_name: '',
    email: '',
    entreprise: '',
    siret: '',
    telephone: '',
    adresse: '',
})
const success = ref(false)
const error = ref('')
const loading = ref(false)

onMounted(() => {
    if (authStore.user) {
        form.value.first_name = authStore.user.first_name ?? ''
        form.value.last_name = authStore.user.last_name ?? ''
        form.value.email = authStore.user.email ?? ''
    }
})

async function save() {
    loading.value = true
    error.value = ''
    success.value = false
    try {
        const res = await fetch('http://localhost:8081/pro/profil', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify(form.value),
        })
        if (res.ok) {
            success.value = true
            await authStore.fetchCurrentUser()
        } else {
            const d = await res.json()
            error.value = d.message ?? 'Erreur lors de la sauvegarde.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
    loading.value = false
}
</script>

<template>
    <div class="profil">
        <div class="page-header">
            <h1 class="page-title">Mon Profil.</h1>
            <p class="page-subtitle">Informations personnelles et entreprise.</p>
        </div>

        <form class="form-card" @submit.prevent="save">
            <div v-if="success" class="alert alert--success">Profil mis à jour.</div>
            <div v-if="error" class="alert alert--error">{{ error }}</div>

            <h3 class="form-section-title">Informations personnelles</h3>
            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">Prénom</label>
                    <input v-model="form.first_name" type="text" class="form-input" />
                </div>
                <div class="form-group">
                    <label class="form-label">Nom</label>
                    <input v-model="form.last_name" type="text" class="form-input" />
                </div>
            </div>
            <div class="form-group">
                <label class="form-label">Email</label>
                <input v-model="form.email" type="email" class="form-input" />
            </div>
            <div class="form-group">
                <label class="form-label">Téléphone</label>
                <input v-model="form.telephone" type="tel" class="form-input" />
            </div>

            <h3 class="form-section-title">Entreprise</h3>
            <div class="form-group">
                <label class="form-label">Nom de l'entreprise</label>
                <input v-model="form.entreprise" type="text" class="form-input" />
            </div>
            <div class="form-group">
                <label class="form-label">SIRET</label>
                <input v-model="form.siret" type="text" class="form-input" placeholder="00000000000000" />
            </div>
            <div class="form-group">
                <label class="form-label">Adresse</label>
                <input v-model="form.adresse" type="text" class="form-input" />
            </div>

            <div class="form-actions">
                <button type="submit" class="btn-primary" :disabled="loading">
                    {{ loading ? 'Enregistrement…' : 'Enregistrer' }}
                </button>
            </div>
        </form>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.form-card { max-width: 600px; background: var(--white); border-radius: 16px; border: 1.5px solid rgba(53,53,53,0.1); padding: 32px; display: flex; flex-direction: column; gap: 16px; }
.form-section-title { font-size: 0.9rem; font-weight: 700; color: var(--charcoal); text-transform: uppercase; letter-spacing: 0.06em; margin: 8px 0 4px; opacity: 0.6; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-actions { padding-top: 8px; }
.btn-primary { padding: 12px 28px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; }
.alert--success { background: var(--green-pale); color: var(--green-dark); }
.alert--error { background: #fee2e2; color: #991b1b; }
@media (max-width: 560px) { .form-row { grid-template-columns: 1fr; } }
</style>
