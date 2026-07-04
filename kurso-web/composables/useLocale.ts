// Global UI-locale state. Full i18n (message catalogues) lands in a later
// stage; this holds the chosen language so the switcher works app-wide now.
export type Locale = 'ru' | 'en'

export interface LocaleOption {
  code: Locale
  label: string
  name: string
}

export const locales: LocaleOption[] = [
  { code: 'ru', label: 'RU', name: 'Русский' },
  { code: 'en', label: 'EN', name: 'English' },
]

export const useLocale = () => useState<Locale>('locale', () => 'ru')
