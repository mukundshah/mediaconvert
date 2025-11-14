<script setup lang="ts">
import type { InputGroupButtonProps } from '.'

import { Button } from '@/components/ui/button'
import { cn } from '@/utils/style'
import { inputGroupButtonVariants } from '.'

const props = withDefaults(
  defineProps<InputGroupButtonProps>(),
  {
    size: 'xs',
    variant: 'ghost',
  },
)
</script>

<template>
  <Button
    :class="cn(inputGroupButtonVariants({ size: props.size }), props.class)"
    :data-size="props.size"
    :variant="props.variant"
  >
    <slot></slot>
  </Button>
</template>
