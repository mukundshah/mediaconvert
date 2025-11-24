<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { Play } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { z } from 'zod'

const formSchema = toTypedSchema(z.object({
  email: z.email(),
  password: z.string().min(8),
}))

const form = useForm({
  validationSchema: formSchema,
})

const onSubmit = form.handleSubmit((values) => {
  console.log('Form submitted!', values)
})
</script>

<template>
  <div class="flex flex-col items-center justify-center h-screen -mt-[52px] space-y-8">
    <div class="flex items-center space-x-2">
      <div class="flex size-8 items-center justify-center rounded-lg bg-primary">
        <Play class="size-4 text-primary-foreground" />
      </div>
      <span class="text-2xl font-bold">MediaConvert</span>
    </div>
    <Card class="w-sm -mt-8">
      <CardHeader>
        <CardTitle>Login</CardTitle>
        <CardDescription>Enter your credentials to access your account.</CardDescription>
      </CardHeader>
      <CardContent>
        <form id="login" @submit.prevent="onSubmit">
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
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </FormField>
            </FieldGroup>

            <!-- Password Input -->
            <FieldGroup>
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
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </FormField>
            </FieldGroup>
          </div>
        </form>
      </CardContent>
      <CardFooter class="flex flex-col space-y-2">
        <Button class="w-full" for="login">
          Login
        </Button>
        <div class="text-sm text-center text-muted-foreground">
          Don't have an account?
          <NuxtLink class="underline" to="/auth/register">
            Register
          </NuxtLink>
        </div>
      </CardFooter>
    </Card>
  </div>
</template>
