// Affiliate capture: when a visitor lands via a partner link (`?ref=code` or
// `?ref=code.tag`), persist the tag in the `kurso_ref` cookie for 90 days. The
// API's clickout handler reads this cookie to attribute outbound clicks, and the
// register endpoint reads it to credit the referrer. Not httpOnly — it is an
// affiliate tag, not a secret.
export default defineNuxtPlugin(() => {
  const route = useRoute()
  const ref = route.query.ref
  const code = Array.isArray(ref) ? ref[0] : ref
  if (typeof code !== 'string' || !code.trim()) return

  const cookie = useCookie('kurso_ref', {
    maxAge: 90 * 24 * 60 * 60,
    path: '/',
    sameSite: 'lax',
  })
  cookie.value = code.trim().slice(0, 64)
})
