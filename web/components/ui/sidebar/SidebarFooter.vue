<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { cn } from '@/utils/style'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <div
    data-sidebar="footer"
    data-slot="sidebar-footer"
    :class="cn('flex flex-col gap-2 p-2', props.class)"
  >
    <slot></slot>
  </div>
</template>
