<script setup lang="ts">
const { t, plural } = useI18n()

useSeoMeta({
  title: () => t('about.seoTitle'),
  description: () => t('about.seoDescription'),
})

interface MonitorRow {
  rank: number
  code: string
  name: string
  bg: string
  rate: string
  best?: boolean
  highlight?: boolean
  tone?: 'ink' | 'muted' | 'faint'
  dim?: boolean
}

const monitorRows: MonitorRow[] = [
  {
    rank: 1,
    code: 'CB',
    name: 'CryptoBridge',
    bg: '#3A4452',
    rate: '81.20',
    best: true,
    highlight: true,
    tone: 'ink',
  },
  { rank: 2, code: 'N', name: 'NetEx24', bg: '#5B3FA0', rate: '80.95', tone: 'muted' },
  { rank: 3, code: '24', name: '24Paybank', bg: '#1F8A5B', rate: '80.74', tone: 'muted' },
  { rank: 4, code: 'BM', name: 'BaksMan', bg: '#8A5A2B', rate: '79.10', tone: 'faint', dim: true },
]

const rateClass: Record<string, string> = {
  ink: 'text-ink',
  muted: 'text-ink-muted',
  faint: 'text-ink-faint',
}
</script>

<template>
  <div class="mx-auto max-w-[1200px] px-4 py-10 md:px-6 md:py-14">
    <!-- ============ HERO ============ -->
    <section class="grid grid-cols-1 items-center gap-10 md:grid-cols-[1.12fr_0.88fr] md:gap-12">
      <div>
        <div
          class="mb-5 inline-flex items-center gap-2 font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint"
        >
          <KStatusDot tone="success" pulse :size="7" />{{ t('about.heroBadgePrefix') }}
          <span class="tnum">2021</span> {{ t('about.heroBadgeSuffix') }}
        </div>
        <h1
          class="text-[34px] font-extrabold leading-[1.05] tracking-[-0.03em] text-ink md:text-[46px]"
        >
          {{ t('about.heroTitle') }}
        </h1>
        <p class="mt-5 max-w-[520px] text-[17px] leading-relaxed text-ink-muted md:text-lg">
          {{ t('about.heroLead') }}
        </p>
        <div class="mt-6 flex items-center gap-5 md:gap-6">
          <div>
            <div class="tnum text-2xl font-extrabold tracking-[-0.02em] text-ink">128</div>
            <div class="mt-0.5 text-xs text-ink-faint">{{ plural(128, 'exchangers') }}</div>
          </div>
          <div class="h-[34px] w-px flex-none bg-line" />
          <div>
            <div class="tnum text-2xl font-extrabold tracking-[-0.02em] text-ink">84 200</div>
            <div class="mt-0.5 text-xs text-ink-faint">{{ t('about.statClicks') }}</div>
          </div>
          <div class="h-[34px] w-px flex-none bg-line" />
          <div>
            <div class="tnum text-2xl font-extrabold tracking-[-0.02em] text-ink">1.2M ₽</div>
            <div class="mt-0.5 text-xs text-ink-faint">{{ t('about.statReserves') }}</div>
          </div>
        </div>
      </div>

      <!-- live monitor panel -->
      <div class="overflow-hidden rounded-2xl border border-line bg-surface-nested shadow-panel">
        <div class="flex items-center justify-between border-b border-line px-[18px] py-[15px]">
          <div>
            <div class="text-sm font-bold text-ink">USDT → Тинькофф</div>
            <div class="tnum mt-0.5 text-xs text-ink-faint">
              <span class="tnum">1 000</span> USDT · {{ t('about.now') }}
            </div>
          </div>
          <div class="inline-flex items-center gap-1.5 text-xs text-ink-faint">
            <KStatusDot tone="success" pulse :size="7" />{{ t('about.updated') }}
            <span class="tnum">4</span> {{ t('about.sec') }}
          </div>
        </div>

        <div
          v-for="row in monitorRows"
          :key="row.rank"
          class="flex items-center gap-[11px] border-b border-line-subtle px-[18px] py-3"
          :class="[row.highlight ? 'bg-brand/[0.07]' : '', row.dim ? 'opacity-55' : '']"
        >
          <span
            class="tnum w-3.5 flex-none text-xs font-bold"
            :class="row.highlight ? 'text-brand-bright' : 'text-ink-faint'"
            >{{ row.rank }}</span
          >
          <span
            class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] text-[11px] font-bold text-white"
            :style="{ background: row.bg }"
            >{{ row.code }}</span
          >
          <div class="min-w-0 flex-1 truncate text-sm font-semibold text-ink">{{ row.name }}</div>
          <span
            v-if="row.best"
            class="rounded-[5px] bg-success/10 px-[7px] py-0.5 text-[10px] font-semibold text-success-bright"
            >{{ t('about.best') }}</span
          >
          <span class="tnum text-[15px] font-bold" :class="rateClass[row.tone || 'ink']">{{
            row.rate
          }}</span>
        </div>

        <div class="flex items-center justify-between bg-well px-[18px] py-[11px]">
          <span class="text-xs text-ink-faint">{{ t('about.bestRateTop') }}</span>
          <span class="tnum text-xs text-ink-faint"
            ><span class="tnum">128</span> {{ plural(128, 'exchangers') }}
            {{ t('about.online') }}</span
          >
        </div>
      </div>
    </section>

    <!-- ============ PRINCIPLE ============ -->
    <section class="mt-16 border-l-2 border-brand pl-6 md:mt-20 md:pl-7">
      <div class="mb-4 font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint">
        {{ t('about.principleLabel') }}
      </div>
      <p
        class="max-w-[900px] text-xl font-semibold leading-[1.5] tracking-[-0.01em] text-ink-faint md:text-[27px]"
      >
        <span class="text-ink">{{ t('about.principle1') }}</span> {{ t('about.principle2') }}
        <span class="text-ink">{{ t('about.principle3') }}</span>
      </p>
    </section>

    <!-- ============ PIPELINE ============ -->
    <section class="mt-16 md:mt-20">
      <div class="mb-6 font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint">
        {{ t('about.howLabel') }}
      </div>
      <div class="relative grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-4 md:gap-[18px]">
        <div class="absolute left-[60px] right-[60px] top-[26px] hidden h-0.5 bg-line md:block" />
        <div class="relative">
          <div
            class="relative z-[1] mb-4 flex h-[52px] w-[52px] items-center justify-center rounded-[15px] border border-line-strong bg-surface"
          >
            <span class="tnum text-[19px] font-bold text-brand-bright">01</span>
          </div>
          <div class="mb-1.5 text-base font-bold text-ink">{{ t('about.step1Title') }}</div>
          <div class="text-[13px] leading-relaxed text-ink-muted">
            <span class="tnum">116</span> {{ t('about.step1Body') }}
            <span class="text-ink">{{ t('about.step1Highlight') }}</span
            >.
          </div>
        </div>
        <div class="relative">
          <div
            class="relative z-[1] mb-4 flex h-[52px] w-[52px] items-center justify-center rounded-[15px] border border-line-strong bg-surface"
          >
            <span class="tnum text-[19px] font-bold text-brand-bright">02</span>
          </div>
          <div class="mb-1.5 text-base font-bold text-ink">{{ t('about.step2Title') }}</div>
          <div class="text-[13px] leading-relaxed text-ink-muted">
            {{ t('about.step2Body1') }}
            <span class="text-ink">{{ t('about.step2Highlight') }}</span>
            {{ t('about.step2Body2') }}
          </div>
        </div>
        <div class="relative">
          <div
            class="relative z-[1] mb-4 flex h-[52px] w-[52px] items-center justify-center rounded-[15px] border border-line-strong bg-surface"
          >
            <span class="tnum text-[19px] font-bold text-brand-bright">03</span>
          </div>
          <div class="mb-1.5 text-base font-bold text-ink">{{ t('about.step3Title') }}</div>
          <div class="text-[13px] leading-relaxed text-ink-muted">
            {{ t('about.step3Body') }} <span class="text-ink">{{ t('about.step3Highlight') }}</span
            >.
          </div>
        </div>
        <div class="relative">
          <div
            class="relative z-[1] mb-4 flex h-[52px] w-[52px] items-center justify-center rounded-[15px] border border-line-strong bg-surface"
          >
            <span class="tnum text-[19px] font-bold text-brand-bright">04</span>
          </div>
          <div class="mb-1.5 text-base font-bold text-ink">{{ t('about.step4Title') }}</div>
          <div class="text-[13px] leading-relaxed text-ink-muted">
            {{ t('about.step4Body') }}
            <span class="text-ink">{{ t('about.step4Highlight') }}</span
            >.
          </div>
        </div>
      </div>
    </section>

    <!-- ============ HONESTY + DEMONSTRATION ============ -->
    <section
      class="mt-16 grid grid-cols-1 items-center gap-8 rounded-3xl border border-line bg-surface-nested p-6 md:mt-20 md:grid-cols-[1.05fr_0.95fr] md:gap-9 md:p-9"
    >
      <div>
        <div
          class="mb-4 inline-flex items-center gap-1.5 font-label text-[11px] uppercase tracking-[0.14em] text-success-bright"
        >
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z" />
          </svg>
          {{ t('about.honestyLabel') }}
        </div>
        <h2 class="text-2xl font-extrabold leading-[1.2] tracking-[-0.02em] text-ink">
          {{ t('about.honestyTitle') }}
        </h2>
        <p class="mt-3 leading-relaxed text-ink-muted">
          {{ t('about.honestyBody1') }}
          <span class="font-semibold text-ink">{{ t('about.honestyHighlight') }}</span>
          {{ t('about.honestyBody2') }}
        </p>
      </div>

      <div class="overflow-hidden rounded-xl border border-line bg-well">
        <div class="flex items-center gap-[11px] border-b border-line-subtle px-[15px] py-[13px]">
          <span
            class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-[#26A17B] text-[11px] font-bold text-white"
            >Co</span
          >
          <div class="min-w-0 flex-1 text-sm font-semibold text-ink">Coino</div>
          <span
            class="rounded-[5px] bg-success/10 px-[7px] py-0.5 text-[10px] font-semibold text-success-bright"
            >{{ t('about.bestRate') }}</span
          >
          <span class="tnum text-[15px] font-bold text-ink">81.24</span>
        </div>
        <div class="flex items-center gap-[11px] border-b border-line-subtle px-[15px] py-[13px]">
          <span
            class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-[#3A4452] text-[11px] font-bold text-white"
            >CB</span
          >
          <div class="flex min-w-0 flex-1 items-center gap-1.5">
            <span class="text-sm font-semibold text-ink">CryptoBridge</span>
            <span
              class="rounded-[5px] bg-brand/15 px-1.5 py-0.5 text-[10px] font-semibold text-brand-bright"
              >{{ t('about.partner') }}</span
            >
          </div>
          <span class="tnum text-[15px] font-bold text-ink-muted">81.20</span>
        </div>
        <div class="flex items-center gap-[11px] px-[15px] py-[13px]">
          <span
            class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-[#5B3FA0] text-[11px] font-bold text-white"
            >N</span
          >
          <div class="flex min-w-0 flex-1 items-center gap-1.5">
            <span class="text-sm font-semibold text-ink">NetEx24</span>
            <span
              class="rounded-[5px] bg-brand/15 px-1.5 py-0.5 text-[10px] font-semibold text-brand-bright"
              >{{ t('about.partner') }}</span
            >
          </div>
          <span class="tnum text-[15px] font-bold text-ink-muted">80.95</span>
        </div>
        <div
          class="border-t border-line-subtle bg-surface-nested px-[15px] py-[11px] text-xs text-ink-faint"
        >
          {{ t('about.coinoNote') }}
        </div>
      </div>
    </section>

    <!-- ============ TIMELINE ============ -->
    <section class="mt-16 md:mt-20">
      <div class="mb-6 font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint">
        {{ t('about.growthLabel') }}
      </div>
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-4 md:gap-5">
        <div class="relative border-t-2 border-brand pt-[18px]">
          <span class="absolute -top-[5px] left-0 h-2 w-2 rounded-full bg-brand" />
          <div class="tnum mb-1.5 text-[22px] font-extrabold tracking-[-0.02em] text-ink">2021</div>
          <div class="mb-1 text-sm font-semibold text-ink">{{ t('about.launchTitle') }}</div>
          <div class="text-[13px] leading-relaxed text-ink-faint">
            {{ t('about.launchBody1') }} <span class="tnum">12</span> {{ plural(12, 'exchangers') }}
            {{ t('about.launchBody2') }}
          </div>
        </div>
        <div class="relative border-t-2 border-line-strong pt-[18px]">
          <span class="absolute -top-[5px] left-0 h-2 w-2 rounded-full bg-ink-ghost" />
          <div class="tnum mb-1.5 text-[22px] font-extrabold tracking-[-0.02em] text-ink">2023</div>
          <div class="mb-1 text-sm font-semibold text-ink">{{ t('about.mapMobileTitle') }}</div>
          <div class="text-[13px] leading-relaxed text-ink-faint">
            {{ t('about.mapMobileBody') }}
          </div>
        </div>
        <div class="relative border-t-2 border-line-strong pt-[18px]">
          <span class="absolute -top-[5px] left-0 h-2 w-2 rounded-full bg-ink-ghost" />
          <div class="tnum mb-1.5 text-[22px] font-extrabold tracking-[-0.02em] text-ink">2024</div>
          <div class="mb-1 text-sm font-semibold text-ink">{{ t('about.arbitrationTitle') }}</div>
          <div class="text-[13px] leading-relaxed text-ink-faint">
            {{ t('about.arbitrationBody') }}
          </div>
        </div>
        <div class="relative border-t-2 border-line-strong pt-[18px]">
          <span class="absolute -top-[5px] left-0 h-2 w-2 rounded-full bg-ink-ghost" />
          <div class="tnum mb-1.5 text-[22px] font-extrabold tracking-[-0.02em] text-ink">2026</div>
          <div class="mb-1 text-sm font-semibold text-ink">
            <span class="tnum">128</span> {{ plural(128, 'exchangers') }}
          </div>
          <div class="text-[13px] leading-relaxed text-ink-faint">
            <span class="tnum">84 200</span> {{ t('about.growthClicks') }}
            <span class="tnum">2 400</span> {{ plural(2400, 'reviews') }}.
          </div>
        </div>
      </div>
    </section>

    <!-- ============ CTA ============ -->
    <section class="mt-16 md:mt-20">
      <div
        class="relative flex flex-wrap items-center justify-between gap-6 overflow-hidden rounded-3xl border border-brand/25 bg-surface-nested p-8 md:p-9"
      >
        <div
          class="pointer-events-none absolute -right-8 -top-16 h-72 w-72 bg-[radial-gradient(circle,rgba(46,125,242,0.16),transparent_70%)]"
        />
        <div class="relative">
          <h2 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">
            {{ t('about.ctaTitle') }}
          </h2>
          <p class="mt-2 text-ink-muted"><span class="tnum">128</span> {{ t('about.ctaBody') }}</p>
        </div>
        <div class="relative flex flex-wrap gap-3">
          <KButton size="lg" @click="navigateTo('/')">{{ t('about.ctaCompare') }}</KButton>
          <KButton variant="secondary" size="lg" @click="navigateTo('/exchangers')">
            {{ t('about.ctaCatalog') }}
          </KButton>
        </div>
      </div>
    </section>
  </div>
</template>
