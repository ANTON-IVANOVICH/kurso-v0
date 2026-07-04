<script setup lang="ts">
import { computed, ref } from 'vue'

// Real registration against the API (bcrypt user). Social sign-up needs the
// OAuth/bot backend and is flagged as upcoming rather than faked.
definePageMeta({ layout: 'auth' })
useSeoMeta({ title: 'Регистрация — Kurso' })

const { register } = useAuth()
const email = ref('')
const password = ref('')
const showPassword = ref(false)
const agreed = ref(false)
const busy = ref(false)
const error = ref('')

async function doRegister() {
  error.value = ''
  if (!email.value.trim() || !password.value) {
    error.value = 'Введите email и пароль'
    return
  }
  if (password.value.length < 8) {
    error.value = 'Пароль минимум 8 символов'
    return
  }
  busy.value = true
  try {
    await register(email.value.trim(), password.value)
    await navigateTo('/account')
  } catch (e) {
    const msg = (e as { data?: { message?: string } })?.data?.message
    error.value = msg || 'Не удалось создать аккаунт'
  } finally {
    busy.value = false
  }
}

function soon() {
  error.value = 'Регистрация через соцсети скоро — пока используйте email'
}

// Lightweight password strength meter (0–3 filled segments).
const strength = computed(() => {
  const p = password.value
  let score = 0
  if (p.length >= 8) score++
  if (/[A-ZА-Я]/.test(p) && /[a-zа-я]/.test(p)) score++
  if (/\d/.test(p) || /[^\w\s]/.test(p)) score++
  return score
})
const strengthLabel = computed(() => ['', 'слабый', 'средний', 'надёжный'][strength.value])
const strengthColor = computed(
  () => ['', 'text-danger', 'text-warning', 'text-success-bright'][strength.value],
)
</script>

<template>
  <div>
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

    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Создайте аккаунт</h1>
    <p class="mb-5 mt-1.5 text-sm text-ink-muted">Алерты, избранное и история обменов</p>

    <!-- socials -->
    <div class="mb-[18px] flex flex-col gap-2.5">
      <button
        type="button"
        class="flex w-full items-center justify-center gap-2.5 rounded-2xl bg-brand py-3.5 text-[15px] font-semibold text-white transition-colors hover:bg-brand-hover"
        @click="soon"
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
        Telegram
        <span class="rounded-md bg-white/20 px-2 py-[3px] text-[10px] font-semibold"
          >рекомендуем</span
        >
      </button>
      <div class="flex gap-2.5">
        <button
          type="button"
          class="flex flex-1 items-center justify-center gap-2.5 rounded-2xl border border-line-strong bg-surface py-3 text-sm font-semibold text-ink transition-colors hover:border-[#3A4047]"
          @click="soon"
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
          @click="soon"
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

    <!-- divider -->
    <div class="mb-[18px] flex items-center gap-3">
      <span class="h-px flex-1 bg-line" />
      <span class="text-xs text-ink-faint">или по email</span>
      <span class="h-px flex-1 bg-line" />
    </div>

    <!-- email -->
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

    <!-- password -->
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

    <!-- strength -->
    <div class="mb-[18px] mt-2.5 flex items-center gap-2">
      <span
        v-for="i in 3"
        :key="i"
        class="h-1 flex-1 rounded-full transition-colors"
        :class="i <= strength ? 'bg-success' : 'bg-line-strong'"
      />
      <span class="w-14 text-right text-xs font-semibold" :class="strengthColor">{{
        strengthLabel
      }}</span>
    </div>

    <!-- consent -->
    <label class="mb-[18px] flex cursor-pointer items-start gap-2.5">
      <input v-model="agreed" type="checkbox" class="peer sr-only" />
      <span
        class="mt-0.5 flex h-5 w-5 flex-none items-center justify-center rounded-md border-2 border-line-strong transition-colors peer-checked:border-brand peer-checked:bg-brand"
      >
        <svg
          v-if="agreed"
          width="12"
          height="12"
          viewBox="0 0 24 24"
          fill="none"
          stroke="#fff"
          stroke-width="3"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M5 12l5 5L20 6" />
        </svg>
      </span>
      <span class="text-[13px] leading-snug text-ink-muted">
        Соглашаюсь с <NuxtLink to="/terms" class="text-brand-bright">условиями</NuxtLink> и
        <NuxtLink to="/privacy" class="text-brand-bright">политикой конфиденциальности</NuxtLink>
      </span>
    </label>

    <p v-if="error" class="mb-2.5 text-[13px] text-danger">{{ error }}</p>
    <KButton block size="lg" class="!rounded-2xl" :disabled="!agreed || busy" @click="doRegister">{{
      busy ? 'Создаём…' : 'Создать аккаунт'
    }}</KButton>

    <!-- trust -->
    <div
      class="mt-4 flex items-center gap-3 rounded-2xl border border-success/25 bg-success/[0.07] px-[15px] py-3.5"
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
      <span class="text-[13px] leading-snug text-ink-muted"
        >Не запрашиваем доступ к кошелькам и приватным ключам</span
      >
    </div>

    <div class="mt-[18px] text-center text-sm text-ink-faint">
      Уже есть аккаунт?
      <NuxtLink to="/login" class="font-semibold text-brand-bright">Войти</NuxtLink>
    </div>
  </div>
</template>
