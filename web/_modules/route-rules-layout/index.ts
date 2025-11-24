import type { RouterContext } from 'rou3'

import { defu } from 'defu'
import { addTypeTemplate, defineNuxtModule } from 'nuxt/kit'
import { createRouter as _createRouter, addRoute, findAllRoutes } from 'rou3'

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

const genString = (name: string) => {
  return JSON.stringify(name)
}

export default defineNuxtModule({
  meta: {
    name: 'route-rules-layout',
  },
  async setup(options, nuxt) {
    const router = createRouter(nuxt.options.routeRules ?? {})

    const getRules = (url: string) => {
      return defu(
        {},
        ...findAllRoutes(router, 'GET', url)
          .map(route => route.data)
          .reverse(),
      ) as Record<string, any>
    }

    addTypeTemplate({
      filename: 'types/route-rules-layout.d.ts',
      getContents: ({ app }) => {
        return [
          `export type LayoutKey = ${Object.keys(app.layouts).map(name => genString(name)).join(' | ') || 'string'}`,
          'declare module \'nitropack/types\' {',
          '  interface NitroRouteConfig {',
          '    appLayout?: LayoutKey | false',
          '  }',
          '}',
          'declare module \'nitropack\' {',
          '  interface NitroRouteConfig {',
          '    appLayout?: LayoutKey | false',
          '  }',
          '}',
        ].join('\n')
      },
    }, { nuxt: true, nitro: true, node: true })

    nuxt.options.experimental.scanPageMeta = true
    nuxt.options.experimental.extraPageMetaExtractionKeys ||= []
    nuxt.options.experimental.extraPageMetaExtractionKeys.push('layout')

    nuxt.hook('pages:resolved', (pages) => {
      pages.forEach((page) => {
        const rules = getRules(page.path)
        if (rules?.appLayout !== undefined) {
          page.meta ||= {}
          page.meta.layout ??= rules.appLayout
        }
      })
    })
  },
})
