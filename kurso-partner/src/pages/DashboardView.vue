<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useDashboardQuery, useReviewsQuery, useRatesQuery } from '../composables/useMerchant'
import { fmtInt, timeAgo } from '../lib/format'
import Button from 'primevue/button'
import Rating from 'primevue/rating'
import StatusDot from '../components/ui/StatusDot.vue'

const router = useRouter()

const dashboard = useDashboardQuery()
const reviews = useReviewsQuery()
const rates = useRatesQuery()

const metrics = computed(() => dashboard.data.value?.metrics)
const traffic = computed(() => dashboard.data.value?.traffic ?? [])

// Most recent feed fetch across the exchanger's rates → "last update".
const lastFetch = computed(() => {
  const list = rates.data.value ?? []
  let latest = ''
  for (const r of list) if (r.fetchedAt && r.fetchedAt > latest) latest = r.fetchedAt
  return latest
})

const unanswered = computed(() => (reviews.data.value ?? []).filter((r) => !r.reply).slice(0, 3))

// Feed health buckets from the live rate rows.
const feed = computed(() => {
  const list = rates.data.value ?? []
  return {
    ok: list.filter((r) => r.feed === 'ok').length,
    delayed: list.filter((r) => r.feed === 'delayed').length,
    down: list.filter((r) => r.feed === 'down').length,
  }
})

// Traffic delta vs the previous day (dashboard series is 7 days).
const clicksDelta = computed(() => {
  const m = metrics.value
  if (!m) return null
  if (m.clicksYesterday === 0) return m.clicksToday > 0 ? 100 : 0
  return Math.round(((m.clicksToday - m.clicksYesterday) / m.clicksYesterday) * 100)
})

const maxTraffic = computed(() => Math.max(1, ...traffic.value.map((t) => t.clicks)))
const weekTotal = computed(() => traffic.value.reduce((s, t) => s + t.clicks, 0))
const dow = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб']

async function refreshFeed() {
  await Promise.all([dashboard.refetch(), rates.refetch(), reviews.refetch()])
}

function toReviews() {
  router.push({ name: 'reviews' })
}
</script>

<template>
  <div>
    <div class="mb-5 flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold tracking-tight md:text-[22px]">Главная</h1>
        <p class="mt-1 text-sm text-ink-faint">
          Последнее обновление фида:
          <span class="tnum">{{ timeAgo(lastFetch) }}</span>
        </p>
      </div>
      <Button
        label="Обновить курсы"
        icon="pi pi-refresh"
        severity="secondary"
        size="small"
        :loading="dashboard.isLoading.value"
        @click="refreshFeed"
      />
    </div>

    <!-- metrics -->
    <div class="mb-5 grid grid-cols-2 gap-3.5 lg:grid-cols-4">
      <div class="rounded-[14px] border border-line bg-surface p-4">
        <div class="mb-2.5 text-[13px] text-ink-muted">Курсы активны</div>
        <div class="tnum text-[26px] font-bold">
          {{ metrics?.ratesActive ?? 0
          }}<span class="text-base text-ink-faint">/{{ metrics?.ratesTotal ?? 0 }}</span>
        </div>
        <div
          class="mt-1.5 flex items-center gap-1.5 text-xs"
          :class="metrics?.ratesStale ? 'text-warn-text' : 'text-ink-faint'"
        >
          <StatusDot :tone="metrics?.ratesStale ? 'warn' : 'success'" :size="6" />
          {{ metrics?.ratesStale ? `${metrics.ratesStale} устарели` : 'фид в норме' }}
        </div>
      </div>

      <div class="rounded-[14px] border border-line bg-surface p-4">
        <div class="mb-2.5 text-[13px] text-ink-muted">Клики сегодня</div>
        <div class="tnum text-[26px] font-bold">{{ fmtInt(metrics?.clicksToday ?? 0) }}</div>
        <div
          class="mt-1.5 text-xs"
          :class="(clicksDelta ?? 0) >= 0 ? 'text-success-text' : 'text-danger-text'"
        >
          {{ clicksDelta === null ? '—' : `${clicksDelta >= 0 ? '+' : ''}${clicksDelta}% к вчера` }}
        </div>
      </div>

      <div class="rounded-[14px] border border-line bg-surface p-4">
        <div class="mb-2.5 text-[13px] text-ink-muted">Рейтинг</div>
        <div class="tnum flex items-baseline gap-1.5 text-[26px] font-bold">
          {{ (metrics?.ratingAvg ?? 0).toFixed(1) }}
          <span class="text-[15px] text-warn-text">★</span>
        </div>
        <div class="tnum mt-1.5 text-xs text-ink-faint">
          {{ fmtInt(metrics?.reviewsCount ?? 0) }} отзывов
        </div>
      </div>

      <div class="rounded-[14px] border border-line bg-surface p-4">
        <div class="mb-2.5 text-[13px] text-ink-muted">Новые отзывы</div>
        <div class="tnum text-[26px] font-bold">{{ metrics?.reviewsUnanswered ?? 0 }}</div>
        <div class="mt-1.5 text-xs text-brand-light">требуют ответа</div>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-[1.3fr_1fr]">
      <!-- reviews to answer -->
      <div class="rounded-[14px] border border-line bg-surface p-[18px]">
        <div class="mb-3.5 flex items-center justify-between">
          <span class="text-[15px] font-bold">Отзывы, требующие ответа</span>
          <button class="text-xs font-semibold text-brand-light" @click="toReviews">Все →</button>
        </div>
        <div v-if="unanswered.length" class="flex flex-col gap-2.5">
          <div
            v-for="r in unanswered"
            :key="r.id"
            class="rounded-xl border border-line bg-well p-3.5"
          >
            <div class="mb-2 flex items-center gap-2.5">
              <span class="text-[13px] font-semibold">{{ r.author }}</span>
              <Rating :model-value="r.rating" readonly />
              <span class="tnum ml-auto text-[11px] text-ink-faint">{{
                timeAgo(r.createdAt)
              }}</span>
            </div>
            <p class="mb-2.5 line-clamp-2 text-[13px] leading-snug text-ink-body">{{ r.body }}</p>
            <Button label="Ответить" size="small" @click="toReviews" />
          </div>
        </div>
        <div v-else class="py-8 text-center text-sm text-ink-faint">
          Все отзывы отвечены — новых нет.
        </div>
      </div>

      <div class="flex flex-col gap-4">
        <!-- feed health -->
        <div class="rounded-[14px] border border-line bg-surface p-[18px]">
          <div class="mb-3.5 text-[15px] font-bold">Здоровье фида</div>
          <div class="flex flex-col gap-3">
            <div class="flex items-center gap-2.5">
              <StatusDot tone="success" :size="8" />
              <span class="flex-1 text-[13px]">Активные фиды</span>
              <span class="tnum text-xs text-ink-faint">{{ feed.ok }} · OK</span>
            </div>
            <div class="flex items-center gap-2.5">
              <StatusDot tone="warn" :size="8" />
              <span class="flex-1 text-[13px]">С задержкой</span>
              <span
                class="tnum text-xs"
                :class="feed.delayed ? 'text-warn-text' : 'text-ink-faint'"
                >{{ feed.delayed }}</span
              >
            </div>
            <div class="flex items-center gap-2.5">
              <StatusDot :tone="feed.down ? 'danger' : 'muted'" :size="8" />
              <span class="flex-1 text-[13px]">Не отвечают</span>
              <span
                class="tnum text-xs"
                :class="feed.down ? 'text-danger-text' : 'text-ink-faint'"
                >{{ feed.down }}</span
              >
            </div>
          </div>
        </div>

        <!-- 7-day traffic -->
        <div class="rounded-[14px] border border-line bg-surface p-[18px]">
          <div class="mb-3.5 flex items-center justify-between">
            <span class="text-[15px] font-bold">Трафик · 7 дней</span>
            <span class="tnum text-xs text-ink-faint">{{ fmtInt(weekTotal) }}</span>
          </div>
          <div class="flex h-20 items-end gap-2">
            <div
              v-for="(t, i) in traffic"
              :key="t.day"
              class="flex-1 rounded-t-[5px]"
              :class="
                i === traffic.length - 1
                  ? 'bg-gradient-to-b from-brand-bright to-brand'
                  : 'bg-line-strong'
              "
              :style="{ height: `${Math.max(6, (t.clicks / maxTraffic) * 100)}%` }"
              :title="`${t.clicks}`"
            />
          </div>
          <div class="mt-2 flex justify-between font-mono text-[10px] text-ink-faint">
            <span v-for="t in traffic" :key="t.day">{{ dow[new Date(t.day).getUTCDay()] }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
