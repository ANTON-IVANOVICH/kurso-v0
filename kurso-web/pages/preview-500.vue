<script setup lang="ts">
// Preview helper: renders the 500 error screen (error.vue with a non-404 code).
// A real 500 only happens on a server fault, so this route lets you review it.
// Safe to delete before production.
throw createError({ statusCode: 500, statusMessage: 'Курсы временно недоступны', fatal: true })
</script>

<template>
  <div />
</template>
