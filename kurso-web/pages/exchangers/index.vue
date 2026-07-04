<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import type { Exchanger } from '~/composables/useApi'

// Каталог обменников — вариант A (таблица на десктопе, компактный список на мобиле).
// Данные — живые (useExchangersQuery); резерв/направления/активы — реальные
// агрегаты из API. Фильтры, поиск, сортировка и пагинация работают на клиенте.

const { state, data } = useExchangersQuery()
const { t, plural } = useI18n()
const all = computed<Exchanger[]>(() => data.value ?? [])
const loading = computed(() => state.value.status === 'pending' && all.value.length === 0)

// --- фильтры / поиск / сортировка / пагинация ---
type FilterKey = 'all' | 'partners' | 'verified'
type SortKey = 'rating' | 'reserve' | 'reviews'

const filter = ref<FilterKey>('all')
const search = ref('')
const sort = ref<SortKey>('rating')
const PAGE = 12
const shown = ref(PAGE)

const filters: { key: FilterKey; i18n: string }[] = [
  { key: 'all', i18n: 'exchangers.fAll' },
  { key: 'partners', i18n: 'exchangers.fPartners' },
  { key: 'verified', i18n: 'exchangers.fVerified' },
]
const sortOptions: { key: SortKey; i18n: string }[] = [
  { key: 'rating', i18n: 'exchangers.sortRating' },
  { key: 'reserve', i18n: 'exchangers.sortReserve' },
  { key: 'reviews', i18n: 'exchangers.sortReviews' },
]
const sortLabel = computed(() => t(sortOptions.find((o) => o.key === sort.value)?.i18n ?? ''))
function cycleSort() {
  const i = sortOptions.findIndex((o) => o.key === sort.value)
  sort.value = sortOptions[(i + 1) % sortOptions.length].key
}
function setFilter(key: FilterKey) {
  filter.value = key
  shown.value = PAGE
}

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  return all.value.filter((e) => {
    if (filter.value === 'partners' && !e.partner) return false
    if (filter.value === 'verified' && !e.isVerified) return false
    if (q && !e.name.toLowerCase().includes(q)) return false
    return true
  })
})
const sorted = computed(() => {
  const list = [...filtered.value]
  list.sort((a, b) => {
    if (sort.value === 'reserve') return reserveNum(b) - reserveNum(a)
    if (sort.value === 'reviews') return b.reviewsCount - a.reviewsCount
    return (b.ratingAvg ?? 0) - (a.ratingAvg ?? 0)
  })
  return list
})
const visible = computed(() => sorted.value.slice(0, shown.value))
const canLoadMore = computed(() => shown.value < sorted.value.length)
const remaining = computed(() => sorted.value.length - shown.value)
const partners = computed(() => all.value.filter((e) => e.partner).slice(0, 3))

function reserveNum(e: Exchanger): number {
  return e.reserveTotal ? parseFloat(e.reserveTotal) : 0
}

// --- VM для строки/карточки ---
interface Row {
  slug: string
  name: string
  initials: string
  color: string
  partner: boolean
  verified: boolean
  rating: string
  reviews: number
  reserve: string
  assets: string[]
  directions: number
  onSince: number | null
  href: string
}
function toRow(e: Exchanger): Row {
  const av = exchangerAvatar(e.slug, e.name)
  return {
    slug: e.slug,
    name: e.name,
    initials: av.initials,
    color: av.color,
    partner: e.partner,
    verified: e.isVerified,
    rating: e.ratingAvg != null ? e.ratingAvg.toFixed(1) : '—',
    reviews: e.reviewsCount,
    reserve: e.reserveTotal ? `${fmtCompact(parseFloat(e.reserveTotal))} ₽` : '—',
    assets: e.assets ?? [],
    directions: e.directionsCount ?? 0,
    onSince: e.onSince || null,
    // Links to the exchanger's detail card; the outbound clickout lives there.
    href: `/exchangers/${e.slug}`,
  }
}
const rows = computed(() => visible.value.map(toRow))
const partnerRows = computed(() => partners.value.map(toRow))
const total = computed(() => all.value.length)

// --- «обновлено N сек назад» (клиентские часы, без SSR-рассинхрона) ---
const mounted = ref(false)
const now = ref(0)
const loadedAt = ref(0)
let timer: ReturnType<typeof setInterval> | undefined
onMounted(() => {
  mounted.value = true
  now.value = Date.now()
  loadedAt.value = Date.now()
  timer = setInterval(() => (now.value = Date.now()), 1000)
})
onBeforeUnmount(() => clearInterval(timer))
const updatedAgo = computed(() => {
  if (!mounted.value) return ''
  const sec = Math.max(0, Math.round((now.value - loadedAt.value) / 1000))
  if (sec < 5) return t('exchangers.updatedNow')
  if (sec < 60) return t('exchangers.updatedSec', { n: sec })
  const min = Math.floor(sec / 60)
  if (min < 60) return t('exchangers.updatedMin', { n: min })
  return t('exchangers.updatedHour', { n: Math.floor(min / 60) })
})

useSeoMeta({
  title: () => `${t('exchangers.titleDesktop')} — Kurso`,
  description: () => t('exchangers.subtitle'),
})

const openBtn =
  'inline-flex flex-none items-center justify-center whitespace-nowrap rounded-lg px-[18px] py-[11px] text-sm font-semibold transition-colors'
</script>

<template>
  <div class="mx-auto max-w-[1320px] px-4 py-6 md:px-6 md:py-8">
    <!-- ===== Заголовок ===== -->
    <div class="mb-5 flex items-end justify-between gap-4 md:mb-6">
      <div>
        <h1 class="text-[22px] font-extrabold tracking-[-0.025em] md:text-[30px]">
          <span class="md:hidden">{{ t('exchangers.titleMobile') }}</span>
          <span class="hidden md:inline">{{ t('exchangers.titleDesktop') }}</span>
        </h1>
        <p class="mt-2 hidden text-[15px] text-ink-muted md:block">
          <span class="tnum font-semibold text-ink">{{ total }}</span>
          {{ t('exchangers.subtitle') }}
        </p>
      </div>
      <ClientOnly>
        <span
          v-if="updatedAgo"
          class="hidden items-center gap-2 whitespace-nowrap text-[13px] text-ink-faint md:inline-flex"
        >
          <span class="h-2 w-2 flex-none rounded-full bg-success animate-kpulse" />
          {{ updatedAgo }}
        </span>
      </ClientOnly>
    </div>

    <!-- ===== Полоса партнёров (десктоп) ===== -->
    <div v-if="partnerRows.length" class="mb-6 hidden gap-3.5 md:grid md:grid-cols-3">
      <div
        v-for="p in partnerRows"
        :key="p.slug"
        class="relative overflow-hidden rounded-2xl border border-[#24324A] bg-surface-nested p-[18px]"
      >
        <div
          class="pointer-events-none absolute -right-6 -top-8 h-[150px] w-[150px] rounded-full"
          style="background: radial-gradient(circle, rgba(46, 125, 242, 0.16), transparent 70%)"
        />
        <div class="relative">
          <div class="mb-[15px] flex items-center gap-3">
            <span
              class="flex h-11 w-11 flex-none items-center justify-center rounded-[13px] text-base font-extrabold text-white"
              :style="{ background: p.color }"
              >{{ p.initials }}</span
            >
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-[7px]">
                <span class="truncate text-[17px] font-bold">{{ p.name }}</span>
                <KBadge tone="brand">{{ t('exchangers.partner') }}</KBadge>
              </div>
              <div class="mt-[3px] text-[13px] text-ink-muted">
                <span class="text-warning-deep">★</span>
                <span class="tnum font-semibold text-ink">{{ p.rating }}</span> ·
                <span class="tnum">{{ p.reviews }}</span> {{ plural(p.reviews, 'reviews') }}
              </div>
            </div>
          </div>
          <div class="mb-[15px] flex rounded-xl border border-line bg-well px-1.5 py-3">
            <div class="flex-1 border-r border-line text-center">
              <div class="tnum text-[15px] font-bold">{{ p.reserve }}</div>
              <div class="mt-0.5 text-[11px] text-ink-faint">
                {{ t('exchangers.reserveLabel') }}
              </div>
            </div>
            <div class="flex-1 border-r border-line text-center">
              <div class="tnum text-[15px] font-bold">{{ p.directions }}</div>
              <div class="mt-0.5 text-[11px] text-ink-faint">
                {{ t('exchangers.directionsLabel') }}
              </div>
            </div>
            <div class="flex-1 text-center">
              <div class="tnum text-[15px] font-bold">{{ p.onSince ?? '—' }}</div>
              <div class="mt-0.5 text-[11px] text-ink-faint">{{ t('exchangers.onKursoStat') }}</div>
            </div>
          </div>
          <NuxtLink
            :to="p.href"
            :class="[openBtn, 'w-full bg-brand text-white hover:bg-brand-hover']"
            >{{ t('exchangers.openCard') }}</NuxtLink
          >
        </div>
      </div>
    </div>

    <!-- ===== Панель инструментов ===== -->
    <!-- десктоп -->
    <div class="mb-4 hidden flex-wrap items-center gap-3 md:flex">
      <label
        class="flex w-[270px] items-center gap-2.5 rounded-md border border-line bg-surface px-3.5 py-2.5 text-ink-faint focus-within:border-brand"
      >
        <svg
          class="flex-none"
          width="15"
          height="15"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
        >
          <circle cx="11" cy="11" r="7" />
          <path d="M21 21l-3.5-3.5" />
        </svg>
        <input
          v-model="search"
          type="text"
          :placeholder="t('exchangers.searchPlaceholder')"
          class="w-full bg-transparent text-sm text-ink placeholder:text-ink-faint focus:outline-none"
        />
      </label>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="f in filters"
          :key="f.key"
          type="button"
          class="rounded-md px-3.5 py-2 text-[13px] font-semibold transition-colors"
          :class="
            filter === f.key
              ? 'bg-brand text-white'
              : 'border border-line-strong bg-surface font-medium text-ink-muted hover:text-ink'
          "
          @click="setFilter(f.key)"
        >
          {{ t(f.i18n) }}
        </button>
      </div>
      <div class="ml-auto flex items-center gap-3.5">
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-md border border-line-strong bg-surface px-3 py-2 text-[13px] text-ink-muted transition-colors hover:text-ink"
          @click="cycleSort"
        >
          {{ t('exchangers.sort') }}: <span class="font-semibold text-ink">{{ sortLabel }}</span>
          <svg
            width="13"
            height="13"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#6E757E"
            stroke-width="2.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M6 9l6 6 6-6" />
          </svg>
        </button>
        <span class="whitespace-nowrap text-[13px] text-ink-faint"
          >1–<span class="tnum">{{ visible.length }}</span> {{ t('exchangers.of') }}
          <span class="tnum">{{ sorted.length }}</span></span
        >
      </div>
    </div>

    <!-- мобайл: чипы + сорт-бар -->
    <div class="md:hidden">
      <div class="scrollx -mx-4 mb-3 flex items-center gap-2 overflow-x-auto px-4">
        <button
          v-for="f in filters"
          :key="f.key"
          type="button"
          class="whitespace-nowrap rounded-[10px] px-3.5 py-2 text-[13px] font-semibold transition-colors"
          :class="
            filter === f.key
              ? 'bg-brand text-white'
              : 'border border-line-strong bg-surface font-medium text-ink-muted'
          "
          @click="setFilter(f.key)"
        >
          {{ t(f.i18n) }}
        </button>
      </div>
      <div class="mb-3 flex items-center justify-between">
        <span class="text-xs text-ink-faint"
          ><span class="tnum text-ink-muted">{{ sorted.length }}</span>
          {{ plural(sorted.length, 'services') }}</span
        >
        <button
          type="button"
          class="inline-flex items-center gap-1.5 text-[13px] text-ink-muted"
          @click="cycleSort"
        >
          {{ sortLabel }}
          <svg
            width="13"
            height="13"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#6E757E"
            stroke-width="2.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M6 9l6 6 6-6" />
          </svg>
        </button>
      </div>
    </div>

    <!-- ===== Пусто / загрузка ===== -->
    <div
      v-if="loading"
      class="rounded-2xl border border-line bg-surface py-20 text-center text-sm text-ink-faint"
    >
      {{ t('exchangers.loading') }}
    </div>
    <div
      v-else-if="sorted.length === 0"
      class="rounded-2xl border border-line bg-surface py-20 text-center text-sm text-ink-faint"
    >
      {{ t('exchangers.notFound') }}
    </div>

    <template v-else>
      <!-- ===== Таблица (десктоп) ===== -->
      <div class="hidden overflow-hidden rounded-2xl border border-line bg-surface md:block">
        <div
          class="grid grid-cols-[2.4fr_1.1fr_1.5fr_0.8fr_auto] gap-3.5 border-b border-line bg-well px-5 py-3 text-[11px] uppercase tracking-[0.05em] text-ink-faint"
        >
          <div>{{ t('exchangers.colExchanger') }}</div>
          <div>{{ t('exchangers.colRating') }}</div>
          <div>{{ t('exchangers.colReserve') }}</div>
          <div>{{ t('exchangers.colDirections') }}</div>
          <div />
        </div>
        <div
          v-for="r in rows"
          :key="r.slug"
          class="grid grid-cols-[2.4fr_1.1fr_1.5fr_0.8fr_auto] items-center gap-3.5 border-b border-line-subtle px-5 py-[15px] last:border-b-0"
          :class="r.partner ? 'bg-brand/[0.06]' : ''"
        >
          <div class="flex items-center gap-3">
            <span
              class="flex h-10 w-10 flex-none items-center justify-center rounded-xl text-[15px] font-extrabold text-white"
              :style="{ background: r.color }"
              >{{ r.initials }}</span
            >
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-[7px] text-[15px] font-semibold">
                {{ r.name }}
                <KBadge v-if="r.partner" tone="brand">{{ t('exchangers.partner') }}</KBadge>
                <KBadge v-if="r.verified" tone="success">
                  <svg
                    width="10"
                    height="10"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="3.2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M20 6 9 17l-5-5" />
                  </svg>
                  {{ t('exchangers.verified') }}
                </KBadge>
              </div>
              <div v-if="r.onSince" class="mt-[3px] text-[11px] text-ink-faint">
                {{ t('exchangers.onKursoSince') }} <span class="tnum">{{ r.onSince }}</span>
              </div>
            </div>
          </div>
          <div>
            <div class="text-sm">
              <span class="text-warning-deep">★</span>
              <span class="tnum font-semibold">{{ r.rating }}</span>
            </div>
            <div class="tnum mt-0.5 text-[11px] text-ink-faint">
              {{ r.reviews }} {{ plural(r.reviews, 'reviews') }}
            </div>
          </div>
          <div>
            <div class="tnum text-[15px] font-semibold">{{ r.reserve }}</div>
            <div v-if="r.assets.length" class="mt-0.5 truncate text-[11px] text-ink-faint">
              {{ r.assets.join(' · ') }}
            </div>
          </div>
          <div class="tnum text-[15px] font-semibold">{{ r.directions }}</div>
          <NuxtLink
            :to="r.href"
            :class="[
              openBtn,
              'px-[18px] py-[9px]',
              r.partner
                ? 'bg-brand text-white hover:bg-brand-hover'
                : 'border border-line-strong bg-surface-raised text-ink hover:border-[#3A4047]',
            ]"
            >{{ t('exchangers.open') }}</NuxtLink
          >
        </div>
      </div>

      <!-- ===== Компактный список (мобайл) ===== -->
      <div class="overflow-hidden rounded-[18px] border border-line bg-surface md:hidden">
        <NuxtLink
          v-for="r in rows"
          :key="r.slug"
          :to="r.href"
          class="flex items-center gap-[11px] border-b border-line-subtle px-[15px] py-[13px] last:border-b-0"
          :class="r.partner ? 'bg-brand/[0.06]' : ''"
        >
          <span
            class="flex h-[38px] w-[38px] flex-none items-center justify-center rounded-[11px] text-[13px] font-extrabold text-white"
            :style="{ background: r.color }"
            >{{ r.initials }}</span
          >
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-1.5">
              <span class="truncate text-[15px] font-semibold">{{ r.name }}</span>
              <KBadge v-if="r.partner" tone="brand">{{ t('exchangers.partner') }}</KBadge>
            </div>
            <div class="mt-0.5 truncate text-xs text-ink-faint">
              <span class="text-warning-deep">★</span>
              <span class="tnum text-ink-muted">{{ r.rating }}</span> ·
              <span class="tnum">{{ r.reviews }}</span> ·
              <span class="tnum">{{ r.reserve }}</span>
            </div>
          </div>
          <svg
            class="flex-none text-ink-ghost"
            width="17"
            height="17"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M9 6l6 6-6 6" />
          </svg>
        </NuxtLink>
      </div>

      <!-- ===== Пагинация ===== -->
      <div v-if="canLoadMore" class="mt-6 flex justify-center">
        <button
          type="button"
          class="rounded-xl border border-line-strong bg-surface px-6 py-3 text-sm font-semibold text-ink transition-colors hover:border-[#3A4047]"
          @click="shown += PAGE"
        >
          {{ t('exchangers.showMore') }} <span class="tnum text-ink-faint">{{ remaining }}</span>
        </button>
      </div>
    </template>
  </div>
</template>
