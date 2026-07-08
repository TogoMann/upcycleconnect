<script setup lang="ts">
withDefaults(defineProps<{ variant?: 'nav' | 'footer' }>(), {
    variant: 'nav',
})
</script>

<template>
    <span class="app-logo" :class="`app-logo--${variant}`">
        <svg class="app-logo-leaf" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path
                d="M4 20C4 20 3.5 13 8 8.5C12 4.5 20 4 20 4C20 4 20 12 15.5 16C11 20.5 4 20 4 20Z"
                fill="currentColor"
                opacity="0.18"
            />
            <path
                d="M4 20C4 20 3.5 13 8 8.5C12 4.5 20 4 20 4C20 4 20 12 15.5 16C11 20.5 4 20 4 20Z"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linejoin="round"
            />
            <path
                d="M5 19C5 19 9 13 17 7"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linecap="round"
            />
        </svg>
        <span class="app-logo-text">upcycleconnect</span>
    </span>
</template>

<style scoped>
.app-logo {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-weight: 800;
    letter-spacing: -0.01em;
    line-height: 1;
}
.app-logo-leaf {
    width: 1.35em;
    height: 1.35em;
    flex-shrink: 0;
}
.app-logo-text {
    white-space: nowrap;
}

.app-logo--nav {
    font-size: 1.1rem;
    color: var(--green-dark);
}
.app-logo--footer {
    font-size: 1.2rem;
    color: var(--white);
}
</style>
