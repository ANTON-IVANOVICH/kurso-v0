<script setup lang="ts">
// Preview helper: renders the 500 error screen (error.vue with a non-404 code).
// A real 500 only happens on a server fault, so this route lets you review it.
// Safe to delete before production.
const { t } = useI18n()
throw createError({ statusCode: 500, statusMessage: t('preview500.statusMessage'), fatal: true })
</script>

<template>
  <div />
</template>
