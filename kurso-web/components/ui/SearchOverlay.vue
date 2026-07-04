<script setup lang="ts">
import { ref } from 'vue'

// Full-screen site search (opened from the mobile header lupe). Scaffold: the
// input filters nothing yet — results below are representative placeholders.
const { open, close } = useSearchOverlay()
useScrollLock(open)

const query = ref('USDT')

const directions = ['USDT → Тинькофф', 'USDT → Сбер', 'USDT → Нал.']

const results = [
  {
    code: 'CB',
    name: 'CryptoBridge',
    partner: true,
    rating: '4.9',
    reviews: '1203',
    rate: '81.20 ₽',
    color: '#3A4452',
  },
  {
    code: 'N',
    name: 'NetEx24',
    partner: false,
    rating: '4.9',
    reviews: '2104',
    rate: '80.95 ₽',
    color: '#5B3FA0',
  },
  {
    code: '24',
    name: '24Paybank',
    partner: false,
    rating: '4.7',
    reviews: '560',
    rate: '80.74 ₽',
    color: '#1F8A5B',
  },
]

function onClose() {
  close()
}
</script>

<template>
  <Transition
    enter-active-class="transition-opacity duration-200 ease-out"
    enter-from-class="opacity-0"
    leave-active-class="transition-opacity duration-150 ease-in"
    leave-to-class="opacity-0"
  >
    <div v-if="open" class="fixed inset-0 z-50 flex flex-col bg-canvas">
      <!-- search bar -->
      <div class="flex items-center gap-3 px-4 pb-4 pt-5">
        <div
          class="flex flex-1 items-center gap-2.5 rounded-2xl border border-brand bg-surface px-3.5 py-3"
        >
          <svg
            class="flex-none text-ink-faint"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          >
            <circle cx="11" cy="11" r="7" />
            <path d="M21 21l-3.5-3.5" />
          </svg>
          <input
            v-model="query"
            type="text"
            placeholder="Валюта, направление или обменник"
            class="w-full bg-transparent text-[15px] text-ink placeholder:text-ink-faint focus:outline-none"
          />
          <button
            v-if="query"
            type="button"
            aria-label="Очистить"
            class="flex h-5 w-5 flex-none items-center justify-center rounded-full bg-surface-chip text-ink-bright"
            @click="query = ''"
          >
            <svg
              width="10"
              height="10"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.6"
              stroke-linecap="round"
            >
              <path d="M6 6l12 12M18 6 6 18" />
            </svg>
          </button>
        </div>
        <button type="button" class="text-[15px] font-semibold text-brand-bright" @click="onClose">
          Отмена
        </button>
      </div>

      <div class="scrolly flex-1 overflow-y-auto px-4 pb-8">
        <!-- popular directions -->
        <div class="mb-2.5 text-[13px] text-ink-faint">Популярные направления</div>
        <div class="scrollx -mx-4 flex gap-2 overflow-x-auto px-4 pb-1">
          <button
            v-for="d in directions"
            :key="d"
            type="button"
            class="whitespace-nowrap rounded-full border border-line bg-surface px-3.5 py-2 text-[13px] text-ink-bright"
          >
            {{ d }}
          </button>
        </div>

        <!-- results -->
        <div class="mb-2.5 mt-6 text-[13px] text-ink-faint">
          Обменники по запросу «<span class="text-ink-bright">{{ query || '…' }}</span
          >»
        </div>
        <div class="flex flex-col gap-2.5">
          <button
            v-for="r in results"
            :key="r.code"
            type="button"
            class="flex items-center gap-3 rounded-2xl border border-line bg-surface px-3.5 py-3 text-left transition-colors hover:border-line-strong"
          >
            <span
              class="flex h-10 w-10 flex-none items-center justify-center rounded-full text-sm font-bold text-white"
              :class="r.code.length > 2 ? 'text-[13px]' : ''"
              :style="{ backgroundColor: r.color }"
              >{{ r.code }}</span
            >
            <span class="min-w-0 flex-1">
              <span class="flex items-center gap-2 text-[15px] font-semibold text-ink">
                {{ r.name }}
                <span
                  v-if="r.partner"
                  class="rounded-[5px] bg-brand/[0.18] px-1.5 py-0.5 text-[10px] font-semibold text-brand-bright"
                  >Партнёр</span
                >
              </span>
              <span class="mt-0.5 block text-xs text-ink-faint">
                <span class="text-warning">★</span> <span class="tnum">{{ r.rating }}</span> ·
                <span class="tnum">{{ r.reviews }}</span>
              </span>
            </span>
            <span class="tnum flex-none whitespace-nowrap text-base font-bold text-ink">
              {{ r.rate }}
            </span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
