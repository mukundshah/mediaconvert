<script setup lang=ts>
import { SIDEBAR_FOOTER_NAVIGATION, SIDEBAR_NAVIIGATION } from '@/constants/navigation'

const { logout } = useAuth()

const handleLogout = async () => {
  await logout()
  await navigateTo({ name: 'auth-login' })
}

const isSearchFocused = ref(false)

const breadcrumbItems = useBreadcrumbItems({ hideRoot: true, overrides: [undefined, undefined, false, false] })

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.ctrlKey && e.key === 'k' && !isSearchFocused.value) {
    e.preventDefault()
    const searchInput = document.getElementById('search-input')
    if (searchInput) {
      searchInput.focus()
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})

useHead({ titleTemplate: '%siteName %separator %s' })
</script>

<template>
  <div class="overflow-hidden">
    <SidebarProvider
      :style=" {
        '--sidebar-width': 'calc(var(--spacing) * 72)',
        '--header-height': 'calc(var(--spacing) * 12)',
      }"
    >
      <Sidebar collapsible="offcanvas">
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                as-child
                class="data-[slot=sidebar-menu-button]:p-1.5!"
              >
                <NuxtLink to="/">
                  <AppIcon />
                  <span class="text-base font-semibold">4Fin</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupContent class="flex flex-col gap-2">
              <SidebarMenu>
                <SidebarMenuItem v-for="item in SIDEBAR_NAVIIGATION" :key="item.title">
                  <SidebarMenuButton
                    as-child
                    :is-active="$route.name?.toString().startsWith(item.to.name)"
                    :tooltip="item.title"
                  >
                    <NuxtLink :to="item.to">
                      <Icon :name="item.icon" />
                      <span>{{ item.title }}</span>
                    </NuxtLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
          <SidebarGroup class="mt-auto">
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem
                  v-for="item in SIDEBAR_FOOTER_NAVIGATION"
                  :key="item.title"
                >
                  <SidebarMenuButton as-child>
                    <NuxtLink :to="item.to">
                      <Icon :name="item.icon" />
                      {{ item.title }}
                    </NuxtLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <!-- <NavUser :user="data.user" /> -->
        </SidebarFooter>
      </Sidebar>

      <SidebarInset>
        <header class="flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)">
          <div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
            <SidebarTrigger class="-ml-1" />
            <Separator
              class="mx-2 data-[orientation=vertical]:h-4"
              orientation="vertical"
            />
            <Breadcrumb>
              <BreadcrumbList>
                <template v-for="(item, index) in breadcrumbItems" :key="index">
                  <BreadcrumbItem>
                    <BreadcrumbPage v-if="item.current" as-child>
                      <h1 class="text-base font-medium">
                        {{ item.label }}
                      </h1>
                    </BreadcrumbPage>
                    <BreadcrumbLink v-else as-child>
                      <NuxtLink :to="item.href">
                        {{ item.label }}
                      </NuxtLink>
                    </BreadcrumbLink>
                  </BreadcrumbItem>
                  <BreadcrumbSeparator v-if="index !== breadcrumbItems.length - 1">
                    <Icon name="lucide:slash" />
                  </BreadcrumbSeparator>
                </template>
              </BreadcrumbList>
            </Breadcrumb>
            <div class="ml-auto flex items-center gap-2">
              <Button class="hidden sm:flex" size="sm">
                <Icon name="lucide:circle-plus" />
                <span>Quick Create</span>
              </Button>
            </div>
          </div>
        </header>
        <div class="flex flex-1 flex-col gap-4 p-4 pt-0">
          <slot></slot>
        </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>
