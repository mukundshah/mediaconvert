<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import type { EmptyMediaVariants } from '.'

import { cn } from '@/utils/style'
import { emptyMediaVariants } from '.'

const props = defineProps<{
  class?: HTMLAttributes['class']
  variant?: EmptyMediaVariants['variant']
}>()
</script>

<template>
  <div
    data-slot="empty-icon"
    :class="cn(emptyMediaVariants({ variant }), props.class)"
    :data-variant="variant"
  >
    <slot></slot>
  </div>
</template>
