import { watch } from 'vue'

// Persists the user's alerts / favorites / clickout history to localStorage so
// they survive reloads (the session itself is a cookie — see useAuth). Loading
// happens on `app:mounted`, i.e. AFTER hydration, so the first client render
// still matches the (empty) server render and no hydration mismatch occurs; the
// stored data then flows in reactively.
const KEYS = ['alerts', 'favorites', 'history'] as const

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook('app:mounted', () => {
    for (const key of KEYS) {
      const state = useState<unknown[]>(key, () => [])
      try {
        const raw = localStorage.getItem(`kurso_${key}`)
        if (raw) state.value = JSON.parse(raw)
      } catch {
        /* corrupt/absent — keep empty */
      }
      watch(
        state,
        (value) => {
          try {
            localStorage.setItem(`kurso_${key}`, JSON.stringify(value))
          } catch {
            /* quota / private mode — ignore */
          }
        },
        { deep: true },
      )
    }
  })
})
