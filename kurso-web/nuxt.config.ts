// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss', '@nuxt/eslint', '@pinia/nuxt'],
  css: ['~/assets/css/main.css'],
  devServer: { port: 3000 },
  runtimeConfig: {
    public: {
      // Base URL of kurso-api. Override in prod via NUXT_PUBLIC_API_BASE.
      apiBase: 'http://localhost:8080',
      // Mapbox GL access token for the map page. Override via NUXT_PUBLIC_MAPBOX_TOKEN.
      // Empty = the map falls back to a static projected preview (no tiles).
      mapboxToken: '',
    },
  },
  // Auto-import components by filename only (no directory prefix), so
  // components/ui/KButton.vue is <KButton>, not <UiKButton>.
  components: [{ path: '~/components', pathPrefix: false }],
  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      title: 'Kurso — агрегатор курсов обменников',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'theme-color', content: '#0A0B0D' },
        { name: 'apple-mobile-web-app-title', content: 'Kurso' },
        { name: 'apple-mobile-web-app-capable', content: 'yes' },
        { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' },
        // Social sharing (Open Graph + Twitter) using the branded cover.
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'Kurso' },
        { property: 'og:title', content: 'Kurso — агрегатор курсов обменников' },
        {
          property: 'og:description',
          content:
            'Лучший курс обмена криптовалюты за секунды: сравнение курсов, алерты, карта пунктов.',
        },
        { property: 'og:image', content: 'https://kurso.io/og-cover.png' },
        { property: 'og:image:width', content: '1200' },
        { property: 'og:image:height', content: '630' },
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:title', content: 'Kurso — агрегатор курсов обменников' },
        {
          name: 'twitter:description',
          content:
            'Лучший курс обмена криптовалюты за секунды: сравнение курсов, алерты, карта пунктов.',
        },
        { name: 'twitter:image', content: 'https://kurso.io/og-cover.png' },
      ],
      link: [
        // Kurso identity favicons + PWA manifest.
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
        { rel: 'icon', type: 'image/png', sizes: '96x96', href: '/favicon-96.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
        { rel: 'manifest', href: '/site.webmanifest' },
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Geist+Mono:wght@400;500;600;700&family=Onest:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600;700&display=swap',
        },
      ],
    },
  },
})
