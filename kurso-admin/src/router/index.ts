import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../pages/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('../components/layout/AdminLayout.vue'),
    children: [
      { path: '', name: 'dashboard', component: () => import('../pages/DashboardView.vue') },
      {
        path: 'exchangers',
        name: 'exchangers',
        component: () => import('../pages/ExchangersView.vue'),
      },
      {
        path: 'exchangers/new',
        name: 'exchanger-new',
        component: () => import('../pages/ExchangerEditorView.vue'),
      },
      {
        path: 'exchangers/:slug',
        name: 'exchanger-edit',
        component: () => import('../pages/ExchangerEditorView.vue'),
      },
      {
        path: 'currencies',
        name: 'currencies',
        component: () => import('../pages/CurrenciesView.vue'),
      },
      {
        path: 'directions',
        name: 'directions',
        component: () => import('../pages/DirectionsView.vue'),
      },
      { path: 'parsers', name: 'parsers', component: () => import('../pages/ParsersView.vue') },
      { path: 'reviews', name: 'reviews', component: () => import('../pages/ReviewsView.vue') },
      {
        path: 'complaints',
        name: 'complaints',
        component: () => import('../pages/ComplaintsView.vue'),
      },
      { path: 'partners', name: 'partners', component: () => import('../pages/PartnersView.vue') },
      { path: 'payouts', name: 'payouts', component: () => import('../pages/PayoutsView.vue') },
      { path: 'export', name: 'export', component: () => import('../pages/ExportView.vue') },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

// Guard: everything but the login page requires a session.
router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isAuthenticated) {
    return { name: 'login', query: to.fullPath !== '/' ? { redirect: to.fullPath } : undefined }
  }
  if (to.name === 'login' && auth.isAuthenticated) {
    return { name: 'dashboard' }
  }
  return true
})
