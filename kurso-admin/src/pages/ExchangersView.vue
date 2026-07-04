<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCatalog } from '../composables/useCatalog'
import type { Exchanger, ExchangerStatus } from '../types/models'
import { toNum, fmtCompact } from '../lib/format'
import PageHeader from '../components/ui/PageHeader.vue'
import ExchangerAvatar from '../components/ui/ExchangerAvatar.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Select from 'primevue/select'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import Rating from 'primevue/rating'
import type { MenuItem } from 'primevue/menuitem'

const catalog = useCatalog()
const router = useRouter()

const search = ref('')
const statusFilter = ref<ExchangerStatus | null>(null)
const selected = ref<Exchanger[]>([])

const statusOptions = [
  { label: 'Все статусы', value: null },
  { label: 'Активные', value: 'active' },
  { label: 'На паузе', value: 'paused' },
  { label: 'Забанены', value: 'banned' },
]

const STATUS_META: Record<ExchangerStatus, { label: string; severity: string }> = {
  active: { label: 'активен', severity: 'success' },
  paused: { label: 'на паузе', severity: 'warn' },
  banned: { label: 'забанен', severity: 'danger' },
}

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  return catalog.exchangers.filter((e) => {
    if (statusFilter.value && e.status !== statusFilter.value) return false
    if (q && !e.name.toLowerCase().includes(q) && !e.slug.toLowerCase().includes(q)) return false
    return true
  })
})

function openEditor(e: Exchanger) {
  router.push({ name: 'exchanger-edit', params: { slug: e.slug } })
}

// Row action menu.
const menu = ref<InstanceType<typeof Menu> | null>(null)
const menuRow = ref<Exchanger | null>(null)
const menuItems = computed<MenuItem[]>(() => {
  const e = menuRow.value
  if (!e) return []
  const items: MenuItem[] = [
    { label: 'Редактировать', icon: 'pi pi-pencil', command: () => openEditor(e) },
  ]
  if (e.status !== 'active')
    items.push({
      label: 'Активировать',
      icon: 'pi pi-check-circle',
      command: () => catalog.patchExchanger(e.slug, { status: 'active' }),
    })
  if (e.status !== 'paused')
    items.push({
      label: 'Поставить на паузу',
      icon: 'pi pi-pause',
      command: () => catalog.patchExchanger(e.slug, { status: 'paused' }),
    })
  if (e.status !== 'banned')
    items.push({
      label: 'Забанить',
      icon: 'pi pi-ban',
      command: () => catalog.patchExchanger(e.slug, { status: 'banned' }),
    })
  return items
})
function toggleMenu(event: Event, row: Exchanger) {
  menuRow.value = row
  menu.value?.toggle(event)
}

function reserveText(e: Exchanger): string {
  return e.reserveTotal ? fmtCompact(toNum(e.reserveTotal)) : '—'
}
</script>

<template>
  <div>
    <PageHeader title="Обменники" subtitle="Каталог обменников · управление статусами и профилями">
      <template #actions>
        <Button
          label="Обменник"
          icon="pi pi-plus"
          size="small"
          @click="router.push({ name: 'exchanger-new' })"
        />
      </template>
    </PageHeader>

    <DataTable
      v-model:selection="selected"
      :value="filtered"
      :loading="catalog.loading"
      data-key="id"
      paginator
      :rows="10"
      :rows-per-page-options="[10, 20, 50]"
      removable-sort
      class="overflow-hidden rounded-[14px] border border-line"
    >
      <template #header>
        <div class="flex flex-wrap items-center gap-2.5">
          <IconField>
            <InputIcon class="pi pi-search" />
            <InputText v-model="search" placeholder="Поиск обменника…" class="w-60" />
          </IconField>
          <Select
            v-model="statusFilter"
            :options="statusOptions"
            option-label="label"
            option-value="value"
            class="w-44"
          />
          <span class="ml-auto text-xs text-ink-faint">
            Показано <span class="tnum">{{ filtered.length }}</span> из
            <span class="tnum">{{ catalog.exchangers.length }}</span>
          </span>
        </div>
      </template>

      <template #empty>
        <div class="py-10 text-center text-sm text-ink-faint">
          {{ catalog.error ? catalog.error : 'Обменники не найдены' }}
        </div>
      </template>

      <Column selection-mode="multiple" header-style="width:3rem" />

      <Column field="name" header="Обменник" sortable>
        <template #body="{ data }">
          <button class="flex items-center gap-2.5 text-left" @click="openEditor(data)">
            <ExchangerAvatar :slug="data.slug" :name="data.name" :size="28" />
            <span class="font-semibold">{{ data.name }}</span>
            <Tag v-if="data.partner" value="Партнёр" severity="info" />
          </button>
        </template>
      </Column>

      <Column field="directionsCount" header="Напр." sortable>
        <template #body="{ data }"
          ><span class="tnum text-ink-muted">{{ data.directionsCount }}</span></template
        >
      </Column>

      <Column header="Резерв">
        <template #body="{ data }"
          ><span class="tnum text-ink-muted">{{ reserveText(data) }}</span></template
        >
      </Column>

      <Column field="ratingAvg" header="Рейтинг" sortable>
        <template #body="{ data }">
          <span v-if="data.ratingAvg != null" class="inline-flex items-center gap-1.5">
            <Rating :model-value="Math.round(data.ratingAvg)" readonly />
            <span class="tnum text-xs text-ink-faint">{{ data.ratingAvg.toFixed(1) }}</span>
          </span>
          <span v-else class="text-ink-faint">—</span>
        </template>
      </Column>

      <Column field="status" header="Статус" sortable>
        <template #body="{ data }">
          <Tag
            :value="STATUS_META[data.status as ExchangerStatus].label"
            :severity="STATUS_META[data.status as ExchangerStatus].severity"
          />
        </template>
      </Column>

      <Column header-style="width:3rem">
        <template #body="{ data }">
          <Button
            icon="pi pi-ellipsis-h"
            severity="secondary"
            text
            rounded
            aria-label="Действия"
            @click="toggleMenu($event, data)"
          />
        </template>
      </Column>
    </DataTable>

    <Menu ref="menu" :model="menuItems" popup />
  </div>
</template>
