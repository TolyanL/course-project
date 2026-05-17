<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useScheduleStore } from '@/stores/schedule'
import { groupsApi } from '@/api/client'
import FilterBar from '@/components/FilterBar.vue'
import ScheduleTable from '@/components/ScheduleTable.vue'
import ProgressSpinner from 'primevue/progressspinner'
import Message from 'primevue/message'
import Button from 'primevue/button'
import type { Group } from '@/types'

const scheduleStore = useScheduleStore()
const groups = ref<Group[]>([])
const currentWeekStart = ref<string | null>(null)
const filterStartDate = ref<string>('')
const filterEndDate = ref<string>('')

const weekStartDate = computed((): string => {
  if (currentWeekStart.value) return currentWeekStart.value
  const today = new Date()
  return today.toISOString().split('T')[0] as string
})

function getMonday(date: Date): Date {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  return new Date(d.setDate(diff))
}

function toDateString(date: Date): string {
  return date.toISOString().split('T')[0] as string
}

async function loadWeekSchedule(dateStr: string) {
  const date = new Date(dateStr)
  const monday = getMonday(date)
  const startDate = toDateString(monday)
  const endDate = toDateString(new Date(monday.getTime() + 6 * 24 * 60 * 60 * 1000))
  currentWeekStart.value = startDate
  filterStartDate.value = startDate
  filterEndDate.value = endDate

  await scheduleStore.fetchSchedule({
    start_date: startDate,
    end_date: endDate,
    group_id: undefined,
  })
}

onMounted(async () => {
  const { data } = await groupsApi.getAll()
  groups.value = data
  await loadWeekSchedule(new Date().toISOString().split('T')[0] as string)
})

async function handleFilterChange(filters: {
  start_date: string
  end_date: string
  group_id: number | null
}) {
  currentWeekStart.value = filters.start_date
  filterStartDate.value = filters.start_date
  filterEndDate.value = filters.end_date
  await scheduleStore.fetchSchedule({
    start_date: filters.start_date,
    end_date: filters.end_date,
    group_id: filters.group_id ?? undefined,
  })
}

function prevWeek() {
  const current = new Date(weekStartDate.value)
  current.setDate(current.getDate() - 7)
  loadWeekSchedule(toDateString(current))
}

function nextWeek() {
  const current = new Date(weekStartDate.value)
  current.setDate(current.getDate() + 7)
  loadWeekSchedule(toDateString(current))
}

function currentWeek() {
  const today = new Date()
  const monday = getMonday(today)
  loadWeekSchedule(toDateString(monday))
}

function formatWeekRange(startDate: string): string {
  const start = new Date(startDate)
  const end = new Date(start)
  end.setDate(end.getDate() + 6)
  return `${start.toLocaleDateString('ru-RU', { day: 'numeric', month: 'numeric' })} - ${end.toLocaleDateString('ru-RU', { day: 'numeric', month: 'numeric', year: 'numeric' })}`
}
</script>

<template>
  <div class="schedule-view container py-3" style="max-width: 1400px">
    <h1 class="h4 mb-3">Расписание</h1>

    <div class="week-nav d-flex align-items-center justify-content-between mb-3 gap-2">
      <Button icon="pi pi-chevron-left" text @click="prevWeek" />
      <span class="week-range">{{ formatWeekRange(weekStartDate) }}</span>
      <Button icon="pi pi-chevron-right" text @click="nextWeek" />
      <Button label="Текущая" severity="secondary" size="small" @click="currentWeek" />
    </div>

    <FilterBar :groups="groups" @filter-change="handleFilterChange" />

    <div v-if="scheduleStore.loading" class="text-center py-4">
      <ProgressSpinner />
    </div>
    <Message v-else-if="scheduleStore.error" severity="error" :closable="false">{{
      scheduleStore.error
    }}</Message>
    <ScheduleTable
      v-else
      :entries="scheduleStore.entries"
      :startDate="filterStartDate"
      :endDate="filterEndDate"
    />
  </div>
</template>

<style scoped>
.week-nav {
  background: white;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.week-range {
  font-weight: 600;
  font-size: 1rem;
  white-space: nowrap;
}

@media (max-width: 576px) {
  .week-nav {
    flex-wrap: wrap;
    justify-content: center !important;
  }

  .week-range {
    order: -1;
    width: 100%;
    text-align: center;
    margin-bottom: 0.5rem;
  }
}
</style>
