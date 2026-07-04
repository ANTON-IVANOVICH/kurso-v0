import type { Config } from 'tailwindcss'

// Kurso design tokens — dark theme.
// Fonts: Onest (sans), Geist Mono (mono / tabular figures), JetBrains Mono (labels).
export default <Partial<Config>>{
  content: [
    './app.vue',
    './error.vue',
    './components/**/*.{vue,js,ts}',
    './pages/**/*.{vue,js,ts}',
    './layouts/**/*.{vue,js,ts}',
  ],
  theme: {
    extend: {
      colors: {
        canvas: '#0A0B0D', // page background
        panel: '#0C0D10', // large panels / app body
        well: '#0E1013', // inputs / deepest wells
        surface: {
          DEFAULT: '#15171B', // cards
          hi: '#15181C', // dropdowns, toasts, modals
          raised: '#1F2329', // secondary buttons, pills
          nested: '#121417', // nested cards
          nav: '#1A1D22', // mobile bottom nav
          chip: '#262B32', // swap button etc.
        },
        line: {
          DEFAULT: '#20242A', // default border
          strong: '#2E343B', // stronger border
          subtle: '#1B1E22', // subtle divider
          faint: '#16191D',
        },
        ink: {
          DEFAULT: '#F4F5F7', // primary text
          bright: '#C9CED4',
          muted: '#A8AEB6', // secondary text
          dim: '#9BA1A8',
          faint: '#6E757E', // muted labels
          ghost: '#4A5159',
          fainter: '#3A4047',
        },
        brand: {
          DEFAULT: '#2E7DF2', // primary blue
          hover: '#3B86F5',
          light: '#4A90F5', // gradient light stop
          bright: '#6BA6FF', // accent text on dark
        },
        success: { DEFAULT: '#2BC58C', bright: '#3DD79C' },
        warning: { DEFAULT: '#E0A954', deep: '#D99A33', light: '#F2D08A' },
        danger: { DEFAULT: '#EC5B5B', soft: '#E08C7A' },
        violet: { DEFAULT: '#8B5CF6', light: '#C0A4F5' },
      },
      fontFamily: {
        sans: ['Onest', 'system-ui', 'sans-serif'],
        mono: ['Geist Mono', 'ui-monospace', 'monospace'],
        label: ['JetBrains Mono', 'ui-monospace', 'monospace'],
      },
      borderRadius: {
        // matched to the design scale
        sm: '7px',
        DEFAULT: '9px',
        md: '10px',
        lg: '11px',
        xl: '12px',
        '2xl': '16px',
        '3xl': '20px',
        full: '9999px',
      },
      boxShadow: {
        glow: '0 8px 24px rgba(46,125,242,0.35)',
        card: '0 24px 60px rgba(0,0,0,0.45)',
        panel: '0 30px 80px rgba(0,0,0,0.5)',
        pop: '0 18px 40px rgba(0,0,0,0.55)',
        modal: '0 30px 70px rgba(0,0,0,0.6)',
        toast: '0 14px 34px rgba(0,0,0,0.5)',
      },
      keyframes: {
        kpulse: {
          '0%,100%': { opacity: '1', transform: 'scale(1)' },
          '50%': { opacity: '.35', transform: 'scale(.82)' },
        },
        kspin: { to: { transform: 'rotate(360deg)' } },
        kping: {
          '0%': { transform: 'scale(1)', opacity: '.6' },
          '100%': { transform: 'scale(2.4)', opacity: '0' },
        },
      },
      animation: {
        kpulse: 'kpulse 1.6s ease-in-out infinite',
        kspin: 'kspin 0.9s linear infinite',
        kping: 'kping 2.2s ease-out infinite',
      },
    },
  },
  plugins: [],
}
