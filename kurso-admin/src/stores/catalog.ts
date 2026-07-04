import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api, ApiError } from '../lib/api'
import type { Currency, Direction, Exchanger } from '../types/models'

// Catalogue data for the admin. Reads come from the live public API
// (`/api/v1/exchangers`, `/currencies`, `/directions`). Mutations (status
// changes, profile edits, new exchangers) have no admin endpoint yet, so they
// apply to an in-memory overlay — the UX is real and reverts on reload. Each
// mutation is one `api.patch/post` away from being persisted.

export const useCatalogStore = defineStore('catalog', () => {
  const exchangers = ref<Exchanger[]>([])
  const currencies = ref<Currency[]>([])
  const directions = ref<Direction[]>([])

  const loading = ref(false)
  const error = ref<string | null>(null)
  const loaded = ref(false)

  const partnersCount = computed(() => exchangers.value.filter((e) => e.partner).length)
  const activeCount = computed(() => exchangers.value.filter((e) => e.status === 'active').length)

  function findExchanger(slug: string): Exchanger | undefined {
    return exchangers.value.find((e) => e.slug === slug)
  }

  async function load(force = false) {
    if (loaded.value && !force) return
    loading.value = true
    error.value = null
    try {
      const [ex, cu, di] = await Promise.all([
        api.get<Exchanger[]>('/api/v1/exchangers'),
        api.get<Currency[]>('/api/v1/currencies'),
        api.get<Direction[]>('/api/v1/directions'),
      ])
      exchangers.value = ex
      currencies.value = cu
      directions.value = di
      loaded.value = true
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : 'Не удалось загрузить каталог'
    } finally {
      loading.value = false
    }
  }

  /** Local-only patch to an exchanger (pending the admin write API). */
  function patchExchanger(slug: string, patch: Partial<Exchanger>) {
    const i = exchangers.value.findIndex((e) => e.slug === slug)
    if (i !== -1) exchangers.value[i] = { ...exchangers.value[i], ...patch }
  }

  return {
    exchangers,
    currencies,
    directions,
    loading,
    error,
    loaded,
    partnersCount,
    activeCount,
    findExchanger,
    load,
    patchExchanger,
  }
})
