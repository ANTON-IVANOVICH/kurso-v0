<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@pinia/colada'
import { useToast } from 'primevue/usetoast'
import { api } from '../lib/api'
import PageHeader from '../components/ui/PageHeader.vue'
import ExchangerAvatar from '../components/ui/ExchangerAvatar.vue'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Rating from 'primevue/rating'
import Toast from 'primevue/toast'

interface AdminReview {
  id: string
  author: string
  rating: number
  title?: string
  body: string
  status: string
  createdAt: string
  exchangerSlug: string
  exchangerName: string
}

const toast = useToast()
const { data, isLoading, refetch } = useQuery({
  key: ['admin-reviews', 'pending'],
  query: () => api.get<AdminReview[]>('/admin/reviews?status=pending'),
})
const queue = computed(() => data.value ?? [])

async function moderate(
  r: AdminReview,
  status: 'published' | 'rejected' | 'needs_info',
  label: string,
) {
  try {
    await api.patch(`/admin/reviews/${r.id}`, { status })
    await refetch()
    toast.add({ severity: 'success', summary: label, detail: r.exchangerName, life: 2200 })
  } catch {
    toast.add({ severity: 'error', summary: 'Не удалось выполнить действие', life: 2500 })
  }
}
</script>

<template>
  <div class="max-w-2xl">
    <Toast />
    <PageHeader title="Отзывы" subtitle="Очередь модерации — новые отзывы ждут проверки">
      <template #actions>
        <Tag :value="`${queue.length} в очереди`" severity="info" />
      </template>
    </PageHeader>

    <div
      v-if="isLoading"
      class="rounded-[14px] border border-line bg-surface py-14 text-center text-ink-faint"
    >
      Загрузка…
    </div>

    <div v-else-if="queue.length" class="flex flex-col gap-3">
      <div v-for="r in queue" :key="r.id" class="rounded-[14px] border border-line bg-surface p-5">
        <div class="mb-3 flex items-center gap-3">
          <ExchangerAvatar :slug="r.exchangerSlug" :name="r.author" :size="30" />
          <div class="flex-1">
            <div class="text-sm font-semibold">{{ r.author }} · {{ r.exchangerName }}</div>
            <div class="text-xs text-ink-faint">
              {{ new Date(r.createdAt).toLocaleString('ru') }}
            </div>
          </div>
          <Rating :model-value="r.rating" readonly />
        </div>
        <p v-if="r.title" class="mb-1 text-sm font-semibold text-ink">{{ r.title }}</p>
        <p class="mb-4 text-[13px] leading-relaxed text-ink-body">{{ r.body }}</p>
        <div class="flex flex-wrap gap-2.5">
          <Button
            label="Опубликовать"
            icon="pi pi-check"
            class="flex-1"
            @click="moderate(r, 'published', 'Опубликован')"
          />
          <Button
            label="Отклонить"
            severity="danger"
            outlined
            class="flex-1"
            @click="moderate(r, 'rejected', 'Отклонён')"
          />
          <Button
            label="Уточнить"
            severity="secondary"
            class="flex-1"
            @click="moderate(r, 'needs_info', 'Запрошено уточнение')"
          />
        </div>
      </div>
    </div>

    <div v-else class="rounded-[14px] border border-line bg-surface py-14 text-center">
      <i class="pi pi-check-circle mb-2 block text-2xl text-success-text" />
      <p class="text-sm text-ink-muted">Очередь модерации пуста</p>
    </div>
  </div>
</template>
