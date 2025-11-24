<script setup lang="ts">
import type { UserData } from '@/types'
import { useSidebar } from '@/components/ui/sidebar'

const { isMobile } = useSidebar()
const { logout } = useAuth()
const colorMode = useColorMode()

const handleLogout = async () => {
  await logout()
  await navigateTo({ name: 'auth-login' })
}

// Fetch user data from the API
const { data: user, pending } = await useAPI<UserData>('/v1/users/me/')
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem v-if="pending">
      <SidebarMenuButton
        class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
        size="lg"
      >
        <div>
          <Skeleton class="h-8 w-8 rounded-lg" />
        </div>
        <div>
          <Skeleton class="h-3 w-24" />
          <Skeleton class="h-3 w-32" />
        </div>
      </SidebarMenuButton>
    </SidebarMenuItem>

    <SidebarMenuItem v-if="user">
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            size="lg"
          >
            <Avatar class="h-8 w-8 rounded-lg grayscale">
              <AvatarImage :alt="`${user.first_name} ${user.last_name}`" :src="user.avatar_url" />
              <AvatarFallback class="rounded-lg">
                {{ user.first_name.charAt(0) + user.last_name.charAt(0) }}
              </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-medium">{{ user.first_name }}</span>
              <span class="text-muted-foreground truncate text-xs">
                {{ user.email }}
              </span>
            </div>
            <IconDotsVertical class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          align="end"
          class="w-(--reka-dropdown-menu-trigger-width) min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          :side-offset="4"
        >
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-8 w-8 rounded-lg">
                <AvatarImage :alt="`${user.first_name} ${user.last_name}`" :src="user.avatar_url" />
                <AvatarFallback class="rounded-lg">
                  {{ user.first_name.charAt(0) + user.last_name.charAt(0) }}
                </AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-medium">{{ `${user.first_name} ${user.last_name}` }}</span>
                <span class="text-muted-foreground truncate text-xs">
                  {{ user.email }}
                </span>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>
              <Icon name="lucide:moon" />
              Theme
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent>
                <DropdownMenuItem @click="colorMode.preference = 'light'">
                  Light
                </DropdownMenuItem>
                <DropdownMenuItem @click="colorMode.preference = 'dark'">
                  Dark
                </DropdownMenuItem>
                <DropdownMenuItem @click="colorMode.preference = 'system'">
                  System
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="handleLogout">
            <Icon name="lucide:log-out" />
            Log out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
