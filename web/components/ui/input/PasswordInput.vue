<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { useVModel } from '@vueuse/core'
import { EyeIcon, EyeOffIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import { cn } from '@/utils/style'
import Input from './Input.vue'

defineOptions({
  inheritAttrs: false,
})

const props = defineProps<{
  defaultValue?: string
  modelValue?: string
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  'update:modelValue': [payload: string]
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue,
})

const isVisible = ref(false)

const toggleVisibility = () => {
  isVisible.value = !isVisible.value
}
</script>

<template>
  <div class="relative">
    <Input
      v-model="modelValue"
      :class="cn('pe-9', props.class)"
      :type="isVisible ? 'text' : 'password'"
      v-bind="$attrs"
    />

    <Button
      size="icon"
      variant="ghost"
      :class="cn('absolute inset-y-0 end-0', props.class)"
      @click="toggleVisibility"
    >
      <EyeIcon v-if="isVisible" aria-hidden="true" :size="16" />
      <EyeOffIcon v-else aria-hidden="true" :size="16" />
    </Button>
  </div>
</template>
