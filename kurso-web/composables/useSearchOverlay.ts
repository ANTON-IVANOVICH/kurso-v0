// Global controller for the site search overlay (opened from the header lupe).
export function useSearchOverlay() {
  const open = useState('searchOverlayOpen', () => false)
  return {
    open,
    openSearch: () => (open.value = true),
    close: () => (open.value = false),
  }
}
