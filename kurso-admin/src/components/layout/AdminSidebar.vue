<script setup lang="ts">
import { useRoute } from 'vue-router'
import { NAV, type NavItem, type BadgeTone } from '../../lib/nav'
import { useCatalog } from '../../composables/useCatalog'
import { useAuthStore } from '../../stores/auth'
import UiIcon from '../ui/UiIcon.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ navigate: [] }>()

const route = useRoute()
const catalog = useCatalog()
const auth = useAuthStore()

function isActive(to: string): boolean {
  if (to === '/') return route.path === '/'
  return route.path === to || route.path.startsWith(to + '/')
}

const BADGE_COLOR: Record<BadgeTone, string> = {
  faint: 'text-ink-faint',
  warn: 'text-warn-text',
  danger: 'text-danger-text',
  brand: 'text-brand-light',
}

function badgeFor(item: NavItem): { text: string; tone: BadgeTone } | null {
  if (item.dynamic === 'exchangers') {
    const n = catalog.exchangers.length
    return n ? { text: String(n), tone: 'faint' } : null
  }
  return item.badge ?? null
}
</script>

<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 flex w-[232px] flex-col border-r border-chrome bg-panel px-2.5 py-3.5 text-[13px] transition-transform md:sticky md:top-0 md:h-screen md:translate-x-0"
    :class="open ? 'translate-x-0' : '-translate-x-full'"
  >
    <!-- brand -->
    <div class="mb-3 flex items-center gap-2.5 border-b border-chrome px-2 pb-3.5 pt-1">
      <span
        class="flex h-[30px] w-[30px] items-center justify-center rounded-lg bg-brand text-[13px] font-extrabold text-white"
        >K</span
      >
      <div>
        <div class="text-[13px] font-bold">Kurso Admin</div>
        <div class="text-[10px] text-ink-faint">v2.4.0</div>
      </div>
    </div>

    <!-- nav -->
    <nav class="min-h-0 flex-1 overflow-y-auto">
      <div v-for="group in NAV" :key="group.title" class="mb-1">
        <div class="px-2.5 py-1.5 text-[10px] uppercase tracking-[0.1em] text-ink-fainter">
          {{ group.title }}
        </div>
        <RouterLink
          v-for="item in group.items"
          :key="item.to"
          :to="item.to"
          class="mb-0.5 flex items-center gap-2.5 rounded-lg px-2.5 py-2 transition-colors"
          :class="
            isActive(item.to)
              ? 'bg-brand/[0.12] font-semibold text-brand-light'
              : 'text-ink-muted hover:bg-white/[0.03] hover:text-ink'
          "
          @click="emit('navigate')"
        >
          <UiIcon v-if="item.icon === 'grid'" name="grid" :size="15" />
          <span>{{ item.label }}</span>
          <span
            v-if="badgeFor(item)"
            class="tnum ml-auto text-[11px]"
            :class="BADGE_COLOR[badgeFor(item)!.tone]"
            >{{ badgeFor(item)!.text }}</span
          >
        </RouterLink>
      </div>
    </nav>

    <!-- user -->
    <div class="mt-2 flex items-center gap-2.5 border-t border-chrome px-1.5 pt-3">
      <span
        class="flex h-8 w-8 flex-none items-center justify-center rounded-lg bg-brand text-[11px] font-bold text-white"
      >
        {{ auth.user?.initials ?? 'AK' }}
      </span>
      <div class="min-w-0 flex-1">
        <div class="truncate text-[12px] font-semibold">
          {{ auth.user?.email ?? 'admin@kurso.io' }}
        </div>
        <div class="text-[10px] capitalize text-ink-faint">
          {{ auth.user?.role ?? 'superadmin' }}
        </div>
      </div>
      <button
        class="flex-none rounded-md p-1.5 text-ink-faint transition-colors hover:bg-white/[0.04] hover:text-danger-text"
        aria-label="Выйти"
        @click="auth.logout()"
      >
        <UiIcon name="logout" :size="16" />
      </button>
    </div>
  </aside>
</template>
