<script setup lang="ts">
import { ref } from 'vue'

// Client onboarding scaffold — 4 steps (design: Auth & Onboarding).
// value-first → soft alert hook → auth → success. No backend wiring yet.
definePageMeta({ layout: false })
useSeoMeta({ title: 'Знакомство — Kurso' })

const { login } = useAuth()
const step = ref(1)
function pick(via: 'telegram' | 'google' | 'email') {
  login({ via })
  step.value = 4
}
const rates = [
  {
    name: 'CryptoBridge',
    initials: 'CB',
    color: '#3A4452',
    rate: '81.20 ₽',
    partner: true,
    best: true,
  },
  { name: 'NetEx24', initials: 'N', color: '#5B3FA0', rate: '80.95 ₽' },
  { name: '24Paybank', initials: '24', color: '#1F8A5B', rate: '80.74 ₽' },
]
const finish = () => navigateTo('/')
</script>

<template>
  <div class="relative flex min-h-[100dvh] flex-col overflow-x-clip bg-canvas">
    <GlowBackdrop />
    <div class="relative z-10 flex flex-1 items-center justify-center px-4 py-8">
      <div class="w-full max-w-[380px]">
        <!-- progress dots (steps 1–3) -->
        <div v-if="step < 4" class="mb-5 flex gap-1.5">
          <span
            v-for="i in 4"
            :key="i"
            class="h-1 flex-1 rounded-full transition-colors"
            :class="i <= step ? 'bg-brand' : 'bg-line-strong'"
          />
        </div>

        <!-- STEP 1 · value first -->
        <div v-if="step === 1">
          <h1 class="text-xl font-extrabold tracking-[-0.02em] text-ink">Лучшие курсы — сразу</h1>
          <p class="mb-4 mt-1.5 text-[13px] text-ink-muted">
            Без регистрации. Выберите пару и смотрите.
          </p>
          <div class="flex flex-col gap-2.5">
            <div
              v-for="r in rates"
              :key="r.name"
              class="flex items-center gap-2.5 rounded-2xl border px-3 py-3"
              :class="r.best ? 'border-brand/30 bg-brand/[0.08]' : 'border-line bg-surface'"
            >
              <span
                class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-full text-[11px] font-bold text-white"
                :class="r.initials.length > 1 ? 'text-[10px]' : ''"
                :style="{ backgroundColor: r.color }"
                >{{ r.initials }}</span
              >
              <div class="flex flex-1 items-center gap-1.5 text-[13px] font-semibold text-ink">
                {{ r.name }}
                <span
                  v-if="r.partner"
                  class="rounded bg-brand/[0.18] px-1.5 py-0.5 text-[9px] font-semibold text-brand-bright"
                  >Партнёр</span
                >
              </div>
              <span class="tnum text-[15px] font-bold text-ink">{{ r.rate }}</span>
            </div>
          </div>
          <KButton block size="lg" class="mt-4 !rounded-2xl" @click="step = 2">Далее</KButton>
        </div>

        <!-- STEP 2 · soft alert hook -->
        <div v-else-if="step === 2">
          <span
            class="mb-3.5 flex h-[46px] w-[46px] items-center justify-center rounded-[13px] bg-brand/[0.14] text-brand-bright"
          >
            <svg
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
              <path d="M10.5 20a1.8 1.8 0 0 0 3 0" />
            </svg>
          </span>
          <h1 class="text-xl font-extrabold tracking-[-0.02em] text-ink">Курс не нравится?</h1>
          <p class="mb-4 mt-1.5 text-sm leading-relaxed text-ink-muted">
            Поставим алерт и пришлём в Telegram, когда USDT → Тинькофф вырастет.
          </p>
          <div
            class="mb-3.5 flex items-baseline gap-2 rounded-2xl border border-brand bg-well px-4 py-3.5"
          >
            <span class="text-[13px] text-ink-faint">Порог</span>
            <span class="tnum ml-auto text-2xl font-bold text-ink">81.50 ₽</span>
          </div>
          <KButton block size="lg" class="!rounded-2xl" @click="step = 3">Настроить алерт</KButton>
          <button
            type="button"
            class="mt-2.5 w-full py-1.5 text-sm font-semibold text-ink-muted transition-colors hover:text-ink"
            @click="finish"
          >
            Позже
          </button>
        </div>

        <!-- STEP 3 · auth -->
        <div v-else-if="step === 3">
          <div class="mb-6 text-center">
            <h1 class="text-[21px] font-extrabold tracking-[-0.02em] text-ink">
              Сохраним ваш алерт
            </h1>
            <p class="mx-auto mt-2 max-w-[300px] text-sm leading-relaxed text-ink-muted">
              Войдите, чтобы получать уведомления. Быстрее всего — через Telegram.
            </p>
          </div>
          <button
            type="button"
            class="mb-3 flex w-full items-center justify-center gap-2.5 rounded-2xl bg-brand py-4 text-base font-bold text-white transition-colors hover:bg-brand-hover"
            @click="pick('telegram')"
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#fff"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 3 2 10.5l6 2.2M22 3l-3 17-8-6.3M22 3 8 12.7m0 0v5.3l3-3.6" />
            </svg>
            Продолжить с Telegram
          </button>
          <div class="mb-[18px] flex gap-2.5">
            <button
              type="button"
              class="flex flex-1 items-center justify-center gap-2 rounded-2xl border border-line-strong bg-surface py-3 text-sm font-semibold text-ink transition-colors hover:border-[#3A4047]"
              @click="pick('google')"
            >
              <span
                class="flex h-5 w-5 items-center justify-center rounded-full bg-white text-xs font-extrabold text-[#1A1A1A]"
                >G</span
              >
              Google
            </button>
            <button
              type="button"
              class="flex flex-1 items-center justify-center rounded-2xl border border-line-strong bg-surface py-3 text-sm font-semibold text-ink-muted transition-colors hover:border-[#3A4047]"
              @click="pick('email')"
            >
              Email
            </button>
          </div>
          <div
            class="flex items-center gap-3 rounded-2xl border border-success/25 bg-success/[0.07] px-[15px] py-3.5"
          >
            <svg
              class="flex-none text-success-bright"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 3 4 6v6c0 5 3.5 7.5 8 9 4.5-1.5 8-4 8-9V6Z" />
              <path d="m9 12 2 2 4-4" />
            </svg>
            <span class="text-[13px] text-ink-muted">Не запрашиваем доступ к кошелькам</span>
          </div>
        </div>

        <!-- STEP 4 · success -->
        <div v-else>
          <div class="mb-[22px] text-center">
            <span
              class="relative mx-auto mb-4 flex h-[76px] w-[76px] items-center justify-center rounded-full bg-success/[0.14]"
            >
              <span class="absolute h-[76px] w-[76px] animate-kpulse rounded-full bg-success/30" />
              <span
                class="relative flex h-[54px] w-[54px] items-center justify-center rounded-full bg-success"
              >
                <svg
                  width="28"
                  height="28"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="#06231A"
                  stroke-width="3"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M5 12l5 5L20 6" />
                </svg>
              </span>
            </span>
            <h1 class="text-[22px] font-extrabold tracking-[-0.02em] text-ink">Алерт создан!</h1>
            <p class="mx-auto mt-2 max-w-[300px] text-sm leading-relaxed text-ink-muted">
              Пришлём в Telegram, когда USDT → Тинькофф достигнет <span class="tnum">81.50 ₽</span>
            </p>
          </div>

          <div class="mb-4 rounded-2xl border border-line bg-surface p-4">
            <div class="mb-2.5 flex items-center justify-between">
              <span class="text-[13px] text-ink-muted">До цели осталось</span>
              <span class="tnum text-[13px] font-semibold text-success-bright">+0.37%</span>
            </div>
            <div class="h-2 overflow-hidden rounded-full bg-well">
              <div
                class="h-full w-[78%] rounded-full bg-[linear-gradient(90deg,#1FB37A,#2BC58C)]"
              />
            </div>
            <div class="mt-2 flex items-center justify-between text-xs text-ink-faint">
              <span class="tnum">сейчас 81.20 ₽</span>
              <span class="tnum">цель 81.50 ₽</span>
            </div>
          </div>

          <div class="mb-2.5 text-[13px] font-bold text-ink">Что дальше</div>
          <div class="mb-4 flex flex-col gap-2.5">
            <NuxtLink
              to="/exchangers"
              class="flex items-center gap-3 rounded-[13px] border border-line bg-surface px-3 py-3 transition-colors hover:border-line-strong"
            >
              <span
                class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-warning/[0.12] text-warning"
              >
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M12 20s-6.5-4-9-8.5A4.5 4.5 0 0 1 12 7a4.5 4.5 0 0 1 9 4.5C18.5 16 12 20 12 20Z"
                  />
                </svg>
              </span>
              <span class="flex-1 text-[13px] text-ink">Добавьте обменники в избранное</span>
              <svg
                class="text-ink-ghost"
                width="15"
                height="15"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M9 6l6 6-6 6" />
              </svg>
            </NuxtLink>
            <NuxtLink
              to="/map"
              class="flex items-center gap-3 rounded-[13px] border border-line bg-surface px-3 py-3 transition-colors hover:border-line-strong"
            >
              <span
                class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-brand/[0.12] text-brand-bright"
              >
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M12 21s-7-5.5-7-11a7 7 0 0 1 14 0c0 5.5-7 11-7 11Z" />
                  <circle cx="12" cy="10" r="2.5" />
                </svg>
              </span>
              <span class="flex-1 text-[13px] text-ink">Найдите обменники на карте</span>
              <svg
                class="text-ink-ghost"
                width="15"
                height="15"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M9 6l6 6-6 6" />
              </svg>
            </NuxtLink>
          </div>

          <KButton block size="lg" class="!rounded-2xl" @click="finish">Перейти к курсам</KButton>
        </div>
      </div>
    </div>
  </div>
</template>
