<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import type { FieldVariants } from '.'

import { cn } from '@/utils/style'
import { fieldVariants } from '.'

const props = defineProps<{
  class?: HTMLAttributes['class']
  orientation?: FieldVariants['orientation']
}>()
</script>

<template>
  <div
    data-slot="field"
    role="group"
    :class="cn(fieldVariants({ orientation }), props.class)"
    :data-orientation="orientation"
  >
    <slot></slot>
  </div>
</template>
