<template>
  <nav class="inline-block text-sm ml-2" aria-label="Breadcrumb">
    <BreadcrumbList class="flex items-center">
      <template v-for="(item, index) in levelList" :key="item.path">
        <BreadcrumbItem class="flex items-center">
          <span
            v-if="item.redirect === 'noRedirect' || index === levelList.length - 1"
            class="text-[var(--c-accent)] cursor-text"
          >
            {{ t(String(item.title || item.name)) }}
          </span>
          <RouterLink
            v-else
            :to="resolveBreadcrumbTarget(item)"
            class="text-[var(--c-text-sub)] cursor-pointer transition-colors hover:text-[var(--c-accent)]"
          >
            {{ t(String(item.title || item.name)) }}
          </RouterLink>
        </BreadcrumbItem>
        <BreadcrumbSeparator v-if="index < levelList.length - 1" class="text-[var(--c-text-sub)]">
          /
        </BreadcrumbSeparator>
      </template>
    </BreadcrumbList>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { BreadcrumbItem, BreadcrumbList, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { menuRouteMap } from '@/router/routes'
import type { AppRouteMeta } from '@/router/routes'

interface BreadcrumbEntry {
  path: string
  title?: string
  name?: string | symbol
  redirect?: string
}

const route = useRoute()
const { t } = useI18n()

const dashboardBreadcrumb: BreadcrumbEntry = {
  path: '/dashboard',
  title: 'dashboard.title'
}

const levelList = computed<BreadcrumbEntry[]>(() => {
  const matched: BreadcrumbEntry[] = route.matched
    .filter(item => {
      const meta = item.meta as AppRouteMeta
      return meta.title && meta.breadcrumb !== false
    })
    .map(item => ({
      path: item.path,
      title: (item.meta as AppRouteMeta).title,
      name: typeof item.name === 'string' || typeof item.name === 'symbol' ? item.name : undefined,
      redirect: typeof item.redirect === 'string' ? item.redirect : undefined
    }))

  const breadcrumbs = matched.length > 0 ? [...matched] : []
  const first = breadcrumbs[0]

  if (!first || first.path !== '/dashboard') {
    breadcrumbs.unshift(dashboardBreadcrumb)
  }

  const meta = route.meta as AppRouteMeta
  if (!meta.hidden || !meta.activeMenu) {
    return breadcrumbs
  }

  const activeMenuRoute = menuRouteMap.get(meta.activeMenu)
  if (!activeMenuRoute || breadcrumbs.some(item => item.path === activeMenuRoute.path)) {
    return breadcrumbs
  }

  breadcrumbs.splice(1, 0, {
    path: activeMenuRoute.path,
    title: activeMenuRoute.title,
    name: activeMenuRoute.name
  })

  return breadcrumbs
})

function resolveBreadcrumbTarget(item: BreadcrumbEntry) {
  if (item.redirect && item.redirect !== 'noRedirect' && item.redirect !== route.path) {
    return item.redirect
  }

  return item.path
}
</script>

