<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const title = ref('')
const category = ref('Communauté')
const content = ref('')
const sending = ref(false)
const error = ref('')

const categories = ['Bricolage', 'Textile', 'Ressources', 'Débutants', 'Communauté']

async function handleSubmit() {
    if (!title.value.trim() || !content.value.trim() || !authStore.isAuthenticated) return
    sending.value = true
    error.value = ''
    try {
        const res = await fetch('http://localhost:8081/thread', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({
                title: title.value,
                category: category.value,
                content: content.value,
            }),
        })
        if (res.ok) {
            const data = await res.json()
            router.push(`/forum/${data.id}`)
        } else {
            error.value = "Erreur lors de la création du sujet."
        }
    } catch {
        error.value = "Erreur réseau."
    }
    sending.value = false
}
</script>

<template>
    <div class="page-content">
        <section class="breadcrumb-bar">
            <div class="container">
                <router-link to="/forum" class="breadcrumb-link">Forum</router-link>
                <span class="breadcrumb-sep">›</span>
                <span class="breadcrumb-current">Nouveau sujet</span>
            </div>
        </section>

        <section class="form-section">
            <div class="container">
                <h1 class="page-title">Créer un nouveau sujet</h1>

                <form class="create-thread-form" @submit.prevent="handleSubmit">
                    <div class="form-group">
                        <label for="title">Titre du sujet</label>
                        <input
                            id="title"
                            v-model="title"
                            type="text"
                            placeholder="De quoi souhaitez-vous discuter ?"
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label for="category">Catégorie</label>
                        <select id="category" v-model="category">
                            <option v-for="cat in categories" :key="cat" :value="cat">
                                {{ cat }}
                            </option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label for="content">Contenu</label>
                        <textarea
                            id="content"
                            v-model="content"
                            rows="10"
                            placeholder="Détaillez votre question ou votre partage..."
                            required
                        ></textarea>
                    </div>

                    <div v-if="error" class="error-msg">{{ error }}</div>

                    <div class="form-actions">
                        <button type="button" class="btn-cancel" @click="router.back()">Annuler</button>
                        <button type="submit" class="btn-submit" :disabled="sending">
                            {{ sending ? 'Création...' : 'Publier le sujet' }}
                        </button>
                    </div>
                </form>
            </div>
        </section>
    </div>
</template>

<style scoped>
.page-content { flex: 1; display: flex; flex-direction: column; }
.container { max-width: 800px; margin: 0 auto; padding: 0 32px; width: 100%; box-sizing: border-box; }
.breadcrumb-bar { padding: 20px 0; border-bottom: 1px solid rgba(53,53,53,0.08); }
.breadcrumb-bar .container { display: flex; align-items: center; gap: 8px; font-size: 0.85rem; }
.breadcrumb-link { color: var(--green-mid); text-decoration: none; }
.breadcrumb-sep { color: var(--charcoal); opacity: 0.4; }
.breadcrumb-current { color: var(--charcoal); opacity: 0.7; }

.form-section { padding: 40px 0 80px; }
.page-title { font-size: 2rem; font-weight: 800; color: var(--charcoal); margin: 0 0 32px; letter-spacing: -0.02em; }

.create-thread-form { display: flex; flex-direction: column; gap: 24px; background: var(--white); padding: 32px; border: 1.5px solid rgba(53,53,53,0.1); border-radius: 12px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-group label { font-size: 0.9rem; font-weight: 700; color: var(--charcoal); }
.form-group input, .form-group select, .form-group textarea {
    padding: 12px 16px;
    font-size: 0.95rem;
    font-family: inherit;
    border: 1.5px solid rgba(53,53,53,0.2);
    border-radius: 8px;
    outline: none;
    transition: border-color 0.2s;
}
.form-group input:focus, .form-group select:focus, .form-group textarea:focus { border-color: var(--green-mid); }
.error-msg { color: #dc2626; font-size: 0.9rem; font-weight: 600; }

.form-actions { display: flex; justify-content: flex-end; gap: 16px; margin-top: 8px; }
.btn-cancel { background: none; border: none; color: var(--charcoal); opacity: 0.6; font-weight: 600; cursor: pointer; padding: 12px 20px; }
.btn-submit { background: var(--green-dark); color: var(--white); border: none; padding: 12px 28px; border-radius: 8px; font-weight: 700; cursor: pointer; transition: background 0.2s; }
.btn-submit:hover { background: var(--green-mid); }
.btn-submit:disabled { opacity: 0.5; }
</style>
