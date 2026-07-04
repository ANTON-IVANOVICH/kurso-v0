<script setup lang="ts">
definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Мои отзывы — Kurso' })
</script>

<template>
  <div class="max-w-2xl">
    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Мои отзывы</h1>
    <p class="mb-5 mt-1 text-sm text-ink-faint">Отзывы, которые вы оставили обменникам</p>

    <div
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
          <path d="M21 15a2 2 0 0 1-2 2H8l-4 4V6a2 2 0 0 1 2-2h13a2 2 0 0 1 2 2Z" />
        </svg>
      </span>
      <p class="max-w-xs text-sm text-ink-muted">
        Вы ещё не оставляли отзывов. Поделитесь опытом на странице обменника — это помогает другим.
      </p>
      <NuxtLink to="/exchangers" class="text-sm font-semibold text-brand-bright"
        >К обменникам →</NuxtLink
      >
    </div>
  </div>
</template>
