<script setup lang="ts">
import { computed, ref } from 'vue'
import { useCatalog } from '../composables/useCatalog'
import PageHeader from '../components/ui/PageHeader.vue'
import UiIcon from '../components/ui/UiIcon.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import ToggleButton from 'primevue/togglebutton'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'

const catalog = useCatalog()
const search = ref('')
const popularOnly = ref(false)

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  return catalog.directions.filter((d) => {
    if (popularOnly.value && !d.isPopular) return false
    if (
      q &&
      !d.slug.toLowerCase().includes(q) &&
      !`${d.fromCode} ${d.toCode} ${d.fromName} ${d.toName}`.toLowerCase().includes(q)
    )
      return false
    return true
  })
})
</script>

<template>
  <div>
    <PageHeader title="Направления" subtitle="Пары обмена, по которым сравниваются курсы" />

    <DataTable
      :value="filtered"
      :loading="catalog.loading"
      data-key="id"
      paginator
      :rows="15"
      class="overflow-hidden rounded-[14px] border border-line"
    >
      <template #header>
        <div class="flex flex-wrap items-center gap-2.5">
          <IconField>
            <InputIcon class="pi pi-search" />
            <InputText v-model="search" placeholder="Поиск направления…" class="w-60" />
          </IconField>
          <ToggleButton
            v-model="popularOnly"
            on-label="Только популярные"
            off-label="Все"
            on-icon="pi pi-star-fill"
            off-icon="pi pi-star"
          />
          <span class="ml-auto text-xs text-ink-faint">
            Показано <span class="tnum">{{ filtered.length }}</span> из
            <span class="tnum">{{ catalog.directions.length }}</span>
          </span>
        </div>
      </template>
      <template #empty>
        <div class="py-8 text-center text-sm text-ink-faint">Направления не найдены</div>
      </template>

      <Column header="Пара">
        <template #body="{ data }">
          <div class="flex items-center gap-2 font-semibold">
            <span>{{ data.fromCode }}</span>
            <UiIcon name="arrowRight" :size="14" class="text-ink-faint" />
            <span>{{ data.toCode }}</span>
          </div>
        </template>
      </Column>
      <Column header="Из">
        <template #body="{ data }"
          ><span class="text-ink-muted">{{ data.fromName }}</span></template
        >
      </Column>
      <Column header="В">
        <template #body="{ data }"
          ><span class="text-ink-muted">{{ data.toName }}</span></template
        >
      </Column>
      <Column field="slug" header="Slug">
        <template #body="{ data }"
          ><span class="tnum text-xs text-ink-faint">{{ data.slug }}</span></template
        >
      </Column>
      <Column header="Популярное">
        <template #body="{ data }">
          <Tag v-if="data.isPopular" value="популярное" severity="info" />
          <span v-else class="text-ink-faint">—</span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>
