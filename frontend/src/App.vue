<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Menubar from 'primevue/menubar'
import Button from 'primevue/button'

const router = useRouter()
const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const role = computed(() => authStore.role)
const name = computed(() => authStore.name)

function logout() {
  authStore.logout()
  router.push('/login')
}

function goToSchedule() {
  router.push('/')
}
</script>

<template>
  <div id="app">
    <header v-if="isAuthenticated" class="app-header">
      <Menubar :model="[]" class="border-0 bg-transparent">
        <template #start>
          <Button
            label="Электронное расписание"
            text
            class="text-white fw-bold"
            @click="goToSchedule"
          />
        </template>
        <template #end>
          <div class="d-flex align-items-center gap-3">
            <span v-if="role === 'admin'" class="text-white">
              <i class="pi pi-crown"></i>
              {{ name }}
            </span>
            <span v-else-if="role === 'teacher'" class="text-white">
              <i class="pi pi-graduation-cap"></i>
              {{ name }}
            </span>
            <Button
              label="Выйти"
              icon="pi pi-sign-out"
              severity="secondary"
              size="small"
              @click="logout"
            />
          </div>
        </template>
      </Menubar>
    </header>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
html,
body {
  height: 100%;
}
body {
  font-family:
    -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background: #f5f5f5;
  color: #333;
}
#app {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.app-header {
  background: #0f6cbd;
  color: white;
  padding: 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.app-header :deep(.p-menubar) {
  background: transparent;
  border: none;
  padding: 0.5rem 1rem;
  max-width: 1200px;
  margin: 0 auto;
}
.main-content {
  flex: 1;
  padding: 1rem;
}

@media (max-width: 768px) {
  .app-header :deep(.p-menubar) {
    padding: 0.5rem 0.75rem;
  }
  .app-header :deep(.p-menubar-start) .p-button {
    font-size: 0.875rem;
    padding: 0.5rem 0.75rem;
  }
  .app-header :deep(.p-menubar-end) .d-flex {
    gap: 0.5rem;
  }
  .app-header :deep(.p-menubar-end) .p-button .p-button-label {
    display: none;
  }
  .app-header :deep(.p-menubar-end) span {
    font-size: 0.8rem;
  }
  .main-content {
    padding: 0.75rem;
  }
}

@media (max-width: 480px) {
  .app-header :deep(.p-menubar-start) .p-button {
    font-size: 0.75rem;
    padding: 0.375rem 0.5rem;
  }
}
</style>
