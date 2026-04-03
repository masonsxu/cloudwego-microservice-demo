<template>
  <div class="user-create">
    <div class="page-header">
      <button class="back-btn" @click="handleBack">
        <ArrowLeft class="h-4 w-4" />
        {{ t('common.back') }}
      </button>
      <h1 class="page-title">{{ t('user.createUser') }}</h1>
    </div>

    <div class="form-content">
      <form class="space-y-6" @submit.prevent="handleSubmit">
        <div class="form-section">
          <div class="card-header">
            <User class="h-5 w-5" />
            <span>{{ t('user.userDetail') }}</span>
          </div>
          <div class="form-grid">
            <div class="form-field">
              <label class="form-label">{{ t('user.username') }} <span class="text-red-500">*</span></label>
              <Input
                v-model="form.username"
                :placeholder="t('user.username')"
                maxlength="20"
                :class="{ 'border-red-500': errors.username }"
              />
              <span v-if="errors.username" class="error-text">{{ errors.username }}</span>
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('auth.password') }} <span class="text-red-500">*</span></label>
              <div class="password-input-wrapper">
                <Input
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  :placeholder="t('auth.password')"
                  maxlength="50"
                  :class="{ 'border-red-500': errors.password }"
                />
                <button type="button" class="password-toggle" @click="showPassword = !showPassword">
                  <Eye v-if="!showPassword" class="h-4 w-4" />
                  <EyeOff v-else class="h-4 w-4" />
                </button>
              </div>
              <span v-if="errors.password" class="error-text">{{ errors.password }}</span>
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.firstName') }}</label>
              <Input v-model="form.first_name" :placeholder="t('user.firstName')" maxlength="50" />
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.lastName') }}</label>
              <Input v-model="form.last_name" :placeholder="t('user.lastName')" maxlength="50" />
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.realName') }}</label>
              <Input v-model="form.real_name" :placeholder="t('user.realName')" maxlength="100" />
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.email') }}</label>
              <Input v-model="form.email" :placeholder="t('user.email')" type="email" :class="{ 'border-red-500': errors.email }" />
              <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.phone') }}</label>
              <Input v-model="form.phone" :placeholder="t('user.phone')" maxlength="20" :class="{ 'border-red-500': errors.phone }" />
              <span v-if="errors.phone" class="error-text">{{ errors.phone }}</span>
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.gender.label') }}</label>
              <RadioGroup v-model="genderString" class="flex gap-4">
                <div class="flex items-center gap-2">
                  <RadioGroupItem value="0" id="gender-0" />
                  <label for="gender-0" class="text-sm">{{ t('user.gender.unknown') }}</label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioGroupItem value="1" id="gender-1" />
                  <label for="gender-1" class="text-sm">{{ t('user.gender.male') }}</label>
                </div>
                <div class="flex items-center gap-2">
                  <RadioGroupItem value="2" id="gender-2" />
                  <label for="gender-2" class="text-sm">{{ t('user.gender.female') }}</label>
                </div>
              </RadioGroup>
            </div>
          </div>
        </div>

        <div class="form-section">
          <div class="card-header">
            <Building2 class="h-5 w-5" />
            <span>{{ t('user.organization') }}</span>
          </div>
          <div class="form-grid single-col">
            <div class="form-field">
              <label class="form-label">{{ t('user.organization') }}</label>
              <Select v-model="form.organization_id">
                <SelectTrigger>
                  <SelectValue :placeholder="t('user.organization')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="org in organizations" :key="org.id" :value="org.id">
                    {{ org.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.employeeId') }}</label>
              <Input v-model="form.employee_id" :placeholder="t('user.employeeId')" maxlength="50" />
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.professionalTitle') }}</label>
              <Input v-model="form.professional_title" :placeholder="t('user.professionalTitle')" maxlength="100" />
            </div>
          </div>
        </div>

        <div class="form-section">
          <div class="card-header">
            <KeyRound class="h-5 w-5" />
            <span>{{ t('user.roles') }}</span>
          </div>
          <div class="form-grid single-col">
            <div class="form-field">
              <label class="form-label">{{ t('user.roles') }}</label>
              <DropdownMenu>
                <DropdownMenuTrigger as-child>
                  <button class="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
                    <span class="truncate mr-2">
                      {{ roleIdsStrings.length ? roles.filter(r => roleIdsStrings.includes(String(r.id))).map(r => r.name).join(', ') : t('user.roles') }}
                    </span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 opacity-50 flex-shrink-0"><path d="m6 9 6 6 6-6"/></svg>
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent class="w-[var(--radix-dropdown-menu-trigger-width)] max-h-[200px] overflow-y-auto">
                  <DropdownMenuItem v-for="role in roles" :key="role.id" class="cursor-pointer" @select.prevent="toggleRole(String(role.id))">
                    <Checkbox :checked="roleIdsStrings.includes(String(role.id))" class="mr-2" />
                    <span class="truncate">{{ role.name }}</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.mustChangePassword') }}</label>
              <Switch v-model="form.must_change_password" />
            </div>
            <div class="form-field">
              <label class="form-label">{{ t('user.accountExpiry') }}</label>
              <Input
                v-model="accountExpiryDisplay"
                type="datetime-local"
                :placeholder="t('user.accountExpiry')"
                @change="handleAccountExpiryChange"
              />
            </div>
          </div>
        </div>

        <div class="form-actions">
          <Button type="button" variant="outline" @click="handleBack">{{ t('common.cancel') }}</Button>
          <Button type="submit" :disabled="submitting">
            {{ t('common.submit') }}
          </Button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { User, Building2, KeyRound, ArrowLeft, Eye, EyeOff } from 'lucide-vue-next'
import { createUser } from '@/api/user'
import { getOrganizationList } from '@/api/organization'
import { getRoleList } from '@/api/role'
import type { CreateUserRequest } from '@/types/user'
import type { Organization } from '@/types/organization'
import type { RoleDefinition } from '@/types/role'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import {
  Select, SelectTrigger, SelectValue, SelectContent, SelectItem
} from '@/components/ui/select'
import {
  DropdownMenu, DropdownMenuTrigger, DropdownMenuContent,
  DropdownMenuItem
} from '@/components/ui/dropdown-menu'
import { Checkbox } from '@/components/ui/checkbox'

const router = useRouter()
const { t } = useI18n()

const submitting = ref(false)
const organizations = ref<Organization[]>([])
const roles = ref<RoleDefinition[]>([])
const showPassword = ref(false)
const accountExpiryDisplay = ref('')
const genderString = ref('0')
const roleIdsStrings = ref<string[]>([])
const errors = reactive<Record<string, string>>({})

const form = reactive<CreateUserRequest>({
  username: '',
  password: '',
  email: '',
  phone: '',
  first_name: '',
  last_name: '',
  real_name: '',
  professional_title: '',
  employee_id: '',
  must_change_password: false,
  account_expiry: undefined,
  gender: 0,
  role_ids: [],
  organization_id: ''
})

onMounted(() => {
  fetchOrganizations()
  fetchRoles()
})

async function fetchOrganizations() {
  try {
    const response = await getOrganizationList({ limit: 1000 })
    organizations.value = response.organizations || []
  } catch (error) {
    console.error('Failed to fetch organizations:', error)
  }
}

async function fetchRoles() {
  try {
    const response = await getRoleList({ limit: 1000 })
    roles.value = response.roles || []
  } catch (error) {
    console.error('Failed to fetch roles:', error)
  }
}

function handleBack() {
  router.back()
}

function handleAccountExpiryChange() {
  if (accountExpiryDisplay.value) {
    form.account_expiry = new Date(accountExpiryDisplay.value).getTime() / 1000
  } else {
    form.account_expiry = undefined
  }
}

function validate(): boolean {
  Object.keys(errors).forEach(key => delete errors[key])
  let valid = true

  if (!form.username) {
    errors.username = '请输入用户名'
    valid = false
  } else if (form.username.length < 3 || form.username.length > 20) {
    errors.username = '用户名长度在 3 到 20 个字符'
    valid = false
  } else if (!/^[a-zA-Z0-9_-]+$/.test(form.username)) {
    errors.username = '用户名只能包含字母、数字、下划线和短横线'
    valid = false
  }

  if (!form.password) {
    errors.password = '请输入密码'
    valid = false
  } else if (form.password.length < 6) {
    errors.password = '密码至少6位'
    valid = false
  }

  if (form.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = '请输入正确的邮箱地址'
    valid = false
  }

  if (form.phone && !/^1[3-9]\d{9}$/.test(form.phone)) {
    errors.phone = '请输入正确的手机号'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  form.gender = Number(genderString.value) as any
  form.role_ids = roleIdsStrings.value

  try {
    submitting.value = true
    await createUser(form)
    toast.success(t('common.operationSuccess'))
    router.push('/users')
  } catch (error: any) {
    console.error('Failed to create user:', error)
    if (error.message) {
      toast.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}

function toggleRole(roleId: string) {
  const idx = roleIdsStrings.value.indexOf(roleId)
  if (idx === -1) {
    roleIdsStrings.value.push(roleId)
  } else {
    roleIdsStrings.value.splice(idx, 1)
  }
}
</script>

<style scoped lang="scss">
.user-create {
  padding: 20px;

  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;

    .back-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 12px;
      border: 1px solid hsl(var(--border));
      border-radius: 8px;
      background: var(--bg-card);
      color: var(--c-text-sub);
      font-size: 14px;
      cursor: pointer;
      transition: all 0.2s ease;

      &:hover {
        background: var(--bg-input);
        color: var(--c-text-main);
      }
    }

    .page-title {
      font-size: 20px;
      font-weight: 700;
      color: var(--c-primary);
      font-family: 'Inter', sans-serif;
      margin: 0;
    }
  }

  .form-content {
    margin-top: 20px;

    .form-section {
      background: var(--bg-card);
      border: 1px solid hsl(var(--border) / 0.6);
      border-radius: 18px;
      box-shadow: var(--shadow-card);
      padding: 20px;

      .card-header {
        display: flex;
        align-items: center;
        gap: 10px;
        color: var(--c-primary);
        font-family: 'Inter', sans-serif;
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 20px;
        padding-bottom: 12px;
        border-bottom: 1px solid hsl(var(--border) / 0.6);
      }

      .form-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 16px;

        &.single-col {
          grid-template-columns: 1fr;
        }
      }

      .form-field {
        display: flex;
        flex-direction: column;
        gap: 6px;

        .form-label {
          color: var(--c-text-sub);
          font-size: 13px;
          font-weight: 500;
        }

        .error-text {
          color: #ef4444;
          font-size: 12px;
        }

        .password-input-wrapper {
          position: relative;
          display: flex;
          align-items: center;

          :deep(input) {
            padding-right: 40px;
          }

          .password-toggle {
            position: absolute;
            right: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            width: 24px;
            height: 24px;
            border: none;
            background: transparent;
            color: var(--c-text-muted);
            cursor: pointer;
            border-radius: 4px;

            &:hover {
              background: var(--bg-card);
              color: var(--c-text-main);
            }
          }
        }
      }
    }

    .form-actions {
      margin-top: 24px;
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      padding: 16px;
      background: var(--bg-card);
      border: 1px solid hsl(var(--border) / 0.6);
      border-radius: 16px;
      box-shadow: var(--shadow-card);

      button {
        min-width: 120px;
      }
    }
  }
}
</style>
