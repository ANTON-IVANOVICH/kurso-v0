import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Tooltip from 'primevue/tooltip'
import 'primeicons/primeicons.css'
import './style.css'
import App from './App.vue'
import { router } from './router'
import { KursoPreset } from './theme'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)

// Re-hydrate any persisted session before the router guard runs, so a refresh
// on a protected page doesn't bounce to /login.
useAuthStore(pinia).restore()

app.use(PrimeVue, {
  theme: {
    preset: KursoPreset,
    options: {
      // The admin is dark-only; <html class="dark"> keeps this always matched.
      darkModeSelector: '.dark',
      // PrimeVue sits between Tailwind's preflight and its utilities (see
      // style.css @layer order) so preflight can't strip component padding but
      // utility classes still win for deliberate overrides.
      cssLayer: { name: 'primevue', order: 'tailwind-base, primevue, tailwind-utilities' },
    },
  },
})
app.use(ToastService)
app.use(ConfirmationService)
app.directive('tooltip', Tooltip)

app.use(router)
app.mount('#app')
