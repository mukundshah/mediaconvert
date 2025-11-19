import tailwindcss from '@tailwindcss/vite'
import { createResolver } from 'nuxt/kit'

const { resolve } = createResolver(import.meta.url)

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },

  modules: [
    '@nuxt/eslint',
    '@nuxt/fonts',
    '@nuxt/icon',
    '@nuxtjs/seo',
    '@nuxtjs/color-mode',
    'shadcn-nuxt',
    'nuxt-gtag',
  ],

  srcDir: 'web',

  dir: {
    modules: resolve('./web/modules'),
  },

  app: {
    rootAttrs: {
      id: 'app',
    },
    head: {
      templateParams: {
        separator: '-',
      },
    },
  },

  site: {
    name: 'Media Convert',
    url: 'https://mediaconvert.app',
  },

  routeRules: {
    '/**': { appLayout: 'app', robots: false, auth: { required: true } },
    '/auth/login': { robots: true },
    '/auth/**': { appLayout: 'auth', robots: false, auth: { required: false, redirectIfLoggedIn: true } },
    '/': { appLayout: 'site', robots: true, auth: { required: false } },
    '/terms': { appLayout: 'site', robots: true, auth: { required: false } },
    '/privacy': { appLayout: 'site', robots: true, auth: { required: false } },
    '/contact': { appLayout: 'site', robots: true, auth: { required: false } },
    '/about': { appLayout: 'site', robots: true, auth: { required: false } },
  },

  css: ['assets/css/tailwind.css'],

  api: {
    protected: true,
    trailingSlash: true,
    authorizationHeader: 'Authorization',
    tokenPrefix: 'Bearer ',
  },

  colorMode: {
    classSuffix: '',
  },

  eslint: {
    config: {
      standalone: false,
    },
  },

  shadcn: {
    prefix: '',
  },

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  nitro: {
    handlers: [
      {
        route: '/health',
        handler: '~/health.ts',
      },
    ],
  },
})
