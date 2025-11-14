<script setup lang="ts">
import type { AvatarImageProps } from 'reka-ui'
import { AvatarImage } from 'reka-ui'

const props = defineProps<AvatarImageProps>()
</script>

<template>
  <AvatarImage
    v-bind="props"
    class="aspect-square size-full"
    data-slot="avatar-image"
  >
    <slot></slot>
  </AvatarImage>
</template>
