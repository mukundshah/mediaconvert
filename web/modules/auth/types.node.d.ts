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

export {}
