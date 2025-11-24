import { addComponent, addImports, createResolver, defineNuxtModule } from 'nuxt/kit'

export default defineNuxtModule({
  meta: {
    name: 'formatted-number',
    configKey: 'formattedNumber',
  },
  async setup() {
    const { resolve } = createResolver(import.meta.url)

    addComponent({
      name: 'FormattedNumber',
      filePath: resolve('./runtime/components/FormattedNumber.vue'),
    })

    addImports({
      name: 'useNumberFormatter',
      from: resolve('./runtime/composables/useNumberFormatter.ts'),
    })
  },
})
