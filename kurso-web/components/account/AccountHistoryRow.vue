<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { timeAgo, useHistory, type HistoryEntry } from '~/composables/useHistory'

const props = defineProps<{ entry: HistoryEntry }>()
const apiBase = useApiBase()
const { record } = useHistory()

// Client clock so the relative time doesn't drift SSR↔client.
const now = ref(props.entry.at)
onMounted(() => (now.value = Date.now()))
const ago = computed(() => timeAgo(props.entry.at, now.value))

const href = computed(
  () => `${apiBase}/go/${props.entry.slug}?direction=${props.entry.directionSlug}`,
)

// "снова" is a real outbound jump — log it again before leaving.
function again() {
  record({
    pair: props.entry.pair,
    exchanger: props.entry.exchanger,
    slug: props.entry.slug,
    directionSlug: props.entry.directionSlug,
    amount: props.entry.amount,
  })
}
</script>

<template>
  <div
    class="grid grid-cols-[1fr_auto] items-center gap-3 border-b border-line px-4 py-3 last:border-0 md:grid-cols-[1.5fr_1fr_1fr_auto]"
  >
    <div class="min-w-0">
      <div class="truncate text-sm font-semibold text-ink">{{ entry.pair }}</div>
      <div class="truncate text-[11px] text-ink-faint">
        {{ entry.exchanger }}<span class="md:hidden"> · {{ ago }}</span>
      </div>
    </div>
    <div class="tnum hidden text-[13px] text-ink-muted md:block">{{ entry.amount }}</div>
    <div class="tnum hidden text-xs text-ink-faint md:block">{{ ago }}</div>
    <a
      :href="href"
      target="_blank"
      rel="noopener nofollow"
      class="flex-none whitespace-nowrap text-[13px] font-semibold text-brand-bright"
      @click="again"
      >снова →</a
    >
  </div>
</template>
