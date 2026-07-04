<script setup lang="ts">
import { computed } from 'vue'

definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: 'Личный кабинет — Kurso' })

const { alerts, activeCount } = useAlerts()
const { favorites } = useFavorites()
const { history } = useHistory()
const { data: exchangers } = useExchangersQuery()

const onlineMap = computed(
  () => new Map((exchangers.value ?? []).map((e) => [e.slug, e.status === 'active'])),
)
const isOnline = (slug: string) => onlineMap.value.get(slug) ?? true

const monthCount = computed(() => {
  const since = Date.now() - 30 * 864e5
  return history.value.filter((h) => h.at >= since).length
})

const topAlerts = computed(() => alerts.value.slice(0, 4))
const topFavorites = computed(() => favorites.value.slice(0, 4))
const recentHistory = computed(() => history.value.slice(0, 3))
</script>

<template>
  <div>
    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">Обзор</h1>
    <p class="mb-6 mt-1 text-sm text-ink-faint">Сводка по вашим алертам, избранному и переходам</p>

    <!-- metrics -->
    <div class="mb-6 grid grid-cols-3 gap-2.5 md:gap-4">
      <div class="rounded-2xl border border-line bg-surface p-3.5 md:rounded-[18px] md:p-5">
        <div class="mb-2 flex items-center justify-between md:mb-3.5">
          <span class="text-[11px] text-ink-faint md:text-sm md:text-ink-muted"
            >Активные алерты</span
          >
          <span
            class="hidden h-[34px] w-[34px] items-center justify-center rounded-[10px] bg-brand/[0.12] text-brand-bright md:flex"
          >
            <svg
              width="17"
              height="17"
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
        </div>
        <div class="tnum text-2xl font-bold tracking-[-0.02em] text-ink md:text-[32px]">
          {{ activeCount }}
        </div>
        <div class="mt-1 text-[11px] text-ink-faint md:mt-1.5 md:text-xs">
          {{ alerts.length ? 'отслеживаются' : 'нет активных' }}
        </div>
      </div>

      <div class="rounded-2xl border border-line bg-surface p-3.5 md:rounded-[18px] md:p-5">
        <div class="mb-2 flex items-center justify-between md:mb-3.5">
          <span class="text-[11px] text-ink-faint md:text-sm md:text-ink-muted">Избранное</span>
          <span
            class="hidden h-[34px] w-[34px] items-center justify-center rounded-[10px] bg-warning/[0.12] text-warning md:flex"
          >
            <svg
              width="17"
              height="17"
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
        </div>
        <div class="tnum text-2xl font-bold tracking-[-0.02em] text-ink md:text-[32px]">
          {{ favorites.length }}
        </div>
        <div class="mt-1 text-[11px] text-ink-faint md:mt-1.5 md:text-xs">обменников</div>
      </div>

      <div class="rounded-2xl border border-line bg-surface p-3.5 md:rounded-[18px] md:p-5">
        <div class="mb-2 flex items-center justify-between md:mb-3.5">
          <span class="text-[11px] text-ink-faint md:text-sm md:text-ink-muted"
            >Переходы за месяц</span
          >
          <span
            class="hidden h-[34px] w-[34px] items-center justify-center rounded-[10px] bg-success/[0.12] text-success-bright md:flex"
          >
            <svg
              width="17"
              height="17"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M7 17 17 7M9 7h8v8" />
            </svg>
          </span>
        </div>
        <div class="tnum text-2xl font-bold tracking-[-0.02em] text-ink md:text-[32px]">
          {{ monthCount }}
        </div>
        <div class="mt-1 text-[11px] text-ink-faint md:mt-1.5 md:text-xs">переходов</div>
      </div>
    </div>

    <!-- alerts -->
    <div class="mb-3 flex items-center justify-between">
      <span class="text-[17px] font-bold text-ink">Мои алерты</span>
      <NuxtLink to="/account/alerts" class="text-[13px] font-semibold text-brand-bright"
        >Все алерты →</NuxtLink
      >
    </div>
    <div v-if="topAlerts.length" class="mb-7 grid gap-3 md:grid-cols-2">
      <AlertRow v-for="a in topAlerts" :key="a.id" :alert="a" />
    </div>
    <div
      v-else
      class="mb-7 flex flex-col items-center gap-3 rounded-2xl border border-dashed border-line-strong bg-surface/50 px-6 py-9 text-center"
    >
      <span
        class="flex h-11 w-11 items-center justify-center rounded-full bg-brand/[0.12] text-brand-bright"
      >
        <svg
          width="20"
          height="20"
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
      <p class="text-sm text-ink-muted">
        Пока нет алертов. Поставьте порог — пришлём, когда курс дойдёт.
      </p>
      <KButton size="sm" @click="navigateTo('/account/alerts/new')">Создать алерт</KButton>
    </div>

    <!-- favorites + history -->
    <div class="grid gap-6 md:grid-cols-[1fr_1.2fr]">
      <div>
        <div class="mb-3 flex items-center justify-between">
          <span class="text-[17px] font-bold text-ink">Избранные обменники</span>
          <NuxtLink
            to="/account/favorites"
            class="text-[13px] font-semibold text-brand-bright md:hidden"
            >Все →</NuxtLink
          >
        </div>
        <div v-if="topFavorites.length">
          <!-- desktop grid -->
          <div class="hidden grid-cols-2 gap-3 md:grid">
            <FavoriteCard
              v-for="f in topFavorites"
              :key="f.slug"
              :favorite="f"
              :online="isOnline(f.slug)"
            />
          </div>
          <!-- mobile carousel -->
          <div class="scrollx -mx-4 flex gap-2.5 overflow-x-auto px-4 pb-1 md:hidden">
            <FavoriteCard
              v-for="f in topFavorites"
              :key="f.slug"
              :favorite="f"
              :online="isOnline(f.slug)"
              vertical
            />
          </div>
        </div>
        <NuxtLink
          v-else
          to="/exchangers"
          class="flex flex-col items-center gap-2 rounded-2xl border border-dashed border-line-strong bg-surface/50 px-6 py-8 text-center transition-colors hover:border-line-strong"
        >
          <span
            class="flex h-10 w-10 items-center justify-center rounded-full bg-warning/[0.12] text-warning"
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
              <path
                d="M12 20s-6.5-4-9-8.5A4.5 4.5 0 0 1 12 7a4.5 4.5 0 0 1 9 4.5C18.5 16 12 20 12 20Z"
              />
            </svg>
          </span>
          <span class="text-sm text-ink-muted">Добавьте обменники в избранное</span>
        </NuxtLink>
      </div>

      <div>
        <div class="mb-3 text-[17px] font-bold text-ink">История переходов</div>
        <div
          v-if="recentHistory.length"
          class="overflow-hidden rounded-2xl border border-line bg-surface"
        >
          <AccountHistoryRow v-for="h in recentHistory" :key="h.id" :entry="h" />
        </div>
        <div
          v-else
          class="rounded-2xl border border-dashed border-line-strong bg-surface/50 px-6 py-8 text-center text-sm text-ink-faint"
        >
          Пока пусто — переходы в обменники появятся здесь.
        </div>
      </div>
    </div>
  </div>
</template>
