<script setup lang="ts">
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import type { ConflictDetail } from '@/types'
import { useAuthStore } from '@/stores/auth'

defineProps<{
  details: ConflictDetail[]
}>()

const emit = defineEmits<{
  forceSave: []
  cancel: []
}>()

const authStore = useAuthStore()
</script>

<template>
  <Dialog :visible="true" header="Конфликт" :modal="true" :closable="false" :style="{ width: '400px' }">
    <div class="conflicts mb-3">
      <div
        v-for="(detail, index) in details"
        :key="index"
        class="conflict-item p-2 rounded mb-2"
        :class="detail.type === 'classroom' ? 'bg-primary-subtle' : 'bg-danger-subtle'"
      >
        <strong>{{ detail.type === 'classroom' ? 'Аудитория' : 'Преподаватель' }}:</strong>
        {{ detail.message }}
      </div>
    </div>
    <div class="d-flex gap-2 justify-content-end">
      <Button label="Отмена" severity="secondary" text @click="emit('cancel')" />
      <Button
        v-if="authStore.role === 'teacher'"
        label="Сохранить принудительно"
        severity="warning"
        @click="emit('forceSave')"
      />
    </div>
  </Dialog>
</template>