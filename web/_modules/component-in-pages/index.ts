import { defineNuxtModule } from 'nuxt/kit'

export default defineNuxtModule({
  meta: {
    name: 'component-in-pages',
  },
  async setup(options, nuxt) {
    nuxt.hook('pages:extend', (pages) => {
      // Step 1: Filter out routes with @components in the file path
      const filtered = pages.filter(route => !route.file?.includes('/@components'))

      // Step 2: Identify layout routes and regular routes
      const layouts = [] as typeof pages
      const regularRoutes = [] as typeof pages

      filtered.forEach((route) => {
        if (route.file?.includes('/@layout')) {
          layouts.push(route)
        } else {
          regularRoutes.push(route)
        }
      })

      // Step 3: Transform layouts and build the hierarchy
      const result = [] as typeof pages
      const processedRoutes = new Set()

      layouts.forEach((layout) => {
        // Rewrite the path: remove /@layout
        const layoutPath = layout.path.replace('/@layout', '')
        const layoutDir = layout.file?.substring(0, layout.file?.lastIndexOf('/') ?? 0)

        // Find children: routes whose file path starts with the layout's directory
        const children = regularRoutes.filter((route) => {
          // Check if the route's file is within the layout's directory or subdirectories
          return route.file.startsWith(`${layoutDir}/`) && route.file !== layout.file
        })

        // Mark children as processed
        children.forEach(child => processedRoutes.add(child))

        // Add the layout with its children
        result.push({
          ...layout,
          path: layoutPath,
          children: children.map(child => ({ ...child })),
        })
      })

      // Step 4: Add remaining routes that weren't nested under a layout
      regularRoutes.forEach((route) => {
        if (!processedRoutes.has(route)) {
          result.push({ ...route })
        }
      })

      // Step 5: Clear the original array and add back only the filtered pages
      pages.length = 0
      pages.push(...result)
    })
  },
})
