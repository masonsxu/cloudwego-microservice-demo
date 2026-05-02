<template>
  <div class="p-5">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-5">
      <div class="flex items-center gap-3">
        <Button variant="ghost" size="sm" @click="handleBack">
          <ArrowLeft class="w-4 h-4 mr-1" />
          {{ t('common.back') }}
        </Button>
        <h1 class="text-xl font-bold font-[Inter] text-primary">
          {{ orgDetail?.name || t('organization.orgDetail') }}
        </h1>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" @click="handleEdit">
          <Pencil class="w-4 h-4 mr-1" />
          {{ t('organization.editOrganization') }}
        </Button>
        <Button variant="destructive" @click="handleDelete">
          <Trash2 class="w-4 h-4 mr-1" />
          {{ t('organization.deleteOrganization') }}
        </Button>
      </div>
    </div>

    <div class="mt-5">
      <DetailPageSkeleton
        v-if="initialLoading"
        :side-span="8"
        :main-span="16"
        :side-cards="2"
        :main-cards="3"
        :show-avatar="true"
        :side-item-counts="[3, 2]"
        :main-item-counts="[6, 4, 2]"
      />
      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-5">
        <!-- Left Column -->
        <div class="space-y-5">
          <!-- Basic Info -->
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <Building2 class="w-4 h-4" />
                <span>{{ t('organization.basicInfo') }}</span>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div class="flex flex-col items-center mb-5">
                <div class="w-[100px] h-[100px] rounded-full bg-gradient-to-br from-primary to-accent flex items-center justify-center text-white text-3xl font-semibold">
                  <img v-if="orgDetail?.logo" :src="orgDetail.logo" class="w-full h-full rounded-full object-cover" />
                  <Building2 v-else class="w-10 h-10" />
                </div>
              </div>
              <div class="space-y-3">
                <div class="flex justify-between text-sm">
                  <span class="text-muted-foreground">{{ t('organization.name') }}</span>
                  <span class="text-foreground font-medium">{{ orgDetail?.name }}</span>
                </div>
                <div class="flex justify-between text-sm">
                  <span class="text-muted-foreground">{{ t('organization.code') }}</span>
                  <span class="text-foreground font-medium">{{ orgDetail?.code || '-' }}</span>
                </div>
                <div class="flex justify-between text-sm">
                  <span class="text-muted-foreground">{{ t('organization.parentOrganization') }}</span>
                  <span class="text-foreground font-medium">
                    <Badge v-if="!orgDetail?.parent_id" variant="outline" class="text-xs">{{ t('organization.root') }}</Badge>
                    <router-link v-else :to="`/organizations/${orgDetail.parent_id}`" class="text-primary hover:underline">{{ orgDetail?.parent?.name }}</router-link>
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Statistics -->
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <BarChart3 class="w-4 h-4" />
                <span>{{ t('organization.statistics') }}</span>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div class="grid grid-cols-2 gap-2.5">
                <div class="text-center p-4 bg-input rounded-lg border border-primary/20">
                  <div class="text-2xl font-semibold text-primary mb-1">{{ orgDetail?.member_count || 0 }}</div>
                  <div class="text-xs text-muted-foreground">{{ t('organization.memberCount') }}</div>
                </div>
                <div class="text-center p-4 bg-input rounded-lg border border-primary/20">
                  <div class="text-2xl font-semibold text-primary mb-1">{{ orgDetail?.department_count || 0 }}</div>
                  <div class="text-xs text-muted-foreground">{{ t('organization.departmentCount') }}</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Right Column -->
        <div class="lg:col-span-2 space-y-5">
          <!-- Detail Info -->
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <Info class="w-4 h-4" />
                <span>{{ t('organization.detailInfo') }}</span>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('organization.facilityType') }}</span>
                  <p class="text-sm text-foreground">{{ orgDetail?.facility_type || '-' }}</p>
                </div>
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('organization.accreditationStatus') }}</span>
                  <p class="text-sm text-foreground">
                    <Badge v-if="orgDetail?.accreditation_status" variant="outline" class="bg-green-500/10 border-green-500/30 text-green-500 text-xs">{{ orgDetail.accreditation_status }}</Badge>
                    <span v-else>-</span>
                  </p>
                </div>
                <div class="space-y-1 col-span-2">
                  <span class="text-xs text-muted-foreground">{{ t('organization.provinceCity') }}</span>
                  <div class="flex flex-wrap gap-1 mt-1">
                    <Badge v-for="city in orgDetail?.province_city" :key="city" variant="secondary" class="text-xs">{{ city }}</Badge>
                    <span v-if="!orgDetail?.province_city || orgDetail.province_city.length === 0" class="text-sm text-muted-foreground">-</span>
                  </div>
                </div>
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('common.createTime') }}</span>
                  <p class="text-sm text-foreground">{{ formatDateTime(orgDetail?.created_at) }}</p>
                </div>
                <div class="space-y-1">
                  <span class="text-xs text-muted-foreground">{{ t('common.updateTime') }}</span>
                  <p class="text-sm text-foreground">{{ formatDateTime(orgDetail?.updated_at) }}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Child Organizations -->
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <FolderOpen class="w-4 h-4" />
                <span>{{ t('organization.childOrganizations') }}</span>
              </div>
            </CardHeader>
            <CardContent class="p-5">
              <div v-if="!orgDetail?.children || orgDetail.children.length === 0" class="text-center py-8 text-muted-foreground">
                {{ t('organization.noChildOrganizations') }}
              </div>
              <div v-else class="space-y-2.5">
                <div v-for="child in orgDetail?.children" :key="child.id" class="flex items-center justify-between p-4 bg-input rounded-lg border border-border/60">
                  <div class="flex items-center gap-4">
                    <Building2 class="w-8 h-8 text-primary" />
                    <div>
                      <div class="font-medium text-foreground">{{ child.name }}</div>
                      <div class="flex gap-4 text-xs text-muted-foreground mt-1">
                        <span>{{ t('organization.memberCount') }}: {{ child.member_count || 0 }}</span>
                        <span>{{ t('organization.departmentCount') }}: {{ child.department_count || 0 }}</span>
                      </div>
                    </div>
                  </div>
                  <Button variant="ghost" size="sm" @click="router.push(`/organizations/${child.id}`)">
                    {{ t('common.view') }}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Logo Management -->
          <Card class="bg-card border-border/60 rounded-[18px] shadow-card">
            <CardHeader class="border-b border-border/60 px-5 py-3.5">
              <div class="flex items-center gap-2 text-primary font-semibold">
                <ImageIcon class="w-4 h-4" />
                <span>{{ t('organization.logoManagement') }}</span>
              </div>
            </CardHeader>
            <CardContent class="p-5 flex justify-center min-h-[250px] items-center">
              <OrgLogoUpload
                :organization-id="orgId"
                :logo-url="orgDetail?.logo"
                @update:logo-url="handleLogoUpdate"
                @upload-success="handleLogoUploadSuccess"
                @remove="handleLogoRemove"
              />
            </CardContent>
          </Card>
        </div>
      </div>
    </div>

    <!-- Edit Dialog -->
    <Dialog v-model:open="editDialogVisible">
      <DialogContent class="max-w-[600px]">
        <DialogHeader>
          <DialogTitle>{{ t('organization.editOrganization') }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>{{ t('organization.name') }}</Label>
            <Input v-model="editFormData.name" :placeholder="t('organization.namePlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.code') }}</Label>
            <Input v-model="editFormData.code" :placeholder="t('organization.codePlaceholder')" maxlength="50" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.parentOrganization') }}</Label>
            <Select v-model="editFormData.parent_id">
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
            <Input v-model="editFormData.facility_type" :placeholder="t('organization.facilityTypePlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.accreditationStatus') }}</Label>
            <Input v-model="editFormData.accreditation_status" :placeholder="t('organization.accreditationStatusPlaceholder')" maxlength="100" />
          </div>
          <div class="space-y-2">
            <Label>{{ t('organization.provinceCity') }}</Label>
            <div class="flex flex-wrap gap-2 max-h-40 overflow-y-auto p-2 border border-input rounded-md">
              <template v-for="province in provinceCityOptions" :key="province.label">
                <template v-for="city in province.children" :key="city.label">
                  <label class="flex items-center gap-1 text-sm">
                    <Checkbox
                      :checked="editFormDataProvinceCity.includes(city.label)"
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
          <Button variant="outline" @click="editDialogVisible = false">{{ t('common.cancel') }}</Button>
          <Button @click="handleUpdate" :disabled="updating">{{ t('common.confirm') }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { ArrowLeft, Pencil, Trash2, Building2, BarChart3, Info, FolderOpen, ImageIcon } from 'lucide-vue-next'
import { organizationApi } from '@/api/organization'
import { formatOptionalTimestamp } from '@/utils/date'
import OrgLogoUpload from '@/components/OrgLogoUpload.vue'
import DetailPageSkeleton from '@/components/skeleton/DetailPageSkeleton.vue'
import type { OrganizationDTO, UpdateOrganizationRequestDTO } from '@/api/organization'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { Select, SelectValue, SelectTrigger, SelectContent, SelectItem } from '@/components/ui/select'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const loading = ref(false)
const initialLoading = ref(true)
const updating = ref(false)
const orgDetail = ref<OrganizationDTO | null>(null)
const organizationTree = ref<any[]>([])

const editDialogVisible = ref(false)

const editFormData = reactive<UpdateOrganizationRequestDTO & { code?: string }>({
  name: '',
  parent_id: 'none',
  facility_type: '',
  accreditation_status: '',
  province_city: []
})

const editFormDataProvinceCity = computed({
  get: () => editFormData.province_city || [],
  set: (val: string[]) => { editFormData.province_city = val },
})

function toggleCitySelection(checked: boolean, city: string) {
  if (checked) {
    editFormData.province_city = [...(editFormData.province_city || []), city]
  } else {
    editFormData.province_city = (editFormData.province_city || []).filter(c => c !== city)
  }
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

onMounted(async () => {
  try {
    await Promise.all([fetchData(), fetchOrganizationTree()])
  } finally {
    initialLoading.value = false
  }
})

function flattenOrgTree(orgs: any[]): any[] {
  const result: any[] = []
  for (const org of orgs) {
    result.push(org)
    if (org.children?.length) result.push(...flattenOrgTree(org.children))
  }
  return result
}

async function fetchOrganizationTree() {
  try {
    const response = await organizationApi.listOrganizations({ page: 1, limit: 1000 })
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
    toast.error(error.message || t('common.operationFailed'))
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
    parent_id: orgDetail.value.parent_id || 'none',
    facility_type: orgDetail.value.facility_type || '',
    accreditation_status: orgDetail.value.accreditation_status || '',
    province_city: orgDetail.value.province_city || []
  })

  editDialogVisible.value = true
}

async function handleUpdate() {
  updating.value = true
  try {
    const updateData: UpdateOrganizationRequestDTO = {
      name: editFormData.name,
      parent_id: editFormData.parent_id !== 'none' ? editFormData.parent_id : undefined,
      facility_type: editFormData.facility_type,
      accreditation_status: editFormData.accreditation_status,
      province_city: editFormData.province_city
    }

    await organizationApi.updateOrganization(orgId, updateData)
    toast.success(t('common.updateSuccess'))
    editDialogVisible.value = false
    fetchData()
  } catch (error: any) {
    console.error('Failed to update organization:', error)
    toast.error(error.message || t('common.operationFailed'))
  } finally {
    updating.value = false
  }
}

async function handleDelete() {
  try {
    await organizationApi.deleteOrganization(orgId)
    toast.success(t('common.deleteSuccess'))
    router.push('/organizations')
  } catch (error: any) {
    if (error !== 'cancel' && error !== 'close') {
      console.error('Failed to delete organization:', error)
      toast.error(error.message || t('common.operationFailed'))
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
  toast.success(t('organization.logoUploadSuccess'))
  fetchData()
}

function handleLogoRemove() {
  if (orgDetail.value) {
    orgDetail.value.logo = ''
  }
  toast.success(t('organization.logoRemoveSuccess'))
}
</script>

<style scoped>
</style>
