<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import type { OTPInputEmits, OTPInputProps } from 'vue-input-otp'

import { reactiveOmit } from '@vueuse/core'
import { useForwardPropsEmits } from 'reka-ui'
import { OTPInput } from 'vue-input-otp'

import { cn } from '@/utils/style'

const props = defineProps<OTPInputProps & { class?: HTMLAttributes['class'] }>()
const emits = defineEmits<OTPInputEmits>()

const delegatedProps = reactiveOmit(props, 'class')

const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <OTPInput
    v-slot="slotProps"
    v-bind="forwarded"
    class="disabled:cursor-not-allowed"
    data-slot="input-otp"
    :container-class="cn('flex items-center gap-2 has-disabled:opacity-50', props.class)"
  >
    <slot v-bind="slotProps"></slot>
  </OTPInput>
</template>
