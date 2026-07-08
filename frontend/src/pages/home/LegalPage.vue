<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
    <div class="page-content">
        <div class="container">
            <h1 class="page-title">{{ t('legal.pageTitle') }}</h1>

            <section id="mentions" class="legal-section">
                <h2 class="section-title">{{ t('legal.mentions.title') }}</h2>
                <p class="section-text">{{ t('legal.mentions.editor') }}</p>
                <p class="section-text">{{ t('legal.mentions.host') }}</p>
                <p class="section-text">{{ t('legal.mentions.contact') }}</p>
            </section>

            <section id="cgu" class="legal-section">
                <h2 class="section-title">{{ t('legal.terms.title') }}</h2>
                <p class="section-text">{{ t('legal.terms.intro') }}</p>
                <p class="section-text">{{ t('legal.terms.usage') }}</p>
                <p class="section-text">{{ t('legal.terms.liability') }}</p>
            </section>

            <section id="confidentialite" class="legal-section">
                <h2 class="section-title">{{ t('legal.privacy.title') }}</h2>
                <p class="section-text">{{ t('legal.privacy.data') }}</p>
                <p class="section-text">{{ t('legal.privacy.rights') }}</p>
                <p class="section-text">{{ t('legal.privacy.contact') }}</p>
            </section>
        </div>
    </div>
</template>

<style scoped>
.page-content {
    flex: 1;
    padding: 56px 0 80px;
}
.container {
    max-width: 780px;
    margin: 0 auto;
    padding: 0 32px;
}
.page-title {
    font-size: clamp(2rem, 4vw, 2.8rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 40px;
}
.legal-section {
    margin-bottom: 40px;
}
.section-title {
    font-size: 1.3rem;
    font-weight: 700;
    color: var(--green-dark);
    margin: 0 0 16px;
}
.section-text {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.8;
    line-height: 1.7;
    margin: 0 0 12px;
}
</style>
