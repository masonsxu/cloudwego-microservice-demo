/**
 * 时间处理工具函数
 *
 * 后端时间戳格式规范：
 * - 所有时间字段均为毫秒级时间戳（number 类型）
 * - 单位：milliseconds（13位整数）
 * - 示例：1766021112386
 */

/**
 * 将毫秒时间戳转换为 Date 对象
 * @param timestamp 毫秒时间戳
 */
export function timestampToDate(timestamp: number): Date {
  return new Date(timestamp)
}

/**
 * 格式化时间戳为本地时间字符串
 * @param timestamp 毫秒时间戳
 * @param format 格式模板，默认 'YYYY-MM-DD HH:mm:ss'
 */
export function formatTimestamp(
  timestamp: number,
  format: string = 'YYYY-MM-DD HH:mm:ss'
): string {
  if (!timestamp) return '-'

  const date = new Date(timestamp)

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化为相对时间（如"3分钟前"）
 * @param timestamp 毫秒时间戳
 */
export function formatRelativeTime(timestamp: number): string {
  if (!timestamp) return '-'

  const now = Date.now()
  const diff = now - timestamp

  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const months = Math.floor(days / 30)
  const years = Math.floor(days / 365)

  if (seconds < 60) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  if (days < 30) return `${Math.floor(days / 7)}周前`
  if (days < 365) return `${months}个月前`

  return `${years}年前`
}

/**
 * 检查时间戳是否过期
 * @param timestamp 毫秒时间戳
 * @param maxAge 最大有效期（毫秒）
 */
export function isTimestampExpired(timestamp: number, maxAge: number): boolean {
  if (!timestamp) return true

  const now = Date.now()
  return now - timestamp > maxAge
}

/**
 * 获取当前时间戳（毫秒）
 */
export function getCurrentTimestamp(): number {
  return Date.now()
}

/**
 * 格式化可选时间戳
 * @param timestamp 可选的毫秒时间戳
 * @param format 格式模板
 */
export function formatOptionalTimestamp(
  timestamp: number | undefined,
  format: string = 'YYYY-MM-DD HH:mm:ss'
): string {
  if (!timestamp) return '-'
  return formatTimestamp(timestamp, format)
}

/**
 * 计算两个时间戳之间的差值（毫秒）
 */
export function timestampDiff(start: number, end: number): number {
  return end - start
}

/**
 * 将时间戳转换为 ISO 8601 格式
 */
export function toISO(timestamp: number): string {
  return new Date(timestamp).toISOString()
}
