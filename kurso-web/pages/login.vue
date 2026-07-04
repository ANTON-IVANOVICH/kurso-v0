<script setup lang="ts">
import { ref } from 'vue'

// Client login scaffold. Telegram flow switches to an inline "waiting for
// confirmation" state (design: Auth & Onboarding). No backend wiring yet.
definePageMeta({ layout: 'auth' })
useSeoMeta({ title: 'Вход — Kurso' })

const mode = ref<'form' | 'telegram'>('form')
const email = ref('')
const password = ref('')
const showPassword = ref(false)
</script>

<template>
  <!-- ===== email / social form ===== -->
  <div v-if="mode === 'form'">
    <span
      class="mb-[18px] flex h-12 w-12 items-center justify-center rounded-[14px] bg-[linear-gradient(150deg,#4A90F5,#2E7DF2)] shadow-glow"
    >
      <svg
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="#fff"
        stroke-width="2.1"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M7 4 L7 20 M7 20 L4 17 M7 20 L10 17 M17 20 L17 4 M17 4 L14 7 M17 4 L20 7" />
      </svg>
    </span>

    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">С возвращением</h1>
    <p class="mb-5 mt-1.5 text-sm text-ink-muted">Войдите, чтобы управлять алертами и избранным</p>

    <div class="mb-[18px] flex flex-col gap-2.5">
      <button
        type="button"
        class="flex w-full items-center justify-center gap-2.5 rounded-2xl bg-brand py-3.5 text-[15px] font-semibold text-white transition-colors hover:bg-brand-hover"
        @click="mode = 'telegram'"
      >
        <svg
          width="19"
          height="19"
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
      <div class="flex gap-2.5">
        <button
          type="button"
          class="flex flex-1 items-center justify-center gap-2.5 rounded-2xl border border-line-strong bg-surface py-3 text-sm font-semibold text-ink transition-colors hover:border-[#3A4047]"
        >
          <span
            class="flex h-5 w-5 items-center justify-center rounded-full bg-white text-xs font-extrabold text-[#1A1A1A]"
            >G</span
          >
          Google
        </button>
        <button
          type="button"
          class="flex flex-1 items-center justify-center gap-2.5 rounded-2xl border border-line-strong bg-surface py-3 text-sm font-semibold text-ink transition-colors hover:border-[#3A4047]"
        >
          <svg width="17" height="17" viewBox="0 0 24 24" fill="currentColor">
            <path
              d="M17.05 12.54c-.02-2.06 1.68-3.05 1.76-3.1-.96-1.4-2.45-1.6-2.98-1.62-1.27-.13-2.48.74-3.12.74-.64 0-1.64-.72-2.7-.7-1.39.02-2.67.8-3.38 2.04-1.44 2.5-.37 6.2 1.04 8.23.69.99 1.5 2.1 2.57 2.06 1.03-.04 1.42-.66 2.67-.66 1.24 0 1.6.66 2.69.64 1.11-.02 1.81-1 2.49-2 .78-1.15 1.1-2.26 1.12-2.32-.02-.01-2.15-.83-2.18-3.26ZM15 6.8c.57-.69.95-1.65.85-2.6-.82.03-1.81.54-2.4 1.23-.52.6-.98 1.58-.86 2.5.91.08 1.84-.46 2.41-1.13Z"
            />
          </svg>
          Apple
        </button>
      </div>
    </div>

    <div class="mb-[18px] flex items-center gap-3">
      <span class="h-px flex-1 bg-line" />
      <span class="text-xs text-ink-faint">или по email</span>
      <span class="h-px flex-1 bg-line" />
    </div>

    <label
      class="mb-2.5 flex items-center gap-2.5 rounded-2xl border border-line bg-surface px-[15px] py-3.5 focus-within:border-brand"
    >
      <svg
        class="flex-none text-ink-faint"
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <rect x="3" y="5" width="18" height="14" rx="2" />
        <path d="m3 7 9 6 9-6" />
      </svg>
      <input
        v-model="email"
        type="email"
        placeholder="you@email.com"
        class="tnum w-full bg-transparent text-[15px] text-ink placeholder:text-ink-faint focus:outline-none"
      />
    </label>

    <label
      class="flex items-center gap-2.5 rounded-2xl border border-line bg-surface px-[15px] py-3.5 focus-within:border-brand"
    >
      <svg
        class="flex-none text-ink-faint"
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <rect x="4" y="10" width="16" height="10" rx="2" />
        <path d="M8 10V7a4 4 0 0 1 8 0v3" />
      </svg>
      <input
        v-model="password"
        :type="showPassword ? 'text' : 'password'"
        placeholder="Пароль"
        class="w-full bg-transparent text-[15px] text-ink placeholder:text-ink-faint focus:outline-none"
      />
      <button
        type="button"
        class="flex-none text-ink-faint transition-colors hover:text-ink-muted"
        :aria-label="showPassword ? 'Скрыть пароль' : 'Показать пароль'"
        @click="showPassword = !showPassword"
      >
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
          <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z" />
          <circle cx="12" cy="12" r="3" />
        </svg>
      </button>
    </label>

    <div class="mb-[18px] mt-2.5 text-right">
      <button type="button" class="text-[13px] text-brand-bright">Забыли пароль?</button>
    </div>

    <KButton block size="lg" class="!rounded-2xl">Войти</KButton>

    <div class="mt-[18px] text-center text-sm text-ink-faint">
      Нет аккаунта?
      <NuxtLink to="/register" class="font-semibold text-brand-bright">Создать</NuxtLink>
    </div>
  </div>

  <!-- ===== telegram confirmation waiting ===== -->
  <div v-else>
    <button
      type="button"
      aria-label="Назад"
      class="mb-2 flex h-[38px] w-[38px] items-center justify-center rounded-full border border-line bg-surface text-ink transition-colors hover:border-line-strong"
      @click="mode = 'form'"
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

    <div class="px-2 text-center">
      <span
        class="mx-auto mb-[22px] flex h-[84px] w-[84px] items-center justify-center rounded-3xl bg-[linear-gradient(150deg,#4A90F5,#2E7DF2)] shadow-[0_12px_30px_rgba(46,125,242,0.4)]"
      >
        <svg
          width="40"
          height="40"
          viewBox="0 0 24 24"
          fill="none"
          stroke="#fff"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M22 3 2 10.5l6 2.2M22 3l-3 17-8-6.3M22 3 8 12.7m0 0v5.3l3-3.6" />
        </svg>
      </span>
      <h1 class="text-[22px] font-extrabold tracking-[-0.02em] text-ink">Подтвердите вход</h1>
      <p class="mx-auto mb-[26px] mt-2 max-w-[300px] text-sm leading-relaxed text-ink-muted">
        Откройте нашего бота в Telegram и нажмите «Подтвердить» — вернётесь сюда автоматически
      </p>
    </div>

    <div class="mb-[26px] flex flex-col gap-3.5">
      <div class="flex items-center gap-3">
        <span class="flex h-7 w-7 flex-none items-center justify-center rounded-full bg-success">
          <svg
            width="15"
            height="15"
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
        <span class="text-sm text-ink"
          >Открыли бота <span class="text-brand-bright">@kurso_bot</span></span
        >
      </div>
      <div class="flex items-center gap-3">
        <span
          class="flex h-7 w-7 flex-none items-center justify-center rounded-full border-2 border-brand"
        >
          <span
            class="h-3.5 w-3.5 animate-kspin rounded-full border-2 border-brand border-t-transparent"
          />
        </span>
        <span class="text-sm text-ink">Нажмите «Подтвердить» в боте</span>
      </div>
      <div class="flex items-center gap-3">
        <span
          class="tnum flex h-7 w-7 flex-none items-center justify-center rounded-full border-2 border-line-strong text-xs font-bold text-ink-faint"
          >3</span
        >
        <span class="text-sm text-ink-faint">Вернитесь в приложение</span>
      </div>
    </div>

    <div class="mb-[18px] flex items-center justify-center gap-2 text-[13px] text-ink-faint">
      <KStatusDot tone="success" pulse :size="7" />Ждём подтверждения…
    </div>

    <KButton block size="lg" class="!rounded-2xl">
      Открыть @kurso_bot
      <svg
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.4"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M7 17 17 7M9 7h8v8" />
      </svg>
    </KButton>
  </div>
</template>
