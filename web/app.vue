<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'
import 'vue-sonner/style.css'
</script>

<template>
  <NuxtRouteAnnouncer />
  <NuxtLoadingIndicator color="#000" />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <Toaster />
</template>
