<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from 'primevue/button'

const router = useRouter()
const authStore = useAuthStore()

function logout() {
  authStore.logout()
  router.push('/login')
}

const menuItems = [
  { label: 'Преподаватели', icon: 'pi pi-users', command: () => router.push('/admin/teachers') },
  { label: 'Предметы', icon: 'pi pi-book', command: () => router.push('/admin/subjects') },
  { label: 'Аудитории', icon: 'pi pi-home', command: () => router.push('/admin/classrooms') },
  { label: 'Группы', icon: 'pi pi-users', command: () => router.push('/admin/groups') },
]
</script>

<template>
  <div class="admin-layout d-flex min-vh-100">
    <aside class="sidebar p-3 d-flex flex-column">
      <h2 class="h5 mb-3 text-white">Админ-панель</h2>
      <nav class="d-flex flex-column gap-2 flex-grow-1">
        <Button
          v-for="item in menuItems"
          :key="item.label"
          :label="item.label"
          :icon="item.icon"
          severity="secondary"
          text
          class="text-white text-start"
          @click="item.command"
        />
      </nav>
      <Button
        label="Выйти"
        icon="pi pi-sign-out"
        severity="danger"
        class="mt-auto"
        @click="logout"
      />
    </aside>
    <main class="content flex-grow-1 p-3 overflow-auto">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.sidebar {
  width: 250px;
  background: #212529;
}
.content {
  background: #f8f9fa;
}
</style>