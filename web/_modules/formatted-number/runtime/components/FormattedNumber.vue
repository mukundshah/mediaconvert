<script setup lang="ts">
import type { Decimal } from 'decimal.js'
import type { NumberFormatterOptions } from '../composables/useNumberFormatter'

import { useNumberFormatter } from '../composables/useNumberFormatter'

const {
  value,
  locale = 'en-US',
  options,
  nullValue = '-',
  zeroValue,
} = defineProps<NumberFormatterOptions & { value: number | string | Decimal | null | undefined }>()

const { format } = useNumberFormatter(
  {
    locale,
    options,
    nullValue,
    zeroValue,
  },
)

const formattedValue = computed(() => format(value))
</script>

<template>
  <span>{{ formattedValue }}</span>
</template>
