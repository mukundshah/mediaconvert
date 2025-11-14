import antfu from '@antfu/eslint-config'

export default antfu({
  type: 'app',

  // Enable stylistic formatting rules
  stylistic: {
    indent: 2,
    quotes: 'single',
  },

  // TypeScript and Vue are auto-detected
  typescript: true,
  vue: true,

  // Disable jsonc, yaml, toml, markdown if not needed
  jsonc: false,
  yaml: false,

  // `.eslintignore` is no longer supported in Flat config, use `ignores` instead
  ignores: [
    '**/node_modules',
    '**/dist',
    '**/.nuxt',
    '**/.output',
    '**/coverage',
  ],
})
