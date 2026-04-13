<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const thread = {
    id: Number(route.params.id),
    titre: 'Comment upcycler des palettes de bois en mobilier de jardin ?',
    categorie: 'Bricolage',
    auteur: 'Sophie M.',
    avatar: 'SM',
    date: '12 avr. 2026 à 14h32',
    contenu: `Bonjour à tous ! Je cherche des conseils pour transformer des palettes de bois récupérées en mobilier de jardin. J'en ai une dizaine de disponibles et j'aimerais créer une table basse, quelques chaises et peut-être un canapé d'extérieur.

Quelqu'un a-t-il déjà réalisé ce type de projet ? Quels outils sont vraiment indispensables ? Et surtout, comment traiter le bois pour qu'il résiste aux intempéries sans utiliser de produits chimiques trop agressifs pour l'environnement ?

Merci d'avance pour vos retours !`,
}

const reponses = [
    {
        id: 1,
        auteur: 'Thomas L.',
        avatar: 'TL',
        date: '12 avr. 2026 à 16h10',
        contenu: `Super projet ! J'ai fait la même chose l'été dernier. Mon conseil principal : vérifiez bien les marquages sur les palettes. Préférez les palettes estampillées "HT" (traitement thermique) et évitez celles marquées "MB" (bromure de méthyle).

Pour le traitement du bois, j'ai utilisé de l'huile de lin qui protège efficacement et reste naturelle. Deux ou trois couches suffisent.`,
    },
    {
        id: 2,
        auteur: 'Camille R.',
        avatar: 'CR',
        date: '12 avr. 2026 à 18h45',
        contenu: `Pour les outils, le minimum syndical : une ponceuse orbitale (les palettes ont souvent des échardes), une scie circulaire, une visseuse et du papier de verre en différents grains. Comptez aussi des vis en inox pour la résistance à la rouille en extérieur.`,
    },
    {
        id: 3,
        auteur: 'Antoine V.',
        avatar: 'AV',
        date: '13 avr. 2026 à 09h20',
        contenu: `N'oubliez pas de laisser les palettes sécher au minimum 2-3 semaines avant de les travailler si elles ont été exposées à la pluie. Un bois humide travaille et vos assemblages pourraient se déformer. Bon courage pour le projet, hâte de voir le résultat !`,
    },
]

const newReply = reactive({ contenu: '' })

function handleReply() {
    console.log('new reply:', newReply.contenu)
    newReply.contenu = ''
}
</script>

<template>
    <div class="page-content">
        <section class="breadcrumb-bar">
            <div class="container">
                <router-link to="/forum" class="breadcrumb-link">Forum</router-link>
                <span class="breadcrumb-sep">›</span>
                <span class="breadcrumb-current">{{ thread.titre }}</span>
            </div>
        </section>

        <section class="thread-section">
            <div class="container">
                <div class="thread-header">
                    <span class="badge-cat">{{ thread.categorie }}</span>
                    <h1 class="thread-title">{{ thread.titre }}</h1>
                </div>

                <div class="post post--original">
                    <div class="post-author">
                        <div class="post-avatar">{{ thread.avatar }}</div>
                        <div class="post-author-info">
                            <span class="post-author-name">{{ thread.auteur }}</span>
                            <span class="post-date">{{ thread.date }}</span>
                        </div>
                    </div>
                    <div class="post-body">
                        <p v-for="(paragraph, i) in thread.contenu.split('\n\n')" :key="i" class="post-paragraph">
                            {{ paragraph }}
                        </p>
                    </div>
                </div>

                <div class="replies-section">
                    <h2 class="replies-title">{{ reponses.length }} réponses</h2>

                    <div v-for="reponse in reponses" :key="reponse.id" class="post">
                        <div class="post-author">
                            <div class="post-avatar">{{ reponse.avatar }}</div>
                            <div class="post-author-info">
                                <span class="post-author-name">{{ reponse.auteur }}</span>
                                <span class="post-date">{{ reponse.date }}</span>
                            </div>
                        </div>
                        <div class="post-body">
                            <p v-for="(paragraph, i) in reponse.contenu.split('\n\n')" :key="i" class="post-paragraph">
                                {{ paragraph }}
                            </p>
                        </div>
                    </div>
                </div>

                <div class="reply-form-section">
                    <h2 class="reply-form-title">Répondre à ce sujet</h2>
                    <form class="reply-form" @submit.prevent="handleReply">
                        <textarea
                            v-model="newReply.contenu"
                            placeholder="Partagez votre expérience ou vos conseils..."
                            class="reply-textarea"
                            rows="6"
                            required
                        />
                        <div class="reply-form-footer">
                            <router-link to="/auth/login" class="reply-login-hint">
                                Connectez-vous pour répondre
                            </router-link>
                            <button type="submit" class="btn-reply">Publier la réponse</button>
                        </div>
                    </form>
                </div>
            </div>
        </section>
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

.breadcrumb-bar {
    padding: 20px 0;
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
}
.breadcrumb-bar .container {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
}
.breadcrumb-link {
    color: var(--green-mid);
    text-decoration: none;
    transition: color 0.2s;
}
.breadcrumb-link:hover {
    color: var(--green-dark);
}
.breadcrumb-sep {
    color: var(--charcoal);
    opacity: 0.4;
}
.breadcrumb-current {
    color: var(--charcoal);
    opacity: 0.7;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 400px;
}

.thread-section {
    flex: 1;
    padding: 40px 0 80px;
}

.thread-header {
    margin-bottom: 32px;
}
.badge-cat {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 5px;
    background: var(--green-pale);
    color: var(--green-dark);
    font-size: 0.75rem;
    font-weight: 700;
    margin-bottom: 12px;
}
.thread-title {
    font-size: clamp(1.5rem, 3vw, 2.2rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.2;
    letter-spacing: -0.02em;
    margin: 0;
}

.post {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 10px;
    padding: 24px;
    margin-bottom: 16px;
}
.post--original {
    background: var(--green-pale);
    border-color: rgba(8, 106, 53, 0.15);
}

.post-author {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
}
.post-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: var(--green-mid);
    color: var(--white);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.72rem;
    font-weight: 700;
    flex-shrink: 0;
}
.post--original .post-avatar {
    background: var(--green-dark);
}
.post-author-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}
.post-author-name {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--charcoal);
}
.post-date {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.55;
}

.post-body {
    padding-left: 52px;
}
.post-paragraph {
    font-size: 0.9rem;
    color: var(--charcoal);
    line-height: 1.7;
    margin: 0 0 14px;
    opacity: 0.88;
}
.post-paragraph:last-child {
    margin-bottom: 0;
}

.replies-section {
    margin-top: 40px;
}
.replies-title {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 20px;
    opacity: 0.7;
}

.reply-form-section {
    margin-top: 48px;
    padding-top: 40px;
    border-top: 1px solid rgba(53, 53, 53, 0.1);
}
.reply-form-title {
    font-size: 1.2rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 20px;
}
.reply-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.reply-textarea {
    width: 100%;
    padding: 16px 18px;
    font-size: 0.92rem;
    font-family: inherit;
    color: var(--charcoal);
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    outline: none;
    resize: vertical;
    transition:
        border-color 0.2s,
        box-shadow 0.2s;
    box-sizing: border-box;
}
.reply-textarea::placeholder {
    color: rgba(53, 53, 53, 0.4);
}
.reply-textarea:focus {
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.12);
}
.reply-form-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
}
.reply-login-hint {
    font-size: 0.85rem;
    color: var(--green-light);
    text-decoration: none;
    transition: color 0.2s;
}
.reply-login-hint:hover {
    color: var(--green-dark);
}
.btn-reply {
    background: var(--green-dark);
    color: var(--white);
    border: none;
    padding: 12px 28px;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-reply:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}

@media (max-width: 700px) {
    .post-body {
        padding-left: 0;
    }
    .breadcrumb-current {
        max-width: 200px;
    }
}
</style>
