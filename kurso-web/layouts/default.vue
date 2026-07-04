<template>
  <div class="relative min-h-screen overflow-x-clip">
    <GlowBackdrop />
    <SiteHeader />
    <main class="relative z-10 pb-28 md:pb-8">
      <slot />
    </main>
    <SiteFooter />
    <MobileBottomNav />

    <!-- Global overlays: currency picker, site search, notifications -->
    <CurrencyPickerDrawer />
    <SearchOverlay />
    <NotificationsDrawer />
  </div>
</template>
