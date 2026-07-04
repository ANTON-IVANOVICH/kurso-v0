<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useProfileQuery } from '../composables/useMerchant'
import { useAuthStore } from '../stores/auth'
import { api, ApiError } from '../lib/api'
import { fmtCompact, toNum } from '../lib/format'
import PageHeader from '../components/ui/PageHeader.vue'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import type { MerchantProfile } from '../types/models'

const profile = useProfileQuery()
const auth = useAuthStore()
const toast = useToast()

const form = ref({ name: '', description: '', websiteUrl: '' })
const busy = ref(false)
const error = ref('')
const canEdit = computed(() => auth.merchant?.role !== 'viewer')

// Seed the form from the loaded profile and reset dirty tracking.
watch(
  () => profile.data.value,
  (p) => {
    if (p)
      form.value = {
        name: p.name,
        description: p.description ?? '',
        websiteUrl: p.websiteUrl ?? '',
      }
  },
  { immediate: true },
)

const dirty = computed(() => {
  const p = profile.data.value
  if (!p) return false
  return (
    form.value.name !== p.name ||
    form.value.description !== (p.description ?? '') ||
    form.value.websiteUrl !== (p.websiteUrl ?? '')
  )
})

function reset() {
  const p = profile.data.value
  if (p)
    form.value = { name: p.name, description: p.description ?? '', websiteUrl: p.websiteUrl ?? '' }
  error.value = ''
}

async function save() {
  error.value = ''
  if (!form.value.name.trim()) {
    error.value = 'Название не может быть пустым'
    return
  }
  busy.value = true
  try {
    const updated = await api.patch<MerchantProfile>('/partner/profile', {
      name: form.value.name.trim(),
      description: form.value.description.trim(),
      websiteUrl: form.value.websiteUrl.trim(),
    })
    await profile.refetch()
    form.value = {
      name: updated.name,
      description: updated.description ?? '',
      websiteUrl: updated.websiteUrl ?? '',
    }
    toast.add({
      severity: 'success',
      summary: 'Профиль сохранён',
      detail: 'Карточка отправлена на повторную модерацию',
      life: 3000,
    })
  } catch (e) {
    error.value = e instanceof ApiError ? e.message : 'Не удалось сохранить'
  } finally {
    busy.value = false
  }
}

const p = computed(() => profile.data.value)
</script>

<template>
  <div>
    <Toast position="bottom-right" />
    <PageHeader
      title="Профиль обменника"
      subtitle="Описание, контакты и условия — видны пользователям каталога"
    />

    <div class="grid gap-4 lg:grid-cols-[1.4fr_1fr]">
      <!-- editable profile -->
      <div class="rounded-[14px] border border-line bg-surface p-[22px]">
        <div class="mb-4 flex items-center gap-4">
          <span
            class="flex h-16 w-16 flex-none items-center justify-center rounded-2xl bg-[#3A4452] text-[22px] font-extrabold text-white"
            >{{ auth.initials }}</span
          >
          <div>
            <div class="text-[15px] font-bold">{{ p?.name ?? '—' }}</div>
            <div class="tnum mt-0.5 text-xs text-ink-faint">/{{ p?.slug ?? '' }}</div>
          </div>
        </div>

        <div class="flex flex-col gap-4">
          <label class="block">
            <span class="mb-1.5 block text-xs text-ink-faint">Название</span>
            <InputText v-model="form.name" :disabled="!canEdit" fluid />
          </label>
          <label class="block">
            <span class="mb-1.5 block text-xs text-ink-faint">Описание</span>
            <Textarea
              v-model="form.description"
              :disabled="!canEdit"
              rows="3"
              auto-resize
              class="w-full"
            />
          </label>
          <label class="block">
            <span class="mb-1.5 block text-xs text-ink-faint">Сайт</span>
            <InputText
              v-model="form.websiteUrl"
              :disabled="!canEdit"
              placeholder="https://…"
              fluid
            />
          </label>
        </div>

        <Message v-if="error" severity="error" size="small" class="mt-3">{{ error }}</Message>
        <Message v-if="!canEdit" severity="secondary" size="small" class="mt-3"
          >Роль «viewer» — редактирование недоступно.</Message
        >

        <div v-if="canEdit" class="mt-4 flex justify-end gap-2 border-t border-line pt-4">
          <Button
            label="Отмена"
            severity="secondary"
            size="small"
            :disabled="!dirty || busy"
            @click="reset"
          />
          <Button label="Сохранить" size="small" :disabled="!dirty" :loading="busy" @click="save" />
        </div>
      </div>

      <!-- read-only metrics + verification -->
      <div class="rounded-[14px] border border-line bg-surface p-[22px]">
        <div class="mb-4 text-[15px] font-bold">Условия и статус</div>
        <div class="flex flex-col">
          <div class="flex items-center justify-between border-b border-line py-3">
            <span class="text-[13px] text-ink-muted">Статус карточки</span>
            <Tag
              :value="p?.status === 'active' ? 'активна' : (p?.status ?? '—')"
              :severity="p?.status === 'active' ? 'success' : 'warn'"
            />
          </div>
          <div class="flex items-center justify-between border-b border-line py-3">
            <span class="text-[13px] text-ink-muted">Модерация</span>
            <Tag
              :value="p?.isVerified ? 'проверен' : 'на модерации'"
              :severity="p?.isVerified ? 'success' : 'warn'"
            />
          </div>
          <div class="flex items-center justify-between border-b border-line py-3">
            <span class="text-[13px] text-ink-muted">Рейтинг</span>
            <span class="tnum text-sm font-semibold">{{ p?.ratingAvg ?? '—' }} ★</span>
          </div>
          <div class="flex items-center justify-between border-b border-line py-3">
            <span class="text-[13px] text-ink-muted">Отзывов</span>
            <span class="tnum text-sm font-semibold">{{ p?.reviewsCount ?? 0 }}</span>
          </div>
          <div class="flex items-center justify-between border-b border-line py-3">
            <span class="text-[13px] text-ink-muted">Суммарный резерв</span>
            <span class="tnum text-sm font-semibold">{{
              p?.reserveTotal ? fmtCompact(toNum(p.reserveTotal)) : '—'
            }}</span>
          </div>
          <div class="flex items-center justify-between py-3">
            <span class="text-[13px] text-ink-muted">Активы</span>
            <span class="tnum text-sm font-semibold">{{
              (p?.assets ?? []).join(' · ') || '—'
            }}</span>
          </div>
        </div>

        <div class="mt-4 flex items-start gap-2.5 rounded-xl border border-line bg-well p-3.5">
          <i class="pi pi-info-circle mt-0.5 text-brand-light" />
          <span class="text-[12px] leading-relaxed text-ink-muted">
            Изменение профиля отправляет карточку на повторную модерацию Kurso.
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
