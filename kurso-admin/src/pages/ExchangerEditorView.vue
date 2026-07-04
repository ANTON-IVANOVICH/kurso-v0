<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useCatalog } from '../composables/useCatalog'
import { api } from '../lib/api'
import type { Exchanger, ExchangerStatus } from '../types/models'
import { PARSER_LOGS, TEST_RUN } from '../lib/fixtures'
import ExchangerAvatar from '../components/ui/ExchangerAvatar.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import InputNumber from 'primevue/inputnumber'
import ToggleSwitch from 'primevue/toggleswitch'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import Toast from 'primevue/toast'

const route = useRoute()
const router = useRouter()
const catalog = useCatalog()
const toast = useToast()

const isNew = computed(() => route.name === 'exchanger-new')
const slug = computed(() => (route.params.slug as string) ?? '')

const loading = ref(false)
const notFound = ref(false)
const source = ref<Exchanger | null>(null)

// Editable form. Fields the API exposes (name/website/description/status) persist
// via the store overlay; type/conditions/parser have no backend yet and are kept
// as local drafts so the editor is fully interactive.
const form = ref({
  name: '',
  website: '',
  type: 'Онлайн',
  description: '',
  status: 'active' as ExchangerStatus,
  min: 3000 as number | null,
  max: 5000000 as number | null,
  exchangeTime: '5–15 мин',
  kyc: false,
  feedUrl: '',
  format: 'XML (ExchML)',
  interval: '10 сек',
})

const typeOptions = ['Онлайн', 'Наличные', 'Смешанный']
const formatOptions = ['XML (ExchML)', 'JSON', 'HTML']
const intervalOptions = ['10 сек', '15 сек', '30 сек', '60 сек']
const statusOptions = [
  { label: 'Активен', value: 'active' },
  { label: 'На паузе', value: 'paused' },
  { label: 'Забанен', value: 'banned' },
]
const STATUS_META: Record<ExchangerStatus, { label: string; severity: string }> = {
  active: { label: 'активен', severity: 'success' },
  paused: { label: 'на паузе', severity: 'warn' },
  banned: { label: 'забанен', severity: 'danger' },
}

function hydrate(e: Exchanger) {
  source.value = e
  form.value.name = e.name
  form.value.website = e.websiteUrl ?? ''
  form.value.description = e.description ?? ''
  form.value.status = e.status
  form.value.feedUrl = e.websiteUrl
    ? `https://${e.websiteUrl.replace(/^https?:\/\//, '')}/api/rates.xml`
    : ''
}

onMounted(async () => {
  if (isNew.value) return
  const local = catalog.findExchanger(slug.value)
  if (local) {
    hydrate(local)
    return
  }
  loading.value = true
  try {
    const e = await api.get<Exchanger>(`/api/v1/exchangers/${slug.value}`)
    hydrate(e)
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
})

function save() {
  if (!form.value.name.trim()) {
    toast.add({ severity: 'warn', summary: 'Укажите название', life: 3000 })
    return
  }
  if (isNew.value) {
    toast.add({
      severity: 'info',
      summary: 'Черновик сохранён локально',
      detail: 'Создание обменников появится с админским API',
      life: 4000,
    })
    return
  }
  catalog.patchExchanger(slug.value, {
    name: form.value.name.trim(),
    websiteUrl: form.value.website || null,
    description: form.value.description || null,
    status: form.value.status,
  })
  toast.add({ severity: 'success', summary: 'Сохранено', life: 2500 })
}

// Test run.
const testing = ref(false)
const testResult = ref<typeof TEST_RUN | null>(null)
function runTest() {
  testing.value = true
  testResult.value = null
  window.setTimeout(() => {
    testResult.value = TEST_RUN
    testing.value = false
  }, 700)
}

const logFilter = ref<'all' | 'WARN' | 'ERR'>('all')
const logs = computed(() =>
  logFilter.value === 'all' ? PARSER_LOGS : PARSER_LOGS.filter((l) => l.level === logFilter.value),
)
const LOG_COLOR = { OK: 'text-success-text', WARN: 'text-warn-text', ERR: 'text-danger-text' }
</script>

<template>
  <div>
    <Toast />

    <!-- breadcrumb -->
    <nav class="mb-4 flex items-center gap-2 text-xs text-ink-faint">
      <RouterLink to="/exchangers" class="hover:text-ink-muted">Каталог</RouterLink>
      <span class="text-line-strong">/</span>
      <RouterLink to="/exchangers" class="hover:text-ink-muted">Обменники</RouterLink>
      <span class="text-line-strong">/</span>
      <span class="text-ink-muted">{{ isNew ? 'Новый обменник' : (source?.name ?? slug) }}</span>
    </nav>

    <div v-if="notFound" class="rounded-[14px] border border-line bg-surface p-10 text-center">
      <p class="text-ink-muted">
        Обменник <span class="tnum">{{ slug }}</span> не найден.
      </p>
      <Button label="К списку" class="mt-4" @click="router.push('/exchangers')" />
    </div>

    <template v-else>
      <!-- header -->
      <div class="mb-5 flex flex-wrap items-center gap-4">
        <ExchangerAvatar
          :slug="source?.slug ?? 'new'"
          :name="form.name || 'Новый'"
          :size="48"
          shape="rounded"
        />
        <div class="flex-1">
          <div class="flex flex-wrap items-center gap-2.5">
            <span class="text-xl font-extrabold tracking-tight">{{
              form.name || 'Новый обменник'
            }}</span>
            <Tag v-if="source?.partner" value="Партнёр" severity="info" />
            <Tag
              :value="STATUS_META[form.status].label"
              :severity="STATUS_META[form.status].severity"
            />
          </div>
          <div v-if="source" class="tnum mt-1 text-xs text-ink-faint">
            {{ source.slug }} · на площадке с {{ source.onSince }}
          </div>
        </div>
        <Button label="Отмена" severity="secondary" @click="router.push('/exchangers')" />
        <Button label="Сохранить" icon="pi pi-check" @click="save" />
      </div>

      <div
        v-if="loading"
        class="rounded-[14px] border border-line bg-surface p-10 text-center text-ink-faint"
      >
        Загрузка…
      </div>

      <Tabs v-else value="profile">
        <TabList>
          <Tab value="profile">Профиль</Tab>
          <Tab value="conditions">Условия</Tab>
          <Tab value="parser">Парсер</Tab>
          <Tab value="logs">Логи</Tab>
        </TabList>
        <TabPanels>
          <!-- PROFILE -->
          <TabPanel value="profile">
            <div class="grid gap-4 md:grid-cols-2">
              <label class="block">
                <span class="mb-1.5 block text-xs text-ink-faint">Название</span>
                <InputText v-model="form.name" fluid />
              </label>
              <label class="block">
                <span class="mb-1.5 block text-xs text-ink-faint">Сайт</span>
                <InputText v-model="form.website" placeholder="cryptobridge.io" fluid />
              </label>
              <label class="block">
                <span class="mb-1.5 block text-xs text-ink-faint">Тип</span>
                <Select v-model="form.type" :options="typeOptions" fluid />
              </label>
              <label class="block">
                <span class="mb-1.5 block text-xs text-ink-faint">Статус</span>
                <Select
                  v-model="form.status"
                  :options="statusOptions"
                  option-label="label"
                  option-value="value"
                  fluid
                />
              </label>
              <label class="block md:col-span-2">
                <span class="mb-1.5 block text-xs text-ink-faint">Описание</span>
                <Textarea v-model="form.description" rows="3" auto-resize fluid />
              </label>
            </div>
          </TabPanel>

          <!-- CONDITIONS -->
          <TabPanel value="conditions">
            <Message severity="secondary" class="mb-4" size="small">
              Условия обмена хранятся локально — соответствующие поля появятся в API позже.
            </Message>
            <div class="grid gap-4 md:grid-cols-2">
              <label class="block">
                <span class="mb-1.5 block text-xs text-ink-faint">Минимум, ₽</span>
                <InputNumber v-model="form.min" :min="0" fluid />
              </label>
              <label class="block">
                <span class="mb-1.5 block text-xs text-ink-faint">Максимум, ₽</span>
                <InputNumber v-model="form.max" :min="0" fluid />
              </label>
              <label class="block">
                <span class="mb-1.5 block text-xs text-ink-faint">Время обмена</span>
                <InputText v-model="form.exchangeTime" fluid />
              </label>
              <div
                class="flex items-center justify-between rounded-[10px] border border-line bg-well px-3.5 py-3"
              >
                <span class="text-sm">KYC {{ form.kyc ? 'требуется' : 'не требуется' }}</span>
                <ToggleSwitch v-model="form.kyc" />
              </div>
            </div>
          </TabPanel>

          <!-- PARSER -->
          <TabPanel value="parser">
            <div class="grid gap-4 lg:grid-cols-[1.1fr_1fr]">
              <div class="flex flex-col gap-4">
                <div class="flex items-center justify-between">
                  <span class="text-sm font-bold">Парсер фида</span>
                  <span class="flex items-center gap-1.5 text-xs text-success-text">
                    <StatusDot tone="success" />healthy · 99.8%
                  </span>
                </div>
                <label class="block">
                  <span class="mb-1.5 block text-xs text-ink-faint">URL фида</span>
                  <InputText v-model="form.feedUrl" class="font-mono text-xs" fluid />
                </label>
                <div class="grid grid-cols-2 gap-4">
                  <label class="block">
                    <span class="mb-1.5 block text-xs text-ink-faint">Формат</span>
                    <Select v-model="form.format" :options="formatOptions" fluid />
                  </label>
                  <label class="block">
                    <span class="mb-1.5 block text-xs text-ink-faint">Интервал</span>
                    <Select v-model="form.interval" :options="intervalOptions" fluid />
                  </label>
                </div>
                <Button
                  :label="testing ? 'Прогон…' : 'Тестовый прогон'"
                  icon="pi pi-play"
                  :loading="testing"
                  fluid
                  @click="runTest"
                />
                <div v-if="testResult" class="rounded-xl border border-success/30 bg-well p-3.5">
                  <div class="mb-2.5 flex items-center gap-2">
                    <i class="pi pi-check-circle text-success-text" />
                    <span class="text-sm font-semibold text-success-text">Прогон успешен</span>
                    <span class="tnum ml-auto text-xs text-ink-faint">
                      {{ testResult.ms }} мс · HTTP {{ testResult.http }}
                    </span>
                  </div>
                  <div class="flex gap-5 text-xs text-ink-muted">
                    <span
                      >Направлений:
                      <span class="tnum font-semibold text-ink">{{
                        testResult.directions
                      }}</span></span
                    >
                    <span
                      >Распознано:
                      <span class="tnum font-semibold text-success-text">{{
                        testResult.recognized
                      }}</span></span
                    >
                    <span
                      >Ошибок:
                      <span class="tnum font-semibold text-ink-faint">{{
                        testResult.errors
                      }}</span></span
                    >
                  </div>
                </div>
              </div>

              <!-- logs preview -->
              <div class="overflow-hidden rounded-[14px] border border-line">
                <div class="border-b border-line px-4 py-3 text-sm font-bold">Последние записи</div>
                <div class="tnum bg-bg px-4 py-3 text-xs leading-[1.9]">
                  <div v-for="(l, i) in PARSER_LOGS.slice(0, 6)" :key="i">
                    <span class="text-ink-fainter">{{ l.time }}</span>
                    <span :class="LOG_COLOR[l.level]"> {{ l.level }} </span>
                    <span class="text-ink-muted">{{ l.text }}</span>
                  </div>
                </div>
              </div>
            </div>
          </TabPanel>

          <!-- LOGS -->
          <TabPanel value="logs">
            <div class="mb-3 flex gap-2">
              <Button
                v-for="f in ['all', 'WARN', 'ERR'] as const"
                :key="f"
                :label="f === 'all' ? 'все' : f.toLowerCase()"
                size="small"
                :severity="logFilter === f ? 'primary' : 'secondary'"
                :outlined="logFilter !== f"
                @click="logFilter = f"
              />
            </div>
            <div
              class="tnum overflow-hidden rounded-[14px] border border-line bg-bg px-4 py-3 text-xs leading-[1.9]"
            >
              <div v-for="(l, i) in logs" :key="i">
                <span class="text-ink-fainter">{{ l.time }}</span>
                <span :class="LOG_COLOR[l.level]"> {{ l.level }} </span>
                <span class="text-ink-muted">{{ l.text }}</span>
              </div>
              <div v-if="!logs.length" class="py-4 text-center text-ink-faint">
                Нет записей уровня {{ logFilter }}
              </div>
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </template>
  </div>
</template>
