<script setup lang="ts">
import { computed, ref } from 'vue'
import { useComplaintsQuery } from '../composables/useMerchant'
import { timeAgo } from '../lib/format'
import PageHeader from '../components/ui/PageHeader.vue'
import Tag from 'primevue/tag'
import Rating from 'primevue/rating'

const complaints = useComplaintsQuery()

type Filter = 'open' | 'all'
const filter = ref<Filter>('open')

const list = computed(() => {
  const all = complaints.data.value ?? []
  return filter.value === 'open' ? all.filter((c) => c.status === 'open') : all
})
const openCount = computed(
  () => (complaints.data.value ?? []).filter((c) => c.status === 'open').length,
)
const resolvedCount = computed(
  () => (complaints.data.value ?? []).filter((c) => c.status !== 'open').length,
)

const STATUS: Record<string, { label: string; severity: string }> = {
  open: { label: 'открыта', severity: 'danger' },
  reviewed: { label: 'рассмотрена', severity: 'success' },
  dismissed: { label: 'отклонена', severity: 'secondary' },
}

const REASONS: Record<string, string> = {
  spam: 'Спам',
  fake: 'Поддельный отзыв',
  offensive: 'Оскорбления',
  other: 'Другое',
}
function reasonLabel(r: string) {
  return REASONS[r] ?? r
}
</script>

<template>
  <div>
    <PageHeader title="Жалобы" subtitle="Жалобы пользователей на отзывы о вашем обменнике">
      <template #actions>
        <div class="flex rounded-[9px] border border-line bg-well p-0.5 text-xs">
          <button
            class="rounded-md px-3 py-1.5 font-medium transition-colors"
            :class="filter === 'open' ? 'bg-brand text-white' : 'text-ink-muted hover:text-ink'"
            @click="filter = 'open'"
          >
            Открытые <span class="tnum ml-1 opacity-70">{{ openCount }}</span>
          </button>
          <button
            class="rounded-md px-3 py-1.5 font-medium transition-colors"
            :class="filter === 'all' ? 'bg-brand text-white' : 'text-ink-muted hover:text-ink'"
            @click="filter = 'all'"
          >
            Все
          </button>
        </div>
      </template>
    </PageHeader>

    <div v-if="complaints.isLoading.value" class="py-10 text-center text-sm text-ink-faint">
      Загрузка…
    </div>
    <div
      v-else-if="!list.length"
      class="rounded-[14px] border border-line bg-surface py-12 text-center text-sm text-ink-faint"
    >
      {{ filter === 'open' ? 'Открытых жалоб нет.' : 'Жалоб пока нет.' }}
    </div>

    <div v-else class="overflow-hidden rounded-[14px] border border-line bg-surface">
      <div v-for="c in list" :key="c.id" class="border-b border-line p-[18px] last:border-b-0">
        <div class="mb-3 flex flex-wrap items-center gap-2.5">
          <Tag :value="reasonLabel(c.reason)" severity="warn" />
          <span class="tnum text-[11px] text-ink-faint">{{ timeAgo(c.createdAt) }}</span>
          <Tag
            class="ml-auto"
            :value="STATUS[c.status].label"
            :severity="STATUS[c.status].severity"
          />
        </div>

        <p v-if="c.details" class="mb-3 text-[13px] leading-relaxed text-ink-body">
          {{ c.details }}
        </p>

        <!-- the reported review -->
        <div class="rounded-xl border border-line bg-well p-3.5">
          <div class="mb-1.5 flex items-center gap-2">
            <span class="text-[12px] font-semibold text-ink-muted">{{ c.author }}</span>
            <Rating :model-value="c.rating" readonly />
            <span class="ml-auto text-[11px] text-ink-faint">отзыв под жалобой</span>
          </div>
          <p class="text-[13px] leading-relaxed text-ink-body">{{ c.reviewBody }}</p>
        </div>
      </div>
    </div>

    <div
      v-if="!complaints.isLoading.value && resolvedCount"
      class="mt-3 flex items-center gap-2.5 rounded-xl border border-line bg-well px-4 py-3.5"
    >
      <i class="pi pi-check-circle text-success-text" />
      <span class="text-[13px] text-ink-muted">
        Рассмотрено модерацией: <span class="tnum font-semibold text-ink">{{ resolvedCount }}</span>
      </span>
    </div>
  </div>
</template>
