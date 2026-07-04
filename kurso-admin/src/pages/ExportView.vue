<script setup lang="ts">
import { useToast } from 'primevue/usetoast'
import { EXPORTS, AUDIT, type ExportSet, type HealthTone } from '../lib/fixtures'
import PageHeader from '../components/ui/PageHeader.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Toast from 'primevue/toast'

const toast = useToast()

const ICON: Record<ExportSet['key'], string> = {
  clicks: 'pi pi-chart-bar text-brand-light',
  revenue: 'pi pi-chart-line text-success-text',
  reviews: 'pi pi-comments text-warn-text',
}

const TAG_SEVERITY: Record<HealthTone, string> = {
  success: 'success',
  warn: 'warn',
  danger: 'danger',
  brand: 'info',
  muted: 'secondary',
}

function download(s: ExportSet) {
  toast.add({ severity: 'info', summary: 'Экспорт CSV', detail: s.title, life: 2200 })
}
</script>

<template>
  <div>
    <Toast />
    <PageHeader title="Экспорт · Логи" subtitle="Выгрузка данных и аудит-лог действий" />

    <div class="grid gap-4 lg:grid-cols-[0.85fr_1.4fr]">
      <!-- export -->
      <div class="rounded-[14px] border border-line bg-surface p-4">
        <div class="text-sm font-bold">Экспорт данных</div>
        <div class="mb-4 mt-0.5 text-xs text-ink-faint">CSV · обновляется в реальном времени</div>
        <div class="flex flex-col gap-2.5">
          <div
            v-for="s in EXPORTS"
            :key="s.key"
            class="flex items-center gap-3 rounded-xl border border-line bg-well px-3.5 py-3"
          >
            <i :class="ICON[s.key]" />
            <div class="flex-1">
              <div class="text-[13px] font-semibold">{{ s.title }}</div>
              <div class="tnum text-[11px] text-ink-faint">{{ s.meta }}</div>
            </div>
            <Button
              icon="pi pi-download"
              text
              rounded
              severity="secondary"
              aria-label="Скачать"
              @click="download(s)"
            />
          </div>
        </div>
      </div>

      <!-- audit log -->
      <div class="overflow-hidden rounded-[14px] border border-line bg-surface">
        <div class="border-b border-line px-4 py-3 text-sm font-bold">Аудит-лог</div>
        <DataTable :value="AUDIT" data-key="time" :show-headers="false" class="text-[13px]">
          <Column header-style="width:4rem">
            <template #body="{ data }"
              ><span class="tnum text-ink-fainter">{{ data.time }}</span></template
            >
          </Column>
          <Column header-style="width:2.5rem">
            <template #body="{ data }">
              <span
                class="flex h-[22px] w-[22px] items-center justify-center rounded-md text-[9px] font-bold text-white"
                :style="{ background: data.actorColor }"
                >{{ data.actor }}</span
              >
            </template>
          </Column>
          <Column header="Событие">
            <template #body="{ data }">
              <span class="text-ink-body">{{ data.pre }}</span
              ><span v-if="data.mark" :style="{ color: data.markColor }">{{ data.mark }}</span
              ><span class="text-ink-body">{{ data.post }}</span>
            </template>
          </Column>
          <Column header-style="width:8rem">
            <template #body="{ data }">
              <Tag :value="data.tag" :severity="TAG_SEVERITY[data.tagTone as HealthTone]" />
            </template>
          </Column>
        </DataTable>
      </div>
    </div>

    <p class="mt-4 text-xs text-ink-fainter">
      Экспорт и аудит-лог — демо до появления соответствующего API.
    </p>
  </div>
</template>
