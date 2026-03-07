export interface AuditLog {
  id: string
  request_id: string
  trace_id: string
  user_id?: string
  username?: string
  organization_id?: string
  action: number
  resource: string
  resource_id?: string
  status_code: number
  success: boolean
  client_ip: string
  user_agent?: string
  request_body?: string
  duration_ms: number
  created_at: number
}

export interface AuditLogStats {
  total_count: number
  success_count: number
  avg_duration_ms: number
}

export interface ListAuditLogsParams {
  page?: number
  limit?: number
  user_id?: string
  action?: number
  resource?: string
  success?: boolean
  start_time?: number
  end_time?: number
}
