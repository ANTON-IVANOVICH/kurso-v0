<script setup lang="ts">
import { ref } from 'vue'

useSeoMeta({
  title: 'Контакты — Kurso',
  description:
    'Свяжитесь с командой Kurso: вопрос, спор по обмену или предложение о сотрудничестве. Обычно отвечаем в течение 6 часов.',
})

const name = ref('')
const email = ref('')
const subject = ref('Общий вопрос')
const message = ref('')
const sent = ref(false)

const subjects = ['Общий вопрос', 'Спор по обмену', 'Сотрудничество', 'Обменникам', 'Другое']

function onSubmit() {
  sent.value = true
}
</script>

<template>
  <div class="mx-auto max-w-[1200px] px-4 py-10 md:px-6 md:py-14">
    <h1 class="text-[28px] font-extrabold tracking-[-0.025em] text-ink md:text-[32px]">
      Свяжитесь с нами
    </h1>
    <p class="mt-2.5 max-w-[560px] leading-relaxed text-ink-muted">
      Вопрос, спор по обмену или предложение о сотрудничестве — напишите, обычно отвечаем в течение
      <span class="tnum">6</span> часов.
    </p>

    <div class="mt-8 grid grid-cols-1 items-start gap-6 md:grid-cols-[1.3fr_1fr]">
      <!-- form -->
      <form class="rounded-2xl border border-line bg-surface p-6" @submit.prevent="onSubmit">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <label class="block">
            <span class="mb-1.5 block text-xs text-ink-faint">Имя</span>
            <input
              v-model="name"
              type="text"
              placeholder="Ваше имя"
              class="w-full rounded-md border border-line bg-well px-[13px] py-[11px] text-sm text-ink placeholder:text-ink-faint focus:border-brand focus:outline-none"
            />
          </label>
          <label class="block">
            <span class="mb-1.5 block text-xs text-ink-faint">Email</span>
            <input
              v-model="email"
              type="email"
              placeholder="you@email.com"
              class="tnum w-full rounded-md border border-line bg-well px-[13px] py-[11px] text-sm text-ink placeholder:text-ink-faint focus:border-brand focus:outline-none"
            />
          </label>
        </div>

        <label class="mt-4 block">
          <span class="mb-1.5 block text-xs text-ink-faint">Тема</span>
          <div class="relative">
            <select
              v-model="subject"
              class="w-full appearance-none rounded-md border border-line bg-well px-[13px] py-[11px] pr-10 text-sm text-ink focus:border-brand focus:outline-none"
            >
              <option v-for="s in subjects" :key="s" :value="s">{{ s }}</option>
            </select>
            <svg
              class="pointer-events-none absolute right-3.5 top-1/2 -translate-y-1/2 text-ink-faint"
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.4"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M6 9l6 6 6-6" />
            </svg>
          </div>
        </label>

        <label class="mt-4 block">
          <span class="mb-1.5 block text-xs text-ink-faint">Сообщение</span>
          <textarea
            v-model="message"
            rows="5"
            placeholder="Опишите вопрос как можно подробнее…"
            class="w-full resize-y rounded-md border border-line bg-well px-[13px] py-[11px] text-sm leading-relaxed text-ink placeholder:text-ink-faint focus:border-brand focus:outline-none"
          />
        </label>

        <div class="mt-5 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <span class="max-w-[280px] text-xs leading-relaxed text-ink-faint">
            Отправляя форму, вы соглашаетесь с
            <NuxtLink to="/privacy" class="text-brand-bright">политикой конфиденциальности</NuxtLink
            >.
          </span>
          <KButton type="submit">Отправить</KButton>
        </div>

        <p v-if="sent" class="mt-4 flex items-center gap-2 text-[13px] text-success-bright">
          <KStatusDot tone="success" />
          Спасибо! Мы получили ваше сообщение и ответим в течение <span class="tnum">6</span> часов.
        </p>
      </form>

      <!-- contacts -->
      <div class="flex flex-col gap-3.5">
        <div class="rounded-2xl border border-line bg-surface p-5">
          <a
            href="mailto:support@kurso.io"
            class="flex items-center gap-3 transition-colors hover:opacity-90"
          >
            <span
              class="flex h-10 w-10 flex-none items-center justify-center rounded-[11px] bg-brand/[0.12] text-brand-bright"
            >
              <svg
                width="19"
                height="19"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <rect x="3" y="5" width="18" height="14" rx="2" />
                <path d="m3 7 9 6 9-6" />
              </svg>
            </span>
            <span>
              <span class="block text-xs text-ink-faint">Email</span>
              <span class="tnum block font-semibold text-ink">support@kurso.io</span>
            </span>
          </a>

          <a
            href="https://t.me/kurso_support"
            target="_blank"
            rel="noopener"
            class="mt-3.5 flex items-center gap-3 transition-colors hover:opacity-90"
          >
            <span
              class="flex h-10 w-10 flex-none items-center justify-center rounded-[11px] bg-brand/[0.12] text-brand-bright"
            >
              <svg
                width="19"
                height="19"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M21 4 3 11l6 2 2 6 4-6 6-9Z" />
              </svg>
            </span>
            <span>
              <span class="block text-xs text-ink-faint">Telegram</span>
              <span class="tnum block font-semibold text-ink">@kurso_support</span>
            </span>
          </a>
        </div>

        <!-- exchanger CTA -->
        <div
          class="relative overflow-hidden rounded-2xl border border-brand/25 bg-surface-nested p-5"
        >
          <div
            class="pointer-events-none absolute -right-6 -top-8 h-36 w-36 bg-[radial-gradient(circle,rgba(46,125,242,0.16),transparent_70%)]"
          />
          <div class="relative">
            <div class="text-base font-bold text-ink">Вы обменник?</div>
            <p class="mt-1.5 text-[13px] leading-relaxed text-ink-muted">
              Добавьте свой сервис в каталог Kurso и получайте целевой трафик.
            </p>
            <KButton block class="mt-4">Добавить обменник</KButton>
          </div>
        </div>

        <!-- response time -->
        <div
          class="flex items-center gap-3 rounded-2xl border border-line bg-surface px-5 py-4 text-[13px] text-ink-muted"
        >
          <KStatusDot tone="success" pulse />
          <span>Среднее время ответа — <span class="font-semibold text-ink">6 часов</span></span>
        </div>
      </div>
    </div>
  </div>
</template>
