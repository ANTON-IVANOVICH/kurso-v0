import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { PiniaColada } from '@pinia/colada'
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
app.use(PiniaColada) // async cache for admin data queries

app.use(PrimeVue, {
  theme: {
    preset: KursoPreset,
    options: {
      darkModeSelector: '.dark',
      cssLayer: { name: 'primevue', order: 'tailwind-base, primevue, tailwind-utilities' },
    },
  },
})
app.use(ToastService)
app.use(ConfirmationService)
app.directive('tooltip', Tooltip)

// Revive the session from the httpOnly refresh cookie, THEN wire the router so
// the first navigation guard already sees the correct auth state.
useAuthStore(pinia)
  .restore()
  .finally(() => {
    app.use(router)
    app.mount('#app')
  })
