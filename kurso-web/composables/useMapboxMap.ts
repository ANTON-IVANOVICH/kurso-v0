import { watch, onScopeDispose, type Ref } from 'vue'
import type { MapPoint } from './useMapPoints'

// Real Mapbox GL integration for the map page (client-only). Encapsulates the
// map lifecycle: a dark base map, a price-pill marker per exchanger, a "you are
// here" marker, and a popup for the selected point wired to route/clickout
// callbacks. mapbox-gl is imported dynamically so it never touches the SSR
// bundle. Only used when a Mapbox token is configured; otherwise the page shows
// a static projected preview instead.

interface Options {
  container: Ref<HTMLElement | null>
  token: string
  points: Ref<MapPoint[]>
  selected: Ref<string | null>
  userLoc: Ref<{ lat: number; lng: number } | null>
  fmtRate: (rate: string | null | undefined) => string
  labels: { partner: string; route: string; go: string }
  onSelect: (slug: string) => void
  onRoute: (p: MapPoint) => void
  onGo: (p: MapPoint) => void
}

const MOSCOW: [number, number] = [37.6173, 55.7558]

export function useMapboxMap(opts: Options) {
  if (!import.meta.client)
    return { ready: () => false, zoomIn() {}, zoomOut() {}, flyTo() {}, locate() {} }

  let map: import('mapbox-gl').Map | null = null
  let mapboxgl: typeof import('mapbox-gl').default | null = null
  const markers = new Map<string, import('mapbox-gl').Marker>()
  const els = new Map<string, HTMLElement>()
  let userMarker: import('mapbox-gl').Marker | null = null
  let popup: import('mapbox-gl').Popup | null = null

  function pinEl(p: MapPoint, active: boolean): HTMLElement {
    const el = document.createElement('button')
    el.type = 'button'
    el.className =
      'kmap-pin' + (p.partner ? ' kmap-pin--partner' : '') + (active ? ' kmap-pin--active' : '')
    const initials = p.name
      .replace(/[^A-Za-zА-Яа-я0-9]/g, '')
      .slice(0, 2)
      .toUpperCase()
    el.innerHTML =
      `<span class="kmap-pin__badge">${initials}</span>` +
      `<span class="kmap-pin__rate">${opts.fmtRate(p.rate).replace(' ₽', '')}</span>`
    el.addEventListener('click', (e) => {
      e.stopPropagation()
      opts.onSelect(p.slug)
    })
    return el
  }

  function popupHtml(p: MapPoint): HTMLElement {
    const wrap = document.createElement('div')
    wrap.className = 'kmap-popup'
    const rating = p.ratingAvg != null ? p.ratingAvg.toFixed(1) : '—'
    wrap.innerHTML = `
      <div class="kmap-popup__head">
        <div class="kmap-popup__name">${p.name}${p.partner ? `<span class="kmap-popup__tag">${opts.labels.partner}</span>` : ''}</div>
        <div class="kmap-popup__rate">${opts.fmtRate(p.rate)}</div>
      </div>
      <div class="kmap-popup__meta">${p.address ?? ''}${p.hours ? ' · ' + p.hours : ''}</div>
      <div class="kmap-popup__meta">★ ${rating} · ${p.reviewsCount}</div>
      <div class="kmap-popup__actions">
        <button class="kmap-btn kmap-btn--ghost" data-act="route">${opts.labels.route}</button>
        <button class="kmap-btn kmap-btn--primary" data-act="go">${opts.labels.go}</button>
      </div>`
    wrap.querySelector('[data-act="route"]')?.addEventListener('click', () => opts.onRoute(p))
    wrap.querySelector('[data-act="go"]')?.addEventListener('click', () => opts.onGo(p))
    return wrap
  }

  async function init() {
    if (map || !opts.container.value) return
    const mod = await import('mapbox-gl')
    mapboxgl = mod.default
    mapboxgl.accessToken = opts.token
    map = new mapboxgl.Map({
      container: opts.container.value,
      style: 'mapbox://styles/mapbox/dark-v11',
      center: MOSCOW,
      zoom: 11.5,
      attributionControl: false,
    })
    map.on('load', () => {
      syncMarkers()
      syncUser()
      fitPoints()
    })
  }

  function syncMarkers() {
    if (!map || !mapboxgl) return
    const seen = new Set<string>()
    for (const p of opts.points.value) {
      seen.add(p.slug)
      const active = p.slug === opts.selected.value
      if (markers.has(p.slug)) {
        const el = els.get(p.slug)!
        el.className =
          'kmap-pin' + (p.partner ? ' kmap-pin--partner' : '') + (active ? ' kmap-pin--active' : '')
      } else {
        const el = pinEl(p, active)
        els.set(p.slug, el)
        markers.set(
          p.slug,
          new mapboxgl.Marker({ element: el, anchor: 'bottom' })
            .setLngLat([p.lng, p.lat])
            .addTo(map),
        )
      }
    }
    for (const [slug, m] of markers)
      if (!seen.has(slug)) {
        m.remove()
        markers.delete(slug)
        els.delete(slug)
      }
  }

  function syncUser() {
    if (!map || !mapboxgl) return
    const u = opts.userLoc.value
    if (!u) return
    if (!userMarker) {
      const el = document.createElement('div')
      el.className = 'kmap-user'
      userMarker = new mapboxgl.Marker({ element: el }).setLngLat([u.lng, u.lat]).addTo(map)
    } else {
      userMarker.setLngLat([u.lng, u.lat])
    }
  }

  function fitPoints() {
    if (!map || !mapboxgl || !opts.points.value.length) return
    const b = new mapboxgl.LngLatBounds()
    for (const p of opts.points.value) b.extend([p.lng, p.lat])
    if (opts.userLoc.value) b.extend([opts.userLoc.value.lng, opts.userLoc.value.lat])
    map.fitBounds(b, {
      padding: { top: 120, bottom: 220, left: 60, right: 60 },
      maxZoom: 14,
      duration: 0,
    })
  }

  function openPopup(slug: string) {
    if (!map || !mapboxgl) return
    const p = opts.points.value.find((x) => x.slug === slug)
    if (!p) return
    popup?.remove()
    popup = new mapboxgl.Popup({ offset: 30, closeButton: false, className: 'kmap-popup-wrap' })
      .setLngLat([p.lng, p.lat])
      .setDOMContent(popupHtml(p))
      .addTo(map)
    map.flyTo({ center: [p.lng, p.lat], zoom: Math.max(map.getZoom(), 13), speed: 0.8 })
  }

  // React to data + selection changes.
  watch(
    opts.points,
    () => {
      syncMarkers()
      if (!markers.size) return
      if (map?.loaded()) fitPoints()
    },
    { deep: false },
  )
  watch(opts.selected, (slug) => {
    syncMarkers()
    if (slug) openPopup(slug)
  })
  watch(opts.userLoc, syncUser)

  if (opts.container.value) init()
  else watch(opts.container, (c) => c && init())

  onScopeDispose(() => {
    popup?.remove()
    map?.remove()
    map = null
  })

  return {
    ready: () => !!map,
    zoomIn: () => map?.zoomIn(),
    zoomOut: () => map?.zoomOut(),
    flyTo: (lat: number, lng: number) => map?.flyTo({ center: [lng, lat], zoom: 14 }),
    locate: () => {
      const u = opts.userLoc.value
      if (u) map?.flyTo({ center: [u.lng, u.lat], zoom: 14 })
    },
  }
}
