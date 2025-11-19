<script setup lang="ts">
import type { PaginationListItemProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import type { ButtonVariants } from '@/components/ui/button'

import { reactiveOmit } from '@vueuse/core'
import { PaginationListItem } from 'reka-ui'

import { buttonVariants } from '@/components/ui/button'
import { cn } from '@/utils/style'

interface Props extends PaginationListItemProps {
  size?: ButtonVariants['size']
  class?: HTMLAttributes['class']
  isActive?: boolean
}

const props = withDefaults(
  defineProps<Props>(),
  {
    size: 'icon',
  },
)

const delegatedProps = reactiveOmit(props, 'class', 'size', 'isActive')
</script>

<template>
  <PaginationListItem
    data-slot="pagination-item"
    v-bind="delegatedProps"
    :class="cn(
      buttonVariants({
        variant: isActive ? 'outline' : 'ghost',
        size,
      }),
      props.class)"
  >
    <slot></slot>
  </PaginationListItem>
</template>
