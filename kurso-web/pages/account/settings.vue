<script setup lang="ts">
import { reactive } from 'vue'

const { t } = useI18n()

definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: () => t('accountSettings.seoTitle') })

const { user, logout } = useAuth()

// Notification channels (functional in-session; server-side prefs land with the backend).
const channels = reactive({ telegram: true, email: false, push: false })

async function signOut() {
  await logout()
  navigateTo('/')
}
async function removeAccount() {
  if (confirm(t('accountSettings.confirmDelete'))) {
    await logout()
    navigateTo('/')
  }
}
</script>

<template>
  <div class="max-w-2xl">
    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">
      {{ t('accountSettings.title') }}
    </h1>
    <p class="mb-6 mt-1 text-sm text-ink-faint">{{ t('accountSettings.subtitle') }}</p>

    <!-- profile -->
    <section class="mb-5 rounded-2xl border border-line bg-surface p-5">
      <div class="mb-4 text-[15px] font-bold text-ink">{{ t('accountSettings.profile') }}</div>
      <div class="mb-3">
        <span class="mb-1.5 block text-xs text-ink-faint">{{ t('accountSettings.name') }}</span>
        <div class="rounded-2xl border border-line bg-well px-[15px] py-3 text-[15px] text-ink">
          {{ user?.name ?? '—' }}
        </div>
      </div>
      <div>
        <span class="mb-1.5 block text-xs text-ink-faint">Email</span>
        <div
          class="tnum rounded-2xl border border-line bg-well px-[15px] py-3 text-[15px] text-ink-muted"
        >
          {{ user?.email ?? '—' }}
        </div>
      </div>
    </section>

    <!-- notifications -->
    <section class="mb-5 rounded-2xl border border-line bg-surface p-5">
      <div class="mb-4 text-[15px] font-bold text-ink">
        {{ t('accountSettings.notificationChannels') }}
      </div>
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
      <div class="mb-4 text-[15px] font-bold text-ink">{{ t('accountSettings.account') }}</div>
      <div class="flex flex-wrap gap-3">
        <KButton variant="secondary" size="sm" @click="signOut">{{
          t('accountSettings.signOut')
        }}</KButton>
        <KButton variant="danger" size="sm" @click="removeAccount">{{
          t('accountSettings.deleteAccount')
        }}</KButton>
      </div>
    </section>
  </div>
</template>
