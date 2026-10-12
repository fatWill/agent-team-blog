import type { Ref } from 'vue'
import type { AlmanacDay } from '../types'
import { apiFetchAlmanacDay, getTuyileErrorMessage } from '../api'

/**
 * 单日黄历数据（SSR 预取）。
 * date 为空时请求不带参数，由后端取北京时间今天。
 */
export function useAlmanacDay(date: Ref<string>) {
  const { data, pending, error, refresh } = useAsyncData<AlmanacDay>(
    () => `tuyile-almanac-${date.value || 'today'}`,
    () => apiFetchAlmanacDay(date.value || undefined),
    { watch: [date] },
  )
  const errorMessage = computed(() => (error.value ? getTuyileErrorMessage(error.value, '黄历加载失败，稍后再试～') : ''))
  return { data, pending, errorMessage, refresh }
}
