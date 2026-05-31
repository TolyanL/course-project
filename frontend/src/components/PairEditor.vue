<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import type { Group, Subject, Classroom, Teacher, Pair, PairFormData, ConflictDetail } from '@/types'
import { PAIR_TIMES } from '@/types'
import { useScheduleStore } from '@/stores/schedule'
import ConflictModal from './ConflictModal.vue'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import Button from 'primevue/button'

const scheduleStore = useScheduleStore()

const props = defineProps<{
  isEdit: boolean
  entry?: Pair
  groups: Group[]
  subjects: Subject[]
  classrooms: Classroom[]
  teachers: Teacher[]
  initialTeacherId?: number
  initialGroupId?: number
  initialDate?: string
}>()

const emit = defineEmits<{
  save: [data: PairFormData]
  cancel: []
}>()

const form = ref({
  group_id: 0,
  date: new Date(),
  subject_id: 0,
  teacher_id: 0,
  classroom_id: 0,
  pair_number: 1,
})

const conflicts = ref<ConflictDetail[]>([])
const showConflictModal = ref(false)

const hideGroupField = computed(() => !!props.initialGroupId)
const hideTeacherField = computed(() => !!props.initialTeacherId)

watch(
  () => props.entry,
  (entry) => {
    if (entry) {
      form.value = {
        group_id: props.initialGroupId ?? 0,
        date: props.initialDate ? new Date(props.initialDate) : new Date(),
        subject_id: entry.subject_id,
        teacher_id: entry.teacher_id,
        classroom_id: entry.classroom_id,
        pair_number: entry.pair_number,
      }
    }
  },
  { immediate: true }
)

watch(
  () => props.initialGroupId,
  (id) => {
    if (id) {
      form.value.group_id = id
    }
  },
  { immediate: true }
)

watch(
  () => props.initialTeacherId,
  (id) => {
    if (id) {
      form.value.teacher_id = id
    }
  },
  { immediate: true }
)

watch(
  () => props.initialDate,
  (date) => {
    if (date && !props.entry) {
      form.value.date = new Date(date)
    }
  },
  { immediate: true }
)

function getDateString(date: Date | string | undefined): string {
  if (!date) return ''
  const d = date instanceof Date ? date : new Date(date as string)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const bookedSlots = computed(() => {
  if (!form.value.teacher_id || !form.value.date) return new Set<number>()
  const dateStr = getDateString(form.value.date)
  const booked = new Set<number>()
  for (const entry of scheduleStore.entries) {
    if (entry.date !== dateStr) continue
    for (const pair of entry.pairs) {
      if (pair.teacher_id !== form.value.teacher_id) continue
      if (props.isEdit && props.entry && pair.id === props.entry.id) continue
      booked.add(pair.pair_number)
    }
  }
  return booked
})

const availablePairOptions = computed(() => {
  return Array.from({ length: 7 }, (_, i) => {
    const num = i + 1
    const isBooked = bookedSlots.value.has(num)
    return {
      label: `${num} (${PAIR_TIMES[num]})${isBooked ? ' — занято' : ''}`,
      value: num,
      disabled: isBooked,
    }
  })
})

watch(bookedSlots, (booked) => {
  if (booked.has(form.value.pair_number)) {
    const free = Array.from({ length: 7 }, (_, i) => i + 1).find((n) => !booked.has(n))
    if (free) form.value.pair_number = free
  }
})

const isValid = computed(() => {
  if (form.value.pair_number < 1 || form.value.pair_number > 7) return false
  if (bookedSlots.value.has(form.value.pair_number)) return false
  return (
    form.value.group_id > 0 &&
    form.value.date &&
    form.value.subject_id > 0 &&
    form.value.teacher_id > 0 &&
    form.value.classroom_id > 0 &&
    form.value.pair_number >= 1 &&
    form.value.pair_number <= 7
  )
})

function handleSave() {
  const dateStr = getDateString(form.value.date)
  
  emit('save', {
    group_id: form.value.group_id,
    date: dateStr,
    subject_id: form.value.subject_id,
    teacher_id: form.value.teacher_id,
    classroom_id: form.value.classroom_id,
    pair_number: form.value.pair_number,
  })
}

function handleForceSave() {
  const dateStr = getDateString(form.value.date)
  
  emit('save', {
    group_id: form.value.group_id,
    date: dateStr,
    subject_id: form.value.subject_id,
    teacher_id: form.value.teacher_id,
    classroom_id: form.value.classroom_id,
    pair_number: form.value.pair_number,
    force_save: true,
  })
  showConflictModal.value = false
}

defineExpose({ showConflict: (details: ConflictDetail[]) => {
  conflicts.value = details
  showConflictModal.value = true
} })
</script>

<template>
  <div class="pair-editor">
    <form @submit.prevent="handleSave">
      <div class="mb-3" v-if="!hideGroupField">
        <label class="form-label">Группа</label>
        <Select v-model="form.group_id" :options="groups" optionLabel="name" optionValue="id" placeholder="Выберите группу" class="w-100" />
      </div>
      <div class="mb-3">
        <label class="form-label">Дата</label>
        <DatePicker v-model="form.date" dateFormat="yy-mm-dd" class="w-100" />
      </div>
      <div class="mb-3">
        <label class="form-label">Предмет</label>
        <Select v-model="form.subject_id" :options="subjects" optionLabel="name" optionValue="id" placeholder="Выберите предмет" class="w-100" />
      </div>
      <div class="mb-3" v-if="!hideTeacherField">
        <label class="form-label">Преподаватель</label>
        <Select v-model="form.teacher_id" :options="teachers" optionLabel="name" optionValue="id" placeholder="Выберите преподавателя" class="w-100" />
      </div>
      <div class="mb-3">
        <label class="form-label">Аудитория</label>
        <Select v-model="form.classroom_id" :options="classrooms" optionLabel="number" optionValue="id" placeholder="Выберите аудиторию" class="w-100" />
      </div>
      <div class="mb-3">
        <label class="form-label">Номер пары</label>
        <Select v-model="form.pair_number" :options="availablePairOptions" optionLabel="label" optionValue="value" class="w-100" />
      </div>
      <div class="d-flex gap-2 justify-content-end mt-4">
        <Button label="Отмена" severity="secondary" text @click="emit('cancel')" />
        <Button label="Сохранить" type="submit" :disabled="!isValid" />
      </div>
    </form>
    <ConflictModal
      v-if="showConflictModal"
      :details="conflicts"
      @force-save="handleForceSave"
      @cancel="showConflictModal = false"
    />
  </div>
</template>