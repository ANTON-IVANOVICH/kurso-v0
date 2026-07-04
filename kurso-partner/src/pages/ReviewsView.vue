<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useReviewsQuery, useDashboardQuery } from '../composables/useMerchant'
import { api, ApiError } from '../lib/api'
import { timeAgo } from '../lib/format'
import type { MerchantReview } from '../types/models'
import PageHeader from '../components/ui/PageHeader.vue'
import Button from 'primevue/button'
import Rating from 'primevue/rating'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'

const reviews = useReviewsQuery()
const dashboard = useDashboardQuery()
const toast = useToast()

type Filter = 'all' | 'unanswered' | 'answered'
const filter = ref<Filter>('unanswered')

const list = computed(() => {
  const all = reviews.data.value ?? []
  if (filter.value === 'unanswered') return all.filter((r) => !r.reply)
  if (filter.value === 'answered') return all.filter((r) => r.reply)
  return all
})
const counts = computed(() => {
  const all = reviews.data.value ?? []
  return {
    all: all.length,
    unanswered: all.filter((r) => !r.reply).length,
    answered: all.filter((r) => r.reply).length,
  }
})

// Per-review draft state (editor open + text + busy).
const drafts = reactive<Record<string, { text: string; busy: boolean; error: string }>>({})
function draft(id: string, initial = '') {
  if (!drafts[id]) drafts[id] = { text: initial, busy: false, error: '' }
  return drafts[id]
}
const editing = ref<Set<string>>(new Set())
function openEditor(r: MerchantReview) {
  draft(r.id, r.reply ?? '').text = r.reply ?? ''
  editing.value = new Set(editing.value).add(r.id)
}
function closeEditor(id: string) {
  const next = new Set(editing.value)
  next.delete(id)
  editing.value = next
}

async function submit(r: MerchantReview) {
  const d = draft(r.id)
  d.error = ''
  const body = d.text.trim()
  if (body.length < 2) {
    d.error = 'Ответ слишком короткий'
    return
  }
  d.busy = true
  try {
    await api.post(`/partner/reviews/${r.id}/reply`, { body })
    await Promise.all([reviews.refetch(), dashboard.refetch()])
    closeEditor(r.id)
    toast.add({ severity: 'success', summary: 'Ответ опубликован', life: 2500 })
  } catch (e) {
    d.error = e instanceof ApiError ? e.message : 'Не удалось опубликовать'
  } finally {
    d.busy = false
  }
}

const FILTERS: { key: Filter; label: string }[] = [
  { key: 'unanswered', label: 'Без ответа' },
  { key: 'answered', label: 'Отвеченные' },
  { key: 'all', label: 'Все' },
]
</script>

<template>
  <div>
    <Toast position="bottom-right" />
    <PageHeader title="Отзывы" subtitle="Отвечайте на отзывы клиентов — один ответ на отзыв">
      <template #actions>
        <div class="flex rounded-[9px] border border-line bg-well p-0.5 text-xs">
          <button
            v-for="f in FILTERS"
            :key="f.key"
            class="rounded-md px-3 py-1.5 font-medium transition-colors"
            :class="filter === f.key ? 'bg-brand text-white' : 'text-ink-muted hover:text-ink'"
            @click="filter = f.key"
          >
            {{ f.label }}
            <span class="tnum ml-1 opacity-70">{{ counts[f.key] }}</span>
          </button>
        </div>
      </template>
    </PageHeader>

    <div v-if="reviews.isLoading.value" class="py-10 text-center text-sm text-ink-faint">
      Загрузка…
    </div>
    <div
      v-else-if="!list.length"
      class="rounded-[14px] border border-line bg-surface py-12 text-center text-sm text-ink-faint"
    >
      {{ filter === 'unanswered' ? 'Все отзывы отвечены.' : 'Отзывов пока нет.' }}
    </div>

    <div v-else class="flex flex-col gap-3">
      <div
        v-for="r in list"
        :key="r.id"
        class="rounded-[14px] border border-line bg-surface p-[18px]"
      >
        <div class="mb-2 flex flex-wrap items-center gap-2.5">
          <span
            class="flex h-7 w-7 flex-none items-center justify-center rounded-full bg-[#5B3FA0] text-[11px] font-bold text-white"
            >{{ r.author.slice(0, 1).toUpperCase() }}</span
          >
          <span class="text-[13px] font-semibold">{{ r.author }}</span>
          <Rating :model-value="r.rating" readonly />
          <span class="tnum ml-auto text-[11px] text-ink-faint">{{ timeAgo(r.createdAt) }}</span>
        </div>
        <p class="text-[13px] leading-relaxed text-ink-body">{{ r.body }}</p>

        <!-- existing reply -->
        <div
          v-if="r.reply && !editing.has(r.id)"
          class="mt-3 rounded-xl border border-line bg-well p-3.5"
        >
          <div class="mb-1.5 flex items-center gap-2">
            <i class="pi pi-reply text-xs text-brand-light" />
            <span class="text-[12px] font-semibold text-brand-light">Ваш ответ</span>
            <span class="tnum ml-auto text-[11px] text-ink-faint">{{ timeAgo(r.replyAt) }}</span>
          </div>
          <p class="text-[13px] leading-relaxed text-ink-body">{{ r.reply }}</p>
          <button class="mt-2 text-[12px] font-semibold text-brand-light" @click="openEditor(r)">
            Редактировать
          </button>
        </div>

        <!-- composer -->
        <div v-else-if="editing.has(r.id)" class="mt-3">
          <Textarea
            v-model="draft(r.id).text"
            rows="3"
            auto-resize
            maxlength="600"
            placeholder="Напишите ответ клиенту…"
            class="w-full"
          />
          <Message v-if="draft(r.id).error" severity="error" size="small" class="mt-2">{{
            draft(r.id).error
          }}</Message>
          <div class="mt-2.5 flex items-center justify-between">
            <span class="tnum text-[12px] text-ink-faint">{{ draft(r.id).text.length }} / 600</span>
            <div class="flex gap-2">
              <Button
                label="Отмена"
                size="small"
                severity="secondary"
                text
                @click="closeEditor(r.id)"
              />
              <Button
                label="Опубликовать ответ"
                size="small"
                :loading="draft(r.id).busy"
                @click="submit(r)"
              />
            </div>
          </div>
        </div>

        <div v-else class="mt-3">
          <Button label="Ответить" size="small" @click="openEditor(r)" />
        </div>
      </div>
    </div>
  </div>
</template>
