<script setup lang="ts">
// Standalone full-screen maintenance page (chrome-less — its own layout).
// Mobile = centered compact layout (design v1); desktop = two-column (v2).
definePageMeta({ layout: false })

const refresh = () => {
  if (import.meta.client) location.reload()
}
const openTelegram = () => {
  if (import.meta.client) window.open('https://t.me/kurso', '_blank', 'noopener')
}

type StepState = 'done' | 'active' | 'pending'
const steps: { label: string; state: StepState; num?: string }[] = [
  { label: 'Остановка', state: 'done' },
  { label: 'Миграция', state: 'done' },
  { label: 'Обновление', state: 'active' },
  { label: 'Запуск', state: 'pending', num: '4' },
]
// Connector before step i+1 is filled once the prior step is done and the next reached.
const isFilled = (i: number) => steps[i].state === 'done' && steps[i + 1].state !== 'pending'
</script>

<template>
  <div class="relative flex min-h-screen flex-col overflow-x-clip bg-canvas">
    <GlowBackdrop tone="warning" />
    <!-- ambient accent glow that bleeds off the top-right corner of the page -->
    <div
      aria-hidden="true"
      class="pointer-events-none absolute right-[-120px] top-[-140px] z-0 h-[520px] w-[620px] bg-[radial-gradient(ellipse_at_center,rgba(217,154,51,0.16),rgba(217,154,51,0)_64%)]"
    />

    <header class="relative z-10 mx-auto w-full max-w-[1100px] px-4 pt-8 md:px-6 md:pt-10">
      <AppLogo :size="34" />
    </header>

    <main
      class="relative z-10 mx-auto flex w-full max-w-[1100px] flex-1 items-center px-4 py-8 md:px-6 md:py-10"
    >
      <div
        class="grid w-full gap-x-12 gap-y-10 md:min-h-[460px] md:grid-cols-[1.05fr_0.95fr] md:items-center"
      >
        <!-- content: centered on mobile, left-aligned on desktop -->
        <div class="flex flex-col items-center text-center md:items-start md:text-left">
          <!-- mobile: circular spinner icon -->
          <span
            class="flex h-[68px] w-[68px] items-center justify-center rounded-full bg-warning/[0.12] text-warning md:hidden"
          >
            <svg
              class="h-8 w-8 animate-kspin"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 2v3M12 19v3M2 12h3M19 12h3M5 5l2 2M17 17l2 2M5 19l2-2M17 7l2-2" />
              <circle cx="12" cy="12" r="4" />
            </svg>
          </span>

          <!-- label (inline spinner on desktop) -->
          <div class="mt-4 flex items-center gap-2.5 md:mt-0">
            <svg
              class="hidden h-[18px] w-[18px] animate-kspin text-warning md:block"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 2v3M12 19v3M2 12h3M19 12h3M5 5l2 2M17 17l2 2M5 19l2-2M17 7l2-2" />
              <circle cx="12" cy="12" r="4" />
            </svg>
            <span class="font-label text-[11px] uppercase tracking-[0.16em] text-warning"
              >Плановые работы</span
            >
          </div>

          <h1 class="mt-4 text-[22px] font-extrabold text-ink md:mt-5 md:text-[30px]">
            Обновляем парсеры курсов
          </h1>
          <p class="mt-3 max-w-[420px] text-[15px] leading-relaxed text-ink-muted">
            Делаем сбор курсов быстрее и точнее. Сайт ненадолго недоступен, но алерты и Telegram-бот
            продолжают работать — уведомим, когда вернёмся.
          </p>

          <!-- mobile: compact countdown pill -->
          <div
            class="mt-6 inline-flex items-center gap-2.5 rounded-xl border border-line bg-surface-hi px-4 py-3 md:hidden"
          >
            <KStatusDot tone="warning" pulse :size="8" />
            <span class="text-[13px] text-ink-muted"
              >Ориентировочно до <span class="tnum font-semibold text-ink">14:30 МСК</span></span
            >
          </div>

          <!-- desktop: 4-step tracker -->
          <div class="mt-8 hidden w-full max-w-[440px] items-start md:flex">
            <template v-for="(s, i) in steps" :key="s.label">
              <div class="flex flex-col items-center gap-2 text-center">
                <span
                  class="flex h-9 w-9 items-center justify-center rounded-full"
                  :class="{
                    'bg-success text-white': s.state === 'done',
                    'border-2 border-warning bg-warning/[0.12] text-warning': s.state === 'active',
                    'border border-line-strong text-ink-faint': s.state === 'pending',
                  }"
                >
                  <svg
                    v-if="s.state === 'done'"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.4"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M5 12l5 5L20 6" />
                  </svg>
                  <svg
                    v-else-if="s.state === 'active'"
                    class="animate-kspin"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M21 12a9 9 0 1 1-6.2-8.5" />
                  </svg>
                  <span v-else class="tnum text-sm font-semibold">{{ s.num }}</span>
                </span>
                <span
                  class="text-[11px] font-medium"
                  :class="{
                    'text-ink-bright': s.state === 'done',
                    'text-warning': s.state === 'active',
                    'text-ink-faint': s.state === 'pending',
                  }"
                  >{{ s.label }}</span
                >
              </div>
              <div
                v-if="i < steps.length - 1"
                class="mx-1.5 mt-[17px] h-[2px] flex-1 rounded-full"
                :class="isFilled(i) ? 'bg-success' : 'bg-line'"
              />
            </template>
          </div>

          <div class="mt-7 flex flex-wrap justify-center gap-3 md:mt-8 md:justify-start">
            <KButton @click="openTelegram">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M22 3 2 10.5l6 2.2M22 3l-3 17-8-6.3" />
              </svg>
              Статус в Telegram
            </KButton>
            <KButton variant="secondary" @click="refresh">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M21 12a9 9 0 1 1-3-6.7M21 4v5h-5" />
              </svg>
              Обновить
            </KButton>
          </div>
        </div>

        <!-- RIGHT — countdown (desktop only) -->
        <div class="relative hidden md:flex md:items-center md:pl-10">
          <div class="w-full rounded-2xl border border-line bg-surface-nested p-6 text-center">
            <div class="font-label text-[11px] uppercase tracking-[0.16em] text-ink-faint">
              Запуск ориентировочно через
            </div>
            <div
              class="tnum mt-3 bg-[linear-gradient(150deg,#F2D08A,#E0A954)] bg-clip-text text-[52px] font-extrabold leading-none text-transparent"
            >
              00:14:32
            </div>
            <div class="mt-3 inline-flex items-center gap-2 text-sm text-ink-muted">
              <KStatusDot tone="warning" pulse :size="7" />до <span class="tnum">14:30</span> МСК
            </div>

            <div class="mt-6 h-1.5 w-full overflow-hidden rounded-full bg-line">
              <div
                class="h-full w-[66%] rounded-full bg-[linear-gradient(90deg,#D99A33,#E0A954)]"
              />
            </div>
            <div class="tnum mt-3 text-xs text-ink-faint">готово 66% · шаг 3 из 4</div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
