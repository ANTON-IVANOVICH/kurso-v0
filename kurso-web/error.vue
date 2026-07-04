<script setup lang="ts">
// Nuxt root error page — fully replaces the app, so it renders its own
// GlowBackdrop + chrome-less full-screen frame (NOT the default layout).
// Mobile = centered compact layout (design v1); desktop = two-column (v2).
const props = defineProps<{ error: { statusCode?: number; message?: string } }>()

const is404 = computed(() => props.error?.statusCode === 404)

const clearTo = (path: string) => clearError({ redirect: path })
const goHome = () => clearTo('/')
const refresh = () => {
  if (import.meta.client) location.reload()
}

type Row = { code?: string; name?: string; rate?: string; bg?: string; notFound?: boolean }
// Decorative "rates" list for the 404 right panel (desktop only).
const rateRows: Row[] = [
  { code: 'СБ', name: 'Сбербанк', rate: '81 240 ₽', bg: '#1E8E4A' },
  { code: 'ТБ', name: 'Т-Банк', rate: '81 050 ₽', bg: '#2E7DF2' },
  { notFound: true },
  { code: 'АБ', name: 'Альфа-Банк', rate: '80 890 ₽', bg: '#E24B4B' },
  { code: 'НЛ', name: 'Наличные', rate: '80 500 ₽', bg: '#C77D2E' },
  { code: 'ЮM', name: 'ЮMoney', rate: '80 320 ₽', bg: '#8B5CF6' },
]

const directions = ['USDT → Сбербанк', 'BTC → Т-Банк', 'TON → Наличные']

type SysState = 'ok' | 'fail' | 'degraded'
const systems: { label: string; state: SysState }[] = [
  { label: 'API курсов', state: 'ok' },
  { label: 'Парсеры фидов', state: 'fail' },
  { label: 'Виджеты', state: 'degraded' },
  { label: 'Telegram-алерты', state: 'ok' },
  { label: 'База резервов', state: 'ok' },
]
const sysTone: Record<SysState, 'success' | 'danger' | 'warning'> = {
  ok: 'success',
  fail: 'danger',
  degraded: 'warning',
}
const sysText: Record<SysState, string> = { ok: 'ОК', fail: 'Сбой', degraded: 'Деградация' }
const sysColor: Record<SysState, string> = {
  ok: 'text-success',
  fail: 'text-danger',
  degraded: 'text-warning',
}
</script>

<template>
  <div class="relative flex min-h-screen flex-col overflow-x-clip bg-canvas">
    <GlowBackdrop :tone="is404 ? 'brand' : 'danger'" />
    <!-- ambient accent glow that bleeds off the top-right corner of the page -->
    <div
      aria-hidden="true"
      class="pointer-events-none absolute right-[-120px] top-[-140px] z-0 h-[520px] w-[620px]"
      :class="
        is404
          ? 'bg-[radial-gradient(ellipse_at_center,rgba(46,125,242,0.18),rgba(46,125,242,0)_64%)]'
          : 'bg-[radial-gradient(ellipse_at_center,rgba(236,91,74,0.15),rgba(236,91,74,0)_64%)]'
      "
    />

    <header class="relative z-10 mx-auto w-full max-w-[1100px] px-4 pt-8 md:px-6 md:pt-10">
      <AppLogo :size="34" />
    </header>

    <main
      class="relative z-10 mx-auto flex w-full max-w-[1100px] flex-1 items-center px-4 py-8 md:px-6 md:py-10"
    >
      <!-- ========================= 404 ========================= -->
      <div
        v-if="is404"
        class="grid w-full gap-x-12 gap-y-10 md:min-h-[460px] md:grid-cols-[1.05fr_0.95fr] md:items-center"
      >
        <!-- content: centered on mobile, left-aligned on desktop -->
        <div class="flex flex-col items-center text-center md:items-start md:text-left">
          <div class="font-label text-[11px] uppercase tracking-[0.16em] text-brand-bright">
            Ошибка <span class="tnum">404</span> · страница не найдена
          </div>

          <div
            class="tnum mt-4 bg-[linear-gradient(150deg,#6BA6FF,#2E7DF2)] bg-clip-text text-[80px] font-extrabold leading-[0.9] text-transparent sm:text-[96px] md:text-[120px]"
          >
            404
          </div>

          <h1 class="mt-2 text-[22px] font-extrabold text-ink md:text-[30px]">
            Этот курс потерялся
          </h1>
          <p class="mt-3 max-w-[420px] text-[15px] leading-relaxed text-ink-muted">
            Возможно, ссылка устарела или направление больше не отслеживается. Актуальные курсы — на
            месте.
          </p>

          <div class="mt-6 flex flex-wrap justify-center gap-3 md:justify-start">
            <KButton @click="goHome">На главную</KButton>
            <KButton variant="secondary" @click="clearTo('/exchangers')">Все обменники</KButton>
          </div>

          <!-- faux search -->
          <div
            class="mt-5 flex w-full max-w-[420px] items-center gap-2.5 rounded-xl border border-line bg-surface-hi px-4 py-3 text-sm text-ink-faint"
          >
            <svg
              class="flex-none"
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="11" cy="11" r="7" />
              <path d="M21 21l-3.5-3.5" />
            </svg>
            <span class="truncate">Найти направление, например «USDT → Сбербанк»…</span>
          </div>

          <div class="mt-6 w-full">
            <div class="font-label text-[11px] uppercase tracking-[0.16em] text-ink-ghost">
              Популярные направления
            </div>
            <div class="mt-3 flex flex-wrap justify-center gap-2 md:justify-start">
              <span
                v-for="d in directions"
                :key="d"
                class="rounded-full border border-line-strong bg-surface-hi px-3 py-1.5 text-[13px] text-ink-bright"
              >
                {{ d }}
              </span>
            </div>
          </div>
        </div>

        <!-- RIGHT (decorative, desktop only) -->
        <div class="relative hidden min-h-[440px] overflow-hidden md:block">
          <div
            class="absolute inset-0 flex rotate-[-7deg] scale-[1.12] flex-col justify-center gap-3 p-10 [-webkit-mask-image:linear-gradient(to_bottom,transparent,#000_22%,#000_78%,transparent)] [mask-image:linear-gradient(to_bottom,transparent,#000_22%,#000_78%,transparent)]"
          >
            <template v-for="(r, i) in rateRows" :key="i">
              <div
                v-if="r.notFound"
                class="flex items-center gap-3 rounded-xl border border-dashed border-brand/50 bg-brand/[0.06] px-4 py-3"
              >
                <span
                  class="flex h-9 w-9 flex-none items-center justify-center rounded-lg border border-dashed border-brand/50 text-sm font-bold text-brand-bright"
                  >?</span
                >
                <span class="text-sm text-brand-bright">направление не найдено</span>
                <span class="tnum ml-auto text-sm text-ink-faint">— ₽</span>
              </div>
              <div
                v-else
                class="flex items-center gap-3 rounded-xl border border-line bg-surface-nested px-4 py-3"
              >
                <span
                  class="flex h-9 w-9 flex-none items-center justify-center rounded-lg text-xs font-bold text-white"
                  :style="{ background: r.bg }"
                  >{{ r.code }}</span
                >
                <span class="text-sm text-ink-bright">{{ r.name }}</span>
                <span class="tnum ml-auto text-sm text-ink">{{ r.rate }}</span>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- ========================= 500 ========================= -->
      <div
        v-else
        class="grid w-full gap-x-12 gap-y-10 md:min-h-[460px] md:grid-cols-[1.05fr_0.95fr] md:items-center"
      >
        <!-- content: centered on mobile, left-aligned on desktop -->
        <div class="flex flex-col items-center text-center md:items-start md:text-left">
          <!-- icon + code: stacked & centered on mobile, inline on desktop -->
          <div class="flex flex-col items-center gap-3 md:flex-row md:gap-4">
            <span
              class="flex h-[68px] w-[68px] items-center justify-center rounded-[18px] bg-danger/[0.12] text-danger md:h-[46px] md:w-[46px] md:rounded-[13px]"
            >
              <svg
                class="h-[30px] w-[30px] md:h-[22px] md:w-[22px]"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M12 9v4M12 17h.01" />
                <path
                  d="M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z"
                />
              </svg>
            </span>
            <span class="tnum text-[36px] font-extrabold leading-none text-danger md:text-[30px]"
              >500</span
            >
          </div>

          <h1 class="mt-5 text-[22px] font-extrabold text-ink md:mt-6 md:text-[30px]">
            Курсы временно недоступны
          </h1>
          <p class="mt-3 max-w-[420px] text-[15px] leading-relaxed text-ink-muted">
            На стороне сервиса произошёл сбой при обновлении котировок. Мы уже разбираемся —
            обновите страницу через пару минут. Ваши алерты продолжают работать.
          </p>

          <div class="mt-6 flex flex-wrap justify-center gap-3 md:justify-start">
            <KButton @click="refresh">
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
            <KButton variant="secondary" @click="clearTo('/maintenance')">Статус-страница</KButton>
          </div>

          <div class="mt-6 inline-flex items-center gap-2 font-label text-xs text-ink-faint">
            <KStatusDot tone="danger" :size="6" />
            <span
              >request-id: <span class="tnum">0x9F-A21C-4E</span> ·
              <span class="tnum">14:06:22</span> МСК</span
            >
          </div>
        </div>

        <!-- RIGHT — system status (desktop only) -->
        <div class="relative hidden md:flex md:items-center md:pl-10">
          <div class="w-full rounded-2xl border border-line bg-surface-nested p-5">
            <div class="flex items-center justify-between border-b border-line-subtle pb-3">
              <span class="text-sm font-semibold text-ink">Состояние систем</span>
              <span class="text-xs text-ink-faint">uptime <span class="tnum">99.94%</span></span>
            </div>
            <div class="mt-2 flex flex-col gap-0.5">
              <div
                v-for="s in systems"
                :key="s.label"
                class="flex items-center justify-between rounded-lg px-3 py-2.5"
                :class="s.state === 'fail' ? 'bg-danger/[0.06]' : ''"
              >
                <span class="text-sm text-ink-bright">{{ s.label }}</span>
                <span
                  class="inline-flex items-center gap-2 text-[13px] font-medium"
                  :class="sysColor[s.state]"
                >
                  <KStatusDot :tone="sysTone[s.state]" :size="7" />{{ sysText[s.state] }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
