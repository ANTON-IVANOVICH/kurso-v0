<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useCatalogStore } from '../stores/catalog'
import { METRICS, PARSERS, TOP_REVENUE, ATTENTION, type HealthTone } from '../lib/fixtures'
import { fmtCompact, fmtRub, fmtInt } from '../lib/format'
import { exchangerAvatar } from '../lib/avatar'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import UiIcon from '../components/ui/UiIcon.vue'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Chip from 'primevue/chip'

const catalog = useCatalogStore()

const exchangersValue = computed(() =>
  catalog.exchangers.length ? String(catalog.exchangers.length) : '128',
)

const metrics = computed(() => [
  { label: 'Обменники', value: exchangersValue.value },
  { label: 'Парсеры', value: String(METRICS.parsersHealthy), sub: `/${METRICS.parsersTotal}` },
  { label: 'Клики/день', value: fmtCompact(METRICS.clicksPerDay) },
  { label: 'Доход/мес', value: fmtRub(METRICS.revenuePerMonth) },
  { label: 'Модерация', value: String(METRICS.moderation), tone: 'brand' as const },
  { label: 'Жалобы', value: String(METRICS.complaints), tone: 'danger' as const },
])

const toneClass: Record<string, string> = {
  brand: 'text-brand-light',
  danger: 'text-danger-text',
}

const health = computed(() =>
  PARSERS.filter((p) =>
    ['main-rates', 'reserves-eth', 'reserves-btc', 'feed-ton'].includes(p.name),
  ),
)

const SEVERITY: Record<HealthTone, string> = {
  success: 'success',
  warn: 'warn',
  danger: 'danger',
  brand: 'info',
  muted: 'secondary',
}

// Literal class strings so Tailwind's content scanner keeps them.
const TONE_TEXT: Record<HealthTone, string> = {
  success: 'text-success-text',
  warn: 'text-warn-text',
  danger: 'text-danger-text',
  brand: 'text-brand-light',
  muted: 'text-ink-muted',
}

// Live "updated Ns ago" ticker (client-only, resets on mount).
const secs = ref(0)
let timer: ReturnType<typeof setInterval> | undefined
onMounted(() => (timer = setInterval(() => (secs.value += 1), 1000)))
onBeforeUnmount(() => clearInterval(timer))
</script>

<template>
  <div>
    <PageHeader title="Дашборд">
      <template #actions>
        <div class="flex items-center gap-2 text-xs text-ink-faint">
          <StatusDot tone="success" pulse />
          живые данные · обновлено <span class="tnum">{{ secs }}с</span>
        </div>
      </template>
    </PageHeader>

    <!-- metric tiles -->
    <div class="mb-4 grid grid-cols-2 gap-2.5 sm:grid-cols-3 lg:grid-cols-6">
      <Card
        v-for="m in metrics"
        :key="m.label"
        class="border border-line"
        :pt="{ body: { class: 'p-3.5' }, content: { class: 'p-0' } }"
      >
        <template #content>
          <div class="mb-1.5 text-[11px] text-ink-faint">{{ m.label }}</div>
          <div class="tnum text-[21px] font-bold" :class="m.tone ? toneClass[m.tone] : ''">
            {{ m.value }}<span v-if="m.sub" class="text-[13px] text-ink-faint">{{ m.sub }}</span>
          </div>
        </template>
      </Card>
    </div>

    <!-- attention -->
    <div class="mb-4 rounded-xl border border-danger/30 bg-danger/[0.07] p-4">
      <div class="mb-3 flex items-center gap-2 text-sm font-bold">
        <UiIcon name="alert" :size="16" class="text-danger-text" />Требуют внимания
      </div>
      <div class="flex flex-wrap gap-2.5">
        <Chip v-for="(a, i) in ATTENTION" :key="i" class="!bg-surface !border !border-line-strong">
          <StatusDot :tone="a.tone" :size="6" />
          <span class="text-xs">{{ a.text }}</span>
          <span v-if="a.detail" class="text-xs text-ink-faint">· {{ a.detail }}</span>
        </Chip>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-[1.5fr_1fr]">
      <!-- parser health -->
      <Card :pt="{ body: { class: 'p-0' }, content: { class: 'p-0' } }" class="border border-line">
        <template #content>
          <div class="flex items-center justify-between border-b border-line px-4 py-3">
            <span class="text-sm font-bold">Здоровье парсеров</span>
            <RouterLink to="/parsers" class="text-xs font-semibold text-brand-light"
              >Все →</RouterLink
            >
          </div>
          <DataTable :value="health" size="small" class="text-[13px]">
            <Column header="Парсер">
              <template #body="{ data }">
                <div class="flex items-center gap-2">
                  <StatusDot :tone="data.dot" /><span>{{ data.name }}</span>
                </div>
              </template>
            </Column>
            <Column header="Посл. запуск">
              <template #body="{ data }">
                <span class="tnum" :class="TONE_TEXT[data.lastRunTone as HealthTone]">{{
                  data.lastRun
                }}</span>
              </template>
            </Column>
            <Column header="Latency">
              <template #body="{ data }"
                ><span class="tnum text-ink-muted">{{ data.latency }}</span></template
              >
            </Column>
            <Column header="Success">
              <template #body="{ data }">
                <Tag :value="data.success" :severity="SEVERITY[data.successTone as HealthTone]" />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <!-- top revenue -->
      <Card :pt="{ body: { class: 'p-0' }, content: { class: 'p-0' } }" class="border border-line">
        <template #content>
          <div class="border-b border-line px-4 py-3 text-sm font-bold">Топ по доходу</div>
          <ul>
            <li
              v-for="row in TOP_REVENUE"
              :key="row.rank"
              class="flex items-center gap-3 border-b border-line-soft px-4 py-2.5 last:border-0"
            >
              <span class="tnum w-3.5 text-xs text-ink-faint">{{ row.rank }}</span>
              <span
                class="flex h-[26px] w-[26px] flex-none items-center justify-center rounded-full text-[10px] font-bold text-white"
                :style="{ background: exchangerAvatar(row.name.toLowerCase(), row.name).color }"
                >{{ exchangerAvatar(row.name.toLowerCase(), row.name).initials }}</span
              >
              <span class="flex-1 text-[13px] font-semibold">{{ row.name }}</span>
              <span class="tnum text-[13px] font-semibold">{{ fmtRub(row.revenue) }}</span>
            </li>
          </ul>
        </template>
      </Card>
    </div>

    <p class="mt-6 flex items-center gap-2 text-xs text-ink-fainter">
      <UiIcon name="alert" :size="13" />
      Метрики дохода, кликов и здоровья парсеров — демо-данные до появления админского API. Число
      обменников — из живого каталога ({{ fmtInt(catalog.exchangers.length) }}).
    </p>
  </div>
</template>
