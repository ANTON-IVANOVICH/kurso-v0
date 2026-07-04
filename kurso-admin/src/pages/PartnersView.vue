<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCatalog } from '../composables/useCatalog'
import type { Exchanger } from '../types/models'
import { toNum, fmtCompact } from '../lib/format'
import PageHeader from '../components/ui/PageHeader.vue'
import ExchangerAvatar from '../components/ui/ExchangerAvatar.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Message from 'primevue/message'

const catalog = useCatalog()
const router = useRouter()

const partners = computed(() => catalog.exchangers.filter((e) => e.partner))

function openEditor(e: Exchanger) {
  router.push({ name: 'exchanger-edit', params: { slug: e.slug } })
}
</script>

<template>
  <div>
    <PageHeader title="Партнёрка" subtitle="Обменники с реферальным соглашением">
      <template #actions>
        <Tag :value="`${partners.length} партнёров`" severity="info" />
      </template>
    </PageHeader>

    <Message severity="secondary" class="mb-4" size="small">
      Партнёрские метрики и revshare-выплаты появятся с партнёрской программой. Ниже — живой список
      обменников-партнёров из каталога.
    </Message>

    <DataTable
      :value="partners"
      :loading="catalog.loading"
      data-key="id"
      class="overflow-hidden rounded-[14px] border border-line"
    >
      <template #empty>
        <div class="py-10 text-center text-sm text-ink-faint">Партнёры не найдены</div>
      </template>

      <Column header="Обменник">
        <template #body="{ data }">
          <button class="flex items-center gap-2.5 text-left" @click="openEditor(data)">
            <ExchangerAvatar :slug="data.slug" :name="data.name" :size="28" />
            <span class="font-semibold">{{ data.name }}</span>
          </button>
        </template>
      </Column>
      <Column field="directionsCount" header="Направления" sortable>
        <template #body="{ data }"
          ><span class="tnum text-ink-muted">{{ data.directionsCount }}</span></template
        >
      </Column>
      <Column header="Резерв">
        <template #body="{ data }">
          <span class="tnum text-ink-muted">{{
            data.reserveTotal ? fmtCompact(toNum(data.reserveTotal)) : '—'
          }}</span>
        </template>
      </Column>
      <Column field="ratingAvg" header="Рейтинг" sortable>
        <template #body="{ data }">
          <span class="tnum">{{ data.ratingAvg != null ? data.ratingAvg.toFixed(1) : '—' }}</span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>
