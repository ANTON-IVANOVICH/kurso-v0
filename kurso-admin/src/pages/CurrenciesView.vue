<script setup lang="ts">
import { ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useCatalogStore } from '../stores/catalog'
import type { CurrencyKind } from '../types/models'
import { MAP_QUEUE, type MapCandidate } from '../lib/fixtures'
import PageHeader from '../components/ui/PageHeader.vue'
import UiIcon from '../components/ui/UiIcon.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import Toast from 'primevue/toast'

const catalog = useCatalogStore()
const toast = useToast()

const search = ref('')
const KIND_META: Record<CurrencyKind, { label: string; severity: string }> = {
  crypto: { label: 'крипта', severity: 'info' },
  fiat: { label: 'фиат', severity: 'success' },
  cash: { label: 'наличные', severity: 'warn' },
}

// Local copy of the unknown-currency mapping queue (demo). Confirming a binding
// pops it off the list.
const queue = ref<MapCandidate[]>([...MAP_QUEUE])
function bind(c: MapCandidate) {
  queue.value = queue.value.filter((x) => x.source !== c.source)
  toast.add({
    severity: 'success',
    summary: `${c.source} привязана`,
    detail: c.suggestion ? `→ ${c.suggestion.name}` : undefined,
    life: 2500,
  })
}
</script>

<template>
  <div>
    <Toast />
    <PageHeader title="Валюты" subtitle="Справочник валют и очередь маппинга от парсеров" />

    <div class="grid gap-4 lg:grid-cols-[1.6fr_1fr]">
      <!-- currencies table -->
      <DataTable
        :value="catalog.currencies"
        :loading="catalog.loading"
        data-key="id"
        paginator
        :rows="12"
        :global-filter-fields="['code', 'name']"
        :filters="{ global: { value: search, matchMode: 'contains' } }"
        class="overflow-hidden rounded-[14px] border border-line"
      >
        <template #header>
          <div class="flex items-center justify-between gap-2.5">
            <IconField>
              <InputIcon class="pi pi-search" />
              <InputText v-model="search" placeholder="Поиск валюты…" class="w-56" />
            </IconField>
            <span class="text-xs text-ink-faint">
              Всего <span class="tnum">{{ catalog.currencies.length }}</span>
            </span>
          </div>
        </template>
        <template #empty>
          <div class="py-8 text-center text-sm text-ink-faint">Валюты не найдены</div>
        </template>

        <Column field="code" header="Код" sortable>
          <template #body="{ data }"
            ><span class="tnum font-semibold">{{ data.code }}</span></template
          >
        </Column>
        <Column field="name" header="Название" sortable />
        <Column field="kind" header="Тип" sortable>
          <template #body="{ data }">
            <Tag
              :value="KIND_META[data.kind as CurrencyKind].label"
              :severity="KIND_META[data.kind as CurrencyKind].severity"
            />
          </template>
        </Column>
        <Column field="network" header="Сеть">
          <template #body="{ data }">
            <span class="tnum text-ink-muted">{{ data.network || '—' }}</span>
          </template>
        </Column>
      </DataTable>

      <!-- mapping queue -->
      <div class="rounded-[14px] border border-line bg-surface p-4">
        <div class="text-sm font-bold">Маппинг неизвестных валют</div>
        <div class="mb-3.5 mt-0.5 text-xs text-ink-faint">
          Автопривязка по похожести · подтвердите
        </div>

        <div
          v-if="!queue.length"
          class="rounded-xl border border-line bg-well py-8 text-center text-sm text-ink-faint"
        >
          <i class="pi pi-check-circle mb-2 block text-lg text-success-text" />
          Очередь пуста
        </div>

        <div class="flex flex-col gap-2.5">
          <div
            v-for="c in queue"
            :key="c.source"
            class="rounded-xl border border-line bg-well p-3.5"
          >
            <div class="flex flex-wrap items-center gap-2.5">
              <span class="tnum text-[13px] font-semibold text-warn-text">{{ c.source }}</span>
              <UiIcon name="arrowRight" :size="15" class="text-ink-faint" />
              <span v-if="c.suggestion" class="inline-flex items-center gap-2">
                <span
                  class="flex h-[22px] w-[22px] items-center justify-center rounded-full text-[10px] font-bold text-white"
                  :style="{ background: c.suggestion.color }"
                  >{{ c.suggestion.symbol }}</span
                >
                <span class="text-[13px] font-semibold">{{ c.suggestion.name }}</span>
                <span class="tnum text-[11px] text-success-text">{{ c.confidence }}%</span>
              </span>
              <span v-else class="text-xs text-ink-faint">нет уверенного совпадения</span>
            </div>
            <div class="mt-3 flex gap-2">
              <Button
                v-if="c.suggestion"
                label="Привязать"
                size="small"
                class="flex-1"
                @click="bind(c)"
              />
              <Button
                :label="c.suggestion ? 'Другая' : 'Выбрать вручную'"
                size="small"
                severity="secondary"
                :class="c.suggestion ? '' : 'flex-1'"
                @click="
                  toast.add({
                    severity: 'info',
                    summary: 'Ручной выбор',
                    detail: c.source,
                    life: 2000,
                  })
                "
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
