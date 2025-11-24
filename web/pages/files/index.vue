<script setup lang="ts">
import { Download, File, Folder, MoreHorizontal, Trash2, Upload } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
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

const files = ref([
  { name: 'videos', size: 0, lastModified: '', type: 'directory', path: 'videos/', items: 24 },
  { name: 'images', size: 0, lastModified: '', type: 'directory', path: 'images/', items: 156 },
  { name: 'documents', size: 0, lastModified: '', type: 'directory', path: 'documents/', items: 8 },
  { name: 'sample.mp4', size: 1024 * 1024 * 50, lastModified: '2023-11-20T10:00:00Z', type: 'file', path: 'sample.mp4' },
  { name: 'logo.png', size: 1024 * 50, lastModified: '2023-11-21T11:00:00Z', type: 'file', path: 'logo.png' },
  { name: 'presentation.pdf', size: 1024 * 1024 * 2, lastModified: '2023-11-19T14:30:00Z', type: 'file', path: 'presentation.pdf' },
  { name: 'audio.mp3', size: 1024 * 1024 * 5, lastModified: '2023-11-18T09:15:00Z', type: 'file', path: 'audio.mp3' },
])
const currentPath = ref('/')

function formatSize(bytes) {
  if (bytes === 0) {
    return '-'
  }
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${Number.parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`
}
</script>

<template>
  <div class="flex-1 space-y-4 p-8 pt-6 container">
    <div class="flex items-center justify-between space-y-2">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">
          Files
        </h2>
        <p class="text-muted-foreground">
          Browse and manage your media files
        </p>
      </div>
      <div class="flex items-center space-x-2">
        <Button>
          <Upload class="mr-2 h-4 w-4" />
          Upload
        </Button>
      </div>
    </div>

    <Card>
      <CardHeader>
        <div class="flex items-center gap-2">
          <div class="flex-1 flex items-center gap-1 rounded-md border bg-muted/50 px-3 py-2">
            <Folder class="h-4 w-4 text-muted-foreground" />
            <input
              v-model="currentPath"
              class="flex-1 bg-transparent text-sm outline-none"
              placeholder="/path/to/folder"
              type="text"
            />
          </div>
          <Button size="sm" variant="outline">
            Go
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead class="w-[50px]" />
              <TableHead>Name</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Size</TableHead>
              <TableHead>Last Modified</TableHead>
              <TableHead class="text-right">
                Actions
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="file in files" :key="file.name">
              <TableCell>
                <Folder v-if="file.type === 'directory'" class="h-5 w-5 text-primary" />
                <File v-else class="h-5 w-5 text-muted-foreground" />
              </TableCell>
              <TableCell class="font-medium">
                {{ file.name }}
              </TableCell>
              <TableCell>
                <Badge v-if="file.type === 'directory'" variant="secondary">
                  {{ file.items }} items
                </Badge>
                <span v-else class="text-sm text-muted-foreground">
                  {{ file.name.split('.').pop()?.toUpperCase() }}
                </span>
              </TableCell>
              <TableCell>{{ formatSize(file.size) }}</TableCell>
              <TableCell>{{ file.lastModified ? new Date(file.lastModified).toLocaleString() : '-' }}</TableCell>
              <TableCell class="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button class="h-8 w-8 p-0" variant="ghost">
                      <span class="sr-only">Open menu</span>
                      <MoreHorizontal class="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem v-if="file.type === 'file'">
                      <Download class="mr-2 h-4 w-4" />
                      Download
                    </DropdownMenuItem>
                    <DropdownMenuItem v-if="file.type === 'directory'">
                      <Folder class="mr-2 h-4 w-4" />
                      Open folder
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem class="text-destructive">
                      <Trash2 class="mr-2 h-4 w-4" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
