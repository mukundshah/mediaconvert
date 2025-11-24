import type { StorageInterface } from '../allauth/types'

interface SessionState {
  session_token: string
  is_authenticated: boolean
  is_onboarded: boolean
}

interface AuthTokenStore extends StorageInterface {
  getAuthenticationStatus: () => boolean
  setAuthenticationStatus: (value: boolean) => void
  getOnboardingStatus: () => boolean
  setOnboardingStatus: (value: boolean) => void
  isAuthenticated: ComputedRef<boolean>
  isOnboarded: ComputedRef<boolean>
}

const storeInstances = new Map<string, AuthTokenStore>()

export const useAuthTokenStore = ({ namespace = 'default' }: { namespace?: string } = {}): AuthTokenStore => {
  const prefix = namespace !== 'default' ? `${namespace}:` : ''

  const state = useStatefulCookie<Partial<SessionState> | null>(`${prefix}sid`, {
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 30, // 30 days
    secure: !import.meta.dev,
    default: () => null,
    encode: (value) => {
      if (import.meta.dev) {
        return encodeURIComponent(JSON.stringify(value))
      }
      return btoa(encodeURIComponent(JSON.stringify(value)))
    },
    decode: (value) => {
      if (!value) return null
      if (import.meta.dev) {
        return JSON.parse(decodeURIComponent(value))
      }
      return JSON.parse(decodeURIComponent(atob(value)))
    },
  })

  return {
    getSessionToken: () => state.value?.session_token ?? null,
    setSessionToken: (value: string | null) => {
      if (!value) {
        state.value = null
      } else {
        state.value = {
          ...(state.value ?? {}),
          session_token: value,
        }
      }
    },
    isAuthenticated: computed(() => state.value?.is_authenticated ?? false),
    getAuthenticationStatus: () => state.value?.is_authenticated ?? false,
    setAuthenticationStatus: (value: boolean) => {
      state.value = {
        ...state.value,
        is_authenticated: value,
      }
    },
    isOnboarded: computed(() => state.value?.is_onboarded ?? false),
    getOnboardingStatus: () => state.value?.is_onboarded ?? false,
    setOnboardingStatus: (value: boolean) => {
      state.value = {
        ...state.value,
        is_onboarded: value,
      }
    },
  }
}

export const getStorage = ({ namespace = 'default' }: { namespace?: string } = {}) => {
  if (!storeInstances.has(namespace)) {
    const store = useAuthTokenStore({ namespace })
    storeInstances.set(namespace, store)
    return store
  }

  return storeInstances.get(namespace)!
}
