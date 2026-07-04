<script setup lang="ts">
import { reactive, ref } from 'vue'

definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Настройки — Kurso' })

const { user, login, logout } = useAuth()

const name = ref(user.value?.name ?? '')
const saved = ref(false)
function saveProfile() {
  if (!name.value.trim()) return
  // Persists to the session cookie (real, survives reloads).
  login({ name: name.value.trim(), email: user.value?.email, via: user.value?.via })
  saved.value = true
  setTimeout(() => (saved.value = false), 2000)
}

// Notification channels (functional in-session; server-side prefs land with the backend).
const channels = reactive({ telegram: true, email: false, push: false })

function signOut() {
  logout()
  navigateTo('/')
}
function removeAccount() {
  if (confirm('Удалить аккаунт и все алерты? Действие необратимо.')) {
    logout()
    navigateTo('/')
  }
}
</script>

<template>
  <div class="max-w-2xl">
    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Настройки</h1>
    <p class="mb-6 mt-1 text-sm text-ink-faint">Профиль, уведомления и безопасность</p>

    <!-- profile -->
    <section class="mb-5 rounded-2xl border border-line bg-surface p-5">
      <div class="mb-4 text-[15px] font-bold text-ink">Профиль</div>
      <label class="mb-3 block">
        <span class="mb-1.5 block text-xs text-ink-faint">Имя</span>
        <input
          v-model="name"
          class="w-full rounded-2xl border border-line bg-well px-[15px] py-3 text-[15px] text-ink outline-none focus:border-brand"
        />
      </label>
      <label class="mb-4 block">
        <span class="mb-1.5 block text-xs text-ink-faint">Email</span>
        <div
          class="tnum rounded-2xl border border-line bg-well px-[15px] py-3 text-[15px] text-ink-muted"
        >
          {{ user?.email ?? '—' }}
        </div>
      </label>
      <div class="flex items-center gap-3">
        <KButton size="sm" @click="saveProfile">Сохранить</KButton>
        <span v-if="saved" class="flex items-center gap-1.5 text-[13px] text-success-bright">
          <KStatusDot tone="success" :size="6" />Сохранено
        </span>
      </div>
    </section>

    <!-- notifications -->
    <section class="mb-5 rounded-2xl border border-line bg-surface p-5">
      <div class="mb-4 text-[15px] font-bold text-ink">Каналы уведомлений</div>
      <div class="flex flex-col gap-2.5">
        <div
          class="flex items-center justify-between rounded-2xl border border-line bg-well px-4 py-3"
        >
          <span class="text-sm text-ink">Telegram</span>
          <KToggle v-model="channels.telegram" />
        </div>
        <div
          class="flex items-center justify-between rounded-2xl border border-line bg-well px-4 py-3"
        >
          <span class="text-sm text-ink">Email</span>
          <KToggle v-model="channels.email" />
        </div>
        <div
          class="flex items-center justify-between rounded-2xl border border-line bg-well px-4 py-3"
        >
          <span class="text-sm text-ink">Push</span>
          <KToggle v-model="channels.push" />
        </div>
      </div>
    </section>

    <!-- security -->
    <section class="rounded-2xl border border-line bg-surface p-5">
      <div class="mb-4 text-[15px] font-bold text-ink">Аккаунт</div>
      <div class="flex flex-wrap gap-3">
        <KButton variant="secondary" size="sm" @click="signOut">Выйти</KButton>
        <KButton variant="danger" size="sm" @click="removeAccount">Удалить аккаунт</KButton>
      </div>
    </section>
  </div>
</template>
