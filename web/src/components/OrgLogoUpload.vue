<template>
  <div class="org-logo-upload">
    <div
      class="logo-uploader"
      @click="handleUploaderClick"
      @keydown.enter="handleUploaderClick"
      tabindex="0"
      role="button"
      :aria-disabled="disabled || uploading"
    >
      <input
        ref="fileInputRef"
        type="file"
        accept="image/*"
        class="hidden"
        @change="onFileChange"
        :disabled="disabled || uploading"
      />
      <div v-if="logoUrl" class="logo-preview">
        <img
          :src="logoUrl"
          class="w-full h-full object-contain"
          alt="Logo"
        />
        <div class="preview-mask" v-if="!disabled">
          <ZoomIn class="w-6 h-6" />
          <span>{{ t('common.preview') }}</span>
        </div>
      </div>
      <div v-else class="upload-placeholder">
        <Plus class="upload-icon w-10 h-10" />
        <div class="upload-text">{{ t('organization.uploadLogo') }}</div>
        <div class="upload-hint">{{ t('organization.uploadLogoHint') }}</div>
      </div>
    </div>

    <div class="logo-actions" v-if="logoUrl && !disabled">
      <Button
        variant="outline"
        size="sm"
        @click="triggerFileSelect"
        :disabled="uploading"
      >
        <Plus class="w-4 h-4 mr-1" />
        {{ t('common.replace') }}
      </Button>
      <Button
        variant="destructive"
        size="sm"
        @click="handleRemove"
        :disabled="uploading"
      >
        <Trash2 class="w-4 h-4 mr-1" />
        {{ t('common.remove') }}
      </Button>
    </div>

    <Progress
      v-if="uploading"
      :model-value="uploadProgress"
      class="upload-progress h-1.5"
    />

    <p v-if="errorMessage" class="logo-error-tip" role="alert">
      {{ errorMessage }}
    </p>

    <Dialog v-model:open="previewVisible">
      <DialogContent class="max-w-[640px]">
        <img
          v-if="logoUrl"
          :src="logoUrl"
          class="w-full h-auto object-contain"
          alt="Logo preview"
        />
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { Plus, ZoomIn, Trash2 } from 'lucide-vue-next'
import { organizationApi } from '@/api/organization'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Dialog, DialogContent } from '@/components/ui/dialog'

interface Props {
  organizationId?: string
  logoUrl?: string
  disabled?: boolean
}

interface Emits {
  (e: 'update:logoUrl', value: string): void
  (e: 'upload-success', logoId: string, logoUrl: string): void
  (e: 'remove'): void
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false
})

const emit = defineEmits<Emits>()

const { t } = useI18n()

const uploading = ref(false)
const uploadProgress = ref(0)
const fileInputRef = ref<HTMLInputElement>()
const errorMessage = ref('')
const previewVisible = ref(false)
let errorTimer: ReturnType<typeof setTimeout> | null = null

function triggerFileSelect() {
  if (props.disabled || uploading.value) return
  fileInputRef.value?.click()
}

function handleUploaderClick() {
  if (props.disabled || uploading.value) return
  if (props.logoUrl) {
    previewVisible.value = true
  } else {
    fileInputRef.value?.click()
  }
}

function showError(msg: string) {
  errorMessage.value = msg
  toast.error(msg)
  if (errorTimer) clearTimeout(errorTimer)
  errorTimer = setTimeout(() => {
    errorMessage.value = ''
  }, 4000)
}

watch(() => props.logoUrl, (newVal) => {
  if (!newVal) {
    uploadProgress.value = 0
  }
})

async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  errorMessage.value = ''

  if (!file.type.startsWith('image/')) {
    showError(t('organization.logoFormatError'))
    return
  }

  if (file.size / 1024 / 1024 >= 2) {
    showError(t('organization.logoSizeError'))
    return
  }

  await handleUpload(file)
}

async function handleUpload(file: File) {
  uploading.value = true
  uploadProgress.value = 0

  try {
    const response = await organizationApi.uploadTemporaryLogo(
      file.name,
      file,
      file.type
    )

    const logo = response.logo
    if (!logo) {
      throw new Error('上传失败：未返回 Logo 信息')
    }

    uploadProgress.value = 100

    const logoUrl = logo.download_url || ''

    emit('update:logoUrl', logoUrl)
    emit('upload-success', logo.id || '', logoUrl)

    toast.success(t('common.uploadSuccess'))

    if (props.organizationId && logo.id) {
      await organizationApi.bindLogoToOrganization(props.organizationId, logo.id)
      toast.success(t('organization.logoBindSuccess'))
    }
  } catch (error: any) {
    console.error('Failed to upload logo:', error)
    toast.error(error.message || t('common.uploadFailed'))
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

function handleRemove() {
  emit('update:logoUrl', '')
  emit('remove')
  toast.success(t('common.removeSuccess'))
}
</script>

<style scoped lang="scss">
.org-logo-upload {
  .logo-uploader {
    position: relative;
    overflow: hidden;
    cursor: pointer;
    border: 1px dashed rgba(212, 175, 55, 0.3);
    border-radius: 8px;
    background-color: rgba(44, 46, 51, 0.3);
    transition: all 0.3s;

    &:hover {
      border-color: var(--c-accent);
    }

    &[aria-disabled="true"] {
      cursor: not-allowed;
      opacity: 0.6;
    }
  }

  .logo-preview {
    width: 200px;
    height: 200px;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;

    .preview-mask {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background-color: rgba(0, 0, 0, 0.6);
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      color: #fff;
      font-size: 14px;
      opacity: 0;
      transition: opacity 0.3s;
      gap: 8px;
    }

    &:hover .preview-mask {
      opacity: 1;
    }
  }

  .upload-placeholder {
    width: 200px;
    height: 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--c-text-sub);

    .upload-icon {
      color: var(--c-accent);
      margin-bottom: 10px;
    }

    .upload-text {
      font-size: 14px;
      color: var(--c-accent);
      margin-bottom: 5px;
    }

    .upload-hint {
      font-size: 12px;
      color: var(--c-text-sub);
    }
  }

  .logo-actions {
    margin-top: 15px;
    display: flex;
    gap: 10px;
  }

  .upload-progress {
    margin-top: 10px;
  }

  .logo-error-tip {
    margin-top: 10px;
    font-size: 13px;
    color: #ef4444;
    text-align: center;
  }
}
</style>
