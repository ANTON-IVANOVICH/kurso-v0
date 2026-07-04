<script setup lang="ts">
// Mobile-only floating bottom navigation with a central brand FAB.
const { openNotifications } = useNotifications()
const { t } = useI18n()
const route = useRoute()

// Active (non-central) item highlights blue with a 2× thicker icon stroke.
const isActive = (to: string) => route.path.startsWith(to)
</script>

<template>
  <nav class="fixed inset-x-0 bottom-0 z-40 px-4 pb-4 md:hidden">
    <div
      class="mx-auto flex max-w-md items-center justify-around rounded-full border border-line bg-surface-nav px-5 py-3.5 shadow-toast"
    >
      <NuxtLink
        to="/exchangers"
        :aria-label="t('bottomNav.exchangers')"
        class="flex"
        :class="isActive('/exchangers') ? 'text-brand' : 'text-[#767D86]'"
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          :stroke-width="isActive('/exchangers') ? 4 : 2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M9 6h11M9 12h11M9 18h11" />
          <circle cx="4.5" cy="6" r="1.3" />
          <circle cx="4.5" cy="12" r="1.3" />
          <circle cx="4.5" cy="18" r="1.3" />
        </svg>
      </NuxtLink>
      <NuxtLink
        to="/map"
        :aria-label="t('bottomNav.map')"
        class="flex"
        :class="isActive('/map') ? 'text-brand' : 'text-[#767D86]'"
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          :stroke-width="isActive('/map') ? 4 : 2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M9 4 3 6.5V20l6-2.5 6 2.5 6-2.5V4l-6 2.5Z" />
          <path d="M9 4v13.5" />
          <path d="M15 6.5V20" />
        </svg>
      </NuxtLink>
      <NuxtLink
        to="/"
        :aria-label="t('bottomNav.exchange')"
        class="flex h-12 w-12 items-center justify-center rounded-full bg-[linear-gradient(150deg,#4A90F5,#2E7DF2)] shadow-[0_6px_18px_rgba(46,125,242,0.5)]"
      >
        <svg
          width="22"
          height="22"
          viewBox="0 0 24 24"
          fill="none"
          stroke="#fff"
          stroke-width="2.2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M7 4 L7 20 M7 20 L4 17 M7 20 L10 17 M17 20 L17 4 M17 4 L14 7 M17 4 L20 7" />
        </svg>
      </NuxtLink>
      <button type="button" :aria-label="t('bottomNav.favorites')" class="flex text-[#767D86]">
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path
            d="M12 20s-6.5-4-9-8.5A4.5 4.5 0 0 1 12 7a4.5 4.5 0 0 1 9 4.5C18.5 16 12 20 12 20Z"
          />
        </svg>
      </button>
      <button
        type="button"
        :aria-label="t('bottomNav.notifications')"
        class="flex text-[#767D86]"
        @click="openNotifications"
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
          <path d="M10.5 20a1.8 1.8 0 0 0 3 0" />
        </svg>
      </button>
    </div>
  </nav>
</template>
