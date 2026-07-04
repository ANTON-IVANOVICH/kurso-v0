<script setup lang="ts">
definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Алерты — Kurso' })

const { alerts, pause, resume, remove } = useAlerts()
</script>

<template>
  <div class="max-w-2xl">
    <div class="mb-5 flex items-center justify-between gap-3">
      <div>
        <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Алерты</h1>
        <p class="mt-1 text-sm text-ink-faint">Пришлём в Telegram, когда курс дойдёт до порога</p>
      </div>
      <KButton pill @click="navigateTo('/account/alerts/new')">+ Новый алерт</KButton>
    </div>

    <div v-if="alerts.length" class="flex flex-col gap-3">
      <AlertRow
        v-for="a in alerts"
        :key="a.id"
        :alert="a"
        actions
        @pause="pause(a.id)"
        @resume="resume(a.id)"
        @remove="remove(a.id)"
      />
    </div>

    <div
      v-else
      class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-line-strong bg-surface/50 px-6 py-12 text-center"
    >
      <span
        class="flex h-12 w-12 items-center justify-center rounded-full bg-brand/[0.12] text-brand-bright"
      >
        <svg
          width="22"
          height="22"
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
      <p class="max-w-xs text-sm text-ink-muted">
        Пока нет алертов. Поставьте порог по нужной паре — пришлём уведомление, когда курс его
        достигнет.
      </p>
      <KButton @click="navigateTo('/account/alerts/new')">Создать первый алерт</KButton>
    </div>
  </div>
</template>
