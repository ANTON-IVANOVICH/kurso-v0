<script setup lang="ts">
import { useToast } from 'primevue/usetoast'
import { PARSERS, PARSER_TALLY, type HealthTone } from '../lib/fixtures'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Toast from 'primevue/toast'

const toast = useToast()

const SEVERITY: Record<HealthTone, string> = {
  success: 'success',
  warn: 'warn',
  danger: 'danger',
  brand: 'info',
  muted: 'secondary',
}
const TONE_TEXT: Record<HealthTone, string> = {
  success: 'text-success-text',
  warn: 'text-warn-text',
  danger: 'text-danger-text',
  brand: 'text-brand-light',
  muted: 'text-ink-muted',
}

function test(name: string) {
  toast.add({ severity: 'info', summary: 'Тестовый прогон', detail: name, life: 2000 })
}
function testAll() {
  toast.add({ severity: 'info', summary: 'Запущен прогон всех парсеров', life: 2500 })
}
</script>

<template>
  <div>
    <Toast />
    <PageHeader title="Парсеры" subtitle="Здоровье фидов и конфигурация">
      <template #actions>
        <Button
          label="Тест всех"
          icon="pi pi-play"
          size="small"
          severity="secondary"
          @click="testAll"
        />
      </template>
    </PageHeader>

    <!-- tally -->
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <Tag :value="`${PARSER_TALLY.healthy} healthy`" severity="success" icon="pi pi-circle-fill" />
      <Tag
        :value="`${PARSER_TALLY.degraded} деградация`"
        severity="warn"
        icon="pi pi-circle-fill"
      />
      <Tag :value="`${PARSER_TALLY.down} упали`" severity="danger" icon="pi pi-circle-fill" />
      <span class="ml-auto text-xs text-ink-faint"
        >Автопрогон каждые <span class="tnum">10с</span></span
      >
    </div>

    <DataTable
      :value="PARSERS"
      data-key="name"
      class="overflow-hidden rounded-[14px] border border-line"
    >
      <Column header="Парсер">
        <template #body="{ data }">
          <div class="tnum flex items-center gap-2">
            <StatusDot :tone="data.dot" />{{ data.name }}
          </div>
        </template>
      </Column>
      <Column field="exchanger" header="Обменник">
        <template #body="{ data }"
          ><span class="text-ink-muted">{{ data.exchanger }}</span></template
        >
      </Column>
      <Column field="format" header="Формат">
        <template #body="{ data }"
          ><span class="tnum text-ink-muted">{{ data.format }}</span></template
        >
      </Column>
      <Column field="interval" header="Интервал">
        <template #body="{ data }"
          ><span class="tnum text-ink-muted">{{ data.interval }}</span></template
        >
      </Column>
      <Column header="Посл. запуск">
        <template #body="{ data }">
          <span class="tnum" :class="TONE_TEXT[data.lastRunTone as HealthTone]">{{
            data.lastRun
          }}</span>
        </template>
      </Column>
      <Column field="latency" header="Latency">
        <template #body="{ data }"
          ><span class="tnum text-ink-muted">{{ data.latency }}</span></template
        >
      </Column>
      <Column header="Success">
        <template #body="{ data }"
          ><Tag :value="data.success" :severity="SEVERITY[data.successTone as HealthTone]"
        /></template>
      </Column>
      <Column header="Статус">
        <template #body="{ data }"
          ><Tag :value="data.status" :severity="SEVERITY[data.statusTone as HealthTone]"
        /></template>
      </Column>
      <Column header-style="width:8rem">
        <template #body="{ data }">
          <div class="flex justify-end gap-1">
            <Button label="тест" size="small" text @click="test(data.name)" />
            <Button label="логи" size="small" text severity="secondary" @click="test(data.name)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <p class="mt-4 text-xs text-ink-fainter">
      Демо-данные до появления рантайма парсеров и админского API.
    </p>
  </div>
</template>
