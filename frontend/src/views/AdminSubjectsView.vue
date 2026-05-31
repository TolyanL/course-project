<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { subjectsApi } from '@/api/client'
import EntityTable from '@/components/EntityTable.vue'
import EntityForm from '@/components/EntityForm.vue'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type { Subject, SubjectFormData } from '@/types'

const confirm = useConfirm()

const subjects = ref<Subject[]>([])
const showForm = ref(false)
const editingItem = ref<Subject | undefined>(undefined)

const columns = [
  { key: 'id', label: 'ID' },
  { key: 'name', label: 'Название' },
]

const fields = [{ key: 'name', label: 'Название', type: 'text' as const, required: true }]

onMounted(async () => {
  await loadSubjects()
})

async function loadSubjects() {
  const { data } = await subjectsApi.getAll()
  subjects.value = data
}

function handleAdd() {
  editingItem.value = undefined
  showForm.value = true
}

function handleEdit(item: Record<string, unknown>) {
  editingItem.value = item as unknown as Subject
  showForm.value = true
}

async function handleDelete(id: number) {
  confirm.require({
    message: 'Вы уверены, что хотите удалить этот предмет?',
    header: 'Подтверждение',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await subjectsApi.delete(id)
        await loadSubjects()
      } catch (e: unknown) {
        const error = e as { response?: { data?: { error?: string } } }
        if (error.response?.data?.error) {
          alert(error.response.data.error)
        }
      }
    },
  })
}

async function handleSubmit(data: Record<string, unknown>) {
  const formData = data as unknown as SubjectFormData
  if (editingItem.value) {
    await subjectsApi.update(editingItem.value.id, formData)
  } else {
    await subjectsApi.create(formData)
  }
  showForm.value = false
  await loadSubjects()
}

function handleCancel() {
  showForm.value = false
  editingItem.value = undefined
}
</script>

<template>
  <div class="admin-subjects-view">
    <h1 class="h4 mb-3">Управление предметами</h1>
    <EntityTable
      :columns="columns"
      :data="subjects"
      @add="handleAdd"
      @edit="handleEdit"
      @delete="handleDelete"
    />
    <EntityForm
      :fields="fields"
      :initial-data="editingItem"
      :is-edit="!!editingItem"
      v-model:visible="showForm"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
    <ConfirmDialog />
  </div>
</template>
