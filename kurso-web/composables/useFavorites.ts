// Favorited exchangers. Real client state: the heart on exchanger cards / the
// detail page toggles membership here, it persists (localStorage) and the
// account favorites screen resolves each slug against the live exchangers list
// for the current online status. Starts empty.

export interface Favorite {
  slug: string
  name: string
  initials: string
  color: string
}

export function useFavorites() {
  const favorites = useState<Favorite[]>('favorites', () => [])

  const isFavorite = (slug: string) => favorites.value.some((f) => f.slug === slug)

  function remove(slug: string) {
    favorites.value = favorites.value.filter((f) => f.slug !== slug)
  }
  function add(f: Favorite) {
    if (!isFavorite(f.slug)) favorites.value = [...favorites.value, f]
  }
  function toggle(f: Favorite) {
    if (isFavorite(f.slug)) remove(f.slug)
    else add(f)
  }

  return { favorites, isFavorite, add, remove, toggle }
}
