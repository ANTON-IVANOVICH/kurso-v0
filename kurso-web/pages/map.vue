<script setup lang="ts">
import { computed, ref } from 'vue'

// Immersive full-screen map (design variant 2). Scaffold: the map surface is a
// styled placeholder with a `mapContainer` ref where Mapbox GL will mount later;
// pins/markers are absolutely positioned for now (percentages stand in for
// lng/lat) and become real Mapbox markers once integrated.
definePageMeta({ layout: false })

useSeoMeta({
  title: 'Карта обменников — Kurso',
  description:
    'Карта офлайн-обменников с наличными: точки на карте, курсы, часы работы и маршруты до ближайших пунктов.',
})

// Mapbox GL will mount into this element.
const mapContainer = ref<HTMLElement | null>(null)

interface Point {
  id: string
  name: string
  initials: string
  color: string
  rate: string
  address: string
  distance: string
  rating: string
  reviews: string
  hours: string
  open: boolean
  partner?: boolean
  best?: boolean
  x: number // left % — placeholder for lng
  y: number // top % — placeholder for lat
}

const points: Point[] = [
  {
    id: 'cb',
    name: 'CryptoBridge',
    initials: 'CB',
    color: '#2E7DF2',
    rate: '81.20',
    address: 'Тверская, 12',
    distance: '0.4 км',
    rating: '4.9',
    reviews: '1203',
    hours: '10:00–21:00',
    open: true,
    partner: true,
    best: true,
    x: 40,
    y: 38,
  },
  {
    id: 'n',
    name: 'NetEx24',
    initials: 'N',
    color: '#5B3FA0',
    rate: '80.95',
    address: 'Кутузовский, 5',
    distance: '1.1 км',
    rating: '4.8',
    reviews: '842',
    hours: '24/7',
    open: true,
    x: 61,
    y: 31,
  },
  {
    id: 'ot',
    name: 'ОбменТочка',
    initials: 'ОТ',
    color: '#1F8A5B',
    rate: '80.70',
    address: 'Арбат, 24',
    distance: '1.8 км',
    rating: '4.6',
    reviews: '318',
    hours: '09:00–20:00',
    open: true,
    x: 50,
    y: 62,
  },
  {
    id: 'cp',
    name: 'CashPoint',
    initials: 'CP',
    color: '#3A414A',
    rate: '80.40',
    address: 'Ленинский, 40',
    distance: '2.3 км',
    rating: '4.5',
    reviews: '210',
    hours: 'закрыто',
    open: false,
    x: 72,
    y: 67,
  },
]

const selectedId = ref('cb')
const selected = computed(() => points.find((p) => p.id === selectedId.value) ?? points[0])
// keep the desktop callout from overflowing the map edges
const calloutLeft = computed(() => Math.min(Math.max(selected.value.x, 24), 76))

const view = ref<'map' | 'list'>('map')
</script>

<template>
  <div class="relative flex h-[100dvh] flex-col overflow-hidden bg-canvas">
    <!-- desktop: standard app header above the map -->
    <div class="hidden md:block">
      <SiteHeader />
    </div>

    <!-- immersive map surface -->
    <div class="relative flex-1 overflow-hidden bg-[#0E1316]">
      <!-- Mapbox GL mounts here; decorative texture stands in until integrated -->
      <div ref="mapContainer" class="absolute inset-0">
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
        <div
          class="absolute right-40 top-[70px] h-40 w-[220px] rounded-[50%_50%_44%_56%] bg-success/[0.08]"
        />
      </div>

      <!-- your location -->
      <div class="absolute left-[46%] top-1/2 z-[2] -translate-x-1/2 -translate-y-1/2">
        <div
          class="absolute left-1/2 top-1/2 h-[18px] w-[18px] -translate-x-1/2 -translate-y-1/2 rounded-full bg-brand-light animate-kping"
        />
        <div
          class="relative h-4 w-4 rounded-full border-[3px] border-white bg-brand-light shadow-[0_0_0_2px_rgba(74,144,245,0.4)]"
        />
      </div>

      <!-- pins -->
      <button
        v-for="p in points"
        :key="p.id"
        type="button"
        class="absolute z-[4] flex -translate-x-1/2 -translate-y-full flex-col items-center focus:outline-none"
        :class="p.id === selectedId ? 'z-[6]' : ''"
        :style="{ left: `${p.x}%`, top: `${p.y}%` }"
        @click="selectedId = p.id"
      >
        <!-- best pin: brand pill with ping -->
        <template v-if="p.best">
          <div
            class="absolute top-2 h-[30px] w-[30px] rounded-full bg-brand animate-kping"
            aria-hidden="true"
          />
          <div
            class="relative inline-flex items-center gap-1.5 rounded-full border-[3px] border-[#0E1316] bg-brand py-1.5 pl-1.5 pr-3 shadow-[0_8px_20px_rgba(46,125,242,0.5)]"
          >
            <span
              class="flex h-[26px] w-[26px] items-center justify-center rounded-full bg-white/20 text-[11px] font-bold text-white"
              >{{ p.initials }}</span
            >
            <span class="tnum text-[13px] font-bold text-white">{{ p.rate }}</span>
          </div>
          <div
            class="-mt-0.5 h-0 w-0 border-l-[6px] border-r-[6px] border-t-[9px] border-l-transparent border-r-transparent border-t-brand"
          />
        </template>
        <!-- normal pin: dark price pill -->
        <template v-else>
          <div
            class="inline-flex items-center gap-1.5 rounded-full border-2 border-[#0E1316] bg-surface-nav py-1 pl-1 pr-2.5 shadow-[0_6px_16px_rgba(0,0,0,0.5)]"
            :class="p.id === selectedId ? '!border-brand/60' : ''"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full text-[9px] font-bold text-white"
              :style="{ backgroundColor: p.color }"
              >{{ p.initials }}</span
            >
            <span
              class="tnum text-xs font-semibold"
              :class="p.open ? 'text-ink' : 'text-ink-faint'"
              >{{ p.rate }}</span
            >
          </div>
          <div class="mt-px h-[5px] w-[5px] rounded-full bg-surface-nav" />
        </template>
      </button>

      <!-- selected callout (desktop) -->
      <div
        class="absolute z-[9] hidden w-[300px] -translate-x-1/2 translate-y-4 md:block"
        :style="{ left: `${calloutLeft}%`, top: `${selected.y}%` }"
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
                {{ selected.address }} · <span class="text-ink-faint">{{ selected.distance }}</span>
              </div>
            </div>
            <div class="flex-none text-right">
              <div class="tnum text-xl font-bold text-ink">{{ selected.rate }} ₽</div>
              <div class="mt-0.5 text-xs text-success-bright">за 1 USDT</div>
            </div>
          </div>
          <div class="mb-3.5 flex items-center gap-2.5 text-[13px]">
            <span class="text-ink-muted"
              ><span class="text-warning">★</span> <span class="tnum">{{ selected.rating }}</span> ·
              <span class="tnum">{{ selected.reviews }}</span></span
            >
            <span
              class="rounded-md border px-2 py-0.5"
              :class="
                selected.open
                  ? 'border-success/30 bg-success/10 text-success-bright'
                  : 'border-line-strong bg-surface-chip text-ink-faint'
              "
              >{{ selected.hours }}</span
            >
          </div>
          <div class="flex gap-2.5">
            <button
              type="button"
              class="inline-flex flex-1 items-center justify-center gap-1.5 rounded-xl border border-line-strong bg-surface-raised py-2.5 text-sm font-semibold text-ink transition-colors hover:border-[#3A4047]"
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
            <KButton class="flex-1">Перейти</KButton>
          </div>
        </div>
        <div
          class="mx-auto -mt-px h-0 w-0 border-l-[7px] border-r-[7px] border-t-[10px] border-l-transparent border-r-transparent border-t-[rgba(18,21,24,0.92)]"
        />
      </div>

      <!-- floating search / filter panel (desktop) -->
      <div
        class="absolute left-5 top-5 z-10 hidden w-[300px] rounded-2xl border border-line bg-[rgba(14,19,22,0.9)] p-3.5 shadow-pop backdrop-blur-md md:block"
      >
        <div
          class="mb-2.5 flex items-center gap-2.5 rounded-xl border border-line-strong bg-surface px-3 py-2.5 text-ink-faint"
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
          <span class="text-sm">Поиск по адресу или названию</span>
        </div>
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
          <span class="tnum font-semibold text-ink">18</span> точек
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

      <!-- zoom + locate controls (desktop) -->
      <div class="absolute right-5 top-1/2 z-10 hidden -translate-y-1/2 flex-col gap-2.5 md:flex">
        <div
          class="flex flex-col overflow-hidden rounded-xl border border-line bg-[rgba(14,19,22,0.9)] shadow-pop backdrop-blur"
        >
          <button
            type="button"
            aria-label="Приблизить"
            class="flex h-[42px] w-[42px] items-center justify-center border-b border-line text-ink transition-colors hover:bg-white/[0.04]"
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
        <div class="min-w-0 flex-1">
          <div class="text-sm font-semibold text-ink">Москва</div>
          <div class="text-xs text-ink-faint">USDT → RUB · от 10 000 ₽</div>
        </div>
        <span
          class="flex h-9 w-9 flex-none items-center justify-center rounded-xl border border-line-strong bg-surface text-ink-muted"
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
            <path d="M4 6h16M7 12h10M10 18h4" />
          </svg>
        </span>
      </div>

      <!-- mobile: counter pill -->
      <div
        class="absolute left-1/2 top-[86px] z-[9] inline-flex -translate-x-1/2 items-center gap-2 rounded-full border border-line bg-[rgba(14,19,22,0.92)] px-3.5 py-2 text-[13px] text-ink-muted shadow-toast backdrop-blur md:hidden"
      >
        <KStatusDot tone="success" pulse :size="7" />Найдено
        <span class="tnum font-semibold text-ink">18</span> точек
      </div>

      <!-- mobile: locate button -->
      <button
        type="button"
        aria-label="Моё местоположение"
        class="absolute right-4 bottom-[168px] z-[9] flex h-11 w-11 items-center justify-center rounded-full border border-line bg-[rgba(14,19,22,0.92)] text-brand-bright shadow-pop backdrop-blur md:hidden"
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

      <!-- bottom horizontal carousel of nearby points -->
      <div
        class="scrollx absolute inset-x-0 bottom-4 z-[11] flex gap-3 overflow-x-auto px-4 md:inset-x-5 md:px-0.5"
      >
        <button
          v-for="p in points"
          :key="p.id"
          type="button"
          class="w-[300px] flex-none rounded-2xl border bg-[rgba(18,21,24,0.94)] p-3.5 text-left shadow-toast backdrop-blur-md transition-colors"
          :class="[
            p.id === selectedId ? 'border-brand/45' : 'border-line hover:border-line-strong',
            p.open ? '' : 'opacity-70',
          ]"
          @click="selectedId = p.id"
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
            <div class="tnum text-base font-bold" :class="p.open ? 'text-ink' : 'text-ink-faint'">
              {{ p.rate }} ₽
            </div>
          </div>
          <div class="mb-2.5 text-[13px] text-ink-muted">
            {{ p.address }} · <span class="text-ink-faint">{{ p.distance }}</span> ·
            <span :class="p.open ? 'text-success-bright' : 'text-ink-ghost'">{{
              p.open ? 'открыто' : 'закрыто'
            }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-xs text-ink-muted"
              ><span class="text-warning">★</span> <span class="tnum">{{ p.rating }}</span></span
            >
            <span
              class="inline-flex items-center gap-1.5 text-[13px] font-semibold text-brand-bright"
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
  </div>
</template>
