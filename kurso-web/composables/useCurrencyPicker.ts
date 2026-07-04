export interface CurrencyOption {
  code: string
  name: string
  subtitle: string
  symbol: string
  color: string
  /** Symbol drawn in dark ink (for light backgrounds, e.g. Тинькофф yellow). */
  dark?: boolean
  /** Extra search terms (RU + EN) so the filter works in both languages. */
  aliases?: string
}

// Currencies and banks offered in the exchange calculator. The same list backs
// the drawer everywhere it opens (home, alert builder, map, widgets).
export const currencyOptions: CurrencyOption[] = [
  {
    code: 'USDT',
    name: 'USDT',
    subtitle: 'Tether',
    symbol: '₮',
    color: '#26A17B',
    aliases: 'юсдт тезер стейбл',
  },
  {
    code: 'BTC',
    name: 'BTC',
    subtitle: 'Bitcoin',
    symbol: '₿',
    color: '#F7931A',
    aliases: 'биткоин битка',
  },
  {
    code: 'ETH',
    name: 'ETH',
    subtitle: 'Ethereum',
    symbol: 'Ξ',
    color: '#627EEA',
    aliases: 'эфир эфириум',
  },
  {
    code: 'TON',
    name: 'TON',
    subtitle: 'Toncoin',
    symbol: 'TON',
    color: '#0098EA',
    aliases: 'тон тонкоин',
  },
  {
    code: 'Тинькофф',
    name: 'Тинькофф',
    subtitle: 'RUB · банк',
    symbol: 'Т',
    color: '#FFDD2D',
    dark: true,
    aliases: 'tinkoff тбанк т-банк рубль',
  },
  {
    code: 'Сбербанк',
    name: 'Сбербанк',
    subtitle: 'RUB · банк',
    symbol: 'С',
    color: '#1FAE54',
    aliases: 'sber сбер рубль',
  },
]

export type PickerSide = 'give' | 'get'

/**
 * Global controller for the currency-picker bottom sheet. Any selector across
 * the product opens the same drawer; it writes the choice back to the
 * calculator's give/get state.
 */
export function useCurrencyPicker() {
  const open = useState('currencyPickerOpen', () => false)
  const side = useState<PickerSide>('currencyPickerSide', () => 'give')
  const give = useState<CurrencyOption>('calcGive', () => currencyOptions[0])
  const get = useState<CurrencyOption>('calcGet', () => currencyOptions[4])

  function openPicker(which: PickerSide) {
    side.value = which
    open.value = true
  }
  function choose(option: CurrencyOption) {
    if (side.value === 'give') give.value = option
    else get.value = option
    open.value = false
  }
  function close() {
    open.value = false
  }
  function swap() {
    const prev = give.value
    give.value = get.value
    get.value = prev
  }

  return { open, side, give, get, currencyOptions, openPicker, choose, close, swap }
}
