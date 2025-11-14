<script lang="ts" setup>
import type { HTMLAttributes } from 'vue'

import { ErrorMessage } from 'vee-validate'
import { toValue } from 'vue'

import { cn } from '@/utils/style'
import { useFormField } from './useFormField'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()

const { name, formMessageId } = useFormField()
</script>

<template>
  <ErrorMessage
    :id="formMessageId"
    as="p"
    data-slot="form-message"
    :class="cn('text-destructive text-sm', props.class)"
    :name="toValue(name)"
  />
</template>
