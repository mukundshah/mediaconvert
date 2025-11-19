<script lang="ts" setup>
import type { LabelProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'

import { Label } from '@/components/ui/label'
import { cn } from '@/utils/style'
import { useFormField } from './useFormField'

const props = defineProps<LabelProps & { class?: HTMLAttributes['class'] }>()

const { error, formItemId } = useFormField()
</script>

<template>
  <Label
    data-slot="form-label"
    :class="cn('data-[error=true]:text-destructive', props.class)"
    :data-error="!!error"
    :for="formItemId"
  >
    <slot></slot>
  </Label>
</template>
