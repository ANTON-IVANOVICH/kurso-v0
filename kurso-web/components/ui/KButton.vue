<script setup lang="ts">
import { computed } from 'vue'

type Variant = 'primary' | 'secondary' | 'ghost' | 'success' | 'danger'
type Size = 'sm' | 'md' | 'lg'

const props = withDefaults(
  defineProps<{
    variant?: Variant
    size?: Size
    pill?: boolean
    block?: boolean
    disabled?: boolean
    // When set, renders an <a> (opens in a new tab) instead of a <button> — used
    // for outbound clickout links (/go/{slug}).
    href?: string
  }>(),
  { variant: 'primary', size: 'md', pill: false, block: false, disabled: false, href: undefined },
)

const variants: Record<Variant, string> = {
  primary: 'bg-brand text-white hover:bg-brand-hover',
  secondary: 'bg-surface-raised text-ink border border-line-strong hover:border-[#3A4047]',
  ghost: 'bg-transparent text-ink-muted hover:text-ink',
  success: 'bg-success/[0.12] text-success-bright border border-success/30 hover:bg-success/20',
  danger: 'bg-danger/[0.12] text-danger border border-danger/30 hover:bg-danger/20',
}
const sizes: Record<Size, string> = {
  sm: 'text-xs px-[13px] py-2 rounded',
  md: 'text-sm px-[18px] py-[11px] rounded-lg',
  lg: 'text-base px-6 py-[13px] rounded-xl',
}

const classes = computed(() => [
  'inline-flex items-center justify-center gap-2 font-semibold transition-colors',
  'disabled:cursor-not-allowed disabled:opacity-60',
  variants[props.variant],
  sizes[props.size],
  props.pill ? '!rounded-full' : '',
  props.block ? 'w-full' : '',
])
</script>

<template>
  <a v-if="href" :href="href" target="_blank" rel="noopener nofollow" :class="classes">
    <slot />
  </a>
  <button v-else :class="classes" :disabled="disabled">
    <slot />
  </button>
</template>
