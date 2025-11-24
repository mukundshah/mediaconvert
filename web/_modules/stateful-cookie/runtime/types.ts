import type { CookieOptions } from 'nuxt/app'

export type StatefulCookieOptions<T> = CookieOptions<T> & { readonly?: boolean, encrypted?: boolean }
