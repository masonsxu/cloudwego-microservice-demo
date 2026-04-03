<template>
  <div class="flex flex-col h-[calc(100vh-108px)]">
    <!-- Page Header -->
    <div class="flex justify-between items-end mb-7">
      <div>
        <h1 class="text-[26px] font-bold font-[Inter] bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent mb-1.5 leading-tight">{{ t('organization.title') }}</h1>
        <p class="text-[13px] text-muted-foreground m-0">管理医疗机构及其层级关系</p>
      </div>
      <button class="flex items-center gap-2 px-5 py-2.5 bg-gradient-to-r from-primary to-accent text-white font-semibold text-sm border-none rounded-[10px] cursor-pointer transition-all duration-300 hover:shadow-[0_0_24px_rgba(63,81,181,0.35)] hover:-translate-y-0.5 shimmer-btn" @click="handleCreate">
        <Plus class="w-4 h-4" />
        {{ t('organization.createOrganization') }}
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-3 gap-4 mb-6">
      <div class="spotlight-card flex items-center gap-4 p-5 bg-card border border-border/60 rounded-[14px] shadow-card cursor-default" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="w-11 h-11 rounded-[10px] bg-primary/12 border border-primary/20 flex items-center justify-center text-primary flex-shrink-0">
          <Building2 class="w-[22px] h-[22px]" />
        </div>
        <div class="flex flex-col gap-0.5">
          <span class="text-2xl font-bold text-foreground font-[Source_Code_Pro] leading-none">{{ pagination.total }}</span>
          <span class="text-xs text-muted-foreground">{{ t('organization.title') }}</span>
        </div>
      </div>
      <div class="spotlight-card flex items-center gap-4 p-5 bg-card border border-border/60 rounded-[14px] shadow-card cursor-default" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="w-11 h-11 rounded-[10px] bg-primary/12 border border-primary/20 flex items-center justify-center text-primary flex-shrink-0">
          <GitBranch class="w-[22px] h-[22px]" />
        </div>
        <div class="flex flex-col gap-0.5">
          <span class="text-2xl font-bold text-foreground font-[Source_Code_Pro] leading-none">{{ tableData.filter(o => !o.parent_id).length }}</span>
          <span class="text-xs text-muted-foreground">根节点组织</span>
        </div>
      </div>
      <div class="spotlight-card flex items-center gap-4 p-5 bg-card border border-border/60 rounded-[14px] shadow-card cursor-default" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
        <div class="w-11 h-11 rounded-[10px] bg-primary/12 border border-primary/20 flex items-center justify-center text-primary flex-shrink-0">
          <User class="w-[22px] h-[22px]" />
        </div>
        <div class="flex flex-col gap-0.5">
          <span class="text-2xl font-bold text-foreground font-[Source_Code_Pro] leading-none">{{ tableData.reduce((sum, o) => sum + (o.member_count || 0), 0) }}</span>
          <span class="text-xs text-muted-foreground">总成员数</span>
        </div>
      </div>
    </div>

    <!-- Search and Filter -->
    <div class="bg-card border border-border/60 rounded-[14px] px-5 py-4 mb-5">
      <div class="flex items-center flex-wrap gap-3">
        <div class="relative w-[280px]">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            v-model="searchForm.search"
            :placeholder="t('organization.name') + ' / ' + t('organization.code')"
            class="pl-9"
            @keyup.enter="handleSearch"
          />
          <button v-if="searchForm.search" @click="searchForm.search = ''; handleSearch()" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
            <X class="w-4 h-4" />
          </button>
        </div>
        <div class="w-[220px]">
          <Select v-model="searchForm.parent_id">
            <SelectTrigger>
              <SelectValue :placeholder="t('organization.parentOrganization')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">全部</SelectItem>
              <template v-for="org in flattenOrgTree(organizationTree)" :key="org.id">
                <SelectItem :value="org.id">{{ org.name }}</SelectItem>
              </template>
            </SelectContent>
          </Select>
        </div>
        <div class="flex gap-2">
          <Button variant="outline" @click="handleReset">
            <RefreshCw class="w-4 h-4" />
            {{ t('common.reset') }}
          </Button>
          <Button @click="handleSearch">
            <Search class="w-4 h-4" />
            {{ t('common.search') }}
          </Button>
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <div class="spotlight-card flex-1 min-h-0 flex flex-col bg-card border border-border/60 rounded-[14px] overflow-hidden shadow-card" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
      <ListPageSkeleton v-if="initialLoading" :columns="6" :rows="8" />
      <template v-else>
        <div class="flex-1 min-h-0 overflow-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="uppercase text-xs tracking-wider">{{ t('organization.name') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[140px]">{{ t('organization.code') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider">{{ t('organization.parentOrganization') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[140px]">{{ t('organization.facilityType') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[110px] text-center">{{ t('organization.memberCount') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[110px] text-center">{{ t('organization.departmentCount') }}</TableHead>
                <TableHead class="uppercase text-xs tracking-wider w-[180px] text-right">{{ t('common.actions') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="row in tableData" :key="row.id">
                <TableCell>
                  <div class="flex items-center gap-2.5">
                    <div class="w-7 h-7 rounded-md bg-primary/10 flex items-center justify-center text-primary flex-shrink-0">
                      <Building2 class="w-3.5 h-3.5" />
                    </div>
                    <span class="font-medium text-foreground">{{ row.name }}</span>
                  </div>
                </TableCell>
                <TableCell>
                  <code class="inline-block px-2 py-0.5 bg-primary/8 border border-primary/20 rounded-md font-[Source_Code_Pro] text-xs text-primary">{{ row.code || '-' }}</code>
                </TableCell>
                <TableCell>
                  <span v-if="row.parent" class="text-muted-foreground text-[13px]">{{ row.parent.name }}</span>
                  <Badge v-else variant="outline" class="bg-green-500/10 border-green-500/30 text-green-500 text-xs">{{ t('organization.root') }}</Badge>
                </TableCell>
                <TableCell>
                  <span class="text-muted-foreground text-[13px]">{{ row.facility_type || '-' }}</span>
                </TableCell>
                <TableCell class="text-center">
                  <span class="inline-block px-2.5 py-0.5 bg-primary/8 rounded-full text-sm font-semibold text-foreground">{{ row.member_count || 0 }}</span>
                </TableCell>
                <TableCell class="text-center">
                  <span class="inline-block px-2.5 py-0.5 bg-primary/8 rounded-full text-sm font-semibold text-foreground">{{ row.department_count || 0 }}</span>
                </TableCell>
                <TableCell class="text-right">
                  <div class="flex items-center justify-end gap-1.5">
                    <button class="flex items-center justify-center w-[30px] h-[30px] rounded-lg border border-border bg-input cursor-pointer transition-all duration-200 text-muted-foreground hover:text-primary hover:border-primary/30 hover:bg-primary/8 hover:-translate-y-0.5" @click="handleView(row)" :title="t('common.view')">
                      <Eye class="w-3.5 h-3.5" />
                    </button>
                    <button class="flex items-center justify-center w-[30px] h-[30px] rounded-lg border border-border bg-input cursor-pointer transition-all duration-200 text-muted-foreground hover:text-amber-600 hover:border-amber-500/40 hover:bg-amber-500/12 hover:-translate-y-0.5" @click="handleEdit(row)" :title="t('common.edit')">
                      <Pencil class="w-3.5 h-3.5" />
                    </button>
                    <button class="flex items-center justify-center w-[30px] h-[30px] rounded-lg border border-border bg-input cursor-pointer transition-all duration-200 text-muted-foreground hover:text-red-500 hover:border-red-500/40 hover:bg-red-500/12 hover:-translate-y-0.5" @click="handleDelete(row)" :title="t('common.delete')">
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <!-- Pagination -->
        <div class="flex justify-end items-center px-5 py-4 border-t border-border">
          <div class="flex items-center gap-4 text-sm">
            <span class="text-muted-foreground">共 {{ pagination.total }} 条</span>
            <Select v-model="paginationLimitString" @update:model-value="(val: string) => handleSizeChange(Number(val))">
              <SelectTrigger class="w-[80px]">
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
              <Button variant="outline" size="sm" :disabled="pagination.page <= 1" @click="handlePageChange(pagination.page - 1)">
                <ChevronLeft class="w-4 h-4" />
              </Button>
              <span class="px-3 text-sm">{{ pagination.page }}</span>
              <Button variant="outline" size="sm" :disabled="pagination.page * pagination.limit >= pagination.total" @click="handlePageChange(pagination.page + 1)">
                <ChevronRight class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:open="dialogVisible">
      <DialogContent class="max-w-[560px]">
        <DialogHeader>
          <DialogTitle>{{ dialogMode === 'create' ? t('organization.createOrganization') : t('organization.editOrganization') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>{{ t('organization.name') }}</Label>
            <Input v-model="formData.name" :placeholder="t('organization.namePlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.code') }}</Label>
            <Input v-model="formData.code" :placeholder="t('organization.codePlaceholder')" maxlength="50" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.parentOrganization') }}</Label>
            <Select v-model="formData.parent_id">
              <SelectTrigger>
                <SelectValue :placeholder="t('organization.parentOrganizationPlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">无</SelectItem>
                <template v-for="org in flattenOrgTree(organizationTree)" :key="org.id">
                  <SelectItem :value="org.id">{{ org.name }}</SelectItem>
                </template>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.facilityType') }}</Label>
            <Input v-model="formData.facility_type" :placeholder="t('organization.facilityTypePlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.accreditationStatus') }}</Label>
            <Input v-model="formData.accreditation_status" :placeholder="t('organization.accreditationStatusPlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.provinceCity') }}</Label>
            <div class="flex flex-wrap gap-2 max-h-40 overflow-y-auto p-2 border border-input rounded-md">
              <template v-for="province in provinceCityOptions" :key="province.label">
                <template v-for="city in province.children" :key="city.label">
                  <label class="flex items-center gap-1 text-sm">
                    <Checkbox
                      :checked="formDataProvinceCity.includes(city.label)"
                      @update:checked="(checked: boolean) => toggleCitySelection(checked, city.label)"
                    />
                    <span class="text-muted-foreground">{{ province.label }} - {{ city.label }}</span>
                  </label>
                </template>
              </template>
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="dialogVisible = false">{{ t('common.cancel') }}</Button>
          <Button @click="handleSubmit" :disabled="submitting">{{ t('common.confirm') }}</Button>
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
import { Building2, GitBranch, User, Eye, Pencil, Trash2, ChevronLeft, ChevronRight, X, RefreshCw, Search, Plus } from 'lucide-vue-next'
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
  set: (val: string) => { pagination.limit = Number(val) },
})
const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const formRef = ref()
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
  if (!formRef.value) return
  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return
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
  })
}

function resetForm() { Object.assign(formData, { name: '', code: '', parent_id: 'none', facility_type: '', accreditation_status: '', province_city: [] }) }
function handleSizeChange(size: number) { pagination.limit = size; pagination.page = 1; fetchData() }
function handlePageChange(page: number) { pagination.page = page; fetchData() }
</script>

<style scoped>
.spotlight-card {
  position: relative;
  overflow: hidden;
}

.spotlight-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(520px circle at var(--mouse-x, -100%) var(--mouse-y, -100%),
    rgba(63, 81, 181, 0.12), transparent 55%);
  pointer-events: none;
  z-index: 1;
  transition: opacity 0.3s;
}
</style>
