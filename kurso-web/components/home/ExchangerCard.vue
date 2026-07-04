<script setup lang="ts">
import { computed } from 'vue'

type State = 'best' | 'normal' | 'stale'
type DeltaTone = 'success' | 'danger' | 'muted'

const props = withDefaults(
  defineProps<{
    state?: State
    name: string
    initials: string
    avatarColor: string
    rating: string
    reviews: string
    time: string
    reserve: string
    extraTag?: string
    partner?: boolean
    note?: string
    receive: string
    rate: string
    delta: string
    deltaTone?: DeltaTone
    href?: string
  }>(),
  {
    state: 'normal',
    partner: false,
    extraTag: undefined,
    note: undefined,
    deltaTone: undefined,
    href: undefined,
  },
)

const accentCls = computed(
  () => ({ best: 'bg-brand', normal: 'bg-line-strong', stale: 'bg-[#8A6A2B]' })[props.state],
)

const shellCls = computed(
  () =>
    ({
      best: 'border border-brand/35 bg-[linear-gradient(100deg,rgba(46,125,242,0.10),rgba(46,125,242,0.03))]',
      normal: 'border border-line bg-surface',
      stale: 'border border-line bg-surface opacity-60',
    })[props.state],
)

const tone = computed<DeltaTone>(
  () =>
    props.deltaTone ??
    (props.state === 'best' ? 'success' : props.state === 'stale' ? 'muted' : 'danger'),
)

const deltaCls = computed(
  () =>
    ({
      success: 'text-success bg-success/[0.14]',
      danger: 'text-danger-soft bg-danger-soft/[0.12]',
      muted: 'text-ink-faint bg-surface-raised',
    })[tone.value],
)

const { t } = useI18n()
const deltaNumeric = computed(() => props.delta.includes('₽'))
const receiveCls = computed(() => (props.state === 'stale' ? 'text-ink-faint' : 'text-ink'))
const ratingCls = computed(() => (props.state === 'stale' ? 'text-ink-faint' : 'text-ink-muted'))
const starCls = computed(() => (props.state === 'stale' ? 'text-[#8A6A2B]' : 'text-warning-deep'))
const chipBg = computed(() => (props.state === 'best' ? 'bg-white/[0.04]' : 'bg-surface-raised'))
const chipBase =
  'whitespace-nowrap rounded-sm border border-line-strong px-[9px] py-1 text-xs text-ink-muted'
const btnVariant = computed(() => (props.state === 'best' ? 'primary' : 'secondary'))
</script>

<template>
  <article class="relative rounded-[18px] p-4 md:p-[20px_22px]" :class="shellCls">
    <span
      class="absolute bottom-4 left-0 top-4 hidden w-[3px] rounded-[3px] md:block"
      :class="accentCls"
    />

    <!-- ===== Desktop ===== -->
    <div class="hidden flex-wrap items-center gap-x-5 gap-y-4 md:flex">
      <div class="flex min-w-0 flex-1 items-center gap-[15px]">
        <span
          class="flex h-12 w-12 flex-none items-center justify-center rounded-[15px] text-base font-bold text-white"
          :style="{ background: avatarColor }"
          >{{ initials }}</span
        >
        <div class="min-w-0">
          <div class="flex items-center gap-2 text-base font-semibold">
            {{ name }}
            <KBadge v-if="partner" tone="brand">{{ t('exchangers.partner') }}</KBadge>
            <span v-if="note" class="text-[11px] text-warning">{{ note }}</span>
          </div>
          <div class="mt-1.5 flex flex-wrap items-center gap-2 text-[13px]" :class="ratingCls">
            <span :class="starCls">★</span>
            <span class="tnum">{{ rating }}</span>
            <span class="text-line-strong">·</span>
            <span class="tnum">{{ reviews }}</span>
            <span class="text-line-strong">·</span>
            <span class="whitespace-nowrap">{{ time }}</span>
          </div>
          <div class="mt-2.5 flex gap-[7px]">
            <span :class="[chipBase, chipBg]"
              >{{ t('home.card.reserve') }} <span class="tnum">{{ reserve }}</span></span
            >
            <span v-if="extraTag" :class="[chipBase, chipBg]">{{ extraTag }}</span>
          </div>
        </div>
      </div>

      <div class="min-w-[150px] text-right">
        <div class="mb-[3px] text-xs text-ink-faint">{{ t('home.card.youReceive') }}</div>
        <div class="tnum text-[25px] font-bold tracking-[-0.01em]" :class="receiveCls">
          {{ receive }}
        </div>
        <div class="mt-1.5 flex items-center justify-end gap-2">
          <span class="tnum text-[13px] text-ink-faint">{{ rate }}</span>
          <span
            class="rounded-[6px] px-2 py-0.5 text-xs font-semibold"
            :class="[deltaCls, { tnum: deltaNumeric }]"
            >{{ delta }}</span
          >
        </div>
      </div>

      <KButton :variant="btnVariant" size="lg" class="flex-none" :href="href">{{
        t('home.card.go')
      }}</KButton>
    </div>

    <!-- ===== Mobile ===== -->
    <div class="md:hidden">
      <div class="flex items-center gap-[11px]">
        <span
          class="flex h-10 w-10 flex-none items-center justify-center rounded-[13px] text-sm font-bold text-white"
          :style="{ background: avatarColor }"
          >{{ initials }}</span
        >
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-[7px] text-[15px] font-semibold">
            {{ name }}
            <KBadge v-if="partner" tone="brand">{{ t('exchangers.partner') }}</KBadge>
            <span v-if="note" class="text-[10px] text-warning">{{ note }}</span>
          </div>
          <div class="mt-[3px] flex items-center gap-1.5 text-xs" :class="ratingCls">
            <span :class="starCls">★</span>
            <span class="tnum">{{ rating }}</span>
            <span class="text-line-strong">·</span>
            <span class="tnum">{{ reviews }}</span>
            <span class="text-line-strong">·</span>
            <span class="whitespace-nowrap">{{ time }}</span>
          </div>
        </div>
      </div>

      <div class="mb-3 mt-3.5 flex items-end justify-between">
        <div>
          <div class="mb-0.5 text-xs text-ink-faint">{{ t('home.card.youReceive') }}</div>
          <div
            class="tnum text-[27px] font-bold leading-none tracking-[-0.01em]"
            :class="receiveCls"
          >
            {{ receive }}
          </div>
        </div>
        <span
          class="rounded-[8px] px-[9px] py-1 text-xs font-semibold"
          :class="[deltaCls, { tnum: deltaNumeric }]"
          >{{ delta }}</span
        >
      </div>

      <div class="mb-3.5 flex flex-wrap items-center gap-[7px]">
        <span :class="[chipBase, chipBg]"
          >{{ t('home.card.reserve') }} <span class="tnum">{{ reserve }}</span></span
        >
        <span v-if="extraTag" :class="[chipBase, chipBg]">{{ extraTag }}</span>
      </div>

      <KButton :variant="btnVariant" size="lg" block :href="href">{{
        t('home.card.goMobile')
      }}</KButton>
    </div>
  </article>
</template>
