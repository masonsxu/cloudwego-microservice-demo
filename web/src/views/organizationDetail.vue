<template>
  <div class="organization-organizationDetail">
    <div class="menu-card">
      <div class="card-header">
        <span class="title">organizationDetail organization - 待开发</span>
        <span class="meta">Detail</span>
      </div>
      <div class="card-body">
        <p>此页面正在开发中...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped lang="scss">
.menu-card {
  background: var(--bg-card);
  border: 1px solid hsl(var(--border) / 0.6);
  border-radius: 18px;
  box-shadow: var(--shadow-card);
  backdrop-filter: blur(12px);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  color: var(--c-primary);
  font-family: 'Inter', sans-serif;
  font-size: 18px;
  font-weight: 600;
}

.meta {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: var(--c-text-sub);
}

.card-body {
  color: var(--c-text-sub);
}
</style>
