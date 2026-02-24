<template>
  <div class="org-list">
    <el-card>
      <template #header>
        <span class="title">{{ t('organization.title') }}</span>
      </template>
      <p>组织管理页面 - 待开发</p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<style scoped lang="scss">
.title {
  color: #D4AF37;
  font-family: 'Cinzel', serif;
}
</style>
