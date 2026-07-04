<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCatalogStore } from '../../stores/catalog'
import { useAuthStore } from '../../stores/auth'
import AdminSidebar from './AdminSidebar.vue'
import UiIcon from '../ui/UiIcon.vue'

const catalog = useCatalogStore()
const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const drawerOpen = ref(false)

onMounted(() => catalog.load())

// Close the mobile drawer whenever the route changes.
watch(
  () => route.fullPath,
  () => (drawerOpen.value = false),
)

// Bounce to /login the moment the session is cleared.
watch(
  () => auth.isAuthenticated,
  (ok) => {
    if (!ok) router.replace({ name: 'login' })
  },
)
</script>

<template>
  <div class="flex min-h-screen bg-bg">
    <AdminSidebar :open="drawerOpen" @navigate="drawerOpen = false" />

    <!-- backdrop for the mobile drawer -->
    <div
      v-if="drawerOpen"
      class="fixed inset-0 z-30 bg-black/50 md:hidden"
      @click="drawerOpen = false"
    />

    <div class="flex min-w-0 flex-1 flex-col">
      <!-- mobile topbar -->
      <div
        class="sticky top-0 z-20 flex items-center gap-3 border-b border-chrome bg-bg/90 px-4 py-3 backdrop-blur md:hidden"
      >
        <button
          class="rounded-lg border border-line-strong bg-surface p-2 text-ink-muted"
          aria-label="Меню"
          @click="drawerOpen = true"
        >
          <UiIcon name="grid" :size="16" />
        </button>
        <span class="text-sm font-bold">Kurso Admin</span>
      </div>

      <main class="mx-auto w-full max-w-[1320px] flex-1 px-4 py-6 md:px-8 md:py-8">
        <RouterView />
      </main>
    </div>
  </div>
</template>
