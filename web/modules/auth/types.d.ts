import type { AuthMeta } from './runtime/types'

declare module 'nitropack/types' {
  interface NitroRouteConfig {
    auth?: AuthMeta
  }
  interface NitroRouteRules {
    auth?: AuthMeta
  }
}

declare module 'nitropack' {
  interface NitroRouteConfig {
    auth?: AuthMeta
  }
  interface NitroRouteRules {
    auth?: AuthMeta
  }
}

declare module '#app' {
  interface PageMeta {
    auth?: AuthMeta
  }
}

declare module 'vue-router' {
  interface RouteMeta {
    auth?: AuthMeta
  }
}

export {}
