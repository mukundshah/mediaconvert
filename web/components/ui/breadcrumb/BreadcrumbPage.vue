<script lang="ts" setup>
import type { HTMLAttributes } from 'vue'

import { cn } from '@/utils/style'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <span
    aria-current="page"
    aria-disabled="true"
    data-slot="breadcrumb-page"
    role="link"
    :class="cn('text-foreground font-normal', props.class)"
  >
    <slot></slot>
  </span>
</template>
