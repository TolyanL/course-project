<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { groupsApi } from '@/api/client'
import EntityTable from '@/components/EntityTable.vue'
import EntityForm from '@/components/EntityForm.vue'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type { Group, GroupFormData } from '@/types'

const confirm = useConfirm()

const groups = ref<Group[]>([])
const showForm = ref(false)
const editingItem = ref<Group | undefined>(undefined)

const columns = [
  { key: 'id', label: 'ID' },
  { key: 'name', label: 'Название' },
]

const fields = [{ key: 'name', label: 'Название группы', type: 'text' as const, required: true }]

onMounted(async () => {
  await loadGroups()
})

async function loadGroups() {
  const { data } = await groupsApi.getAll()
  groups.value = data
}

function handleAdd() {
  editingItem.value = undefined
  showForm.value = true
}

function handleEdit(item: Record<string, unknown>) {
  editingItem.value = item as unknown as Group
  showForm.value = true
}

async function handleDelete(id: number) {
  confirm.require({
    message: 'Вы уверены, что хотите удалить эту группу?',
    header: 'Подтверждение',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      await groupsApi.delete(id)
      await loadGroups()
    },
  })
}

async function handleSubmit(data: Record<string, unknown>) {
  const formData = data as unknown as GroupFormData
  if (editingItem.value) {
    await groupsApi.update(editingItem.value.id, formData)
  } else {
    await groupsApi.create(formData)
  }
  showForm.value = false
  await loadGroups()
}

function handleCancel() {
  showForm.value = false
  editingItem.value = undefined
}
</script>

<template>
  <div class="admin-groups-view">
    <h1 class="h4 mb-3">Управление группами</h1>
    <EntityTable
      :columns="columns"
      :data="groups"
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
