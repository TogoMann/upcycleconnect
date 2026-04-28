<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
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
    city_id: '',
    image_url: '',
})

const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)

function onFileChange(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) {
        if (file.size > 5 * 1024 * 1024) {
            errors.global = 'L\'image ne doit pas dépasser 5 Mo'
            return
        }
        imageFile.value = file
        imagePreview.value = URL.createObjectURL(file)
    }
}

const errors = reactive({
    name: '',
    description: '',
    price: '',
    category: '',
    city_id: '',
    global: '',
})

const submitting = ref(false)

onMounted(async () => {
    console.log('NouvelleAnnonce mounted, fetching cities...')
    try {
        await clientStore.fetchCities()
        console.log('Cities in store:', clientStore.cities)
    } catch (err) {
        console.error('Failed to fetch cities:', err)
    }
})

function validate(): boolean {
    errors.name = form.name.trim() ? '' : 'Le titre est requis'
    errors.description = form.description.trim() ? '' : 'La description est requise'
    errors.price = form.price && Number(form.price) > 0 ? '' : 'Un prix valide est requis'
    errors.category = form.category ? '' : 'La catégorie est requise'
    errors.city_id = form.city_id ? '' : 'La ville est requise'
    return !errors.name && !errors.description && !errors.price && !errors.category && !errors.city_id
}

async function handleSubmit() {
    if (!validate()) return
    submitting.value = true
    errors.global = ''
    try {
        if (imageFile.value) {
            form.image_url = await clientStore.uploadImage(imageFile.value)
        }

        await clientStore.createAnnonce({
            name: form.name.trim(),
            description: form.description.trim(),
            price: Number(form.price),
            category: form.category,
            city_id: Number(form.city_id),
            image_url: form.image_url,
        })
        await clientStore.fetchAnnonces()
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
                <label class="form-label">Ville</label>
                <select
                    v-model="form.city_id"
                    class="form-input"
                    :class="{ 'form-input--error': errors.city_id }"
                >
                    <option value="" disabled>Sélectionnez une ville</option>
                    <option v-for="city in clientStore.cities" :key="city.id" :value="city.id">
                        {{ city.name }} ({{ city.zip_code }})
                    </option>
                </select>
                <span v-if="errors.city_id" class="form-error">{{ errors.city_id }}</span>
            </div>

            <div class="form-group">
                <label class="form-label">Photo de l'objet (max 5 Mo)</label>
                <div class="image-upload-zone" :class="{ 'has-image': imagePreview }">
                    <input type="file" accept="image/*" class="file-input" @change="onFileChange" />
                    <div v-if="!imagePreview" class="upload-placeholder">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                            <circle cx="8.5" cy="8.5" r="1.5" />
                            <polyline points="21 15 16 10 5 21" />
                        </svg>
                        <span>Cliquez pour ajouter une photo</span>
                    </div>
                    <img v-else :src="imagePreview" class="preview-img" />
                    <button v-if="imagePreview" type="button" class="btn-remove-img" @click="imagePreview = null; imageFile = null">×</button>
                </div>
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
.image-upload-zone {
    position: relative;
    height: 180px;
    background: var(--cream);
    border: 1.5px dashed rgba(53, 53, 53, 0.2);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    transition: all 0.2s;
}
.image-upload-zone:hover {
    border-color: var(--green-mid);
}
.image-upload-zone.has-image {
    border-style: solid;
}
.file-input {
    position: absolute;
    inset: 0;
    opacity: 0;
    cursor: pointer;
    z-index: 2;
}
.upload-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    color: rgba(53, 53, 53, 0.4);
}
.upload-placeholder svg {
    width: 40px;
    height: 40px;
}
.upload-placeholder span {
    font-size: 0.85rem;
    font-weight: 500;
}
.preview-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.btn-remove-img {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 28px;
    height: 28px;
    background: rgba(229, 62, 62, 0.9);
    color: var(--white);
    border: none;
    border-radius: 50%;
    font-size: 1.2rem;
    line-height: 1;
    cursor: pointer;
    z-index: 3;
    display: flex;
    align-items: center;
    justify-content: center;
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
