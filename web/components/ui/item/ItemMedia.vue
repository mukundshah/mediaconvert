<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import type { ItemMediaVariants } from '.'

import { cn } from '@/utils/style'
import { itemMediaVariants } from '.'

const props = defineProps<{
  class?: HTMLAttributes['class']
  variant?: ItemMediaVariants['variant']
}>()
</script>

<template>
  <div
    data-slot="item-media"
    :class="cn(itemMediaVariants({ variant }), props.class)"
    :data-variant="props.variant"
  >
    <slot></slot>
  </div>
</template>
