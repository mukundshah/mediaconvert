<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input, PasswordInput } from '@/components/ui/input'
import { cn } from '@/utils/style'

const props = defineProps<{
  class?: string
}>()

// Form state
const step = ref<'EMAIL' | 'PASSWORD' | 'CODE'>('EMAIL')
const isPending = ref(false)
const isMagicCodePending = ref(false)
const _seePassword = ref(false)
const email = ref('')
const isExistingUser = ref<boolean | null>(null)
const signupAllowed = ref(false)
const isLoginByCodeAvailable = ref(false)

// Get available auth methods
const { checkEmail } = useEmailManagement()
const { login, signup, requestLoginCode, confirmLoginCode } = useAuth()
const { connect: redirectToProvider } = useSocialAuth()
const { getConfiguration } = useAuthConfiguration()

const { data: configuration } = await useAsyncData('auth-configuration', () => getConfiguration())

// Computed properties for UI
const formTitle = computed(() => {
  if (step.value === 'EMAIL') return 'Sign in to 4Fin'
  return isExistingUser.value ? 'Welcome back' : 'Create your account'
})

const formDescription = computed(() => {
  if (step.value === 'EMAIL') return 'Enter your email to continue'
  return isExistingUser.value ? 'Sign in to access your account' : 'Complete your registration'
})

// Initialize form with dynamic validation
const { handleSubmit, resetForm, setFieldError } = useForm<{
  email: string
  password: string
  code: string
}>({
  // validationSchema: schema,
})

// Single submit handler for all steps
const onSubmit = handleSubmit(async (values) => {
  try {
    isPending.value = true

    switch (step.value) {
      case 'EMAIL': {
        const response = await checkEmail({ email: values.email })
        email.value = response.email

        // Determine if user exists based on whether they have a password
        isExistingUser.value = response.login_by_password || (response.login_by_code && !response.signup_allowed)
        signupAllowed.value = response.signup_allowed
        isLoginByCodeAvailable.value = response.login_by_code

        if (!isExistingUser.value) {
          // Handle signup not allowed
          if (!signupAllowed.value) {
            toast.error('Signup is not allowed for this instance')
            return
          }

          if (response.login_by_code) {
            await signup({ email: values.email })
            await requestLoginCode({ email: values.email })
            toast.success('Verification code sent to your email')
            step.value = 'CODE'
            break
          } else {
            throw new Error('No suitable authentication method available')
          }
        } else {
          // For existing users, determine the auth method based on available options
          if (response.login_by_password) {
            step.value = 'PASSWORD'
          } else if (response.login_by_code) {
            await requestLoginCode({ email: values.email })
            toast.success('Verification code sent to your email')
            step.value = 'CODE'
          } else {
            throw new Error('No suitable authentication method available')
          }
        }
        break
      }

      case 'PASSWORD': {
        const authData = {
          email: email.value,
          password: values.password,
        }

        if (isExistingUser.value) {
          await login(authData)
          toast.success('Logged in successfully')
        } else {
          await signup(authData)
          toast.success('Account created successfully')
        }
        resetForm()
        // Take user to dashboard after successful login
        await navigateTo({ name: 'dashboard' })
        break
      }

      case 'CODE': {
        await confirmLoginCode({ code: values.code })
        toast.success('Logged in successfully')
        // Take user to dashboard after successful login
        await navigateTo({ name: 'dashboard' })
        break
      }
    }
  } catch (error: any) {
    // Handle validation errors (400 status with errors array)
    if (error.status === 400 && error.data.errors && Array.isArray(error.data.errors)) {
      error.data.errors.forEach((validationError: any) => {
        if (validationError.param && validationError.message) {
          setFieldError(validationError.param, validationError.message)
        }
      })
    } else if (error.status === 409) {
      // Handle conflict errors (e.g., user already logged in)
      toast.error('You might already be logged in')
      // Redirect to dashboard since they're already authenticated
      await navigateTo({ name: 'dashboard' })
    } else {
      // Handle other errors with toast
      toast.error(error.message || 'Authentication failed')
    }
  } finally {
    isPending.value = false
  }
})

// Switch to magic code method if available
const switchToMagicCode = async () => {
  if (!isLoginByCodeAvailable.value) return
  try {
    isMagicCodePending.value = true
    await requestLoginCode({ email: email.value })
    toast.success('Verification code sent to your email')
    step.value = 'CODE'
  } catch (error: any) {
    // Handle validation errors (400 status with errors array)
    if (error.status === 400 && error.data.errors && Array.isArray(error.data.errors)) {
      error.data.errors.forEach((validationError: any) => {
        if (validationError.param && validationError.message) {
          setFieldError(validationError.param, validationError.message)
        }
      })
    } else {
      // Handle other errors with toast
      toast.error(error.message || 'Failed to request verification code')
    }
  } finally {
    isMagicCodePending.value = false
  }
}
</script>

<template>
  <div :class="cn('flex flex-col gap-6', props.class)">
    <Card class="border-gray-200">
      <CardHeader class="text-center pb-4">
        <div class="flex items-center justify-between">
          <Button
            v-if="step !== 'EMAIL'"
            size="sm"
            type="button"
            variant="ghost"
            @click="step = 'EMAIL'; resetForm()"
          >
            <Icon class="h-4 w-4" name="lucide:arrow-left" />
          </Button>
          <div v-else class="w-8"></div>
          <CardTitle class="text-xl font-bold text-gray-900">
            {{ formTitle }}
          </CardTitle>
          <div class="w-8"></div>
        </div>
        <CardDescription class="text-gray-600">
          {{ formDescription }}
        </CardDescription>
      </CardHeader>
      <CardContent class="pt-0">
        <form @submit.prevent="onSubmit">
          <div class="grid gap-6">
            <!-- Email Input -->
            <FieldGroup>
              <FormField
                v-slot="{ componentField, errors }"
                name="email"
                :validate-on-blur="false"
              >
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="form-email">
                    Email
                  </FieldLabel>
                  <Input
                    id="form-email"
                    v-bind="componentField"
                    placeholder="m@example.com"
                    type="email"
                    :disabled="step !== 'EMAIL'"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </FormField>
            </FieldGroup>

            <!-- Password Input -->
            <FieldGroup v-if="step === 'PASSWORD'">
              <FormField
                v-slot="{ componentField, errors }"
                name="password"
                :validate-on-blur="false"
              >
                <Field :data-invalid="!!errors.length">
                  <div class="flex items-center justify-between">
                    <FieldLabel for="form-password">
                      Password
                    </FieldLabel>
                    <Button class="p-0 text-xs h-min" type="button" variant="link">
                      Forgot password?
                    </Button>
                  </div>
                  <PasswordInput
                    id="form-password"
                    v-bind="componentField"
                    placeholder="Enter your password"
                    :disabled="step !== 'PASSWORD'"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </FormField>
            </FieldGroup>

            <FieldGroup v-if="step === 'CODE'">
              <FormField
                v-slot="{ componentField, errors }"
                name="code"
                :validate-on-blur="false"
              >
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="form-code">
                    Verification Code
                  </FieldLabel>
                  <InputOTP
                    id="form-code"
                    v-bind="componentField"
                    :aria-invalid="!!errors.length"
                    :maxlength="6"
                  >
                    <InputOTPGroup>
                      <InputOTPSlot :index="0" />
                      <InputOTPSlot :index="1" />
                      <InputOTPSlot :index="2" />
                      <InputOTPSlot :index="3" />
                      <InputOTPSlot :index="4" />
                      <InputOTPSlot :index="5" />
                    </InputOTPGroup>
                  </InputOTP>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </FormField>
            </FieldGroup>

            <!-- Submit Button -->
            <Button
              type="submit"
              :disabled="isPending || isMagicCodePending"
            >
              <Spinner v-if="isPending" />
              {{ isPending ? 'Processing...' : 'Continue' }}
            </Button>

            <!-- Switch to Magic Code -->
            <template v-if="step === 'PASSWORD' && isLoginByCodeAvailable">
              <div class="relative text-center text-sm">
                <div class="absolute left-0 right-0 top-1/2 z-0 h-px bg-gray-300"></div>
                <span class="relative z-10 bg-white px-2 text-gray-500">
                  OR
                </span>
                <div class="absolute left-0 right-0 top-1/2 z-0 h-px bg-gray-300"></div>
              </div>
              <Button
                type="button"
                variant="outline"
                :disabled="isMagicCodePending"
                @click="switchToMagicCode"
              >
                <Spinner v-if="isMagicCodePending" />
                <Icon v-else class="mr-2 h-4 w-4" name="lucide:mail" />
                {{ isMagicCodePending ? 'Sending...' : 'Email sign-in code' }}
              </Button>
            </template>

            <!-- Social Sign In -->
            <div v-if="step === 'EMAIL' && configuration?.socialaccount?.providers?.length" class="flex flex-col gap-4">
              <div class="relative text-center text-sm">
                <div class="absolute left-0 right-0 top-1/2 z-0 h-px bg-gray-300"></div>
                <span class="relative z-10 bg-white px-2 text-gray-500">
                  OR
                </span>
                <div class="absolute left-0 right-0 top-1/2 z-0 h-px bg-gray-300"></div>
              </div>

              <Button
                v-for="provider in configuration?.socialaccount?.providers?.filter((provider) => provider.flows.includes('provider_redirect')) || []"
                :key="provider.id"
                type="button"
                variant="outline"
                @click="redirectToProvider({ provider: provider.id })"
              >
                <Icon class="mr-2 h-4 w-4" :name="`social-auth:${provider.id}`" />
                Continue with {{ provider.name }}
              </Button>
            </div>
          </div>
        </form>
      </CardContent>
    </Card>
    <div class="text-center text-xs text-gray-500 [&_a]:underline [&_a]:decoration-dotted [&_a]:hover:decoration-solid [&_a]:hover:text-gray-700 -mx-6">
      By clicking continue, you agree to our <NuxtLink to="/terms">
        Terms of Service
      </NuxtLink>
      and <NuxtLink to="/privacy">
        Privacy Policy
      </NuxtLink>.
    </div>
  </div>
</template>
