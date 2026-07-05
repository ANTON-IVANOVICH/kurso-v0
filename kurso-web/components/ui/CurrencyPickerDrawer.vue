<script setup lang="ts">
import { computed, ref, watch } from 'vue'

// Bottom sheet for choosing the buy/sell currency. Slides up over a dimmed
// backdrop; on desktop it stays centred and width-capped. Wired to the global
// currency-picker state so any selector opens the same drawer.
const { t } = useI18n()
const { open, side, give, get, currencyOptions, choose, close } = useCurrencyPicker()
useScrollLock(open)

const query = ref('')
const current = computed(() => (side.value === 'give' ? give.value : get.value))

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return currencyOptions
  return currencyOptions.filter((o) =>
    `${o.name} ${o.subtitle} ${o.aliases ?? ''}`.toLowerCase().includes(q),
  )
})

function onClose() {
  query.value = ''
  close()
}
function onChoose(option: (typeof currencyOptions)[number]) {
  query.value = ''
  choose(option)
}

// Swipe-to-dismiss from the grab handle. Pointer Events cover mouse + touch.
const dragging = ref(false)
const dragOffset = ref(0)
let startY = 0

function onDragStart(e: PointerEvent) {
  dragging.value = true
  startY = e.clientY
  try {
    ;(e.currentTarget as HTMLElement).setPointerCapture?.(e.pointerId)
  } catch {
    // Ignore invalid-pointer errors (e.g. synthetic events); capture is optional.
  }
}
function onDragMove(e: PointerEvent) {
  if (!dragging.value) return
  dragOffset.value = Math.max(0, e.clientY - startY)
}
function onDragEnd() {
  if (!dragging.value) return
  dragging.value = false
  if (dragOffset.value > 90) onClose()
  else dragOffset.value = 0
}

watch(open, (v) => {
  if (!v) {
    dragOffset.value = 0
    dragging.value = false
  }
})
</script>

<template>
  <!-- backdrop -->
  <Transition
    enter-active-class="transition-opacity duration-200 ease-out"
    enter-from-class="opacity-0"
    leave-active-class="transition-opacity duration-200 ease-in"
    leave-to-class="opacity-0"
  >
    <div v-if="open" class="fixed inset-0 z-50 bg-[#040507]/60" @click="onClose" />
  </Transition>

  <!-- sheet -->
  <Transition
    enter-active-class="transition-transform duration-300 ease-out"
    enter-from-class="translate-y-full"
    leave-active-class="transition-transform duration-200 ease-in"
    leave-to-class="translate-y-full"
  >
    <div v-if="open" class="fixed inset-x-0 bottom-0 z-50 flex justify-center">
      <div
        class="flex w-full max-w-[700px] flex-col rounded-t-[26px] border border-b-0 border-line bg-surface pb-6 shadow-[0_-24px_60px_rgba(0,0,0,0.6)] will-change-transform"
        :class="dragging ? '' : 'transition-transform duration-200 ease-out'"
        :style="dragOffset ? { transform: `translateY(${dragOffset}px)` } : undefined"
      >
        <!-- grab handle + title: the whole zone is the drag affordance -->
        <div
          class="flex-none cursor-grab touch-none select-none px-4 pt-4 active:cursor-grabbing sm:px-6"
          @pointerdown="onDragStart"
          @pointermove="onDragMove"
          @pointerup="onDragEnd"
          @pointercancel="onDragEnd"
        >
          <div class="mb-4 flex justify-center">
            <span class="h-[5px] w-10 rounded-full bg-line-strong" />
          </div>
          <div class="mb-3 text-center text-sm font-semibold text-ink-muted">
            {{ side === 'give' ? t('currencyPicker.give') : t('currencyPicker.get') }}
          </div>
        </div>

        <!-- search — pinned; never scrolls with the list -->
        <div class="flex-none px-4 sm:px-6">
          <label
            class="mb-4 flex items-center gap-2.5 rounded-2xl border border-line bg-well px-[15px] py-3 text-ink-muted focus-within:border-brand"
          >
            <svg
              class="flex-none"
              width="17"
              height="17"
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
              :placeholder="t('currencyPicker.searchPlaceholder')"
              class="w-full bg-transparent text-sm text-ink placeholder:text-ink-faint focus:outline-none"
            />
          </label>
        </div>

        <!-- list — fixed height so filtering never shifts the search bar -->
        <div class="scrolly h-[44vh] overflow-y-auto px-4 sm:px-6">
          <div v-if="filtered.length" class="grid grid-cols-1 gap-2.5 sm:grid-cols-2">
            <button
              v-for="o in filtered"
              :key="o.code"
              type="button"
              class="flex items-center gap-3 rounded-2xl border bg-well px-3.5 py-3 text-left transition-colors"
              :class="
                o.code === current.code ? 'border-brand' : 'border-line hover:border-line-strong'
              "
              @click="onChoose(o)"
            >
              <span
                class="flex h-[38px] w-[38px] flex-none items-center justify-center rounded-full font-bold"
                :class="o.symbol.length > 1 ? 'text-[11px]' : 'text-base'"
                :style="{ backgroundColor: o.color, color: o.dark ? '#111' : '#fff' }"
                >{{ o.symbol }}</span
              >
              <span class="min-w-0 flex-1">
                <span class="block truncate text-[15px] font-semibold text-ink">{{ o.name }}</span>
                <span class="block text-xs text-ink-faint">{{ o.subtitle }}</span>
              </span>
              <svg
                class="flex-none text-ink-ghost sm:hidden"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M9 6l6 6-6 6" />
              </svg>
            </button>
          </div>
          <div v-else class="py-10 text-center text-sm text-ink-faint">
            {{ t('currencyPicker.noResults') }}
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>
