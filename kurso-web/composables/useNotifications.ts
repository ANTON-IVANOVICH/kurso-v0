// Global controller for the notifications feed (opened from the header bell).
export function useNotifications() {
  const open = useState('notificationsOpen', () => false)
  return {
    open,
    openNotifications: () => (open.value = true),
    close: () => (open.value = false),
  }
}
