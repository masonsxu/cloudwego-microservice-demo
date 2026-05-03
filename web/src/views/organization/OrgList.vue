<template>
  <div class="space-y-6">
    <!-- 页头 -->
    <header class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-[22px] font-semibold leading-tight tracking-[-0.012em] text-[color:var(--color-ink-strong)]">
          {{ t('organization.title') }}
        </h1>
        <p class="mt-1 text-[13px] text-[color:var(--color-ink-muted)]">
          {{ t('organization.manageOrgs') }}
        </p>
      </div>
      <Button @click="handleCreate">
        <Plus class="h-4 w-4" />
        {{ t('organization.createOrganization') }}
      </Button>
    </header>

    <!-- 工具条 -->
    <div class="flex flex-wrap items-center gap-2">
      <div class="relative flex-1 min-w-[260px] max-w-[360px]">
        <Search class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[color:var(--color-ink-subtle)]" />
        <Input
          v-model="searchForm.search"
          :placeholder="t('organization.searchPlaceholder')"
          class="pl-9 pr-9"
          @keyup.enter="handleSearch"
        />
        <button
          v-if="searchForm.search"
          class="absolute right-2.5 top-1/2 flex h-5 w-5 -translate-y-1/2 items-center justify-center rounded-xs text-[color:var(--color-ink-subtle)] transition-colors hover:text-[color:var(--color-ink)]"
          @click="searchForm.search = ''; handleSearch()"
        >
          <X class="h-3 w-3" />
        </button>
      </div>

      <Select v-model="searchForm.parent_id" @update:modelValue="handleSearch">
        <SelectTrigger class="w-[220px]">
          <SelectValue :placeholder="t('organization.parentOrganization')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">{{ t('common.all') }}</SelectItem>
          <SelectItem v-for="org in flattenOrgTree(organizationTree)" :key="org.id" :value="org.id">
            {{ org.name }}
          </SelectItem>
        </SelectContent>
      </Select>

      <Button variant="ghost" size="sm" @click="handleReset">
        <RefreshCw class="h-3.5 w-3.5" />
        {{ t('common.reset') }}
      </Button>
    </div>

    <!-- 表格卡片 -->
    <div class="overflow-hidden rounded-md border border-subtle bg-canvas">
      <ListPageSkeleton v-if="initialLoading" :columns="6" :rows="8" />
      <template v-else>
        <Table>
          <TableHeader>
            <TableRow class="border-subtle hover:bg-transparent">
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)]">{{ t('organization.name') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px]">{{ t('organization.code') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)]">{{ t('organization.parentOrganization') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px]">{{ t('organization.facilityType') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[100px] text-right">{{ t('organization.memberCount') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[100px] text-right">{{ t('organization.departmentCount') }}</TableHead>
              <TableHead class="h-10 text-[12px] font-medium uppercase tracking-[0.04em] text-[color:var(--color-ink-muted)] w-[140px] text-right">{{ t('common.actions') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="row in tableData"
              :key="row.id"
              class="h-11 border-subtle transition-colors duration-[var(--duration-fast)] hover:bg-sunken"
            >
              <TableCell>
                <div class="flex items-center gap-2">
                  <span class="flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-sm bg-[color:var(--color-primary-soft)] text-[color:var(--color-primary-active)]">
                    <Building2 class="h-3 w-3" />
                  </span>
                  <span class="font-medium text-ink">{{ row.name }}</span>
                </div>
              </TableCell>
              <TableCell>
                <code v-if="row.code" class="font-mono text-[12px] text-[color:var(--color-ink-muted)]">{{ row.code }}</code>
                <span v-else class="text-[color:var(--color-ink-subtle)]">—</span>
              </TableCell>
              <TableCell>
                <span v-if="row.parent" class="text-[color:var(--color-ink-muted)]">{{ row.parent.name }}</span>
                <Badge v-else variant="success">{{ t('organization.root') }}</Badge>
              </TableCell>
              <TableCell class="text-[color:var(--color-ink-muted)]">
                {{ row.facility_type || '—' }}
              </TableCell>
              <TableCell class="tabular text-right text-ink">{{ row.member_count || 0 }}</TableCell>
              <TableCell class="tabular text-right text-ink">{{ row.department_count || 0 }}</TableCell>
              <TableCell class="text-right">
                <div class="inline-flex items-center gap-0.5">
                  <Button variant="ghost" size="icon" :title="t('common.view')" @click="handleView(row)">
                    <Eye class="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" :title="t('common.edit')" @click="handleEdit(row)">
                    <Pencil class="h-3.5 w-3.5" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    :title="t('common.delete')"
                    class="hover:text-[color:var(--color-danger)]"
                    @click="handleDelete(row)"
                  >
                    <Trash2 class="h-3.5 w-3.5" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>

            <TableRow v-if="!loading && tableData.length === 0">
              <TableCell colspan="7" class="h-[200px] text-center">
                <div class="flex flex-col items-center justify-center gap-2 py-8">
                  <Building2 class="h-8 w-8 text-[color:var(--color-ink-subtle)]" />
                  <p class="text-[14px] font-medium text-ink">{{ t('common.noData') }}</p>
                  <p class="text-[12px] text-[color:var(--color-ink-muted)]">{{ t('organization.emptyHint') }}</p>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>

        <!-- 分页 -->
        <div class="flex items-center justify-between border-t border-subtle px-4 py-2.5">
          <span class="tabular text-[13px] text-[color:var(--color-ink-muted)]">
            {{ t('common.total', { total: pagination.total }) }}
          </span>
          <div class="flex items-center gap-2">
            <Select v-model="paginationLimitString">
              <SelectTrigger class="h-8 w-[80px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="10">10</SelectItem>
                <SelectItem value="20">20</SelectItem>
                <SelectItem value="50">50</SelectItem>
                <SelectItem value="100">100</SelectItem>
              </SelectContent>
            </Select>
            <div class="flex items-center gap-1">
              <Button variant="outline" size="icon" :disabled="pagination.page <= 1" @click="handlePageChange(pagination.page - 1)">
                <ChevronLeft class="h-4 w-4" />
              </Button>
              <span class="tabular min-w-[40px] text-center text-[13px] font-medium text-ink">{{ pagination.page }}</span>
              <Button variant="outline" size="icon" :disabled="pagination.page * pagination.limit >= pagination.total" @click="handlePageChange(pagination.page + 1)">
                <ChevronRight class="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- 创建/编辑对话框 -->
    <Dialog v-model:open="dialogVisible">
      <DialogContent class="max-w-[520px]">
        <DialogHeader>
          <DialogTitle>{{ dialogMode === 'create' ? t('organization.createOrganization') : t('organization.editOrganization') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-2">
          <div class="space-y-1.5">
            <Label>{{ t('organization.name') }}</Label>
            <Input v-model="formData.name" :placeholder="t('organization.namePlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-1.5">
            <Label>{{ t('organization.code') }}</Label>
            <Input v-model="formData.code" :placeholder="t('organization.codePlaceholder')" maxlength="50" />
          </div>
          <div class="space-y-1.5">
            <Label>{{ t('organization.parentOrganization') }}</Label>
            <Select v-model="formData.parent_id">
              <SelectTrigger>
                <SelectValue :placeholder="t('organization.parentOrganizationPlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">{{ t('common.none') }}</SelectItem>
                <SelectItem v-for="org in flattenOrgTree(organizationTree)" :key="org.id" :value="org.id">
                  {{ org.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-1.5">
            <Label>{{ t('organization.facilityType') }}</Label>
            <Input v-model="formData.facility_type" :placeholder="t('organization.facilityTypePlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-1.5">
            <Label>{{ t('organization.accreditationStatus') }}</Label>
            <Input v-model="formData.accreditation_status" :placeholder="t('organization.accreditationStatusPlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-1.5">
            <Label>{{ t('organization.provinceCity') }}</Label>
            <div class="flex max-h-40 flex-wrap gap-2 overflow-y-auto rounded-sm border border-default bg-canvas p-2">
              <template v-for="province in provinceCityOptions" :key="province.label">
                <template v-for="city in province.children" :key="city.label">
                  <label class="flex items-center gap-1.5 text-[13px] cursor-pointer">
                    <Checkbox
                      :checked="formDataProvinceCity.includes(city.label)"
                      @update:checked="(checked: boolean) => toggleCitySelection(checked, city.label)"
                    />
                    <span class="text-[color:var(--color-ink-muted)]">{{ province.label }} - {{ city.label }}</span>
                  </label>
                </template>
              </template>
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="dialogVisible = false">{{ t('common.cancel') }}</Button>
          <Button :disabled="submitting" @click="handleSubmit">{{ t('common.confirm') }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Building2, Eye, Pencil, Trash2, ChevronLeft, ChevronRight, X, RefreshCw, Search, Plus } from 'lucide-vue-next'
import { organizationApi } from '@/api/organization'
import type { OrganizationDTO, CreateOrganizationRequestDTO } from '@/api/organization'
import ListPageSkeleton from '@/components/skeleton/ListPageSkeleton.vue'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'

const router = useRouter()
const { t } = useI18n()

const initialLoading = ref(true)
const loading = ref(false)
const submitting = ref(false)
const tableData = ref<OrganizationDTO[]>([])
const organizationTree = ref<any[]>([])

const searchForm = reactive({ search: '', parent_id: 'all' })
const pagination = reactive({ page: 1, limit: 20, total: 0 })
const paginationLimitString = computed({
  get: () => String(pagination.limit),
  set: (val: string) => { pagination.limit = Number(val); fetchData() },
})
const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const currentOrg = ref<OrganizationDTO | null>(null)

const formData = reactive<CreateOrganizationRequestDTO & { code?: string }>({
  name: '', parent_id: 'none', facility_type: '', accreditation_status: '', province_city: []
})

const formDataProvinceCity = computed({
  get: () => formData.province_city || [],
  set: (val: string[]) => { formData.province_city = val },
})

function toggleCitySelection(checked: boolean, city: string) {
  if (checked) {
    formData.province_city = [...(formData.province_city || []), city]
  } else {
    formData.province_city = (formData.province_city || []).filter(c => c !== city)
  }
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

function flattenOrgTree(orgs: any[]): any[] {
  const result: any[] = []
  for (const org of orgs) {
    result.push(org)
    if (org.children?.length) result.push(...flattenOrgTree(org.children))
  }
  return result
}

function buildTree(orgs: OrganizationDTO[], parentId: string | null = null): any[] {
  return orgs.filter(o => o.parent_id === parentId).map(o => ({ ...o, children: buildTree(orgs, o.id || null) }))
}

async function fetchData() {
  loading.value = true
  try {
    const response = await organizationApi.listOrganizations({
      parent_id: searchForm.parent_id !== 'all' ? searchForm.parent_id : undefined,
      page: pagination.page,
      limit: pagination.limit
    })
    tableData.value = response.organizations || []
    pagination.total = response.page?.total || 0
  } catch (error: any) {
    toast.error(error.message || t('common.operationFailed'))
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

async function fetchOrganizationTree() {
  try {
    const response = await organizationApi.listOrganizations({ page: 1, limit: 1000 })
    organizationTree.value = buildTree(response.organizations || [])
  } catch {}
}

function handleSearch() { pagination.page = 1; fetchData() }
function handleReset() { searchForm.search = ''; searchForm.parent_id = 'all'; pagination.page = 1; fetchData() }
function handleCreate() { dialogMode.value = 'create'; currentOrg.value = null; resetForm(); dialogVisible.value = true }
function handleView(row: OrganizationDTO) { router.push(`/organizations/${row.id}`) }
function handleEdit(row: OrganizationDTO) {
  dialogMode.value = 'edit'; currentOrg.value = row
  Object.assign(formData, { name: row.name, code: row.code, parent_id: row.parent_id || 'none', facility_type: row.facility_type || '', accreditation_status: row.accreditation_status || '', province_city: row.province_city || [] })
  dialogVisible.value = true
}

async function handleDelete(row: OrganizationDTO) {
  if (!confirm(`${t('organization.name')}: ${row.name}\n${t('common.deleteConfirm')}`)) return
  try {
    if (row.id) {
      await organizationApi.deleteOrganization(row.id)
      toast.success(t('common.deleteSuccess'))
      fetchData(); fetchOrganizationTree()
    }
  } catch (error: any) {
    toast.error(error.message || t('common.operationFailed'))
  }
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (dialogMode.value === 'create') {
      await organizationApi.createOrganization({
        ...formData,
        parent_id: formData.parent_id !== 'none' ? formData.parent_id : undefined,
      } as CreateOrganizationRequestDTO)
      toast.success(t('common.createSuccess'))
    } else if (currentOrg.value?.id) {
      await organizationApi.updateOrganization(currentOrg.value.id, { name: formData.name, parent_id: formData.parent_id !== 'none' ? formData.parent_id : undefined, facility_type: formData.facility_type, accreditation_status: formData.accreditation_status, province_city: formData.province_city })
      toast.success(t('common.updateSuccess'))
    }
    dialogVisible.value = false; fetchData(); fetchOrganizationTree()
  } catch (error: any) {
    toast.error(error.message || t('common.operationFailed'))
  } finally {
    submitting.value = false
  }
}

function resetForm() { Object.assign(formData, { name: '', code: '', parent_id: 'none', facility_type: '', accreditation_status: '', province_city: [] }) }
function handlePageChange(page: number) { pagination.page = page; fetchData() }
</script>
