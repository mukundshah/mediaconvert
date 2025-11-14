<script setup lang="ts">
import type { IconCategory } from '@/constants/icons'
import { Input } from '@/components/ui/input'
import { AVAILABLE_ICONS } from '@/constants/icons'

defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const open = ref(false)
const searchQuery = ref('')

// Filter icons based on search query
const filteredCategories = computed(() => {
  if (!searchQuery.value.trim()) {
    return AVAILABLE_ICONS
  }

  const query = searchQuery.value.toLowerCase().trim()
  return AVAILABLE_ICONS
    .map((category: IconCategory) => ({
      ...category,
      icons: category.icons.filter(icon =>
        icon.toLowerCase().includes(query)
        || category.category.toLowerCase().includes(query),
      ),
    }))
    .filter((category: IconCategory) => category.icons.length > 0)
})

function selectIcon(icon: string) {
  emit('update:modelValue', icon)
  open.value = false
  searchQuery.value = ''
}
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <button
        class="size-10 bg-[#ECF4E9] text-[#1E4841] rounded-full flex items-center justify-center text-lg shrink-0 hover:bg-[#D9E8D4] transition-colors cursor-pointer"
        type="button"
      >
        <Icon :name="modelValue" />
      </button>
    </PopoverTrigger>
    <PopoverContent align="start" class="w-82 p-0" side="bottom">
      <div class="py-3 space-y-3 overflow-clip">
        <div class="space-y-2 px-3">
          <h4 class="font-semibold text-sm text-[#1E4841] sr-only">
            Select Icon
          </h4>
          <div class="relative">
            <Input
              v-model="searchQuery"
              class="w-full pr-8"
              placeholder="Search icons..."
              type="search"
            />
            <Icon
              class="absolute right-2.5 top-1/2 -translate-y-1/2 size-4 text-[#6B7271] pointer-events-none"
              name="lucide:search"
            />
          </div>
        </div>
        <div class="h-96 overflow-y-auto space-y-4 px-3 pr-0">
          <div
            v-for="category in filteredCategories"
            :key="category.category"
            class="space-y-2"
          >
            <h5 class="text-xs font-medium text-[#6B7271] uppercase tracking-wide">
              {{ category.category }}
            </h5>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="icon in category.icons"
                :key="icon"
                class="size-10 bg-[#F5F5F5] hover:bg-[#ECF4E9] text-[#1E4841] rounded-lg flex items-center justify-center text-lg transition-colors"
                type="button"
                :class="{ 'bg-[#ECF4E9] ring-2 ring-[#1E4841]': modelValue === icon }"
                @click="selectIcon(icon)"
              >
                <Icon :name="icon" />
              </button>
            </div>
          </div>
          <div
            v-if="filteredCategories.length === 0"
            class="py-8 text-center text-sm text-[#6B7271]"
          >
            No icons found
          </div>
        </div>
      </div>
    </PopoverContent>
  </Popover>
</template>
