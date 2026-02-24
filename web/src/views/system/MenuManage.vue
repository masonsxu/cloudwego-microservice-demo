<template>
  <div class="menu-manage">
    <el-card>
      <template #header>
        <span class="title">{{ t('system.menuManage') }}</span>
      </template>
      <p>菜单管理页面 - 待开发</p>
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
