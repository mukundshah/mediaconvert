import type {
  ConfigurationResponse,
  ConfirmLoginCodeRequest,
  EmailAddress,
  EmailAddressRequest,
  EmailPrimaryRequest,
  EmailVerificationRequest,
  LoginByCodeRequest,
  LoginRequest,
  PasswordChangeRequest,
  PasswordResetConfirmRequest,
  PasswordResetRequest,
  ProviderDisconnectRequest,
  ProviderRedirectRequest,
  ReauthenticateRequest,
  SignupRequest,
} from '../allauth/types'

// @ts-expect-error vfs
import authconfig from '#build/auth.config'

import { getClient } from '../allauth/client'

import { getStorage } from './storage'

export const useAuthConfiguration = ({ namespace }: { namespace?: string } = {}) => {
  const config = useRuntimeConfig()
  const storage = getStorage({ namespace })
  const client = getClient({ storage, apiBaseURL: config.public.api.baseURL })

  const getConfiguration = async (): Promise<ConfigurationResponse['data']> => {
    const { data } = await client.getConfiguration()
    return data
  }

  return {
    getConfiguration,
  }
}

export const useAuth = ({ namespace }: { namespace?: string } = {}) => {
  const config = useRuntimeConfig()

  const storage = getStorage({ namespace })
  const client = getClient({ storage, apiBaseURL: config.public.api.baseURL })

  const login = async (data: LoginRequest) => {
    const { data: response, meta: { is_authenticated } } = await client.login(data)

    storage.setAuthenticationStatus(is_authenticated)
    if ('user' in response) {
      storage.setOnboardingStatus(response.user.is_onboarded)
    }

    return response
  }

  const requestLoginCode = async (data: LoginByCodeRequest) => {
    const { data: response } = await client.requestLoginCode(data)
    return response
  }

  const confirmLoginCode = async (data: ConfirmLoginCodeRequest) => {
    const { data: response, meta: { is_authenticated } } = await client.confirmLoginCode(data)

    storage.setAuthenticationStatus(is_authenticated)
    if ('user' in response) {
      storage.setOnboardingStatus(response.user.is_onboarded)
    }

    return response
  }

  const logout = async () => {
    const { data: response } = await client.logout()

    storage.setAuthenticationStatus(false)
    storage.setOnboardingStatus(false)

    return response
  }

  const signup = async (data: SignupRequest) => {
    const { data: response, meta: { is_authenticated } } = await client.signup(data)

    storage.setAuthenticationStatus(is_authenticated)
    if ('user' in response) {
      storage.setOnboardingStatus(response.user.is_onboarded)
    }

    return response
  }

  const reauthenticate = async (data: ReauthenticateRequest) => {
    const { data: response, meta: { is_authenticated } } = await client.reauthenticate(data)

    storage.setAuthenticationStatus(is_authenticated)
    if ('user' in response) {
      storage.setOnboardingStatus(response.user.is_onboarded)
    }

    return response
  }

  const syncAuthenticationStatus = async () => {
    const { data: response, meta: { is_authenticated } } = await client.getAuthenticationStatus()

    storage.setAuthenticationStatus(is_authenticated)
    if ('user' in response) {
      storage.setOnboardingStatus(response.user.is_onboarded)
    }

    return response
  }

  return {
    get token() {
      return storage.getSessionToken()
    },
    isAuthenticated: computed(() => storage.isAuthenticated.value),
    isOnboarded: computed(() => storage.isOnboarded.value || !authconfig.onboarding?.enabled),
    login,
    requestLoginCode,
    confirmLoginCode,
    logout,
    signup,
    reauthenticate,
    syncAuthenticationStatus,
  }
}

export const useSocialAuth = ({ namespace }: { namespace?: string } = {}) => {
  const storage = getStorage({ namespace })
  const client = getClient({ storage })
  const url = useRequestURL()

  const getProviderAccounts = async () => {
    const { data: response } = await client.listProviderAccounts()
    return response
  }

  const connect = async (data: Omit<ProviderRedirectRequest, 'callback_url'> & { callback_url?: string }) => {
    const callback_url = data.callback_url ?? `${url.origin}/auth/callback/${data.provider}`
    await client.providerRedirect({ callback_url, ...data })
  }

  const callback = async (provider: string, data: Record<string, any>) => {
    const { data: response, meta: { is_authenticated } } = await client.providerCallback(provider, data)

    storage.setAuthenticationStatus(is_authenticated)
    if ('user' in response) {
      storage.setOnboardingStatus(response.user.is_onboarded)
    }

    return response
  }

  const disconnect = async (data: ProviderDisconnectRequest) => {
    const { data: response } = await client.disconnectProviderAccount(data)
    return response
  }

  return {
    connect,
    callback,
    disconnect,
    getProviderAccounts,
  }
}

export const useMFA = ({ namespace }: { namespace?: string } = {}) => {
  const storage = getStorage({ namespace })
  const _client = getClient({ storage })

  // TODO: Implement MFA

  return { }
}

export const useEmailManagement = ({ namespace }: { namespace?: string } = {}) => {
  const storage = getStorage({ namespace })
  const client = getClient({ storage })

  const getEmailAddresses = async (): Promise<EmailAddress[]> => {
    const { data: response } = await client.listEmailAddresses()
    return response
  }

  const addEmail = async (data: EmailAddressRequest) => {
    const { data: response } = await client.addEmailAddress(data)
    return response
  }

  const checkEmail = async (data: EmailAddressRequest) => {
    const { data: response } = await client.checkEmail(data)
    return response
  }

  const deleteEmail = async (data: EmailAddressRequest) => {
    const { data: response } = await client.removeEmailAddress(data)
    return response
  }

  const markEmailAsPrimary = async (data: Omit<EmailPrimaryRequest, 'primary'>) => {
    const { data: response } = await client.changePrimaryEmailAddress({ ...data, primary: true })
    return response
  }

  const requestEmailVerification = async (data: EmailAddressRequest) => {
    const response = await client.requestEmailVerification(data)
    return ('data' in response) ? response.data : { success: true }
  }

  const verifyEmail = async (data: EmailVerificationRequest) => {
    const response = await client.verifyEmail(data)
    return ('data' in response) ? response.data : { success: true }
  }

  return {
    addEmail,
    checkEmail,
    deleteEmail,
    markEmailAsPrimary,
    getEmailAddresses,
    requestEmailVerification,
    verifyEmail,
  }
}

// export const usePhoneManagement = ({ namespace }: { namespace?: string } = {}) => {
//   const storage = getStorage({ namespace })
//   const client = getClient({ storage })

//   const getPhoneNumbers = async (): Promise<PhoneNumber[]> => {
//     const { data: response } = await client.listPhoneNumbers()
//     return response
//   }
// }

export const usePasswordManagement = ({ namespace }: { namespace?: string } = {}) => {
  const storage = getStorage({ namespace })
  const client = getClient({ storage })

  const requestPasswordReset = async (data: PasswordResetRequest) => {
    const response = await client.requestPassword(data)
    return ('data' in response) ? response.data : { success: true }
  }

  const getPasswordResetInfo = async (key: string) => {
    const { data: response } = await client.getPasswordResetInfo(key)
    return response
  }

  const resetPassword = async (data: PasswordResetConfirmRequest) => {
    const { data: response } = await client.resetPassword(data)
    return response
  }

  const changePassword = async (data: PasswordChangeRequest) => {
    const response = await client.changePassword(data)
    return ('data' in response) ? response.data : { success: true }
  }

  return {
    requestPasswordReset,
    resetPassword,
    getPasswordResetInfo,
    changePassword,
  }
}

export const usePermissions = ({ namespace }: { namespace?: string } = {}) => {
  // const config = useRuntimeConfig()
  // const storage = getStorage({ namespace })
  // const client = getClient({ storage, apiBaseURL: config.public.api.baseURL })

  const config = useRuntimeConfig()
  const { token } = useAuth({ namespace })

  const prefix = namespace !== 'default' ? `${namespace}:` : ''

  const _roles = useState<string[]>(`${prefix}roles`, () => [])
  const _permissions = useState<Record<string, Record<string, boolean>>>(`${prefix}permissions`, () => ({}))

  const _loadPermissions = async () => {
    const response = await $fetch('/api/auth/permissions', {
      baseURL: config.public.api.baseURL,
      onRequest: (ctx) => {
        if (token) {
          ctx.options.headers.set('X-Session-Token', token)
        }
      },
    })
    return response
  }

  // const getPermissions = async () => {
  //   const { data: response } = await client.listPermissions()
  //   return response
  // }

  const hasRole = (role: string) => {
    if (!_roles.value?.includes(role)) return false
    return _roles.value.includes(role)
  }

  const hasAnyRole = (roles: string[]) => {
    if (!_roles.value?.length) return false
    return roles.some(role => _roles.value!.includes(role))
  }

  const hasAllRoles = (roles: string[]) => {
    if (!_roles.value.length) return false
    return roles.every(role => _roles.value.includes(role))
  }

  const _hasFullAccess = () => {
    return true
    // return hasAnyRole(authconfig.fullAccessRoles)
  }

  type Permission = `${string}.${string}` | `${string}`

  const hasPermission = (permission: Permission) => {
    // Check for full access roles first
    if (_hasFullAccess()) return true

    if (Object.keys(_permissions.value || {}).length === 0) return false

    if (permission.includes('.')) {
      const [resource, action] = permission.split('.')
      if (!resource || !action) return false
      return _permissions.value?.[resource]?.[action] ?? false
    }

    // Resource-only check - return true if any action is allowed for this resource
    return Object.keys(_permissions.value?.[permission] || {}).some(action => _permissions.value?.[permission]?.[action] === true)
  }

  const hasAnyPermission = (permissions: Permission[]) => {
    // Check for full access roles first
    if (_hasFullAccess()) return true

    if (!_permissions.value) return false
    return permissions.some(permission => hasPermission(permission))
  }

  const hasAllPermissions = (permissions: Permission[]) => {
    // Check for full access roles first
    if (_hasFullAccess()) return true

    if (!_permissions.value) return false
    return permissions.every(permission => hasPermission(permission))
  }

  return {
    hasRole,
    hasAnyRole,
    hasAllRoles,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
  }
}
