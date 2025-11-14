<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { Loader2Icon } from 'lucide-vue-next'

import { cn } from '@/utils/style'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <Loader2Icon aria-label="Loading" role="status" :class="cn('size-4 animate-spin', props.class)" />
</template>
