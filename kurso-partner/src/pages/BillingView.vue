<script setup lang="ts">
import { computed } from 'vue'
import { useBillingQuery } from '../composables/useMerchant'
import { fmtInt, fmtRubFull } from '../lib/format'
import PageHeader from '../components/ui/PageHeader.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'

const billing = useBillingQuery()
const data = computed(() => billing.data.value)
const payouts = computed(() => data.value?.payouts ?? [])

const monthName = new Intl.DateTimeFormat('ru-RU', { month: 'long' })

function periodLabel(startIso: string): string {
  const d = new Date(startIso)
  return `${monthName.format(d)} ${d.getFullYear()}`
}

const STATUS: Record<string, { label: string; severity: string }> = {
  pending: { label: 'ожидает', severity: 'warn' },
  paid: { label: 'оплачен', severity: 'info' },
  cancelled: { label: 'отменён', severity: 'secondary' },
}
</script>

<template>
  <div>
    <PageHeader
      title="Биллинг"
      subtitle="Партнёрская модель CPA — оплата за подтверждённый переход"
    />

    <div class="grid gap-4 lg:grid-cols-[340px_1fr] lg:items-start">
      <div class="flex flex-col gap-4">
        <!-- running balance (real, from this month's clicks) -->
        <div class="relative overflow-hidden rounded-[14px] border border-line bg-surface p-[22px]">
          <div
            class="pointer-events-none absolute -right-8 -top-10 h-44 w-44 rounded-full"
            style="background: radial-gradient(circle, rgba(46, 125, 242, 0.16), transparent 70%)"
          />
          <div class="relative">
            <div class="mb-2 text-[13px] text-ink-muted">Накоплено в этом месяце</div>
            <div class="tnum text-[34px] font-bold tracking-tight">
              {{ fmtRubFull(data?.currentMonth.estimated ?? 0) }}
            </div>
            <div class="mt-2 text-xs text-ink-faint">
              оценка по
              <span class="tnum">{{ fmtInt(data?.currentMonth.clicks ?? 0) }}</span> кликам · счёт
              формируется в конце месяца
            </div>
          </div>
        </div>

        <!-- tariff -->
        <div class="rounded-[14px] border border-line bg-surface p-5">
          <div class="mb-3 text-xs uppercase tracking-wide text-ink-faint">Тариф</div>
          <div class="mb-1 flex items-baseline gap-2">
            <span class="text-lg font-extrabold">Партнёрский</span>
            <span class="text-[13px] font-semibold text-brand-light">CPA</span>
          </div>
          <p class="mb-4 text-[13px] leading-relaxed text-ink-muted">
            Комиссия за подтверждённый переход. Без абонентской платы.
          </p>
          <div class="flex justify-between border-t border-line py-2.5">
            <span class="text-[13px] text-ink-muted">Ставка за переход</span>
            <span class="tnum text-[13px] font-semibold">{{ data?.perClick ?? 0 }} ₽</span>
          </div>
          <div class="flex justify-between border-t border-line py-2.5">
            <span class="text-[13px] text-ink-muted">Переходов в этом месяце</span>
            <span class="tnum text-[13px] font-semibold">{{
              fmtInt(data?.currentMonth.clicks ?? 0)
            }}</span>
          </div>
        </div>
      </div>

      <!-- payout history (real partner_payouts) -->
      <DataTable
        :value="payouts"
        :loading="billing.isLoading.value"
        data-key="id"
        class="overflow-hidden rounded-[14px] border border-line"
      >
        <template #header>
          <span class="text-[15px] font-bold">История платежей</span>
        </template>
        <template #empty>
          <div class="py-12 text-center text-sm text-ink-faint">
            Счетов пока нет. Первый счёт сформируется по итогам месяца.
          </div>
        </template>

        <Column header="Период">
          <template #body="{ data: row }">
            <span class="text-sm capitalize">{{ periodLabel(row.periodStart) }}</span>
          </template>
        </Column>
        <Column header="Переходы">
          <template #body="{ data: row }">
            <span class="tnum text-sm text-ink-muted">{{ fmtInt(row.clicksCount) }}</span>
          </template>
        </Column>
        <Column header="Сумма">
          <template #body="{ data: row }">
            <span class="tnum text-sm font-semibold"
              >{{ fmtInt(Number(row.amount)) }}
              {{ row.currency === 'RUB' ? '₽' : row.currency }}</span
            >
          </template>
        </Column>
        <Column header="Статус">
          <template #body="{ data: row }">
            <Tag :value="STATUS[row.status].label" :severity="STATUS[row.status].severity" />
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>
