import request from '../request'
import type { AuditLog, AuditLogStats, ListAuditLogsParams } from '@/types/audit'

export type AuditLogDTO = AuditLog

export const auditApi = {
  // 获取审计日志列表
  getAuditLogs: (params?: ListAuditLogsParams) =>
    request<{ audit_logs: AuditLog[]; page: any; stats?: AuditLogStats }>({
      method: 'GET',
      url: '/api/v1/identity/audit-logs',
      params,
    }),
}
