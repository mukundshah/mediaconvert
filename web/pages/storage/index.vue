<script setup lang="ts">
import {
  Copy,
  Eye,
  EyeOff,
  FolderOpen,
  Key,
  MoreHorizontal,
  Plus,
  Trash2,
} from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'

const buckets = ref([
  {
    id: 'bucket-1',
    name: 'media-input',
    provider: 'AWS S3',
    region: 'us-east-1',
    status: 'active',
    files: 1247,
    size: '42.8 GB',
    createdAt: '2024-01-15',
  },
  {
    id: 'bucket-2',
    name: 'media-output',
    provider: 'AWS S3',
    region: 'us-east-1',
    status: 'active',
    files: 892,
    size: '28.3 GB',
    createdAt: '2024-01-15',
  },
  {
    id: 'bucket-3',
    name: 'media-archive',
    provider: 'Google Cloud Storage',
    region: 'us-central1',
    status: 'active',
    files: 3421,
    size: '156.2 GB',
    createdAt: '2024-02-01',
  },
])

const credentials = ref([
  {
    id: 'cred-1',
    name: 'AWS Production',
    provider: 'AWS',
    accessKeyId: 'AKIA••••••••••••XAMPLE',
    status: 'active',
    lastUsed: '2 hours ago',
    createdAt: '2024-01-10',
  },
  {
    id: 'cred-2',
    name: 'GCP Production',
    provider: 'Google Cloud',
    accessKeyId: 'gcp-prod-••••••••',
    status: 'active',
    lastUsed: '5 hours ago',
    createdAt: '2024-01-12',
  },
  {
    id: 'cred-3',
    name: 'AWS Development',
    provider: 'AWS',
    accessKeyId: 'AKIA••••••••••••DEV01',
    status: 'inactive',
    lastUsed: '3 days ago',
    createdAt: '2024-01-08',
  },
])

const showSecret = ref({})

function toggleSecret(id) {
  showSecret.value[id] = !showSecret.value[id]
}

function copyToClipboard(text) {
  navigator.clipboard.writeText(text)
}
</script>

<template>
  <div class="flex-1 space-y-4 p-8 pt-6 container">
    <div class="flex items-center justify-between space-y-2">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">
          Storage & Credentials
        </h2>
        <p class="text-muted-foreground">
          Manage your storage buckets and access credentials
        </p>
      </div>
    </div>

    <Tabs class="space-y-4" default-value="buckets">
      <TabsList>
        <TabsTrigger value="buckets">
          <FolderOpen class="mr-2 h-4 w-4" />
          Buckets
        </TabsTrigger>
        <TabsTrigger value="credentials">
          <Key class="mr-2 h-4 w-4" />
          Credentials
        </TabsTrigger>
      </TabsList>

      <!-- Buckets Tab -->
      <TabsContent class="space-y-4" value="buckets">
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <div>
                <CardTitle>Storage Buckets</CardTitle>
                <CardDescription>
                  Configure and manage your storage buckets
                </CardDescription>
              </div>
              <Button>
                <Plus class="mr-2 h-4 w-4" />
                Add Bucket
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Provider</TableHead>
                  <TableHead>Region</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Files</TableHead>
                  <TableHead>Size</TableHead>
                  <TableHead class="text-right">
                    Actions
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="bucket in buckets" :key="bucket.id">
                  <TableCell class="font-medium">
                    {{ bucket.name }}
                  </TableCell>
                  <TableCell>{{ bucket.provider }}</TableCell>
                  <TableCell>{{ bucket.region }}</TableCell>
                  <TableCell>
                    <Badge
                      class="capitalize"
                      :variant="bucket.status === 'active' ? 'default' : 'secondary'"
                    >
                      {{ bucket.status }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ bucket.files.toLocaleString() }}</TableCell>
                  <TableCell>{{ bucket.size }}</TableCell>
                  <TableCell class="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger as-child>
                        <Button class="h-8 w-8 p-0" variant="ghost">
                          <span class="sr-only">Open menu</span>
                          <MoreHorizontal class="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem>View details</DropdownMenuItem>
                        <DropdownMenuItem>Edit configuration</DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem class="text-destructive">
                          <Trash2 class="mr-2 h-4 w-4" />
                          Delete bucket
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </TabsContent>

      <!-- Credentials Tab -->
      <TabsContent class="space-y-4" value="credentials">
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <div>
                <CardTitle>Access Credentials</CardTitle>
                <CardDescription>
                  Manage API keys and access credentials for cloud providers
                </CardDescription>
              </div>
              <Button>
                <Plus class="mr-2 h-4 w-4" />
                Add Credential
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Provider</TableHead>
                  <TableHead>Access Key</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Last Used</TableHead>
                  <TableHead class="text-right">
                    Actions
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="cred in credentials" :key="cred.id">
                  <TableCell class="font-medium">
                    {{ cred.name }}
                  </TableCell>
                  <TableCell>{{ cred.provider }}</TableCell>
                  <TableCell>
                    <div class="flex items-center gap-2">
                      <code class="text-xs">
                        {{ showSecret[cred.id] ? cred.accessKeyId.replace(/•/g, 'X') : cred.accessKeyId }}
                      </code>
                      <Button
                        class="h-6 w-6 p-0"
                        size="sm"
                        variant="ghost"
                        @click="toggleSecret(cred.id)"
                      >
                        <component :is="showSecret[cred.id] ? EyeOff : Eye" class="h-3 w-3" />
                      </Button>
                      <Button
                        class="h-6 w-6 p-0"
                        size="sm"
                        variant="ghost"
                        @click="copyToClipboard(cred.accessKeyId)"
                      >
                        <Copy class="h-3 w-3" />
                      </Button>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge
                      class="capitalize"
                      :variant="cred.status === 'active' ? 'default' : 'secondary'"
                    >
                      {{ cred.status }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ cred.lastUsed }}</TableCell>
                  <TableCell class="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger as-child>
                        <Button class="h-8 w-8 p-0" variant="ghost">
                          <span class="sr-only">Open menu</span>
                          <MoreHorizontal class="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem>View details</DropdownMenuItem>
                        <DropdownMenuItem>Rotate key</DropdownMenuItem>
                        <DropdownMenuItem>
                          {{ cred.status === 'active' ? 'Deactivate' : 'Activate' }}
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem class="text-destructive">
                          <Trash2 class="mr-2 h-4 w-4" />
                          Delete credential
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
