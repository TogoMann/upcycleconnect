<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
    <div class="error-page">
        <div class="error-content">
            <span class="error-code">404</span>
            <h1 class="error-title">Page introuvable.</h1>
            <p class="error-desc">
                La page que vous cherchez n'existe pas ou a été déplacée.
            </p>
            <div class="error-actions">
                <router-link to="/" class="btn-home">Retour à l'accueil</router-link>
                <button class="btn-back" @click="router.back()">Page précédente</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.error-page {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 80px 32px;
    text-align: center;
}

.error-content {
    max-width: 480px;
}

.error-code {
    display: block;
    font-size: clamp(5rem, 15vw, 9rem);
    font-weight: 800;
    color: var(--green-pale);
    line-height: 1;
    letter-spacing: -0.04em;
    margin-bottom: 16px;
}

.error-title {
    font-size: clamp(1.8rem, 4vw, 2.8rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.1;
    letter-spacing: -0.025em;
    margin: 0 0 16px;
}

.error-desc {
    font-size: 1rem;
    color: var(--charcoal);
    opacity: 0.65;
    line-height: 1.6;
    margin: 0 0 36px;
}

.error-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
    flex-wrap: wrap;
}

.btn-home {
    display: inline-block;
    padding: 13px 28px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    text-decoration: none;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-home:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}

.btn-back {
    display: inline-block;
    padding: 13px 28px;
    background: transparent;
    color: var(--charcoal);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        border-color 0.2s,
        color 0.2s;
}
.btn-back:hover {
    border-color: var(--green-mid);
    color: var(--green-mid);
}
</style>
