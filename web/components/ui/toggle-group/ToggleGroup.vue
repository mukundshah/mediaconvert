<script setup lang="ts">
import type { VariantProps } from 'class-variance-authority'
import type { ToggleGroupRootEmits, ToggleGroupRootProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import type { toggleVariants } from '@/components/ui/toggle'

import { reactiveOmit } from '@vueuse/core'
import { ToggleGroupRoot, useForwardPropsEmits } from 'reka-ui'
import { provide } from 'vue'

import { cn } from '@/utils/style'

type ToggleGroupVariants = VariantProps<typeof toggleVariants>

interface Props extends ToggleGroupRootProps {
  class?: HTMLAttributes['class']
  variant?: ToggleGroupVariants['variant']
  size?: ToggleGroupVariants['size']
  spacing?: number
}

const props = withDefaults(
  defineProps<Props>(),
  {
    spacing: 0,
  },
)

const emits = defineEmits<ToggleGroupRootEmits>()

provide('toggleGroup', {
  variant: props.variant,
  size: props.size,
  spacing: props.spacing,
})

const delegatedProps = reactiveOmit(props, 'class', 'size', 'variant')
const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <ToggleGroupRoot
    v-slot="slotProps"
    data-slot="toggle-group"
    v-bind="forwarded"
    :class="cn('group/toggle-group flex w-fit items-center gap-[--spacing(var(--gap))] rounded-md data-[spacing=default]:data-[variant=outline]:shadow-xs', props.class)"
    :data-size="size"
    :data-spacing="spacing"
    :data-variant="variant"
    :style="{
      '--gap': spacing,
    }"
  >
    <slot v-bind="slotProps"></slot>
  </ToggleGroupRoot>
</template>
