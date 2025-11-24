<script setup lang="ts">
import { ArrowLeft, Save } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const router = useRouter()

const pipeline = ref({
  name: '',
  description: '',
  format: 'yaml',
  content: `name: my-pipeline
steps:
  - operation: transcode
    input: \${input}
    output: \${output}/video.mp4
    params:
      codec: h264
      quality: 23`,
})

const save = () => {
  // Mock save
  console.log('Saving pipeline', pipeline.value)
  router.push('/pipelines')
}
</script>

<template>
  <div class="h-full flex-1 flex-col space-y-8 p-8 md:flex">
    <div class="flex items-center justify-between space-y-2">
      <div class="flex items-center space-x-2">
        <Button variant="ghost" size="icon" @click="router.back()">
          <ArrowLeft class="h-4 w-4" />
        </Button>
        <div>
          <h2 class="text-2xl font-bold tracking-tight">
            New Pipeline
          </h2>
          <p class="text-muted-foreground">
            Define a new processing workflow.
          </p>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <Button @click="save">
          <Save class="mr-2 h-4 w-4" />
          Save Pipeline
        </Button>
      </div>
    </div>

    <div class="grid gap-4 py-4">
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="name" class="text-right">
          Name
        </Label>
        <Input id="name" v-model="pipeline.name" class="col-span-3" placeholder="e.g. video-compress" />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="description" class="text-right">
          Description
        </Label>
        <Input id="description" v-model="pipeline.description" class="col-span-3" placeholder="Short description" />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="format" class="text-right">
          Format
        </Label>
        <Select v-model="pipeline.format">
          <SelectTrigger class="col-span-3">
            <SelectValue placeholder="Select format" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="yaml">
              YAML
            </SelectItem>
            <SelectItem value="json">
              JSON
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="grid grid-cols-4 gap-4">
        <Label for="content" class="text-right pt-2">
          Definition
        </Label>
        <Textarea id="content" v-model="pipeline.content" class="col-span-3 font-mono" rows="15" />
      </div>
    </div>
  </div>
</template>
