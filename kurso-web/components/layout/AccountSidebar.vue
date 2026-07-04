<script setup lang="ts">
const route = useRoute()
const { user, logout } = useAuth()
const { activeCount } = useAlerts()

interface Item {
  to: string
  label: string
  icon: string
  badge?: () => number
}

// Icon inner-SVG (24×24, stroke=currentColor). Static/trusted content.
const nav: Item[] = [
  {
    to: '/account',
    label: 'Обзор',
    icon: '<rect x="3" y="3" width="7" height="9" rx="1.5"/><rect x="14" y="3" width="7" height="5" rx="1.5"/><rect x="14" y="12" width="7" height="9" rx="1.5"/><rect x="3" y="16" width="7" height="5" rx="1.5"/>',
  },
  {
    to: '/account/alerts',
    label: 'Алерты',
    icon: '<path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M10.5 20a1.8 1.8 0 0 0 3 0"/>',
    badge: () => activeCount.value,
  },
  {
    to: '/account/favorites',
    label: 'Избранное',
    icon: '<path d="M12 20s-6.5-4-9-8.5A4.5 4.5 0 0 1 12 7a4.5 4.5 0 0 1 9 4.5C18.5 16 12 20 12 20Z"/>',
  },
  {
    to: '/account/history',
    label: 'История',
    icon: '<circle cx="12" cy="12" r="9"/><path d="M12 8v4l3 2"/>',
  },
  {
    to: '/account/reviews',
    label: 'Отзывы',
    icon: '<path d="M21 15a2 2 0 0 1-2 2H8l-4 4V6a2 2 0 0 1 2-2h13a2 2 0 0 1 2 2Z"/>',
  },
  {
    to: '/account/partner',
    label: 'Партнёрка',
    icon: '<circle cx="12" cy="12" r="9"/><path d="M9.5 9a2.5 2.5 0 0 1 4.5 1.5c0 1.5-2.5 2-2.5 3.5M12 17h.01"/>',
  },
  {
    to: '/account/settings',
    label: 'Настройки',
    icon: '<circle cx="12" cy="12" r="3"/><path d="M19 12a7 7 0 0 0-.1-1l2-1.5-2-3.5-2.4 1a7 7 0 0 0-1.7-1L14.5 3h-5l-.3 2.5a7 7 0 0 0-1.7 1l-2.4-1-2 3.5 2 1.5a7 7 0 0 0 0 2l-2 1.5 2 3.5 2.4-1a7 7 0 0 0 1.7 1l.3 2.5h5l.3-2.5a7 7 0 0 0 1.7-1l2.4 1 2-3.5-2-1.5a7 7 0 0 0 .1-1Z"/>',
  },
]

const isActive = (to: string) =>
  to === '/account' ? route.path === '/account' : route.path.startsWith(to)

function onLogout() {
  logout()
  navigateTo('/login')
}
</script>

<template>
  <aside class="sticky top-6 flex flex-col gap-1">
    <!-- user card -->
    <div class="mb-2 flex items-center gap-3 px-2 pb-4">
      <span
        class="flex h-10 w-10 flex-none items-center justify-center rounded-full bg-[#3A4452] text-[15px] font-bold text-white"
        >{{ user?.initials ?? 'K' }}</span
      >
      <div class="min-w-0">
        <div class="truncate text-sm font-semibold text-ink">{{ user?.name ?? 'Гость' }}</div>
        <div class="text-xs text-ink-faint">{{ user?.plan ?? '' }} план</div>
      </div>
    </div>

    <NuxtLink
      v-for="item in nav"
      :key="item.to"
      :to="item.to"
      class="flex items-center gap-[11px] rounded-lg px-3.5 py-[11px] text-sm transition-colors"
      :class="
        isActive(item.to)
          ? 'bg-brand/[0.12] font-semibold text-brand-bright'
          : 'text-ink-muted hover:text-ink'
      "
    >
      <!-- eslint-disable-next-line vue/no-v-html -->
      <svg
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        v-html="item.icon"
      />
      <span class="flex-1">{{ item.label }}</span>
      <span
        v-if="item.badge && item.badge()"
        class="tnum rounded-md bg-brand/15 px-1.5 py-0.5 text-[11px] font-semibold text-brand-bright"
        >{{ item.badge() }}</span
      >
    </NuxtLink>

    <button
      type="button"
      class="mt-2 flex items-center gap-[11px] rounded-lg px-3.5 py-[11px] text-sm text-ink-faint transition-colors hover:text-danger"
      @click="onLogout"
    >
      <svg
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9" />
      </svg>
      Выйти
    </button>
  </aside>
</template>
