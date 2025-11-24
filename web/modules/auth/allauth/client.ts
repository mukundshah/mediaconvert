import type {
  AuthenticatedResponse,
  AuthenticationResponse,
  AuthenticatorsResponse,
  ConfigurationResponse,
  ConfirmLoginCodeRequest,
  EmailAddressesResponse,
  EmailAddressRequest,
  EmailCheckResponse,
  EmailPrimaryRequest,
  EmailVerificationInfoResponse,
  EmailVerificationRequest,
  LoginByCodeRequest,
  LoginRequest,
  MFAAuthenticateRequest,
  MFATrustRequest,
  NotAuthenticatedResponse,
  NoTOTPAuthenticatorResponse,
  PasswordChangeRequest,
  PasswordResetConfirmRequest,
  PasswordResetInfoResponse,
  PasswordResetRequest,
  PhoneNumberResponse,
  PhoneVerificationRequest,
  ProviderAccountsResponse,
  ProviderDisconnectRequest,
  ProviderRedirectRequest,
  ProviderSignupRequest,
  ProviderSignupResponse,
  ProviderTokenRequest,
  ReauthenticateRequest,
  SensitiveRecoveryCodesAuthenticatorResponse,
  SessionsResponse,
  SignupRequest,
  StorageInterface,
  TOTPActivateRequest,
  TOTPAuthenticatorResponse,
  WebAuthnCredentialCreationOptions,
  WebAuthnCredentialRequestOptions,
  WebAuthnLoginRequest,
  WebAuthnSignupRequest,
} from './types'

import { joinURL, withQuery } from 'ufo'

// Singleton instance
let clientInstance: AllAuthClient | null = null

/**
 * AllauthClient provides methods to interact with the django-allauth headless API.
 * It supports both browser and app clients with automatic token management.
 */
export class AllAuthClient {
  private storage: StorageInterface
  private apiBaseURL: string

  constructor(
    {
      apiBaseURL = '/',
      storage,
    }: {
      apiBaseURL?: string
      storage: StorageInterface
    },
  ) {
    this.apiBaseURL = joinURL(apiBaseURL, '/v1')
    this.storage = storage
  }

  // ============================================================================
  // Private Helper Methods
  // ============================================================================

  private async fetch(
    url: string,
    options: RequestInit = {},
  ): Promise<Response> {
    const headers = new Headers(options.headers || {})

    // Set default content type if not form data and not already set
    if (!options.body || !(options.body instanceof FormData)) {
      if (!headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
      }
    }

    // Add session token if available
    const sessionToken = await this.storage.getSessionToken()
    if (sessionToken) {
      headers.set('X-Session-Token', sessionToken)
    }

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: options.credentials || 'include',
      mode: options.mode || 'cors',
    })

    // Handle session token from response
    try {
      const clonedResponse = response.clone()
      const data = await clonedResponse.json()
      if (data?.meta?.session_token) {
        await this.storage.setSessionToken(data.meta.session_token)
      }
    } catch {
      // Response might not be JSON
    }

    // Handle 410 Gone (session expired)
    if (response.status === 410) {
      await this.storage.setSessionToken(null)
    }

    return response
  }

  private async request<T>(
    path: string,
    options: RequestInit = {},
  ): Promise<T> {
    const url = joinURL(this.apiBaseURL, path)
    const response = await this.fetch(url, options)
    const data = await response.json()

    // For error responses, throw the error so TanStack Query treats it as a failure
    if (!response.ok) {
      throw data
    }

    return data as T
  }

  // ============================================================================
  // Configuration
  // ============================================================================

  async getConfiguration(): Promise<ConfigurationResponse> {
    return this.request<ConfigurationResponse>('/auth/config')
  }

  // ============================================================================
  // Authentication - Core
  // ============================================================================

  async getAuthenticationStatus(): Promise<
    AuthenticatedResponse | NotAuthenticatedResponse
  > {
    try {
      return await this.request<AuthenticatedResponse | NotAuthenticatedResponse>(
        '/auth/session',
      )
    } catch (error: any) {
      if (error.status === 401) {
        return error
      }
      throw error
    }
  }

  async login(data: LoginRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/login',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async logout(): Promise<NotAuthenticatedResponse> {
    try {
      return await this.request<NotAuthenticatedResponse>('/auth/session', {
        method: 'DELETE',
      })
    } catch (error: any) {
      if (error.status === 401) {
        return error
      }
      throw error
    }
  }

  async signup(data: SignupRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/signup',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async reauthenticate(data: ReauthenticateRequest): Promise<
    AuthenticatedResponse
 > {
    return this.request<AuthenticatedResponse>(
      '/auth/reauthenticate',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  // ============================================================================
  // Authentication - Login by Code
  // ============================================================================

  async requestLoginCode(data: LoginByCodeRequest): Promise<
    AuthenticationResponse
  > {
    try {
      return await this.request<AuthenticationResponse>(
        '/auth/code/request',
        {
          method: 'POST',
          body: JSON.stringify(data),
        },
      )
    } catch (error: any) {
      if (error.status === 401) {
        return error
      }
      throw error
    }
  }

  async confirmLoginCode(data: ConfirmLoginCodeRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/code/confirm',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  // ============================================================================
  // Email Management
  // ============================================================================

  async listEmailAddresses(): Promise<EmailAddressesResponse> {
    return this.request<EmailAddressesResponse>('/account/email')
  }

  async addEmailAddress(data: EmailAddressRequest): Promise<
    EmailAddressesResponse
 > {
    return this.request<EmailAddressesResponse>(
      '/account/email',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async checkEmail(data: EmailAddressRequest): Promise<
    EmailCheckResponse
 > {
    return this.request<EmailCheckResponse>(
      '/auth/email/check',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async removeEmailAddress(data: EmailAddressRequest): Promise<
    EmailAddressesResponse
 > {
    return this.request<EmailAddressesResponse>(
      '/account/email',
      {
        method: 'DELETE',
        body: JSON.stringify(data),
      },
    )
  }

  async changePrimaryEmailAddress(data: EmailPrimaryRequest): Promise<
    EmailAddressesResponse
 > {
    return this.request<EmailAddressesResponse>(
      '/account/email',
      {
        method: 'PATCH',
        body: JSON.stringify(data),
      },
    )
  }

  async requestEmailVerification(data: EmailAddressRequest): Promise<
    { status: 200 }
 > {
    return this.request<{ status: 200 }>(
      '/account/email',
      {
        method: 'PUT',
        body: JSON.stringify(data),
      },
    )
  }

  async getEmailVerificationInfo(key: string): Promise<
    EmailVerificationInfoResponse
 > {
    return this.request<EmailVerificationInfoResponse>(
      withQuery('/auth/email/verify', { key }),
    )
  }

  async verifyEmail(data: EmailVerificationRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/email/verify',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async resendEmailVerification(): Promise<{ status: 200 }> {
    return this.request<{ status: 200 }>(
      '/auth/email/verify/resend',
      {
        method: 'POST',
      },
    )
  }

  // ============================================================================
  // Phone Management
  // ============================================================================

  async getPhoneNumber(): Promise<PhoneNumberResponse> {
    return this.request<PhoneNumberResponse>('/account/phone')
  }

  async updatePhoneNumber(phone: string): Promise<
    PhoneNumberResponse
 > {
    return this.request<PhoneNumberResponse>(
      '/account/phone',
      {
        method: 'PUT',
        body: JSON.stringify({ phone }),
      },
    )
  }

  async removePhoneNumber(): Promise<{ status: 200 }> {
    return this.request<{ status: 200 }>(
      '/account/phone',
      {
        method: 'DELETE',
      },
    )
  }

  async verifyPhone(data: PhoneVerificationRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/phone/verify',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async resendPhoneVerification(): Promise<{ status: 200 }> {
    return this.request<{ status: 200 }>(
      '/auth/phone/verify/resend',
      {
        method: 'POST',
      },
    )
  }

  // ============================================================================
  // Password Management
  // ============================================================================

  async requestPassword(data: PasswordResetRequest): Promise<
    { status: 200 } | AuthenticationResponse
  > {
    try {
      return await this.request<{ status: 200 } | AuthenticationResponse>(
        '/auth/password/request',
        {
          method: 'POST',
          body: JSON.stringify(data),
        },
      )
    } catch (error: any) {
      if (error.status === 401) {
        return error
      }
      throw error
    }
  }

  async getPasswordResetInfo(key: string): Promise<
    PasswordResetInfoResponse
 > {
    return this.request<PasswordResetInfoResponse>(
      withQuery('/auth/password/reset', { key }),
    )
  }

  async resetPassword(data: PasswordResetConfirmRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    try {
      return await this.request<AuthenticatedResponse | AuthenticationResponse>(
        '/auth/password/reset',
        {
          method: 'POST',
          body: JSON.stringify(data),
        },
      )
    } catch (error: any) {
      if (error.status === 401) {
        return error
      }
      throw error
    }
  }

  async changePassword(data: PasswordChangeRequest): Promise<
    { status: 200 }
 > {
    return this.request<{ status: 200 }>(
      '/account/password/change',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  // ============================================================================
  // Social/Provider Authentication
  // ============================================================================

  async listProviderAccounts(): Promise<ProviderAccountsResponse> {
    return this.request<ProviderAccountsResponse>('/account/providers')
  }

  async disconnectProviderAccount(data: ProviderDisconnectRequest): Promise<
    ProviderAccountsResponse
 > {
    return this.request<ProviderAccountsResponse>(
      '/account/providers',
      {
        method: 'DELETE',
        body: JSON.stringify(data),
      },
    )
  }

  async providerRedirect(data: ProviderRedirectRequest): Promise<void> {
    const url = joinURL(this.apiBaseURL, '/auth/provider/redirect')
    const form = document.createElement('form')
    form.method = 'POST'
    form.action = url

    const fields = { process: 'login', ...data }
    Object.entries(fields).forEach(([key, value]) => {
      const input = document.createElement('input')
      input.type = 'hidden'
      input.name = key
      input.value = value
      form.appendChild(input)
    })

    document.body.appendChild(form)
    form.submit()
  }

  async providerCallback(provider: string, data: Record<string, any>): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      joinURL('/auth/provider/callback', provider),
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async providerToken(data: ProviderTokenRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/provider/token',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async getProviderSignup(): Promise<ProviderSignupResponse> {
    return this.request<ProviderSignupResponse>(
      '/auth/provider/signup',
    )
  }

  async providerSignup(data: ProviderSignupRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/provider/signup',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  // ============================================================================
  // Multi-Factor Authentication
  // ============================================================================

  async listAuthenticators(): Promise<AuthenticatorsResponse> {
    return this.request<AuthenticatorsResponse>('/account/authenticators')
  }

  async getTOTPAuthenticator(): Promise<
    TOTPAuthenticatorResponse | NoTOTPAuthenticatorResponse
  > {
    return this.request<TOTPAuthenticatorResponse | NoTOTPAuthenticatorResponse>(
      '/account/authenticators/totp',
    )
  }

  async activateTOTP(data: TOTPActivateRequest): Promise<
    TOTPAuthenticatorResponse
 > {
    return this.request<TOTPAuthenticatorResponse>(
      '/account/authenticators/totp',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async deactivateTOTP(): Promise<{ status: 200 }> {
    return this.request<{ status: 200 }>(
      '/account/authenticators/totp',
      {
        method: 'DELETE',
      },
    )
  }

  async listRecoveryCodes(): Promise<
    SensitiveRecoveryCodesAuthenticatorResponse | { status: 404 }
  > {
    return this.request<
      SensitiveRecoveryCodesAuthenticatorResponse | { status: 404 }
    >('/account/authenticators/recovery-codes')
  }

  async regenerateRecoveryCodes(): Promise<
    SensitiveRecoveryCodesAuthenticatorResponse
 > {
    return this.request<SensitiveRecoveryCodesAuthenticatorResponse>(
      '/account/authenticators/recovery-codes',
      {
        method: 'POST',
      },
    )
  }

  async mfaAuthenticate(data: MFAAuthenticateRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/2fa/authenticate',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async mfaReauthenticate(): Promise<AuthenticatedResponse> {
    return this.request<AuthenticatedResponse>(
      '/auth/2fa/reauthenticate',
      {
        method: 'POST',
      },
    )
  }

  async mfaTrust(data: MFATrustRequest): Promise<
    AuthenticatedResponse
  > {
    const url = joinURL(this.apiBaseURL, '/auth/2fa/trust')
    const response = await this.fetch(url, {
      method: 'POST',
      body: JSON.stringify(data),
    })
    const result = await response.json()

    // For error responses, throw the error so TanStack Query treats it as a failure
    if (!response.ok && result.errors) {
      throw result
    }

    return result as AuthenticatedResponse
  }

  // ============================================================================
  // WebAuthn
  // ============================================================================

  async getWebAuthnSignupOptions(): Promise<
    WebAuthnCredentialCreationOptions
 > {
    return this.request<WebAuthnCredentialCreationOptions>(
      '/auth/webauthn/signup',
    )
  }

  async webAuthnSignup(data: WebAuthnSignupRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/webauthn/signup',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async getWebAuthnLoginOptions(): Promise<
    WebAuthnCredentialRequestOptions
 > {
    return this.request<WebAuthnCredentialRequestOptions>(
      '/auth/webauthn/login',
    )
  }

  async webAuthnLogin(data: WebAuthnLoginRequest): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/webauthn/login',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
    )
  }

  async getWebAuthnAuthenticateOptions(): Promise<
    WebAuthnCredentialRequestOptions
 > {
    return this.request<WebAuthnCredentialRequestOptions>(
      '/auth/webauthn/authenticate',
    )
  }

  async webAuthnAuthenticate(credential: string): Promise<
    AuthenticatedResponse | AuthenticationResponse
  > {
    return this.request<AuthenticatedResponse | AuthenticationResponse>(
      '/auth/webauthn/authenticate',
      {
        method: 'POST',
        body: JSON.stringify({ credential }),
      },
    )
  }

  async getWebAuthnReauthenticateOptions(): Promise<
    WebAuthnCredentialRequestOptions
 > {
    return this.request<WebAuthnCredentialRequestOptions>(
      '/auth/webauthn/reauthenticate',
    )
  }

  async webAuthnReauthenticate(credential: string): Promise<
    AuthenticatedResponse
 > {
    return this.request<AuthenticatedResponse>(
      '/auth/webauthn/reauthenticate',
      {
        method: 'POST',
        body: JSON.stringify({ credential }),
      },
    )
  }

  async listWebAuthnCredentials(): Promise<AuthenticatorsResponse> {
    return this.request<AuthenticatorsResponse>(
      '/account/authenticators/webauthn',
    )
  }

  async deleteWebAuthnCredential(id: string): Promise<
    AuthenticatorsResponse
 > {
    return this.request<AuthenticatorsResponse>(
      '/account/authenticators/webauthn',
      {
        method: 'DELETE',
        body: JSON.stringify({ id }),
      },
    )
  }

  // ============================================================================
  // Session Management
  // ============================================================================

  async listSessions(): Promise<SessionsResponse> {
    return this.request<SessionsResponse>('/auth/sessions')
  }

  async deleteSession({ id }: { id?: number } = {}): Promise<SessionsResponse> {
    const path = id ? '/auth/sessions' : '/auth/session'
    return this.request<SessionsResponse>(path, {
      method: 'DELETE',
      body: id ? JSON.stringify({ id }) : undefined,
    })
  }
}

// // ============================================================================
// // Singleton Management
// // ============================================================================

// export function initializeClient(config: {
//   baseUrl?: string
//   csrfTokenEndpoint?: string
//   clientType?: ClientType
//   storage?: StorageInterface
// }): AllAuthClient {
//   if (!clientInstance) {
//     const { baseUrl = '', csrfTokenEndpoint, clientType = 'browser', storage } = config
//     clientInstance = new AllAuthClient(
//       baseUrl,
//       csrfTokenEndpoint,
//       clientType,
//       storage || getStorage(clientType, baseUrl),
//     )
//   }
//   return clientInstance
// }

// export function getClient(): AllAuthClient {
//   if (!clientInstance) {
//     throw new Error(
//       'AllauthClient not initialized. Please wrap your app with AllauthProvider or call initializeClient first.',
//     )
//   }
//   return clientInstance
// }

export const getClient = ({
  apiBaseURL = '/',
  storage,
}: {
  apiBaseURL?: string
  storage: StorageInterface
}) => {
  if (!clientInstance) {
    clientInstance = new AllAuthClient({ apiBaseURL, storage })
  }
  return clientInstance
}
