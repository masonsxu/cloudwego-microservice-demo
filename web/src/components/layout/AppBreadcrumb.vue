<template>
  <el-breadcrumb separator="/" class="app-breadcrumb">
    <transition-group name="breadcrumb">
      <el-breadcrumb-item v-for="(item, index) in levelList" :key="item.path">
        <span v-if="item.redirect === 'noRedirect' || index === levelList.length - 1" class="no-redirect">
          {{ t(String(item.meta?.title || item.name)) }}
        </span>
        <a v-else @click.prevent="handleLink(item)">
          {{ t(String(item.meta?.title || item.name)) }}
        </a>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter, type RouteLocationMatched } from 'vue-router'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const levelList = ref<RouteLocationMatched[]>([])

function getBreadcrumb() {
  let matched = route.matched.filter(item => item.meta && item.meta.title)
  const first = matched[0]

  if (first && first.path !== '/dashboard') {
    matched = [{ path: '/dashboard', meta: { title: 'dashboard.title' } } as any].concat(matched)
  }

  levelList.value = matched.filter(item => {
    return item.meta && item.meta.title && item.meta.breadcrumb !== false
  })
}

function handleLink(item: RouteLocationMatched) {
  const { redirect, path } = item
  if (redirect) {
    router.push(redirect as string)
    return
  }
  router.push(path)
}

watch(
  () => route.path,
  () => getBreadcrumb(),
  { immediate: true }
)
</script>

<style scoped lang="scss">
.app-breadcrumb {
  display: inline-block;
  font-size: 14px;
  line-height: 60px;
  margin-left: 8px;

  .no-redirect {
    color: #D4AF37;
    cursor: text;
  }

  a {
    color: #8B9bb4;
    cursor: pointer;
    transition: color 0.3s ease;

    &:hover {
      color: #D4AF37;
    }
  }
}

.breadcrumb-enter-active,
.breadcrumb-leave-active {
  transition: all 0.3s ease;
}

.breadcrumb-enter-from,
.breadcrumb-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
