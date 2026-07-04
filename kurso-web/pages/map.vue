<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { MapPoint } from '~/composables/useMapPoints'

// Immersive full-screen map (design variant 2). Real exchanger cash-desks come
// from the API with their live rate for the direction; a configured Mapbox token
// renders an interactive dark map with markers + popups, otherwise the page falls
// back to a static projected preview. The mobile bottom nav stays visible here.
definePageMeta({ layout: false })
useSeoMeta({
  title: 'Карта обменников — Kurso',
  description:
    'Карта офлайн-обменников с наличными: точки на карте, курсы, часы работы и маршруты до ближайших пунктов.',
})

const DIRECTION = 'usdt-tinkoff'
const { points } = useMapPoints(DIRECTION)
const base = useApiBase()
const config = useRuntimeConfig()
const token = (config.public.mapboxToken as string) || ''
const hasToken = !!token

const mapContainer = ref<HTMLElement | null>(null)
const selectedSlug = ref<string | null>(null)
const view = ref<'map' | 'list'>('map')
const userLoc = ref<{ lat: number; lng: number } | null>(null)
const search = ref('')

const selected = computed(
  () => points.value.find((p) => p.slug === selectedSlug.value) ?? points.value[0] ?? null,
)
// Default the selection to the best-rate point once data arrives.
watch(points, (list) => {
  if (!selectedSlug.value && list.length) selectedSlug.value = list[0].slug
})

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return points.value
  return points.value.filter(
    (p) => p.name.toLowerCase().includes(q) || (p.address ?? '').toLowerCase().includes(q),
  )
})

function fmtRate(rate: string | null | undefined): string {
  if (!rate) return '— ₽'
  const n = parseFloat(rate)
  return Number.isFinite(n) ? `${fmtNumber(n, 2)} ₽` : '— ₽'
}

// Haversine distance user→point, as a short "X км" label (blank without a fix).
function distanceKm(p: MapPoint): string {
  if (!userLoc.value) return ''
  const R = 6371
  const dLat = ((p.lat - userLoc.value.lat) * Math.PI) / 180
  const dLng = ((p.lng - userLoc.value.lng) * Math.PI) / 180
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos((userLoc.value.lat * Math.PI) / 180) *
      Math.cos((p.lat * Math.PI) / 180) *
      Math.sin(dLng / 2) ** 2
  const km = R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
  return km < 10 ? `${km.toFixed(1)} км` : `${Math.round(km)} км`
}

const isOpen = (p: MapPoint) => !!p.hours && p.hours !== 'закрыто'

function go(p: MapPoint) {
  // Real clickout + referral redirect through the API.
  window.location.href = `${base}/go/${p.slug}?direction=${DIRECTION}`
}
function route(p: MapPoint) {
  // Open external turn-by-turn directions (Yandex Maps — RU audience).
  const dest = `${p.lat},${p.lng}`
  const from = userLoc.value ? `${userLoc.value.lat},${userLoc.value.lng}~` : ''
  window.open(`https://yandex.ru/maps/?rtext=${from}${dest}&rtt=auto`, '_blank', 'noopener')
}

// --- static fallback projection (no Mapbox token): place real points inside the
// container by their lng/lat within a padded bbox, so relative positions hold ---
const projection = computed(() => {
  const list = points.value
  if (!list.length) return new Map<string, { x: number; y: number }>()
  const lats = list.map((p) => p.lat)
  const lngs = list.map((p) => p.lng)
  let minLat = Math.min(...lats),
    maxLat = Math.max(...lats),
    minLng = Math.min(...lngs),
    maxLng = Math.max(...lngs)
  const padLat = (maxLat - minLat) * 0.35 + 0.004
  const padLng = (maxLng - minLng) * 0.35 + 0.004
  minLat -= padLat
  maxLat += padLat
  minLng -= padLng
  maxLng += padLng
  const m = new Map<string, { x: number; y: number }>()
  for (const p of list) {
    m.set(p.slug, {
      x: ((p.lng - minLng) / (maxLng - minLng)) * 100,
      y: (1 - (p.lat - minLat) / (maxLat - minLat)) * 100,
    })
  }
  return m
})
const pos = (slug: string) => projection.value.get(slug) ?? { x: 50, y: 50 }
const calloutLeft = computed(() => Math.min(Math.max(pos(selected.value?.slug ?? '').x, 24), 76))

// --- Mapbox (real map) ---
let mb: ReturnType<typeof useMapboxMap> | null = null
onMounted(() => {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (p) => (userLoc.value = { lat: p.coords.latitude, lng: p.coords.longitude }),
      () => {},
      { timeout: 6000 },
    )
  }
  if (hasToken) {
    mb = useMapboxMap({
      container: mapContainer,
      token,
      points,
      selected: selectedSlug,
      userLoc,
      fmtRate,
      onSelect: (slug) => (selectedSlug.value = slug),
      onRoute: route,
      onGo: go,
    })
  }
})

function select(slug: string) {
  selectedSlug.value = slug
}
function zoomIn() {
  mb?.zoomIn()
}
function zoomOut() {
  mb?.zoomOut()
}
function locate() {
  if (!userLoc.value && navigator.geolocation) {
    navigator.geolocation.getCurrentPosition((p) => {
      userLoc.value = { lat: p.coords.latitude, lng: p.coords.longitude }
      mb?.locate()
    })
  } else {
    mb?.locate()
  }
}
</script>

<template>
  <div class="relative flex h-[100dvh] flex-col overflow-hidden bg-canvas">
    <!-- desktop: standard app header above the map -->
    <div class="hidden md:block">
      <SiteHeader />
    </div>

    <div class="relative flex-1 overflow-hidden bg-[#0E1316]">
      <!-- Mapbox GL mounts here when a token is configured -->
      <div v-if="hasToken" ref="mapContainer" class="absolute inset-0" />

      <!-- fallback: decorative surface + projected real pins (no token) -->
      <template v-else>
        <div class="absolute inset-0">
          <div
            class="absolute inset-0 bg-[repeating-linear-gradient(0deg,transparent_0_47px,rgba(255,255,255,0.022)_47px_49px),repeating-linear-gradient(90deg,transparent_0_47px,rgba(255,255,255,0.022)_47px_49px)]"
          />
          <div
            class="absolute left-[34%] top-[-40px] h-[820px] w-[90px] rotate-[24deg] bg-white/[0.03]"
          />
          <div
            class="absolute left-[-40px] top-[44%] h-16 w-[1400px] rotate-[-7deg] bg-white/[0.03]"
          />
          <div
            class="absolute bottom-[-90px] right-[-50px] h-[320px] w-[460px] rotate-[8deg] rounded-[46%_54%_60%_40%] bg-brand/10"
          />
        </div>

        <!-- your location (centred; approximate without a map) -->
        <div class="absolute left-[46%] top-1/2 z-[2] -translate-x-1/2 -translate-y-1/2">
          <div
            class="absolute left-1/2 top-1/2 h-[18px] w-[18px] -translate-x-1/2 -translate-y-1/2 animate-kping rounded-full bg-brand-light"
          />
          <div
            class="relative h-4 w-4 rounded-full border-[3px] border-white bg-brand-light shadow-[0_0_0_2px_rgba(74,144,245,0.4)]"
          />
        </div>

        <!-- pins -->
        <button
          v-for="p in filtered"
          :key="p.slug"
          type="button"
          class="absolute z-[4] flex -translate-x-1/2 -translate-y-full flex-col items-center focus:outline-none"
          :class="p.slug === selectedSlug ? 'z-[6]' : ''"
          :style="{ left: `${pos(p.slug).x}%`, top: `${pos(p.slug).y}%` }"
          @click="select(p.slug)"
        >
          <template v-if="p.partner">
            <div
              class="absolute top-2 h-[30px] w-[30px] animate-kping rounded-full bg-brand"
              aria-hidden="true"
            />
            <div
              class="relative inline-flex items-center gap-1.5 rounded-full border-[3px] border-[#0E1316] bg-brand py-1.5 pl-1.5 pr-3 shadow-[0_8px_20px_rgba(46,125,242,0.5)]"
            >
              <span
                class="flex h-[26px] w-[26px] items-center justify-center rounded-full bg-white/20 text-[11px] font-bold text-white"
                >{{ p.name.slice(0, 2).toUpperCase() }}</span
              >
              <span class="tnum text-[13px] font-bold text-white">{{
                fmtRate(p.rate).replace(' ₽', '')
              }}</span>
            </div>
            <div
              class="-mt-0.5 h-0 w-0 border-l-[6px] border-r-[6px] border-t-[9px] border-l-transparent border-r-transparent border-t-brand"
            />
          </template>
          <template v-else>
            <div
              class="inline-flex items-center gap-1.5 rounded-full border-2 border-[#0E1316] bg-surface-nav py-1 pl-1 pr-2.5 shadow-[0_6px_16px_rgba(0,0,0,0.5)]"
              :class="p.slug === selectedSlug ? '!border-brand/60' : ''"
            >
              <span
                class="flex h-5 w-5 items-center justify-center rounded-full bg-[#3A4452] text-[9px] font-bold text-white"
                >{{ p.name.slice(0, 2).toUpperCase() }}</span
              >
              <span
                class="tnum text-xs font-semibold"
                :class="isOpen(p) ? 'text-ink' : 'text-ink-faint'"
                >{{ fmtRate(p.rate).replace(' ₽', '') }}</span
              >
            </div>
            <div class="mt-px h-[5px] w-[5px] rounded-full bg-surface-nav" />
          </template>
        </button>

        <!-- selected callout (desktop, fallback only) -->
        <div
          v-if="selected"
          class="absolute z-[9] hidden w-[300px] -translate-x-1/2 translate-y-4 md:block"
          :style="{ left: `${calloutLeft}%`, top: `${pos(selected.slug).y}%` }"
        >
          <div
            class="rounded-[18px] border border-[#2A3138] bg-[rgba(18,21,24,0.92)] p-4 shadow-modal backdrop-blur-md"
          >
            <div class="mb-2 flex items-start justify-between gap-3">
              <div class="min-w-0">
                <div class="flex items-center gap-2 text-[17px] font-bold text-ink">
                  {{ selected.name }}
                  <span
                    v-if="selected.partner"
                    class="flex-none rounded-md bg-brand/15 px-1.5 py-0.5 text-[11px] font-semibold text-brand-bright"
                    >Партнёр</span
                  >
                </div>
                <div class="mt-1.5 text-[13px] text-ink-muted">
                  {{ selected.address
                  }}<span v-if="distanceKm(selected)">
                    · <span class="text-ink-faint">{{ distanceKm(selected) }}</span></span
                  >
                </div>
              </div>
              <div class="flex-none text-right">
                <div class="tnum text-xl font-bold text-ink">{{ fmtRate(selected.rate) }}</div>
                <div class="mt-0.5 text-xs text-success-bright">за 1 USDT</div>
              </div>
            </div>
            <div class="mb-3.5 flex items-center gap-2.5 text-[13px]">
              <span class="text-ink-muted"
                ><span class="text-warning">★</span>
                <span class="tnum">{{ selected.ratingAvg?.toFixed(1) ?? '—' }}</span> ·
                <span class="tnum">{{ selected.reviewsCount }}</span></span
              >
              <span
                class="rounded-md border px-2 py-0.5"
                :class="
                  isOpen(selected)
                    ? 'border-success/30 bg-success/10 text-success-bright'
                    : 'border-line-strong bg-surface-chip text-ink-faint'
                "
                >{{ selected.hours ?? '—' }}</span
              >
            </div>
            <div class="flex gap-2.5">
              <button
                type="button"
                class="inline-flex flex-1 items-center justify-center gap-1.5 rounded-xl border border-line-strong bg-surface-raised py-2.5 text-sm font-semibold text-ink transition-colors hover:border-[#3A4047]"
                @click="route(selected)"
              >
                <svg
                  width="15"
                  height="15"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="#6BA6FF"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M3 11l18-8-8 18-2-8-8-2Z" />
                </svg>
                Маршрут
              </button>
              <KButton class="flex-1" @click="go(selected)">Перейти</KButton>
            </div>
          </div>
          <div
            class="mx-auto -mt-px h-0 w-0 border-l-[7px] border-r-[7px] border-t-[10px] border-l-transparent border-r-transparent border-t-[rgba(18,21,24,0.92)]"
          />
        </div>
      </template>

      <!-- floating search / filter panel (desktop) -->
      <div
        class="absolute left-5 top-5 z-10 hidden w-[300px] rounded-2xl border border-line bg-[rgba(14,19,22,0.9)] p-3.5 shadow-pop backdrop-blur-md md:block"
      >
        <label
          class="mb-2.5 flex items-center gap-2.5 rounded-xl border border-line-strong bg-surface px-3 py-2.5 text-ink-faint focus-within:border-brand"
        >
          <svg
            class="flex-none"
            width="16"
            height="16"
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
            placeholder="Поиск по адресу или названию"
            class="w-full bg-transparent text-sm text-ink placeholder:text-ink-faint focus:outline-none"
          />
        </label>
        <div class="flex flex-wrap gap-2">
          <span
            class="inline-flex items-center gap-1.5 rounded-lg border border-line-strong bg-surface px-3 py-2 text-[13px] font-medium text-ink"
          >
            <svg
              width="13"
              height="13"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#2E7DF2"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 21s-7-5.5-7-11a7 7 0 0 1 14 0c0 5.5-7 11-7 11Z" />
              <circle cx="12" cy="10" r="2.5" />
            </svg>
            Москва
          </span>
          <span
            class="inline-flex items-center gap-1.5 rounded-lg border border-line-strong bg-surface py-2 pl-1.5 pr-3 text-[13px] text-ink"
          >
            <span
              class="flex h-[18px] w-[18px] items-center justify-center rounded-full bg-[#26A17B] text-[10px] font-bold text-white"
              >₮</span
            >
            USDT
            <svg
              width="11"
              height="11"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#6E757E"
              stroke-width="2.2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M5 12h14M13 6l6 6-6 6" />
            </svg>
            <span
              class="flex h-[18px] w-[18px] items-center justify-center rounded-full bg-[#3A4452] text-[10px] font-bold text-white"
              >₽</span
            >
            RUB
          </span>
        </div>
        <div
          class="mt-3 flex items-center gap-2 border-t border-line pt-2.5 text-[13px] text-ink-faint"
        >
          <KStatusDot tone="success" pulse :size="8" />Найдено
          <span class="tnum font-semibold text-ink">{{ filtered.length }}</span> точек
        </div>
      </div>

      <!-- view toggle (desktop) -->
      <div
        class="absolute right-5 top-5 z-10 hidden rounded-xl border border-line bg-[rgba(14,19,22,0.9)] p-1 shadow-pop backdrop-blur-md md:flex"
      >
        <button
          type="button"
          class="inline-flex items-center gap-1.5 rounded-lg px-3.5 py-2 text-[13px] font-semibold transition-colors"
          :class="view === 'map' ? 'bg-brand text-white' : 'text-ink-muted hover:text-ink'"
          @click="view = 'map'"
        >
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M9 4 3 6.5V20l6-2.5 6 2.5 6-2.5V4l-6 2.5Z" />
            <path d="M9 4v13.5" />
            <path d="M15 6.5V20" />
          </svg>
          Карта
        </button>
        <button
          type="button"
          class="inline-flex items-center gap-1.5 rounded-lg px-3.5 py-2 text-[13px] font-semibold transition-colors"
          :class="view === 'list' ? 'bg-brand text-white' : 'text-ink-muted hover:text-ink'"
          @click="view = 'list'"
        >
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M8 6h13M8 12h13M8 18h13" />
            <path d="M3 6h.01M3 12h.01M3 18h.01" />
          </svg>
          Список
        </button>
      </div>

      <!-- list overlay (desktop, when toggled) -->
      <div
        v-if="view === 'list'"
        class="absolute right-5 top-[74px] z-[8] hidden max-h-[calc(100%-220px)] w-[340px] flex-col gap-2 overflow-y-auto rounded-2xl border border-line bg-[rgba(14,19,22,0.94)] p-3 shadow-pop backdrop-blur-md md:flex"
      >
        <button
          v-for="p in filtered"
          :key="p.slug"
          type="button"
          class="rounded-xl border p-3 text-left transition-colors"
          :class="
            p.slug === selectedSlug
              ? 'border-brand/45 bg-brand/[0.06]'
              : 'border-line hover:border-line-strong'
          "
          @click="select(p.slug)"
        >
          <div class="flex items-center justify-between gap-2">
            <span class="text-sm font-semibold text-ink">{{ p.name }}</span>
            <span class="tnum text-sm font-bold text-ink">{{ fmtRate(p.rate) }}</span>
          </div>
          <div class="mt-1 text-xs text-ink-muted">{{ p.address }} · {{ p.hours }}</div>
        </button>
      </div>

      <!-- zoom + locate controls (desktop) -->
      <div class="absolute right-5 top-1/2 z-10 hidden -translate-y-1/2 flex-col gap-2.5 md:flex">
        <div
          class="flex flex-col overflow-hidden rounded-xl border border-line bg-[rgba(14,19,22,0.9)] shadow-pop backdrop-blur"
        >
          <button
            type="button"
            aria-label="Приблизить"
            class="flex h-[42px] w-[42px] items-center justify-center border-b border-line text-ink transition-colors hover:bg-white/[0.04]"
            @click="zoomIn"
          >
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.2"
              stroke-linecap="round"
            >
              <path d="M12 5v14M5 12h14" />
            </svg>
          </button>
          <button
            type="button"
            aria-label="Отдалить"
            class="flex h-[42px] w-[42px] items-center justify-center text-ink transition-colors hover:bg-white/[0.04]"
            @click="zoomOut"
          >
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.2"
              stroke-linecap="round"
            >
              <path d="M5 12h14" />
            </svg>
          </button>
        </div>
        <button
          type="button"
          aria-label="Моё местоположение"
          class="flex h-[42px] w-[42px] items-center justify-center rounded-xl border border-line bg-[rgba(14,19,22,0.9)] text-brand-bright shadow-pop backdrop-blur transition-colors hover:bg-white/[0.04]"
          @click="locate"
        >
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <circle cx="12" cy="12" r="3.5" />
            <path d="M12 2v3M12 19v3M2 12h3M19 12h3" />
          </svg>
        </button>
      </div>

      <!-- mobile: floating search pill -->
      <div
        class="absolute inset-x-4 top-4 z-10 flex items-center gap-2.5 rounded-2xl border border-line bg-[rgba(14,19,22,0.92)] px-3.5 py-3 shadow-pop backdrop-blur-md md:hidden"
      >
        <svg
          class="flex-none text-ink-faint"
          width="17"
          height="17"
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
          placeholder="Москва · USDT → RUB"
          class="min-w-0 flex-1 bg-transparent text-sm font-semibold text-ink placeholder:text-ink-muted focus:outline-none"
        />
        <span
          class="tnum flex-none rounded-lg border border-line-strong bg-surface px-2 py-1 text-xs text-ink-muted"
          >{{ filtered.length }}</span
        >
      </div>

      <!-- mobile: locate button -->
      <button
        type="button"
        aria-label="Моё местоположение"
        class="absolute bottom-[236px] right-4 z-[9] flex h-11 w-11 items-center justify-center rounded-full border border-line bg-[rgba(14,19,22,0.92)] text-brand-bright shadow-pop backdrop-blur md:hidden"
        @click="locate"
      >
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="12" cy="12" r="3.5" />
          <path d="M12 2v3M12 19v3M2 12h3M19 12h3" />
        </svg>
      </button>

      <!-- bottom horizontal carousel (kept clear of the mobile bottom nav) -->
      <div
        class="scrollx absolute inset-x-0 bottom-[96px] z-[11] flex gap-3 overflow-x-auto px-4 md:inset-x-5 md:bottom-5 md:px-0.5"
      >
        <button
          v-for="p in filtered"
          :key="p.slug"
          type="button"
          class="w-[300px] flex-none rounded-2xl border bg-[rgba(18,21,24,0.94)] p-3.5 text-left shadow-toast backdrop-blur-md transition-colors"
          :class="[
            p.slug === selectedSlug ? 'border-brand/45' : 'border-line hover:border-line-strong',
            isOpen(p) ? '' : 'opacity-70',
          ]"
          @click="select(p.slug)"
        >
          <div class="mb-1.5 flex items-center justify-between gap-2">
            <div class="flex items-center gap-2 text-[15px] font-semibold text-ink">
              {{ p.name }}
              <span
                v-if="p.partner"
                class="rounded-[5px] bg-brand/[0.18] px-1.5 py-0.5 text-[10px] font-semibold text-brand-bright"
                >Партнёр</span
              >
            </div>
            <div
              class="tnum text-base font-bold"
              :class="isOpen(p) ? 'text-ink' : 'text-ink-faint'"
            >
              {{ fmtRate(p.rate) }}
            </div>
          </div>
          <div class="mb-2.5 text-[13px] text-ink-muted">
            {{ p.address
            }}<span v-if="distanceKm(p)">
              · <span class="text-ink-faint">{{ distanceKm(p) }}</span></span
            >
            ·
            <span :class="isOpen(p) ? 'text-success-bright' : 'text-ink-ghost'">{{
              isOpen(p) ? 'открыто' : 'закрыто'
            }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-ink-muted"
              ><span class="text-warning">★</span>
              <span class="tnum">{{ p.ratingAvg?.toFixed(1) ?? '—' }}</span></span
            >
            <span
              class="inline-flex items-center gap-1.5 text-[13px] font-semibold text-brand-bright"
              @click.stop="route(p)"
            >
              <svg
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M3 11l18-8-8 18-2-8-8-2Z" />
              </svg>
              Маршрут
            </span>
          </div>
        </button>
      </div>
    </div>

    <!-- mobile bottom navigation stays on the map page -->
    <MobileBottomNav />
  </div>
</template>

<style>
/* Global (unscoped) so Mapbox's imperatively-created marker/popup DOM is styled.
   Loading the CSS here keeps it off every other page. */
@import 'mapbox-gl/dist/mapbox-gl.css';

.kmap-pin {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px 4px 4px;
  border: 2px solid #0e1316;
  border-radius: 999px;
  background: #14181d;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.5);
  cursor: pointer;
  font-family: inherit;
}
.kmap-pin__badge {
  display: flex;
  height: 20px;
  width: 20px;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: #3a4452;
  color: #fff;
  font-size: 9px;
  font-weight: 700;
}
.kmap-pin__rate {
  font-family: 'Geist Mono', ui-monospace, monospace;
  font-size: 12px;
  font-weight: 600;
  color: #f4f5f7;
}
.kmap-pin--partner {
  background: #2e7df2;
  border-color: #0e1316;
  box-shadow: 0 8px 20px rgba(46, 125, 242, 0.5);
}
.kmap-pin--partner .kmap-pin__badge {
  background: rgba(255, 255, 255, 0.2);
}
.kmap-pin--partner .kmap-pin__rate {
  color: #fff;
}
.kmap-pin--active {
  outline: 2px solid rgba(46, 125, 242, 0.7);
  outline-offset: 1px;
}
.kmap-user {
  height: 16px;
  width: 16px;
  border-radius: 999px;
  border: 3px solid #fff;
  background: #4a90f5;
  box-shadow: 0 0 0 2px rgba(74, 144, 245, 0.4);
}
.kmap-popup .mapboxgl-popup-content,
.kmap-popup-wrap .mapboxgl-popup-content {
  background: transparent;
  padding: 0;
  box-shadow: none;
}
.kmap-popup {
  width: 260px;
  border-radius: 16px;
  border: 1px solid #2a3138;
  background: rgba(18, 21, 24, 0.95);
  padding: 14px;
  backdrop-filter: blur(8px);
  color: #f4f5f7;
}
.kmap-popup__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 6px;
}
.kmap-popup__name {
  font-size: 15px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 6px;
}
.kmap-popup__tag {
  border-radius: 6px;
  background: rgba(46, 125, 242, 0.15);
  padding: 1px 6px;
  font-size: 10px;
  font-weight: 600;
  color: #6ba6ff;
}
.kmap-popup__rate {
  font-family: 'Geist Mono', ui-monospace, monospace;
  font-size: 16px;
  font-weight: 700;
}
.kmap-popup__meta {
  font-size: 12px;
  color: #a8aeb6;
  margin-top: 2px;
}
.kmap-popup__actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}
.kmap-btn {
  flex: 1;
  border-radius: 10px;
  padding: 8px 0;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
}
.kmap-btn--ghost {
  border: 1px solid #2e343b;
  background: #1b1e22;
  color: #f4f5f7;
}
.kmap-btn--primary {
  border: none;
  background: #2e7df2;
  color: #fff;
}
.mapboxgl-popup-tip {
  border-top-color: rgba(18, 21, 24, 0.95) !important;
}
</style>
