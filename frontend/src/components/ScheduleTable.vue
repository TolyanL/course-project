<script setup lang="ts">
import { computed } from 'vue'
import Button from 'primevue/button'
import type { ScheduleEntry, Pair } from '@/types'
import { PAIR_TIMES } from '@/types'

const props = defineProps<{
  entries: ScheduleEntry[]
  startDate?: string
  endDate?: string
  showActions?: boolean
  currentTeacherId?: number
  adminMode?: boolean
}>()

const emit = defineEmits<{
  edit: [pair: Pair]
  delete: [id: number]
}>()

const DAYS = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота', 'Воскресенье']

interface DayData {
  date: string
  dayName: string
  isSunday: boolean
  isSaturday: boolean
  isToday: boolean
  pairs: Pair[]
  allPairs: (Pair | null)[]
  pairTypes: ('pair' | 'before' | 'after' | 'none')[]
}

interface WeekData {
  startDate: string
  endDate: string
  days: DayData[]
}

function getMonday(date: Date): Date {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

function addDays(date: Date, days: number): Date {
  const result = new Date(date)
  result.setDate(result.getDate() + days)
  return result
}

function formatDateShort(date: Date): string {
  return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'numeric' })
}

function formatDateISO(date: Date): string {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const weeks = computed<WeekData[]>(() => {
  const hasEntries = props.entries && props.entries.length > 0

  if (props.startDate && props.endDate) {
    return generateWeekFromRange(props.startDate, props.endDate, hasEntries ? props.entries : [])
  }

  if (hasEntries) {
    const dates = [...new Set(props.entries.map((e) => e.date))].sort()
    if (dates.length === 0) return generateEmptyWeek(props.entries)

    const firstDate = dates[0]!
    const start = new Date(firstDate)

    const monday = getMonday(start)
    const sunday = addDays(monday, 6)
    const todayStr = formatDateISO(new Date())

    const weeksMap: Map<string, WeekData> = new Map()

    let currentMonday = new Date(monday)
    while (currentMonday <= sunday) {
      const weekKey = formatDateISO(currentMonday)
      const weekSunday = addDays(currentMonday, 6)

      const days: DayData[] = []
      for (let i = 0; i < 7; i++) {
        const dayDate = addDays(currentMonday, i)
        const dayDateStr = formatDateISO(dayDate)
        const dayOfWeek = dayDate.getDay()
        const dayIndex = dayOfWeek === 0 ? 6 : dayOfWeek - 1
        const isSunday = dayOfWeek === 0
        const isSaturday = dayOfWeek === 6
        const isToday = dayDateStr === todayStr

        const dayEntries = props.entries.filter((e) => e.date === dayDateStr)
        const pairs: Pair[] = dayEntries.flatMap((e) => e.pairs)

        const firstPairNumber = pairs.length > 0 ? Math.min(...pairs.map((p) => p.pair_number)) : 7
        const lastPairNumber = pairs.length > 0 ? Math.max(...pairs.map((p) => p.pair_number)) : 0
        const maxPairNumber = 7

        const allPairs: (Pair | null)[] = []
        const pairTypes: ('pair' | 'before' | 'after' | 'none')[] = []
        for (let pn = 1; pn <= maxPairNumber; pn++) {
          if (pn < firstPairNumber) {
            allPairs.push(null)
            pairTypes.push('before')
          } else if (pn > lastPairNumber) {
            allPairs.push(null)
            pairTypes.push('after')
          } else {
            const pair = pairs.find((p) => p.pair_number === pn)
            allPairs.push(pair ?? null)
            pairTypes.push(pair ? 'pair' : 'none')
          }
        }

        days.push({
          date: dayDateStr,
          dayName: DAYS[dayIndex] ?? '',
          isSunday,
          isSaturday,
          isToday,
          pairs,
          allPairs,
          pairTypes,
        })
      }

      weeksMap.set(weekKey, {
        startDate: formatDateShort(currentMonday),
        endDate: formatDateShort(weekSunday),
        days,
      })

      currentMonday = addDays(currentMonday, 7)
    }

    return Array.from(weeksMap.values())
  }

  return generateEmptyWeek(props.entries)
})

function getTime(pairNumber: number): string {
  return PAIR_TIMES[pairNumber] || ''
}

function generateEmptyWeek(entries: ScheduleEntry[] = []): WeekData[] {
  const today = new Date()
  const monday = getMonday(today)
  const sunday = addDays(monday, 6)

  return [
    {
      startDate: formatDateShort(monday),
      endDate: formatDateShort(sunday),
      days: generateDays(monday, entries),
    },
  ]
}

function generateWeekFromRange(
  startDate: string,
  endDate: string,
  entries: ScheduleEntry[],
): WeekData[] {
  const monday = new Date(startDate)
  const sunday = new Date(endDate)

  return [
    {
      startDate: formatDateShort(monday),
      endDate: formatDateShort(sunday),
      days: generateDays(monday, entries),
    },
  ]
}

function generateDays(monday: Date, entries: ScheduleEntry[]): DayData[] {
  const days: DayData[] = []
  const todayStr = formatDateISO(new Date())

  for (let i = 0; i < 7; i++) {
    const dayDate = addDays(monday, i)
    const dayDateStr = formatDateISO(dayDate)
    const dayOfWeek = dayDate.getDay()
    const dayIndex = dayOfWeek === 0 ? 6 : dayOfWeek - 1
    const isSunday = dayOfWeek === 0
    const isSaturday = dayOfWeek === 6
    const isToday = dayDateStr === todayStr

    const dayEntries = entries.filter((e) => e.date === dayDateStr)
    const pairs: Pair[] = dayEntries.flatMap((e) => e.pairs)

    const firstPairNumber = pairs.length > 0 ? Math.min(...pairs.map((p) => p.pair_number)) : 7
    const lastPairNumber = pairs.length > 0 ? Math.max(...pairs.map((p) => p.pair_number)) : 0
    const maxPairNumber = 7

    const allPairs: (Pair | null)[] = []
    const pairTypes: ('pair' | 'before' | 'after' | 'none')[] = []

    for (let pn = 1; pn <= maxPairNumber; pn++) {
      if (pn < firstPairNumber) {
        allPairs.push(null)
        pairTypes.push('before')
      } else if (pn > lastPairNumber) {
        allPairs.push(null)
        pairTypes.push('after')
      } else {
        const pair = pairs.find((p) => p.pair_number === pn)
        allPairs.push(pair ?? null)
        pairTypes.push(pair ? 'pair' : 'after')
      }
    }

    days.push({
      date: dayDateStr,
      dayName: DAYS[dayIndex] ?? '',
      isSunday,
      isSaturday,
      isToday,
      pairs,
      allPairs,
      pairTypes,
    })
  }
  return days
}
</script>

<template>
  <div class="schedule-weeks">
    <div v-if="weeks.length === 0" class="text-center py-4 text-muted">Нет пар</div>

    <div v-else class="weeks-list">
      <div v-for="(week, index) in weeks" :key="index" class="week-section">
        <div class="week-header">Неделя {{ week.startDate }} - {{ week.endDate }}</div>

        <div class="week-grid">
          <div
            v-for="day in week.days"
            :key="day.date"
            class="day-card"
            :class="{
              'day-card--empty':
                (day.isSunday || day.isSaturday) && day.allPairs.every((p) => p === null),
              'day-card--today': day.isToday,
              'day-card--saturday': day.isSaturday && !day.isToday,
              'day-card--sunday': day.isSunday && !day.isToday,
            }"
          >
            <div class="day-header">
              <span class="day-name"><i class="pi pi-calendar me-1"></i>{{ day.dayName }}</span>
              <span class="day-date" :class="{ 'text-white': day.isToday }">{{ day.date }}</span>
            </div>

            <div class="day-pairs">
              <template v-if="day.isSunday">
                <div class="sunday-message"><i class="pi pi-moon me-2"></i>Выходной</div>
              </template>
              <template v-else-if="day.pairs.length === 0">
                <div class="no-pairs-message"><i class="pi pi-inbox me-2"></i>Нет пар</div>
              </template>
              <template v-else>
                <div
                  v-for="(pair, index) in day.allPairs"
                  :key="index"
                  :class="pair ? 'pair-item' : 'pair-item pair-item--empty'"
                  v-show="pair || day.pairTypes[index] === 'before'"
                >
                  <template v-if="pair">
                    <div class="pair-time">
                      <span class="pair-number">{{ pair.pair_number }}</span>
                      <span class="pair-time-text">{{ getTime(pair.pair_number) }}</span>
                    </div>
                    <div class="pair-info">
                      <div class="pair-subject">
                        <i class="pi pi-book"></i>
                        {{ pair.subject.name }}
                      </div>
                      <div class="pair-details">
                        <span class="detail-item">
                          <i class="pi pi-map-marker"></i>
                          {{ pair.classroom.number }}
                        </span>
                        <span class="detail-item">
                          <i class="pi pi-user"></i>
                          {{ pair.teacher.name }}
                        </span>
                      </div>
                    </div>
                    <div v-if="showActions && (adminMode || Number(pair.teacher_id) === Number(currentTeacherId))" class="pair-actions">
                      <Button
                        icon="pi pi-pencil"
                        severity="info"
                        text
                        rounded
                        size="small"
                        @click="emit('edit', pair)"
                      />
                      <Button
                        icon="pi pi-trash"
                        severity="danger"
                        text
                        rounded
                        size="small"
                        @click="emit('delete', pair.id)"
                      />
                    </div>
                  </template>
                  <template v-else-if="day.pairTypes[index] === 'before'">
                    <div class="pair-time">
                      <span class="pair-number">{{ index + 1 }}</span>
                      <span class="pair-time-text">{{ getTime(index + 1) }}</span>
                    </div>
                    <div class="pair-info pair-info--empty">
                      <span>-</span>
                    </div>
                  </template>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.schedule-weeks {
  width: 100%;
}

.weeks-list {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.week-section {
  width: 100%;
}

.week-header {
  font-size: 1.125rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #0f6cbd;
}

.week-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0.75rem;
}

.sunday-message,
.no-pairs-message {
  text-align: center;
  color: #6c757d;
  padding: 2rem 0.5rem;
  font-size: 0.875rem;
}

.day-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.day-card--empty {
  opacity: 0.7;
}

.day-card--today {
  border: 2px solid #198754;
  box-shadow: 0 4px 12px rgba(25, 135, 84, 0.2);
}

.day-card--today .day-header {
  background: #198754;
}

.day-card--saturday .day-header {
  background: #6c757d;
}

.day-card--sunday .day-header {
  background: #dc3545;
}

.day-header {
  background: #0f6cbd;
  color: white;
  padding: 0.75rem;
  text-align: center;
}

.day-name {
  display: block;
  font-weight: 600;
  font-size: 0.875rem;
}

.day-date {
  display: block;
  font-size: 0.75rem;
  opacity: 0.9;
}

.day-pairs {
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.pair-item {
  padding: 0.5rem;
  border-radius: 8px;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.pair-item--empty {
  opacity: 0.5;
  background: #e9ecef;
}

.pair-info--empty {
  color: #6c757d;
  font-style: italic;
}

.pair-time {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.pair-number {
  background: #0f6cbd;
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
}

.pair-time-text {
  font-size: 0.75rem;
  color: #666;
}

.pair-info {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.pair-subject {
  font-weight: 500;
  font-size: 0.875rem;
  color: #333;
}

.pair-details {
  font-size: 0.75rem;
  color: #666;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.detail-item i {
  font-size: 0.7rem;
  color: #0f6cbd;
}

.pair-actions {
  display: flex;
  gap: 0.25rem;
  justify-content: flex-end;
  margin-top: 0.25rem;
}

.day-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 0.875rem;
}

@media (max-width: 1200px) {
  .week-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .week-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 0.5rem;
  }

  .day-card {
    min-height: 150px;
  }

  .day-header {
    padding: 0.5rem;
  }

  .day-name {
    font-size: 0.8rem;
  }

  .pair-item {
    padding: 0.375rem;
  }

  .pair-subject {
    font-size: 0.8rem;
  }
}

@media (max-width: 480px) {
  .week-grid {
    grid-template-columns: 1fr;
  }

  .day-card {
    min-height: auto;
  }
}
</style>
