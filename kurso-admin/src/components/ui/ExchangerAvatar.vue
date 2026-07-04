<script setup lang="ts">
import { computed } from 'vue'
import { exchangerAvatar } from '../../lib/avatar'

const props = withDefaults(
  defineProps<{ slug: string; name: string; size?: number; shape?: 'circle' | 'rounded' }>(),
  { size: 28, shape: 'circle' },
)

const av = computed(() => exchangerAvatar(props.slug, props.name))
const fontSize = computed(() => {
  const n = av.value.initials.length
  return Math.round(props.size * (n > 2 ? 0.3 : n > 1 ? 0.36 : 0.42))
})
const radius = computed(() =>
  props.shape === 'circle' ? '50%' : `${Math.round(props.size * 0.26)}px`,
)
</script>

<template>
  <span
    class="flex flex-none items-center justify-center font-bold text-white"
    :style="{
      width: `${size}px`,
      height: `${size}px`,
      background: av.color,
      borderRadius: radius,
      fontSize: `${fontSize}px`,
    }"
    >{{ av.initials }}</span
  >
</template>
