<script setup lang="ts">
definePageMeta({ layout: 'legal', legalTitle: 'Политика конфиденциальности' })

const { t } = useI18n()

useSeoMeta({
  title: () => t('privacy.seoTitle'),
  description: () => t('privacy.seoDesc'),
})

interface Section {
  id: string
  num: string
  title: string
  heading?: string
  paras: string[]
  bullets?: string[]
  email?: string
  cookiesLink?: boolean
}

const sections = computed<Section[]>(() => [
  {
    id: 'privacy-intro',
    num: '1',
    title: t('privacy.introTitle'),
    paras: [t('privacy.introP1'), t('privacy.introP2'), t('privacy.introP3')],
  },
  {
    id: 'privacy-collect',
    num: '2',
    title: t('privacy.collectTitle'),
    heading: t('privacy.collectHeading'),
    paras: [t('privacy.collectP1')],
    bullets: [
      t('privacy.collectB1'),
      t('privacy.collectB2'),
      t('privacy.collectB3'),
      t('privacy.collectB4'),
      t('privacy.collectB5'),
    ],
  },
  {
    id: 'privacy-basis',
    num: '3',
    title: t('privacy.basisTitle'),
    heading: t('privacy.basisHeading'),
    paras: [
      t('privacy.basisP1'),
      t('privacy.basisP2'),
      t('privacy.basisP3'),
      t('privacy.basisP4'),
      t('privacy.basisP5'),
    ],
  },
  {
    id: 'privacy-use',
    num: '4',
    title: t('privacy.useTitle'),
    paras: [t('privacy.useP1')],
    bullets: [
      t('privacy.useB1'),
      t('privacy.useB2'),
      t('privacy.useB3'),
      t('privacy.useB4'),
      t('privacy.useB5'),
      t('privacy.useB6'),
      t('privacy.useB7'),
    ],
  },
  {
    id: 'privacy-cookies',
    num: '5',
    title: t('privacy.cookiesTitle'),
    paras: [t('privacy.cookiesP1'), t('privacy.cookiesP2')],
    cookiesLink: true,
  },
  {
    id: 'privacy-third-party',
    num: '6',
    title: t('privacy.thirdTitle'),
    paras: [t('privacy.thirdP1')],
    bullets: [
      t('privacy.thirdB1'),
      t('privacy.thirdB2'),
      t('privacy.thirdB3'),
      t('privacy.thirdB4'),
    ],
  },
  {
    id: 'privacy-transfer',
    num: '7',
    title: t('privacy.transferTitle'),
    paras: [t('privacy.transferP1'), t('privacy.transferP2')],
  },
  {
    id: 'privacy-storage',
    num: '8',
    title: t('privacy.storageTitle'),
    paras: [t('privacy.storageP1'), t('privacy.storageP2'), t('privacy.storageP3')],
  },
  {
    id: 'privacy-rights',
    num: '9',
    title: t('privacy.rightsTitle'),
    paras: [t('privacy.rightsP1')],
    bullets: [
      t('privacy.rightsB1'),
      t('privacy.rightsB2'),
      t('privacy.rightsB3'),
      t('privacy.rightsB4'),
      t('privacy.rightsB5'),
      t('privacy.rightsB6'),
      t('privacy.rightsB7'),
    ],
  },
  {
    id: 'privacy-children',
    num: '10',
    title: t('privacy.childrenTitle'),
    paras: [t('privacy.childrenP1'), t('privacy.childrenP2')],
  },
  {
    id: 'privacy-changes',
    num: '11',
    title: t('privacy.changesTitle'),
    paras: [t('privacy.changesP1'), t('privacy.changesP2')],
  },
  {
    id: 'privacy-contacts',
    num: '12',
    title: t('privacy.contactsTitle'),
    paras: [t('privacy.contactsP1'), t('privacy.contactsP2')],
    email: 'privacy@kurso.io',
  },
])

const { activeId } = useScrollSpy(sections.value.map((s) => s.id))
</script>

<template>
  <div class="mx-auto max-w-[1100px] px-4 md:px-6">
    <!-- header -->
    <div class="pt-8 md:pt-12">
      <div class="font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint">
        {{ t('privacy.eyebrow') }}
      </div>
      <h1
        class="mt-3 text-[26px] font-extrabold tracking-[-0.025em] text-ink md:text-[34px] md:font-black"
      >
        {{ t('privacy.pageTitle') }}
      </h1>
      <div class="mt-2.5 text-[13px] text-ink-faint">
        {{ t('privacy.updated') }} <span class="tnum">{{ t('privacy.updatedDate') }}</span>
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
          <p v-if="section.cookiesLink" class="mt-2.5 leading-relaxed text-ink-muted">
            {{ t('privacy.cookiesManagePre') }}
            <NuxtLink to="/cookies" class="tnum text-brand-bright">Cookies</NuxtLink>
            {{ t('privacy.cookiesManagePost') }}
          </p>
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
