<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { MOD_REVIEWS, MOD_TOTAL, type ModReview, type HealthTone } from '../lib/fixtures'
import PageHeader from '../components/ui/PageHeader.vue'
import ExchangerAvatar from '../components/ui/ExchangerAvatar.vue'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Rating from 'primevue/rating'
import Toast from 'primevue/toast'

const toast = useToast()

const queue = ref<ModReview[]>([...MOD_REVIEWS])
const done = ref(0)
const current = computed(() => queue.value[0] ?? null)
const position = computed(() => Math.min(done.value + 1, MOD_TOTAL))

const TAG_SEVERITY: Record<HealthTone, string> = {
  success: 'success',
  warn: 'warn',
  danger: 'danger',
  brand: 'info',
  muted: 'secondary',
}

function resolve(verb: string, severity: 'success' | 'warn' | 'danger') {
  const r = current.value
  if (!r) return
  toast.add({ severity, summary: `${verb}: ${r.id}`, detail: r.exchanger, life: 2200 })
  queue.value.shift()
  done.value += 1
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Enter' && current.value) {
    e.preventDefault()
    resolve('Опубликован', 'success')
  }
}
onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
</script>

<template>
  <div class="max-w-2xl">
    <Toast />
    <PageHeader title="Отзывы" subtitle="Очередь модерации">
      <template #actions>
        <Tag :value="`${MOD_TOTAL} в очереди`" severity="info" />
      </template>
    </PageHeader>

    <div v-if="current" class="rounded-[14px] border border-line bg-surface p-5">
      <div class="mb-3 flex items-center gap-3">
        <ExchangerAvatar
          :slug="current.exchanger.toLowerCase()"
          :name="current.author"
          :size="30"
        />
        <div class="flex-1">
          <div class="text-sm font-semibold">{{ current.author }} · {{ current.exchanger }}</div>
          <div class="text-xs text-ink-faint">
            {{ current.direction }} · {{ current.amount }} · {{ current.ago }}
          </div>
        </div>
        <Rating :model-value="current.rating" readonly />
      </div>

      <p class="mb-3 text-[13px] leading-relaxed text-ink-body">{{ current.body }}</p>

      <div class="mb-4 flex flex-wrap gap-2">
        <Tag
          v-for="(t, i) in current.tags"
          :key="i"
          :value="t.label"
          :severity="TAG_SEVERITY[t.tone]"
        />
      </div>

      <div class="flex flex-wrap gap-2.5">
        <Button
          label="Опубликовать"
          icon="pi pi-check"
          class="flex-1"
          @click="resolve('Опубликован', 'success')"
        />
        <Button
          label="Отклонить"
          severity="danger"
          outlined
          class="flex-1"
          @click="resolve('Отклонён', 'danger')"
        />
        <Button
          label="Уточнить"
          severity="secondary"
          class="flex-1"
          @click="resolve('Запрошено уточнение', 'warn')"
        />
      </div>

      <div class="mt-3 text-center text-xs text-ink-faint">
        <span class="tnum">{{ position }}</span> из <span class="tnum">{{ MOD_TOTAL }}</span> ·
        следующий по Enter
      </div>
    </div>

    <div v-else class="rounded-[14px] border border-line bg-surface py-14 text-center">
      <i class="pi pi-check-circle mb-2 block text-2xl text-success-text" />
      <p class="text-sm text-ink-muted">Демо-очередь разобрана</p>
    </div>

    <p class="mt-4 text-xs text-ink-fainter">Демо-очередь до появления API отзывов и модерации.</p>
  </div>
</template>
