import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'

// Kurso admin theme: Aura recoloured to the console mockup — brand-blue primary
// and the dark neutral ramp from the design (page #0A0B0D, cards #15171B, inset
// fields #0E1013, borders #20242A). Dark-only; `darkModeSelector: '.dark'` is
// always satisfied (see index.html <html class="dark">).
export const KursoPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#eef5ff',
      100: '#d9e8ff',
      200: '#b6d2ff',
      300: '#84b3ff',
      400: '#4A90F5',
      500: '#2E7DF2',
      600: '#1f66d8',
      700: '#1b52ad',
      800: '#1c478a',
      900: '#1c3d70',
      950: '#132546',
    },
    colorScheme: {
      dark: {
        surface: {
          0: '#ffffff',
          50: '#F4F5F7',
          100: '#C9CED4',
          200: '#A8AEB6',
          300: '#6E757E',
          400: '#4A5159',
          500: '#2E343B',
          600: '#20242A',
          700: '#1B1E22',
          800: '#15171B',
          900: '#0E1013',
          950: '#0A0B0D',
        },
        // Cards / panels / menus sit on #15171B, one step lighter than the page.
        content: {
          background: '{surface.800}',
          hoverBackground: '{surface.700}',
          borderColor: '{surface.600}',
        },
        overlay: {
          select: { background: '{surface.800}', borderColor: '{surface.600}' },
          popover: { background: '{surface.800}', borderColor: '{surface.600}' },
          modal: { background: '{surface.800}', borderColor: '{surface.600}' },
        },
        // Inputs use the darkest inset (#0E1013) like the mockup fields.
        formField: {
          background: '{surface.900}',
          borderColor: '{surface.600}',
          hoverBorderColor: '{surface.500}',
        },
      },
    },
  },
})
