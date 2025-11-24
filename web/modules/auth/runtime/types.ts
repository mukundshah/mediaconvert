export interface AuthMeta {
  required?: boolean
  namespace?: string
  roles?: string[]
  permissions?: string[]
  loginRoute?: string
  redirectIfLoggedIn?: string | false
  redirectIfNotAllowed?: string | false
}
