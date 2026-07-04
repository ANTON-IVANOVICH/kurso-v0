// Sidebar navigation for the merchant cabinet, matching the mockup order.
// Badge counts are filled from live cabinet data (unanswered reviews, stale
// feeds, open complaints) rather than hardcoded.

export type IconName = 'grid' | 'arrowRight' | 'chat' | 'bars' | 'user' | 'flag' | 'card'

export type BadgeTone = 'warn' | 'danger' | 'brand'

export interface NavItem {
  to: string
  label: string
  icon: IconName
  badge?: 'ratesStale' | 'reviewsUnanswered' | 'complaintsOpen'
}

export const NAV: NavItem[] = [
  { to: '/', label: 'Главная', icon: 'grid' },
  { to: '/rates', label: 'Курсы', icon: 'arrowRight', badge: 'ratesStale' },
  { to: '/reviews', label: 'Отзывы', icon: 'chat', badge: 'reviewsUnanswered' },
  { to: '/traffic', label: 'Трафик', icon: 'bars' },
  { to: '/profile', label: 'Профиль', icon: 'user' },
  { to: '/complaints', label: 'Жалобы', icon: 'flag', badge: 'complaintsOpen' },
  { to: '/billing', label: 'Биллинг', icon: 'card' },
]
