<template>
  <div class="flex min-h-screen flex-col bg-gray-50 transition-colors dark:bg-gray-900">
    <!-- 顶部导航 -->
    <header class="sticky top-0 z-30 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex h-14 max-w-5xl items-center px-4">
        <NuxtLink
          to="/renovation"
          class="flex items-center gap-2 text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
          </svg>
          返回装修
        </NuxtLink>
        <span class="mx-3 text-gray-300 dark:text-gray-600">/</span>
        <span class="text-sm font-medium text-gray-900 dark:text-gray-100">成本预算</span>
      </div>
    </header>

    <main class="mx-auto w-full max-w-5xl flex-1 px-4 py-6">
      <!-- 汇总统计卡片 -->
      <div class="mb-6 grid grid-cols-2 gap-3 sm:gap-4">
        <!-- 总预算 -->
        <div class="rounded-xl border border-indigo-200 bg-gradient-to-br from-indigo-50 to-white p-4 shadow-sm dark:border-indigo-800 dark:from-indigo-950/40 dark:to-gray-800">
          <div class="mb-1 flex items-center gap-1.5 text-xs text-indigo-500 dark:text-indigo-400">
            <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
            </svg>
            总预算
          </div>
          <div class="text-lg font-bold text-indigo-700 dark:text-indigo-300 sm:text-xl">{{ formatMoney(summary.totalBudget) }}</div>
        </div>
        <!-- 已支出 -->
        <div
          class="rounded-xl border p-4 shadow-sm"
          :class="summary.totalActual > summary.totalBudget
            ? 'border-red-200 bg-gradient-to-br from-red-50 to-white dark:border-red-800 dark:from-red-950/40 dark:to-gray-800'
            : 'border-emerald-200 bg-gradient-to-br from-emerald-50 to-white dark:border-emerald-800 dark:from-emerald-950/40 dark:to-gray-800'"
        >
          <div
            class="mb-1 flex items-center gap-1.5 text-xs"
            :class="summary.totalActual > summary.totalBudget ? 'text-red-500 dark:text-red-400' : 'text-emerald-500 dark:text-emerald-400'"
          >
            <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            已支出
          </div>
          <div
            class="text-lg font-bold sm:text-xl"
            :class="summary.totalActual > summary.totalBudget ? 'text-red-600 dark:text-red-400' : 'text-emerald-600 dark:text-emerald-400'"
          >
            {{ formatMoney(summary.totalActual) }}
          </div>
        </div>
      </div>

      <!-- 每个大类独立一个表格 -->
      <div class="space-y-5">
        <div v-for="category in budgetData" :key="category.key">
          <!-- 大类标题 -->
          <div class="mb-2 flex items-center gap-2">
            <span
              class="inline-block h-5 w-1 rounded-full"
              :class="categoryBarClass(category.key)"
            />
            <span
              class="inline-flex items-center rounded-md px-2.5 py-1 text-sm font-bold"
              :class="categoryBadgeClass(category.key)"
            >
              {{ category.name }}
            </span>
            <span class="text-xs text-gray-400 dark:text-gray-500">
              （{{ category.items.length }} 项，预算 {{ formatMoney(categoryBudgetSum(category)) }}）
            </span>
            <button
              class="ml-auto flex h-6 w-6 items-center justify-center rounded-md text-gray-300 transition-colors hover:bg-gray-100 hover:text-gray-500 dark:text-gray-600 dark:hover:bg-gray-700 dark:hover:text-gray-300"
              title="下载该分类截图"
              @click="downloadCategory(category)"
            >
              <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
              </svg>
            </button>
          </div>

          <!-- 表格 -->
          <div class="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <div class="overflow-x-auto">
              <table class="w-full min-w-[520px] text-sm">
                <thead>
                  <tr class="border-b border-gray-200 dark:border-gray-700" :class="categoryHeaderBgClass(category.key)">
                    <th class="whitespace-nowrap px-3 py-2.5 text-left font-semibold text-gray-600 dark:text-gray-300 sm:px-4">项目名称</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-right font-semibold text-gray-600 dark:text-gray-300 sm:px-4">预算金额</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-right font-semibold text-gray-600 dark:text-gray-300 sm:px-4">实际支出</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-right font-semibold text-gray-600 dark:text-gray-300 sm:px-4">差额</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-left font-semibold text-gray-600 dark:text-gray-300 sm:px-4">责任人/备注</th>
                  </tr>
                </thead>
                <tbody>
                  <!-- 数据行 -->
                  <tr
                    v-for="(item, itemIdx) in category.items"
                    :key="itemIdx"
                    class="border-b border-gray-100 transition-colors hover:bg-gray-50/50 dark:border-gray-700/50 dark:hover:bg-gray-700/30"
                    :class="itemIdx % 2 === 1 ? 'bg-gray-50/30 dark:bg-gray-800/50' : ''"
                  >
                    <td class="whitespace-nowrap px-3 py-2.5 font-medium text-gray-900 dark:text-gray-100 sm:px-4">
                      {{ item.name }}
                    </td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono text-gray-700 dark:text-gray-300 sm:px-4">
                      {{ formatMoney(item.budget) }}
                    </td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono sm:px-4" :class="actualClass(item)">
                      {{ item.actual !== null ? formatMoney(item.actual) : '—' }}
                    </td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono sm:px-4" :class="diffClass(item)">
                      {{ diffText(item) }}
                    </td>
                    <td class="px-3 py-2.5 text-gray-500 dark:text-gray-400 sm:px-4">
                      {{ item.remark || '—' }}
                    </td>
                  </tr>

                  <!-- 小计行 -->
                  <tr class="border-t border-gray-200 dark:border-gray-600" :class="categorySubtotalBgClass(category.key)">
                    <td class="px-3 py-2.5 sm:px-4">
                      <span class="text-xs font-bold" :class="categorySubtotalTextClass(category.key)">小计</span>
                    </td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono font-bold text-gray-700 dark:text-gray-300 sm:px-4">
                      {{ formatMoney(categoryBudgetSum(category)) }}
                    </td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono font-bold sm:px-4" :class="categoryActualClass(category)">
                      {{ categoryActualSum(category) > 0 ? formatMoney(categoryActualSum(category)) : '—' }}
                    </td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono font-bold sm:px-4" :class="categoryDiffClass(category)">
                      {{ categoryDiffText(category) }}
                    </td>
                    <td class="px-3 py-2.5 sm:px-4" />
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- 总计卡片 -->
      <div class="mt-6 overflow-hidden rounded-xl border-2 border-gray-300 bg-gray-100 shadow-sm dark:border-gray-600 dark:bg-gray-700">
        <div class="overflow-x-auto">
          <table class="w-full min-w-[520px] text-sm">
            <tbody>
              <tr>
                <td class="px-3 py-3 font-bold text-gray-900 dark:text-gray-100 sm:px-4">总计</td>
                <td class="whitespace-nowrap px-3 py-3 text-right font-mono font-bold text-gray-900 dark:text-gray-100 sm:px-4">
                  {{ formatMoney(summary.totalBudget) }}
                </td>
                <td class="whitespace-nowrap px-3 py-3 text-right font-mono font-bold sm:px-4" :class="summary.totalActual > summary.totalBudget ? 'text-red-600 dark:text-red-400' : 'text-emerald-600 dark:text-emerald-400'">
                  {{ formatMoney(summary.totalActual) }}
                </td>
                <td class="whitespace-nowrap px-3 py-3 text-right font-mono font-bold sm:px-4" :class="summary.totalDiff >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'">
                  {{ summary.totalDiff >= 0 ? '-' + formatMoney(summary.totalDiff) : '+' + formatMoney(Math.abs(summary.totalDiff)) }}
                </td>
                <td class="px-3 py-3 sm:px-4" />
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 底部说明 -->
      <p class="mt-4 text-center text-xs text-gray-400 dark:text-gray-500">
        差额 = 预算 − 实际支出 · 绿色为节省 · 红色为超支 · 数据更新于 2026 年 4 月
      </p>
    </main>

    <!-- 右下角悬浮按钮 -->
    <div class="fixed bottom-6 right-6 z-50 flex flex-col gap-3">
      <button
        class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-800 text-white shadow-lg transition-all hover:scale-110 hover:bg-gray-700 dark:bg-gray-200 dark:text-gray-800 dark:hover:bg-gray-300"
        title="分享"
      >
        <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z" />
        </svg>
      </button>
      <button
        class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-800 text-white shadow-lg transition-all hover:scale-110 hover:bg-gray-700 dark:bg-gray-200 dark:text-gray-800 dark:hover:bg-gray-300"
        title="下载预算图片"
        @click="downloadAll"
      >
        <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
        </svg>
      </button>
    </div>

    <!-- 图片预览弹窗 -->
    <Teleport to="body">
      <Transition name="preview-fade">
        <div
          v-if="previewVisible"
          class="fixed inset-0 z-[100] flex items-center justify-center bg-black/70 backdrop-blur-sm p-6"
        >
          <div class="relative flex max-h-[90vh] max-w-[92vw] flex-col rounded-2xl bg-white shadow-2xl dark:bg-gray-800 overflow-hidden">
            <!-- 右上角关闭按钮 -->
            <button
              class="absolute right-3 top-3 z-10 flex h-8 w-8 items-center justify-center rounded-full bg-black/40 text-white backdrop-blur-sm transition hover:bg-black/60"
              @click="closePreview"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>

            <!-- 图片区域（可滚动） -->
            <div class="overflow-y-auto">
              <img
                v-if="previewUrl"
                :src="previewUrl"
                class="block w-full"
                alt="预算图片预览"
                draggable="false"
              />
              <div v-else class="flex h-48 items-center justify-center">
                <svg class="h-8 w-8 animate-spin text-gray-400" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
              </div>
            </div>

            <!-- 底部 tips -->
            <div class="flex items-center justify-center gap-1.5 border-t border-gray-100 px-4 py-2.5 dark:border-gray-700">
              <svg class="h-3.5 w-3.5 flex-shrink-0 text-gray-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18v-5.25m0 0a6.01 6.01 0 001.5-.189m-1.5.189a6.01 6.01 0 01-1.5-.189m3.75 7.478a12.06 12.06 0 01-4.5 0m3.75 2.383a14.406 14.406 0 01-3 0M14.25 18v-.192c0-.983.658-1.823 1.508-2.316a7.5 7.5 0 10-7.517 0c.85.493 1.509 1.333 1.509 2.316V18" />
              </svg>
              <span class="text-xs text-gray-400 dark:text-gray-500">长按图片保存到相册</span>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
// ====== 类型定义 ======
interface BudgetItem {
  name: string
  budget: number
  actual: number | null
  remark: string | null
}

interface BudgetCategory {
  key: string
  name: string
  items: BudgetItem[]
}

// ====== 静态数据 ======
const budgetData: BudgetCategory[] = [
  {
    key: 'partial',
    name: '局改',
    items: [
      { name: '保温层拆除+修复', budget: 2500, actual: 2500, remark: '杨工' },
      { name: '厨房拆墙+修复', budget: 6000, actual: 6000, remark: '杨工' },
      { name: '全屋零线+水电', budget: 2000, actual: null, remark: '杨工' },
      { name: '瓷砖', budget: 1000, actual: null, remark: null },
    ],
  },
  {
    key: 'hard',
    name: '硬装',
    items: [
      { name: '艺术漆', budget: 6800, actual: 6800, remark: '光明京东家居（义峰）' },
      { name: '美缝', budget: 2000, actual: null, remark: null },
      { name: '全屋定制', budget: 35000, actual: null, remark: null },
      { name: '阳台柜+防护网+花台', budget: 5500, actual: null, remark: null },
      { name: '完美系统三联轨', budget: 4000, actual: null, remark: null },
      { name: '全屋纱窗', budget: 1500, actual: null, remark: null },
      { name: '换门', budget: 800, actual: null, remark: null },
    ],
  },
  {
    key: 'furniture',
    name: '家私',
    items: [
      { name: '沙发', budget: 4000, actual: null, remark: null },
      { name: '床垫×3', budget: 5000, actual: null, remark: null },
      { name: '次卧小桌子', budget: 1500, actual: null, remark: null },
      { name: '主卧大桌子', budget: 2000, actual: null, remark: null },
      { name: '电视地柜', budget: 2000, actual: null, remark: null },
      { name: '餐厅桌椅', budget: 250, actual: null, remark: null },
      { name: '窗帘', budget: 2000, actual: null, remark: null },
      { name: '主卧椅子', budget: 1000, actual: null, remark: null },
      { name: '次卧椅子', budget: 200, actual: null, remark: null },
    ],
  },
  {
    key: 'appliance',
    name: '家电',
    items: [
      { name: '电视机', budget: 3500, actual: null, remark: null },
      { name: '洗衣机', budget: 5000, actual: null, remark: null },
      { name: '冰箱', budget: 6500, actual: null, remark: null },
      { name: '洗碗机', budget: 8000, actual: null, remark: null },
      { name: '净饮一体机', budget: 8000, actual: null, remark: null },
      { name: '热水器', budget: 3000, actual: null, remark: null },
      { name: '零冷水', budget: 1000, actual: null, remark: '回流泵' },
      { name: '全屋净水', budget: 2000, actual: null, remark: '大白瓶' },
      { name: '扫地机', budget: 5500, actual: null, remark: null },
      { name: '洗地机', budget: 3000, actual: null, remark: null },
      { name: '智能锁', budget: 2500, actual: null, remark: null },
    ],
  },
  {
    key: 'smart',
    name: '全屋智能',
    items: [
      { name: '智能灯', budget: 3000, actual: null, remark: null },
      { name: '电动窗帘', budget: 2000, actual: null, remark: null },
      { name: '智能开关', budget: 1000, actual: null, remark: null },
      { name: '空调伴侣', budget: 200, actual: null, remark: null },
      { name: '网关', budget: 800, actual: null, remark: null },
      { name: '路由器', budget: 800, actual: null, remark: null },
      { name: '智能马桶', budget: 5000, actual: null, remark: null },
      { name: '人在传感器', budget: 300, actual: null, remark: null },
    ],
  },
  {
    key: 'other',
    name: '其他',
    items: [
      { name: '抽拉式数显水龙头', budget: 700, actual: null, remark: null },
      { name: '验房', budget: 230, actual: null, remark: null },
    ],
  },
]

// ====== 工具函数 ======
/** 金额格式化为 ¥X,XXX */
function formatMoney(value: number): string {
  return '¥' + value.toLocaleString('zh-CN')
}

/** 差额文字 */
function diffText(item: BudgetItem): string {
  if (item.actual === null) return '—'
  const diff = item.budget - item.actual
  if (diff === 0) return '¥0'
  return diff > 0 ? formatMoney(diff) : '-' + formatMoney(Math.abs(diff))
}

/** 差额颜色 class */
function diffClass(item: BudgetItem): string {
  if (item.actual === null) return 'text-gray-300 dark:text-gray-600'
  const diff = item.budget - item.actual
  if (diff > 0) return 'text-emerald-600 dark:text-emerald-400'
  if (diff < 0) return 'text-red-600 dark:text-red-400'
  return 'text-gray-500 dark:text-gray-400'
}

/** 实际支出颜色 class */
function actualClass(item: BudgetItem): string {
  if (item.actual === null) return 'text-gray-300 dark:text-gray-600'
  if (item.actual > item.budget) return 'text-red-600 dark:text-red-400'
  return 'text-emerald-600 dark:text-emerald-400'
}

// ====== 分类样式 ======
function categoryBarClass(key: string): string {
  const map: Record<string, string> = {
    partial: 'bg-blue-500',
    hard: 'bg-orange-500',
    furniture: 'bg-emerald-500',
    appliance: 'bg-violet-500',
    smart: 'bg-rose-500',
    other: 'bg-amber-500',
  }
  return map[key] || 'bg-gray-400'
}

function categoryBadgeClass(key: string): string {
  const map: Record<string, string> = {
    partial: 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300',
    hard: 'bg-orange-50 text-orange-700 dark:bg-orange-900/30 dark:text-orange-300',
    furniture: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300',
    appliance: 'bg-violet-50 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300',
    smart: 'bg-rose-50 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300',
    other: 'bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
  }
  return map[key] || ''
}

function categoryHeaderBgClass(key: string): string {
  const map: Record<string, string> = {
    partial: 'bg-blue-50/60 dark:bg-blue-900/15',
    hard: 'bg-orange-50/60 dark:bg-orange-900/15',
    furniture: 'bg-emerald-50/60 dark:bg-emerald-900/15',
    appliance: 'bg-violet-50/60 dark:bg-violet-900/15',
    smart: 'bg-rose-50/60 dark:bg-rose-900/15',
    other: 'bg-amber-50/60 dark:bg-amber-900/15',
  }
  return map[key] || 'bg-gray-50 dark:bg-gray-800/80'
}

function categorySubtotalBgClass(key: string): string {
  const map: Record<string, string> = {
    partial: 'bg-blue-50/50 dark:bg-blue-900/10',
    hard: 'bg-orange-50/50 dark:bg-orange-900/10',
    furniture: 'bg-emerald-50/50 dark:bg-emerald-900/10',
    appliance: 'bg-violet-50/50 dark:bg-violet-900/10',
    smart: 'bg-rose-50/50 dark:bg-rose-900/10',
    other: 'bg-amber-50/50 dark:bg-amber-900/10',
  }
  return map[key] || ''
}

function categorySubtotalTextClass(key: string): string {
  const map: Record<string, string> = {
    partial: 'text-blue-700 dark:text-blue-300',
    hard: 'text-orange-700 dark:text-orange-300',
    furniture: 'text-emerald-700 dark:text-emerald-300',
    appliance: 'text-violet-700 dark:text-violet-300',
    smart: 'text-rose-700 dark:text-rose-300',
    other: 'text-amber-700 dark:text-amber-300',
  }
  return map[key] || ''
}

// ====== 分类汇总 ======
function categoryBudgetSum(cat: BudgetCategory): number {
  return cat.items.reduce((sum, item) => sum + item.budget, 0)
}

function categoryActualSum(cat: BudgetCategory): number {
  return cat.items.reduce((sum, item) => sum + (item.actual ?? 0), 0)
}

function categoryDiffText(cat: BudgetCategory): string {
  const actualSum = categoryActualSum(cat)
  if (actualSum === 0) return '—'
  const diff = categoryBudgetSum(cat) - actualSum
  if (diff === 0) return '¥0'
  return diff > 0 ? formatMoney(diff) : '-' + formatMoney(Math.abs(diff))
}

/** 分类实际支出颜色 class */
function categoryActualClass(cat: BudgetCategory): string {
  const actualSum = categoryActualSum(cat)
  if (actualSum === 0) return 'text-gray-300 dark:text-gray-600'
  if (actualSum > categoryBudgetSum(cat)) return 'text-red-600 dark:text-red-400'
  return 'text-emerald-600 dark:text-emerald-400'
}

/** 分类差额颜色 class */
function categoryDiffClass(cat: BudgetCategory): string {
  const actualSum = categoryActualSum(cat)
  if (actualSum === 0) return 'text-gray-300 dark:text-gray-600'
  const diff = categoryBudgetSum(cat) - actualSum
  if (diff > 0) return 'text-emerald-600 dark:text-emerald-400'
  if (diff < 0) return 'text-red-600 dark:text-red-400'
  return 'text-gray-500 dark:text-gray-400'
}

// ====== 汇总统计 ======
const summary = computed(() => {
  let totalBudget = 0
  let totalActual = 0

  for (const cat of budgetData) {
    for (const item of cat.items) {
      totalBudget += item.budget
      if (item.actual !== null) {
        totalActual += item.actual
      }
    }
  }

  const totalDiff = totalBudget - totalActual

  return { totalBudget, totalActual, totalDiff }
})

// ====== 预览弹窗 ======
const previewVisible = ref(false)
const previewUrl = ref<string | null>(null)

function closePreview() {
  previewVisible.value = false
  previewUrl.value = null
}

// ====== Canvas 绘制 & 下载 ======
import QRCode from 'qrcode'

const AVATAR_URL = 'https://assets.fatwill.cloud/uploads/1775045182650-e9d2upw9.jpeg'
const QR_URL = 'https://fatwill.cloud/renovation/budget'
const IMG_WIDTH = 900
const DPR = 2 // 高清渲染

// 分类颜色映射（用于 Canvas 绘制）
const categoryColors: Record<string, { main: string; light: string; text: string }> = {
  partial: { main: '#3b82f6', light: '#eff6ff', text: '#1e40af' },
  hard: { main: '#f97316', light: '#fff7ed', text: '#9a3412' },
  furniture: { main: '#10b981', light: '#ecfdf5', text: '#065f46' },
  appliance: { main: '#8b5cf6', light: '#f5f3ff', text: '#5b21b6' },
  smart: { main: '#f43f5e', light: '#fff1f2', text: '#9f1239' },
  other: { main: '#f59e0b', light: '#fffbeb', text: '#92400e' },
}

/** 生成二维码 Canvas */
async function generateQRCode(url: string, size: number): Promise<HTMLCanvasElement> {
  const canvas = document.createElement('canvas')
  await QRCode.toCanvas(canvas, url, {
    width: size,
    margin: 1,
    color: { dark: '#333333', light: '#ffffff' },
  })
  return canvas
}

/** 加载图片（带跨域处理） */
function loadImage(url: string): Promise<HTMLImageElement | null> {
  return new Promise((resolve) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => resolve(img)
    img.onerror = () => resolve(null)
    img.src = url
  })
}

/** 绘制圆形图片 */
function drawCircleImage(ctx: CanvasRenderingContext2D, img: HTMLImageElement | null, x: number, y: number, size: number) {
  ctx.save()
  ctx.beginPath()
  ctx.arc(x + size / 2, y + size / 2, size / 2, 0, Math.PI * 2)
  ctx.closePath()
  ctx.clip()
  if (img) {
    ctx.drawImage(img, x, y, size, size)
  } else {
    // 兜底灰色圆形
    ctx.fillStyle = '#d1d5db'
    ctx.fillRect(x, y, size, size)
  }
  ctx.restore()
}

/** 绘制圆角矩形 */
function drawRoundRect(ctx: CanvasRenderingContext2D, x: number, y: number, w: number, h: number, r: number) {
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.lineTo(x + w - r, y)
  ctx.arcTo(x + w, y, x + w, y + r, r)
  ctx.lineTo(x + w, y + h - r)
  ctx.arcTo(x + w, y + h, x + w - r, y + h, r)
  ctx.lineTo(x + r, y + h)
  ctx.arcTo(x, y + h, x, y + h - r, r)
  ctx.lineTo(x, y + r)
  ctx.arcTo(x, y, x + r, y, r)
  ctx.closePath()
}

/** 纯金额数字（不含¥） */
function moneyNum(value: number): string {
  return value.toLocaleString('zh-CN')
}

/** 核心绘制函数 */
async function generateBudgetImage(categories: BudgetCategory[], title: string, subtitle: string): Promise<void> {
  // 预加载资源
  const [avatarImg, qrCanvas] = await Promise.all([
    loadImage(AVATAR_URL),
    generateQRCode(QR_URL, 80 * DPR),
  ])

  // 计算画布高度
  const headerH = 80
  const catTitleH = 40
  const tableHeaderH = 32
  const rowH = 30
  const catGap = 20
  const footerH = 100
  const padding = 32

  let contentH = headerH + padding
  for (const cat of categories) {
    contentH += catTitleH + tableHeaderH + cat.items.length * rowH + rowH /* 小计 */ + catGap
  }
  contentH += footerH + padding

  const canvas = document.createElement('canvas')
  canvas.width = IMG_WIDTH * DPR
  canvas.height = contentH * DPR
  const ctx = canvas.getContext('2d')!
  ctx.scale(DPR, DPR)

  // 白色背景
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, IMG_WIDTH, contentH)

  // ---- Header ----
  ctx.fillStyle = '#f8f9fa'
  ctx.fillRect(0, 0, IMG_WIDTH, headerH)
  // 底部分割线
  ctx.strokeStyle = '#e5e7eb'
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(0, headerH)
  ctx.lineTo(IMG_WIDTH, headerH)
  ctx.stroke()

  // 头像
  const avatarSize = 48
  const avatarX = padding
  const avatarY = (headerH - avatarSize) / 2
  drawCircleImage(ctx, avatarImg, avatarX, avatarY, avatarSize)

  // 标题
  ctx.fillStyle = '#111827'
  ctx.font = 'bold 18px "PingFang SC", "Microsoft YaHei", sans-serif'
  ctx.textBaseline = 'middle'
  ctx.fillText(title, avatarX + avatarSize + 14, headerH / 2 - 10)

  // 副标题
  ctx.fillStyle = '#9ca3af'
  ctx.font = '12px "PingFang SC", "Microsoft YaHei", sans-serif'
  ctx.fillText(subtitle, avatarX + avatarSize + 14, headerH / 2 + 12)

  // 日期
  const today = new Date().toISOString().slice(0, 10)
  ctx.fillStyle = '#9ca3af'
  ctx.font = '12px "PingFang SC", "Microsoft YaHei", sans-serif'
  ctx.textAlign = 'right'
  ctx.fillText(today, IMG_WIDTH - padding, headerH / 2)
  ctx.textAlign = 'left'

  // ---- 表格区域 ----
  let curY = headerH + padding
  const colWidths = [240, 130, 130, 130, 206] // 项目名称、预算、实际支出、差额、备注
  const tableX = padding
  const tableW = IMG_WIDTH - padding * 2

  for (const cat of categories) {
    const colors = categoryColors[cat.key] || categoryColors.other

    // 分类标题行
    drawRoundRect(ctx, tableX, curY, tableW, catTitleH, 8)
    ctx.fillStyle = colors.light
    ctx.fill()
    // 左侧色条
    ctx.fillStyle = colors.main
    drawRoundRect(ctx, tableX, curY, 4, catTitleH, 2)
    ctx.fill()
    // 分类名
    ctx.fillStyle = colors.text
    ctx.font = 'bold 14px "PingFang SC", "Microsoft YaHei", sans-serif'
    ctx.textBaseline = 'middle'
    ctx.fillText(cat.name, tableX + 16, curY + catTitleH / 2)
    // 小计金额
    const catBudget = categoryBudgetSum(cat)
    ctx.fillStyle = '#6b7280'
    ctx.font = '12px "PingFang SC", "Microsoft YaHei", sans-serif'
    ctx.textAlign = 'right'
    ctx.fillText(`${cat.items.length} 项 · 预算 ¥${moneyNum(catBudget)}`, tableX + tableW - 12, curY + catTitleH / 2)
    ctx.textAlign = 'left'
    curY += catTitleH

    // 表头
    ctx.fillStyle = '#f9fafb'
    ctx.fillRect(tableX, curY, tableW, tableHeaderH)
    ctx.strokeStyle = '#e5e7eb'
    ctx.lineWidth = 0.5
    ctx.beginPath()
    ctx.moveTo(tableX, curY + tableHeaderH)
    ctx.lineTo(tableX + tableW, curY + tableHeaderH)
    ctx.stroke()

    const headers = ['项目名称', '预算', '实际支出', '差额', '备注']
    ctx.fillStyle = '#6b7280'
    ctx.font = 'bold 11px "PingFang SC", "Microsoft YaHei", sans-serif'
    ctx.textBaseline = 'middle'
    let colX = tableX + 12
    for (let i = 0; i < headers.length; i++) {
      if (i >= 1 && i <= 3) {
        ctx.textAlign = 'right'
        ctx.fillText(headers[i], colX + colWidths[i] - 12, curY + tableHeaderH / 2)
        ctx.textAlign = 'left'
      } else {
        ctx.fillText(headers[i], colX, curY + tableHeaderH / 2)
      }
      colX += colWidths[i]
    }
    curY += tableHeaderH

    // 数据行
    for (let ri = 0; ri < cat.items.length; ri++) {
      const item = cat.items[ri]
      // 斑马纹
      if (ri % 2 === 1) {
        ctx.fillStyle = '#f9fafb'
        ctx.fillRect(tableX, curY, tableW, rowH)
      }
      // 底部线
      ctx.strokeStyle = '#f3f4f6'
      ctx.lineWidth = 0.5
      ctx.beginPath()
      ctx.moveTo(tableX, curY + rowH)
      ctx.lineTo(tableX + tableW, curY + rowH)
      ctx.stroke()

      ctx.font = '12px "PingFang SC", "Microsoft YaHei", sans-serif'
      ctx.textBaseline = 'middle'
      const midY = curY + rowH / 2
      colX = tableX + 12

      // 项目名称
      ctx.fillStyle = '#111827'
      ctx.fillText(item.name, colX, midY)
      colX += colWidths[0]

      // 预算
      ctx.fillStyle = '#374151'
      ctx.textAlign = 'right'
      ctx.fillText(`¥${moneyNum(item.budget)}`, colX + colWidths[1] - 12, midY)
      colX += colWidths[1]

      // 实际支出
      if (item.actual !== null) {
        ctx.fillStyle = item.actual > item.budget ? '#dc2626' : '#059669'
        ctx.fillText(`¥${moneyNum(item.actual)}`, colX + colWidths[2] - 12, midY)
      } else {
        ctx.fillStyle = '#d1d5db'
        ctx.fillText('—', colX + colWidths[2] - 12, midY)
      }
      colX += colWidths[2]

      // 差额
      if (item.actual !== null) {
        const diff = item.budget - item.actual
        if (diff > 0) {
          ctx.fillStyle = '#059669'
          ctx.fillText(`¥${moneyNum(diff)}`, colX + colWidths[3] - 12, midY)
        } else if (diff < 0) {
          ctx.fillStyle = '#dc2626'
          ctx.fillText(`-¥${moneyNum(Math.abs(diff))}`, colX + colWidths[3] - 12, midY)
        } else {
          ctx.fillStyle = '#6b7280'
          ctx.fillText('¥0', colX + colWidths[3] - 12, midY)
        }
      } else {
        ctx.fillStyle = '#d1d5db'
        ctx.fillText('—', colX + colWidths[3] - 12, midY)
      }
      colX += colWidths[3]

      // 备注
      ctx.textAlign = 'left'
      ctx.fillStyle = '#9ca3af'
      ctx.fillText(item.remark || '—', colX, midY)

      curY += rowH
    }

    // 小计行
    ctx.fillStyle = colors.light
    ctx.fillRect(tableX, curY, tableW, rowH)
    ctx.strokeStyle = '#e5e7eb'
    ctx.lineWidth = 0.5
    ctx.beginPath()
    ctx.moveTo(tableX, curY)
    ctx.lineTo(tableX + tableW, curY)
    ctx.stroke()

    ctx.font = 'bold 12px "PingFang SC", "Microsoft YaHei", sans-serif'
    ctx.textBaseline = 'middle'
    const subMidY = curY + rowH / 2
    colX = tableX + 12

    // 小计文字
    ctx.fillStyle = colors.text
    ctx.fillText('小计', colX, subMidY)
    colX += colWidths[0]

    // 预算小计
    ctx.fillStyle = '#374151'
    ctx.textAlign = 'right'
    ctx.fillText(`¥${moneyNum(catBudget)}`, colX + colWidths[1] - 12, subMidY)
    colX += colWidths[1]

    // 实际支出小计
    const catActual = categoryActualSum(cat)
    if (catActual > 0) {
      ctx.fillStyle = catActual > catBudget ? '#dc2626' : '#059669'
      ctx.fillText(`¥${moneyNum(catActual)}`, colX + colWidths[2] - 12, subMidY)
    } else {
      ctx.fillStyle = '#d1d5db'
      ctx.fillText('—', colX + colWidths[2] - 12, subMidY)
    }
    colX += colWidths[2]

    // 差额小计
    if (catActual > 0) {
      const catDiff = catBudget - catActual
      if (catDiff > 0) {
        ctx.fillStyle = '#059669'
        ctx.fillText(`¥${moneyNum(catDiff)}`, colX + colWidths[3] - 12, subMidY)
      } else if (catDiff < 0) {
        ctx.fillStyle = '#dc2626'
        ctx.fillText(`-¥${moneyNum(Math.abs(catDiff))}`, colX + colWidths[3] - 12, subMidY)
      } else {
        ctx.fillStyle = '#6b7280'
        ctx.fillText('¥0', colX + colWidths[3] - 12, subMidY)
      }
    } else {
      ctx.fillStyle = '#d1d5db'
      ctx.fillText('—', colX + colWidths[3] - 12, subMidY)
    }
    ctx.textAlign = 'left'

    curY += rowH + catGap
  }

  // ---- Footer ----
  const footerY = curY
  ctx.fillStyle = '#f8f9fa'
  ctx.fillRect(0, footerY, IMG_WIDTH, footerH)
  ctx.strokeStyle = '#e5e7eb'
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(0, footerY)
  ctx.lineTo(IMG_WIDTH, footerY)
  ctx.stroke()

  // 总预算 / 已支出
  let totalBudget = 0
  let totalActual = 0
  for (const cat of categories) {
    totalBudget += categoryBudgetSum(cat)
    totalActual += categoryActualSum(cat)
  }

  ctx.font = 'bold 14px "PingFang SC", "Microsoft YaHei", sans-serif'
  ctx.textBaseline = 'middle'
  ctx.fillStyle = '#111827'
  ctx.fillText(`总预算  ¥${moneyNum(totalBudget)}`, padding, footerY + footerH / 2 - 12)

  ctx.fillStyle = totalActual > totalBudget ? '#dc2626' : '#059669'
  ctx.font = 'bold 13px "PingFang SC", "Microsoft YaHei", sans-serif'
  ctx.fillText(`已支出  ¥${moneyNum(totalActual)}`, padding, footerY + footerH / 2 + 12)

  // 二维码
  const qrSize = 70
  const qrX = IMG_WIDTH - padding - qrSize
  const qrY = footerY + (footerH - qrSize - 14) / 2
  ctx.drawImage(qrCanvas, qrX, qrY, qrSize, qrSize)

  // 二维码下方文字
  ctx.fillStyle = '#9ca3af'
  ctx.font = '9px "PingFang SC", "Microsoft YaHei", sans-serif'
  ctx.textAlign = 'center'
  ctx.fillText('扫码查看详情', qrX + qrSize / 2, qrY + qrSize + 11)
  ctx.textAlign = 'left'

  // 显示预览弹窗
  previewUrl.value = canvas.toDataURL('image/png')
  previewVisible.value = true
}

/** 下载全部 */
async function downloadAll(): Promise<void> {
  await generateBudgetImage(budgetData, 'fatwill 的预算', '装修成本预算清单')
}

/** 下载单个分类 */
async function downloadCategory(cat: BudgetCategory): Promise<void> {
  await generateBudgetImage([cat], `fatwill 的预算 · ${cat.name}`, `${cat.name}成本明细`)
}

// ====== SEO ======
useSeoMeta({
  title: '成本预算 - 装修',
  description: '装修成本预算明细表，包含局改、硬装、家私、家电、全屋智能、其他六大分类',
})
</script>

<style scoped>
.preview-fade-enter-active,
.preview-fade-leave-active {
  transition: opacity 0.25s ease;
}
.preview-fade-enter-from,
.preview-fade-leave-to {
  opacity: 0;
}
</style>