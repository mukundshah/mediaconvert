<script setup lang="ts">
import type { PrimitiveProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import type { SidebarMenuButtonVariants } from '.'

import { Primitive } from 'reka-ui'

import { cn } from '@/utils/style'
import { sidebarMenuButtonVariants } from '.'

export interface SidebarMenuButtonProps extends PrimitiveProps {
  variant?: SidebarMenuButtonVariants['variant']
  size?: SidebarMenuButtonVariants['size']
  isActive?: boolean
  class?: HTMLAttributes['class']
}

const props = withDefaults(
  defineProps<SidebarMenuButtonProps>(),
  {
    as: 'button',
    variant: 'default',
    size: 'default',
  },
)
</script>

<template>
  <Primitive
    data-sidebar="menu-button"
    data-slot="sidebar-menu-button"
    :as="as"
    :as-child="asChild"
    :class="cn(sidebarMenuButtonVariants({ variant, size }), props.class)"
    :data-active="isActive"
    :data-size="size"
    v-bind="$attrs"
  >
    <slot></slot>
  </Primitive>
</template>
