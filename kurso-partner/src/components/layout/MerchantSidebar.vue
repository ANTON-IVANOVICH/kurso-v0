<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { NAV, type NavItem, type BadgeTone } from '../../lib/nav'
import { useDashboardQuery, useComplaintsQuery } from '../../composables/useMerchant'
import UiIcon from '../ui/UiIcon.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ navigate: [] }>()

const route = useRoute()
const auth = useAuthStore()

// Live badge sources — Colada singletons shared with the pages.
const dashboard = useDashboardQuery()
const complaints = useComplaintsQuery()

const openComplaints = computed(
  () => (complaints.data.value ?? []).filter((c) => c.status === 'open').length,
)

function isActive(to: string): boolean {
  if (to === '/') return route.path === '/'
  return route.path === to || route.path.startsWith(to + '/')
}

const BADGE_COLOR: Record<BadgeTone, string> = {
  warn: 'text-warn-text',
  danger: 'text-danger-text',
  brand: 'text-brand-light',
}

function badgeFor(item: NavItem): { text: string; tone: BadgeTone } | null {
  const m = dashboard.data.value?.metrics
  if (item.badge === 'ratesStale' && m?.ratesStale)
    return { text: String(m.ratesStale), tone: 'warn' }
  if (item.badge === 'reviewsUnanswered' && m?.reviewsUnanswered)
    return { text: String(m.reviewsUnanswered), tone: 'brand' }
  if (item.badge === 'complaintsOpen' && openComplaints.value)
    return { text: String(openComplaints.value), tone: 'danger' }
  return null
}
</script>

<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 flex w-[240px] flex-col border-r border-chrome bg-panel px-3 py-4 text-[13px] transition-transform md:sticky md:top-0 md:h-screen md:translate-x-0"
    :class="open ? 'translate-x-0' : '-translate-x-full'"
  >
    <!-- exchanger identity -->
    <div class="mb-3 flex items-center gap-2.5 border-b border-chrome px-1.5 pb-4 pt-1">
      <span
        class="flex h-9 w-9 flex-none items-center justify-center rounded-[10px] bg-[#3A4452] text-[13px] font-extrabold text-white"
      >
        {{ auth.initials }}
      </span>
      <div class="min-w-0">
        <div class="truncate text-[14px] font-bold">
          {{ auth.merchant?.exchangerName ?? 'Кабинет' }}
        </div>
        <div class="text-[11px] text-ink-faint">кабинет обменника</div>
      </div>
    </div>

    <!-- nav -->
    <nav class="min-h-0 flex-1 overflow-y-auto">
      <RouterLink
        v-for="item in NAV"
        :key="item.to"
        :to="item.to"
        class="mb-0.5 flex items-center gap-2.5 rounded-[10px] px-3 py-2.5 transition-colors"
        :class="
          isActive(item.to)
            ? 'bg-brand/[0.12] font-semibold text-brand-light'
            : 'text-ink-muted hover:bg-white/[0.03] hover:text-ink'
        "
        @click="emit('navigate')"
      >
        <UiIcon :name="item.icon" :size="16" />
        <span>{{ item.label }}</span>
        <span
          v-if="badgeFor(item)"
          class="tnum ml-auto rounded-md px-1.5 py-0.5 text-[11px] font-semibold"
          :class="BADGE_COLOR[badgeFor(item)!.tone]"
          >{{ badgeFor(item)!.text }}</span
        >
      </RouterLink>
    </nav>

    <!-- user / logout -->
    <div class="mt-2 flex items-center gap-2.5 border-t border-chrome px-1.5 pt-3">
      <span
        class="flex h-8 w-8 flex-none items-center justify-center rounded-lg bg-brand text-[11px] font-bold text-white"
      >
        {{ auth.initials }}
      </span>
      <div class="min-w-0 flex-1">
        <div class="truncate text-[12px] font-semibold">
          {{ auth.merchant?.email ?? 'partner@kurso.io' }}
        </div>
        <div class="text-[10px] capitalize text-ink-faint">
          {{ auth.merchant?.role ?? 'owner' }}
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
