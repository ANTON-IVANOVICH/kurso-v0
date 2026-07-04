// Sample data for the admin sections that have no backend yet (parser health,
// revenue, moderation queue, complaints, payouts, audit log, exports). It mirrors
// the design mockups so the console is fully navigable; every consumer treats it
// as demo data. Once the admin API lands these become live queries.

export type HealthTone = 'success' | 'warn' | 'danger' | 'muted' | 'brand'

export interface ParserRow {
  name: string
  exchanger: string
  format: 'XML' | 'JSON' | 'HTML'
  interval: string
  lastRun: string
  latency: string
  success: string
  successTone: HealthTone
  lastRunTone: HealthTone
  dot: HealthTone
  status: string
  statusTone: HealthTone
}

export const PARSERS: ParserRow[] = [
  {
    name: 'main-rates',
    exchanger: 'CryptoBridge',
    format: 'XML',
    interval: '10с',
    lastRun: '12 сек',
    lastRunTone: 'muted',
    latency: '240мс',
    success: '99.8%',
    successTone: 'success',
    dot: 'success',
    status: 'healthy',
    statusTone: 'success',
  },
  {
    name: 'netex-feed',
    exchanger: 'NetEx24',
    format: 'JSON',
    interval: '15с',
    lastRun: '8 сек',
    lastRunTone: 'muted',
    latency: '180мс',
    success: '99.4%',
    successTone: 'success',
    dot: 'success',
    status: 'healthy',
    statusTone: 'success',
  },
  {
    name: 'reserves-eth',
    exchanger: 'NetEx24',
    format: 'JSON',
    interval: '30с',
    lastRun: '30 сек',
    lastRunTone: 'muted',
    latency: '410мс',
    success: '99.1%',
    successTone: 'success',
    dot: 'success',
    status: 'healthy',
    statusTone: 'success',
  },
  {
    name: 'reserves-btc',
    exchanger: 'BaksMan',
    format: 'XML',
    interval: '30с',
    lastRun: '9 мин',
    lastRunTone: 'warn',
    latency: '2.1с',
    success: '88.4%',
    successTone: 'warn',
    dot: 'warn',
    status: 'деградация',
    statusTone: 'warn',
  },
  {
    name: 'feed-ton',
    exchanger: 'TON-Exchange',
    format: 'JSON',
    interval: '10с',
    lastRun: 'timeout',
    lastRunTone: 'danger',
    latency: '—',
    success: '0%',
    successTone: 'danger',
    dot: 'danger',
    status: 'упал · 9м',
    statusTone: 'danger',
  },
  {
    name: 'alfa-cash',
    exchanger: '24Paybank',
    format: 'HTML',
    interval: '60с',
    lastRun: '14 мин',
    lastRunTone: 'danger',
    latency: '—',
    success: '0%',
    successTone: 'danger',
    dot: 'danger',
    status: 'упал · 14м',
    statusTone: 'danger',
  },
]

export const PARSER_TALLY = { healthy: 116, degraded: 9, down: 3, total: 128 }

export interface RevenueRow {
  rank: number
  name: string
  revenue: number
}

export const TOP_REVENUE: RevenueRow[] = [
  { rank: 1, name: 'CryptoBridge', revenue: 412_000 },
  { rank: 2, name: 'NetEx24', revenue: 388_000 },
  { rank: 3, name: '24Paybank', revenue: 214_000 },
  { rank: 4, name: 'BaksMan', revenue: 96_000 },
]

export interface AttentionItem {
  tone: HealthTone
  text: string
  detail?: string
}

export const ATTENTION: AttentionItem[] = [
  { tone: 'danger', text: '3 парсера упали', detail: 'TON, BTC-cash, Alfa' },
  { tone: 'warn', text: '2 жалобы без ответа', detail: '24ч+' },
  { tone: 'brand', text: '3 неизвестные валюты в очереди' },
]

export interface ModReview {
  id: string
  author: string
  exchanger: string
  direction: string
  amount: string
  ago: string
  rating: number
  body: string
  tags: { label: string; tone: HealthTone }[]
}

export const MOD_REVIEWS: ModReview[] = [
  {
    id: 'R-8842',
    author: 'Аноним',
    exchanger: 'CryptoBridge',
    direction: 'USDT → Тинькофф',
    amount: '≈ 250 000 ₽',
    ago: '18 мин назад',
    rating: 1,
    body: 'Курс на сайте отличался от заявленного на Kurso, в итоге получил меньше. Будьте внимательны!',
    tags: [
      { label: '⚑ возможный спор', tone: 'warn' },
      { label: 'негатив', tone: 'muted' },
    ],
  },
  {
    id: 'R-8843',
    author: 'Дмитрий В.',
    exchanger: 'NetEx24',
    direction: 'BTC → Сбербанк',
    amount: '≈ 480 000 ₽',
    ago: '32 мин назад',
    rating: 5,
    body: 'Быстро, без задержек. Оператор на связи, всё пришло в течение 8 минут. Рекомендую.',
    tags: [{ label: 'позитив', tone: 'success' }],
  },
  {
    id: 'R-8844',
    author: 'Аноним',
    exchanger: '24Paybank',
    direction: 'ETH → Наличные USD',
    amount: '≈ 1 200 $',
    ago: '54 мин назад',
    rating: 3,
    body: 'Курс нормальный, но пришлось ждать подтверждения почти полчаса. Могло быть и лучше.',
    tags: [{ label: 'нейтрально', tone: 'muted' }],
  },
  {
    id: 'R-8845',
    author: 'Ольга К.',
    exchanger: 'BaksMan',
    direction: 'USDT → Наличные RUB',
    amount: '≈ 90 000 ₽',
    ago: '1 ч назад',
    rating: 2,
    body: 'Заявленный резерв не совпал, предложили обменять частично. В итоге пришлось искать другой обменник.',
    tags: [
      { label: '⚑ спор о резерве', tone: 'warn' },
      { label: 'негатив', tone: 'muted' },
    ],
  },
]

export const MOD_TOTAL = 14

export interface MapCandidate {
  source: string
  suggestion?: { code: string; name: string; symbol: string; color: string; dark?: boolean }
  confidence?: number
}

export const MAP_QUEUE: MapCandidate[] = [
  {
    source: 'USDT.TRC20',
    suggestion: { code: 'USDT', name: 'USDT', symbol: '₮', color: '#26A17B' },
    confidence: 98,
  },
  {
    source: 'SBER RUB',
    suggestion: { code: 'SBER', name: 'Сбербанк', symbol: 'С', color: '#1FAE54' },
    confidence: 95,
  },
  { source: 'XMONERO' },
]

export interface Complaint {
  id: string
  client: string
  exchanger: string
  direction: string
  amount: string
  ago: string
  tag: string
  body: string
  reply?: string
}

export const COMPLAINTS: Complaint[] = [
  {
    id: 'C-1042',
    client: 'Игорь П.',
    exchanger: 'CryptoBridge',
    direction: 'BTC → Наличные',
    amount: '420 000 ₽',
    ago: '18ч',
    tag: 'спор о курсе',
    body: 'Курс на сайте отличался от итогового при подтверждении. Разница ≈ 1.2%, оператор сослался на «обновление курса».',
    reply: 'Готовы вернуть комиссию в качестве компенсации.',
  },
  {
    id: 'C-1043',
    client: 'Марина С.',
    exchanger: 'NetEx24',
    direction: 'USDT → Тинькофф',
    amount: '150 000 ₽',
    ago: '1д 4ч',
    tag: 'задержка выплаты',
    body: 'Обмен подтверждён, но средства не поступили в течение 2 часов. Поддержка отвечала с задержками.',
  },
  {
    id: 'C-1044',
    client: 'Аноним',
    exchanger: '24Paybank',
    direction: 'ETH → Наличные USD',
    amount: '3 400 $',
    ago: '2д',
    tag: 'некорректный резерв',
    body: 'Резерв на карточке не соответствовал доступному, обмен отменили после оплаты. Требую разбор.',
  },
]

export interface Payout {
  recipient: string
  kind: string
  method: string
  amount: number
}

export const PAYOUTS: Payout[] = [
  { recipient: '@cryptoblog', kind: 'партнёр · revshare', method: 'USDT TRC20', amount: 412_000 },
  { recipient: 'NetEx24', kind: 'обменник · бонус', method: 'Банк. перевод', amount: 388_000 },
  { recipient: '@exchange_ru', kind: 'партнёр · revshare', method: 'USDT TRC20', amount: 214_000 },
  { recipient: '@cryptotop', kind: 'партнёр · revshare', method: 'USDT TRC20', amount: 186_000 },
]

export const PAYOUTS_TOTAL = 1_200_000

export interface AuditEntry {
  time: string
  actor: string
  actorColor: string
  pre: string
  mark?: string
  markColor?: string
  post?: string
  tag: string
  tagTone: HealthTone
}

export const AUDIT: AuditEntry[] = [
  {
    time: '12:04',
    actor: 'AK',
    actorColor: '#2E7DF2',
    pre: 'Опубликован отзыв ',
    mark: '#R-8841',
    markColor: '#6BA6FF',
    post: ' о CryptoBridge',
    tag: 'опубликовано',
    tagTone: 'success',
  },
  {
    time: '11:58',
    actor: 'MS',
    actorColor: '#1F8A5B',
    pre: 'Выплата ',
    mark: '₽388k',
    markColor: '#6BA6FF',
    post: ' → NetEx24 подтверждена',
    tag: 'выплачено',
    tagTone: 'success',
  },
  {
    time: '11:42',
    actor: 'AK',
    actorColor: '#8A5A2B',
    pre: 'Валюта ',
    mark: 'USDT.TRC20',
    markColor: '#6BA6FF',
    post: ' привязана к USDT',
    tag: 'маппинг',
    tagTone: 'brand',
  },
  {
    time: '11:30',
    actor: 'SYS',
    actorColor: '#3A4452',
    pre: 'Парсер ',
    mark: 'feed-ton',
    markColor: '#EC7B7A',
    post: ' помечен как упавший',
    tag: 'авто',
    tagTone: 'danger',
  },
  {
    time: '11:15',
    actor: 'AK',
    actorColor: '#3A4452',
    pre: 'Обменник ',
    mark: 'ExchPro',
    markColor: '#6BA6FF',
    post: ' скрыт из выдачи',
    tag: 'скрыт',
    tagTone: 'muted',
  },
]

export interface ExportSet {
  key: 'clicks' | 'revenue' | 'reviews'
  title: string
  meta: string
}

export const EXPORTS: ExportSet[] = [
  { key: 'clicks', title: 'Клики и переходы', meta: '84.2k строк · сегодня' },
  { key: 'revenue', title: 'Доход по обменникам', meta: '128 строк · месяц' },
  { key: 'reviews', title: 'Отзывы и жалобы', meta: '2.4k строк · всё время' },
]

export interface ParserLogLine {
  time: string
  level: 'OK' | 'WARN' | 'ERR'
  text: string
}

export const PARSER_LOGS: ParserLogLine[] = [
  { time: '12:04:18', level: 'OK', text: 'fetch 200 · 24 rates · 240ms' },
  { time: '12:04:08', level: 'OK', text: 'fetch 200 · 24 rates · 218ms' },
  { time: '12:03:58', level: 'OK', text: 'fetch 200 · 24 rates · 251ms' },
  { time: '12:03:48', level: 'WARN', text: 'reserve BTC stale (9m), kept last value' },
  { time: '12:03:38', level: 'OK', text: 'fetch 200 · 24 rates · 233ms' },
  { time: '12:03:28', level: 'ERR', text: 'direction TON-TINKOFF: unknown currency «XTON»' },
  { time: '12:03:18', level: 'OK', text: 'fetch 200 · 24 rates · 229ms' },
]

/** Result shown after a "Тестовый прогон" in the editor. */
export const TEST_RUN = {
  ok: true,
  ms: 240,
  http: 200,
  directions: 24,
  recognized: 24,
  errors: 0,
}

/** Dashboard headline metrics (exchangers count is overridden with live data). */
export const METRICS = {
  parsersHealthy: 116,
  parsersTotal: 128,
  clicksPerDay: 84_200,
  revenuePerMonth: 2_400_000,
  moderation: 14,
  complaints: 3,
}
