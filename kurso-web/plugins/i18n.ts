import { watch } from 'vue'
import type { Locale } from '~/composables/useLocale'

// Persists the chosen language in a cookie (survives reload, read on SSR) and
// keeps <html lang> in sync. Runs once, universally; useLocale() is shared
// app-wide via useState so all components react to the switch.
export default defineNuxtPlugin(() => {
  const locale = useLocale()
  const cookie = useCookie<Locale>('kurso_lang', {
    maxAge: 60 * 60 * 24 * 365,
    sameSite: 'lax',
    path: '/',
  })

  // Initialise from the cookie (SSR reads the request cookie; client reads it too).
  if (cookie.value && cookie.value !== locale.value) {
    locale.value = cookie.value
  }
  // Persist future changes.
  watch(locale, (v) => {
    cookie.value = v
  })

  // Reactive document language for a11y / SEO.
  useHead({ htmlAttrs: { lang: locale } })
})
