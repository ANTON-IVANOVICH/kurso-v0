<script setup lang="ts">
import { computed } from 'vue'

// Layout for nested content/legal pages (terms, privacy, cookies).
// Desktop: full site chrome. Mobile: a nested screen — back bar instead of the
// site header, no bottom navigation, and a slim legal footer.
const { t } = useI18n()
const route = useRoute()
const title = computed(() => (route.meta.legalTitle as string | undefined) ?? '')

// Slim legal footer per the Content & Legal design. "О проекте" only appears on
// desktop (the mobile design omits it); the current page is highlighted.
const footerLinks = [
  { key: 'terms', to: '/terms' },
  { key: 'privacy', to: '/privacy' },
  { key: 'cookies', to: '/cookies' },
  { key: 'about', to: '/about', desktopOnly: true },
  { key: 'contact', to: '/contact' },
]
</script>

<template>
  <div class="relative min-h-screen overflow-x-clip">
    <GlowBackdrop />

    <!-- desktop: full site header -->
    <div class="hidden md:block">
      <SiteHeader />
    </div>
    <!-- mobile: nested-screen back bar -->
    <div class="md:hidden">
      <LegalBackBar :title="title" />
    </div>

    <main class="relative z-10 pb-12 md:pb-8">
      <slot />
    </main>

    <!-- slim legal footer: copyright + inline links (design Content & Legal) -->
    <footer
      class="relative z-10 mt-8 border-t border-line-subtle px-4 py-6 md:mt-16 md:px-6 md:py-7"
    >
      <div
        class="mx-auto flex max-w-[1100px] flex-col gap-4 md:flex-row md:items-center md:justify-between"
      >
        <span class="tnum order-last text-[11px] text-ink-faint md:order-none md:text-xs"
          >© 2026 Kurso</span
        >
        <nav class="flex flex-wrap gap-x-[22px] gap-y-2.5 text-[13px]">
          <NuxtLink
            v-for="link in footerLinks"
            :key="link.to"
            :to="link.to"
            class="transition-colors"
            :class="[
              route.path === link.to ? 'font-semibold text-ink' : 'text-ink-muted hover:text-ink',
              link.desktopOnly ? 'hidden md:inline' : '',
            ]"
            >{{ t(`legalNav.${link.key}`) }}</NuxtLink
          >
        </nav>
      </div>
    </footer>
  </div>
</template>
