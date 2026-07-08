<script setup lang="ts">
import { reactive, ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from '@/stores/client'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const clientStore = useClientStore()
const authStore = useAuthStore()

const backRoute = computed(() => {
    return authStore.userRole === 'pro' ? '/pro/annonces' : '/particulier/annonces'
})

const CATEGORIES = ['Mobilier', 'Décoration', 'Vêtements', 'Jouet', 'Electronique', 'Outils'] as const
const ETATS = ['Neuf', 'Bon état', 'Abimé', 'Cassé'] as const
const SIZES = computed(() => [
    { value: 'S', label: t('client.nouvelleAnnonce.sizeSmall') },
    { value: 'M', label: t('client.nouvelleAnnonce.sizeMedium') },
    { value: 'L', label: t('client.nouvelleAnnonce.sizeLarge') },
])
const SIZE_RANK: Record<string, number> = { S: 1, M: 2, L: 3 }

const categoryLabels = computed<Record<string, string>>(() => ({
    'Mobilier': t('client.nouvelleAnnonce.categoryFurniture'),
    'Décoration': t('client.nouvelleAnnonce.categoryDecoration'),
    'Vêtements': t('client.nouvelleAnnonce.categoryClothing'),
    'Jouet': t('client.nouvelleAnnonce.categoryToy'),
    'Electronique': t('client.nouvelleAnnonce.categoryElectronics'),
    'Outils': t('client.nouvelleAnnonce.categoryTools'),
}))
const etatLabels = computed<Record<string, string>>(() => ({
    'Neuf': t('client.nouvelleAnnonce.stateNew'),
    'Bon état': t('client.nouvelleAnnonce.stateGood'),
    'Abimé': t('client.nouvelleAnnonce.stateDamaged'),
    'Cassé': t('client.nouvelleAnnonce.stateBroken'),
}))

const form = reactive({
    name: '',
    description: '',
    price: '',
    category: '' as string,
    handoffMode: 'main_propre' as 'main_propre' | 'casier',
    itemSize: '' as string,
    city_id: '',
    siteId: '',
    lockerId: '',
    physicalState: '' as string,
    address: '',
    image_url: '',
    isDon: false,
    weight: '',
    street: '',
    zipCode: '',
    city: '',
    department: '',
})

const cityQuery = ref('')
const showCitySuggestions = ref(false)

const citySuggestions = computed(() => {
    const q = cityQuery.value.trim().toLowerCase()
    if (!q) return []
    return clientStore.cities.filter((c: any) => 
        c.name.toLowerCase().includes(q) || 
        c.zip_code.includes(q)
    ).slice(0, 8)
})

function selectCity(city: any) {
    form.city_id = String(city.id)
    form.city = city.name
    form.zipCode = city.zip_code
    cityQuery.value = `${city.name} (${city.zip_code})`
    showCitySuggestions.value = false
    
    if (city.zip_code && city.zip_code.length >= 2) {
        const deptNum = city.zip_code.slice(0, 2)
        form.department = getDepartmentName(deptNum)
    }
}

function getDepartmentName(num: string): string {
    const depts: Record<string, string> = {
        '01': 'Ain', '02': 'Aisne', '03': 'Allier', '04': 'Alpes-de-Haute-Provence', '05': 'Hautes-Alpes',
        '06': 'Alpes-Maritimes', '07': 'Ardèche', '08': 'Ardennes', '09': 'Ariège', '10': 'Aube',
        '11': 'Aude', '12': 'Aveyron', '13': 'Bouches-du-Rhône', '14': 'Calvados', '15': 'Cantal',
        '16': 'Charente', '17': 'Charente-Maritime', '18': 'Cher', '19': 'Corrèze', '2A': 'Corse-du-Sud',
        '2B': 'Haute-Corse', '21': 'Côte-d\'Or', '22': 'Côtes-d\'Armor', '23': 'Creuse', '24': 'Dordogne',
        '25': 'Doubs', '26': 'Drôme', '27': 'Eure', '28': 'Eure-et-Loir', '29': 'Finistère',
        '30': 'Gard', '31': 'Haute-Garonne', '32': 'Gers', '33': 'Gironde', '34': 'Hérault',
        '35': 'Ille-et-Vilaine', '36': 'Indre', '37': 'Indre-et-Loire', '38': 'Isère', '39': 'Jura',
        '40': 'Landes', '41': 'Loir-et-Cher', '42': 'Loire', '43': 'Haute-Loire', '44': 'Loire-Atlantique',
        '45': 'Loiret', '46': 'Lot', '47': 'Lot-et-Garonne', '48': 'Lozère', '49': 'Maine-et-Loire',
        '50': 'Manche', '51': 'Marne', '52': 'Haute-Marne', '53': 'Mayenne', '54': 'Meurthe-et-Moselle',
        '55': 'Meuse', '56': 'Morbihan', '57': 'Moselle', '58': 'Nièvre', '59': 'Nord',
        '60': 'Oise', '61': 'Orne', '62': 'Pas-de-Calais', '63': 'Puy-de-Dôme', '64': 'Pyrénées-Atlantiques',
        '65': 'Hautes-Pyrénées', '66': 'Pyrénées-Orientales', '67': 'Bas-Rhin', '68': 'Haut-Rhin', '69': 'Rhône',
        '70': 'Haute-Saône', '71': 'Saône-et-Loire', '72': 'Sarthe', '73': 'Savoie', '74': 'Haute-Savoie',
        '75': 'Paris', '76': 'Seine-Maritime', '77': 'Seine-et-Marne', '78': 'Yvelines', '79': 'Deux-Sèvres',
        '80': 'Somme', '81': 'Tarn', '82': 'Tarn-et-Garonne', '83': 'Var', '84': 'Vaucluse',
        '85': 'Vendée', '86': 'Vienne', '87': 'Haute-Vienne', '88': 'Vosges', '89': 'Yonne',
        '90': 'Territoire de Belfort', '91': 'Essonne', '92': 'Hauts-de-Seine', '93': 'Seine-Saint-Denis',
        '94': 'Val-de-Marne', '95': 'Val-d\'Oise'
    }
    return depts[num] ? `${num} - ${depts[num]}` : num
}

const sitesWithLockers = ref<any[]>([])
const sitesLoading = ref(false)
const sitesError = ref('')

async function loadSites() {
    sitesWithLockers.value = []
    form.siteId = ''
    form.lockerId = ''
    if (form.handoffMode !== 'casier' || !form.city_id) return

    sitesLoading.value = true
    sitesError.value = ''
    try {
        sitesWithLockers.value = await clientStore.fetchSitesWithLockers(Number(form.city_id))
    } catch (e: any) {
        sitesError.value = e.message
    } finally {
        sitesLoading.value = false
    }
}

watch(() => [form.handoffMode, form.city_id], loadSites)
watch(() => form.siteId, () => { form.lockerId = '' })

function lockerIsUsable(locker: any): boolean {
    if (locker.status !== 'Available') return false
    if (!form.itemSize) return true
    return (SIZE_RANK[locker.size] || 0) >= (SIZE_RANK[form.itemSize] || 0)
}

function siteIsUsable(site: any): boolean {
    return site.lockers.some((l: any) => lockerIsUsable(l))
}

function availableCount(site: any): number {
    return site.lockers.filter((l: any) => l.status === 'Available').length
}

function rawId(value: any): string {
    if (value && typeof value === 'object') return String(value.Int64 ?? value.id ?? '')
    return String(value ?? '')
}

const selectedSite = computed(() => {
    return sitesWithLockers.value.find((s: any) => rawId(s.site_id) === form.siteId) || null
})

watch(() => form.itemSize, () => {
    if (selectedSite.value && !siteIsUsable(selectedSite.value)) {
        form.siteId = ''
        form.lockerId = ''
        return
    }
    const currentLocker = selectedSite.value?.lockers.find((l: any) => rawId(l.id) === form.lockerId)
    if (currentLocker && !lockerIsUsable(currentLocker)) {
        form.lockerId = ''
    }
})

function siteLabel(site: any, index: number): string {
    const type = site.type_site || t('client.nouvelleAnnonce.collectionPoint')
    return `${type} ${index + 1} — ${site.address}`
}

const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)

function onFileChange(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) {
        if (file.size > 5 * 1024 * 1024) {
            errors.global = t('client.nouvelleAnnonce.errorImageTooLarge')
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
    siteId: '',
    lockerId: '',
    itemSize: '',
    physicalState: '',
    address: '',
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
    errors.name = form.name.trim() ? '' : t('client.nouvelleAnnonce.errorTitleRequired')
    errors.description = form.description.trim() ? '' : t('client.nouvelleAnnonce.errorDescriptionRequired')
    errors.price = form.isDon || (form.price && Number(form.price) > 0) ? '' : t('client.nouvelleAnnonce.errorPriceRequired')
    errors.category = form.category ? '' : t('client.nouvelleAnnonce.errorCategoryRequired')
    errors.city_id = form.city_id ? '' : t('client.nouvelleAnnonce.errorCityRequired')
    errors.address = form.handoffMode === 'main_propre' && !form.street.trim() ? 'L\'adresse de rue est requise.' : ''
    errors.itemSize = form.handoffMode === 'casier' && !form.itemSize ? t('client.nouvelleAnnonce.errorSizeRequired') : ''
    errors.siteId = form.handoffMode === 'casier' && !form.siteId ? t('client.nouvelleAnnonce.errorSiteRequired') : ''
    errors.lockerId = form.handoffMode === 'casier' && !form.lockerId ? t('client.nouvelleAnnonce.errorLockerRequired') : ''
    errors.physicalState = form.handoffMode === 'casier' && !form.physicalState ? t('client.nouvelleAnnonce.errorStateRequired') : ''
    
    return !errors.name && !errors.description && !errors.price && !errors.category && !errors.city_id
        && !errors.address && !errors.itemSize && !errors.siteId && !errors.lockerId && !errors.physicalState
}

async function handleSubmit() {
    if (!validate()) return
    submitting.value = true
    errors.global = ''
    try {
        if (imageFile.value) {
            form.image_url = await clientStore.uploadImage(imageFile.value)
        }

        const formattedAddress = form.handoffMode === 'main_propre'
            ? `${form.street.trim()}, ${form.zipCode} ${form.city}, ${form.department.trim()}`
            : ''

        await clientStore.createAnnonce({
            name: form.name.trim(),
            description: form.description.trim(),
            price: form.isDon ? 0 : Number(form.price),
            category: form.category,
            city_id: Number(form.city_id),
            image_url: form.image_url,
            handoff_mode: form.handoffMode,
            address: formattedAddress,
            weight: form.weight ? Number(form.weight) : 0,
            locker_id: form.handoffMode === 'casier' ? Number(form.lockerId) : undefined,
            physical_state: form.handoffMode === 'casier' ? form.physicalState : undefined,
            size: form.handoffMode === 'casier' ? form.itemSize : undefined,
        })
        await clientStore.fetchAnnonces()
        router.push(backRoute.value)
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
            <router-link :to="backRoute" class="back-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="19" y1="12" x2="5" y2="12" />
                    <polyline points="12 19 5 12 12 5" />
                </svg>
                {{ t('client.nouvelleAnnonce.backToListings') }}
            </router-link>
            <h1 class="page-title">{{ t('client.nouvelleAnnonce.pageTitle') }}</h1>
        </div>

        <form class="form-card" @submit.prevent="handleSubmit">
            <div class="form-group">
                <label class="form-label">{{ t('client.nouvelleAnnonce.titleLabel') }}</label>
                <input
                    v-model="form.name"
                    type="text"
                    class="form-input"
                    :class="{ 'form-input--error': errors.name }"
                    :placeholder="t('client.nouvelleAnnonce.titlePlaceholder')"
                />
                <span v-if="errors.name" class="form-error">{{ errors.name }}</span>
            </div>

            <div class="form-group">
                <label class="form-label">{{ t('client.nouvelleAnnonce.descriptionLabel') }}</label>
                <textarea
                    v-model="form.description"
                    class="form-input form-textarea"
                    :class="{ 'form-input--error': errors.description }"
                    :placeholder="t('client.nouvelleAnnonce.descriptionPlaceholder')"
                    rows="5"
                ></textarea>
                <span v-if="errors.description" class="form-error">{{ errors.description }}</span>
            </div>

            <div class="form-group">
                <label class="form-label">{{ t('client.nouvelleAnnonce.categoryLabel') }}</label>
                <select
                    v-model="form.category"
                    class="form-input"
                    :class="{ 'form-input--error': errors.category }"
                >
                    <option value="" disabled>{{ t('client.nouvelleAnnonce.selectCategory') }}</option>
                    <option v-for="cat in CATEGORIES" :key="cat" :value="cat">{{ categoryLabels[cat] }}</option>
                </select>
                <span v-if="errors.category" class="form-error">{{ errors.category }}</span>
            </div>

            <div class="form-group">
                <label class="form-label">{{ t('client.nouvelleAnnonce.handoffModeLabel') }}</label>
                <div class="radio-row">
                    <label class="form-radio-label">
                        <input type="radio" value="main_propre" v-model="form.handoffMode" />
                        {{ t('client.nouvelleAnnonce.handoffInPerson') }}
                    </label>
                    <label class="form-radio-label">
                        <input type="radio" value="casier" v-model="form.handoffMode" />
                        {{ t('client.nouvelleAnnonce.handoffLocker') }}
                    </label>
                </div>
            </div>

            <div class="form-group" v-if="form.handoffMode === 'casier'">
                <label class="form-label">{{ t('client.nouvelleAnnonce.itemSizeLabel') }}</label>
                <select v-model="form.itemSize" class="form-input" :class="{ 'form-input--error': errors.itemSize }">
                    <option value="" disabled>{{ t('client.nouvelleAnnonce.selectSize') }}</option>
                    <option v-for="s in SIZES" :key="s.value" :value="s.value">{{ s.label }}</option>
                </select>
                <span v-if="errors.itemSize" class="form-error">{{ errors.itemSize }}</span>
            </div>

            <div class="form-group form-group--row" style="position: relative;">
                <div class="city-field" style="position: relative;">
                    <label class="form-label">{{ t('client.nouvelleAnnonce.cityLabel') }}</label>
                    <div style="position: relative;">
                        <input
                            v-model="cityQuery"
                            type="text"
                            class="form-input"
                            :class="{ 'form-input--error': errors.city_id }"
                            placeholder="Saisissez une ville ou un code postal..."
                            @focus="showCitySuggestions = true"
                            @blur="setTimeout(() => showCitySuggestions = false, 250)"
                        />
                        <ul v-if="showCitySuggestions && citySuggestions.length > 0" class="suggestions-list">
                            <li v-for="c in citySuggestions" :key="c.id" @mousedown="selectCity(c)">
                                {{ c.name }} ({{ c.zip_code }})
                            </li>
                        </ul>
                    </div>
                    <span v-if="errors.city_id" class="form-error">{{ errors.city_id }}</span>
                </div>

                <div v-if="form.handoffMode === 'casier' && form.city_id" class="sites-summary">
                    <span class="sites-summary-title">{{ t('client.nouvelleAnnonce.availableLockersTitle') }}</span>
                    <div v-if="sitesLoading" class="locker-hint">{{ t('client.nouvelleAnnonce.loading') }}</div>
                    <div v-else-if="sitesWithLockers.length === 0" class="locker-hint">{{ t('client.nouvelleAnnonce.noSitesInCity') }}</div>
                    <ul v-else class="sites-summary-list">
                        <li v-for="(site, index) in sitesWithLockers" :key="rawId(site.site_id)">
                            <span class="sites-summary-name">{{ siteLabel(site, index) }}</span>
                            <span class="sites-summary-count">
                                {{ t('client.nouvelleAnnonce.availableCount', { available: availableCount(site), total: site.lockers.length }) }}
                            </span>
                        </li>
                    </ul>
                </div>
            </div>

            <template v-if="form.handoffMode === 'casier'">
                <div class="form-group">
                    <label class="form-label">{{ t('client.nouvelleAnnonce.collectionPointLabel') }}</label>
                    <div v-if="!form.city_id" class="locker-hint">{{ t('client.nouvelleAnnonce.selectCityFirst') }}</div>
                    <div v-else-if="sitesLoading" class="locker-hint">{{ t('client.nouvelleAnnonce.searchingSites') }}</div>
                    <div v-else-if="sitesError" class="locker-hint locker-hint--error">{{ sitesError }}</div>
                    <div v-else-if="sitesWithLockers.length === 0" class="locker-hint">{{ t('client.nouvelleAnnonce.noCollectionPoints') }}</div>
                    <select v-else v-model="form.siteId" class="form-input" :class="{ 'form-input--error': errors.siteId }">
                        <option value="" disabled>{{ t('client.nouvelleAnnonce.selectCollectionPoint') }}</option>
                        <option
                            v-for="(site, index) in sitesWithLockers"
                            :key="rawId(site.site_id)"
                            :value="rawId(site.site_id)"
                            :disabled="!siteIsUsable(site)"
                        >
                            {{ siteLabel(site, index) }}{{ !siteIsUsable(site) ? t('client.nouvelleAnnonce.full') : '' }}
                        </option>
                    </select>
                    <span v-if="errors.siteId" class="form-error">{{ errors.siteId }}</span>
                </div>

                <div class="form-group">
                    <label class="form-label">{{ t('client.nouvelleAnnonce.physicalStateLabel') }}</label>
                    <select v-model="form.physicalState" class="form-input" :class="{ 'form-input--error': errors.physicalState }">
                        <option value="" disabled>{{ t('client.nouvelleAnnonce.select') }}</option>
                        <option v-for="e in ETATS" :key="e" :value="e">{{ etatLabels[e] }}</option>
                    </select>
                    <span v-if="errors.physicalState" class="form-error">{{ errors.physicalState }}</span>
                </div>

                <div class="form-group" v-if="form.siteId">
                    <label class="form-label">{{ t('client.nouvelleAnnonce.lockerLabel') }}</label>
                    <div v-if="!selectedSite" class="locker-hint">{{ t('client.nouvelleAnnonce.selectCollectionPointFirst') }}</div>
                    <select v-else v-model="form.lockerId" class="form-input" :class="{ 'form-input--error': errors.lockerId }">
                        <option value="" disabled>{{ t('client.nouvelleAnnonce.selectLocker') }}</option>
                        <option
                            v-for="l in selectedSite.lockers"
                            :key="rawId(l.id)"
                            :value="rawId(l.id)"
                            :disabled="!lockerIsUsable(l)"
                        >
                            {{ t('client.nouvelleAnnonce.locker', { label: l.label, size: l.size }) }}
                            {{ l.status !== 'Available' ? t('client.nouvelleAnnonce.lockerFull') : (!lockerIsUsable(l) ? t('client.nouvelleAnnonce.lockerTooSmall') : '') }}
                        </option>
                    </select>
                    <span v-if="errors.lockerId" class="form-error">{{ errors.lockerId }}</span>
                </div>
            </template>

            <template v-if="form.handoffMode === 'main_propre'">
                <div class="form-group">
                    <label class="form-label">Adresse / Rue</label>
                    <input
                        v-model="form.street"
                        type="text"
                        class="form-input"
                        :class="{ 'form-input--error': errors.address }"
                        placeholder="Numéro et nom de rue (ex. 12 rue de la Paix)"
                    />
                    <span v-if="errors.address" class="form-error">{{ errors.address }}</span>
                </div>

                <div class="form-group">
                    <label class="form-label">Département / Région</label>
                    <input
                        v-model="form.department"
                        type="text"
                        class="form-input"
                        placeholder="Département ou région"
                    />
                </div>
            </template>

            <div class="form-group">
                <label class="form-label">{{ t('client.nouvelleAnnonce.photoLabel') }}</label>
                <div class="image-upload-zone" :class="{ 'has-image': imagePreview }">
                    <input type="file" accept="image/*" class="file-input" @change="onFileChange" />
                    <div v-if="!imagePreview" class="upload-placeholder">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
                            <circle cx="8.5" cy="8.5" r="1.5" />
                            <polyline points="21 15 16 10 5 21" />
                        </svg>
                        <span>{{ t('client.nouvelleAnnonce.clickToAddPhoto') }}</span>
                    </div>
                    <img v-else :src="imagePreview" class="preview-img" />
                    <button v-if="imagePreview" type="button" class="btn-remove-img" @click="imagePreview = null; imageFile = null">×</button>
                </div>
            </div>

            <div class="form-group form-group--checkbox">
                <label class="form-checkbox-label">
                    <input type="checkbox" v-model="form.isDon" />
                    {{ t('client.nouvelleAnnonce.donateCheckbox') }}
                </label>
            </div>

            <div class="form-group" v-if="!form.isDon">
                <label class="form-label">{{ t('client.nouvelleAnnonce.priceLabel') }}</label>
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

            <div class="form-group">
                <label class="form-label">{{ t('client.nouvelleAnnonce.weightLabel') }}</label>
                <input
                    v-model="form.weight"
                    type="number"
                    step="0.1"
                    min="0"
                    class="form-input"
                    placeholder="ex. 2.5"
                />
            </div>

            <div v-if="errors.global" class="error-banner">{{ errors.global }}</div>

            <div class="form-actions">
                <router-link :to="backRoute" class="btn-cancel">{{ t('client.nouvelleAnnonce.cancel') }}</router-link>
                <button type="submit" class="btn-submit" :disabled="submitting">
                    {{ submitting ? t('client.nouvelleAnnonce.publishing') : t('client.nouvelleAnnonce.publish') }}
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
.form-checkbox-label {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--charcoal);
    cursor: pointer;
}
.radio-row {
    display: flex;
    gap: 20px;
}
.form-radio-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.88rem;
    font-weight: 500;
    color: var(--charcoal);
    cursor: pointer;
}
.form-error {
    font-size: 0.78rem;
    color: #e53e3e;
    font-weight: 500;
}
.locker-hint {
    font-size: 0.85rem;
    color: var(--charcoal);
    opacity: 0.6;
    padding: 10px 0;
}
.locker-hint--error {
    color: #e53e3e;
    opacity: 1;
}
.form-group--row {
    display: flex;
    flex-direction: row;
    align-items: flex-start;
    gap: 20px;
    flex-wrap: wrap;
}
.city-field {
    flex: 1;
    min-width: 220px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.sites-summary {
    flex: 1;
    min-width: 220px;
    background: var(--cream);
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    border-radius: 8px;
    padding: 12px 14px;
}
.sites-summary-title {
    display: block;
    font-size: 0.78rem;
    font-weight: 700;
    color: var(--charcoal);
    opacity: 0.6;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    margin-bottom: 8px;
}
.sites-summary-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.sites-summary-list li {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    font-size: 0.82rem;
}
.sites-summary-name {
    color: var(--charcoal);
    opacity: 0.8;
}
.sites-summary-count {
    color: var(--green-dark);
    font-weight: 700;
    white-space: nowrap;
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

.suggestions-list {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    border-radius: 8px;
    margin-top: 4px;
    padding: 0;
    list-style: none;
    max-height: 200px;
    overflow-y: auto;
    z-index: 10;
    box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}
.suggestions-list li {
    padding: 10px 14px;
    cursor: pointer;
    font-size: 0.88rem;
    transition: background 0.2s, color 0.2s;
    color: var(--charcoal);
    text-align: left;
}
.suggestions-list li:hover {
    background: var(--green-pale);
    color: var(--green-dark);
}
</style>
