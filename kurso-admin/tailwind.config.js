/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        // Surfaces, darkest → lightest (the admin is dark-only by design).
        bg: '#0A0B0D', // page backdrop
        panel: '#0C0D10', // app frame behind cards
        well: '#0E1013', // inset fields / inner rows
        surface: '#15171B', // card background
        chrome: '#15181C', // window chrome / sidebar dividers
        // Borders.
        line: '#20242A', // default card border
        'line-soft': '#1B1E22', // table row divider
        'line-strong': '#2E343B', // chip / control border
        brand: {
          DEFAULT: '#2E7DF2',
          bright: '#4A90F5', // gradient top-stop
          light: '#6BA6FF', // brand text on dark
        },
        ink: {
          DEFAULT: '#F4F5F7', // primary text
          body: '#C9CED4', // body copy
          muted: '#A8AEB6', // secondary
          faint: '#6E757E', // labels, meta
          fainter: '#4A5159', // section eyebrows, timestamps
        },
        success: { DEFAULT: '#2BC58C', text: '#3DD79C' },
        warn: { DEFAULT: '#D99A33', text: '#E0A954' },
        danger: { DEFAULT: '#EC5B5B', text: '#EC7B7A' },
      },
      fontFamily: {
        sans: ['Onest', 'system-ui', 'sans-serif'],
        mono: ['Geist Mono', 'ui-monospace', 'monospace'],
        label: ['JetBrains Mono', 'ui-monospace', 'monospace'],
      },
      keyframes: {
        kpulse: { '0%,100%': { opacity: '1' }, '50%': { opacity: '.35' } },
      },
      animation: {
        kpulse: 'kpulse 1.8s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
