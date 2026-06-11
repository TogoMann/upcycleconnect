<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { API_BASE } from '@/config'

const authStore = useAuthStore()

const HOUR_HEIGHT = 60
const DAY_START = 7
const DAY_END = 21
const DAY_NAMES = ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim']

const viewMode = ref<'month' | 'week'>('month')
const currentDate = ref(new Date())
const selectedDate = ref<string | null>(null)
const filterType = ref('')

const rawItems = ref<any[]>([])
const isLoading = ref(false)

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

const normalizedItems = computed(() =>
    rawItems.value.map((item: any) => ({
        id: item.id,
        title: item.titre || item.title || '',
        type: item.type || '',
        responsable: item.responsable || '',
        date: item.date || '',
        start_time: item.heure_debut || item.start_time || '00:00',
        end_time: item.heure_fin || item.end_time || '00:00',
        participants: item.participants ?? 0,
    }))
)

const filteredItems = computed(() =>
    filterType.value
        ? normalizedItems.value.filter(i => i.type === filterType.value)
        : normalizedItems.value
)

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
    return filteredItems.value.filter(e => e.date === iso)
}

const TYPE: Record<string, { bg: string; border: string; text: string; label: string }> = {
    formation: { bg: '#e8f8f5', border: '#1abc9c', text: '#0e6655', label: 'Formation' },
    atelier:   { bg: '#dbeafe', border: '#3b82f6', text: '#1e40af', label: 'Atelier' },
    depot:     { bg: '#fef3c7', border: '#f59e0b', text: '#92400e', label: 'Dépôt' },
    collecte:  { bg: '#ede9fe', border: '#8b5cf6', text: '#5b21b6', label: 'Collecte' },
    workshop:  { bg: '#dbeafe', border: '#3b82f6', text: '#1e40af', label: 'Atelier' },
    event:     { bg: '#fef9e7', border: '#f39c12', text: '#9a6200', label: 'Événement' },
    personal:  { bg: '#ebf5fb', border: '#3498db', text: '#1a6fa0', label: 'Personnel' },
}

function ts(type: string) {
    return TYPE[type] || { bg: '#f5f5f5', border: '#aaa', text: '#555', label: 'Autre' }
}

function parseHour(t: string): number {
    const [h, m] = (t || '00:00').substring(0, 5).split(':').map(Number)
    return h + m / 60
}

function eventStyle(ev: any): Record<string, string> {
    const start = Math.max(parseHour(ev.start_time), DAY_START)
    const end = Math.min(parseHour(ev.end_time), DAY_END)
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

const periodLabel = computed(() => {
    if (viewMode.value === 'month') {
        return currentDate.value.toLocaleDateString('fr-FR', { month: 'long', year: 'numeric' })
    }
    const days = weekDays.value
    const s = days[0].toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' })
    const e = days[6].toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
    return `${s} – ${e}`
})

function selectDay(date: Date) {
    const iso = toISO(date)
    selectedDate.value = selectedDate.value === iso ? null : iso
}

const selectedEvents = computed(() =>
    selectedDate.value ? filteredItems.value.filter(e => e.date === selectedDate.value) : []
)

const selectedDateLabel = computed(() => {
    if (!selectedDate.value) return ''
    return new Date(selectedDate.value + 'T12:00:00').toLocaleDateString('fr-FR', {
        weekday: 'long', day: 'numeric', month: 'long', year: 'numeric',
    })
})

function ft(t: string) {
    return (t || '').substring(0, 5)
}

const gridHeight = (DAY_END - DAY_START) * HOUR_HEIGHT

const allTypes = computed(() => {
    const set = new Set(rawItems.value.map((i: any) => i.type).filter(Boolean))
    return Array.from(set) as string[]
})

onMounted(async () => {
    isLoading.value = true
    try {
        const res = await fetch(`${API_BASE}/admin/plannings`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) rawItems.value = await res.json()
    } finally {
        isLoading.value = false
    }
})
</script>

<template>
    <div class="plannings">
        <div class="cal-header">
            <div class="left-group">
                <div class="nav-group">
                    <button class="nav-btn" @click="navigate(-1)">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="15 18 9 12 15 6"/></svg>
                    </button>
                    <span class="period-label">{{ periodLabel }}</span>
                    <button class="nav-btn" @click="navigate(1)">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
                    </button>
                </div>
                <div class="filter-chips">
                    <button
                        class="filter-chip"
                        :class="{ active: filterType === '' }"
                        @click="filterType = ''"
                    >Tous</button>
                    <button
                        v-for="t in allTypes"
                        :key="t"
                        class="filter-chip"
                        :class="{ active: filterType === t }"
                        :style="filterType === t ? { backgroundColor: ts(t).border, color: 'white', borderColor: ts(t).border } : {}"
                        @click="filterType = filterType === t ? '' : t"
                    >{{ ts(t).label }}</button>
                </div>
            </div>

            <div class="control-group">
                <button class="today-btn" @click="goToday">Aujourd'hui</button>
                <div class="view-toggle">
                    <button :class="{ active: viewMode === 'month' }" @click="switchView('month')">Mois</button>
                    <button :class="{ active: viewMode === 'week' }" @click="switchView('week')">Semaine</button>
                </div>
            </div>
        </div>

        <div v-if="isLoading" class="cal-loading">Chargement…</div>

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
                                +{{ eventsForDate(cell.date).length - 3 }} autres
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="selectedDate" class="day-panel">
                <div class="day-panel-head">
                    <span class="day-panel-label">{{ selectedDateLabel }}</span>
                    <button class="close-btn" @click="selectedDate = null">×</button>
                </div>
                <div v-if="selectedEvents.length === 0" class="day-panel-empty">Aucun créneau ce jour.</div>
                <div v-else class="day-panel-list">
                    <div v-for="ev in selectedEvents" :key="ev.id" class="day-ev-row" :style="{ borderLeft: `3px solid ${ts(ev.type).border}` }">
                        <div class="day-ev-time">{{ ft(ev.start_time) }} – {{ ft(ev.end_time) }}</div>
                        <div class="day-ev-info">
                            <span class="day-ev-title">{{ ev.title }}</span>
                            <div class="day-ev-meta">
                                <span class="day-ev-type" :style="{ color: ts(ev.type).text }">{{ ts(ev.type).label }}</span>
                                <span v-if="ev.responsable" class="day-ev-resp">· {{ ev.responsable }}</span>
                                <span v-if="ev.participants" class="day-ev-part">· {{ ev.participants }} participant(s)</span>
                            </div>
                        </div>
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
                                @click="selectDay(day)"
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
                                    <span v-if="ev.responsable" class="wev-resp">{{ ev.responsable }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="selectedDate" class="day-panel">
                <div class="day-panel-head">
                    <span class="day-panel-label">{{ selectedDateLabel }}</span>
                    <button class="close-btn" @click="selectedDate = null">×</button>
                </div>
                <div v-if="selectedEvents.length === 0" class="day-panel-empty">Aucun créneau ce jour.</div>
                <div v-else class="day-panel-list">
                    <div v-for="ev in selectedEvents" :key="ev.id" class="day-ev-row" :style="{ borderLeft: `3px solid ${ts(ev.type).border}` }">
                        <div class="day-ev-time">{{ ft(ev.start_time) }} – {{ ft(ev.end_time) }}</div>
                        <div class="day-ev-info">
                            <span class="day-ev-title">{{ ev.title }}</span>
                            <div class="day-ev-meta">
                                <span class="day-ev-type" :style="{ color: ts(ev.type).text }">{{ ts(ev.type).label }}</span>
                                <span v-if="ev.responsable" class="day-ev-resp">· {{ ev.responsable }}</span>
                                <span v-if="ev.participants" class="day-ev-part">· {{ ev.participants }} participant(s)</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>

<style scoped>
.plannings {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.cal-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 20px;
}

.left-group {
    display: flex;
    flex-direction: column;
    gap: 10px;
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
    text-transform: capitalize;
    min-width: 180px;
    text-align: center;
}

.filter-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}

.filter-chip {
    padding: 5px 12px;
    border: 1.5px solid rgba(53,53,53,0.12);
    border-radius: 20px;
    background: none;
    font-size: 0.78rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    font-family: inherit;
    color: var(--charcoal);
}
.filter-chip:hover { border-color: var(--green-mid); color: var(--green-dark); }
.filter-chip.active { background: var(--green-dark); color: white; border-color: var(--green-dark); }

.control-group {
    display: flex;
    align-items: center;
    gap: 10px;
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
    font-family: inherit;
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
    font-family: inherit;
    transition: all 0.15s;
}
.view-toggle button.active { background: var(--green-dark); color: white; }

.cal-loading { text-align: center; padding: 60px; opacity: 0.5; }

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
    opacity: 0.5;
}

.month-grid { display: grid; grid-template-columns: repeat(7, 1fr); }

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
.month-cell--other { opacity: 0.3; }
.month-cell--today .cell-num { background: var(--green-dark); color: white; border-radius: 50%; width: 26px; height: 26px; display: flex; align-items: center; justify-content: center; }
.month-cell--selected { background: rgba(215,236,225,0.4); }

.cell-num {
    font-size: 0.85rem;
    font-weight: 700;
    margin-bottom: 4px;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.cell-events { display: flex; flex-direction: column; gap: 2px; }

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
.chip-time { font-size: 0.65rem; opacity: 0.7; flex-shrink: 0; }
.chip-title { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.event-more { font-size: 0.65rem; opacity: 0.5; font-weight: 600; padding: 1px 4px; }

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
.day-panel-label { font-size: 0.9rem; font-weight: 700; text-transform: capitalize; }
.close-btn { background: none; border: none; font-size: 1.3rem; cursor: pointer; color: rgba(53,53,53,0.4); line-height: 1; }
.day-panel-empty { padding: 20px 18px; font-size: 0.875rem; opacity: 0.5; }
.day-panel-list { display: flex; flex-direction: column; }

.day-ev-row {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px 18px;
    padding-left: 15px;
    border-bottom: 1px solid rgba(53,53,53,0.05);
}
.day-ev-row:last-child { border-bottom: none; }
.day-ev-time { font-size: 0.78rem; font-weight: 600; opacity: 0.5; min-width: 100px; flex-shrink: 0; padding-top: 2px; }
.day-ev-info { flex: 1; display: flex; flex-direction: column; gap: 3px; }
.day-ev-title { font-size: 0.9rem; font-weight: 700; }
.day-ev-meta { display: flex; gap: 6px; align-items: center; font-size: 0.75rem; }
.day-ev-type { font-weight: 600; }
.day-ev-resp { opacity: 0.6; }
.day-ev-part { opacity: 0.6; }

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
.wh-gutter { width: 60px; flex-shrink: 0; }
.wh-day {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 10px 0;
    gap: 4px;
    border-left: 1px solid rgba(53,53,53,0.06);
}
.wh-name { font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; opacity: 0.5; }
.wh-num { font-size: 1.1rem; font-weight: 800; }
.today-circle { background: var(--green-dark); color: white; border-radius: 50%; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; }
.wh-day--today .wh-name { color: var(--green-dark); opacity: 1; }

.week-body { display: flex; overflow-y: auto; max-height: 620px; }
.week-times { width: 60px; flex-shrink: 0; }
.wt-label {
    font-size: 0.72rem;
    opacity: 0.4;
    text-align: right;
    padding-right: 10px;
    transform: translateY(-9px);
    display: flex;
    align-items: flex-start;
    justify-content: flex-end;
    box-sizing: border-box;
}

.week-grid-area { flex: 1; position: relative; }
.hour-lines { position: absolute; top: 0; left: 0; right: 0; pointer-events: none; }
.hour-line { border-top: 1px solid rgba(53,53,53,0.06); box-sizing: border-box; }

.day-cols {
    position: absolute;
    top: 0; left: 0; right: 0; bottom: 0;
    display: grid;
    grid-template-columns: repeat(7, 1fr);
}
.day-col { position: relative; border-left: 1px solid rgba(53,53,53,0.06); cursor: pointer; }
.day-col--today { background: rgba(215,236,225,0.12); }

.week-ev {
    position: absolute;
    left: 3px; right: 3px;
    border-radius: 5px;
    padding: 3px 6px;
    font-size: 0.72rem;
    overflow: hidden;
    z-index: 1;
    display: flex;
    flex-direction: column;
    gap: 1px;
    cursor: default;
}
.wev-time { font-weight: 700; opacity: 0.7; flex-shrink: 0; }
.wev-title { font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.wev-resp { font-size: 0.65rem; opacity: 0.7; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
