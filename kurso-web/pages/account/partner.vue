<script setup lang="ts">
definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Партнёрка — Kurso' })
</script>

<template>
  <div class="max-w-2xl">
    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Партнёрка</h1>
    <p class="mb-5 mt-1 text-sm text-ink-faint">Зарабатывайте, приводя пользователей в Kurso</p>

    <div class="rounded-2xl border border-line bg-surface p-6">
      <span
        class="mb-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-brand/[0.12] text-brand-bright"
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
          <circle cx="12" cy="12" r="9" />
          <path d="M9.5 9a2.5 2.5 0 0 1 4.5 1.5c0 1.5-2.5 2-2.5 3.5M12 17h.01" />
        </svg>
      </span>
      <div class="text-lg font-bold text-ink">Партнёрская программа скоро</div>
      <p class="mt-1.5 max-w-md text-sm leading-relaxed text-ink-muted">
        Реферальные ссылки, доход с переходов и выплаты появятся в личном кабинете. Оставьте алерт
        или добавьте обменники в избранное — а мы сообщим о запуске.
      </p>
      <div class="mt-4 flex flex-wrap gap-3">
        <KButton size="sm" @click="navigateTo('/account/alerts/new')">Создать алерт</KButton>
        <KButton variant="secondary" size="sm" @click="navigateTo('/exchangers')"
          >К обменникам</KButton
        >
      </div>
    </div>
  </div>
</template>
