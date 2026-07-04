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
    component: () => import('../components/layout/MerchantLayout.vue'),
    children: [
      { path: '', name: 'dashboard', component: () => import('../pages/DashboardView.vue') },
      { path: 'rates', name: 'rates', component: () => import('../pages/RatesView.vue') },
      { path: 'reviews', name: 'reviews', component: () => import('../pages/ReviewsView.vue') },
      { path: 'traffic', name: 'traffic', component: () => import('../pages/TrafficView.vue') },
      { path: 'profile', name: 'profile', component: () => import('../pages/ProfileView.vue') },
      {
        path: 'complaints',
        name: 'complaints',
        component: () => import('../pages/ComplaintsView.vue'),
      },
      { path: 'billing', name: 'billing', component: () => import('../pages/BillingView.vue') },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

// Guard: everything but the login page requires a merchant session.
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
