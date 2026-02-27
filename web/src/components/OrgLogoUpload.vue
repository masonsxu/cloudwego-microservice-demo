<template>
  <div class="org-logo-upload">
    <el-upload
      class="logo-uploader"
      :show-file-list="false"
      :before-upload="beforeUpload"
      :http-request="handleUpload"
      :disabled="uploading"
      accept="image/*"
    >
      <div v-if="logoUrl" class="logo-preview">
        <el-image
          :src="logoUrl"
          fit="contain"
          style="width: 100%; height: 100%;"
          :preview-src-list="[logoUrl]"
        />
        <div class="preview-mask" v-if="!disabled">
          <el-icon><ZoomIn /></el-icon>
          <span>{{ t('common.preview') }}</span>
        </div>
      </div>
      <div v-else class="upload-placeholder">
        <el-icon class="upload-icon"><Plus /></el-icon>
        <div class="upload-text">{{ t('organization.uploadLogo') }}</div>
        <div class="upload-hint">{{ t('organization.uploadLogoHint') }}</div>
      </div>
    </el-upload>

    <div class="logo-actions" v-if="logoUrl && !disabled">
      <el-button
        type="danger"
        size="small"
        @click="handleRemove"
        :disabled="uploading"
      >
        <el-icon><Delete /></el-icon>
        {{ t('common.remove') }}
      </el-button>
    </div>

    <el-progress
      v-if="uploading"
      :percentage="uploadProgress"
      :stroke-width="6"
      class="upload-progress"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, ZoomIn, Delete } from '@element-plus/icons-vue'
import { organizationApi } from '@/api/organization'

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

watch(() => props.logoUrl, (newVal) => {
  if (!newVal) {
    uploadProgress.value = 0
  }
})

async function beforeUpload(file: File) {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error(t('organization.logoFormatError'))
    return false
  }

  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    ElMessage.error(t('organization.logoSizeError'))
    return false
  }

  return true
}

async function handleUpload(options: any) {
  const file = options.file as File

  uploading.value = true
  uploadProgress.value = 0

  try {
    const response = await organizationApi.uploadTemporaryLogo({
      fileName: file.name,
      fileContent: file,
      mimeType: file.type
    })

    const logo = response.logo
    if (!logo) {
      throw new Error('上传失败：未返回 Logo 信息')
    }

    uploadProgress.value = 100

    const logoUrl = logo.download_url || ''

    emit('update:logoUrl', logoUrl)
    emit('upload-success', logo.id || '', logoUrl)

    ElMessage.success(t('common.uploadSuccess'))

    if (props.organizationId && logo.id) {
      await organizationApi.bindLogoToOrganization(props.organizationId, logo.id)
      ElMessage.success(t('organization.logoBindSuccess'))
    }
  } catch (error: any) {
    console.error('Failed to upload logo:', error)
    ElMessage.error(error.message || t('common.uploadFailed'))
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

function handleRemove() {
  emit('update:logoUrl', '')
  emit('remove')
  ElMessage.success(t('common.removeSuccess'))
}
</script>

<style scoped lang="scss">
.org-logo-upload {
  .logo-uploader {
    :deep(.el-upload) {
      position: relative;
      overflow: hidden;
      cursor: pointer;
      border: 1px dashed rgba(212, 175, 55, 0.3);
      border-radius: 8px;
      background-color: rgba(44, 46, 51, 0.3);
      transition: all 0.3s;

      &:hover {
        border-color: #D4AF37;
      }
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

      .el-icon {
        font-size: 24px;
      }
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
    color: #8B9bb4;

    .upload-icon {
      font-size: 40px;
      color: #D4AF37;
      margin-bottom: 10px;
    }

    .upload-text {
      font-size: 14px;
      color: #D4AF37;
      margin-bottom: 5px;
    }

    .upload-hint {
      font-size: 12px;
      color: #8B9bb4;
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
}
</style>
