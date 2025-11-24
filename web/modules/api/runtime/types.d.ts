import type { ModuleOptions as APIModuleOptions } from '../index'
import type { $API } from './types'

export type { $APIOptions, KeysOf, PickFrom, UseAPIOptions } from './types'

declare module '#app' {
  interface NuxtApp {
    $api: $API
  }
}

declare module 'nuxt/schema' {
  interface PublicRuntimeConfig {
    api: APIModuleOptions
  }
}

export {}
