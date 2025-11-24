import type { RouterContext } from 'rou3'
import type { RouteLocationRaw } from 'vue-router'

import { defu } from 'defu'
import { addImports, addPlugin, addTemplate, createResolver, defineNuxtModule } from 'nuxt/kit'
import { createRouter as _createRouter, addRoute, findAllRoutes } from 'rou3'

interface ModuleOptions {
  loginRoute?: string
  fullAccessRoles?: string[]
  redirectIfNotAllowed?: string | false
  onboarding?: { enabled?: false } | { enabled: true, route: RouteLocationRaw }
}

const createRouter = <T extends Record<string, any> = Record<string, string>>(routes: string[] | Record<string, T>): RouterContext<T> => {
  const router = _createRouter<T>()
  if (Array.isArray(routes)) {
    for (const route of routes) {
      addRoute(router, 'GET', route, { path: route } as unknown as T)
    }
  } else {
    for (const [route, data] of Object.entries(routes)) {
      addRoute(router, 'GET', route, data)
    }
  }
  return router
}

export default defineNuxtModule<ModuleOptions>({
  meta: {
    name: 'auth',
    configKey: 'auth',
  },
  defaults: {
    fullAccessRoles: [],
    loginRoute: '/auth/login',
    redirectIfNotAllowed: false,
    onboarding: { enabled: false },
  },
  async setup(options, nuxt) {
    nuxt.options.runtimeConfig.public.auth = defu(nuxt.options.runtimeConfig.public.auth, {
      google: { enabled: false },
      magicCode: { enabled: false },
    })

    const { resolve } = createResolver(import.meta.url)

    const router = createRouter(nuxt.options.routeRules ?? {})
    const getRules = (url: string) => {
      const _rules = defu(
        {},
        ...findAllRoutes(router, 'GET', url)
          .map(route => route.data)
          .reverse(),
      ) as Record<string, any>

      if (_rules.auth) {
        _rules.auth.loginRoute ??= options.loginRoute
        _rules.auth.redirectIfNotAllowed ??= options.redirectIfNotAllowed

        if (options.fullAccessRoles?.length) {
          _rules.auth.roles = (_rules.auth.roles || []).concat(options.fullAccessRoles)
        }
      }

      return _rules
    }

    await addImports([
      { name: 'useAuth', from: resolve('./runtime/composable') },
      { name: 'useAuthConfiguration', from: resolve('./runtime/composable') },
      { name: 'useSocialAuth', from: resolve('./runtime/composable') },
      { name: 'useMFA', from: resolve('./runtime/composable') },
      { name: 'useEmailManagement', from: resolve('./runtime/composable') },
      { name: 'usePasswordManagement', from: resolve('./runtime/composable') },
      { name: 'usePermissions', from: resolve('./runtime/composable') },
    ])

    await addTemplate({
      filename: 'auth.config.mjs',
      getContents: () => `export default ${JSON.stringify(options, null, 2)}`,
    })

    nuxt.hook('modules:done', () => {
      addPlugin({ src: resolve('./runtime/plugin') }, { append: true })
    })

    nuxt.hook('prepare:types', ({ references, nodeReferences }) => {
      nodeReferences.push({
        path: resolve('./types.node.d.ts'),
      })
      references.push({
        path: resolve('./types.d.ts'),
      })
    })

    nuxt.hook('pages:resolved', (pages) => {
      pages.forEach((page) => {
        const rules = getRules(page.path)
        if (rules?.auth) {
          page.meta ||= {}
          page.meta.auth ??= rules.auth
        }
      })
    })
  },
})
