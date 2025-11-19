<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import type { ButtonGroupVariants } from '.'

import { cn } from '@/utils/style'
import { buttonGroupVariants } from '.'

const props = defineProps<{
  class?: HTMLAttributes['class']
  orientation?: ButtonGroupVariants['orientation']
}>()
</script>

<template>
  <div
    data-slot="button-group"
    role="group"
    :class="cn(buttonGroupVariants({ orientation: props.orientation }), props.class)"
    :data-orientation="props.orientation"
  >
    <slot></slot>
  </div>
</template>
