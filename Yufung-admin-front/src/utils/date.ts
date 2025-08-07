import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

// 扩展dayjs插件
dayjs.extend(relativeTime);

/**
 * 格式化日期时间
 * @param date 日期字符串或Date对象
 * @param format 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的日期字符串
 */
export function formatDateTime(date: string | Date, format: string = 'YYYY-MM-DD HH:mm:ss'): string {
  if (!date) return '';
  return dayjs(date).format(format);
}

/**
 * 格式化日期
 * @param date 日期字符串或Date对象
 * @param format 格式化模板，默认为 'YYYY-MM-DD'
 * @returns 格式化后的日期字符串
 */
export function formatDate(date: string | Date, format: string = 'YYYY-MM-DD'): string {
  if (!date) return '';
  return dayjs(date).format(format);
}

/**
 * 格式化时间
 * @param date 日期字符串或Date对象
 * @param format 格式化模板，默认为 'HH:mm:ss'
 * @returns 格式化后的时间字符串
 */
export function formatTime(date: string | Date, format: string = 'HH:mm:ss'): string {
  if (!date) return '';
  return dayjs(date).format(format);
}

/**
 * 获取相对时间描述
 * @param date 日期字符串或Date对象
 * @returns 相对时间描述，如 '2小时前'
 */
export function getRelativeTime(date: string | Date): string {
  if (!date) return '';
  return dayjs(date).fromNow();
}

/**
 * 检查是否为有效日期
 * @param date 日期字符串或Date对象
 * @returns 是否为有效日期
 */
export function isValidDate(date: string | Date): boolean {
  return dayjs(date).isValid();
} 