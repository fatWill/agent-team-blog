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
        <!-- 右侧操作按钮 -->
        <div class="ml-auto flex shrink-0 items-center gap-2">
          <template v-if="editing">
            <button
              class="flex shrink-0 items-center gap-1 rounded-lg bg-indigo-500 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-indigo-600 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="saving"
              @click="handleSave"
            >
              <svg v-if="saving" class="h-3.5 w-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
              </svg>
              保存
            </button>
            <button
              class="shrink-0 rounded-lg border border-gray-300 px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-100 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
              @click="handleCancel"
            >
              取消
            </button>
          </template>
          <button
            v-else-if="isLoggedIn"
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-indigo-500 transition-colors hover:bg-indigo-50 dark:hover:bg-indigo-900/20"
            title="编辑"
            @click="enterEditMode"
          >
            <svg class="h-[18px] w-[18px]" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <main class="mx-auto w-full max-w-5xl flex-1 px-4 py-6">
      <!-- 汇总统计卡片 -->
      <div class="mb-6 grid grid-cols-2 gap-3 sm:gap-4">
        <div class="rounded-xl border border-indigo-200 bg-gradient-to-br from-indigo-50 to-white p-4 shadow-sm dark:border-indigo-800 dark:from-indigo-950/40 dark:to-gray-800">
          <div class="mb-1 flex items-center gap-1.5 text-xs text-indigo-500 dark:text-indigo-400">
            <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
            </svg>
            总预算
          </div>
          <div class="text-lg font-bold text-indigo-700 dark:text-indigo-300 sm:text-xl">{{ formatMoney(computedSummary.totalBudget) }}</div>
        </div>
        <div
          class="rounded-xl border p-4 shadow-sm"
          :class="computedSummary.totalActual > computedSummary.totalBudget
            ? 'border-red-200 bg-gradient-to-br from-red-50 to-white dark:border-red-800 dark:from-red-950/40 dark:to-gray-800'
            : 'border-emerald-200 bg-gradient-to-br from-emerald-50 to-white dark:border-emerald-800 dark:from-emerald-950/40 dark:to-gray-800'"
        >
          <div
            class="mb-1 flex items-center gap-1.5 text-xs"
            :class="computedSummary.totalActual > computedSummary.totalBudget ? 'text-red-500 dark:text-red-400' : 'text-emerald-500 dark:text-emerald-400'"
          >
            <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            已支出
          </div>
          <div
            class="text-lg font-bold sm:text-xl"
            :class="computedSummary.totalActual > computedSummary.totalBudget ? 'text-red-600 dark:text-red-400' : 'text-emerald-600 dark:text-emerald-400'"
          >
            {{ formatMoney(computedSummary.totalActual) }}
          </div>
        </div>
      </div>

      <!-- 每个大类独立一个表格 -->
      <div class="space-y-5">
        <div v-for="category in displayCategories" :key="category.key">
          <div class="mb-2 flex items-center gap-2">
            <span class="inline-block h-5 w-1 rounded-full" :class="categoryBarClass(category.key)" />
            <span class="inline-flex items-center rounded-md px-2.5 py-1 text-sm font-bold" :class="categoryBadgeClass(category.key)">
              {{ category.name }}
            </span>
            <span class="text-xs text-gray-400 dark:text-gray-500">
              （{{ category.items.length }} 项，预算 {{ formatMoney(categorySumAmount(category)) }}）
            </span>
            <button
              v-if="!editing"
              class="ml-auto flex h-6 w-6 items-center justify-center rounded-md text-gray-300 transition-colors hover:bg-gray-100 hover:text-gray-500 dark:text-gray-600 dark:hover:bg-gray-700 dark:hover:text-gray-300"
              title="下载该分类截图"
              @click="downloadCategory(category)"
            >
              <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
              </svg>
            </button>
          </div>
          <div class="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <div class="overflow-x-auto">
              <table class="w-full min-w-[540px] text-sm">
                <thead>
                  <tr class="border-b border-gray-200 dark:border-gray-700" :class="categoryHeaderBgClass(category.key)">
                    <th class="whitespace-nowrap px-3 py-2.5 text-left font-semibold text-gray-600 dark:text-gray-300 sm:px-4">项目名称</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-right font-semibold text-gray-600 dark:text-gray-300 sm:px-4">预算金额</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-right font-semibold text-gray-600 dark:text-gray-300 sm:px-4">实际支出</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-right font-semibold text-gray-600 dark:text-gray-300 sm:px-4">差额</th>
                    <th class="whitespace-nowrap px-3 py-2.5 text-left font-semibold text-gray-600 dark:text-gray-300 sm:px-4">备注</th>
                    <th v-if="editing" class="w-10 px-2 py-2.5" />
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(item, itemIdx) in category.items"
                    :key="item._uid"
                    class="border-b border-gray-100 transition-colors hover:bg-gray-50/50 dark:border-gray-700/50 dark:hover:bg-gray-700/30"
                    :class="itemIdx % 2 === 1 ? 'bg-gray-50/30 dark:bg-gray-800/50' : ''"
                  >
                    <template v-if="!editing">
                      <td class="whitespace-nowrap px-3 py-2.5 font-medium text-gray-900 dark:text-gray-100 sm:px-4">{{ item.itemName }}</td>
                      <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono text-gray-700 dark:text-gray-300 sm:px-4">{{ formatMoney(item.amount) }}</td>
                      <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono sm:px-4" :class="actualClassByItem(item)">{{ item.actual != null ? formatMoney(item.actual) : '—' }}</td>
                      <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono sm:px-4" :class="diffClassByItem(item)">{{ diffTextByItem(item) }}</td>
                      <td class="px-3 py-2.5 text-gray-500 dark:text-gray-400 sm:px-4">{{ item.remark || '—' }}</td>
                    </template>
                    <template v-else>
                      <td class="px-2 py-1.5 sm:px-3">
                        <input v-model="item.itemName" class="w-full rounded border border-gray-200 px-2 py-1 text-sm text-gray-900 outline-none focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100" placeholder="项目名称" />
                      </td>
                      <td class="px-2 py-1.5 sm:px-3">
                        <input v-model.number="item.amount" type="number" min="0" step="0.01" class="w-full rounded border border-gray-200 px-2 py-1 text-right text-sm font-mono text-gray-900 outline-none focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100" placeholder="0" />
                      </td>
                      <td class="px-2 py-1.5 sm:px-3">
                        <input v-model.number="item.actual" type="number" min="0" step="0.01" class="w-full rounded border border-gray-200 px-2 py-1 text-right text-sm font-mono text-gray-900 outline-none focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100" placeholder="—" />
                      </td>
                      <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono sm:px-4" :class="diffClassByItem(item)">{{ diffTextByItem(item) }}</td>
                      <td class="px-2 py-1.5 sm:px-3">
                        <input v-model="item.remark" class="w-full rounded border border-gray-200 px-2 py-1 text-sm text-gray-900 outline-none focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100" placeholder="备注" />
                      </td>
                      <td class="px-2 py-1.5">
                        <button class="flex h-7 w-7 items-center justify-center rounded text-red-400 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20" title="删除此行" @click="removeDraftRow(item._uid)">
                          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                          </svg>
                        </button>
                      </td>
                    </template>
                  </tr>
                  <tr v-if="editing">
                    <td :colspan="6" class="px-3 py-2.5">
                      <button class="flex items-center gap-1.5 rounded-lg px-3 py-2 text-xs font-medium text-emerald-600 transition-colors hover:bg-emerald-50 dark:text-emerald-400 dark:hover:bg-emerald-900/20" @click="addDraftRow(category.key)">
                        <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
                        新增一行
                      </button>
                    </td>
                  </tr>
                  <tr class="border-t border-gray-200 dark:border-gray-600" :class="categorySubtotalBgClass(category.key)">
                    <td class="px-3 py-2.5 sm:px-4"><span class="text-xs font-bold" :class="categorySubtotalTextClass(category.key)">小计</span></td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono font-bold text-gray-700 dark:text-gray-300 sm:px-4">{{ formatMoney(categorySumAmount(category)) }}</td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono font-bold sm:px-4" :class="categoryActualClass(category)">{{ categorySumActual(category) > 0 ? formatMoney(categorySumActual(category)) : '—' }}</td>
                    <td class="whitespace-nowrap px-3 py-2.5 text-right font-mono font-bold sm:px-4" :class="categoryDiffClass(category)">{{ categoryDiffText(category) }}</td>
                    <td class="px-3 py-2.5 sm:px-4" />
                    <td v-if="editing" />
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
          <table class="w-full min-w-[540px] text-sm">
            <tbody>
              <tr>
                <td class="px-3 py-3 font-bold text-gray-900 dark:text-gray-100 sm:px-4">总计</td>
                <td class="whitespace-nowrap px-3 py-3 text-right font-mono font-bold text-gray-900 dark:text-gray-100 sm:px-4">{{ formatMoney(computedSummary.totalBudget) }}</td>
                <td class="whitespace-nowrap px-3 py-3 text-right font-mono font-bold sm:px-4" :class="computedSummary.totalActual > computedSummary.totalBudget ? 'text-red-600 dark:text-red-400' : 'text-emerald-600 dark:text-emerald-400'">{{ formatMoney(computedSummary.totalActual) }}</td>
                <td class="whitespace-nowrap px-3 py-3 text-right font-mono font-bold sm:px-4" :class="computedSummary.totalDiff >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'">{{ computedSummary.totalDiff >= 0 ? '-' + formatMoney(computedSummary.totalDiff) : '+' + formatMoney(Math.abs(computedSummary.totalDiff)) }}</td>
                <td class="px-3 py-3 sm:px-4" />
                <td v-if="editing" />
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <p class="mt-4 text-center text-xs text-gray-400 dark:text-gray-500">差额 = 预算 − 实际支出 · 绿色为节省 · 红色为超支</p>
    </main>

    <!-- 右下角悬浮按钮 -->
    <div v-if="!editing" class="fixed bottom-6 right-6 z-50 flex flex-col gap-3">
      <button class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-800 text-white shadow-lg transition-all hover:scale-110 hover:bg-gray-700 dark:bg-gray-200 dark:text-gray-800 dark:hover:bg-gray-300" title="分享" @click="downloadAll">
        <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z" /></svg>
      </button>
      <button class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-800 text-white shadow-lg transition-all hover:scale-110 hover:bg-gray-700 dark:bg-gray-200 dark:text-gray-800 dark:hover:bg-gray-300" title="下载预算图片" @click="downloadAll">
        <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" /></svg>
      </button>
    </div>

    <!-- 图片预览弹窗 -->
    <Teleport to="body">
      <Transition name="preview-fade">
        <div v-if="previewVisible" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/70 backdrop-blur-sm p-6">
          <div class="relative flex max-h-[90vh] max-w-[92vw] flex-col rounded-2xl bg-white shadow-2xl dark:bg-gray-800 overflow-hidden">
            <button class="absolute right-3 top-3 z-10 flex h-8 w-8 items-center justify-center rounded-full bg-black/40 text-white backdrop-blur-sm transition hover:bg-black/60" @click="closePreview">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
            </button>
            <div class="overflow-y-auto">
              <img v-if="previewUrl" :src="previewUrl" class="block w-full" alt="预算图片预览" draggable="false" />
              <div v-else class="flex h-48 items-center justify-center">
                <svg class="h-8 w-8 animate-spin text-gray-400" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" /><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" /></svg>
              </div>
            </div>
            <div class="flex items-center justify-center gap-1.5 border-t border-gray-100 px-4 py-2.5 dark:border-gray-700">
              <svg class="h-3.5 w-3.5 flex-shrink-0 text-gray-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 18v-5.25m0 0a6.01 6.01 0 001.5-.189m-1.5.189a6.01 6.01 0 01-1.5-.189m3.75 7.478a12.06 12.06 0 01-4.5 0m3.75 2.383a14.406 14.406 0 01-3 0M14.25 18v-.192c0-.983.658-1.823 1.508-2.316a7.5 7.5 0 10-7.517 0c.85.493 1.509 1.333 1.509 2.316V18" /></svg>
              <span class="text-xs text-gray-400 dark:text-gray-500">长按图片保存到相册</span>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { showSuccess, showError, showConfirm } from '~/shared/utils/ui'
import QRCode from 'qrcode'

// ====== 类型 ======
interface BudgetItemAPI { id?: number; category: string; itemName: string; amount: number; actual?: number | null; remark: string; sortOrder: number }
interface DraftItem { _uid: string; id?: number; category: string; itemName: string; amount: number; actual: number | null; remark: string; sortOrder: number }
interface DisplayCategory { key: string; name: string; items: DraftItem[] }

// ====== 常量 ======
const CATEGORY_ORDER = ['局改', '硬装', '家私', '家电', '全屋智能', '其他', '额外支出']
const CATEGORY_KEY_MAP: Record<string, string> = { '局改': 'partial', '硬装': 'hard', '家私': 'furniture', '家电': 'appliance', '全屋智能': 'smart', '其他': 'other', '额外支出': 'extra' }

const HARDCODED_DATA: DraftItem[] = [
  { _uid: 'h1', category: '局改', itemName: '保温层拆除+修复', amount: 2500, actual: 2500, remark: '杨工', sortOrder: 1 },
  { _uid: 'h2', category: '局改', itemName: '厨房拆墙+修复', amount: 6000, actual: 6000, remark: '杨工', sortOrder: 2 },
  { _uid: 'h3', category: '局改', itemName: '全屋零线+水电', amount: 2000, actual: null, remark: '杨工', sortOrder: 3 },
  { _uid: 'h4', category: '局改', itemName: '瓷砖', amount: 1000, actual: null, remark: '', sortOrder: 4 },
  { _uid: 'h5', category: '硬装', itemName: '艺术漆', amount: 6800, actual: 6800, remark: '光明京东家居（义峰）', sortOrder: 5 },
  { _uid: 'h6', category: '硬装', itemName: '美缝', amount: 2000, actual: 1700, remark: '', sortOrder: 6 },
  { _uid: 'h7', category: '硬装', itemName: '全屋定制', amount: 35000, actual: null, remark: '', sortOrder: 7 },
  { _uid: 'h8', category: '硬装', itemName: '阳台柜+防护网+花台', amount: 5500, actual: 5800, remark: '', sortOrder: 8 },
  { _uid: 'h9', category: '硬装', itemName: '完美系统三联轨', amount: 4000, actual: 3700, remark: '', sortOrder: 9 },
  { _uid: 'h10', category: '硬装', itemName: '全屋纱窗', amount: 1500, actual: 1500, remark: '', sortOrder: 10 },
  { _uid: 'h11', category: '硬装', itemName: '换门', amount: 800, actual: 0, remark: '', sortOrder: 11 },
  { _uid: 'h12', category: '家私', itemName: '沙发', amount: 4000, actual: 3230.85, remark: '', sortOrder: 12 },
  { _uid: 'h13', category: '家私', itemName: '床垫×3', amount: 5000, actual: null, remark: '', sortOrder: 13 },
  { _uid: 'h14', category: '家私', itemName: '次卧小桌子', amount: 1500, actual: 0, remark: '', sortOrder: 14 },
  { _uid: 'h15', category: '家私', itemName: '主卧大桌子', amount: 2000, actual: 0, remark: '', sortOrder: 15 },
  { _uid: 'h16', category: '家私', itemName: '电视地柜', amount: 2000, actual: 832.37, remark: '', sortOrder: 16 },
  { _uid: 'h17', category: '家私', itemName: '餐厅桌椅', amount: 2500, actual: 3596, remark: '', sortOrder: 17 },
  { _uid: 'h18', category: '家私', itemName: '窗帘', amount: 2000, actual: null, remark: '', sortOrder: 18 },
  { _uid: 'h19', category: '家私', itemName: '主卧椅子', amount: 1000, actual: 1974.31, remark: '', sortOrder: 19 },
  { _uid: 'h20', category: '家私', itemName: '次卧椅子', amount: 200, actual: null, remark: '', sortOrder: 20 },
  { _uid: 'h21', category: '家电', itemName: '电视机', amount: 3500, actual: null, remark: '', sortOrder: 21 },
  { _uid: 'h22', category: '家电', itemName: '洗衣机', amount: 5000, actual: 7200, remark: '', sortOrder: 22 },
  { _uid: 'h23', category: '家电', itemName: '冰箱', amount: 6500, actual: 5699, remark: '', sortOrder: 23 },
  { _uid: 'h24', category: '家电', itemName: '洗碗机', amount: 8000, actual: 7500, remark: '', sortOrder: 24 },
  { _uid: 'h25', category: '家电', itemName: '直饮机', amount: 8000, actual: 7800, remark: '', sortOrder: 25 },
  { _uid: 'h26', category: '家电', itemName: '零冷水', amount: 1000, actual: 0, remark: '回流泵', sortOrder: 26 },
  { _uid: 'h27', category: '家电', itemName: '全屋净水', amount: 2000, actual: 0, remark: '大白瓶', sortOrder: 27 },
  { _uid: 'h28', category: '家电', itemName: '扫地机', amount: 5500, actual: 4200, remark: '', sortOrder: 28 },
  { _uid: 'h29', category: '家电', itemName: '智能锁', amount: 2500, actual: 1660.80, remark: '', sortOrder: 29 },
  { _uid: 'h30', category: '家电', itemName: '电动晾衣架', amount: 1500, actual: 1040.82, remark: '', sortOrder: 30 },
  { _uid: 'h31', category: '家电', itemName: '智能马桶', amount: 5000, actual: 2616.60, remark: '', sortOrder: 31 },
  { _uid: 'h32', category: '全屋智能', itemName: '智能灯', amount: 3000, actual: 3097.2, remark: '', sortOrder: 32 },
  { _uid: 'h33', category: '全屋智能', itemName: '电动窗帘', amount: 2000, actual: 1460.30, remark: '', sortOrder: 33 },
  { _uid: 'h34', category: '全屋智能', itemName: '智能开关', amount: 1000, actual: 462.80, remark: '', sortOrder: 34 },
  { _uid: 'h35', category: '全屋智能', itemName: '空调伴侣', amount: 200, actual: 0, remark: '', sortOrder: 35 },
  { _uid: 'h36', category: '全屋智能', itemName: '网关', amount: 800, actual: 0, remark: '', sortOrder: 36 },
  { _uid: 'h37', category: '全屋智能', itemName: '路由器', amount: 800, actual: 738.99, remark: '', sortOrder: 37 },
  { _uid: 'h38', category: '全屋智能', itemName: '人在传感器', amount: 300, actual: null, remark: '', sortOrder: 38 },
  { _uid: 'h39', category: '其他', itemName: '抽拉式数显水龙头', amount: 700, actual: null, remark: '', sortOrder: 39 },
  { _uid: 'h40', category: '其他', itemName: '验房', amount: 230, actual: 210, remark: '', sortOrder: 40 },
  { _uid: 'h41', category: '其他', itemName: '淋浴系统', amount: 2500, actual: 1082.24, remark: '', sortOrder: 41 },
  { _uid: 'h42', category: '额外支出', itemName: '公牛轨道插座', amount: 0, actual: 1008.72, remark: '', sortOrder: 42 },
  { _uid: 'h43', category: '额外支出', itemName: '洗衣液', amount: 0, actual: 100.20, remark: '', sortOrder: 43 },
  { _uid: 'h44', category: '额外支出', itemName: '灯带', amount: 0, actual: 168.30, remark: '', sortOrder: 44 },
  { _uid: 'h45', category: '额外支出', itemName: '小米室外摄像头4', amount: 0, actual: 190.33, remark: '', sortOrder: 45 },
  { _uid: 'h46', category: '额外支出', itemName: '小米摄像头4pro', amount: 0, actual: 247.07, remark: '', sortOrder: 46 },
  { _uid: 'h47', category: '额外支出', itemName: '洗碗粉三件套', amount: 0, actual: 207.29, remark: '', sortOrder: 47 },
  { _uid: 'h48', category: '额外支出', itemName: '纸巾', amount: 0, actual: 169.75, remark: '', sortOrder: 48 },
  { _uid: 'h49', category: '额外支出', itemName: 'UPS', amount: 0, actual: 385.23, remark: '', sortOrder: 49 },
  { _uid: 'h50', category: '额外支出', itemName: 'NAS', amount: 0, actual: 2238.11, remark: '', sortOrder: 50 },
  { _uid: 'h51', category: '额外支出', itemName: '硬盘4T', amount: 0, actual: 1262.80, remark: '', sortOrder: 51 },
  { _uid: 'h52', category: '额外支出', itemName: '冰箱内嵌插座', amount: 0, actual: 21.60, remark: '', sortOrder: 52 },
  { _uid: 'h53', category: '额外支出', itemName: '洗手台保护性拆除', amount: 0, actual: 100, remark: '', sortOrder: 53 },
  { _uid: 'h54', category: '额外支出', itemName: '铝槽', amount: 0, actual: 70, remark: '', sortOrder: 54 },
  { _uid: 'h55', category: '额外支出', itemName: '箭牌安装', amount: 0, actual: 40, remark: '', sortOrder: 55 },
  { _uid: 'h56', category: '额外支出', itemName: '本高穿线盒', amount: 0, actual: 58.26, remark: '', sortOrder: 56 },
]

// ====== 状态 ======
const isLoggedIn = ref(false)
const editing = ref(false)
const saving = ref(false)
const migrated = ref(false)
const originalItems = ref<DraftItem[]>([])
const draftItems = ref<DraftItem[]>([])

let uidCounter = 0
function genUid(): string { return `uid_${++uidCounter}_${Date.now()}` }

function apiItemsToDraft(items: BudgetItemAPI[]): DraftItem[] {
  return items.map(item => ({ _uid: genUid(), id: item.id, category: item.category, itemName: item.itemName, amount: item.amount, actual: item.actual ?? null, remark: item.remark || '', sortOrder: item.sortOrder }))
}
function deepCloneItems(items: DraftItem[]): DraftItem[] {
  return items.map(item => ({ ...item, _uid: genUid() }))
}

// SSR-safe: 用硬编码做初始渲染数据
originalItems.value = deepCloneItems(HARDCODED_DATA)
draftItems.value = deepCloneItems(HARDCODED_DATA)

onMounted(async () => {
  try { await $fetch('/api/auth/check'); isLoggedIn.value = true } catch { isLoggedIn.value = false }
  try {
    const res = await $fetch<{ items: BudgetItemAPI[] }>('/api/renovation/budget/items')
    if (res.items && res.items.length > 0) { originalItems.value = apiItemsToDraft(res.items); migrated.value = true }
    else { originalItems.value = deepCloneItems(HARDCODED_DATA); migrated.value = false }
  } catch { originalItems.value = deepCloneItems(HARDCODED_DATA); migrated.value = false }
  draftItems.value = deepCloneItems(originalItems.value)

  // 一次性数据回填：仅登录态 + 已迁移 + sessionStorage 未标记时执行
  if (isLoggedIn.value && migrated.value) {
    const RESTORE_KEY = 'budget-actual-restored-v1'
    if (!sessionStorage.getItem(RESTORE_KEY)) {
      // 构建硬编码 actual 映射（key = itemName, value = actual）
      const hardcodedActualMap = new Map<string, number>()
      for (const hd of HARDCODED_DATA) {
        if (hd.actual !== null && hd.actual > 0) {
          hardcodedActualMap.set(hd.itemName, hd.actual)
        }
      }
      // 找出接口返回 actual 为 0 或 null，但硬编码里有非 0 actual 的项
      const needRestore: { id: number; category: string; item_name: string; amount: number; actual: number; remark: string; sort_order: number }[] = []
      for (const item of originalItems.value) {
        if (item.id && (item.actual === null || item.actual === 0)) {
          const hardActual = hardcodedActualMap.get(item.itemName)
          if (hardActual && hardActual > 0) {
            needRestore.push({
              id: item.id,
              category: item.category,
              item_name: item.itemName,
              amount: item.amount,
              actual: hardActual,
              remark: item.remark,
              sort_order: item.sortOrder,
            })
          }
        }
      }
      if (needRestore.length > 0) {
        try {
          const restoreRes = await $fetch<{ code: number; data: { items: BudgetItemAPI[] } }>(
            '/api/renovation/budget/items/batch',
            { method: 'POST', body: { creates: [], updates: needRestore, deletes: [] } }
          )
          if (restoreRes.data?.items) {
            originalItems.value = apiItemsToDraft(restoreRes.data.items)
            draftItems.value = deepCloneItems(originalItems.value)
          }
          showSuccess(`已自动恢复 ${needRestore.length} 条实际支出数据`)
        } catch (e: any) {
          showError('自动恢复实际支出数据失败：' + (e?.message || '未知错误'))
        }
      }
      // 无论成功失败都标记，防止重复触发
      sessionStorage.setItem(RESTORE_KEY, 'true')
    }
  }
})

// ====== computed ======
const displayCategories = computed<DisplayCategory[]>(() => {
  const source = editing.value ? draftItems.value : originalItems.value
  const grouped: Record<string, DraftItem[]> = {}
  for (const item of source) { if (!grouped[item.category]) grouped[item.category] = []; grouped[item.category].push(item) }
  const result: DisplayCategory[] = []
  const seen = new Set<string>()
  for (const catName of CATEGORY_ORDER) {
    if (grouped[catName]) { seen.add(catName); result.push({ key: CATEGORY_KEY_MAP[catName] || catName, name: catName, items: grouped[catName].sort((a, b) => a.sortOrder - b.sortOrder) }) }
  }
  for (const catName of Object.keys(grouped)) { if (!seen.has(catName)) result.push({ key: catName.toLowerCase(), name: catName, items: grouped[catName].sort((a, b) => a.sortOrder - b.sortOrder) }) }
  return result
})

const computedSummary = computed(() => {
  const source = editing.value ? draftItems.value : originalItems.value
  let totalBudget = 0, totalActual = 0
  for (const item of source) { totalBudget += item.amount || 0; if (item.actual != null) totalActual += item.actual || 0 }
  return { totalBudget, totalActual, totalDiff: totalBudget - totalActual }
})

// ====== 编辑操作 ======
function enterEditMode() { draftItems.value = deepCloneItems(originalItems.value); editing.value = true }
function removeDraftRow(uid: string) { draftItems.value = draftItems.value.filter(i => i._uid !== uid) }
function addDraftRow(categoryKey: string) {
  const catName = Object.entries(CATEGORY_KEY_MAP).find(([, v]) => v === categoryKey)?.[0] || categoryKey
  const catItems = draftItems.value.filter(i => i.category === catName)
  const maxSort = catItems.length > 0 ? Math.max(...catItems.map(i => i.sortOrder)) : 0
  draftItems.value.push({ _uid: genUid(), category: catName, itemName: '', amount: 0, actual: null, remark: '', sortOrder: maxSort + 1 })
}

function hasDraftChanges(): boolean {
  if (draftItems.value.length !== originalItems.value.length) return true
  const ds = [...draftItems.value].sort((a, b) => a.sortOrder - b.sortOrder)
  const os = [...originalItems.value].sort((a, b) => a.sortOrder - b.sortOrder)
  for (let i = 0; i < ds.length; i++) {
    const d = ds[i], o = os[i]
    if (!o || d.category !== o.category || d.itemName !== o.itemName || d.amount !== o.amount || d.actual !== o.actual || d.remark !== o.remark) return true
  }
  return false
}

function handleCancel() {
  if (hasDraftChanges()) {
    showConfirm({ title: '确认放弃修改？', content: '所有未保存的改动将丢失', okText: '放弃', cancelText: '继续编辑', danger: true, onOk: () => { draftItems.value = deepCloneItems(originalItems.value); editing.value = false } })
  } else { editing.value = false }
}

async function handleSave() {
  for (let i = 0; i < draftItems.value.length; i++) {
    const item = draftItems.value[i]
    if (!item.category.trim()) { showError(`第 ${i + 1} 行「类别」不能为空`); return }
    if (!item.itemName.trim()) { showError(`第 ${i + 1} 行「项目名称」不能为空`); return }
    if (item.amount < 0) { showError(`第 ${i + 1} 行「金额」不能为负数`); return }
  }
  saving.value = true
  try {
    let body: { creates: any[]; updates: any[]; deletes: any[] }
    if (!migrated.value) {
      body = { creates: draftItems.value.map((item, idx) => ({ category: item.category, item_name: item.itemName, amount: item.amount, actual: item.actual, remark: item.remark, sort_order: idx + 1 })), updates: [], deletes: [] }
    } else {
      const creates: any[] = [], updates: any[] = [], deletes: any[] = []
      const origMap = new Map<number, DraftItem>()
      for (const item of originalItems.value) { if (item.id) origMap.set(item.id, item) }
      const draftIdSet = new Set<number>()
      for (let idx = 0; idx < draftItems.value.length; idx++) {
        const item = draftItems.value[idx]
        if (!item.id) { creates.push({ category: item.category, item_name: item.itemName, amount: item.amount, actual: item.actual, remark: item.remark, sort_order: idx + 1 }) }
        else {
          draftIdSet.add(item.id)
          const orig = origMap.get(item.id)
          if (orig && (orig.category !== item.category || orig.itemName !== item.itemName || orig.amount !== item.amount || orig.actual !== item.actual || orig.remark !== item.remark || orig.sortOrder !== idx + 1)) {
            updates.push({ id: item.id, category: item.category, item_name: item.itemName, amount: item.amount, actual: item.actual, remark: item.remark, sort_order: idx + 1 })
          }
        }
      }
      for (const item of originalItems.value) { if (item.id && !draftIdSet.has(item.id)) deletes.push({ id: item.id }) }
      body = { creates, updates, deletes }
    }
    const res = await $fetch<{ code: number; msg: string; data: { items: BudgetItemAPI[] } }>('/api/renovation/budget/items/batch', { method: 'POST', body })
    if (res.data?.items) { originalItems.value = apiItemsToDraft(res.data.items); migrated.value = true }
    draftItems.value = deepCloneItems(originalItems.value)
    editing.value = false
    showSuccess('保存成功')
  } catch (err: any) { showError(err?.data?.statusMessage || err?.data?.msg || err?.message || '保存失败') }
  finally { saving.value = false }
}

// ====== 工具函数 ======
function formatMoney(value: number): string { return '¥' + value.toLocaleString('zh-CN') }
function diffTextByItem(item: DraftItem): string { if (item.actual == null) return '—'; const d = item.amount - item.actual; if (d === 0) return '¥0'; return d > 0 ? formatMoney(d) : '-' + formatMoney(Math.abs(d)) }
function diffClassByItem(item: DraftItem): string { if (item.actual == null) return 'text-gray-300 dark:text-gray-600'; const d = item.amount - item.actual; if (d > 0) return 'text-emerald-600 dark:text-emerald-400'; if (d < 0) return 'text-red-600 dark:text-red-400'; return 'text-gray-500 dark:text-gray-400' }
function actualClassByItem(item: DraftItem): string { if (item.actual == null) return 'text-gray-300 dark:text-gray-600'; if (item.actual > item.amount) return 'text-red-600 dark:text-red-400'; return 'text-emerald-600 dark:text-emerald-400' }

// ====== 分类样式 ======
function categoryBarClass(key: string): string { return ({ partial: 'bg-blue-500', hard: 'bg-orange-500', furniture: 'bg-emerald-500', appliance: 'bg-violet-500', smart: 'bg-rose-500', other: 'bg-amber-500', extra: 'bg-cyan-500' } as Record<string, string>)[key] || 'bg-gray-400' }
function categoryBadgeClass(key: string): string { return ({ partial: 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300', hard: 'bg-orange-50 text-orange-700 dark:bg-orange-900/30 dark:text-orange-300', furniture: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300', appliance: 'bg-violet-50 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300', smart: 'bg-rose-50 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300', other: 'bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300', extra: 'bg-cyan-50 text-cyan-700 dark:bg-cyan-900/30 dark:text-cyan-300' } as Record<string, string>)[key] || '' }
function categoryHeaderBgClass(key: string): string { return ({ partial: 'bg-blue-50/60 dark:bg-blue-900/15', hard: 'bg-orange-50/60 dark:bg-orange-900/15', furniture: 'bg-emerald-50/60 dark:bg-emerald-900/15', appliance: 'bg-violet-50/60 dark:bg-violet-900/15', smart: 'bg-rose-50/60 dark:bg-rose-900/15', other: 'bg-amber-50/60 dark:bg-amber-900/15', extra: 'bg-cyan-50/60 dark:bg-cyan-900/15' } as Record<string, string>)[key] || 'bg-gray-50 dark:bg-gray-800/80' }
function categorySubtotalBgClass(key: string): string { return ({ partial: 'bg-blue-50/50 dark:bg-blue-900/10', hard: 'bg-orange-50/50 dark:bg-orange-900/10', furniture: 'bg-emerald-50/50 dark:bg-emerald-900/10', appliance: 'bg-violet-50/50 dark:bg-violet-900/10', smart: 'bg-rose-50/50 dark:bg-rose-900/10', other: 'bg-amber-50/50 dark:bg-amber-900/10', extra: 'bg-cyan-50/50 dark:bg-cyan-900/10' } as Record<string, string>)[key] || '' }
function categorySubtotalTextClass(key: string): string { return ({ partial: 'text-blue-700 dark:text-blue-300', hard: 'text-orange-700 dark:text-orange-300', furniture: 'text-emerald-700 dark:text-emerald-300', appliance: 'text-violet-700 dark:text-violet-300', smart: 'text-rose-700 dark:text-rose-300', other: 'text-amber-700 dark:text-amber-300', extra: 'text-cyan-700 dark:text-cyan-300' } as Record<string, string>)[key] || '' }

function categorySumAmount(cat: DisplayCategory): number { return cat.items.reduce((s, i) => s + (i.amount || 0), 0) }
function categorySumActual(cat: DisplayCategory): number { return cat.items.reduce((s, i) => s + (i.actual ?? 0), 0) }
function categoryDiffText(cat: DisplayCategory): string { const a = categorySumActual(cat); if (a === 0) return '—'; const d = categorySumAmount(cat) - a; if (d === 0) return '¥0'; return d > 0 ? formatMoney(d) : '-' + formatMoney(Math.abs(d)) }
function categoryActualClass(cat: DisplayCategory): string { const a = categorySumActual(cat); if (a === 0) return 'text-gray-300 dark:text-gray-600'; return a > categorySumAmount(cat) ? 'text-red-600 dark:text-red-400' : 'text-emerald-600 dark:text-emerald-400' }
function categoryDiffClass(cat: DisplayCategory): string { const a = categorySumActual(cat); if (a === 0) return 'text-gray-300 dark:text-gray-600'; const d = categorySumAmount(cat) - a; if (d > 0) return 'text-emerald-600 dark:text-emerald-400'; if (d < 0) return 'text-red-600 dark:text-red-400'; return 'text-gray-500 dark:text-gray-400' }

// ====== 预览弹窗 ======
const previewVisible = ref(false)
const previewUrl = ref<string | null>(null)
function closePreview() { previewVisible.value = false; previewUrl.value = null }

// ====== Canvas 绘制 ======
const AVATAR_URL = 'https://assets.fatwill.cloud/uploads/1775045182650-e9d2upw9.jpeg'
const QR_URL = 'https://fatwill.cloud/renovation/budget'
const IMG_WIDTH = 900
const DPR = 2
const categoryColorsMap: Record<string, { main: string; light: string; text: string }> = { partial: { main: '#3b82f6', light: '#eff6ff', text: '#1e40af' }, hard: { main: '#f97316', light: '#fff7ed', text: '#9a3412' }, furniture: { main: '#10b981', light: '#ecfdf5', text: '#065f46' }, appliance: { main: '#8b5cf6', light: '#f5f3ff', text: '#5b21b6' }, smart: { main: '#f43f5e', light: '#fff1f2', text: '#9f1239' }, other: { main: '#f59e0b', light: '#fffbeb', text: '#92400e' }, extra: { main: '#06b6d4', light: '#ecfeff', text: '#155e75' } }

interface CanvasCat { key: string; name: string; items: { name: string; budget: number; actual: number | null; remark: string | null }[] }
function toCanvasCats(cats: DisplayCategory[]): CanvasCat[] { return cats.map(c => ({ key: c.key, name: c.name, items: c.items.map(i => ({ name: i.itemName, budget: i.amount, actual: i.actual, remark: i.remark || null })) })) }
function ccBudgetSum(cat: CanvasCat): number { return cat.items.reduce((s, i) => s + i.budget, 0) }
function ccActualSum(cat: CanvasCat): number { return cat.items.reduce((s, i) => s + (i.actual ?? 0), 0) }

function loadImage(url: string): Promise<HTMLImageElement | null> { return new Promise(r => { const img = new Image(); img.crossOrigin = 'anonymous'; img.onload = () => r(img); img.onerror = () => r(null); img.src = url }) }
function drawCircleImage(ctx: CanvasRenderingContext2D, img: HTMLImageElement | null, x: number, y: number, size: number) { ctx.save(); ctx.beginPath(); ctx.arc(x + size / 2, y + size / 2, size / 2, 0, Math.PI * 2); ctx.closePath(); ctx.clip(); if (img) ctx.drawImage(img, x, y, size, size); else { ctx.fillStyle = '#d1d5db'; ctx.fillRect(x, y, size, size) }; ctx.restore() }
function drawRoundRect(ctx: CanvasRenderingContext2D, x: number, y: number, w: number, h: number, r: number) { ctx.beginPath(); ctx.moveTo(x + r, y); ctx.lineTo(x + w - r, y); ctx.arcTo(x + w, y, x + w, y + r, r); ctx.lineTo(x + w, y + h - r); ctx.arcTo(x + w, y + h, x + w - r, y + h, r); ctx.lineTo(x + r, y + h); ctx.arcTo(x, y + h, x, y + h - r, r); ctx.lineTo(x, y + r); ctx.arcTo(x, y, x + r, y, r); ctx.closePath() }
function moneyNum(v: number): string { return v.toLocaleString('zh-CN') }

async function generateBudgetImage(categories: CanvasCat[], title: string, subtitle: string): Promise<void> {
  const [avatarImg, qrCanvas] = await Promise.all([loadImage(AVATAR_URL), QRCode.toCanvas(document.createElement('canvas'), QR_URL, { width: 80 * DPR, margin: 1, color: { dark: '#333333', light: '#ffffff' } }).then(c => c)])
  const headerH = 80, catTitleH = 40, tableHeaderH = 32, rowH = 30, catGap = 20, footerH = 100, padding = 32
  let contentH = headerH + padding; for (const cat of categories) contentH += catTitleH + tableHeaderH + cat.items.length * rowH + rowH + catGap; contentH += footerH + padding
  const canvas = document.createElement('canvas'); canvas.width = IMG_WIDTH * DPR; canvas.height = contentH * DPR
  const ctx = canvas.getContext('2d')!; ctx.scale(DPR, DPR)
  ctx.fillStyle = '#ffffff'; ctx.fillRect(0, 0, IMG_WIDTH, contentH)
  ctx.fillStyle = '#f8f9fa'; ctx.fillRect(0, 0, IMG_WIDTH, headerH)
  ctx.strokeStyle = '#e5e7eb'; ctx.lineWidth = 1; ctx.beginPath(); ctx.moveTo(0, headerH); ctx.lineTo(IMG_WIDTH, headerH); ctx.stroke()
  const avatarSize = 48, avatarX = padding, avatarY = (headerH - avatarSize) / 2
  drawCircleImage(ctx, avatarImg, avatarX, avatarY, avatarSize)
  ctx.fillStyle = '#111827'; ctx.font = 'bold 18px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.textBaseline = 'middle'; ctx.fillText(title, avatarX + avatarSize + 14, headerH / 2 - 10)
  ctx.fillStyle = '#9ca3af'; ctx.font = '12px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.fillText(subtitle, avatarX + avatarSize + 14, headerH / 2 + 12)
  ctx.textAlign = 'right'; ctx.fillText(new Date().toISOString().slice(0, 10), IMG_WIDTH - padding, headerH / 2); ctx.textAlign = 'left'
  let curY = headerH + padding; const colWidths = [240, 130, 130, 130, 206]; const tableX = padding; const tableW = IMG_WIDTH - padding * 2
  for (const cat of categories) {
    const colors = categoryColorsMap[cat.key] || categoryColorsMap.other; const catBudget = ccBudgetSum(cat)
    drawRoundRect(ctx, tableX, curY, tableW, catTitleH, 8); ctx.fillStyle = colors.light; ctx.fill()
    ctx.fillStyle = colors.main; drawRoundRect(ctx, tableX, curY, 4, catTitleH, 2); ctx.fill()
    ctx.fillStyle = colors.text; ctx.font = 'bold 14px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.textBaseline = 'middle'; ctx.fillText(cat.name, tableX + 16, curY + catTitleH / 2)
    ctx.fillStyle = '#6b7280'; ctx.font = '12px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.textAlign = 'right'; ctx.fillText(`${cat.items.length} 项 · 预算 ¥${moneyNum(catBudget)}`, tableX + tableW - 12, curY + catTitleH / 2); ctx.textAlign = 'left'; curY += catTitleH
    ctx.fillStyle = '#f9fafb'; ctx.fillRect(tableX, curY, tableW, tableHeaderH); ctx.strokeStyle = '#e5e7eb'; ctx.lineWidth = 0.5; ctx.beginPath(); ctx.moveTo(tableX, curY + tableHeaderH); ctx.lineTo(tableX + tableW, curY + tableHeaderH); ctx.stroke()
    const headers = ['项目名称', '预算', '实际支出', '差额', '备注']; ctx.fillStyle = '#6b7280'; ctx.font = 'bold 11px "PingFang SC", "Microsoft YaHei", sans-serif'
    let colX = tableX + 12; for (let i = 0; i < headers.length; i++) { if (i >= 1 && i <= 3) { ctx.textAlign = 'right'; ctx.fillText(headers[i], colX + colWidths[i] - 12, curY + tableHeaderH / 2); ctx.textAlign = 'left' } else ctx.fillText(headers[i], colX, curY + tableHeaderH / 2); colX += colWidths[i] }; curY += tableHeaderH
    for (let ri = 0; ri < cat.items.length; ri++) {
      const item = cat.items[ri]; if (ri % 2 === 1) { ctx.fillStyle = '#f9fafb'; ctx.fillRect(tableX, curY, tableW, rowH) }
      ctx.strokeStyle = '#f3f4f6'; ctx.lineWidth = 0.5; ctx.beginPath(); ctx.moveTo(tableX, curY + rowH); ctx.lineTo(tableX + tableW, curY + rowH); ctx.stroke()
      ctx.font = '12px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.textBaseline = 'middle'; const midY = curY + rowH / 2; colX = tableX + 12
      ctx.fillStyle = '#111827'; ctx.fillText(item.name, colX, midY); colX += colWidths[0]
      ctx.fillStyle = '#374151'; ctx.textAlign = 'right'; ctx.fillText(`¥${moneyNum(item.budget)}`, colX + colWidths[1] - 12, midY); colX += colWidths[1]
      if (item.actual !== null) { ctx.fillStyle = item.actual > item.budget ? '#dc2626' : '#059669'; ctx.fillText(`¥${moneyNum(item.actual)}`, colX + colWidths[2] - 12, midY) } else { ctx.fillStyle = '#d1d5db'; ctx.fillText('—', colX + colWidths[2] - 12, midY) }; colX += colWidths[2]
      if (item.actual !== null) { const diff = item.budget - item.actual; if (diff > 0) { ctx.fillStyle = '#059669'; ctx.fillText(`¥${moneyNum(diff)}`, colX + colWidths[3] - 12, midY) } else if (diff < 0) { ctx.fillStyle = '#dc2626'; ctx.fillText(`-¥${moneyNum(Math.abs(diff))}`, colX + colWidths[3] - 12, midY) } else { ctx.fillStyle = '#6b7280'; ctx.fillText('¥0', colX + colWidths[3] - 12, midY) } } else { ctx.fillStyle = '#d1d5db'; ctx.fillText('—', colX + colWidths[3] - 12, midY) }; colX += colWidths[3]
      ctx.textAlign = 'left'; ctx.fillStyle = '#9ca3af'; ctx.fillText(item.remark || '—', colX, midY); curY += rowH
    }
    ctx.fillStyle = colors.light; ctx.fillRect(tableX, curY, tableW, rowH); ctx.strokeStyle = '#e5e7eb'; ctx.lineWidth = 0.5; ctx.beginPath(); ctx.moveTo(tableX, curY); ctx.lineTo(tableX + tableW, curY); ctx.stroke()
    ctx.font = 'bold 12px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.textBaseline = 'middle'; const subMidY = curY + rowH / 2; colX = tableX + 12
    ctx.fillStyle = colors.text; ctx.fillText('小计', colX, subMidY); colX += colWidths[0]
    ctx.fillStyle = '#374151'; ctx.textAlign = 'right'; ctx.fillText(`¥${moneyNum(catBudget)}`, colX + colWidths[1] - 12, subMidY); colX += colWidths[1]
    const catActual = ccActualSum(cat); if (catActual > 0) { ctx.fillStyle = catActual > catBudget ? '#dc2626' : '#059669'; ctx.fillText(`¥${moneyNum(catActual)}`, colX + colWidths[2] - 12, subMidY) } else { ctx.fillStyle = '#d1d5db'; ctx.fillText('—', colX + colWidths[2] - 12, subMidY) }; colX += colWidths[2]
    if (catActual > 0) { const cd = catBudget - catActual; if (cd > 0) { ctx.fillStyle = '#059669'; ctx.fillText(`¥${moneyNum(cd)}`, colX + colWidths[3] - 12, subMidY) } else if (cd < 0) { ctx.fillStyle = '#dc2626'; ctx.fillText(`-¥${moneyNum(Math.abs(cd))}`, colX + colWidths[3] - 12, subMidY) } else { ctx.fillStyle = '#6b7280'; ctx.fillText('¥0', colX + colWidths[3] - 12, subMidY) } } else { ctx.fillStyle = '#d1d5db'; ctx.fillText('—', colX + colWidths[3] - 12, subMidY) }
    ctx.textAlign = 'left'; curY += rowH + catGap
  }
  const footerY = curY; ctx.fillStyle = '#f8f9fa'; ctx.fillRect(0, footerY, IMG_WIDTH, footerH); ctx.strokeStyle = '#e5e7eb'; ctx.lineWidth = 1; ctx.beginPath(); ctx.moveTo(0, footerY); ctx.lineTo(IMG_WIDTH, footerY); ctx.stroke()
  let tb = 0, ta = 0; for (const cat of categories) { tb += ccBudgetSum(cat); ta += ccActualSum(cat) }
  ctx.font = 'bold 14px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.textBaseline = 'middle'; ctx.fillStyle = '#111827'; ctx.fillText(`总预算  ¥${moneyNum(tb)}`, padding, footerY + footerH / 2 - 12)
  ctx.fillStyle = ta > tb ? '#dc2626' : '#059669'; ctx.font = 'bold 13px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.fillText(`已支出  ¥${moneyNum(ta)}`, padding, footerY + footerH / 2 + 12)
  const qrSize = 70, qrX = IMG_WIDTH - padding - qrSize, qrY = footerY + (footerH - qrSize - 14) / 2
  ctx.drawImage(qrCanvas, qrX, qrY, qrSize, qrSize)
  ctx.fillStyle = '#9ca3af'; ctx.font = '9px "PingFang SC", "Microsoft YaHei", sans-serif'; ctx.textAlign = 'center'; ctx.fillText('扫码查看详情', qrX + qrSize / 2, qrY + qrSize + 11); ctx.textAlign = 'left'
  previewUrl.value = canvas.toDataURL('image/png'); previewVisible.value = true
}

async function downloadAll() { await generateBudgetImage(toCanvasCats(displayCategories.value), 'fatwill 的预算', '装修成本预算清单') }
async function downloadCategory(cat: DisplayCategory) { await generateBudgetImage(toCanvasCats([cat]), `fatwill 的预算 · ${cat.name}`, `${cat.name}成本明细`) }

useSeoMeta({ title: '成本预算 - 装修', description: '装修成本预算明细表，包含局改、硬装、家私、家电、全屋智能、其他六大分类' })
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
