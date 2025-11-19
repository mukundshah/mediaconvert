<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { cn } from '@/utils/style'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <div data-slot="item-group" role="list" :class="cn('group/item-group flex flex-col', props.class)">
    <slot></slot>
  </div>
</template>
