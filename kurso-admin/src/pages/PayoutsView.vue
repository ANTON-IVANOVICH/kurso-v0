<script setup lang="ts">
import { computed, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { PAYOUTS, type Payout } from '../lib/fixtures'
import { fmtRub } from '../lib/format'
import PageHeader from '../components/ui/PageHeader.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Toast from 'primevue/toast'

const toast = useToast()
const rows = ref<Payout[]>([...PAYOUTS])
const total = computed(() => rows.value.reduce((s, p) => s + p.amount, 0))

function pay(p: Payout) {
  rows.value = rows.value.filter((x) => x.recipient !== p.recipient)
  toast.add({
    severity: 'success',
    summary: `Выплата ${fmtRub(p.amount)}`,
    detail: p.recipient,
    life: 2500,
  })
}
</script>

<template>
  <div>
    <Toast />
    <PageHeader title="Выплаты" subtitle="Партнёрам и обменникам · к подтверждению">
      <template #actions>
        <Tag :value="fmtRub(total)" severity="warn" />
      </template>
    </PageHeader>

    <DataTable
      :value="rows"
      data-key="recipient"
      class="overflow-hidden rounded-[14px] border border-line"
    >
      <template #empty>
        <div class="py-10 text-center text-sm text-ink-faint">
          <i class="pi pi-check-circle mb-2 block text-lg text-success-text" />Все выплаты
          подтверждены
        </div>
      </template>

      <Column header="Получатель">
        <template #body="{ data }">
          <div>
            <div class="font-semibold">{{ data.recipient }}</div>
            <div class="text-[11px] text-ink-faint">{{ data.kind }}</div>
          </div>
        </template>
      </Column>
      <Column field="method" header="Метод">
        <template #body="{ data }"
          ><span class="tnum text-ink-muted">{{ data.method }}</span></template
        >
      </Column>
      <Column header="Сумма">
        <template #body="{ data }"
          ><span class="tnum font-semibold">{{ fmtRub(data.amount) }}</span></template
        >
      </Column>
      <Column header-style="width:7rem">
        <template #body="{ data }">
          <Button label="Выплатить" size="small" @click="pay(data)" />
        </template>
      </Column>
    </DataTable>

    <p class="mt-4 text-xs text-ink-fainter">Демо-данные до появления API выплат.</p>
  </div>
</template>
