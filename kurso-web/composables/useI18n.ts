import { messages } from '~/i18n/messages'

// Lightweight, SSR-safe i18n over the global `useLocale` state. `t()` resolves a
// dotted key against the active locale's catalogue (falling back to ru, then the
// key itself) and interpolates {placeholders}. `plural()` picks the count-word
// form (Russian 3-form / English 2-form). Reading `locale.value` inside makes
// every template usage reactive to the language switch.
function resolve(dict: unknown, key: string): unknown {
  return key.split('.').reduce<unknown>((o, k) => (o as Record<string, unknown>)?.[k], dict)
}

export function useI18n() {
  const locale = useLocale()

  function t(key: string, params?: Record<string, string | number>): string {
    const active = messages[locale.value] ?? messages.ru
    let val = resolve(active, key)
    if (typeof val !== 'string') val = resolve(messages.ru, key)
    let str = typeof val === 'string' ? val : key
    if (params) {
      for (const [k, v] of Object.entries(params)) str = str.replaceAll(`{${k}}`, String(v))
    }
    return str
  }

  /** The count-word for `n` (word only, e.g. "обменников"). */
  function plural(n: number, key: keyof typeof messages.ru.plurals): string {
    const active = messages[locale.value] ?? messages.ru
    const forms = active.plurals[key] as readonly string[]
    if (locale.value === 'ru') {
      return pluralRu(n, forms[0], forms[1], forms[2])
    }
    return n === 1 ? forms[0] : forms[1]
  }

  return { t, plural, locale }
}
