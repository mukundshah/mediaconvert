<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import type { InputGroupVariants } from '.'

import { cn } from '@/utils/style'
import { inputGroupAddonVariants } from '.'

const props = withDefaults(
  defineProps<{
    align?: InputGroupVariants['align']
    class?: HTMLAttributes['class']
  }>(),
  {
    align: 'inline-start',
  },
)

const handleInputGroupAddonClick = (e: MouseEvent) => {
  const currentTarget = e.currentTarget as HTMLElement | null
  const target = e.target as HTMLElement | null
  if (target && target.closest('button')) {
    return
  }
  if (currentTarget && currentTarget?.parentElement) {
    currentTarget.parentElement?.querySelector('input')?.focus()
  }
}
</script>

<template>
  <div
    data-slot="input-group-addon"
    role="group"
    :class="cn(inputGroupAddonVariants({ align: props.align }), props.class)"
    :data-align="props.align"
    @click="handleInputGroupAddonClick"
  >
    <slot></slot>
  </div>
</template>
