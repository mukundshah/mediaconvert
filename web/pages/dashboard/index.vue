<script setup lang="ts">
import {
  ArrowUpRight,
  CheckCircle2,
  Clock,
  MoreHorizontal,
  Play,
  XCircle,
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
  DropdownMenuLabel,
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

const jobs = ref([
  {
    id: 'job-123',
    pipeline: 'Video Compression',
    status: 'processing',
    startedAt: '2 mins ago',
    duration: '45s',
    progress: 45,
  },
  {
    id: 'job-124',
    pipeline: 'Image Resizing',
    status: 'completed',
    startedAt: '1 hour ago',
    duration: '1m 20s',
    progress: 100,
  },
  {
    id: 'job-125',
    pipeline: 'PDF Extraction',
    status: 'failed',
    startedAt: '2 hours ago',
    duration: '15s',
    progress: 10,
  },
  {
    id: 'job-126',
    pipeline: 'Video Thumbnail',
    status: 'pending',
    startedAt: 'Just now',
    duration: '-',
    progress: 0,
  },
  {
    id: 'job-127',
    pipeline: 'Audio Transcode',
    status: 'completed',
    startedAt: '3 hours ago',
    duration: '4m 12s',
    progress: 100,
  },
])

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Good morning'
  if (hour < 18) return 'Good afternoon'
  return 'Good evening'
})
</script>

<template>
  <div class="flex-1 space-y-4 p-8 pt-6 container">
    <div class="flex items-center justify-between space-y-2">
      <h2 class="text-3xl font-bold tracking-tight">
        {{ greeting }}, Mukund
      </h2>
      <div class="flex items-center space-x-2">
        <Button>
          <Play class="mr-2 h-4 w-4" />
          New Job
        </Button>
      </div>
    </div>

    <!-- Recent Jobs Table -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>Recent Jobs</CardTitle>
            <CardDescription>
              Latest processing tasks and their status
            </CardDescription>
          </div>
          <Button as-child size="sm" variant="outline">
            <NuxtLink to="/jobs">
              View All
              <ArrowUpRight class="ml-2 h-4 w-4" />
            </NuxtLink>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Job ID</TableHead>
              <TableHead>Pipeline</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Started</TableHead>
              <TableHead>Duration</TableHead>
              <TableHead class="text-right">
                Actions
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="job in jobs" :key="job.id">
              <TableCell class="font-medium">
                {{ job.id }}
              </TableCell>
              <TableCell>{{ job.pipeline }}</TableCell>
              <TableCell>
                <Badge
                  class="capitalize"
                  :variant="job.status === 'completed' ? 'default' : job.status === 'failed' ? 'destructive' : 'secondary'"
                >
                  <component
                    :is="job.status === 'completed' ? CheckCircle2 : job.status === 'failed' ? XCircle : Clock"
                    class="mr-1 h-3 w-3"
                  />
                  {{ job.status }}
                </Badge>
              </TableCell>
              <TableCell>{{ job.startedAt }}</TableCell>
              <TableCell>{{ job.duration }}</TableCell>
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
                    <DropdownMenuItem>View logs</DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem class="text-destructive">
                      Cancel job
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
