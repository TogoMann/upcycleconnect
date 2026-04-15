<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from '@/stores/client'

const router = useRouter()
const clientStore = useClientStore()

const CATEGORIES = ['Mobilier', 'Décoration', 'Vêtements', 'Jouet', 'Electronique', 'Outils'] as const

const form = reactive({
    name: '',
    description: '',
    price: '',
    category: '' as string,
})

const errors = reactive({
    name: '',
    description: '',
    price: '',
    category: '',
    global: '',
})

const submitting = ref(false)

function validate(): boolean {
    errors.name = form.name.trim() ? '' : 'Le titre est requis'
    errors.description = form.description.trim() ? '' : 'La description est requise'
    errors.price = form.price && Number(form.price) > 0 ? '' : 'Un prix valide est requis'
    errors.category = form.category ? '' : 'La catégorie est requise'
    return !errors.name && !errors.description && !errors.price && !errors.category
}

async function handleSubmit() {
    if (!validate()) return
    submitting.value = true
    errors.global = ''
    try {
        await clientStore.createAnnonce({
            name: form.name.trim(),
            description: form.description.trim(),
            price: Number(form.price),
            category: form.category,
        })
        router.push('/particulier/annonces')
    } catch (e: any) {
        errors.global = e.message
    } finally {
        submitting.value = false
    }
}
</script>

<template>
    <div class="page">
        <div class="page-header">
            <router-link to="/particulier/annonces" class="back-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="19" y1="12" x2="5" y2="12" />
                    <polyline points="12 19 5 12 12 5" />
                </svg>
                Mes annonces
            </router-link>
            <h1 class="page-title">Nouvelle annonce.</h1>
        </div>

        <form class="form-card" @submit.prevent="handleSubmit">
            <div class="form-group">
                <label class="form-label">Titre de l'annonce</label>
                <input
                    v-model="form.name"
                    type="text"
                    class="form-input"
                    :class="{ 'form-input--error': errors.name }"
                    placeholder="ex. Chaise en bois restaurée"
                />
                <span v-if="errors.name" class="form-error">{{ errors.name }}</span>
            </div>

            <div class="form-group">
                <label class="form-label">Description</label>
                <textarea
                    v-model="form.description"
                    class="form-input form-textarea"
                    :class="{ 'form-input--error': errors.description }"
                    placeholder="Décrivez votre objet, son état, ses dimensions, son histoire…"
                    rows="5"
                ></textarea>
                <span v-if="errors.description" class="form-error">{{ errors.description }}</span>
            </div>

            <div class="form-group">
                <label class="form-label">Catégorie</label>
                <select
                    v-model="form.category"
                    class="form-input"
                    :class="{ 'form-input--error': errors.category }"
                >
                    <option value="" disabled>Sélectionnez une catégorie</option>
                    <option v-for="cat in CATEGORIES" :key="cat" :value="cat">{{ cat }}</option>
                </select>
                <span v-if="errors.category" class="form-error">{{ errors.category }}</span>
            </div>

            <div class="form-group">
                <label class="form-label">Prix (€)</label>
                <div class="input-prefix-wrap">
                    <span class="input-prefix">€</span>
                    <input
                        v-model="form.price"
                        type="number"
                        step="0.01"
                        min="0"
                        class="form-input form-input--prefixed"
                        :class="{ 'form-input--error': errors.price }"
                        placeholder="0.00"
                    />
                </div>
                <span v-if="errors.price" class="form-error">{{ errors.price }}</span>
            </div>

            <div v-if="errors.global" class="error-banner">{{ errors.global }}</div>

            <div class="form-actions">
                <router-link to="/particulier/annonces" class="btn-cancel">Annuler</router-link>
                <button type="submit" class="btn-submit" :disabled="submitting">
                    {{ submitting ? 'Publication…' : 'Publier l\'annonce' }}
                </button>
            </div>
        </form>
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
    max-width: 600px;
}

.page-header {
    margin-bottom: 32px;
}
.back-link {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 0.83rem;
    color: var(--green-mid);
    text-decoration: none;
    font-weight: 500;
    margin-bottom: 16px;
    transition: color 0.2s;
}
.back-link:hover {
    color: var(--green-dark);
}
.back-link svg {
    width: 16px;
    height: 16px;
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0;
    line-height: 1.08;
}

.form-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    padding: 32px;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.form-label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--charcoal);
}
.form-input {
    width: 100%;
    padding: 13px 16px;
    font-size: 0.9rem;
    font-family: inherit;
    color: var(--charcoal);
    background: var(--cream);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
}
.form-input::placeholder {
    color: rgba(53, 53, 53, 0.38);
}
.form-input:focus {
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.1);
}
.form-input--error {
    border-color: #e53e3e;
}
.form-textarea {
    resize: vertical;
    min-height: 120px;
}
.input-prefix-wrap {
    position: relative;
    display: flex;
    align-items: center;
}
.input-prefix {
    position: absolute;
    left: 14px;
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.45;
    pointer-events: none;
}
.form-input--prefixed {
    padding-left: 30px;
}
.form-error {
    font-size: 0.78rem;
    color: #e53e3e;
    font-weight: 500;
}
.error-banner {
    background: rgba(229, 62, 62, 0.08);
    border: 1px solid rgba(229, 62, 62, 0.25);
    border-radius: 8px;
    padding: 12px 16px;
    font-size: 0.83rem;
    color: #e53e3e;
}

.form-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    padding-top: 8px;
    border-top: 1px solid rgba(53, 53, 53, 0.08);
}
.btn-cancel {
    padding: 11px 20px;
    background: transparent;
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    color: rgba(53, 53, 53, 0.6);
    text-decoration: none;
    transition: border-color 0.2s, color 0.2s;
    cursor: pointer;
}
.btn-cancel:hover {
    border-color: rgba(53, 53, 53, 0.4);
    color: var(--charcoal);
}
.btn-submit {
    padding: 11px 24px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
}
.btn-submit:hover:not(:disabled) {
    background: var(--green-mid);
}
.btn-submit:disabled {
    opacity: 0.65;
    cursor: not-allowed;
}
</style>
