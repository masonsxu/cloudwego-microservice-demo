<template>
  <div class="timestamp-demo">
    <h2>时间戳处理示例</h2>

    <div class="demo-section">
      <h3>原始数据（后端返回）</h3>
      <pre>{{ JSON.stringify(userData, null, 2) }}</pre>
    </div>

    <div class="demo-section">
      <h3>格式化显示</h3>
      <table>
        <tr>
          <th>字段</th>
          <th>原始值</th>
          <th>格式化</th>
          <th>相对时间</th>
        </tr>
        <tr>
          <td>创建时间</td>
          <td>{{ userData.created_at }}</td>
          <td>{{ formatTime(userData.created_at) }}</td>
          <td>{{ relativeTime(userData.created_at) }}</td>
        </tr>
        <tr>
          <td>更新时间</td>
          <td>{{ userData.updated_at || '-' }}</td>
          <td>{{ formatTime(userData.updated_at) }}</td>
          <td>{{ relativeTime(userData.updated_at) }}</td>
        </tr>
        <tr>
          <td>最后登录</td>
          <td>{{ userData.last_login_time || '-' }}</td>
          <td>{{ formatTime(userData.last_login_time) }}</td>
          <td>{{ relativeTime(userData.last_login_time) }}</td>
        </tr>
      </table>
    </div>

    <div class="demo-section">
      <h3>工具函数说明</h3>
      <ul>
        <li><code>formatTimestamp(timestamp)</code> - 格式化为 'YYYY-MM-DD HH:mm:ss'</li>
        <li><code>formatRelativeTime(timestamp)</code> - 格式化为相对时间（如 '3分钟前'）</li>
        <li><code>formatOptionalTimestamp(timestamp)</code> - 处理可选字段，空值返回 '-'</li>
        <li><code>isTimestampExpired(timestamp, maxAge)</code> - 检查是否过期</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { formatRelativeTime, formatOptionalTimestamp } from '@/utils/date'

// 模拟后端返回的用户数据
const userData = {
  id: 'e91b4a72-6c22-415e-80a9-f268a128f227',
  username: 'superadmin',
  created_at: 1766021112386,  // 毫秒时间戳
  updated_at: 1771919315301,
  last_login_time: 1771919315300,
  account_expiry: undefined  // 可选字段
}

function formatTime(timestamp: number | undefined) {
  return formatOptionalTimestamp(timestamp)
}

function relativeTime(timestamp: number | undefined) {
  if (!timestamp) return '-'
  return formatRelativeTime(timestamp)
}
</script>

<style scoped>
.timestamp-demo {
  padding: 20px;
}

.demo-section {
  margin-bottom: 30px;
  padding: 15px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
}

.demo-section h3 {
  margin-top: 0;
  color: #333;
}

pre {
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

table th,
table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

table th {
  background-color: #f2f2f2;
  font-weight: bold;
}

code {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
}
</style>
