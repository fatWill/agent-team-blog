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
        <div class="rounded-xl border border-gray-100 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
          <div class="mb-1 text-xs text-gray-500 dark:text-gray-400">总预算</div>
          <div class="text-lg font-bold text-gray-900 dark:text-gray-100 sm:text-xl">{{ formatMoney(summary.totalBudget) }}</div>
        </div>
        <!-- 已支出 -->
        <div class="rounded-xl border border-gray-100 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
          <div class="mb-1 text-xs text-gray-500 dark:text-gray-400">已支出</div>
          <div class="text-lg font-bold text-gray-900 dark:text-gray-100 sm:text-xl">{{ formatMoney(summary.totalActual) }}</div>
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

// ====== SEO ======
useSeoMeta({
  title: '成本预算 - 装修',
  description: '装修成本预算明细表，包含局改、硬装、家私、家电、全屋智能五大分类',
})
</script>