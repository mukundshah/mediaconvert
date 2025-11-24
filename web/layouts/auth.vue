<script setup lang="ts">
const colorMode = useColorMode()

useHead({ titleTemplate: '%siteName %separator %s' })
</script>

<template>
  <div class="magicpattern">
    <header class="py-2 px-2 sm:px-4 flex items-center justify-end">
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button size="icon" variant="ghost">
            <Icon class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" name="lucide:moon" />
            <Icon class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" name="lucide:sun" />
            <span class="sr-only">Toggle theme</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem @click="colorMode.preference = 'light'">
            Light
          </DropdownMenuItem>
          <DropdownMenuItem @click="colorMode.preference = 'dark'">
            Dark
          </DropdownMenuItem>
          <DropdownMenuItem @click="colorMode.preference = 'system'">
            System
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </header>
    <slot></slot>
  </div>
</template>

<style>
.magicpattern {
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center center;
  background-repeat: repeat;
  background-image: url('assets/svgs/magicpattern.svg');
}

.dark .magicpattern {
  background-image: url('assets/svgs/magicpattern-dark.svg');
}
</style>
