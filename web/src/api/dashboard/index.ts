import request from '../request'

export interface DashboardStats {
  userCount: number
  orgCount: number
  roleCount: number
  activityCount: number
}

/**
 * 获取仪表盘统计数据
 * 通过调用各个模块的列表接口来获取统计数据
 */
export function getDashboardStats() {
  // 并行请求用户、组织、角色列表的统计数据
  return Promise.all([
    // 用户总数
    request<{ page: { total: number } }>({
      url: '/api/v1/identity/users',
      method: 'GET',
      params: { limit: 1 }
    }),
    // 组织总数
    request<{ page: { total: number } }>({
      url: '/api/v1/identity/organizations',
      method: 'GET',
      params: { limit: 1 }
    }),
    // 角色总数
    request<{ page: { total: number } }>({
      url: '/api/v1/permission/roles',
      method: 'GET',
      params: { limit: 1 }
    })
  ]).then(([usersRes, orgsRes, rolesRes]) => ({
    userCount: usersRes.page?.total || 0,
    orgCount: orgsRes.page?.total || 0,
    roleCount: rolesRes.page?.total || 0,
    activityCount: 0 // 暂时设为 0，后续可以添加活动记录功能
  }))
}
