<script setup lang="ts">
import type { Favorite } from '~/composables/useFavorites'

const { t } = useI18n()

withDefaults(defineProps<{ favorite: Favorite; online?: boolean; vertical?: boolean }>(), {
  online: true,
  vertical: false,
})
</script>

<template>
  <NuxtLink
    :to="`/exchangers/${favorite.slug}`"
    class="rounded-2xl border border-line bg-surface p-3.5 transition-colors hover:border-line-strong"
    :class="[
      vertical ? 'block w-[130px] flex-none' : 'flex items-center gap-2.5',
      !online ? 'opacity-60' : '',
    ]"
  >
    <span
      class="flex flex-none items-center justify-center rounded-full font-bold text-white"
      :class="[
        vertical ? 'mb-2.5 h-9 w-9 text-[13px]' : 'h-[34px] w-[34px] text-xs',
        favorite.initials.length > 1 ? '!text-[11px]' : '',
      ]"
      :style="{ backgroundColor: favorite.color }"
      >{{ favorite.initials }}</span
    >
    <div class="min-w-0" :class="vertical ? '' : 'flex-1'">
      <div class="truncate text-sm font-semibold text-ink">{{ favorite.name }}</div>
      <div
        class="mt-0.5 flex items-center gap-1.5 text-[11px]"
        :class="online ? 'text-success-bright' : 'text-ink-faint'"
      >
        <span
          class="h-1.5 w-1.5 rounded-full"
          :style="{ backgroundColor: online ? '#2BC58C' : '#5A616A' }"
        />{{ online ? t('favoriteCard.online') : t('favoriteCard.offline') }}
      </div>
    </div>
  </NuxtLink>
</template>
