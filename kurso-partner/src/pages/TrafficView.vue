<script setup lang="ts">
import { computed, ref } from 'vue'
import { useTrafficQuery } from '../composables/useMerchant'
import { useAuthStore } from '../stores/auth'
import { fmtInt } from '../lib/format'
import PageHeader from '../components/ui/PageHeader.vue'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const auth = useAuthStore()
const days = ref(14)
const traffic = useTrafficQuery(() => days.value)

const data = computed(() => traffic.data.value)
const series = computed(() => data.value?.series ?? [])
const directions = computed(() => data.value?.directions ?? [])
const total = computed(() => data.value?.total ?? 0)
const maxDay = computed(() => Math.max(1, ...series.value.map((s) => s.clicks)))
const maxDir = computed(() => Math.max(1, ...directions.value.map((d) => d.clicks)))

// Total over the previous window of the same length, for a delta figure.
const half = computed(() => {
  const s = series.value
  const n = Math.floor(s.length / 2)
  const prev = s.slice(0, n).reduce((a, b) => a + b.clicks, 0)
  const curr = s.slice(s.length - n).reduce((a, b) => a + b.clicks, 0)
  if (prev === 0) return curr > 0 ? 100 : 0
  return Math.round(((curr - prev) / prev) * 100)
})

const TABS = [7, 14, 30]

function exportCsv() {
  const header = 'direction,from,to,clicks\n'
  const body = directions.value
    .map((d) => `${d.directionSlug},${d.fromCode},${d.toCode},${d.clicks}`)
    .join('\n')
  const blob = new Blob([header + body], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `traffic-${auth.merchant?.exchangerSlug ?? 'exchanger'}-${days.value}d.csv`
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<template>
  <div>
    <PageHeader title="Статистика трафика" subtitle="Клики по направлениям из каталога Kurso">
      <template #actions>
        <div class="flex rounded-[9px] border border-line bg-well p-0.5 text-xs">
          <button
            v-for="t in TABS"
            :key="t"
            class="rounded-md px-3 py-1.5 font-medium transition-colors"
            :class="days === t ? 'bg-brand text-white' : 'text-ink-muted hover:text-ink'"
            @click="days = t"
          >
            {{ t }}д
          </button>
        </div>
        <Button
          label="Экспорт CSV"
          icon="pi pi-download"
          severity="secondary"
          size="small"
          :disabled="!directions.length"
          @click="exportCsv"
        />
      </template>
    </PageHeader>

    <!-- headline + daily chart -->
    <div class="mb-4 rounded-[14px] border border-line bg-surface p-[22px]">
      <div class="mb-4 flex items-end gap-2">
        <div class="tnum text-[30px] font-bold">{{ fmtInt(total) }}</div>
        <div class="pb-2 text-[13px]" :class="half >= 0 ? 'text-success-text' : 'text-danger-text'">
          {{ half >= 0 ? '+' : '' }}{{ half }}% к прошлому периоду
        </div>
      </div>
      <div v-if="series.length" class="flex h-[148px] items-end gap-[5px]">
        <div
          v-for="(s, i) in series"
          :key="s.day"
          class="flex-1 rounded-t-[3px]"
          :class="
            i === series.length - 1
              ? 'bg-gradient-to-b from-brand-bright to-brand'
              : 'bg-line-strong'
          "
          :style="{ height: `${Math.max(4, (s.clicks / maxDay) * 100)}%` }"
          :title="`${s.day}: ${s.clicks}`"
        />
      </div>
      <div v-else class="py-10 text-center text-sm text-ink-faint">Нет данных за период</div>
    </div>

    <!-- by direction -->
    <DataTable
      :value="directions"
      :loading="traffic.isLoading.value"
      data-key="directionSlug"
      class="overflow-hidden rounded-[14px] border border-line"
    >
      <template #header>
        <span class="text-[15px] font-bold">Клики по направлениям</span>
      </template>
      <template #empty>
        <div class="py-10 text-center text-sm text-ink-faint">
          Переходов за выбранный период не было. Клики появляются, когда пользователи переходят к
          вам из каталога.
        </div>
      </template>

      <Column header="Направление">
        <template #body="{ data: d }">
          <span class="text-sm font-medium">{{ d.fromCode }} → {{ d.toCode }}</span>
        </template>
      </Column>
      <Column header="Доля">
        <template #body="{ data: d }">
          <div class="h-2 w-40 overflow-hidden rounded-full bg-line">
            <div class="h-full bg-brand" :style="{ width: `${(d.clicks / maxDir) * 100}%` }" />
          </div>
        </template>
      </Column>
      <Column header="Клики" header-style="width:7rem">
        <template #body="{ data: d }">
          <span class="tnum text-sm font-semibold">{{ fmtInt(d.clicks) }}</span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>
