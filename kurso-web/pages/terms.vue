<script setup lang="ts">
definePageMeta({ layout: 'legal' })

const { t } = useI18n()

const route = useRoute()
watchEffect(() => {
  route.meta.legalTitle = t('terms.title')
})

useSeoMeta({
  title: () => t('terms.seoTitle'),
  description: () => t('terms.seoDescription'),
})

interface Section {
  id: string
  num: string
  title: string
  heading?: string
  paras: string[]
  bullets?: string[]
  email?: string
}

const sections = computed<Section[]>(() => [
  {
    id: 'terms-general',
    num: '1',
    title: t('terms.generalTitle'),
    paras: [
      t('terms.generalP1'),
      t('terms.generalP2'),
      t('terms.generalP3'),
      t('terms.generalP4'),
      t('terms.generalP5'),
    ],
  },
  {
    id: 'terms-definitions',
    num: '2',
    title: t('terms.definitionsTitle'),
    paras: [t('terms.definitionsP1')],
    bullets: [
      t('terms.definitionsB1'),
      t('terms.definitionsB2'),
      t('terms.definitionsB3'),
      t('terms.definitionsB4'),
      t('terms.definitionsB5'),
      t('terms.definitionsB6'),
      t('terms.definitionsB7'),
      t('terms.definitionsB8'),
    ],
  },
  {
    id: 'terms-service',
    num: '3',
    title: t('terms.serviceTitle'),
    paras: [
      t('terms.serviceP1'),
      t('terms.serviceP2'),
      t('terms.serviceP3'),
      t('terms.serviceP4'),
      t('terms.serviceP5'),
    ],
  },
  {
    id: 'terms-data',
    num: '4',
    title: t('terms.dataTitle'),
    paras: [
      t('terms.dataP1'),
      t('terms.dataP2'),
      t('terms.dataP3'),
      t('terms.dataP4'),
      t('terms.dataP5'),
    ],
  },
  {
    id: 'terms-liability',
    num: '5',
    title: t('terms.liabilityTitle'),
    paras: [
      t('terms.liabilityP1'),
      t('terms.liabilityP2'),
      t('terms.liabilityP3'),
      t('terms.liabilityP4'),
      t('terms.liabilityP5'),
    ],
  },
  {
    id: 'terms-account',
    num: '6',
    title: t('terms.accountTitle'),
    paras: [t('terms.accountP1'), t('terms.accountP2'), t('terms.accountP3'), t('terms.accountP4')],
  },
  {
    id: 'terms-partners',
    num: '7',
    title: t('terms.partnersTitle'),
    paras: [
      t('terms.partnersP1'),
      t('terms.partnersP2'),
      t('terms.partnersP3'),
      t('terms.partnersP4'),
    ],
  },
  {
    id: 'terms-reviews',
    num: '8',
    title: t('terms.reviewsTitle'),
    paras: [
      t('terms.reviewsP1'),
      t('terms.reviewsP2'),
      t('terms.reviewsP3'),
      t('terms.reviewsP4'),
      t('terms.reviewsP5'),
    ],
  },
  {
    id: 'terms-prohibited',
    num: '9',
    title: t('terms.prohibitedTitle'),
    paras: [t('terms.prohibitedP1')],
    bullets: [
      t('terms.prohibitedB1'),
      t('terms.prohibitedB2'),
      t('terms.prohibitedB3'),
      t('terms.prohibitedB4'),
      t('terms.prohibitedB5'),
      t('terms.prohibitedB6'),
      t('terms.prohibitedB7'),
    ],
  },
  {
    id: 'terms-ip',
    num: '10',
    title: t('terms.ipTitle'),
    paras: [t('terms.ipP1'), t('terms.ipP2'), t('terms.ipP3'), t('terms.ipP4')],
  },
  {
    id: 'terms-changes',
    num: '11',
    title: t('terms.changesTitle'),
    paras: [t('terms.changesP1'), t('terms.changesP2'), t('terms.changesP3'), t('terms.changesP4')],
  },
  {
    id: 'terms-law',
    num: '12',
    title: t('terms.lawTitle'),
    paras: [t('terms.lawP1'), t('terms.lawP2'), t('terms.lawP3'), t('terms.lawP4')],
  },
  {
    id: 'terms-contacts',
    num: '13',
    title: t('terms.contactsTitle'),
    paras: [t('terms.contactsP1'), t('terms.contactsP2')],
    email: 'legal@kurso.io',
  },
])

const { activeId } = useScrollSpy(sections.value.map((s) => s.id))
</script>

<template>
  <div class="mx-auto max-w-[1100px] px-4 md:px-6">
    <!-- header -->
    <div class="pt-8 md:pt-12">
      <div class="font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint">
        {{ t('terms.eyebrow') }}
      </div>
      <h1
        class="mt-3 text-[26px] font-extrabold tracking-[-0.025em] text-ink md:text-[34px] md:font-black"
      >
        {{ t('terms.title') }}
      </h1>
      <div class="mt-2.5 text-[13px] text-ink-faint">
        {{ t('terms.updatedLabel') }} <span class="tnum">{{ t('terms.updatedDate') }}</span> ·
        {{ t('terms.effectiveNote') }}
      </div>

      <div
        class="mt-5 flex max-w-[680px] items-start gap-2.5 rounded-xl border border-warning/25 bg-warning/[0.08] px-3.5 py-3"
      >
        <svg
          class="mt-0.5 flex-none text-warning"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="12" cy="12" r="10" />
          <path d="M12 8v4M12 16h.01" />
        </svg>
        <span class="text-[13px] leading-relaxed text-warning-light">
          {{ t('terms.disclaimer') }}
        </span>
      </div>
    </div>

    <!-- content: sidebar + article -->
    <div class="mt-8 pb-4 md:grid md:grid-cols-[230px_1fr] md:gap-11">
      <LegalToc :sections="sections" :active-id="activeId" />

      <article class="max-w-[680px]">
        <section
          v-for="section in sections"
          :id="section.id"
          :key="section.id"
          class="scroll-mt-24 [&:not(:last-child)]:mb-9"
        >
          <h2 class="text-[17px] font-bold text-ink md:text-[19px]">
            {{ section.num }}. {{ section.heading ?? section.title }}
          </h2>
          <p
            v-for="(para, i) in section.paras"
            :key="i"
            class="mt-2.5 leading-relaxed text-ink-muted"
          >
            {{ para }}
          </p>
          <ul v-if="section.bullets" class="mt-3 space-y-2">
            <li
              v-for="(bullet, bi) in section.bullets"
              :key="bi"
              class="flex gap-2.5 leading-relaxed text-ink-muted"
            >
              <span class="mt-[9px] h-1 w-1 flex-none rounded-full bg-ink-ghost" />
              <span>{{ bullet }}</span>
            </li>
          </ul>
          <p v-if="section.email" class="mt-2.5 leading-relaxed text-ink-muted">
            <a :href="`mailto:${section.email}`" class="tnum text-brand-bright">{{
              section.email
            }}</a>
          </p>
        </section>
      </article>
    </div>
  </div>
</template>
