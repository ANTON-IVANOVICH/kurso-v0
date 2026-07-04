<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

// Affiliate cabinet: a real, unique referral link (clicks are attributed through
// the kurso_ref cookie → clickouts.ref_code) with live click/registration stats.
// Earnings are an honest estimate from real clicks until the confirmed-exchange
// revshare pipeline lands.
definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Партнёрка — Kurso' })

const { data, pending, error, load } = usePartner()
const origin = useRequestURL().origin

const PRESET_TAGS = ['youtube', 'telegram', 'vc']
const activeTag = ref<string>('')
const customTags = ref<string[]>([])
const addingTag = ref(false)
const newTag = ref('')
const copied = ref(false)

const link = computed(() => {
  const code = data.value?.code ?? '…'
  return `${origin}/?ref=${activeTag.value ? `${code}.${activeTag.value}` : code}`
})
const allTags = computed(() => [...PRESET_TAGS, ...customTags.value])

async function copyLink() {
  try {
    await navigator.clipboard.writeText(link.value)
    copied.value = true
    setTimeout(() => (copied.value = false), 1600)
  } catch {
    /* clipboard blocked — no-op */
  }
}
function addTag() {
  const t = newTag.value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9_-]/g, '')
  if (t && !allTags.value.includes(t)) customTags.value.push(t)
  newTag.value = ''
  addingTag.value = false
  if (t) activeTag.value = t
}

// clicks-by-tag lookup for the "sources" panel
const maxTagClicks = computed(() => Math.max(1, ...(data.value?.byTag ?? []).map((t) => t.clicks)))
const maxDay = computed(() => Math.max(1, ...(data.value?.series ?? []).map((s) => s.clicks)))
function tagLabel(tag: string) {
  return tag ? `ref / ${tag}` : 'прямая ссылка'
}

onMounted(load)
</script>

<template>
  <div>
    <div class="mb-5 flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Партнёрская программа</h1>
        <p class="mt-1 text-sm text-ink-faint">Зарабатывайте на переходах к обменникам</p>
      </div>
      <span
        class="inline-flex items-center gap-2 rounded-[10px] border border-line bg-surface px-3.5 py-2 text-[13px] text-ink-muted"
      >
        <span class="tnum font-bold text-ink">{{ data?.revsharePct ?? 30 }}%</span> revshare ·
        cookie <span class="tnum font-bold text-ink">{{ data?.cookieDays ?? 90 }}</span> дней
      </span>
    </div>

    <!-- link + balance -->
    <div class="mb-4 grid gap-4 lg:grid-cols-[1.6fr_1fr]">
      <div class="rounded-2xl border border-line bg-surface p-[18px]">
        <div class="mb-2.5 text-[13px] text-ink-faint">Ваша партнёрская ссылка</div>
        <div class="flex flex-wrap items-center gap-2.5">
          <div
            class="tnum min-w-0 flex-1 overflow-hidden text-ellipsis whitespace-nowrap rounded-xl border border-line bg-well px-[15px] py-3 text-sm text-ink-body"
          >
            {{ link }}
          </div>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-xl bg-brand px-4 py-3 text-[13px] font-semibold text-white transition-colors hover:bg-brand-hover"
            @click="copyLink"
          >
            <svg
              width="15"
              height="15"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <rect x="9" y="9" width="11" height="11" rx="2" />
              <path d="M5 15V5a2 2 0 0 1 2-2h10" />
            </svg>
            {{ copied ? 'Скопировано' : 'Копировать' }}
          </button>
          <NuxtLink
            to="/account/widgets"
            class="rounded-xl border border-line-strong bg-surface-raised px-4 py-3 text-[13px] font-semibold text-ink transition-colors hover:border-[#3A4047]"
            >Виджеты</NuxtLink
          >
        </div>
        <div class="mt-3 flex flex-wrap gap-2">
          <button
            type="button"
            class="rounded-lg border px-2.5 py-1.5 text-xs transition-colors"
            :class="
              activeTag === ''
                ? 'border-brand bg-brand/10 text-brand-bright'
                : 'border-line bg-well text-ink-muted hover:text-ink'
            "
            @click="activeTag = ''"
          >
            без метки
          </button>
          <button
            v-for="t in allTags"
            :key="t"
            type="button"
            class="rounded-lg border px-2.5 py-1.5 text-xs transition-colors"
            :class="
              activeTag === t
                ? 'border-brand bg-brand/10 text-brand-bright'
                : 'border-line bg-well text-ink-muted hover:text-ink'
            "
            @click="activeTag = t"
          >
            UTM: {{ t }}
          </button>
          <template v-if="addingTag">
            <input
              v-model="newTag"
              placeholder="метка"
              class="w-24 rounded-lg border border-line-strong bg-well px-2.5 py-1.5 text-xs text-ink focus:border-brand focus:outline-none"
              @keyup.enter="addTag"
              @blur="addTag"
            />
          </template>
          <button
            v-else
            type="button"
            class="rounded-lg border border-dashed border-line-strong bg-well px-2.5 py-1.5 text-xs text-brand-bright"
            @click="addingTag = true"
          >
            + метка
          </button>
        </div>
      </div>

      <div
        class="flex flex-col rounded-2xl border border-brand/30 bg-[linear-gradient(160deg,rgba(46,125,242,0.16),rgba(46,125,242,0.04))] p-[18px]"
      >
        <div class="mb-2 text-[13px] text-ink-muted">Накоплено (оценка)</div>
        <div class="tnum text-[34px] font-extrabold tracking-[-0.02em] text-ink">
          ₽ {{ fmtNumber(data?.estimatedRub ?? 0, 0) }}
        </div>
        <div class="mt-1 text-xs text-ink-faint">
          оценка по переходам · выплаты — после подтверждённых обменов
        </div>
        <button
          type="button"
          disabled
          class="mt-auto cursor-not-allowed rounded-xl bg-brand/40 py-3 text-[15px] font-bold text-white/70"
          title="Доступно после первых подтверждённых обменов"
        >
          Запросить выплату
        </button>
      </div>
    </div>

    <!-- metrics -->
    <div class="mb-4 grid grid-cols-2 gap-3.5 lg:grid-cols-4">
      <div class="rounded-2xl border border-line bg-surface p-4">
        <div class="mb-2 text-[13px] text-ink-muted">Клики</div>
        <div class="tnum text-2xl font-bold text-ink">{{ fmtNumber(data?.clicks ?? 0, 0) }}</div>
        <div class="mt-1 text-xs text-ink-faint">за 30 дней</div>
      </div>
      <div class="rounded-2xl border border-line bg-surface p-4">
        <div class="mb-2 text-[13px] text-ink-muted">Регистрации</div>
        <div class="tnum text-2xl font-bold text-ink">
          {{ fmtNumber(data?.registrations ?? 0, 0) }}
        </div>
        <div class="mt-1 text-xs text-ink-faint">по вашей ссылке</div>
      </div>
      <div class="rounded-2xl border border-line bg-surface p-4">
        <div class="mb-2 text-[13px] text-ink-muted">Обмены</div>
        <div class="tnum text-2xl font-bold text-ink-faint">—</div>
        <div class="mt-1 text-xs text-ink-faint">учёт скоро</div>
      </div>
      <div class="rounded-2xl border border-line bg-surface p-4">
        <div class="mb-2 text-[13px] text-ink-muted">Заработано (оценка)</div>
        <div class="tnum text-2xl font-bold text-ink">
          ₽{{ fmtNumber(data?.estimatedRub ?? 0, 0) }}
        </div>
        <div class="mt-1 text-xs text-success-bright">по переходам</div>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-[1.5fr_1fr]">
      <!-- clicks chart -->
      <div class="rounded-2xl border border-line bg-surface p-[18px]">
        <div class="mb-4 flex items-center justify-between">
          <span class="text-[15px] font-bold text-ink">Клики по дням</span>
          <span class="text-xs text-ink-faint">30 дней</span>
        </div>
        <div v-if="pending" class="py-10 text-center text-sm text-ink-faint">Загрузка…</div>
        <div v-else-if="error" class="py-10 text-center text-sm text-ink-faint">
          Не удалось загрузить статистику
        </div>
        <template v-else>
          <div class="flex h-[130px] items-end gap-[3px]">
            <div
              v-for="(s, i) in data?.series ?? []"
              :key="s.day"
              class="flex-1 rounded-t-[3px]"
              :class="
                i === (data?.series.length ?? 0) - 1
                  ? 'bg-gradient-to-b from-brand-bright to-brand'
                  : 'bg-line-strong'
              "
              :style="{ height: `${Math.max(3, (s.clicks / maxDay) * 100)}%` }"
              :title="`${s.day}: ${s.clicks}`"
            />
          </div>
          <p v-if="!data?.clicks" class="mt-3 text-center text-xs text-ink-faint">
            Пока нет переходов — поделитесь ссылкой, чтобы начать.
          </p>
        </template>
      </div>

      <!-- sources by UTM -->
      <div class="rounded-2xl border border-line bg-surface p-[18px]">
        <div class="mb-3.5 text-[15px] font-bold text-ink">Источники · UTM</div>
        <div v-if="(data?.byTag ?? []).length" class="flex flex-col gap-3">
          <div v-for="t in data?.byTag ?? []" :key="t.tag" class="flex items-center gap-2.5">
            <span class="flex-1 truncate text-[13px] text-ink">{{ tagLabel(t.tag) }}</span>
            <div class="h-[7px] w-[90px] overflow-hidden rounded-full bg-line">
              <div
                class="h-full bg-brand"
                :style="{ width: `${(t.clicks / maxTagClicks) * 100}%` }"
              />
            </div>
            <span class="tnum w-10 text-right text-xs text-ink-muted">{{ t.clicks }}</span>
          </div>
        </div>
        <p v-else class="py-6 text-center text-sm text-ink-faint">
          Добавьте UTM-метку к ссылке, чтобы видеть источники.
        </p>
      </div>
    </div>
  </div>
</template>
