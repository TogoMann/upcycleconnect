<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useClientStore } from '@/stores/client'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const clientStore = useClientStore()

const viewMode = ref<'month' | 'week'>('month')
const currentDate = ref(new Date())
const selectedDate = ref<string | null>(null)
const showModal = ref(false)

const HOUR_HEIGHT = 60
const DAY_START = 7
const DAY_END = 21

const DAY_NAMES = computed(() => [
    t('client.planning.dayMon'), t('client.planning.dayTue'), t('client.planning.dayWed'),
    t('client.planning.dayThu'), t('client.planning.dayFri'), t('client.planning.daySat'), t('client.planning.daySun'),
])

function toISO(d: Date): string {
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
}

function isToday(d: Date): boolean {
    return toISO(d) === toISO(new Date())
}

function navigate(dir: number) {
    const d = new Date(currentDate.value)
    if (viewMode.value === 'month') {
        d.setMonth(d.getMonth() + dir)
    } else {
        d.setDate(d.getDate() + dir * 7)
    }
    currentDate.value = d
}

function goToday() {
    currentDate.value = new Date()
}

function switchView(v: 'month' | 'week') {
    viewMode.value = v
}

const monthCells = computed(() => {
    const year = currentDate.value.getFullYear()
    const month = currentDate.value.getMonth()
    const first = new Date(year, month, 1)
    const daysInMonth = new Date(year, month + 1, 0).getDate()
    const startDow = (first.getDay() + 6) % 7

    const cells: { date: Date; current: boolean }[] = []

    for (let i = startDow - 1; i >= 0; i--) {
        cells.push({ date: new Date(year, month, -i), current: false })
    }
    for (let d = 1; d <= daysInMonth; d++) {
        cells.push({ date: new Date(year, month, d), current: true })
    }
    while (cells.length < 42) {
        const next = cells.length - daysInMonth - startDow + 1
        cells.push({ date: new Date(year, month + 1, next), current: false })
    }

    return cells
})

const weekDays = computed(() => {
    const d = new Date(currentDate.value)
    const dow = (d.getDay() + 6) % 7
    d.setDate(d.getDate() - dow)
    return Array.from({ length: 7 }, (_, i) => {
        const day = new Date(d)
        day.setDate(d.getDate() + i)
        return day
    })
})

const timeSlots = computed(() =>
    Array.from({ length: DAY_END - DAY_START }, (_, i) =>
        `${String(DAY_START + i).padStart(2, '0')}:00`
    )
)

function eventsForDate(date: Date) {
    const iso = toISO(date)
    return (clientStore.planning || []).filter((e: any) => e.date === iso)
}

const TYPE = computed<Record<string, { bg: string; border: string; text: string; label: string }>>(() => ({
    depot:    { bg: '#ebf5fb', border: '#3498db', text: '#1a6fa0', label: t('client.planning.typeDeposit') },
    workshop: { bg: '#f4ecf7', border: '#9b59b6', text: '#6c3483', label: t('client.planning.typeWorkshop') },
    event:    { bg: '#fef9e7', border: '#f39c12', text: '#9a6200', label: t('client.planning.typeEvent') },
    personal: { bg: '#e8f8f5', border: '#1abc9c', text: '#0e6655', label: t('client.planning.typePersonal') },
}))

function ts(type: string) {
    return TYPE.value[type] || { bg: '#f0f0f0', border: '#aaa', text: '#555', label: t('client.planning.typeOther') }
}

function parseHour(t: string): number {
    const [h, m] = (t || '00:00').substring(0, 5).split(':').map(Number)
    return h + m / 60
}

function eventStyle(ev: any): Record<string, string> {
    const start = Math.max(parseHour(ev.start_time), DAY_START)
    const end = Math.min(ev.end_time ? parseHour(ev.end_time) : (parseHour(ev.start_time) + 1), DAY_END)
    const top = (start - DAY_START) * HOUR_HEIGHT
    const height = Math.max((end - start) * HOUR_HEIGHT, 24)
    const s = ts(ev.type)
    return {
        top: `${top}px`,
        height: `${height}px`,
        backgroundColor: s.bg,
        borderLeft: `3px solid ${s.border}`,
        color: s.text,
    }
}

const dateLocale = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

const periodLabel = computed(() => {
    if (viewMode.value === 'month') {
        return currentDate.value.toLocaleDateString(dateLocale.value, { month: 'long', year: 'numeric' })
    }
    const days = weekDays.value
    const s = days[0].toLocaleDateString(dateLocale.value, { day: 'numeric', month: 'short' })
    const e = days[6].toLocaleDateString(dateLocale.value, { day: 'numeric', month: 'short', year: 'numeric' })
    return `${s} – ${e}`
})

function selectDay(date: Date) {
    const iso = toISO(date)
    selectedDate.value = selectedDate.value === iso ? null : iso
}

const selectedEvents = computed(() =>
    selectedDate.value
        ? (clientStore.planning || []).filter((e: any) => e.date === selectedDate.value)
        : []
)

const selectedDateLabel = computed(() => {
    if (!selectedDate.value) return ''
    return new Date(selectedDate.value + 'T12:00:00').toLocaleDateString(dateLocale.value, {
        weekday: 'long', day: 'numeric', month: 'long', year: 'numeric',
    })
})

const form = reactive({
    title: '', description: '', date: '', start_time: '', end_time: '', all_day: false,
})

function openModal(preDate?: string) {
    form.title = ''
    form.description = ''
    form.date = preDate || toISO(new Date())
    form.start_time = ''
    form.end_time = ''
    form.all_day = false
    showModal.value = true
}

async function submitEvent() {
    try {
        const payload = { ...form }
        if (payload.all_day) {
            payload.start_time = '00:00'
            payload.end_time = '23:59'
        }
        if (!payload.start_time || !payload.end_time) {
            throw new Error(t('client.planning.errorTimeRequired'))
        }
        await clientStore.createPersonalEvent(payload)
        showModal.value = false
    } catch (e: any) {
        alert(e.message)
    }
}

async function deleteItem(item: any) {
    if (item.type === 'personal' && confirm(t('client.planning.confirmDeletePersonal'))) {
        await clientStore.deletePersonalEvent(item.id)
    } else if (item.type === 'depot' && confirm(t('client.planning.confirmDeleteDeposit'))) {
        await clientStore.deleteEntry(item.id)
        await clientStore.fetchPlanning()
    }
}

function ft(t: string) {
    return (t || '').substring(0, 5)
}

const todayISO = toISO(new Date())
const gridHeight = (DAY_END - DAY_START) * HOUR_HEIGHT

onMounted(() => clientStore.fetchPlanning())
</script>

<template>
    <div class="planning">
        <div class="cal-header">
            <div class="nav-group">
                <button class="nav-btn" @click="navigate(-1)">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="15 18 9 12 15 6"/></svg>
                </button>
                <span class="period-label">{{ periodLabel }}</span>
                <button class="nav-btn" @click="navigate(1)">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
                </button>
            </div>

            <div class="control-group">
                <button class="today-btn" @click="goToday">{{ t('client.planning.today') }}</button>
                <div class="view-toggle">
                    <button :class="{ active: viewMode === 'month' }" @click="switchView('month')">{{ t('client.planning.month') }}</button>
                    <button :class="{ active: viewMode === 'week' }" @click="switchView('week')">{{ t('client.planning.week') }}</button>
                </div>
                <router-link to="/particulier/conteneurs/deposer" class="btn-secondary">{{ t('client.planning.newDeposit') }}</router-link>
                <button class="btn-primary" @click="openModal()">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                    {{ t('client.planning.add') }}
                </button>
            </div>
        </div>

        <div v-if="clientStore.isLoading" class="cal-loading">{{ t('client.planning.loading') }}</div>

        <template v-else-if="viewMode === 'month'">
            <div class="month-view">
                <div class="month-daynames">
                    <div v-for="d in DAY_NAMES" :key="d" class="dayname">{{ d }}</div>
                </div>
                <div class="month-grid">
                    <div
                        v-for="cell in monthCells"
                        :key="toISO(cell.date)"
                        class="month-cell"
                        :class="{
                            'month-cell--other': !cell.current,
                            'month-cell--today': isToday(cell.date),
                            'month-cell--selected': selectedDate === toISO(cell.date),
                        }"
                        @click="selectDay(cell.date)"
                    >
                        <div class="cell-num">{{ cell.date.getDate() }}</div>
                        <div class="cell-events">
                            <div
                                v-for="ev in eventsForDate(cell.date).slice(0, 3)"
                                :key="ev.id + ev.type"
                                class="event-chip"
                                :style="{ backgroundColor: ts(ev.type).bg, color: ts(ev.type).text, borderLeft: `2px solid ${ts(ev.type).border}` }"
                            >
                                <span class="chip-time">{{ ft(ev.start_time) }}</span>
                                <span class="chip-title">{{ ev.title }}</span>
                            </div>
                            <div v-if="eventsForDate(cell.date).length > 3" class="event-more">
                                {{ t('client.planning.more', { count: eventsForDate(cell.date).length - 3 }) }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="selectedDate" class="day-panel">
                <div class="day-panel-head">
                    <span class="day-panel-label">{{ selectedDateLabel }}</span>
                    <div class="day-panel-actions">
                        <button class="btn-add-small" @click="openModal(selectedDate ?? undefined)">{{ t('client.planning.addEvent') }}</button>
                        <button class="close-btn" @click="selectedDate = null">×</button>
                    </div>
                </div>
                <div v-if="selectedEvents.length === 0" class="day-panel-empty">
                    {{ t('client.planning.noEventsToday') }}
                    <span class="link" @click="openModal(selectedDate ?? undefined)">{{ t('client.planning.addLink') }}</span>
                </div>
                <div v-else class="day-panel-list">
                    <div v-for="ev in selectedEvents" :key="ev.id + ev.type" class="day-ev-row" :style="{ borderLeft: `3px solid ${ts(ev.type).border}` }">
                        <div class="day-ev-time">
                            <span v-if="ev.end_time">{{ ft(ev.start_time) }} – {{ ft(ev.end_time) }}</span>
                            <span v-else>{{ ft(ev.start_time) }}</span>
                        </div>
                        <div class="day-ev-info">
                            <span class="day-ev-title">{{ ev.title }}</span>
                            <span class="day-ev-type" :style="{ color: ts(ev.type).text }">{{ ts(ev.type).label }}</span>
                        </div>
                        <span v-if="ev.location" class="day-ev-loc">{{ ev.location }}</span>
                        <button v-if="ev.type === 'personal' || ev.type === 'depot'" class="btn-del" @click="deleteItem(ev)">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/></svg>
                        </button>
                    </div>
                </div>
            </div>
        </template>

        <template v-else>
            <div class="week-view">
                <div class="week-head">
                    <div class="wh-gutter"></div>
                    <div
                        v-for="(day, i) in weekDays"
                        :key="i"
                        class="wh-day"
                        :class="{ 'wh-day--today': isToday(day) }"
                    >
                        <span class="wh-name">{{ DAY_NAMES[i] }}</span>
                        <span class="wh-num" :class="{ 'today-circle': isToday(day) }">{{ day.getDate() }}</span>
                    </div>
                </div>

                <div class="week-body">
                    <div class="week-times">
                        <div v-for="slot in timeSlots" :key="slot" class="wt-label" :style="{ height: HOUR_HEIGHT + 'px' }">
                            {{ slot }}
                        </div>
                    </div>

                    <div class="week-grid-area" :style="{ height: gridHeight + 'px' }">
                        <div class="hour-lines">
                            <div v-for="slot in timeSlots" :key="slot" class="hour-line" :style="{ height: HOUR_HEIGHT + 'px' }"></div>
                        </div>

                        <div class="day-cols">
                            <div
                                v-for="(day, i) in weekDays"
                                :key="i"
                                class="day-col"
                                :class="{ 'day-col--today': isToday(day) }"
                                @click="openModal(toISO(day))"
                            >
                                <div
                                    v-for="ev in eventsForDate(day)"
                                    :key="ev.id + ev.type"
                                    class="week-ev"
                                    :style="eventStyle(ev)"
                                    @click.stop
                                >
                                    <span class="wev-time">{{ ft(ev.start_time) }}</span>
                                    <span class="wev-title">{{ ev.title }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </template>

        <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
            <div class="modal">
                <div class="modal-head">
                    <h2 class="modal-title">{{ t('client.planning.newEventTitle') }}</h2>
                    <button class="modal-close" @click="showModal = false">×</button>
                </div>
                <form class="modal-body" @submit.prevent="submitEvent">
                    <div class="fg">
                        <label>{{ t('client.planning.titleLabel') }}</label>
                        <input v-model="form.title" type="text" required :placeholder="t('client.planning.titlePlaceholder')" />
                    </div>
                    <div class="fg">
                        <label>{{ t('client.planning.descriptionLabel') }}</label>
                        <textarea v-model="form.description" rows="2"></textarea>
                    </div>
                    <label class="cb-label">
                        <input v-model="form.all_day" type="checkbox" />
                        {{ t('client.planning.allDay') }}
                    </label>
                    <div class="fg">
                        <label>{{ t('client.planning.dateLabel') }}</label>
                        <input v-model="form.date" type="date" required :min="todayISO" />
                    </div>
                    <div v-if="!form.all_day" class="fg-row">
                        <div class="fg">
                            <label>{{ t('client.planning.startLabel') }}</label>
                            <input v-model="form.start_time" type="time" :required="!form.all_day" />
                        </div>
                        <div class="fg">
                            <label>{{ t('client.planning.endLabel') }}</label>
                            <input v-model="form.end_time" type="time" :required="!form.all_day" />
                        </div>
                    </div>
                    <div class="modal-acts">
                        <button type="button" class="btn-cancel" @click="showModal = false">{{ t('client.planning.cancel') }}</button>
                        <button type="submit" class="btn-submit">{{ t('client.planning.save') }}</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped>
.planning {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.cal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 20px;
}

.nav-group {
    display: flex;
    align-items: center;
    gap: 8px;
}

.nav-btn {
    background: none;
    border: 1.5px solid rgba(53,53,53,0.15);
    border-radius: 8px;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: var(--charcoal);
    transition: all 0.15s;
}
.nav-btn:hover { border-color: var(--green-mid); color: var(--green-dark); }

.period-label {
    font-size: 1.1rem;
    font-weight: 800;
    color: var(--charcoal);
    text-transform: capitalize;
    min-width: 180px;
    text-align: center;
}

.control-group {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
}

.today-btn {
    padding: 7px 14px;
    border: 1.5px solid rgba(53,53,53,0.15);
    border-radius: 8px;
    background: none;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    color: var(--charcoal);
}
.today-btn:hover { border-color: var(--green-mid); color: var(--green-dark); }

.view-toggle {
    display: flex;
    border: 1.5px solid rgba(53,53,53,0.15);
    border-radius: 8px;
    overflow: hidden;
}
.view-toggle button {
    padding: 7px 14px;
    background: none;
    border: none;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    color: var(--charcoal);
    transition: all 0.15s;
}
.view-toggle button.active {
    background: var(--green-dark);
    color: white;
}

.btn-primary, .btn-secondary {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    border-radius: 8px;
    font-size: 0.82rem;
    font-weight: 600;
    text-decoration: none;
    cursor: pointer;
    border: none;
    font-family: inherit;
    transition: all 0.15s;
}
.btn-primary { background: var(--green-dark); color: white; }
.btn-primary:hover { background: var(--green-mid); }
.btn-secondary { background: var(--green-pale); color: var(--green-dark); }
.btn-secondary:hover { background: var(--green-light); color: white; }

.cal-loading {
    text-align: center;
    padding: 60px;
    opacity: 0.5;
}


.month-view {
    border: 1.5px solid rgba(53,53,53,0.1);
    border-radius: 14px;
    overflow: hidden;
    background: var(--white);
}

.month-daynames {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    background: rgba(53,53,53,0.03);
    border-bottom: 1.5px solid rgba(53,53,53,0.08);
}
.dayname {
    text-align: center;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    padding: 10px 0;
    color: var(--charcoal);
    opacity: 0.5;
}

.month-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
}

.month-cell {
    min-height: 100px;
    padding: 8px;
    border-right: 1px solid rgba(53,53,53,0.06);
    border-bottom: 1px solid rgba(53,53,53,0.06);
    cursor: pointer;
    transition: background 0.15s;
}
.month-cell:nth-child(7n) { border-right: none; }
.month-cell:hover { background: rgba(215,236,225,0.25); }
.month-cell--other { opacity: 0.35; }
.month-cell--today .cell-num {
    background: var(--green-dark);
    color: white;
    border-radius: 50%;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
}
.month-cell--selected { background: rgba(215,236,225,0.4); }

.cell-num {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--charcoal);
    margin-bottom: 4px;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.cell-events {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.event-chip {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 2px 5px;
    border-radius: 4px;
    font-size: 0.7rem;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.chip-time {
    font-size: 0.65rem;
    opacity: 0.7;
    flex-shrink: 0;
}
.chip-title {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.event-more {
    font-size: 0.65rem;
    color: var(--charcoal);
    opacity: 0.5;
    font-weight: 600;
    padding: 1px 4px;
}


.day-panel {
    margin-top: 16px;
    background: var(--white);
    border: 1.5px solid rgba(53,53,53,0.1);
    border-radius: 12px;
    overflow: hidden;
}

.day-panel-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 18px;
    border-bottom: 1px solid rgba(53,53,53,0.08);
    background: rgba(53,53,53,0.02);
}

.day-panel-label {
    font-size: 0.9rem;
    font-weight: 700;
    text-transform: capitalize;
}

.day-panel-actions {
    display: flex;
    align-items: center;
    gap: 8px;
}

.btn-add-small {
    background: var(--green-dark);
    color: white;
    border: none;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 0.78rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
}

.close-btn {
    background: none;
    border: none;
    font-size: 1.3rem;
    cursor: pointer;
    color: rgba(53,53,53,0.4);
    padding: 0 4px;
    line-height: 1;
}

.day-panel-empty {
    padding: 20px 18px;
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.5;
    display: flex;
    gap: 8px;
    align-items: center;
}
.link {
    color: var(--green-dark);
    font-weight: 600;
    cursor: pointer;
    opacity: 1;
}

.day-panel-list {
    display: flex;
    flex-direction: column;
}

.day-ev-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 18px;
    border-bottom: 1px solid rgba(53,53,53,0.05);
    padding-left: 15px;
}
.day-ev-row:last-child { border-bottom: none; }

.day-ev-time {
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--charcoal);
    opacity: 0.5;
    min-width: 100px;
    flex-shrink: 0;
}

.day-ev-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
}
.day-ev-title { font-size: 0.9rem; font-weight: 700; }
.day-ev-type { font-size: 0.72rem; font-weight: 600; }
.day-ev-loc { font-size: 0.78rem; opacity: 0.5; }

.btn-del {
    background: none;
    border: none;
    cursor: pointer;
    color: rgba(53,53,53,0.3);
    padding: 4px;
    border-radius: 6px;
    transition: all 0.15s;
    display: flex;
}
.btn-del:hover { color: #e53e3e; background: rgba(229,62,62,0.08); }


.week-view {
    background: var(--white);
    border: 1.5px solid rgba(53,53,53,0.1);
    border-radius: 14px;
    overflow: hidden;
}

.week-head {
    display: flex;
    border-bottom: 1.5px solid rgba(53,53,53,0.08);
    background: rgba(53,53,53,0.02);
}

.wh-gutter {
    width: 60px;
    flex-shrink: 0;
}

.wh-day {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 10px 0;
    gap: 4px;
    border-left: 1px solid rgba(53,53,53,0.06);
}

.wh-name {
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--charcoal);
    opacity: 0.5;
}

.wh-num {
    font-size: 1.1rem;
    font-weight: 800;
    color: var(--charcoal);
}

.today-circle {
    background: var(--green-dark);
    color: white;
    border-radius: 50%;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.wh-day--today .wh-name { color: var(--green-dark); opacity: 1; }

.week-body {
    display: flex;
    overflow-y: auto;
    max-height: 620px;
}

.week-times {
    width: 60px;
    flex-shrink: 0;
    padding-top: 0;
}

.wt-label {
    font-size: 0.72rem;
    color: var(--charcoal);
    opacity: 0.4;
    text-align: right;
    padding-right: 10px;
    transform: translateY(-9px);
    display: flex;
    align-items: flex-start;
    justify-content: flex-end;
    box-sizing: border-box;
}

.week-grid-area {
    flex: 1;
    position: relative;
}

.hour-lines {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    pointer-events: none;
}

.hour-line {
    border-top: 1px solid rgba(53,53,53,0.06);
    box-sizing: border-box;
}

.day-cols {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: grid;
    grid-template-columns: repeat(7, 1fr);
}

.day-col {
    position: relative;
    border-left: 1px solid rgba(53,53,53,0.06);
    cursor: pointer;
}
.day-col--today { background: rgba(215,236,225,0.12); }

.week-ev {
    position: absolute;
    left: 3px;
    right: 3px;
    border-radius: 5px;
    padding: 3px 6px;
    font-size: 0.72rem;
    overflow: hidden;
    cursor: default;
    z-index: 1;
    display: flex;
    flex-direction: column;
    gap: 1px;
}

.wev-time {
    font-weight: 700;
    opacity: 0.7;
    flex-shrink: 0;
}
.wev-title {
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}


.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.45);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
}

.modal {
    background: var(--white);
    width: 100%;
    max-width: 480px;
    border-radius: 16px;
    box-shadow: 0 20px 40px rgba(0,0,0,0.15);
    overflow: hidden;
}

.modal-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 22px;
    border-bottom: 1px solid rgba(53,53,53,0.1);
}
.modal-title { margin: 0; font-size: 1.15rem; font-weight: 800; }
.modal-close {
    background: none;
    border: none;
    font-size: 1.4rem;
    cursor: pointer;
    color: rgba(53,53,53,0.4);
    line-height: 1;
    padding: 0;
}

.modal-body { padding: 20px 22px; }

.fg {
    margin-bottom: 14px;
}
.fg label {
    display: block;
    font-size: 0.82rem;
    font-weight: 700;
    margin-bottom: 5px;
}
.fg input, .fg textarea {
    width: 100%;
    padding: 9px 11px;
    border: 1.5px solid rgba(53,53,53,0.12);
    border-radius: 8px;
    font-family: inherit;
    font-size: 0.9rem;
    box-sizing: border-box;
}
.fg textarea { resize: none; }

.cb-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    margin-bottom: 14px;
}

.fg-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
}

.modal-acts {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 20px;
}
.btn-cancel {
    background: none;
    border: none;
    font-weight: 600;
    font-size: 0.875rem;
    cursor: pointer;
    color: var(--charcoal);
    opacity: 0.6;
    font-family: inherit;
}
.btn-submit {
    background: var(--green-dark);
    color: white;
    border: none;
    padding: 9px 22px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
}
</style>
