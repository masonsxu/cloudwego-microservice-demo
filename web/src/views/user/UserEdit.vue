<template>
  <div class="user-edit">
    <el-page-header @back="handleBack" :title="t('common.back')">
      <template #content>
        <span class="page-title">{{ t('user.editUser') }}</span>
      </template>
    </el-page-header>

    <div class="form-content">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="140px"
        v-loading="loading"
      >
        <el-card class="form-section">
          <template #header>
            <div class="card-header">
              <el-icon><User /></el-icon>
              <span>{{ t('user.userDetail') }}</span>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.username')">
                <el-input
                  :model-value="userDetail?.username"
                  disabled
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.status')">
                <el-tag :type="getStatusType(userDetail?.status)">
                  {{ t(`user.status.${getStatusKey(userDetail?.status)}`) }}
                </el-tag>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.firstName')" prop="first_name">
                <el-input
                  v-model="form.first_name"
                  :placeholder="t('user.firstName')"
                  maxlength="50"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.lastName')" prop="last_name">
                <el-input
                  v-model="form.last_name"
                  :placeholder="t('user.lastName')"
                  maxlength="50"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.realName')" prop="real_name">
                <el-input
                  v-model="form.real_name"
                  :placeholder="t('user.realName')"
                  maxlength="100"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.email')" prop="email">
                <el-input
                  v-model="form.email"
                  :placeholder="t('user.email')"
                  type="email"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.phone')" prop="phone">
                <el-input
                  v-model="form.phone"
                  :placeholder="t('user.phone')"
                  maxlength="20"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :md="12">
              <el-form-item :label="t('user.gender')" prop="gender">
                <el-radio-group v-model="form.gender">
                  <el-radio :value="0">{{ t('user.gender.unknown') }}</el-radio>
                  <el-radio :value="1">{{ t('user.gender.male') }}</el-radio>
                  <el-radio :value="2">{{ t('user.gender.female') }}</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <el-card class="form-section">
          <template #header>
            <div class="card-header">
              <el-icon><OfficeBuilding /></el-icon>
              <span>{{ t('user.organization') }}</span>
            </div>
          </template>

          <el-form-item :label="t('user.organization')" prop="organization_id">
            <el-select
              v-model="form.organization_id"
              :placeholder="t('user.organization')"
              filterable
              clearable
              style="width: 100%"
            >
              <el-option
                v-for="org in organizations"
                :key="org.id"
                :label="org.name"
                :value="org.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item :label="t('user.employeeId')" prop="employee_id">
            <el-input
              v-model="form.employee_id"
              :placeholder="t('user.employeeId')"
              maxlength="50"
            />
          </el-form-item>

          <el-form-item :label="t('user.professionalTitle')" prop="professional_title">
            <el-input
              v-model="form.professional_title"
              :placeholder="t('user.professionalTitle')"
              maxlength="100"
            />
          </el-form-item>

          <el-form-item :label="t('user.licenseNumber')" prop="license_number">
            <el-input
              v-model="form.license_number"
              :placeholder="t('user.licenseNumber')"
              maxlength="50"
            />
          </el-form-item>

          <el-form-item :label="t('user.specialties')" prop="specialties">
            <el-select
              v-model="form.specialties"
              :placeholder="t('user.specialties')"
              multiple
              filterable
              allow-create
              style="width: 100%"
            >
            </el-select>
          </el-form-item>
        </el-card>

        <el-card class="form-section">
          <template #header>
            <div class="card-header">
              <el-icon><Key /></el-icon>
              <span>{{ t('user.roles') }}</span>
            </div>
          </template>

          <el-form-item :label="t('user.roles')" prop="role_ids">
            <el-select
              v-model="form.role_ids"
              :placeholder="t('user.roles')"
              multiple
              style="width: 100%"
            >
              <el-option
                v-for="role in roles"
                :key="role.id"
                :label="role.name"
                :value="role.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item :label="t('user.accountExpiry')" prop="account_expiry">
            <el-date-picker
              v-model="accountExpiryDate"
              type="datetime"
              :placeholder="t('user.accountExpiry')"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="X"
              @change="handleAccountExpiryChange"
            />
          </el-form-item>
        </el-card>

        <el-form-item class="form-actions">
          <el-button @click="handleBack">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ t('common.save') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, OfficeBuilding, Key } from '@element-plus/icons-vue'
import { getUserDetail, updateUser } from '@/api/user'
import { getOrganizationList } from '@/api/organization'
import { getRoleList } from '@/api/role'
import type { UpdateUserRequest, UserProfile } from '@/types/user'
import type { Organization } from '@/types/organization'
import type { RoleDefinition } from '@/types/role'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)
const userDetail = ref<UserProfile | null>(null)
const organizations = ref<Organization[]>([])
const roles = ref<RoleDefinition[]>([])
const accountExpiryDate = ref<number>()

const form = reactive<UpdateUserRequest>({
  email: '',
  phone: '',
  first_name: '',
  last_name: '',
  real_name: '',
  professional_title: '',
  license_number: '',
  specialties: [],
  employee_id: '',
  account_expiry: undefined,
  gender: 0,
  role_ids: [],
  organization_id: ''
})

const rules: FormRules<UpdateUserRequest> = {
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

onMounted(() => {
  fetchUserDetail()
  fetchOrganizations()
  fetchRoles()
})

async function fetchUserDetail() {
  loading.value = true
  try {
    const userId = route.params.id as string
    const response = await getUserDetail(userId)
    userDetail.value = response.user

    // 填充表单
    form.email = response.user.email || ''
    form.phone = response.user.phone || ''
    form.first_name = response.user.first_name || ''
    form.last_name = response.user.last_name || ''
    form.real_name = response.user.real_name || ''
    form.professional_title = response.user.professional_title || ''
    form.license_number = response.user.license_number || ''
    form.specialties = response.user.specialties || []
    form.employee_id = response.user.employee_id || ''
    form.account_expiry = response.user.account_expiry
    form.gender = response.user.gender || 0
    form.role_ids = response.user.role_ids || []
    form.organization_id = response.user.organization?.id || ''

    if (response.user.account_expiry) {
      accountExpiryDate.value = response.user.account_expiry
    }
  } catch (error: any) {
    console.error('Failed to fetch user detail:', error)
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

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

function handleAccountExpiryChange(value: number) {
  form.account_expiry = value
}

async function handleSubmit() {
  if (!formRef.value || !userDetail.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    await updateUser(userDetail.value.id, form)
    ElMessage.success(t('common.operationSuccess'))
    router.push(`/users/${userDetail.value.id}`)
  } catch (error: any) {
    console.error('Failed to update user:', error)
    if (error.message) {
      ElMessage.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}

function getStatusKey(status?: number): string {
  if (!status) return 'unknown'
  const statusMap: Record<number, string> = {
    1: 'active',
    2: 'inactive',
    3: 'suspended',
    4: 'locked'
  }
  return statusMap[status] || 'unknown'
}

function getStatusType(status?: number): string {
  if (!status) return 'info'
  const typeMap: Record<number, string> = {
    1: 'success',
    2: 'info',
    3: 'warning',
    4: 'danger'
  }
  return typeMap[status] || 'info'
}
</script>

<style scoped lang="scss">
.user-edit {
  .page-title {
    font-size: 18px;
    font-weight: 600;
    color: #D4AF37;
    font-family: 'Cinzel', serif;
  }

  .form-content {
    margin-top: 20px;

    .form-section {
      background: linear-gradient(145deg, rgba(30, 32, 36, 0.9), rgba(20, 20, 22, 0.95));
      border: 1px solid rgba(255, 255, 255, 0.05);
      border-radius: 20px;
      box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
      margin-bottom: 20px;

      .card-header {
        display: flex;
        align-items: center;
        gap: 10px;
        color: #D4AF37;
        font-family: 'Cinzel', serif;
        font-size: 16px;
        font-weight: 600;
      }

      :deep(.el-card__header) {
        border-bottom: 1px solid rgba(212, 175, 55, 0.2);
      }

      :deep(.el-form-item__label) {
        color: #8B9bb4;
      }
    }

    .form-actions {
      margin-top: 30px;
      text-align: center;

      button {
        margin: 0 10px;
        min-width: 120px;
      }
    }
  }
}
</style>
