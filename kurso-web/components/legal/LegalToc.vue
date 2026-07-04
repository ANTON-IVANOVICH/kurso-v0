<script setup lang="ts">
// Desktop-only "Содержание" sidebar for legal pages. Anchor links to sections,
// with the in-view section highlighted (scroll-spy driven by the parent page).
defineProps<{
  sections: { id: string; num: string; title: string }[]
  activeId: string
}>()
</script>

<template>
  <aside class="hidden md:block">
    <div class="sticky top-24">
      <div class="font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint">
        Содержание
      </div>
      <nav class="mt-3.5 flex flex-col gap-0.5 text-sm">
        <a
          v-for="s in sections"
          :key="s.id"
          :href="`#${s.id}`"
          class="rounded-lg px-3 py-2 leading-snug transition-colors"
          :class="
            s.id === activeId
              ? 'bg-brand/10 font-semibold text-brand-bright'
              : 'text-ink-muted hover:text-ink'
          "
          >{{ s.num }}. {{ s.title }}</a
        >
      </nav>
    </div>
  </aside>
</template>
