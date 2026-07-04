<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRatesQuery, useDashboardQuery } from '../composables/useMerchant'
import { api, ApiError } from '../lib/api'
import { fmtRate, timeAgo } from '../lib/format'
import type { FeedStatus, MerchantRate } from '../types/models'
import PageHeader from '../components/ui/PageHeader.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import StatusDot from '../components/ui/StatusDot.vue'
import { useToast } from 'primevue/usetoast'
import Toast from 'primevue/toast'

const rates = useRatesQuery()
const dashboard = useDashboardQuery()
const toast = useToast()

const rows = computed(() => rates.data.value ?? [])
const refreshing = ref<Record<string, boolean>>({})
const bulkBusy = ref(false)

const FEED: Record<FeedStatus, { tone: 'success' | 'warn' | 'danger'; label: string }> = {
  ok: { tone: 'success', label: 'OK' },
  delayed: { tone: 'warn', label: 'задержка' },
  down: { tone: 'danger', label: 'нет данных' },
}

const downCount = computed(() => rows.value.filter((r) => r.feed !== 'ok').length)

async function refreshOne(row: MerchantRate) {
  refreshing.value = { ...refreshing.value, [row.directionId]: true }
  try {
    await api.post(`/partner/rates/${row.directionId}/refresh`)
    await Promise.all([rates.refetch(), dashboard.refetch()])
    toast.add({
      severity: 'success',
      summary: 'Фид обновлён',
      detail: `${row.fromCode} → ${row.toCode}`,
      life: 2500,
    })
  } catch (e) {
    const msg = e instanceof ApiError ? e.message : 'Не удалось обновить'
    toast.add({ severity: 'error', summary: 'Ошибка', detail: msg, life: 3500 })
  } finally {
    refreshing.value = { ...refreshing.value, [row.directionId]: false }
  }
}

async function refreshAll() {
  bulkBusy.value = true
  try {
    for (const row of rows.value) {
      try {
        await api.post(`/partner/rates/${row.directionId}/refresh`)
      } catch {
        /* skip rows without a configured rate */
      }
    }
    await Promise.all([rates.refetch(), dashboard.refetch()])
    toast.add({ severity: 'success', summary: 'Курсы обновлены', life: 2500 })
  } finally {
    bulkBusy.value = false
  }
}
</script>

<template>
  <div>
    <Toast position="bottom-right" />
    <PageHeader
      title="Курсы · направления"
      subtitle="Текущие курсы и состояние фида по вашим направлениям"
    >
      <template #actions>
        <Button
          label="Обновить все"
          icon="pi pi-refresh"
          severity="secondary"
          size="small"
          :loading="bulkBusy"
          @click="refreshAll"
        />
      </template>
    </PageHeader>

    <DataTable
      :value="rows"
      :loading="rates.isLoading.value"
      data-key="directionId"
      class="overflow-hidden rounded-[14px] border border-line"
    >
      <template #empty>
        <div class="py-10 text-center text-sm text-ink-faint">
          {{
            rates.error.value
              ? 'Не удалось загрузить курсы'
              : 'Курсы по направлениям пока не настроены'
          }}
        </div>
      </template>

      <Column header="Направление">
        <template #body="{ data }">
          <span class="text-sm font-medium">{{ data.fromCode }} → {{ data.toCode }}</span>
        </template>
      </Column>

      <Column header="Курс">
        <template #body="{ data }">
          <span
            class="tnum text-sm font-semibold"
            :class="data.feed === 'down' ? 'text-danger-text' : ''"
            >{{ fmtRate(data.rate) }}</span
          >
        </template>
      </Column>

      <Column header="Фид">
        <template #body="{ data }">
          <span
            class="flex items-center gap-2 text-xs"
            :class="`text-${FEED[data.feed as FeedStatus].tone}-text`"
          >
            <StatusDot :tone="FEED[data.feed as FeedStatus].tone" :size="7" />
            {{ FEED[data.feed as FeedStatus].label }}
            <span class="tnum text-ink-faint">· {{ timeAgo(data.fetchedAt) }}</span>
          </span>
        </template>
      </Column>

      <Column header-style="width:8rem">
        <template #body="{ data }">
          <Button
            :label="data.feed === 'down' ? 'Тест' : 'Обновить'"
            size="small"
            text
            :loading="refreshing[data.directionId]"
            @click="refreshOne(data)"
          />
        </template>
      </Column>
    </DataTable>

    <div
      v-if="downCount"
      class="mt-3 flex items-center gap-2.5 rounded-xl border border-line bg-well px-4 py-3.5"
    >
      <i class="pi pi-exclamation-triangle text-warn-text" />
      <span class="flex-1 text-[13px] text-ink-muted">
        Диагностика: <span class="tnum">{{ downCount }}</span> направлений требуют внимания —
        обновите фид.
      </span>
      <Button
        label="Обновить все"
        size="small"
        severity="secondary"
        :loading="bulkBusy"
        @click="refreshAll"
      />
    </div>
  </div>
</template>
