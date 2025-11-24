import { addImports, createResolver, defineNuxtModule } from 'nuxt/kit'

export default defineNuxtModule({
  meta: {
    name: 'stateful-cookie',
  },
  async setup(options, nuxt) {
    const { resolve } = createResolver(import.meta.url)

    await addImports([{ name: 'useStatefulCookie', from: resolve('./runtime/composable') }])

    nuxt.hook('prepare:types', ({ references }) => {
      references.push({
        path: resolve('./runtime/types.ts'),
      })
    })
  },
})
