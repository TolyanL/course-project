<script setup lang="ts">
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'

defineProps<{
  columns: { key: string; label: string }[]
  data: Record<string, unknown>[]
}>()

const emit = defineEmits<{
  edit: [item: Record<string, unknown>]
  delete: [id: number]
  add: []
}>()
</script>

<template>
  <div class="entity-table">
    <div class="table-header mb-3">
      <Button label="Добавить" icon="pi pi-plus" @click="emit('add')" />
    </div>
    <DataTable :value="data" :paginator="data.length > 10" :rows="10" stripedRows responsiveLayout="scroll">
      <Column v-for="col in columns" :key="col.key" :field="col.key" :header="col.label" />
      <Column header="Действия" :exportable="false" style="width: 200px">
        <template #body="{ data }">
          <div class="d-flex gap-2">
            <Button
              icon="pi pi-pencil"
              severity="info"
              text
              rounded
              title="Редактировать"
              @click="emit('edit', data)"
            />
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              rounded
              title="Удалить"
              @click="emit('delete', data.id as number)"
            />
          </div>
        </template>
      </Column>
      <template #empty>
        <div class="text-center py-4 text-muted">Нет данных</div>
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.entity-table {
  background: white;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.entity-table :deep(.p-datatable) {
  overflow-x: auto;
}

.entity-table :deep(.p-datatable-wrapper) {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

@media (max-width: 768px) {
  .entity-table {
    padding: 0.75rem;
  }

  .table-header {
    display: flex;
    justify-content: flex-end;
  }
}
</style>