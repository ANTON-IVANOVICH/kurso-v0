<template>
  <div class="relative flex min-h-[100dvh] flex-col overflow-x-clip bg-canvas">
    <GlowBackdrop />
    <div class="relative z-10 flex flex-1 items-center justify-center px-4 py-8 md:py-12">
      <div class="w-full max-w-[420px]">
        <slot />
      </div>
    </div>
  </div>
</template>
