<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/utils/style'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <div class="relative w-full overflow-auto" data-slot="table-container">
    <table data-slot="table" :class="cn('w-full caption-bottom text-sm', props.class)">
      <slot></slot>
    </table>
  </div>
</template>
