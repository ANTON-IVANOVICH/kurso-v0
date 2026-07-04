<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { timeAgo } from '~/composables/useHistory'
import {
  useReviewsData,
  submitReview,
  reportReview,
  type ReviewItem,
} from '~/composables/useReviews'

const props = defineProps<{ slug: string; exchangerName: string }>()
const { t, plural } = useI18n()

const { data, refresh } = useReviewsData(() => props.slug)
const summary = computed(
  () => data.value?.summary ?? { average: 0, count: 0, histogram: [0, 0, 0, 0, 0] },
)
const allReviews = computed<ReviewItem[]>(() => data.value?.reviews ?? [])

// Client clock for stable relative timestamps.
const now = ref(0)
onMounted(() => (now.value = Date.now()))
const ago = (iso: string) => (now.value ? timeAgo(new Date(iso).getTime(), now.value) : '')

// Histogram rows 5★ → 1★.
const rows = computed(() =>
  [5, 4, 3, 2, 1].map((star) => {
    const n = summary.value.histogram[star - 1] ?? 0
    return { star, n, pct: summary.value.count ? Math.round((n / summary.value.count) * 100) : 0 }
  }),
)

// Filter.
const filter = ref<'all' | 'positive' | 'negative'>('all')
const filters = [
  { key: 'all', label: t('detail.rvAll') },
  { key: 'positive', label: t('detail.rvPositive') },
  { key: 'negative', label: t('detail.rvNegative') },
] as const
const shown = computed(() =>
  allReviews.value.filter((r) =>
    filter.value === 'positive'
      ? r.rating >= 4
      : filter.value === 'negative'
        ? r.rating <= 2
        : true,
  ),
)

// --- submit form ---
const showForm = ref(false)
const form = reactive({ rating: 5, author: '', body: '' })
const submitting = ref(false)
const notice = ref('')
const error = ref('')

async function submit() {
  error.value = ''
  if (form.body.trim().length < 10) {
    error.value = 'Отзыв слишком короткий (минимум 10 символов)'
    return
  }
  submitting.value = true
  try {
    const created = await submitReview(props.slug, {
      author: form.author.trim(),
      rating: form.rating,
      body: form.body.trim(),
    })
    if (created.status === 'published') {
      await refresh()
      notice.value = 'Спасибо! Ваш отзыв опубликован.'
    } else {
      notice.value = 'Отзыв отправлен на модерацию — появится после проверки.'
    }
    form.body = ''
    form.author = ''
    form.rating = 5
    showForm.value = false
  } catch {
    error.value = 'Не удалось отправить отзыв. Попробуйте позже.'
  } finally {
    submitting.value = false
  }
}

// --- report ---
const reported = reactive<Record<string, boolean>>({})
async function report(id: string) {
  if (reported[id]) return
  reported[id] = true
  try {
    await reportReview(id, 'inappropriate')
  } catch {
    /* keep the optimistic "reported" state */
  }
}

const avatar = (name: string) => exchangerAvatar(name, name)
</script>

<template>
  <div class="rounded-[20px] border border-line bg-surface p-[22px]">
    <div class="mb-5 flex items-center justify-between">
      <span class="text-[17px] font-bold">{{ t('detail.reviews') }}</span>
      <button
        type="button"
        class="rounded-[10px] bg-brand px-[18px] py-2.5 text-sm font-semibold text-white transition-colors hover:bg-brand-hover"
        @click="showForm = !showForm"
      >
        {{ t('detail.leaveReview') }}
      </button>
    </div>

    <!-- notice -->
    <div
      v-if="notice"
      class="mb-5 flex items-center gap-2.5 rounded-2xl border border-success/25 bg-success/[0.07] px-4 py-3 text-[13px] text-success-bright"
    >
      <KStatusDot tone="success" :size="7" />{{ notice }}
    </div>

    <!-- submit form -->
    <form
      v-if="showForm"
      class="mb-5 rounded-2xl border border-line bg-well p-4"
      @submit.prevent="submit"
    >
      <div class="mb-3 flex items-center gap-3">
        <span class="text-[13px] text-ink-muted">Оценка</span>
        <div class="flex gap-1">
          <button
            v-for="s in 5"
            :key="s"
            type="button"
            class="text-xl leading-none transition-colors"
            :class="s <= form.rating ? 'text-warning-deep' : 'text-line-strong'"
            :aria-label="`${s}`"
            @click="form.rating = s"
          >
            ★
          </button>
        </div>
      </div>
      <input
        v-model="form.author"
        placeholder="Имя (необязательно)"
        class="mb-2.5 w-full rounded-xl border border-line bg-surface px-3.5 py-2.5 text-sm text-ink placeholder:text-ink-faint focus:border-brand focus:outline-none"
      />
      <textarea
        v-model="form.body"
        rows="3"
        placeholder="Как прошёл обмен? Курс, скорость, поддержка…"
        class="mb-1 w-full resize-none rounded-xl border border-line bg-surface px-3.5 py-2.5 text-sm text-ink placeholder:text-ink-faint focus:border-brand focus:outline-none"
      />
      <p v-if="error" class="mb-2 text-[13px] text-danger">{{ error }}</p>
      <div class="mt-2 flex items-center gap-2.5">
        <KButton size="sm" :disabled="submitting" @click="submit">{{
          submitting ? 'Отправка…' : 'Отправить'
        }}</KButton>
        <button
          type="button"
          class="text-[13px] text-ink-faint transition-colors hover:text-ink-muted"
          @click="showForm = false"
        >
          Отмена
        </button>
      </div>
    </form>

    <!-- summary -->
    <div class="mb-5 flex items-center gap-7 border-b border-line pb-5">
      <div class="flex-none text-center">
        <div class="tnum text-[48px] font-extrabold leading-none tracking-[-0.02em]">
          {{ summary.count ? summary.average.toFixed(1) : '—' }}
        </div>
        <div class="mt-1.5 text-[15px] tracking-widest text-warning-deep">★★★★★</div>
        <div class="tnum mt-1.5 text-xs text-ink-faint">
          {{ summary.count }} {{ plural(summary.count, 'reviews') }}
        </div>
      </div>
      <div class="flex flex-1 flex-col gap-[7px]">
        <div v-for="d in rows" :key="d.star" class="flex items-center gap-2.5">
          <span class="w-3.5 text-xs text-ink-faint">{{ d.star }}</span>
          <div class="h-[7px] flex-1 overflow-hidden rounded bg-surface-raised">
            <div
              class="h-full rounded"
              :class="d.star >= 4 ? 'bg-success' : d.star === 3 ? 'bg-warning-deep' : 'bg-danger'"
              :style="{ width: `${d.pct}%` }"
            />
          </div>
          <span class="tnum w-8 text-right text-xs text-ink-faint">{{ d.pct }}%</span>
        </div>
      </div>
    </div>

    <!-- filters -->
    <div class="scrollx -mx-1 mb-[18px] flex gap-2 overflow-x-auto px-1">
      <button
        v-for="f in filters"
        :key="f.key"
        type="button"
        class="whitespace-nowrap rounded-full px-3.5 py-[7px] text-[13px] transition-colors"
        :class="
          filter === f.key
            ? 'bg-brand font-medium text-white'
            : 'border border-line-strong bg-surface-raised text-ink-muted hover:text-ink'
        "
        @click="filter = f.key"
      >
        {{ f.label }}
      </button>
    </div>

    <!-- empty -->
    <div
      v-if="!shown.length"
      class="rounded-2xl border border-dashed border-line-strong bg-well py-10 text-center text-sm text-ink-faint"
    >
      {{ allReviews.length ? 'Нет отзывов в этой категории' : 'Пока нет отзывов — станьте первым' }}
    </div>

    <!-- review cards -->
    <div
      v-for="(rv, i) in shown"
      :key="rv.id"
      class="border-line pb-[18px]"
      :class="i < shown.length - 1 ? 'mb-[18px] border-b' : ''"
    >
      <div class="mb-2.5 flex items-center gap-3">
        <span
          class="flex h-10 w-10 flex-none items-center justify-center rounded-full text-sm font-bold text-white"
          :style="{ background: avatar(rv.author).color }"
          >{{ avatar(rv.author).initials }}</span
        >
        <div class="min-w-0 flex-1">
          <div class="text-[15px] font-semibold">{{ rv.author }}</div>
          <div class="text-xs text-ink-faint">{{ ago(rv.createdAt) }}</div>
        </div>
        <span class="text-sm tracking-widest text-warning-deep"
          >{{ '★'.repeat(rv.rating)
          }}<span class="text-line-strong">{{ '★'.repeat(5 - rv.rating) }}</span></span
        >
      </div>
      <p class="mb-2 text-sm leading-relaxed text-ink-bright">{{ rv.body }}</p>
      <button
        type="button"
        class="text-xs transition-colors"
        :class="reported[rv.id] ? 'text-ink-faint' : 'text-ink-faint hover:text-danger'"
        :disabled="reported[rv.id]"
        @click="report(rv.id)"
      >
        {{ reported[rv.id] ? 'Жалоба отправлена' : 'Пожаловаться' }}
      </button>
    </div>
  </div>
</template>
