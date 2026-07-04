<template>
  <div class="relative min-h-screen overflow-x-clip">
    <GlowBackdrop />
    <SiteHeader />
    <main class="relative z-10 pb-28 md:pb-12">
      <div class="mx-auto max-w-[1200px] px-4 py-6 md:px-6">
        <div class="gap-6 md:grid md:grid-cols-[240px_minmax(0,1fr)]">
          <AccountSidebar class="hidden md:block" />
          <div class="min-w-0"><slot /></div>
        </div>
      </div>
    </main>
    <MobileBottomNav />

    <!-- Global overlays reused from the default shell (currency picker powers the
         alert builder; notifications/search stay reachable). -->
    <CurrencyPickerDrawer />
    <SearchOverlay />
    <NotificationsDrawer />
  </div>
</template>
