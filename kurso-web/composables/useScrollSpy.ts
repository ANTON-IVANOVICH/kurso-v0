import { onBeforeUnmount, onMounted, ref } from 'vue'

/**
 * Tracks which section (by element id) is currently in view, for a table of
 * contents. Returns the id of the topmost visible section. Client-only —
 * on the server it simply returns the first id.
 */
export function useScrollSpy(ids: string[], topOffset = 120) {
  const activeId = ref(ids[0] ?? '')

  if (import.meta.client) {
    let observer: IntersectionObserver | null = null
    const visible = new Set<string>()

    onMounted(() => {
      observer = new IntersectionObserver(
        (entries) => {
          for (const entry of entries) {
            if (entry.isIntersecting) visible.add(entry.target.id)
            else visible.delete(entry.target.id)
          }
          // First section (in document order) that is currently visible.
          const first = ids.find((id) => visible.has(id))
          if (first) activeId.value = first
        },
        { rootMargin: `-${topOffset}px 0px -55% 0px`, threshold: 0 },
      )
      for (const id of ids) {
        const el = document.getElementById(id)
        if (el) observer.observe(el)
      }
    })

    onBeforeUnmount(() => observer?.disconnect())
  }

  return { activeId }
}
