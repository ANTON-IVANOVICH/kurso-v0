<script setup lang="ts">
import { computed, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { COMPLAINTS, type Complaint } from '../lib/fixtures'
import PageHeader from '../components/ui/PageHeader.vue'
import ExchangerAvatar from '../components/ui/ExchangerAvatar.vue'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Toast from 'primevue/toast'

const toast = useToast()

const queue = ref<Complaint[]>([...COMPLAINTS])
const total = COMPLAINTS.length
const done = ref(0)
const current = computed(() => queue.value[0] ?? null)
const position = computed(() => Math.min(done.value + 1, total))

function resolve(verb: string, severity: 'success' | 'warn' | 'danger') {
  const c = current.value
  if (!c) return
  toast.add({ severity, summary: `${verb}`, detail: `#${c.id} · ${c.exchanger}`, life: 2200 })
  queue.value.shift()
  done.value += 1
}
</script>

<template>
  <div class="max-w-2xl">
    <Toast />
    <PageHeader title="Жалобы" subtitle="Арбитраж споров клиент ↔ обменник">
      <template #actions>
        <Tag :value="`${total} открыты`" severity="danger" />
      </template>
    </PageHeader>

    <div v-if="current" class="rounded-[14px] border border-line bg-surface p-5">
      <div class="mb-3 flex items-center gap-3">
        <ExchangerAvatar :slug="current.client" :name="current.client" :size="30" />
        <div class="flex-1">
          <div class="text-sm font-semibold">{{ current.client }} → {{ current.exchanger }}</div>
          <div class="tnum text-xs text-ink-faint">
            {{ current.direction }} · {{ current.amount }} · #{{ current.id }} · {{ current.ago }}
          </div>
        </div>
        <Tag :value="current.tag" severity="danger" />
      </div>

      <p class="mb-3 text-[13px] leading-relaxed text-ink-body">{{ current.body }}</p>

      <div
        v-if="current.reply"
        class="mb-4 flex items-center gap-2.5 rounded-[10px] border border-line bg-well px-3 py-2.5"
      >
        <ExchangerAvatar
          :slug="current.exchanger.toLowerCase()"
          :name="current.exchanger"
          :size="18"
          shape="rounded"
        />
        <span class="text-xs leading-snug text-ink-muted">Обменник: «{{ current.reply }}»</span>
      </div>

      <div class="flex flex-wrap gap-2.5">
        <Button
          label="В пользу клиента"
          icon="pi pi-check"
          class="flex-1"
          @click="resolve('Решено в пользу клиента', 'success')"
        />
        <Button
          label="Запросить лог"
          severity="secondary"
          class="flex-1"
          @click="resolve('Запрошен лог', 'warn')"
        />
        <Button
          label="Отклонить"
          severity="danger"
          outlined
          class="flex-1"
          @click="resolve('Отклонена', 'danger')"
        />
      </div>

      <div class="mt-3 text-center text-xs text-ink-faint">
        <span class="tnum">{{ position }}</span> из <span class="tnum">{{ total }}</span>
      </div>
    </div>

    <div v-else class="rounded-[14px] border border-line bg-surface py-14 text-center">
      <i class="pi pi-check-circle mb-2 block text-2xl text-success-text" />
      <p class="text-sm text-ink-muted">Все споры разобраны</p>
    </div>

    <p class="mt-4 text-xs text-ink-fainter">Демо-очередь до появления API жалоб.</p>
  </div>
</template>
