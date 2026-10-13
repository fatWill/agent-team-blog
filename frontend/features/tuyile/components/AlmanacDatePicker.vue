<template>
  <ConfigProvider :locale="zhCN" :theme="antdTheme">
    <DatePicker
      class="tuyile-hidden-picker"
      :value="value"
      :open="open"
      :allow-clear="false"
      :input-read-only="true"
      :disabled-date="disabledDate"
      :show-today="true"
      placement="bottomRight"
      popup-class-name="tuyile-picker-popup"
      aria-hidden="true"
      :tabindex="-1"
      @update:open="onOpenChange"
      @change="onChange"
    />
  </ConfigProvider>
</template>

<script setup lang="ts">
/**
 * 快速选日（antd DatePicker：日 / 月 / 年面板可切换）。
 * 只渲染弹层，输入框隐藏并铺满父级触发按钮，由父组件控制 open。
 * 由父组件按需异步加载并包在 <ClientOnly> 中，避免 antd 进入首屏 chunk 及 SSR。
 */
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import DatePicker from 'ant-design-vue/es/date-picker'
import ConfigProvider from 'ant-design-vue/es/config-provider'
import theme from 'ant-design-vue/es/theme'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { ALMANAC_MAX_DATE, ALMANAC_MIN_DATE } from '../api'

dayjs.locale('zh-cn')
const props = defineProps<{ modelValue: string, open: boolean, dark: boolean }>()
const emit = defineEmits<{
  (e: 'update:modelValue', v: string): void
  (e: 'update:open', v: boolean): void
}>()

const MIN = dayjs(ALMANAC_MIN_DATE)
const MAX = dayjs(ALMANAC_MAX_DATE)

const value = computed(() => dayjs(props.modelValue))

const antdTheme = computed(() => ({
  algorithm: props.dark ? theme.darkAlgorithm : theme.defaultAlgorithm,
  token: { colorPrimary: '#e11d48', colorLink: '#e11d48', colorLinkHover: '#f43f5e', borderRadius: 8 },
}))

function disabledDate(d: Dayjs): boolean {
  return d.isBefore(MIN, 'day') || d.isAfter(MAX, 'day')
}

function onOpenChange(v: boolean) {
  emit('update:open', v)
}

function onChange(d: Dayjs | string | null) {
  if (!d || typeof d === 'string') return
  const s = d.format('YYYY-MM-DD')
  if (s < ALMANAC_MIN_DATE || s > ALMANAC_MAX_DATE) return
  emit('update:modelValue', s)
  emit('update:open', false)
}
</script>

<style scoped>
.tuyile-hidden-picker {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  pointer-events: none;
}
</style>
