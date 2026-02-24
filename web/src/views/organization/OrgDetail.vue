<template>
  <div class="org-detail">
    <el-card>
      <template #header>
        <span class="title">组织详情</span>
      </template>
      <p>组织 ID: {{ orgId }}</p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
const orgId = route.params.id as string
</script>

<style scoped lang="scss">
.title {
  color: #D4AF37;
  font-family: 'Cinzel', serif;
}
</style>
