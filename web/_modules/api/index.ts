import type { RouteLocationRaw } from 'vue-router'

import { defu } from 'defu'
import { addImports, addPlugin, addTemplate, createResolver, defineNuxtModule } from 'nuxt/kit'

export interface ModuleOptions {
  baseURL: string
  trailingSlash: boolean
  protected: boolean
  unauthorized?:
    | {
      redirect: RouteLocationRaw
      statusCodes: number[]
      strategy: 'redirect'
    }
    | {
      statusCodes: number[]
      strategy: 'error'
    }
  authNamespace: string | 'default'
  tokenPrefix?: string
  authorizationHeader?: string
}

const defaults: ModuleOptions = {
  baseURL: '/',
  trailingSlash: false,
  protected: false,
  authNamespace: 'default',
  authorizationHeader: 'Authorization',
  tokenPrefix: 'Token ',
  unauthorized: {
    statusCodes: [401, 403],
    strategy: 'error',
  },
}

export default defineNuxtModule<ModuleOptions>({
  meta: {
    name: 'api',
    configKey: 'api',
  },
  defaults,
  async setup(options, nuxt) {
    nuxt.options.runtimeConfig.public.api = defu(nuxt.options.runtimeConfig.public.api, options)

    const { resolve } = createResolver(import.meta.url)

    await addPlugin({ src: resolve('./runtime/plugin') })
    await addImports([{ name: 'useAPI', from: resolve('./runtime/composable') }])

    await addTemplate({
      filename: 'api.config.mjs',
      getContents: () => `export default ${JSON.stringify(options, null, 2)}`,
    })

    nuxt.hook('prepare:types', ({ references }) => {
      references.push({
        path: resolve('./runtime/types.d.ts'),
      })
    })
  },
})
