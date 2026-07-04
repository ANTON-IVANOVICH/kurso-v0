<script setup lang="ts">
const route = useRoute()
const { openSearch } = useSearchOverlay()
const { openNotifications } = useNotifications()
const { t } = useI18n()
const nav = [
  { key: 'nav.exchange', to: '/' },
  { key: 'nav.exchangers', to: '/exchangers' },
  { key: 'nav.map', to: '/map' },
]
const isActive = (to: string) => (to === '/' ? route.path === '/' : route.path.startsWith(to))
</script>

<template>
  <header class="relative z-20">
    <div class="mx-auto flex max-w-[1200px] items-center gap-7 px-4 py-4 md:px-6">
      <NuxtLink to="/" class="flex-none"><AppLogo :size="30" /></NuxtLink>

      <nav class="hidden items-center gap-1 text-[15px] md:flex">
        <NuxtLink
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="rounded-md px-4 py-2 transition-colors"
          :class="
            isActive(item.to)
              ? 'bg-brand/[0.12] font-semibold text-brand-bright'
              : 'text-ink-muted hover:text-ink'
          "
          >{{ t(item.key) }}</NuxtLink
        >
      </nav>

      <div class="ml-auto flex items-center gap-2.5">
        <!-- interface language -->
        <LangSwitcher class="hidden md:block" />
        <LangSwitcher compact class="md:hidden" />

        <button
          type="button"
          :aria-label="t('header.searchAria')"
          class="flex h-[38px] w-[38px] items-center justify-center rounded-full border border-line bg-surface text-ink-muted md:hidden"
          @click="openSearch"
        >
          <svg
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          >
            <circle cx="11" cy="11" r="7" />
            <path d="M21 21l-3.5-3.5" />
          </svg>
        </button>
        <button
          type="button"
          :aria-label="t('header.notificationsAria')"
          class="flex h-[38px] w-[38px] items-center justify-center rounded-full border border-line bg-surface text-ink-muted md:hidden"
          @click="openNotifications"
        >
          <svg
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
            <path d="M10.5 20a1.8 1.8 0 0 0 3 0" />
          </svg>
        </button>

        <KButton pill class="hidden md:inline-flex" @click="navigateTo('/login')">{{
          t('header.login')
        }}</KButton>
      </div>
    </div>
  </header>
</template>
