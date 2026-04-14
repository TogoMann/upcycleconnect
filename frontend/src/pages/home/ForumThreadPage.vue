<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const threadId = Number(route.params.id)

interface Thread {
    id: number
    created_by: number | null
    title: string
    content: string
    upvotes: number
    downvotes: number
    created_at: string | null
    last_post_at: string | null
}

interface Post {
    id: number
    thread_id: number | null
    created_by: number | null
    content: string
    upvotes: number
    downvotes: number
    created_at: string | null
}

const thread = ref<Thread | null>(null)
const posts = ref<Post[]>([])
const replyContent = ref('')
const sending = ref(false)
const loading = ref(true)

onMounted(async () => {
    try {
        const [threadRes, postsRes] = await Promise.all([
            fetch(`http://localhost:8081/thread/${threadId}`),
            fetch(`http://localhost:8081/thread/${threadId}/posts`),
        ])
        if (threadRes.ok) thread.value = await threadRes.json()
        if (postsRes.ok) posts.value = await postsRes.json()
    } catch {}
    loading.value = false
})

async function handleReply() {
    if (!replyContent.value.trim() || !authStore.isAuthenticated) return
    sending.value = true
    try {
        const res = await fetch('http://localhost:8081/post/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({
                thread_id: threadId,
                created_by: authStore.user?.id,
                content: replyContent.value,
            }),
        })
        if (res.ok) {
            const newPost = await res.json()
            posts.value.push(newPost)
            replyContent.value = ''
        }
    } catch {}
    sending.value = false
}

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function authorLabel(id: number | null): string {
    return id ? `Utilisateur #${id}` : 'Anonyme'
}

function avatarChars(id: number | null): string {
    return id ? `U${id}` : '??'
}
</script>

<template>
    <div class="page-content">
        <section class="breadcrumb-bar">
            <div class="container">
                <router-link to="/forum" class="breadcrumb-link">Forum</router-link>
                <span class="breadcrumb-sep">›</span>
                <span class="breadcrumb-current">{{ thread?.title ?? '…' }}</span>
            </div>
        </section>

        <section class="thread-section">
            <div class="container">
                <div v-if="loading" class="loading">Chargement…</div>

                <template v-else-if="thread">
                    <div class="thread-header">
                        <h1 class="thread-title">{{ thread.title }}</h1>
                    </div>

                    <div class="post post--original">
                        <div class="post-author">
                            <div class="post-avatar">{{ avatarChars(thread.created_by) }}</div>
                            <div class="post-author-info">
                                <span class="post-author-name">{{ authorLabel(thread.created_by) }}</span>
                                <span class="post-date">{{ fmtDate(thread.created_at) }}</span>
                            </div>
                        </div>
                        <div class="post-body">
                            <p v-for="(para, i) in thread.content.split('\n\n')" :key="i" class="post-paragraph">
                                {{ para }}
                            </p>
                        </div>
                        <div class="post-votes">
                            <span class="vote vote--up">▲ {{ thread.upvotes }}</span>
                            <span class="vote vote--down">▼ {{ thread.downvotes }}</span>
                        </div>
                    </div>

                    <div class="replies-section">
                        <h2 class="replies-title">{{ posts.length }} réponse{{ posts.length !== 1 ? 's' : '' }}</h2>

                        <div v-for="p in posts" :key="p.id" class="post">
                            <div class="post-author">
                                <div class="post-avatar">{{ avatarChars(p.created_by) }}</div>
                                <div class="post-author-info">
                                    <span class="post-author-name">{{ authorLabel(p.created_by) }}</span>
                                    <span class="post-date">{{ fmtDate(p.created_at) }}</span>
                                </div>
                            </div>
                            <div class="post-body">
                                <p v-for="(para, i) in p.content.split('\n\n')" :key="i" class="post-paragraph">
                                    {{ para }}
                                </p>
                            </div>
                            <div class="post-votes">
                                <span class="vote vote--up">▲ {{ p.upvotes }}</span>
                                <span class="vote vote--down">▼ {{ p.downvotes }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="reply-form-section">
                        <h2 class="reply-form-title">Répondre à ce sujet</h2>
                        <form class="reply-form" @submit.prevent="handleReply">
                            <textarea
                                v-model="replyContent"
                                placeholder="Partagez votre expérience ou vos conseils..."
                                class="reply-textarea"
                                rows="6"
                                :disabled="!authStore.isAuthenticated"
                                required
                            />
                            <div class="reply-form-footer">
                                <router-link v-if="!authStore.isAuthenticated" to="/auth/login" class="reply-login-hint">
                                    Connectez-vous pour répondre
                                </router-link>
                                <span v-else class="reply-login-hint reply-login-hint--ok">
                                    Connecté en tant que {{ authStore.user?.username }}
                                </span>
                                <button
                                    type="submit"
                                    class="btn-reply"
                                    :disabled="!authStore.isAuthenticated || sending"
                                >
                                    {{ sending ? 'Envoi…' : 'Publier la réponse' }}
                                </button>
                            </div>
                        </form>
                    </div>
                </template>

                <div v-else class="empty-state">Discussion introuvable.</div>
            </div>
        </section>
    </div>
</template>

<style scoped>
.page-content { flex: 1; display: flex; flex-direction: column; }
.container { max-width: 1060px; margin: 0 auto; padding: 0 32px; }
.breadcrumb-bar { padding: 20px 0; border-bottom: 1px solid rgba(53,53,53,0.08); }
.breadcrumb-bar .container { display: flex; align-items: center; gap: 8px; font-size: 0.85rem; }
.breadcrumb-link { color: var(--green-mid); text-decoration: none; transition: color 0.2s; }
.breadcrumb-link:hover { color: var(--green-dark); }
.breadcrumb-sep { color: var(--charcoal); opacity: 0.4; }
.breadcrumb-current { color: var(--charcoal); opacity: 0.7; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 400px; }
.thread-section { flex: 1; padding: 40px 0 80px; }
.thread-header { margin-bottom: 32px; }
.thread-title { font-size: clamp(1.5rem, 3vw, 2.2rem); font-weight: 800; color: var(--charcoal); line-height: 1.2; letter-spacing: -0.02em; margin: 0; }
.post { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 10px; padding: 24px; margin-bottom: 16px; }
.post--original { background: var(--green-pale); border-color: rgba(8,106,53,0.15); }
.post-author { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.post-avatar { width: 40px; height: 40px; border-radius: 50%; background: var(--green-mid); color: var(--white); display: flex; align-items: center; justify-content: center; font-size: 0.7rem; font-weight: 700; flex-shrink: 0; }
.post--original .post-avatar { background: var(--green-dark); }
.post-author-info { display: flex; flex-direction: column; gap: 2px; }
.post-author-name { font-size: 0.9rem; font-weight: 700; color: var(--charcoal); }
.post-date { font-size: 0.78rem; color: var(--charcoal); opacity: 0.55; }
.post-body { padding-left: 52px; }
.post-paragraph { font-size: 0.9rem; color: var(--charcoal); line-height: 1.7; margin: 0 0 14px; opacity: 0.88; }
.post-paragraph:last-child { margin-bottom: 0; }
.post-votes { padding-left: 52px; margin-top: 12px; display: flex; gap: 12px; }
.vote { font-size: 0.78rem; font-weight: 600; opacity: 0.6; }
.vote--up { color: var(--green-dark); }
.vote--down { color: #dc2626; }
.replies-section { margin-top: 40px; }
.replies-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 20px; opacity: 0.7; }
.reply-form-section { margin-top: 48px; padding-top: 40px; border-top: 1px solid rgba(53,53,53,0.1); }
.reply-form-title { font-size: 1.2rem; font-weight: 700; color: var(--charcoal); margin: 0 0 20px; }
.reply-form { display: flex; flex-direction: column; gap: 16px; }
.reply-textarea { width: 100%; padding: 16px 18px; font-size: 0.92rem; font-family: inherit; color: var(--charcoal); background: var(--white); border: 1.5px solid rgba(53,53,53,0.25); border-radius: 8px; outline: none; resize: vertical; transition: border-color 0.2s, box-shadow 0.2s; box-sizing: border-box; }
.reply-textarea:disabled { background: rgba(53,53,53,0.04); opacity: 0.6; }
.reply-textarea::placeholder { color: rgba(53,53,53,0.4); }
.reply-textarea:focus { border-color: var(--green-mid); box-shadow: 0 0 0 3px rgba(52,137,91,0.12); }
.reply-form-footer { display: flex; align-items: center; justify-content: space-between; gap: 16px; flex-wrap: wrap; }
.reply-login-hint { font-size: 0.85rem; color: var(--green-light); text-decoration: none; transition: color 0.2s; }
.reply-login-hint:hover { color: var(--green-dark); }
.reply-login-hint--ok { color: var(--green-mid); }
.btn-reply { background: var(--green-dark); color: var(--white); border: none; padding: 12px 28px; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; font-family: inherit; transition: background 0.2s, transform 0.15s; }
.btn-reply:hover:not(:disabled) { background: var(--green-mid); transform: translateY(-1px); }
.btn-reply:disabled { opacity: 0.5; cursor: default; }
.loading, .empty-state { opacity: 0.5; font-size: 0.9rem; padding: 40px 0; }
@media (max-width: 700px) { .post-body { padding-left: 0; } .post-votes { padding-left: 0; } .breadcrumb-current { max-width: 200px; } }
</style>
