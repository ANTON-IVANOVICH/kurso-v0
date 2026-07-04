<script setup lang="ts">
// Items with `to` navigate; those without are sections not built yet and render
// as muted labels (so the footer matches the design without dead 404 links).
const route = useRoute()
const { t } = useI18n()
interface Link {
  key: string
  to?: string
}
const columns: { titleKey: string; links: Link[] }[] = [
  {
    titleKey: 'footer.colProduct',
    links: [
      { key: 'nav.exchange', to: '/' },
      { key: 'nav.exchangers', to: '/exchangers' },
      { key: 'nav.map', to: '/map' },
      { key: 'footer.alerts' },
    ],
  },
  {
    titleKey: 'footer.colCompany',
    links: [
      { key: 'footer.about', to: '/about' },
      { key: 'footer.contacts', to: '/contact' },
      { key: 'footer.blog' },
      { key: 'footer.forPartners' },
    ],
  },
  {
    titleKey: 'footer.colLegal',
    links: [
      { key: 'footer.terms', to: '/terms' },
      { key: 'footer.privacy', to: '/privacy' },
      { key: 'footer.cookies', to: '/cookies' },
    ],
  },
]
</script>

<template>
  <footer class="relative z-10 mt-24 border-t border-line-subtle bg-canvas">
    <div class="mx-auto max-w-[1200px] px-4 py-14 md:px-6">
      <div class="grid grid-cols-2 gap-10 md:grid-cols-[1.7fr_1fr_1fr_1fr]">
        <div class="col-span-2 md:col-span-1">
          <AppLogo :size="32" />
          <p class="mt-4 max-w-xs text-sm leading-relaxed text-ink-faint">
            {{ t('footer.tagline') }}
          </p>
        </div>
        <div v-for="col in columns" :key="col.titleKey">
          <div class="text-[11px] uppercase tracking-[0.08em] text-ink-faint">
            {{ t(col.titleKey) }}
          </div>
          <ul class="mt-3.5 space-y-2.5">
            <li v-for="link in col.links" :key="link.key">
              <NuxtLink
                v-if="link.to"
                :to="link.to"
                class="text-sm transition-colors"
                :class="
                  route.path === link.to
                    ? 'font-semibold text-ink'
                    : 'text-ink-muted hover:text-ink'
                "
                >{{ t(link.key) }}</NuxtLink
              >
              <span v-else class="text-sm text-ink-faint">{{ t(link.key) }}</span>
            </li>
          </ul>
        </div>
      </div>
      <div
        class="mt-12 flex flex-col gap-3 border-t border-line-subtle pt-7 text-xs text-ink-faint md:flex-row md:items-center md:justify-between"
      >
        <span class="tnum">{{ t('footer.rights') }}</span>
        <span>{{ t('footer.riskNote') }}</span>
      </div>
    </div>
  </footer>
</template>
