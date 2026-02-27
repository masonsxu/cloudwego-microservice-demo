<template>
  <div class="org-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ t('organization.title') }}</h1>
        <p class="page-subtitle">管理医疗机构及其层级关系</p>
      </div>
      <button class="create-btn shimmer-btn" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        {{ t('organization.createOrganization') }}
      </button>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon">
          <el-icon size="22"><OfficeBuilding /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ pagination.total }}</span>
          <span class="stat-label">{{ t('organization.title') }}</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon">
          <el-icon size="22"><Connection /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ tableData.filter(o => !o.parent_id).length }}</span>
          <span class="stat-label">根节点组织</span>
        </div>
      </div>
      <div class="stat-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="stat-icon">
          <el-icon size="22"><User /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ tableData.reduce((sum, o) => sum + (o.member_count || 0), 0) }}</span>
          <span class="stat-label">总成员数</span>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选栏 -->
    <div class="search-section">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchForm.search"
            :placeholder="t('organization.name') + ' / ' + t('organization.code')"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            style="width: 280px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-tree-select
            v-model="searchForm.parent_id"
            :data="organizationTree"
            :props="{ label: 'name', value: 'id' }"
            :placeholder="t('organization.parentOrganization')"
            clearable
            check-strictly
            @change="handleSearch"
            style="width: 220px"
          />
        </el-form-item>
        <el-form-item>
          <el-button @click="handleReset" class="reset-btn">
            <el-icon><Refresh /></el-icon>
            {{ t('common.reset') }}
          </el-button>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ t('common.search') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 数据表格卡片 -->
    <div class="table-card spotlight-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <el-table
        :data="tableData"
        v-loading="loading"
        row-key="id"
        class="modern-table"
        :style="{ width: '100%' }"
      >
        <el-table-column prop="name" :label="t('organization.name')" min-width="220">
          <template #default="{ row }">
            <div class="org-name-cell">
              <div class="org-icon">
                <el-icon size="14"><OfficeBuilding /></el-icon>
              </div>
              <span class="org-name">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="code" :label="t('organization.code')" width="140">
          <template #default="{ row }">
            <code class="code-badge">{{ row.code || '-' }}</code>
          </template>
        </el-table-column>
        <el-table-column :label="t('organization.parentOrganization')" min-width="180">
          <template #default="{ row }">
            <span v-if="row.parent" class="text-sub">{{ row.parent.name }}</span>
            <el-tag v-else size="small" class="root-tag">{{ t('organization.root') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="facility_type" :label="t('organization.facilityType')" width="140">
          <template #default="{ row }">
            <span class="text-sub">{{ row.facility_type || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('organization.memberCount')" width="110" align="center">
          <template #default="{ row }">
            <span class="count-badge">{{ row.member_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('organization.departmentCount')" width="110" align="center">
          <template #default="{ row }">
            <span class="count-badge">{{ row.department_count || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="180" align="right" fixed="right">
          <template #default="{ row }">
            <div class="action-group">
              <button class="action-btn view-btn" @click="handleView(row)" :title="t('common.view')">
                <el-icon><View /></el-icon>
              </button>
              <button class="action-btn edit-btn" @click="handleEdit(row)" :title="t('common.edit')">
                <el-icon><Edit /></el-icon>
              </button>
              <button class="action-btn delete-btn" @click="handleDelete(row)" :title="t('common.delete')">
                <el-icon><Delete /></el-icon>
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          :background="true"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'create' ? t('organization.createOrganization') : t('organization.editOrganization')"
      width="560px"
      @closed="handleDialogClosed"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="130px">
        <el-form-item :label="t('organization.name')" prop="name">
          <el-input v-model="formData.name" :placeholder="t('organization.namePlaceholder')" maxlength="100" show-word-limit />
        </el-form-item>
        <el-form-item :label="t('organization.code')" prop="code">
          <el-input v-model="formData.code" :placeholder="t('organization.codePlaceholder')" maxlength="50" />
        </el-form-item>
        <el-form-item :label="t('organization.parentOrganization')" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="organizationTree"
            :props="{ label: 'name', value: 'id' }"
            :placeholder="t('organization.parentOrganizationPlaceholder')"
            clearable check-strictly style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="t('organization.facilityType')" prop="facility_type">
          <el-input v-model="formData.facility_type" :placeholder="t('organization.facilityTypePlaceholder')" maxlength="100" />
        </el-form-item>
        <el-form-item :label="t('organization.accreditationStatus')" prop="accreditation_status">
          <el-input v-model="formData.accreditation_status" :placeholder="t('organization.accreditationStatusPlaceholder')" maxlength="100" />
        </el-form-item>
        <el-form-item :label="t('organization.provinceCity')" prop="province_city">
          <el-cascader
            v-model="formData.province_city"
            :options="provinceCityOptions"
            :props="{ value: 'label', label: 'label', children: 'children', multiple: true }"
            :placeholder="t('organization.provinceCityPlaceholder')"
            clearable collapse-tags collapse-tags-tooltip style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus, View, Edit, Delete, OfficeBuilding, Connection, User } from '@element-plus/icons-vue'
import { organizationApi } from '@/api/organization'
import type { OrganizationDTO, CreateOrganizationRequestDTO } from '@/api/organization'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<OrganizationDTO[]>([])
const organizationTree = ref<any[]>([])

const searchForm = reactive({ search: '', parent_id: '' })
const pagination = reactive({ page: 1, limit: 20, total: 0 })
const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const currentOrg = ref<OrganizationDTO | null>(null)

const formData = reactive<CreateOrganizationRequestDTO & { code?: string }>({
  name: '', parent_id: '', facility_type: '', accreditation_status: '', province_city: []
})

const formRules: FormRules = {
  name: [
    { required: true, message: t('organization.nameRequired'), trigger: 'blur' },
    { min: 2, max: 100, message: t('organization.nameLengthLimit'), trigger: 'blur' }
  ]
}

const provinceCityOptions = [
  { label: '北京市', children: [{ label: '北京市' }] },
  { label: '上海市', children: [{ label: '上海市' }] },
  { label: '广东省', children: [{ label: '广州市' }, { label: '深圳市' }, { label: '珠海市' }, { label: '东莞市' }] },
  { label: '浙江省', children: [{ label: '杭州市' }, { label: '宁波市' }, { label: '温州市' }] },
  { label: '江苏省', children: [{ label: '南京市' }, { label: '苏州市' }, { label: '无锡市' }] }
]

onMounted(() => {
  fetchData()
  fetchOrganizationTree()
})

function onMouseMove(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  el.style.setProperty('--mouse-x', `${e.clientX - rect.left}px`)
  el.style.setProperty('--mouse-y', `${e.clientY - rect.top}px`)
}

function onMouseLeave(e: MouseEvent) {
  const el = e.currentTarget as HTMLElement
  el.style.removeProperty('--mouse-x')
  el.style.removeProperty('--mouse-y')
}

function buildTree(orgs: OrganizationDTO[], parentId: string | null = null): any[] {
  return orgs.filter(o => o.parent_id === parentId).map(o => ({ ...o, children: buildTree(orgs, o.id || null) }))
}

async function fetchData() {
  loading.value = true
  try {
    const response = await organizationApi.listOrganizations({
      parentId: searchForm.parent_id || undefined,
      page: pagination.page,
      pageSize: pagination.limit
    })
    tableData.value = response.organizations || []
    pagination.total = response.page?.total || 0
  } catch (error: any) {
    ElMessage.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
  }
}

async function fetchOrganizationTree() {
  try {
    const response = await organizationApi.listOrganizations({ page: 1, pageSize: 1000 })
    organizationTree.value = buildTree(response.organizations || [])
  } catch {}
}

function handleSearch() { pagination.page = 1; fetchData() }
function handleReset() { searchForm.search = ''; searchForm.parent_id = ''; pagination.page = 1; fetchData() }
function handleCreate() { dialogMode.value = 'create'; currentOrg.value = null; resetForm(); dialogVisible.value = true }
function handleView(row: OrganizationDTO) { router.push(`/organizations/${row.id}`) }
function handleEdit(row: OrganizationDTO) {
  dialogMode.value = 'edit'; currentOrg.value = row
  Object.assign(formData, { name: row.name, code: row.code, parent_id: row.parent_id || '', facility_type: row.facility_type || '', accreditation_status: row.accreditation_status || '', province_city: row.province_city || [] })
  dialogVisible.value = true
}

async function handleDelete(row: OrganizationDTO) {
  try {
    await ElMessageBox.confirm(`${t('organization.name')}: ${row.name}`, t('common.deleteConfirm'), {
      confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning', distinguishCancelAndClose: true
    })
    if (row.id) {
      await organizationApi.deleteOrganization(row.id)
      ElMessage.success(t('common.deleteSuccess'))
      fetchData(); fetchOrganizationTree()
    }
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') ElMessage.error(error.message || t('common.operationFailed'))
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (dialogMode.value === 'create') {
        await organizationApi.createOrganization(formData as CreateOrganizationRequestDTO)
        ElMessage.success(t('common.createSuccess'))
      } else if (currentOrg.value?.id) {
        await organizationApi.updateOrganization(currentOrg.value.id, { name: formData.name, parent_id: formData.parent_id || undefined, facility_type: formData.facility_type, accreditation_status: formData.accreditation_status, province_city: formData.province_city })
        ElMessage.success(t('common.updateSuccess'))
      }
      dialogVisible.value = false; fetchData(); fetchOrganizationTree()
    } catch (error: any) {
      ElMessage.error(error.message || t('common.operationFailed'))
    } finally {
      submitting.value = false
    }
  })
}

function handleDialogClosed() { formRef.value?.resetFields(); resetForm() }
function resetForm() { Object.assign(formData, { name: '', code: '', parent_id: '', facility_type: '', accreditation_status: '', province_city: [] }) }
function handleSizeChange(size: number) { pagination.limit = size; pagination.page = 1; fetchData() }
function handlePageChange(page: number) { pagination.page = page; fetchData() }
</script>

<style scoped lang="scss">
.org-list {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 28px;

    .header-left {
      .page-title {
        font-size: 26px;
        font-weight: 700;
        font-family: 'Cinzel', serif;
        background: linear-gradient(to right, #D4AF37, #F2D288, #D4AF37);
        background-clip: text;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin: 0 0 6px;
        line-height: 1.2;
      }

      .page-subtitle {
        color: var(--c-text-sub);
        font-size: 13px;
        margin: 0;
      }
    }

    .create-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 10px 20px;
      background: linear-gradient(135deg, #D4AF37 0%, #C4A033 100%);
      color: #000;
      font-weight: 600;
      font-size: 14px;
      border: none;
      border-radius: 10px;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 0 24px rgba(212, 175, 55, 0.4);
        transform: translateY(-1px);
      }
    }
  }

  .stats-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 24px;

    .stat-card {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 20px 24px;
      background: var(--bg-card);
      border: 1px solid var(--c-border-accent);
      border-radius: 14px;
      box-shadow: var(--shadow-card);
      cursor: default;

      .stat-icon {
        width: 44px;
        height: 44px;
        border-radius: 10px;
        background: rgba(212, 175, 55, 0.12);
        border: 1px solid var(--c-border-accent);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--c-accent);
        flex-shrink: 0;
      }

      .stat-info {
        display: flex;
        flex-direction: column;
        gap: 2px;

        .stat-value {
          font-size: 24px;
          font-weight: 700;
          color: var(--c-text-main);
          font-family: 'JetBrains Mono', monospace;
          line-height: 1;
        }

        .stat-label {
          font-size: 12px;
          color: var(--c-text-sub);
        }
      }
    }
  }

  .search-section {
    background: var(--bg-card);
    border: 1px solid var(--c-border-accent);
    border-radius: 14px;
    padding: 16px 20px;
    margin-bottom: 20px;

    .search-form {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 0;
      margin: 0;

      :deep(.el-form-item) {
        margin-bottom: 0;
        margin-right: 12px;
      }
    }
  }

  .table-card {
    background: var(--bg-card);
    border: 1px solid var(--c-border-accent);
    border-radius: 14px;
    overflow: hidden;
    box-shadow: var(--shadow-card);

    .modern-table {
      :deep(.el-table__header-wrapper) {
        th.el-table__cell {
          padding: 14px 12px;
          font-size: 12px;
          letter-spacing: 0.05em;
          text-transform: uppercase;
        }
      }

      :deep(.el-table__body-wrapper) {
        td.el-table__cell {
          padding: 14px 12px;
        }
      }
    }

    .org-name-cell {
      display: flex;
      align-items: center;
      gap: 10px;

      .org-icon {
        width: 28px;
        height: 28px;
        border-radius: 6px;
        background: rgba(212, 175, 55, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--c-accent);
        flex-shrink: 0;
      }

      .org-name {
        font-weight: 500;
        color: var(--c-text-main);
      }
    }

    .code-badge {
      display: inline-block;
      padding: 2px 8px;
      background: rgba(212, 175, 55, 0.08);
      border: 1px solid var(--c-border-accent);
      border-radius: 6px;
      font-family: 'JetBrains Mono', monospace;
      font-size: 12px;
      color: var(--c-accent);
    }

    .root-tag {
      background: rgba(103, 194, 58, 0.1);
      border-color: rgba(103, 194, 58, 0.3);
      color: #67C23A;
    }

    .text-sub {
      color: var(--c-text-sub);
      font-size: 13px;
    }

    .count-badge {
      display: inline-block;
      padding: 2px 10px;
      background: rgba(212, 175, 55, 0.08);
      border-radius: 9999px;
      font-size: 13px;
      font-weight: 600;
      color: var(--c-text-main);
    }

    .action-group {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      gap: 6px;

      .action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 30px;
        height: 30px;
        border-radius: 7px;
        border: 1px solid var(--c-border);
        background: var(--bg-input);
        cursor: pointer;
        transition: all 0.2s ease;
        color: var(--c-text-sub);
        font-size: 14px;

        &:hover {
          transform: translateY(-1px);
        }

        &.view-btn:hover {
          color: var(--c-accent);
          border-color: var(--c-border-accent);
          background: rgba(212, 175, 55, 0.08);
        }

        &.edit-btn:hover {
          color: #E6A23C;
          border-color: rgba(230, 162, 60, 0.4);
          background: rgba(230, 162, 60, 0.08);
        }

        &.delete-btn:hover {
          color: #F56C6C;
          border-color: rgba(245, 108, 108, 0.4);
          background: rgba(245, 108, 108, 0.08);
        }
      }
    }

    .pagination-wrapper {
      display: flex;
      justify-content: flex-end;
      padding: 16px 20px;
      border-top: 1px solid var(--c-border);
    }
  }
}

/* spotlight 效果 */
.spotlight-card {
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(500px circle at var(--mouse-x, -100%) var(--mouse-y, -100%),
      rgba(212, 175, 55, 0.05), transparent 50%);
    pointer-events: none;
    z-index: 1;
    transition: opacity 0.3s;
  }
}
</style>
