<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

// Widget builder: pick a widget type, tune it, see a LIVE preview rendered from
// real rates, and copy a ready-to-embed snippet (HTML / React / Markdown /
// WordPress). The preview data is real (top exchangers for the pair); the embed
// script URL is the target contract for the hosted widget bundle.
definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Виджеты — Kurso' })

const PAIR = 'usdt-tinkoff'
const PAIR_LABEL = 'USDT → Тинькофф'
const { ranked, best, bestRow } = useDirectionRates(PAIR)
const { data: partner, load } = usePartner()
onMounted(load)
const refCode = computed(() => partner.value?.code ?? 'blogger')

interface WidgetType {
  key: 'rates' | 'calculator' | 'badge' | 'card' | 'button'
  name: string
  sub: string
  icon: string
}
const TYPES: WidgetType[] = [
  {
    key: 'rates',
    name: 'Таблица курсов',
    sub: 'Топ обменников по паре',
    icon: '<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 9h18M3 14h18M9 4v16"/>',
  },
  {
    key: 'calculator',
    name: 'Калькулятор',
    sub: 'Отдаю / получаю',
    icon: '<rect x="4" y="2" width="16" height="20" rx="2"/><path d="M8 6h8M8 11h2M14 11h2M8 16h6"/>',
  },
  {
    key: 'badge',
    name: 'Бейдж',
    sub: '«Лучший курс»',
    icon: '<path d="M12 3 4 6v6c0 5 3.5 7.5 8 9 4.5-1.5 8-4 8-9V6Z"/><path d="m9 12 2 2 4-4"/>',
  },
  {
    key: 'card',
    name: 'Карточка',
    sub: 'Обменника с курсом',
    icon: '<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 10h6v9"/>',
  },
  {
    key: 'button',
    name: 'Кнопка',
    sub: '«Лучший курс»',
    icon: '<rect x="3" y="8" width="18" height="8" rx="4"/><path d="M7 12h2"/>',
  },
]

const sel = ref(0)
const type = computed(() => TYPES[sel.value])

// settings
const rows = ref(3)
const amount = ref(1500)
const size = ref<'S' | 'M' | 'L'>('M')
const cardSlug = ref('')
const buttonLabel = ref('Лучший курс на Kurso')
const theme = ref<'dark' | 'light' | 'auto'>('dark')
const accent = ref('#2E7DF2')
const showReserve = ref(true)
const showLogo = ref(true)
const autoRefresh = ref(false)
const codeTab = ref<'html' | 'react' | 'markdown' | 'wordpress'>('html')
const copied = ref(false)

const ACCENTS = ['#2E7DF2', '#2BC58C', '#D99A33', '#8B5CF6']

// real preview data
const topRows = computed(() => ranked.value.slice(0, rows.value))
const bestRate = computed(() => (best.value != null ? fmtNumber(best.value, 2) : '—'))
const bestName = computed(() => bestRow.value?.exchangerName ?? '—')
const receive = computed(() => (best.value != null ? fmtNumber(amount.value * best.value, 0) : '—'))
const card = computed(() => {
  const slug = cardSlug.value
  const r = ranked.value.find((x) => x.row.exchangerSlug === slug) ?? ranked.value[0]
  return r?.row ?? null
})
const cardOptions = computed(() => ranked.value.map((r) => r.row))
watch(cardOptions, (opts) => {
  if (!cardSlug.value && opts.length) cardSlug.value = opts[0].exchangerSlug
})

const light = computed(() => theme.value === 'light')

// generated embed code
const params = computed<Record<string, string>>(() => {
  const base: Record<string, string> = {
    'data-type': type.value.key,
    'data-pair': 'USDT-TINKOFF',
    'data-theme': theme.value,
    'data-accent': accent.value,
    'data-reserve': String(showReserve.value),
    'data-logo': String(showLogo.value),
    'data-ref': refCode.value,
  }
  if (type.value.key === 'rates') base['data-rows'] = String(rows.value)
  if (type.value.key === 'calculator') base['data-amount'] = String(amount.value)
  if (type.value.key === 'badge') base['data-size'] = size.value.toLowerCase()
  if (type.value.key === 'card') base['data-exchanger'] = cardSlug.value
  if (type.value.key === 'button') base['data-label'] = buttonLabel.value
  return base
})

const code = computed(() => {
  const p = params.value
  const attrs = Object.entries(p)
    .map(([k, v]) => `${k}="${v}"`)
    .join(' ')
  if (codeTab.value === 'html') {
    // Build the embed tag name from a variable so the SFC parser never sees a
    // literal closing script tag that would end this block early.
    const tag = 'script'
    return `<!-- Виджет Kurso · ${type.value.name} -->\n<${tag} src="https://widget.kurso.io/v2/embed.js" async></${tag}>\n<div class="kurso-widget" ${attrs}></div>`
  }
  if (codeTab.value === 'react') {
    const props = Object.entries(p)
      .map(([k, v]) => `${k.replace(/^data-/, '')}="${v}"`)
      .join(' ')
    return `import { KursoWidget } from '@kurso/widgets'\n\n<KursoWidget ${props} />`
  }
  if (codeTab.value === 'markdown') {
    const q = Object.entries(p)
      .map(([k, v]) => `${k.replace(/^data-/, '')}=${encodeURIComponent(v)}`)
      .join('&')
    return `[![Kurso · ${type.value.name}](https://widget.kurso.io/v2/img?${q})](https://kurso.io/?ref=${refCode.value})`
  }
  const sc = Object.entries(p)
    .map(([k, v]) => `${k.replace(/^data-/, '')}="${v}"`)
    .join(' ')
  return `[kurso_widget ${sc}]`
})

async function copyCode() {
  try {
    await navigator.clipboard.writeText(code.value)
    copied.value = true
    setTimeout(() => (copied.value = false), 1600)
  } catch {
    /* clipboard blocked */
  }
}

const CODE_TABS: { key: typeof codeTab.value; label: string }[] = [
  { key: 'html', label: 'HTML' },
  { key: 'react', label: 'React' },
  { key: 'markdown', label: 'Markdown' },
  { key: 'wordpress', label: 'WordPress' },
]
</script>

<template>
  <div>
    <!-- breadcrumb -->
    <div class="mb-4 flex items-center gap-2 text-xs text-ink-faint">
      <NuxtLink to="/account/partner" class="hover:text-ink-muted">Партнёрка</NuxtLink>
      <svg
        width="12"
        height="12"
        viewBox="0 0 24 24"
        fill="none"
        stroke="#3A414A"
        stroke-width="2.4"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M9 6l6 6-6 6" />
      </svg>
      <span class="font-semibold text-ink-muted">Виджеты</span>
    </div>

    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Конструктор виджетов</h1>
    <p class="mb-6 mt-1 max-w-xl text-sm leading-relaxed text-ink-faint">
      Выберите вид виджета — превью, настройки и код обновятся. Данные в превью — живые курсы Kurso.
    </p>

    <!-- catalog -->
    <div class="mb-8 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
      <button
        v-for="(t, i) in TYPES"
        :key="t.key"
        type="button"
        class="relative rounded-2xl border p-4 text-left transition-all"
        :class="
          i === sel
            ? 'border-brand shadow-[0_0_0_3px_rgba(46,125,242,0.16)]'
            : 'border-line hover:border-line-strong'
        "
        @click="sel = i"
      >
        <span
          class="mb-3.5 flex h-10 w-10 items-center justify-center rounded-[11px]"
          :class="i === sel ? 'bg-brand text-white' : 'bg-surface-raised text-ink-muted'"
        >
          <!-- eslint-disable-next-line vue/no-v-html -->
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            v-html="t.icon"
          />
        </span>
        <div class="text-[15px] font-bold tracking-[-0.01em] text-ink">{{ t.name }}</div>
        <div class="mt-0.5 text-xs text-ink-faint">{{ t.sub }}</div>
        <span
          v-if="i === sel"
          class="absolute right-3.5 top-3.5 flex h-5 w-5 items-center justify-center rounded-full bg-brand"
        >
          <svg
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#fff"
            stroke-width="3"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M5 12l5 5L20 6" />
          </svg>
        </span>
      </button>
    </div>

    <!-- constructor: settings + preview -->
    <div class="mb-5 flex flex-col gap-5 lg:flex-row">
      <!-- settings -->
      <div
        class="flex w-full flex-col gap-[18px] rounded-2xl border border-line bg-surface p-5 lg:w-[340px] lg:flex-none"
      >
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Направление</div>
          <div
            class="flex items-center gap-2.5 rounded-xl border border-line bg-well px-3 py-2.5 text-sm font-semibold text-ink"
          >
            <span
              class="flex h-6 w-6 items-center justify-center rounded-full bg-[#26A17B] text-[11px] font-bold text-white"
              >₮</span
            >
            {{ PAIR_LABEL }}
          </div>
        </div>

        <!-- type-specific control -->
        <div v-if="type.key === 'rates'">
          <div class="mb-2 flex justify-between text-xs uppercase tracking-[0.05em] text-ink-faint">
            <span>Количество строк</span
            ><span class="tnum font-semibold text-brand-bright">{{ rows }}</span>
          </div>
          <input v-model.number="rows" type="range" min="2" max="6" class="w-full accent-brand" />
        </div>
        <div v-else-if="type.key === 'calculator'">
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">
            Сумма по умолчанию
          </div>
          <input
            v-model.number="amount"
            type="number"
            min="1"
            class="tnum w-full rounded-xl border border-line bg-well px-3 py-2.5 text-sm text-ink focus:border-brand focus:outline-none"
          />
        </div>
        <div v-else-if="type.key === 'badge'">
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Размер</div>
          <div class="flex gap-1.5 rounded-xl border border-line bg-well p-1.5">
            <button
              v-for="s in ['S', 'M', 'L'] as const"
              :key="s"
              type="button"
              class="flex-1 rounded-lg py-2 text-sm font-semibold transition-colors"
              :class="size === s ? 'bg-surface-chip text-ink' : 'text-ink-faint hover:text-ink'"
              @click="size = s"
            >
              {{ s }}
            </button>
          </div>
        </div>
        <div v-else-if="type.key === 'card'">
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Обменник</div>
          <select
            v-model="cardSlug"
            class="w-full rounded-xl border border-line bg-well px-3 py-2.5 text-sm font-semibold text-ink focus:border-brand focus:outline-none"
          >
            <option v-for="o in cardOptions" :key="o.exchangerSlug" :value="o.exchangerSlug">
              {{ o.exchangerName }}
            </option>
          </select>
        </div>
        <div v-else>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Текст кнопки</div>
          <input
            v-model="buttonLabel"
            class="w-full rounded-xl border border-line bg-well px-3 py-2.5 text-sm text-ink focus:border-brand focus:outline-none"
          />
        </div>

        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Тема</div>
          <div class="flex gap-1.5 rounded-xl border border-line bg-well p-1.5">
            <button
              v-for="th in ['dark', 'light', 'auto'] as const"
              :key="th"
              type="button"
              class="flex-1 rounded-lg py-2 text-[13px] font-semibold capitalize transition-colors"
              :class="theme === th ? 'bg-surface-chip text-ink' : 'text-ink-faint hover:text-ink'"
              @click="theme = th"
            >
              {{ { dark: 'Тёмная', light: 'Светлая', auto: 'Авто' }[th] }}
            </button>
          </div>
        </div>
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Акцент</div>
          <div class="flex gap-2.5">
            <button
              v-for="c in ACCENTS"
              :key="c"
              type="button"
              class="h-[30px] w-[30px] rounded-full transition-transform"
              :class="accent === c ? 'ring-2 ring-white ring-offset-2 ring-offset-surface' : ''"
              :style="{ background: c }"
              @click="accent = c"
            />
          </div>
        </div>
        <div class="flex flex-col gap-3">
          <label class="flex items-center justify-between text-sm text-ink"
            ><span>Показывать резерв</span><KToggle v-model="showReserve"
          /></label>
          <label class="flex items-center justify-between text-sm text-ink"
            ><span>Логотип Kurso</span><KToggle v-model="showLogo"
          /></label>
          <label class="flex items-center justify-between text-sm text-ink"
            ><span>Авто-обновление</span><KToggle v-model="autoRefresh"
          /></label>
        </div>
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">UTM-метка</div>
          <div
            class="tnum rounded-xl border border-line bg-well px-3 py-2.5 text-[13px] text-ink-body"
          >
            ref={{ refCode }}
          </div>
        </div>
      </div>

      <!-- preview -->
      <div
        class="flex min-w-0 flex-1 flex-col overflow-hidden rounded-2xl border border-line bg-well"
      >
        <div class="flex items-center justify-between border-b border-line px-[18px] py-3">
          <span class="font-label text-xs uppercase tracking-[0.05em] text-ink-faint"
            >Живое превью</span
          >
          <span class="flex items-center gap-1.5 text-xs text-ink-faint"
            ><KStatusDot tone="success" :size="7" pulse />обновляется</span
          >
        </div>
        <div
          class="flex min-h-[380px] flex-1 items-center justify-center overflow-auto p-10"
          :class="light ? 'bg-[#F4F5F7]' : ''"
          style="
            background-image: radial-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px);
            background-size: 18px 18px;
          "
        >
          <!-- RATES TABLE -->
          <div
            v-if="type.key === 'rates'"
            class="w-[380px] overflow-hidden rounded-2xl border shadow-modal"
            :class="
              light
                ? 'border-[#E3E6EA] bg-white text-[#111]'
                : 'border-line bg-surface-raised text-ink'
            "
          >
            <div
              class="flex items-center justify-between border-b px-4 py-3.5"
              :class="light ? 'border-[#EEF0F3]' : 'border-line'"
            >
              <div class="flex items-center gap-2 text-sm font-bold">
                <span
                  class="flex h-[22px] w-[22px] items-center justify-center rounded-full bg-[#26A17B] text-[10px] font-bold text-white"
                  >₮</span
                >{{ PAIR_LABEL }}
              </div>
              <span v-if="showLogo" class="text-[11px] text-ink-faint">Kurso</span>
            </div>
            <div
              v-for="(r, i) in topRows"
              :key="r.row.exchangerSlug"
              class="flex items-center justify-between border-b px-4 py-3 last:border-b-0"
              :class="[
                light ? 'border-[#F1F3F5]' : 'border-line-faint',
                i === 0 ? (light ? 'bg-brand/[0.05]' : 'bg-brand/[0.06]') : '',
              ]"
            >
              <div class="flex items-center gap-2.5">
                <span
                  class="flex h-[22px] w-[22px] items-center justify-center rounded-full text-[9px] font-bold text-white"
                  :style="{ background: i === 0 ? accent : '#3A4452' }"
                  >{{ r.row.exchangerName.slice(0, 2).toUpperCase() }}</span
                ><span class="text-[13px] font-semibold">{{ r.row.exchangerName }}</span>
              </div>
              <span class="tnum text-sm font-bold">{{ fmtNumber(r.rate, 2) }} ₽</span>
            </div>
            <div class="p-3">
              <div
                class="rounded-[9px] py-2.5 text-center text-[13px] font-semibold text-white"
                :style="{ background: accent }"
              >
                Все курсы на Kurso →
              </div>
            </div>
          </div>

          <!-- CALCULATOR -->
          <div
            v-else-if="type.key === 'calculator'"
            class="w-[340px] overflow-hidden rounded-2xl border shadow-modal"
            :class="
              light
                ? 'border-[#E3E6EA] bg-white text-[#111]'
                : 'border-line bg-surface-raised text-ink'
            "
          >
            <div
              class="flex items-center justify-between border-b px-4 py-3.5"
              :class="light ? 'border-[#EEF0F3]' : 'border-line'"
            >
              <span class="text-sm font-bold">Калькулятор обмена</span
              ><span v-if="showLogo" class="text-[11px] text-ink-faint">Kurso</span>
            </div>
            <div class="p-4">
              <div
                class="rounded-xl border px-3 py-2.5"
                :class="light ? 'border-[#E3E6EA] bg-[#F7F8FA]' : 'border-line bg-well'"
              >
                <div class="mb-1.5 text-[11px] text-ink-faint">Отдаёте</div>
                <div class="flex items-center gap-2.5">
                  <span
                    class="flex h-[22px] w-[22px] items-center justify-center rounded-full bg-[#26A17B] text-[10px] font-bold text-white"
                    >₮</span
                  ><span class="text-sm font-semibold">USDT</span
                  ><span class="tnum flex-1 text-right text-xl font-bold">{{
                    fmtNumber(amount, 0)
                  }}</span>
                </div>
              </div>
              <div class="my-2 flex justify-center">
                <span
                  class="flex h-8 w-8 items-center justify-center rounded-[10px] text-white"
                  :style="{ background: accent }"
                  ><svg
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M7 4v16M7 20l-3-3M7 20l3-3M17 20V4M17 4l-3 3M17 4l3 3" /></svg
                ></span>
              </div>
              <div
                class="rounded-xl border px-3 py-2.5"
                :class="light ? 'border-[#E3E6EA] bg-[#F7F8FA]' : 'border-line bg-well'"
              >
                <div class="mb-1.5 text-[11px] text-ink-faint">Получаете</div>
                <div class="flex items-center gap-2.5">
                  <span
                    class="flex h-[22px] w-[22px] items-center justify-center rounded-full bg-[#FFDD2D] text-[10px] font-extrabold text-[#111]"
                    >Т</span
                  ><span class="text-sm font-semibold">Тинькофф</span
                  ><span class="tnum flex-1 text-right text-xl font-bold">{{ receive }}</span>
                </div>
              </div>
              <div
                class="mt-3 flex items-center justify-between rounded-[10px] border px-3 py-2.5 text-xs"
                :class="light ? 'border-[#E3E6EA] bg-[#F7F8FA]' : 'border-line bg-well'"
              >
                <span class="text-ink-faint">Лучший · {{ bestName }}</span
                ><span class="tnum font-bold">{{ bestRate }} ₽</span>
              </div>
              <div
                class="mt-3 rounded-[11px] py-3 text-center text-sm font-semibold text-white"
                :style="{ background: accent }"
              >
                Обменять на Kurso →
              </div>
            </div>
          </div>

          <!-- BADGE -->
          <div
            v-else-if="type.key === 'badge'"
            class="inline-flex items-center gap-3 rounded-2xl border shadow-modal"
            :class="[
              light
                ? 'border-[#E3E6EA] bg-white text-[#111]'
                : 'border-line bg-surface-raised text-ink',
              size === 'S' ? 'px-3.5 py-3' : size === 'L' ? 'px-6 py-5' : 'px-[18px] py-[15px]',
            ]"
          >
            <span
              class="flex flex-none items-center justify-center rounded-xl"
              :class="size === 'S' ? 'h-8 w-8' : size === 'L' ? 'h-14 w-14' : 'h-[42px] w-[42px]'"
              :style="{ background: `linear-gradient(150deg,${accent},${accent})` }"
              ><svg
                :width="size === 'L' ? 28 : 21"
                :height="size === 'L' ? 28 : 21"
                viewBox="0 0 24 24"
                fill="none"
                stroke="#fff"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path
                  d="M7 4 L7 20 M7 20 L4 17 M7 20 L10 17 M17 20 L17 4 M17 4 L14 7 M17 4 L20 7"
                /></svg
            ></span>
            <div>
              <div class="text-xs text-ink-faint">{{ PAIR_LABEL }} · лучший курс</div>
              <div
                class="tnum font-extrabold tracking-[-0.02em]"
                :class="size === 'L' ? 'text-3xl' : 'text-2xl'"
              >
                {{ bestRate }} ₽
              </div>
              <div class="mt-0.5 text-xs text-ink-muted">
                <span class="text-warning">★</span>
                <span class="tnum font-semibold">{{ bestRow?.ratingAvg?.toFixed(1) ?? '—' }}</span>
                · обновлено на Kurso
              </div>
            </div>
          </div>

          <!-- CARD -->
          <div
            v-else-if="type.key === 'card'"
            class="w-[320px] overflow-hidden rounded-2xl border shadow-modal"
            :class="
              light
                ? 'border-[#E3E6EA] bg-white text-[#111]'
                : 'border-line bg-surface-raised text-ink'
            "
          >
            <div class="flex items-center gap-3 px-4 pb-3.5 pt-4">
              <span
                class="flex h-11 w-11 flex-none items-center justify-center rounded-xl bg-[#3A4452] text-[15px] font-extrabold text-white"
                >{{ (card?.exchangerName ?? '?').slice(0, 2).toUpperCase() }}</span
              >
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-1.5">
                  <span class="text-base font-bold">{{ card?.exchangerName ?? '—' }}</span
                  ><span
                    v-if="card?.partner"
                    class="rounded-[5px] bg-brand/15 px-1.5 py-0.5 text-[10px] font-semibold text-brand-bright"
                    >Партнёр</span
                  >
                </div>
                <div class="mt-0.5 text-xs text-ink-muted">
                  <span class="text-warning">★</span>
                  <span class="tnum font-semibold">{{ card?.ratingAvg?.toFixed(1) ?? '—' }}</span> ·
                  <span class="tnum">{{ card?.reviewsCount ?? 0 }}</span>
                </div>
              </div>
            </div>
            <div class="flex items-end justify-between px-4 pb-3">
              <div>
                <div class="tnum text-[26px] font-extrabold tracking-[-0.02em]">
                  {{ card ? fmtNumber(parseFloat(card.rate), 2) : '—' }} ₽
                </div>
                <div class="mt-0.5 text-xs font-semibold text-success-bright">{{ PAIR_LABEL }}</div>
              </div>
              <div
                v-if="showReserve && card?.reserve"
                class="tnum text-right text-xs text-ink-faint"
              >
                резерв<br /><span class="font-semibold text-ink-muted"
                  >{{ fmtCompact(parseFloat(card.reserve)) }} ₽</span
                >
              </div>
            </div>
            <div class="px-4 pb-4">
              <div
                class="rounded-[11px] py-3 text-center text-sm font-semibold text-white"
                :style="{ background: accent }"
              >
                Перейти на сайт →
              </div>
            </div>
          </div>

          <!-- BUTTON -->
          <div v-else class="flex flex-col items-center gap-3">
            <div
              class="inline-flex items-center gap-2.5 rounded-2xl px-5 py-4 text-base font-semibold text-white shadow-[0_16px_40px_rgba(46,125,242,0.35)]"
              :style="{ background: accent }"
            >
              <span
                class="flex h-[26px] w-[26px] items-center justify-center rounded-full bg-white/[0.18] text-[11px] font-extrabold"
                >₮</span
              >
              {{ buttonLabel }}
              <span class="tnum font-extrabold">{{ bestRate }} ₽</span>
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.4"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M7 17 17 7M9 7h8v8" />
              </svg>
            </div>
            <div class="flex items-center gap-1.5 text-xs text-ink-faint">
              <KStatusDot tone="success" :size="7" />обновлено только что · на Kurso
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- code -->
    <div class="overflow-hidden rounded-2xl border border-line bg-surface">
      <div class="flex flex-wrap items-center gap-1 border-b border-line px-3 pt-3">
        <button
          v-for="t in CODE_TABS"
          :key="t.key"
          type="button"
          class="rounded-t-lg px-4 py-2.5 text-[13px] font-semibold transition-colors"
          :class="codeTab === t.key ? 'bg-well text-ink' : 'text-ink-muted hover:text-ink'"
          @click="codeTab = t.key"
        >
          {{ t.label }}
        </button>
        <button
          type="button"
          class="mb-1.5 ml-auto inline-flex items-center gap-2 rounded-lg border border-line-strong bg-surface-raised px-3.5 py-2 text-xs font-semibold text-ink transition-colors hover:border-[#3A4047]"
          @click="copyCode"
        >
          <svg
            width="13"
            height="13"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <rect x="9" y="9" width="11" height="11" rx="2" />
            <path d="M5 15V5a2 2 0 0 1 2-2h10" />
          </svg>
          {{ copied ? 'Скопировано' : 'Копировать' }}
        </button>
      </div>
      <pre
        class="overflow-x-auto bg-well px-5 py-5 font-mono text-[13px] leading-relaxed text-ink-body"
      ><code>{{ code }}</code></pre>
    </div>
  </div>
</template>
