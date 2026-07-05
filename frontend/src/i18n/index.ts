import { createI18n } from 'vue-i18n'
import fr from './locales/fr.json'
import en from './locales/en.json'

function detectLocale(): string {
    const stored = localStorage.getItem('locale')
    if (stored === 'fr' || stored === 'en') return stored
    const browser = navigator.language.slice(0, 2)
    return browser === 'en' ? 'en' : 'fr'
}

const i18n = createI18n({
    legacy: false,
    locale: detectLocale(),
    fallbackLocale: 'fr',
    messages: { fr, en },
})

export default i18n
