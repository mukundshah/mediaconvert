<script setup lang="ts">
import type { DropdownMenuLabelProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'

import { reactiveOmit } from '@vueuse/core'
import { DropdownMenuLabel, useForwardProps } from 'reka-ui'

import { cn } from '@/utils/style'

const props = defineProps<DropdownMenuLabelProps & { class?: HTMLAttributes['class'], inset?: boolean }>()

const delegatedProps = reactiveOmit(props, 'class', 'inset')
const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <DropdownMenuLabel
    data-slot="dropdown-menu-label"
    v-bind="forwardedProps"
    :class="cn('px-2 py-1.5 text-sm font-medium data-[inset]:pl-8', props.class)"
    :data-inset="inset ? '' : undefined"
  >
    <slot></slot>
  </DropdownMenuLabel>
</template>
