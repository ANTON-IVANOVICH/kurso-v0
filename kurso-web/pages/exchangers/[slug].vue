<script setup lang="ts">
import { computed } from 'vue'

// Страница обменника — вариант 3 (аналитика и доверие) + блок отзывов из варианта 1.
// Идентичность и позиция по направлениям — живые данные; аналитические метрики
// считаются из реальных курсов (спред к рынку, ранг), без выдуманных чисел.
const route = useRoute()
const slug = computed(() => String(route.params.slug))
const apiBase = useApiBase()
const { t, plural } = useI18n()

const { data: exchanger, state: exState } = useExchangerQuery(slug)
const { data: board } = useExchangerBoard(slug)

const notFound = computed(() => exState.value.status === 'success' && !exchanger.value)
const av = computed(() =>
  exchanger.value ? exchangerAvatar(exchanger.value.slug, exchanger.value.name) : null,
)
const entries = computed(() => board.value ?? [])
const primary = computed(() => entries.value[0] ?? null)
const avgSpread = computed(() =>
  entries.value.length
    ? entries.value.reduce((a, e) => a + e.spreadPct, 0) / entries.value.length
    : null,
)

// primary direction market band (min…max) + this exchanger's position within it
const primaryRates = computed(() =>
  primary.value ? primary.value.rows.map((r) => parseFloat(r.rate)).filter(Number.isFinite) : [],
)
const primaryMin = computed(() => (primaryRates.value.length ? Math.min(...primaryRates.value) : 0))
const primaryMax = computed(() => (primaryRates.value.length ? Math.max(...primaryRates.value) : 0))
const myPos = computed(() => {
  const mn = primaryMin.value
  const mx = primaryMax.value
  if (!primary.value || mx <= mn) return 100
  return ((primary.value.myRate - mn) / (mx - mn)) * 100
})

function goHref(dir?: string): string {
  return `${apiBase}/go/${slug.value}${dir ? `?direction=${dir}` : ''}`
}
function fmtRate(n: number): string {
  return `${fmtNumber(n, n >= 1000 ? 0 : 2)} ₽`
}
function fmtPct(p: number): string {
  return `${p >= 0 ? '+' : '−'}${Math.abs(p).toFixed(1)}%`
}
function fmtReserve(s: string | null | undefined): string {
  return s ? `${fmtCompact(parseFloat(s))} ₽` : '—'
}

useSeoMeta({
  title: () => `${exchanger.value?.name ?? t('detail.breadcrumb')} · Kurso`,
  description: () => `${exchanger.value?.name ?? ''} — ${t('detail.allDirections')} · Kurso`,
})
</script>

<template>
  <div class="mx-auto max-w-[1320px] px-4 py-6 md:px-6 md:py-8">
    <!-- not found -->
    <div v-if="notFound" class="rounded-2xl border border-line bg-surface py-24 text-center">
      <div class="text-lg font-semibold">{{ t('detail.notFound') }}</div>
      <NuxtLink to="/exchangers" class="mt-3 inline-block text-sm text-brand-bright">{{
        t('detail.toAll')
      }}</NuxtLink>
    </div>

    <template v-else-if="exchanger">
      <!-- breadcrumb -->
      <div class="mb-4 flex items-center gap-2 text-[13px] text-ink-faint">
        <NuxtLink to="/exchangers" class="hover:text-ink-muted">{{
          t('detail.breadcrumb')
        }}</NuxtLink>
        <span class="text-line-strong">/</span>
        <span class="text-ink-muted">{{ exchanger.name }}</span>
      </div>

      <!-- identity -->
      <div class="mb-5 flex flex-wrap items-center gap-4">
        <div
          class="flex h-[52px] w-[52px] flex-none items-center justify-center rounded-[14px] text-lg font-extrabold text-white"
          :style="{ background: av?.color }"
        >
          {{ av?.initials }}
        </div>
        <div class="min-w-0 flex-1">
          <div class="flex flex-wrap items-center gap-2.5">
            <h1 class="text-[22px] font-extrabold tracking-[-0.02em]">{{ exchanger.name }}</h1>
            <KBadge v-if="exchanger.partner" tone="brand">{{ t('exchangers.partner') }}</KBadge>
            <KBadge v-if="exchanger.isVerified" tone="success">{{
              t('exchangers.verified')
            }}</KBadge>
          </div>
          <div class="mt-1 text-[13px] text-ink-muted">
            <span class="text-warning-deep">★</span>
            <span class="tnum font-semibold text-ink">{{
              exchanger.ratingAvg != null ? exchanger.ratingAvg.toFixed(1) : '—'
            }}</span>
            · <span class="tnum">{{ exchanger.reviewsCount }}</span>
            {{ plural(exchanger.reviewsCount, 'reviews')
            }}<template v-if="exchanger.onSince">
              · {{ t('detail.onKursoSince') }}
              <span class="tnum">{{ exchanger.onSince }}</span></template
            >
          </div>
        </div>
        <button
          type="button"
          class="hidden rounded-xl border border-line-strong bg-surface-raised px-[18px] py-[11px] text-sm font-semibold text-ink-muted transition-colors hover:text-ink sm:inline-flex"
        >
          {{ t('detail.track') }}
        </button>
        <a
          :href="goHref(primary?.direction.slug)"
          target="_blank"
          rel="noopener nofollow"
          class="inline-flex items-center gap-2 rounded-xl bg-brand px-[22px] py-[11px] text-sm font-semibold text-white transition-colors hover:bg-brand-hover"
        >
          {{ t('detail.goToSite') }}
          <svg
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M7 17 17 7M9 7h8v8" />
          </svg>
        </a>
      </div>

      <!-- primary rate + ranking -->
      <div class="mb-5 grid gap-5 lg:grid-cols-[1fr_340px]">
        <!-- primary direction: rate + market band -->
        <div class="rounded-[20px] border border-line bg-surface p-[22px]">
          <template v-if="primary">
            <div class="mb-4 flex flex-wrap items-center gap-2">
              <span
                v-for="code in [primary.direction.fromCode, primary.direction.toCode]"
                :key="code"
                class="flex h-6 w-6 flex-none items-center justify-center rounded-full text-[10px] font-bold"
                :style="{
                  background: currencyBadge(code).color,
                  color: currencyBadge(code).dark ? '#111' : '#fff',
                }"
                >{{ currencyBadge(code).symbol }}</span
              >
              <span class="ml-1 text-sm font-semibold"
                >{{ primary.direction.fromCode }} → {{ primary.direction.toName }}</span
              >
            </div>
            <div class="flex items-baseline gap-3">
              <span class="tnum text-[34px] font-extrabold tracking-[-0.02em]">{{
                fmtRate(primary.myRate)
              }}</span>
              <span
                class="inline-flex items-center gap-1 text-sm font-semibold"
                :class="primary.spreadPct >= 0 ? 'text-success-bright' : 'text-danger'"
              >
                <span class="tnum">{{ fmtPct(primary.spreadPct) }} {{ t('detail.toMarket') }}</span>
              </span>
            </div>

            <!-- market band: min … avg … max with our position -->
            <div class="mt-6">
              <div class="relative h-2 rounded-full bg-well">
                <div
                  class="absolute -top-1 h-4 w-[3px] -translate-x-1/2 rounded-full bg-brand"
                  :style="{ left: `${myPos}%` }"
                />
              </div>
              <div class="tnum mt-2 flex justify-between text-[11px] text-ink-faint">
                <span>{{ t('detail.marketMin') }} {{ fmtRate(primaryMin) }}</span>
                <span>{{ t('detail.marketMid') }} {{ fmtRate(primary.marketAvg) }}</span>
                <span>{{ t('detail.marketMax') }} {{ fmtRate(primaryMax) }}</span>
              </div>
            </div>

            <div
              class="mt-5 flex items-center gap-2.5 rounded-xl border px-4 py-3 text-[13px]"
              :class="
                primary.spreadPct >= 0
                  ? 'border-success/25 bg-success/[0.07] text-ink-bright'
                  : 'border-line bg-well text-ink-muted'
              "
            >
              <span
                class="h-2 w-2 flex-none rounded-full"
                :class="primary.spreadPct >= 0 ? 'bg-success-bright' : 'bg-ink-faint'"
              />
              <span v-if="primary.spreadPct >= 0"
                >{{ t('detail.aboveMarket') }}
                <span class="font-semibold text-success-bright">{{
                  t('detail.aboveMarketStrong')
                }}</span>
                {{ t('detail.aboveMarketTail') }}</span
              >
              <span v-else>{{ t('detail.belowMarket') }}</span>
            </div>
          </template>
          <div v-else class="py-10 text-center text-sm text-ink-faint">
            {{ t('detail.noActiveRates') }}
          </div>
        </div>

        <!-- ranking for the primary direction -->
        <div class="overflow-hidden rounded-[20px] border border-line bg-surface">
          <div class="border-b border-line px-5 py-[18px]">
            <div class="mb-1.5 text-[13px] text-ink-faint">
              {{ t('detail.position')
              }}<template v-if="primary">
                {{ t('detail.positionBy') }} {{ primary.direction.fromCode }} →
                {{ primary.direction.toName }}</template
              >
            </div>
            <div v-if="primary" class="flex items-baseline gap-2">
              <span class="tnum text-[28px] font-extrabold text-brand-bright"
                >#{{ primary.rank }}</span
              >
              <span class="text-[13px] text-ink-muted">{{
                t('detail.ofExchangers', { n: primary.total })
              }}</span>
            </div>
          </div>
          <div
            v-for="(r, i) in primary?.rows.slice(0, 6) ?? []"
            :key="r.exchangerSlug"
            class="flex items-center gap-[11px] border-b border-line-subtle px-5 py-[11px] last:border-b-0"
            :class="r.exchangerSlug === slug ? 'bg-brand/[0.08]' : ''"
          >
            <span
              class="tnum w-4 flex-none text-[13px]"
              :class="r.exchangerSlug === slug ? 'font-bold text-brand-bright' : 'text-ink-faint'"
              >{{ i + 1 }}</span
            >
            <span
              class="flex h-[26px] w-[26px] flex-none items-center justify-center rounded-[7px] text-[10px] font-bold text-white"
              :style="{ background: exchangerAvatar(r.exchangerSlug, r.exchangerName).color }"
              >{{ exchangerAvatar(r.exchangerSlug, r.exchangerName).initials }}</span
            >
            <span class="flex-1 truncate text-[13px] font-medium">
              {{ r.exchangerName }}
              <span
                v-if="r.exchangerSlug === slug"
                class="ml-1.5 rounded-[5px] bg-brand/20 px-1.5 py-0.5 text-[11px] font-semibold text-brand-bright"
                >{{ t('detail.you') }}</span
              >
            </span>
            <span class="tnum flex-none text-[13px] font-bold">{{
              fmtRate(parseFloat(r.rate))
            }}</span>
          </div>
        </div>
      </div>

      <!-- real metric cards -->
      <div class="mb-5 grid grid-cols-2 gap-3 md:grid-cols-4 md:gap-4">
        <div class="rounded-2xl border border-line bg-surface p-[18px]">
          <div class="mb-2.5 text-xs text-ink-faint">{{ t('detail.avgSpread') }}</div>
          <div
            class="tnum text-[22px] font-extrabold md:text-[26px]"
            :class="(avgSpread ?? 0) >= 0 ? 'text-success-bright' : 'text-danger'"
          >
            {{ avgSpread != null ? fmtPct(avgSpread) : '—' }}
          </div>
          <div class="mt-2 text-xs text-ink-faint">{{ t('detail.avgSpreadNote') }}</div>
        </div>
        <div class="rounded-2xl border border-line bg-surface p-[18px]">
          <div class="mb-2.5 text-xs text-ink-faint">{{ t('detail.rating') }}</div>
          <div class="tnum text-[22px] font-extrabold md:text-[26px]">
            {{ exchanger.ratingAvg != null ? exchanger.ratingAvg.toFixed(1) : '—' }}
          </div>
          <div class="tnum mt-2 text-xs text-ink-faint">
            {{ exchanger.reviewsCount }} {{ plural(exchanger.reviewsCount, 'reviews') }}
          </div>
        </div>
        <div class="rounded-2xl border border-line bg-surface p-[18px]">
          <div class="mb-2.5 text-xs text-ink-faint">{{ t('detail.directionsCount') }}</div>
          <div class="tnum text-[22px] font-extrabold md:text-[26px]">
            {{ exchanger.directionsCount ?? entries.length }}
          </div>
          <div class="mt-2 truncate text-xs text-ink-faint">
            {{ (exchanger.assets ?? []).join(' · ') || '—' }}
          </div>
        </div>
        <div class="rounded-2xl border border-line bg-surface p-[18px]">
          <div class="mb-2.5 text-xs text-ink-faint">{{ t('detail.totalReserve') }}</div>
          <div class="tnum text-[22px] font-extrabold md:text-[26px]">
            {{ fmtReserve(exchanger.reserveTotal) }}
          </div>
          <div class="mt-2 text-xs text-ink-faint">{{ t('detail.byActiveRates') }}</div>
        </div>
      </div>

      <!-- directions table -->
      <div
        v-if="entries.length"
        class="mb-5 overflow-hidden rounded-[20px] border border-line bg-surface"
      >
        <div class="flex items-center justify-between border-b border-line px-[22px] py-[18px]">
          <span class="text-[17px] font-bold">{{ t('detail.allDirections') }}</span>
          <span class="tnum text-[13px] text-ink-faint"
            >{{ entries.length }} {{ plural(entries.length, 'directions') }}</span
          >
        </div>
        <div
          class="hidden grid-cols-[1.7fr_1fr_1fr_1fr_auto] gap-4 border-b border-line px-[22px] py-3 text-xs uppercase tracking-[0.05em] text-ink-faint md:grid"
        >
          <div>{{ t('detail.colDirection') }}</div>
          <div>{{ t('detail.colRate') }}</div>
          <div>{{ t('detail.colVsMarket') }}</div>
          <div>{{ t('detail.colReserve') }}</div>
          <div />
        </div>
        <div
          v-for="e in entries"
          :key="e.direction.slug"
          class="flex flex-wrap items-center gap-x-4 gap-y-2 border-b border-line-subtle px-4 py-3.5 last:border-b-0 md:grid md:grid-cols-[1.7fr_1fr_1fr_1fr_auto] md:px-[22px]"
        >
          <div class="flex min-w-0 flex-1 items-center gap-2 md:flex-none">
            <span
              v-for="code in [e.direction.fromCode, e.direction.toCode]"
              :key="code"
              class="flex h-7 w-7 flex-none items-center justify-center rounded-full text-[11px] font-bold"
              :style="{
                background: currencyBadge(code).color,
                color: currencyBadge(code).dark ? '#111' : '#fff',
              }"
              >{{ currencyBadge(code).symbol }}</span
            >
            <span class="ml-1 truncate text-sm font-medium"
              >{{ e.direction.fromCode }} → {{ e.direction.toName }}</span
            >
          </div>
          <div class="tnum text-[15px] font-bold">{{ fmtRate(e.myRate) }}</div>
          <div
            class="tnum text-[13px] font-semibold"
            :class="e.spreadPct >= 0 ? 'text-success-bright' : 'text-danger'"
          >
            {{ fmtPct(e.spreadPct) }}
          </div>
          <div class="tnum text-sm text-ink-muted">{{ fmtReserve(e.myRow.reserve) }}</div>
          <a
            :href="goHref(e.direction.slug)"
            target="_blank"
            rel="noopener nofollow"
            class="ml-auto inline-flex flex-none items-center justify-center rounded-lg bg-brand px-4 py-2 text-[13px] font-semibold text-white transition-colors hover:bg-brand-hover md:ml-0"
            >{{ t('detail.exchange') }}</a
          >
        </div>
      </div>

      <!-- reviews -->
      <ReviewsBlock :slug="slug" :exchanger-name="exchanger.name" />
    </template>

    <!-- loading -->
    <div
      v-else
      class="rounded-2xl border border-line bg-surface py-24 text-center text-sm text-ink-faint"
    >
      {{ t('detail.loading') }}
    </div>
  </div>
</template>
