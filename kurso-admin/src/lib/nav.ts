// Sidebar navigation, grouped exactly as the dashboard mockup. `dynamic` badges
// are filled from live store data (e.g. the real exchanger count); the rest are
// demo counts from fixtures until their backends exist.

export type BadgeTone = 'faint' | 'warn' | 'danger' | 'brand'

export interface NavItem {
  to: string
  label: string
  icon?: 'grid'
  badge?: { text: string; tone: BadgeTone }
  dynamic?: 'exchangers'
}

export interface NavGroup {
  title: string
  items: NavItem[]
}

export const NAV: NavGroup[] = [
  {
    title: 'Обзор',
    items: [{ to: '/', label: 'Дашборд', icon: 'grid' }],
  },
  {
    title: 'Каталог',
    items: [
      { to: '/exchangers', label: 'Обменники', dynamic: 'exchangers' },
      { to: '/currencies', label: 'Валюты', badge: { text: '3 новые', tone: 'warn' } },
      { to: '/directions', label: 'Направления' },
      { to: '/parsers', label: 'Парсеры', badge: { text: '3 ↓', tone: 'danger' } },
    ],
  },
  {
    title: 'Модерация',
    items: [
      { to: '/reviews', label: 'Отзывы', badge: { text: '14', tone: 'brand' } },
      { to: '/complaints', label: 'Жалобы', badge: { text: '3', tone: 'danger' } },
    ],
  },
  {
    title: 'Монетизация · Данные',
    items: [
      { to: '/partners', label: 'Партнёрка' },
      { to: '/payouts', label: 'Выплаты' },
      { to: '/export', label: 'Экспорт · Логи' },
    ],
  },
]
