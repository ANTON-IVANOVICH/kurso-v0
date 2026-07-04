import { onScopeDispose, watch, type Ref } from 'vue'

/** Locks body scroll while `open` is true. Client-only, self-cleaning. */
export function useScrollLock(open: Ref<boolean>) {
  if (!import.meta.client) return

  watch(open, (v) => {
    document.body.style.overflow = v ? 'hidden' : ''
  })
  onScopeDispose(() => {
    document.body.style.overflow = ''
  })
}
