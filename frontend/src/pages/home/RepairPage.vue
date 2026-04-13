<script setup lang="ts">
import { ref } from 'vue'

const typeObjet = ref('')
const codePostal = ref('')

const typesObjets = [
    'Électronique',
    'Mobilier',
    'Vêtement',
    'Électroménager',
    'Jouet',
    'Outil',
    'Autre',
]

const codesPostaux = [
    '75001', '75002', '75003', '75004', '75005', '75006', '75007', '75008',
    '75009', '75010', '75011', '75012', '75013', '75014', '75015', '75016',
    '75017', '75018', '75019', '75020', '69001', '13001', '31000', '33000',
    '59000', '67000', '06000', '44000', '34000', '76000',
]

const handleRecherche = () => {
    console.log('Recherche :', { typeObjet: typeObjet.value, codePostal: codePostal.value })
}
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">Réparer vos objets.</h1>
            </div>
        </section>

        <section class="search-section">
            <div class="container">
                <div class="search-card">
                    <div class="search-card-header">
                        <h2 class="search-card-title">Trouver une ressource de réparation</h2>
                    </div>

                    <div class="search-card-body">
                        <div class="search-fields">
                            <div class="field-group">
                                <label class="field-label">Type d'objet</label>
                                <div class="select-wrap">
                                    <select v-model="typeObjet" class="field-select">
                                        <option value="" disabled>
                                            Électronique, Mobilier, Vêtement...
                                        </option>
                                        <option v-for="t in typesObjets" :key="t" :value="t">
                                            {{ t }}
                                        </option>
                                    </select>
                                    <span class="select-arrow">&#8964;</span>
                                </div>
                            </div>

                            <div class="field-group field-group--sm">
                                <label class="field-label">Code Postal</label>
                                <div class="select-wrap">
                                    <select v-model="codePostal" class="field-select">
                                        <option value="" disabled>75001</option>
                                        <option v-for="cp in codesPostaux" :key="cp" :value="cp">
                                            {{ cp }}
                                        </option>
                                    </select>
                                    <span class="select-arrow">&#8964;</span>
                                </div>
                            </div>
                        </div>

                        <button class="btn-search" @click="handleRecherche">Rechercher</button>
                    </div>
                </div>
            </div>
        </section>

        <div class="results-placeholder" />
    </div>
</template>

<style scoped>
.page-content {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}

.hero {
    padding: 72px 0 48px;
    text-align: center;
}
.hero-title {
    font-size: clamp(2.8rem, 6vw, 4.4rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0;
}

.search-section {
    padding: 0 0 40px;
}

.search-card {
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 24px rgba(8, 106, 53, 0.1);
}

.search-card-header {
    background: var(--green-mid);
    padding: 24px 28px;
}
.search-card-title {
    font-size: 1.35rem;
    font-weight: 700;
    color: var(--white);
    margin: 0;
    letter-spacing: -0.01em;
}

.search-card-body {
    background: var(--green-light);
    padding: 28px 28px 28px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.search-fields {
    display: flex;
    gap: 14px;
    align-items: flex-end;
}

.field-group {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.field-group--sm {
    flex: 0 0 220px;
}

.field-label {
    font-size: 1.15rem;
    font-weight: 700;
    color: rgba(255, 255, 255, 0.85);
    letter-spacing: -0.01em;
}

.select-wrap {
    position: relative;
    display: flex;
    align-items: center;
}
.field-select {
    width: 100%;
    padding: 16px 48px 16px 18px;
    font-size: 1.05rem;
    font-weight: 400;
    color: var(--charcoal);
    background: var(--cream);
    border: 2px solid rgba(255, 255, 255, 0.4);
    border-radius: 8px;
    appearance: none;
    -webkit-appearance: none;
    cursor: pointer;
    font-family: inherit;
    outline: none;
    transition:
        border-color 0.2s,
        box-shadow 0.2s;
}
.field-select:focus {
    border-color: var(--green-dark);
    box-shadow: 0 0 0 3px rgba(8, 106, 53, 0.15);
}
.field-select option[value=''][disabled] {
    color: rgba(53, 53, 53, 0.45);
}
.select-arrow {
    position: absolute;
    right: 16px;
    font-size: 1.3rem;
    color: var(--charcoal);
    pointer-events: none;
    line-height: 1;
}

.btn-search {
    width: 100%;
    padding: 22px;
    background: var(--green-mid);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 1.35rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    letter-spacing: -0.01em;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-search:hover {
    background: var(--green-dark);
    transform: translateY(-1px);
}
.btn-search:active {
    transform: translateY(0);
}

.results-placeholder {
    flex: 1;
    min-height: 120px;
}

@media (max-width: 700px) {
    .search-fields {
        flex-direction: column;
    }
    .field-group--sm {
        flex: none;
        width: 100%;
    }
    .hero-title {
        font-size: 2.4rem;
    }
}
</style>
