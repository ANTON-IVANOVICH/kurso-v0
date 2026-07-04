<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import UiIcon from '../components/ui/UiIcon.vue'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import InputOtp from 'primevue/inputotp'
import Button from 'primevue/button'
import Message from 'primevue/message'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const email = ref('')
const password = ref('')
const otp = ref('')
const busy = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (!email.value.trim() || !password.value) {
    error.value = 'Введите email и пароль'
    return
  }
  busy.value = true
  try {
    await auth.login(email.value.trim(), password.value, otp.value || undefined)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.replace(redirect)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Не удалось войти'
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="relative min-h-screen overflow-hidden">
    <div
      class="pointer-events-none absolute left-1/2 top-[-300px] h-[640px] w-[1000px] -translate-x-1/2"
      style="
        background: radial-gradient(
          ellipse at center,
          rgba(46, 125, 242, 0.14),
          rgba(46, 125, 242, 0) 62%
        );
      "
    />

    <div class="relative mx-auto flex min-h-screen max-w-md flex-col justify-center px-6 py-16">
      <div class="mb-8 flex items-center gap-3">
        <div
          class="flex h-11 w-11 items-center justify-center rounded-xl text-white shadow-[0_8px_24px_rgba(46,125,242,0.35)]"
          style="background: linear-gradient(150deg, #4a90f5, #2e7df2)"
        >
          <UiIcon name="logo" :size="20" :stroke="2.1" />
        </div>
        <div>
          <div class="text-lg font-extrabold tracking-tight">Kurso Admin</div>
          <div class="text-xs text-ink-faint">Панель управления каталогом</div>
        </div>
      </div>

      <form
        class="rounded-2xl border border-line bg-surface p-6 shadow-[0_30px_80px_rgba(0,0,0,0.5)]"
        @submit.prevent="submit"
      >
        <div class="mb-5 flex items-center gap-2.5">
          <span class="h-0.5 w-4 rounded-sm bg-brand" />
          <span class="eyebrow">Вход</span>
        </div>

        <label class="mb-4 block">
          <span class="mb-1.5 block text-xs text-ink-faint">Email</span>
          <InputText
            v-model="email"
            type="email"
            autocomplete="username"
            placeholder="admin@kurso.io"
            fluid
          />
        </label>

        <label class="mb-4 block">
          <span class="mb-1.5 block text-xs text-ink-faint">Пароль</span>
          <Password
            v-model="password"
            :feedback="false"
            toggle-mask
            autocomplete="current-password"
            placeholder="••••••••"
            fluid
          />
        </label>

        <div class="mb-5">
          <span class="mb-1.5 block text-xs text-ink-faint">
            Код 2FA <span class="text-ink-fainter">· если включён</span>
          </span>
          <InputOtp v-model="otp" :length="6" integer-only />
        </div>

        <Message v-if="error" severity="error" class="mb-4" size="small">{{ error }}</Message>

        <Button
          type="submit"
          :label="busy ? 'Вход…' : 'Войти'"
          icon="pi pi-sign-in"
          :loading="busy"
          fluid
        />

        <p class="mt-4 text-center text-[11px] leading-relaxed text-ink-fainter">
          Тестовый доступ: <span class="tnum text-ink-faint">admin@kurso.io</span> · пароль
          <span class="tnum text-ink-faint">admin12345</span>
        </p>
      </form>
    </div>
  </div>
</template>
