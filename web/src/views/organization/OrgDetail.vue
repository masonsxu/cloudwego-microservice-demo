<template>
  <div class="org-detail">
    <el-page-header @back="handleBack" :title="t('common.back')">
      <template #content>
        <span class="page-title">
          {{ orgDetail?.name || t('organization.orgDetail') }}
        </span>
      </template>
      <template #extra>
        <el-button-group>
          <el-button @click="handleEdit">
            <el-icon><Edit /></el-icon>
            {{ t('organization.editOrganization') }}
          </el-button>
          <el-button type="danger" @click="handleDelete">
            <el-icon><Delete /></el-icon>
            {{ t('organization.deleteOrganization') }}
          </el-button>
        </el-button-group>
      </template>
    </el-page-header>

    <div class="detail-content" v-loading="loading">
      <el-row :gutter="20">
        <!-- 基本信息 -->
        <el-col :xs="24" :lg="8">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <el-icon><OfficeBuilding /></el-icon>
                <span>{{ t('organization.basicInfo') }}</span>
              </div>
            </template>
            <div class="org-avatar">
              <el-avatar :size="100" :src="orgDetail?.logo">
                <el-icon><OfficeBuilding /></el-icon>
              </el-avatar>
            </div>
            <el-descriptions :column="1" class="org-info">
              <el-descriptions-item :label="t('organization.name')">
                {{ orgDetail?.name }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('organization.code')">
                {{ orgDetail?.code || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('organization.parentOrganization')">
                <el-tag v-if="!orgDetail?.parent_id" type="info" size="small">
                  {{ t('organization.root') }}
                </el-tag>
                <router-link
                  v-else
                  :to="`/organizations/${orgDetail.parent_id}`"
                  class="link"
                >
                  {{ orgDetail?.parent?.name }}
                </router-link>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- 统计信息 -->
          <el-card class="stat-card">
            <template #header>
              <div class="card-header">
                <el-icon><DataAnalysis /></el-icon>
                <span>{{ t('organization.statistics') }}</span>
              </div>
            </template>
            <el-row :gutter="10">
              <el-col :span="12">
                <div class="stat-item">
                  <div class="stat-value">{{ orgDetail?.member_count || 0 }}</div>
                  <div class="stat-label">{{ t('organization.memberCount') }}</div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="stat-item">
                  <div class="stat-value">{{ orgDetail?.department_count || 0 }}</div>
                  <div class="stat-label">{{ t('organization.departmentCount') }}</div>
                </div>
              </el-col>
            </el-row>
          </el-card>
        </el-col>

        <!-- 详细信息 -->
        <el-col :xs="24" :lg="16">
          <el-card class="detail-card">
            <template #header>
              <div class="card-header">
                <el-icon><InfoFilled /></el-icon>
                <span>{{ t('organization.detailInfo') }}</span>
              </div>
            </template>
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="t('organization.facilityType')">
                {{ orgDetail?.facility_type || '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('organization.accreditationStatus')">
                <el-tag v-if="orgDetail?.accreditation_status" type="success" size="small">
                  {{ orgDetail.accreditation_status }}
                </el-tag>
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('organization.provinceCity')" :span="2">
                <el-tag
                  v-for="city in orgDetail?.province_city"
                  :key="city"
                  class="mr-1"
                >
                  {{ city }}
                </el-tag>
                <span v-if="!orgDetail?.province_city || orgDetail.province_city.length === 0">-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.createTime')">
                {{ formatDateTime(orgDetail?.created_at) }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.updateTime')">
                {{ formatDateTime(orgDetail?.updated_at) }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- 子组织列表 -->
          <el-card class="children-card">
            <template #header>
              <div class="card-header">
                <el-icon><Files /></el-icon>
                <span>{{ t('organization.childOrganizations') }}</span>
              </div>
            </template>
            <el-empty
              v-if="!orgDetail?.children || orgDetail.children.length === 0"
              :description="t('organization.noChildOrganizations')"
            />
            <div v-else class="children-list">
              <div
                v-for="child in orgDetail?.children"
                :key="child.id"
                class="child-item"
              >
                <div class="child-info">
                  <el-icon class="child-icon"><OfficeBuilding /></el-icon>
                  <div class="child-details">
                    <div class="child-name">{{ child.name }}</div>
                    <div class="child-meta">
                      <span>{{ t('organization.memberCount') }}: {{ child.member_count || 0 }}</span>
                      <span>{{ t('organization.departmentCount') }}: {{ child.department_count || 0 }}</span>
                    </div>
                  </div>
                </div>
                <el-button
                  link
                  type="primary"
                  @click="router.push(`/organizations/${child.id}`)"
                >
                  {{ t('common.view') }}
                </el-button>
              </div>
            </div>
          </el-card>

          <!-- Logo 管理 -->
          <el-card class="logo-card">
            <template #header>
              <div class="card-header">
                <el-icon><Picture /></el-icon>
                <span>{{ t('organization.logoManagement') }}</span>
              </div>
            </template>
            <div class="logo-section">
              <OrgLogoUpload
                :organization-id="orgId"
                :logo-url="orgDetail?.logo"
                @update:logo-url="handleLogoUpdate"
                @upload-success="handleLogoUploadSuccess"
                @remove="handleLogoRemove"
              />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 编辑组织对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="t('organization.editOrganization')"
      width="600px"
      @closed="handleEditDialogClosed"
    >
      <el-form
        ref="editFormRef"
        :model="editFormData"
        :rules="editFormRules"
        label-width="140px"
      >
        <el-form-item :label="t('organization.name')" prop="name">
          <el-input
            v-model="editFormData.name"
            :placeholder="t('organization.namePlaceholder')"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('organization.code')" prop="code">
          <el-input
            v-model="editFormData.code"
            :placeholder="t('organization.codePlaceholder')"
            maxlength="50"
          />
        </el-form-item>
        <el-form-item :label="t('organization.parentOrganization')" prop="parent_id">
          <el-tree-select
            v-model="editFormData.parent_id"
            :data="organizationTree"
            :props="{ label: 'name', value: 'id' }"
            :placeholder="t('organization.parentOrganizationPlaceholder')"
            clearable
            check-strictly
          />
        </el-form-item>
        <el-form-item :label="t('organization.facilityType')" prop="facility_type">
          <el-input
            v-model="editFormData.facility_type"
            :placeholder="t('organization.facilityTypePlaceholder')"
            maxlength="100"
          />
        </el-form-item>
        <el-form-item :label="t('organization.accreditationStatus')" prop="accreditation_status">
          <el-input
            v-model="editFormData.accreditation_status"
            :placeholder="t('organization.accreditationStatusPlaceholder')"
            maxlength="100"
          />
        </el-form-item>
        <el-form-item :label="t('organization.provinceCity')" prop="province_city">
          <el-cascader
            v-model="editFormData.province_city"
            :options="provinceCityOptions"
            :props="{
              value: 'label',
              label: 'label',
              children: 'children',
              multiple: true
            }"
            :placeholder="t('organization.provinceCityPlaceholder')"
            clearable
            collapse-tags
            collapse-tags-tooltip
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="updating">
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Edit,
  Delete,
  OfficeBuilding,
  InfoFilled,
  DataAnalysis,
  Files,
  Picture
} from '@element-plus/icons-vue'
import { organizationApi } from '@/api/organization'
import { formatOptionalTimestamp } from '@/utils/date'
import OrgLogoUpload from '@/components/OrgLogoUpload.vue'
import type { OrganizationDTO, UpdateOrganizationRequestDTO } from '@/api/organization'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const loading = ref(false)
const updating = ref(false)
const orgDetail = ref<OrganizationDTO | null>(null)
const organizationTree = ref<any[]>([])

const editDialogVisible = ref(false)
const editFormRef = ref<FormInstance>()

const editFormData = reactive<UpdateOrganizationRequestDTO & { code?: string }>({
  name: '',
  parent_id: '',
  facility_type: '',
  accreditation_status: '',
  province_city: []
})

const editFormRules: FormRules = {
  name: [
    { required: true, message: t('organization.nameRequired'), trigger: 'blur' },
    { min: 2, max: 100, message: t('organization.nameLengthLimit'), trigger: 'blur' }
  ]
}

const provinceCityOptions = [
  {
    label: '北京市',
    children: [
      { label: '北京市' }
    ]
  },
  {
    label: '上海市',
    children: [
      { label: '上海市' }
    ]
  },
  {
    label: '广东省',
    children: [
      { label: '广州市' },
      { label: '深圳市' },
      { label: '珠海市' },
      { label: '东莞市' }
    ]
  },
  {
    label: '浙江省',
    children: [
      { label: '杭州市' },
      { label: '宁波市' },
      { label: '温州市' }
    ]
  },
  {
    label: '江苏省',
    children: [
      { label: '南京市' },
      { label: '苏州市' },
      { label: '无锡市' }
    ]
  }
]

const orgId = route.params.id as string

onMounted(() => {
  fetchData()
  fetchOrganizationTree()
})

async function fetchOrganizationTree() {
  try {
    const response = await organizationApi.listOrganizations({ page: 1, pageSize: 1000 })
    const orgs = response.organizations || []
    organizationTree.value = buildTree(orgs as OrganizationDTO[])
  } catch (error) {
    console.error('Failed to fetch organization tree:', error)
  }
}

function buildTree(orgs: OrganizationDTO[], parentId: string | null = null): any[] {
  return orgs
    .filter(org => org.parent_id === parentId || (parentId === null && !org.parent_id))
    .map(org => ({
      ...org,
      children: buildTree(orgs, org.id || null)
    }))
}

async function fetchData() {
  loading.value = true
  try {
    const response = await organizationApi.getOrganization(orgId)
    orgDetail.value = response.organization || null
  } catch (error: any) {
    console.error('Failed to fetch organization detail:', error)
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

function handleBack() {
  router.back()
}

function handleEdit() {
  if (!orgDetail.value) return

  Object.assign(editFormData, {
    name: orgDetail.value.name,
    code: orgDetail.value.code,
    parent_id: orgDetail.value.parent_id || '',
    facility_type: orgDetail.value.facility_type || '',
    accreditation_status: orgDetail.value.accreditation_status || '',
    province_city: orgDetail.value.province_city || []
  })

  editDialogVisible.value = true
}

async function handleUpdate() {
  if (!editFormRef.value) return

  await editFormRef.value.validate(async (valid) => {
    if (!valid) return

    updating.value = true
    try {
      const updateData: UpdateOrganizationRequestDTO = {
        name: editFormData.name,
        parent_id: editFormData.parent_id || undefined,
        facility_type: editFormData.facility_type,
        accreditation_status: editFormData.accreditation_status,
        province_city: editFormData.province_city
      }

      await organizationApi.updateOrganization(orgId, updateData)
      ElMessage.success(t('common.updateSuccess'))
      editDialogVisible.value = false
      fetchData()
    } catch (error: any) {
      console.error('Failed to update organization:', error)
      ElMessage.error(error.message || t('common.operationFailed'))
    } finally {
      updating.value = false
    }
  })
}

function handleEditDialogClosed() {
  editFormRef.value?.resetFields()
}

async function handleDelete() {
  try {
    await ElMessageBox.confirm(
      `${t('organization.name')}: ${orgDetail.value?.name}`,
      t('common.deleteConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        distinguishCancelAndClose: true
      }
    )

    await organizationApi.deleteOrganization(orgId)
    ElMessage.success(t('common.deleteSuccess'))
    router.push('/organizations')
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('Failed to delete organization:', error)
      ElMessage.error(error.message || t('common.operationFailed'))
    }
  }
}

function formatDateTime(timestamp?: number): string {
  return formatOptionalTimestamp(timestamp, 'YYYY-MM-DD HH:mm:ss')
}

function handleLogoUpdate(logoUrl: string) {
  if (orgDetail.value) {
    orgDetail.value.logo = logoUrl
  }
}

function handleLogoUploadSuccess(_logoId: string, _logoUrl: string) {
  ElMessage.success(t('organization.logoUploadSuccess'))
  fetchData()
}

function handleLogoRemove() {
  if (orgDetail.value) {
    orgDetail.value.logo = ''
  }
  ElMessage.success(t('organization.logoRemoveSuccess'))
}
</script>

<style scoped lang="scss">
.org-detail {
  .page-title {
    color: #D4AF37;
    font-family: 'Cinzel', serif;
    font-size: 20px;
    font-weight: 600;
  }

  .detail-content {
    margin-top: 20px;

    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #D4AF37;
      font-weight: 600;
    }

    .info-card {
      margin-bottom: 20px;

      .org-avatar {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin-bottom: 20px;
      }

      .org-info {
        margin-top: 20px;
      }
    }

    .stat-card {
      margin-bottom: 20px;

      .stat-item {
        text-align: center;
        padding: 15px;
        background-color: rgba(44, 46, 51, 0.3);
        border-radius: 8px;
        border: 1px solid rgba(212, 175, 55, 0.2);

        .stat-value {
          font-size: 28px;
          font-weight: 600;
          color: #D4AF37;
          margin-bottom: 5px;
        }

        .stat-label {
          font-size: 12px;
          color: #8B9bb4;
        }
      }
    }

    .detail-card {
      margin-bottom: 20px;
    }

    .children-card {
      margin-bottom: 20px;

      .children-list {
        .child-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 15px;
          background-color: rgba(44, 46, 51, 0.3);
          border-radius: 8px;
          border: 1px solid rgba(212, 175, 55, 0.1);
          margin-bottom: 10px;

          &:last-child {
            margin-bottom: 0;
          }

          .child-info {
            display: flex;
            align-items: center;
            gap: 15px;

            .child-icon {
              font-size: 32px;
              color: #D4AF37;
            }

            .child-details {
              .child-name {
                font-size: 16px;
                font-weight: 500;
                margin-bottom: 5px;
              }

              .child-meta {
                display: flex;
                gap: 15px;
                font-size: 12px;
                color: #8B9bb4;
              }
            }
          }
        }
      }
    }

    .logo-card {
      .logo-section {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 250px;

        .current-logo {
          display: flex;
          justify-content: center;
          align-items: center;
        }
      }
    }

    .mr-1 {
      margin-right: 4px;
    }

    .link {
      color: #D4AF37;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}
</style>
