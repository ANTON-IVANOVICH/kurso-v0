<script setup lang="ts">
// Notifications feed (opened from the header bell). Scaffold with representative
// items: triggered alerts, new best rates, review replies, reserve updates.
const { open, close } = useNotifications()
useScrollLock(open)

type NotifType = 'alert' | 'rate' | 'review' | 'reserve'
interface Notif {
  type: NotifType
  title: string
  time: string
  unread?: boolean
}

const items: Notif[] = [
  { type: 'alert', title: 'Сработал алерт', time: '2 мин назад', unread: true },
  { type: 'rate', title: 'Новый лучший курс', time: '1 ч назад' },
  { type: 'review', title: 'Ответ на ваш отзыв', time: '3 ч назад' },
  { type: 'reserve', title: 'Резерв пополнен', time: 'вчера' },
]

function onClose() {
  close()
}
</script>

<template>
  <!-- backdrop -->
  <Transition
    enter-active-class="transition-opacity duration-200 ease-out"
    enter-from-class="opacity-0"
    leave-active-class="transition-opacity duration-150 ease-in"
    leave-to-class="opacity-0"
  >
    <div v-if="open" class="fixed inset-0 z-50 bg-[#040507]/60" @click="onClose" />
  </Transition>

  <!-- panel -->
  <Transition
    enter-active-class="transition-transform duration-300 ease-out"
    enter-from-class="translate-x-full"
    leave-active-class="transition-transform duration-200 ease-in"
    leave-to-class="translate-x-full"
  >
    <div
      v-if="open"
      class="fixed inset-y-0 right-0 z-50 flex w-full max-w-[420px] flex-col border-l border-line bg-canvas"
    >
      <!-- header -->
      <div class="flex items-center justify-between gap-3 px-4 pb-3 pt-5">
        <button
          type="button"
          aria-label="Закрыть"
          class="flex h-[38px] w-[38px] flex-none items-center justify-center rounded-full border border-line bg-surface text-ink-muted transition-colors hover:border-line-strong"
          @click="onClose"
        >
          <svg
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M15 5l-7 7 7 7" />
          </svg>
        </button>
        <span class="text-[17px] font-bold text-ink">Уведомления</span>
        <button type="button" class="text-[13px] font-semibold text-brand-bright">Прочитать</button>
      </div>

      <!-- feed -->
      <div class="scrolly flex flex-1 flex-col gap-2.5 overflow-y-auto px-4 pb-8">
        <div
          v-for="(n, i) in items"
          :key="i"
          class="flex gap-3 rounded-2xl border p-3.5"
          :class="n.unread ? 'border-success/30 bg-success/[0.08]' : 'border-line bg-surface'"
        >
          <span
            class="flex h-10 w-10 flex-none items-center justify-center rounded-full"
            :class="{
              'bg-success': n.type === 'alert',
              'bg-brand': n.type === 'rate',
              'bg-surface-raised': n.type === 'review' || n.type === 'reserve',
            }"
          >
            <!-- alert: target -->
            <svg
              v-if="n.type === 'alert'"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#0A0B0D"
              stroke-width="2.4"
            >
              <circle cx="12" cy="12" r="8" />
              <circle cx="12" cy="12" r="3" />
            </svg>
            <!-- rate: bell -->
            <svg
              v-else-if="n.type === 'rate'"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#fff"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
              <path d="M10.5 20a1.8 1.8 0 0 0 3 0" />
            </svg>
            <!-- review: chat -->
            <svg
              v-else-if="n.type === 'review'"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#fff"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M21 15a2 2 0 0 1-2 2H8l-4 4V6a2 2 0 0 1 2-2h13a2 2 0 0 1 2 2Z" />
            </svg>
            <!-- reserve: up arrow -->
            <svg
              v-else
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#2BC58C"
              stroke-width="2.2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 19V5M12 5l-6 6M12 5l6 6" />
            </svg>
          </span>

          <div class="min-w-0 flex-1">
            <div class="text-[15px] font-semibold text-ink">{{ n.title }}</div>
            <div class="mt-0.5 text-[13px] leading-snug text-ink-muted">
              <template v-if="n.type === 'alert'"
                >USDT → Тинькофф достиг <span class="tnum">81.50 ₽</span> в CryptoBridge</template
              >
              <template v-else-if="n.type === 'rate'"
                >ETH → Альфа: <span class="tnum">312 400 ₽</span> — на
                <span class="tnum">0.4%</span> выше</template
              >
              <template v-else-if="n.type === 'review'"
                >NetEx24 ответил на ваш отзыв об обмене</template
              >
              <template v-else>24Paybank: USDT — резерв <span class="tnum">18.2M ₽</span></template>
            </div>
            <div class="tnum mt-1.5 text-xs text-ink-faint">{{ n.time }}</div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>
