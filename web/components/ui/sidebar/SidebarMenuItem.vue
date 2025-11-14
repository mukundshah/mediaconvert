<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { cn } from '@/utils/style'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <li
    data-sidebar="menu-item"
    data-slot="sidebar-menu-item"
    :class="cn('group/menu-item relative', props.class)"
  >
    <slot></slot>
  </li>
</template>
