<script setup lang="ts">
import { computed } from 'vue'

definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Избранное — Kurso' })

const { favorites, remove } = useFavorites()
const { data: exchangers } = useExchangersQuery()
const onlineMap = computed(
  () => new Map((exchangers.value ?? []).map((e) => [e.slug, e.status === 'active'])),
)
</script>

<template>
  <div class="max-w-3xl">
    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Избранное</h1>
    <p class="mb-5 mt-1 text-sm text-ink-faint">Обменники, за которыми вы следите</p>

    <div v-if="favorites.length" class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="f in favorites" :key="f.slug" class="group relative">
        <FavoriteCard :favorite="f" :online="onlineMap.get(f.slug) ?? true" />
        <button
          type="button"
          aria-label="Убрать из избранного"
          class="absolute right-2.5 top-2.5 flex h-7 w-7 items-center justify-center rounded-full bg-canvas/70 text-danger opacity-0 transition-opacity hover:bg-canvas group-hover:opacity-100"
          @click.prevent="remove(f.slug)"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
            <path
              d="M12 20s-6.5-4-9-8.5A4.5 4.5 0 0 1 12 7a4.5 4.5 0 0 1 9 4.5C18.5 16 12 20 12 20Z"
            />
          </svg>
        </button>
      </div>
    </div>

    <NuxtLink
      v-else
      to="/exchangers"
      class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-line-strong bg-surface/50 px-6 py-12 text-center transition-colors hover:border-line-strong"
    >
      <span
        class="flex h-12 w-12 items-center justify-center rounded-full bg-warning/[0.12] text-warning"
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
          <path
            d="M12 20s-6.5-4-9-8.5A4.5 4.5 0 0 1 12 7a4.5 4.5 0 0 1 9 4.5C18.5 16 12 20 12 20Z"
          />
        </svg>
      </span>
      <p class="max-w-xs text-sm text-ink-muted">
        Пусто. Откройте каталог обменников и нажмите на сердечко, чтобы добавить сюда.
      </p>
      <span class="text-sm font-semibold text-brand-bright">К обменникам →</span>
    </NuxtLink>
  </div>
</template>
