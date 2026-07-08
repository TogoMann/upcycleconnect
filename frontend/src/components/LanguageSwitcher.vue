<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { API_BASE } from '@/config'
import { useAuthStore } from '@/stores/auth'

const { locale } = useI18n()
const authStore = useAuthStore()

async function setLocale(lang: string) {
    locale.value = lang
    localStorage.setItem('locale', lang)

    if (authStore.isAuthenticated && authStore.user) {
        try {
            await fetch(`${API_BASE}/users/${authStore.user.id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${authStore.token}`,
                },
                body: JSON.stringify({ language_preference: lang }),
            })
        } catch {
            // preference stays local-only if the sync fails
        }
    }
}
</script>

<template>
    <div class="lang-switcher">
        <button
            class="lang-btn"
            :class="{ 'lang-btn--active': locale === 'fr' }"
            @click="setLocale('fr')"
        >
            FR
        </button>
        <span class="lang-sep">/</span>
        <button
            class="lang-btn"
            :class="{ 'lang-btn--active': locale === 'en' }"
            @click="setLocale('en')"
        >
            EN
        </button>
    </div>
</template>

<style scoped>
.lang-switcher {
    display: inline-flex;
    align-items: center;
    gap: 6px;
}
.lang-btn {
    background: none;
    border: none;
    padding: 2px 4px;
    font-size: 0.8rem;
    font-weight: 600;
    color: inherit;
    opacity: 0.5;
    cursor: pointer;
    font-family: inherit;
    transition: opacity 0.2s;
}
.lang-btn:hover {
    opacity: 0.8;
}
.lang-btn--active {
    opacity: 1;
    text-decoration: underline;
}
.lang-sep {
    opacity: 0.3;
    font-size: 0.8rem;
}
</style>
