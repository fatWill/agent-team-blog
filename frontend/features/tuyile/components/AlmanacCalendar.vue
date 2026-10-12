<template>
  <div class="rounded-xl border border-rose-100 bg-white p-3 shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-4">
    <div class="mb-3 flex items-center justify-between">
      <button
        type="button"
        class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-rose-50 hover:text-rose-600 disabled:opacity-30 dark:text-gray-400 dark:hover:bg-gray-700"
        aria-label="上个月"
        :disabled="month <= MIN_MONTH"
        @click="month = addMonths(month, -1)"
      >
        ‹
      </button>
      <div class="tuyile-serif text-base font-bold text-gray-800 dark:text-gray-100">
        {{ monthTitle }}
      </div>
      <button
        type="button"
        class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-rose-50 hover:text-rose-600 disabled:opacity-30 dark:text-gray-400 dark:hover:bg-gray-700"
        aria-label="下个月"
        :disabled="month >= MAX_MONTH"
        @click="month = addMonths(month, 1)"
      >
        ›
      </button>
    </div>

    <div class="mb-1 grid grid-cols-7 text-center text-xs text-gray-400 dark:text-gray-500">
      <div v-for="(w, i) in WEEK_HEAD" :key="w" :class="i >= 5 ? 'text-rose-400' : ''">{{ w }}</div>
    </div>

    <div v-if="loading" class="grid grid-cols-7 gap-1">
      <div v-for="i in 35" :key="i" class="h-12 animate-pulse rounded-lg bg-gray-100 dark:bg-gray-700/60 sm:h-14" />
    </div>
    <div v-else-if="errorMsg" class="py-8 text-center text-sm text-gray-500 dark:text-gray-400">
      {{ errorMsg }}
      <button type="button" class="ml-2 text-rose-600 hover:underline dark:text-rose-400" @click="load(month, true)">重试</button>
    </div>
    <div v-else class="grid grid-cols-7 gap-1">
      <div v-for="i in offset" :key="`pad-${i}`" />
      <button
        v-for="cell in cells"
        :key="cell.date"
        type="button"
        class="relative flex h-12 flex-col items-center justify-center rounded-lg border text-center transition-all duration-150 sm:h-14"
        :class="cellClass(cell.date)"
        :aria-label="`${cell.date} ${cell.lunarDate}`"
        :aria-pressed="cell.date === modelValue"
        @click="emit('update:modelValue', cell.date)"
      >
        <span class="text-sm font-semibold leading-none">{{ Number(cell.date.slice(8)) }}</span>
        <span class="mt-1 max-w-full truncate px-0.5 text-[10px] leading-none" :class="subClass(cell)">
          {{ cell.festivals[0] || cell.cellText }}
        </span>
        <span
          v-if="cell.zhiShenType === '黄道'"
          class="absolute right-1 top-1 h-1.5 w-1.5 rounded-full bg-amber-400"
          title="黄道吉日"
        />
      </button>
    </div>

    <div class="mt-3 flex items-center gap-3 text-[11px] text-gray-400 dark:text-gray-500">
      <span class="flex items-center gap-1"><span class="h-1.5 w-1.5 rounded-full bg-amber-400" />黄道日</span>
      <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-sm border border-rose-400" />今天</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AlmanacMonthItem } from '../types'
import { apiFetchAlmanacMonth, getTuyileErrorMessage } from '../api'
import { addMonths, mondayOffset } from '../utils'

const props = defineProps<{ modelValue: string, today: string }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()

const WEEK_HEAD = ['一', '二', '三', '四', '五', '六', '日']
const MIN_MONTH = '1900-01'
const MAX_MONTH = '2100-12'

const month = ref(props.modelValue.slice(0, 7))
const cells = ref<AlmanacMonthItem[]>([])
const loading = ref(false)
const errorMsg = ref('')
const cache = new Map<string, AlmanacMonthItem[]>()
let reqSeq = 0

const offset = computed(() => mondayOffset(month.value))
const monthTitle = computed(() => {
  const [y, m] = month.value.split('-').map(Number)
  return `${y} 年 ${m} 月`
})

async function load(m: string, force = false) {
  const cached = cache.get(m)
  if (cached && !force) {
    cells.value = cached
    errorMsg.value = ''
    return
  }
  const seq = ++reqSeq
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await apiFetchAlmanacMonth(m)
    if (seq !== reqSeq) return
    cache.set(m, res.list)
    cells.value = res.list
  } catch (err: unknown) {
    if (seq !== reqSeq) return
    errorMsg.value = getTuyileErrorMessage(err, '月历加载失败')
  } finally {
    if (seq === reqSeq) loading.value = false
  }
}

watch(month, m => load(m), { immediate: true })
// 外部切日跨月时，日历跟随
watch(() => props.modelValue, (v) => {
  if (v.slice(0, 7) !== month.value) month.value = v.slice(0, 7)
})

function cellClass(date: string): string {
  if (date === props.modelValue) {
    return 'border-rose-500 bg-rose-600 text-white shadow-sm dark:border-rose-500 dark:bg-rose-600'
  }
  const base = 'text-gray-700 hover:bg-rose-50 dark:text-gray-200 dark:hover:bg-gray-700'
  return date === props.today ? `border-rose-400 ${base}` : `border-transparent ${base}`
}

function subClass(cell: AlmanacMonthItem): string {
  if (cell.date === props.modelValue) return 'text-rose-100'
  if (cell.festivals.length || cell.jieQi) return 'text-rose-500 dark:text-rose-400'
  return 'text-gray-400 dark:text-gray-500'
}
</script>
