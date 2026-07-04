// Guards the account/alert section. With the auth scaffold a demo user is
// present, so this passes; once real auth lands, an unauthenticated visitor is
// bounced to /login with a return path.
export default defineNuxtRouteMiddleware((to) => {
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated.value) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }
})
