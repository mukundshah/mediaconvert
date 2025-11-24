<script setup lang="ts">
import { Edit, MoreHorizontal, Play, Plus, Trash2 } from 'lucide-vue-next'
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

const pipelines = ref([
  {
    id: 'video-compress',
    name: 'Video Compression',
    description: 'Compresses video to H.264',
    steps: 2,
    status: 'active',
    lastRun: '2 hours ago',
  },
  {
    id: 'image-resize',
    name: 'Image Resizing',
    description: 'Resizes images to 800x600',
    steps: 1,
    status: 'active',
    lastRun: '5 hours ago',
  },
  {
    id: 'pdf-extract',
    name: 'PDF Text Extraction',
    description: 'Extracts text from PDF files',
    steps: 1,
    status: 'inactive',
    lastRun: '3 days ago',
  },
  {
    id: 'audio-transcode',
    name: 'Audio Transcoding',
    description: 'Converts audio to MP3 format',
    steps: 2,
    status: 'active',
    lastRun: '1 day ago',
  },
])
</script>

<template>
  <div class="flex-1 space-y-4 p-8 pt-6 container">
    <div class="flex items-center justify-between space-y-2">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">
          Pipelines
        </h2>
        <p class="text-muted-foreground">
          Manage your media processing workflows
        </p>
      </div>
      <div class="flex items-center space-x-2">
        <Button as-child>
          <NuxtLink to="/pipelines/new">
            <Plus class="mr-2 h-4 w-4" />
            New Pipeline
          </NuxtLink>
        </Button>
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>All Pipelines</CardTitle>
        <CardDescription>
          Configure and manage processing workflows
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Steps</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Last Run</TableHead>
              <TableHead class="text-right">
                Actions
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="pipeline in pipelines" :key="pipeline.id">
              <TableCell class="font-medium">
                {{ pipeline.name }}
              </TableCell>
              <TableCell>{{ pipeline.description }}</TableCell>
              <TableCell>{{ pipeline.steps }}</TableCell>
              <TableCell>
                <Badge
                  class="capitalize"
                  :variant="pipeline.status === 'active' ? 'default' : 'secondary'"
                >
                  {{ pipeline.status }}
                </Badge>
              </TableCell>
              <TableCell>{{ pipeline.lastRun }}</TableCell>
              <TableCell class="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button class="h-8 w-8 p-0" variant="ghost">
                      <span class="sr-only">Open menu</span>
                      <MoreHorizontal class="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem>
                      <Play class="mr-2 h-4 w-4" />
                      Run pipeline
                    </DropdownMenuItem>
                    <DropdownMenuItem>
                      <Edit class="mr-2 h-4 w-4" />
                      Edit
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
