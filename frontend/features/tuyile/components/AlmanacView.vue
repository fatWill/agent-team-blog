<template>
  <div class="mx-auto max-w-3xl px-4 pb-16 pt-16">
    <!-- 顶部运势 banner -->
    <FortuneBanner class="mb-4" @click="emit('chat')" />

    <!-- 日期切换条 -->
    <div class="mb-4 flex items-center justify-between gap-1.5 sm:gap-2">
      <button
        type="button"
        class="flex h-9 shrink-0 items-center whitespace-nowrap rounded-lg border border-gray-200 bg-white px-2.5 text-sm text-gray-600 transition-all duration-200 hover:border-rose-300 hover:text-rose-600 disabled:opacity-40 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 sm:px-3"
        aria-label="前一天"
        :disabled="date <= ALMANAC_MIN_DATE"
        @click="shift(-1)"
      >
        ‹<span class="ml-1 hidden sm:inline">前一天</span>
      </button>
      <div class="flex min-w-0 items-center gap-1.5 sm:gap-2">
        <button
          type="button"
          class="flex h-9 min-w-0 items-center gap-1 whitespace-nowrap rounded-lg border border-gray-200 bg-white px-3 text-sm font-medium text-gray-700 transition-all duration-200 hover:border-rose-300 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-200"
          :aria-expanded="showCalendar"
          @click="showCalendar = !showCalendar"
        >
          <span aria-hidden="true">📅</span><span class="truncate">{{ formatCnDate(date) }}</span>
        </button>
        <button
          v-if="date !== today"
          type="button"
          class="h-9 shrink-0 whitespace-nowrap rounded-lg bg-rose-50 px-3 text-sm text-rose-600 transition-all duration-200 hover:bg-rose-100 dark:bg-rose-500/10 dark:text-rose-300"
          @click="emit('update:date', today)"
        >
          今天
        </button>
        <div class="relative shrink-0">
          <button
            type="button"
            class="flex h-9 items-center gap-1 whitespace-nowrap rounded-lg border border-rose-200 bg-white px-2.5 text-sm text-rose-600 transition-all duration-200 hover:border-rose-400 hover:bg-rose-50 dark:border-rose-500/30 dark:bg-gray-800 dark:text-rose-300 dark:hover:bg-gray-700"
            aria-label="快速选择年月日"
            title="快速选择年 / 月 / 日"
            :aria-expanded="pickerOpen"
            @click="openPicker"
          >
            <span aria-hidden="true">🗓️</span><span class="hidden sm:inline">选日期</span>
          </button>
          <ClientOnly>
            <LazyTuyileDatePicker
              v-if="pickerMounted"
              v-model:open="pickerOpen"
              :model-value="date"
              :dark="isDark"
              @update:model-value="pickDate"
            />
          </ClientOnly>
        </div>
      </div>
      <button
        type="button"
        class="flex h-9 shrink-0 items-center whitespace-nowrap rounded-lg border border-gray-200 bg-white px-2.5 text-sm text-gray-600 transition-all duration-200 hover:border-rose-300 hover:text-rose-600 disabled:opacity-40 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 sm:px-3"
        aria-label="后一天"
        :disabled="date >= ALMANAC_MAX_DATE"
        @click="shift(1)"
      >
        <span class="mr-1 hidden sm:inline">后一天</span>›
      </button>
    </div>

    <Transition name="tuyile-fade">
      <AlmanacCalendar
        v-if="showCalendar"
        class="mb-4"
        :model-value="date"
        :today="today"
        @update:model-value="pickDate"
      />
    </Transition>

    <!-- 骨架屏 -->
    <div v-if="pending && !data" class="space-y-4">
      <div class="h-56 animate-pulse rounded-xl bg-white shadow-sm dark:bg-gray-800" />
      <div class="h-40 animate-pulse rounded-xl bg-white shadow-sm dark:bg-gray-800" />
      <div class="h-64 animate-pulse rounded-xl bg-white shadow-sm dark:bg-gray-800" />
    </div>

    <div v-else-if="errorMessage && !data" class="rounded-xl bg-white p-10 text-center shadow-sm dark:bg-gray-800">
      <p class="mb-3 text-gray-600 dark:text-gray-300">{{ errorMessage }}</p>
      <button type="button" class="rounded-lg bg-rose-600 px-4 py-2 text-sm text-white hover:bg-rose-700" @click="refresh()">重新加载</button>
    </div>

    <div v-else-if="data" class="space-y-4 transition-opacity duration-200" :class="pending ? 'opacity-60' : ''">
      <!-- 头部：公历 / 农历 -->
      <section class="tuyile-paper relative overflow-hidden rounded-xl border border-rose-200/70 p-5 text-center shadow-sm dark:border-rose-500/20 sm:p-6">
        <div class="mb-1 text-sm text-gray-500 dark:text-gray-400">
          {{ data.date.slice(0, 4) }}年{{ Number(data.date.slice(5, 7)) }}月 · {{ data.week }}
        </div>
        <div class="tuyile-serif text-7xl font-black leading-none text-rose-700 dark:text-rose-400 sm:text-8xl">
          {{ Number(data.date.slice(8)) }}
        </div>
        <div class="tuyile-serif mt-3 text-2xl font-bold tracking-widest text-gray-800 dark:text-gray-100">
          {{ data.lunarDate }}
        </div>
        <div class="mt-2 text-sm text-gray-600 dark:text-gray-300">
          农历{{ data.lunarYearCn }}年 · {{ data.ganZhiYear }}{{ data.zodiacYear }}年 · {{ data.ganZhiMonth }}月 · {{ data.ganZhiDay }}日
        </div>
        <div class="mt-3 flex flex-wrap justify-center gap-2">
          <span v-if="data.jieQi" class="rounded-full bg-emerald-500/90 px-3 py-0.5 text-xs font-medium text-white">今日{{ data.jieQi }}</span>
          <span v-else-if="data.currentJieQi" class="rounded-full bg-emerald-50 px-3 py-0.5 text-xs text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300">节气 · {{ data.currentJieQi }}</span>
          <span v-for="f in data.festivals" :key="f" class="rounded-full bg-rose-500/90 px-3 py-0.5 text-xs font-medium text-white">{{ f }}</span>
          <span
            class="rounded-full px-3 py-0.5 text-xs font-medium"
            :class="data.zhiShenType === '黄道' ? 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300' : 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300'"
          >
            {{ data.zhiShen }} · {{ data.zhiShenType }}日
          </span>
        </div>
      </section>

      <!-- 宜 / 忌 -->
      <section class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div class="rounded-xl border border-emerald-100 bg-white p-4 shadow-sm dark:border-emerald-500/20 dark:bg-gray-800">
          <div class="mb-3 flex items-center gap-3">
            <span class="tuyile-serif flex h-12 w-12 items-center justify-center rounded-full bg-emerald-600 text-2xl font-black text-white">宜</span>
            <span class="text-xs text-gray-400">{{ data.yi.length }} 项</span>
          </div>
          <div class="flex flex-wrap gap-1.5">
            <span v-for="it in data.yi" :key="it" class="rounded-md bg-emerald-50 px-2 py-0.5 text-sm text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300">{{ it }}</span>
            <span v-if="!data.yi.length" class="text-sm text-gray-400">诸事不宜</span>
          </div>
        </div>
        <div class="rounded-xl border border-rose-100 bg-white p-4 shadow-sm dark:border-rose-500/20 dark:bg-gray-800">
          <div class="mb-3 flex items-center gap-3">
            <span class="tuyile-serif flex h-12 w-12 items-center justify-center rounded-full bg-rose-600 text-2xl font-black text-white">忌</span>
            <span class="text-xs text-gray-400">{{ data.ji.length }} 项</span>
          </div>
          <div class="flex flex-wrap gap-1.5">
            <span v-for="it in data.ji" :key="it" class="rounded-md bg-rose-50 px-2 py-0.5 text-sm text-rose-700 dark:bg-rose-500/10 dark:text-rose-300">{{ it }}</span>
            <span v-if="!data.ji.length" class="text-sm text-gray-400">百无禁忌</span>
          </div>
        </div>
      </section>

      <!-- 冲煞 / 值神 等 -->
      <section class="rounded-xl border border-gray-100 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-5">
        <h3 class="tuyile-serif mb-3 flex items-center gap-2 text-base font-bold text-gray-800 dark:text-gray-100">
          <span class="h-4 w-1 rounded-full bg-rose-500" />日神吉凶
        </h3>
        <dl class="grid grid-cols-2 gap-x-4 gap-y-3 text-sm sm:grid-cols-3">
          <div v-for="item in infoItems" :key="item.label" :class="item.wide ? 'col-span-2 sm:col-span-3' : ''">
            <dt class="text-xs text-gray-400 dark:text-gray-500">{{ item.label }}</dt>
            <dd class="mt-0.5 font-medium leading-relaxed text-gray-800 dark:text-gray-100">{{ item.value || '—' }}</dd>
          </div>
        </dl>
      </section>

      <!-- 彭祖百忌 -->
      <section class="rounded-xl border border-gray-100 bg-white p-4 text-center shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-5">
        <h3 class="mb-2 text-xs text-gray-400 dark:text-gray-500">彭祖百忌</h3>
        <p class="tuyile-serif text-lg font-semibold tracking-wider text-gray-800 dark:text-gray-100">
          <span v-for="(line, i) in pengZuLines" :key="i" class="block sm:inline sm:px-2">{{ line }}</span>
        </p>
      </section>

      <!-- 吉神 / 凶神 -->
      <section class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div class="rounded-xl border border-gray-100 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
          <h3 class="mb-2 text-sm font-semibold text-emerald-600 dark:text-emerald-400">吉神宜趋</h3>
          <p class="text-sm leading-relaxed text-gray-700 dark:text-gray-300">{{ joinOr(data.jiShen) }}</p>
        </div>
        <div class="rounded-xl border border-gray-100 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
          <h3 class="mb-2 text-sm font-semibold text-rose-600 dark:text-rose-400">凶神宜忌</h3>
          <p class="text-sm leading-relaxed text-gray-700 dark:text-gray-300">{{ joinOr(data.xiongSha) }}</p>
        </div>
      </section>

      <!-- 方位 -->
      <section class="rounded-xl border border-gray-100 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800 sm:p-5">
        <h3 class="tuyile-serif mb-3 flex items-center gap-2 text-base font-bold text-gray-800 dark:text-gray-100">
          <span class="h-4 w-1 rounded-full bg-amber-500" />吉神方位
        </h3>
        <div class="grid grid-cols-3 gap-2 sm:grid-cols-5">
          <div v-for="d in directions" :key="d.label" class="rounded-lg bg-amber-50/70 p-2 text-center dark:bg-amber-500/10">
            <div class="text-xs text-gray-500 dark:text-gray-400">{{ d.label }}</div>
            <div class="tuyile-serif mt-0.5 font-bold text-amber-700 dark:text-amber-300">{{ d.value || '—' }}</div>
          </div>
        </div>
        <div class="mt-3 rounded-lg bg-gray-50 px-3 py-2 text-sm text-gray-600 dark:bg-gray-900/40 dark:text-gray-300">
          <span class="text-gray-400">胎神占方：</span>{{ data.taiShen || '—' }}
        </div>
      </section>

      <!-- 时辰 -->
      <AlmanacHours v-if="data.hours.length" :hours="data.hours" :current-index="currentHourIndex" />

      <!-- 底部声明 -->
      <div class="pt-2 text-center">
        <p class="text-xs text-gray-400 dark:text-gray-500">🎲 黄历与运势内容仅供娱乐，图一乐～</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ALMANAC_MAX_DATE, ALMANAC_MIN_DATE } from '../api'
import { useAlmanacDay } from '../composables/useAlmanacDay'
import { addDays, beijingHour, formatCnDate, hourToShiChenIndex } from '../utils'
import AlmanacCalendar from './AlmanacCalendar.vue'
import AlmanacHours from './AlmanacHours.vue'
import FortuneBanner from './FortuneBanner.vue'

// antd DatePicker 仅在首次点击时异步加载（独立 chunk，不进首屏，不参与 SSR）
const LazyTuyileDatePicker = defineAsyncComponent(() => import('./AlmanacDatePicker.vue'))

const props = defineProps<{ date: string, today: string }>()
const emit = defineEmits<{
  (e: 'update:date', v: string): void
  (e: 'chat'): void
}>()

const { isDark } = useTheme()
const showCalendar = ref(false)
const pickerOpen = ref(false)
const pickerMounted = ref(false)

function openPicker() {
  showCalendar.value = false
  pickerMounted.value = true
  // 首次挂载后再打开，确保弹层定位正确
  nextTick(() => {
    pickerOpen.value = !pickerOpen.value
  })
}
const dateRef = computed(() => props.date)
const { data, pending, errorMessage, refresh } = useAlmanacDay(dateRef)

function shift(delta: number) {
  emit('update:date', addDays(props.date, delta))
}

function pickDate(v: string) {
  emit('update:date', v)
  showCalendar.value = false
  pickerOpen.value = false
}

const currentHourIndex = computed(() => (props.date === props.today ? hourToShiChenIndex(beijingHour()) : -1))

const infoItems = computed(() => {
  const d = data.value
  if (!d) return []
  return [
    { label: '冲煞', value: d.chongDesc, wide: true },
    { label: '值神', value: `${d.zhiShen}（${d.zhiShenType}）` },
    { label: '建除十二神', value: `${d.jianChu}日` },
    { label: '星宿', value: `${d.xiu} · ${d.xiuLuck}（${d.xiuGong}）` },
    { label: '年纳音', value: d.naYinYear },
    { label: '月纳音', value: d.naYinMonth },
    { label: '日纳音（五行）', value: d.naYinDay },
    { label: '六曜', value: d.liuYao },
    { label: '物候', value: d.wuHou },
    { label: '月相', value: d.yueXiang },
  ]
})

const pengZuLines = computed(() => (data.value?.pengZu || '—').split(/\s+/).filter(Boolean))

const directions = computed(() => {
  const d = data.value
  if (!d) return []
  return [
    { label: '喜神', value: d.xiShen },
    { label: '财神', value: d.caiShen },
    { label: '福神', value: d.fuShen },
    { label: '阳贵', value: d.yangGui },
    { label: '阴贵', value: d.yinGui },
  ]
})

function joinOr(list: string[]): string {
  return list.length ? list.join('  ') : '—'
}
</script>
