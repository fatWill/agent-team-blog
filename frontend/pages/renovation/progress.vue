<template>
  <div class="flex min-h-screen flex-col bg-gray-50 transition-colors dark:bg-gray-900">
    <!-- 顶部导航 -->
    <header class="sticky top-0 z-30 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex h-14 max-w-4xl items-center px-4">
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
        <span class="text-sm font-medium text-gray-900 dark:text-gray-100">踩坑日记</span>
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

    <main class="mx-auto w-full max-w-4xl flex-1 px-4 py-6">
      <!-- 加载态 -->
      <div v-if="loading" class="space-y-4">
        <div v-for="i in 3" :key="i" class="animate-pulse rounded-xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
          <div class="mb-3 h-5 w-16 rounded bg-gray-200 dark:bg-gray-700" />
          <div class="mb-2 h-4 w-3/4 rounded bg-gray-200 dark:bg-gray-700" />
          <div class="h-4 w-1/2 rounded bg-gray-200 dark:bg-gray-700" />
        </div>
      </div>

      <!-- 空态 -->
      <div v-else-if="displayItems.length === 0 && !editing" class="py-20 text-center">
        <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-blue-50 dark:bg-blue-900/20">
          <span class="text-3xl">📋</span>
        </div>
        <h2 class="text-lg font-bold text-gray-900 dark:text-gray-100">暂无踩坑记录</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">装修的坑还没开始记录，希望一切顺利 ✨</p>
        <button
          v-if="isLoggedIn"
          class="mt-6 inline-flex items-center gap-1.5 rounded-lg bg-indigo-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-indigo-600"
          @click="enterEditModeWithNewRow"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
          新增第一条踩坑记录
        </button>
      </div>

      <!-- 卡片列表 -->
      <div v-else class="space-y-4">
        <div
          v-for="(item, idx) in displayItems"
          :key="item._uid"
          class="relative rounded-xl border p-5 shadow-sm transition-all"
          :class="editing
            ? 'border-blue-200 bg-blue-50/50 dark:border-blue-800/50 dark:bg-blue-900/10'
            : 'border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800'"
        >
          <!-- 顶部：分类徽章 + 删除按钮 -->
          <div class="mb-3 flex items-start justify-between gap-3">
            <!-- 展示态：分类徽章 -->
            <span
              v-if="!editing"
              class="inline-flex shrink-0 items-center rounded-full px-2.5 py-1 text-xs font-medium"
              :class="categoryBadgeClass(item.category)"
            >
              {{ item.category }}
            </span>
            <!-- 编辑态：分类下拉 -->
            <select
              v-else
              v-model="item.category"
              class="shrink-0 rounded-lg border border-gray-200 bg-white px-2.5 py-1 text-xs font-medium text-gray-700 outline-none focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-200"
            >
              <option v-for="cat in CATEGORIES" :key="cat" :value="cat">{{ cat }}</option>
            </select>
            <!-- 删除按钮（编辑态） -->
            <button
              v-if="editing"
              class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-gray-300 transition-colors hover:bg-red-50 hover:text-red-500 dark:text-gray-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
              title="删除此条"
              @click="removeDraftRow(item._uid)"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
              </svg>
            </button>
            <!-- 序号标记（展示态） -->
            <span v-else class="shrink-0 text-xs text-gray-300 dark:text-gray-600">#{{ idx + 1 }}</span>
          </div>

          <!-- 踩坑点 -->
          <div class="mb-3">
            <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-gray-500">踩坑点</label>
            <textarea
              v-if="editing"
              v-model="item.pitfall"
              rows="2"
              class="w-full resize-none rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              :style="{ fieldSizing: 'content' }"
              placeholder="遇到了什么坑？"
            />
            <p v-else class="whitespace-pre-wrap break-words text-sm text-gray-800 dark:text-gray-200">{{ item.pitfall || '—' }}</p>
          </div>

          <!-- 注意事项/解决方案 -->
          <div class="mb-3">
            <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-gray-500">注意事项 / 解决方案</label>
            <textarea
              v-if="editing"
              v-model="item.solution"
              rows="2"
              class="w-full resize-none rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              :style="{ fieldSizing: 'content' }"
              placeholder="怎么解决的？"
            />
            <p v-else class="whitespace-pre-wrap break-words text-sm text-gray-800 dark:text-gray-200">{{ item.solution || '—' }}</p>
          </div>

          <!-- 备注 -->
          <div v-if="editing || item.remark">
            <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-gray-500">备注</label>
            <textarea
              v-if="editing"
              v-model="item.remark"
              rows="1"
              class="w-full resize-none rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              :style="{ fieldSizing: 'content' }"
              placeholder="补充说明（可选）"
            />
            <p v-else class="whitespace-pre-wrap break-words text-xs text-gray-500 dark:text-gray-400">{{ item.remark }}</p>
          </div>
        </div>

        <!-- 新增一行按钮 -->
        <button
          v-if="editing"
          class="flex w-full items-center justify-center gap-2 rounded-xl border-2 border-dashed border-gray-300 py-4 text-sm font-medium text-gray-400 transition-colors hover:border-indigo-400 hover:text-indigo-500 dark:border-gray-600 dark:text-gray-500 dark:hover:border-indigo-500 dark:hover:text-indigo-400"
          @click="addDraftRow"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
          新增一行
        </button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { showSuccess, showError, showConfirm } from '~/shared/utils/ui'

// ====== 类型 ======
interface PitfallItemAPI {
  id: number
  category: string
  pitfall: string
  solution: string
  remark: string
  sortOrder: number
  createdAt: string
  updatedAt: string
}

interface DraftItem {
  _uid: string
  id?: number
  category: string
  pitfall: string
  solution: string
  remark: string
  sortOrder: number
}

// ====== 分类（从接口获取，SSR 兼容） ======
const { data: categoriesData } = await useAsyncData('pitfall-categories', () =>
  $fetch<{ categories: string[] }>('/api/renovation/pitfall/categories'),
)
const CATEGORIES = computed(() => categoriesData.value?.categories ?? [])

// ====== 状态 ======
const { isLoggedIn } = useAuth()
const editing = ref(false)
const saving = ref(false)
const loading = ref(true)
const originalItems = ref<DraftItem[]>([])
const draftItems = ref<DraftItem[]>([])

let uidCounter = 0
function genUid(): string { return `uid_${++uidCounter}_${Date.now()}` }

function apiItemsToDraft(items: PitfallItemAPI[]): DraftItem[] {
  return items.map(item => ({
    _uid: genUid(),
    id: item.id,
    category: item.category,
    pitfall: item.pitfall || '',
    solution: item.solution || '',
    remark: item.remark || '',
    sortOrder: item.sortOrder,
  }))
}

function deepCloneItems(items: DraftItem[]): DraftItem[] {
  return items.map(item => ({ ...item, _uid: genUid() }))
}

// SSR-safe: 初始为空数组
originalItems.value = []
draftItems.value = []

onMounted(async () => {
  try {
    const res = await $fetch<{ items: PitfallItemAPI[] }>('/api/renovation/pitfall/items')
    if (res.items) {
      originalItems.value = apiItemsToDraft(res.items)
    }
  } catch {
    originalItems.value = []
  }
  draftItems.value = deepCloneItems(originalItems.value)
  loading.value = false
})

// ====== computed ======
const displayItems = computed<DraftItem[]>(() => {
  const source = editing.value ? draftItems.value : originalItems.value
  return [...source].sort((a, b) => a.sortOrder - b.sortOrder)
})

// ====== 编辑操作 ======
function enterEditMode() {
  draftItems.value = deepCloneItems(originalItems.value)
  editing.value = true
}

function enterEditModeWithNewRow() {
  draftItems.value = deepCloneItems(originalItems.value)
  editing.value = true
  addDraftRow()
}

function removeDraftRow(uid: string) {
  draftItems.value = draftItems.value.filter(i => i._uid !== uid)
}

function addDraftRow() {
  const maxSort = draftItems.value.length > 0
    ? Math.max(...draftItems.value.map(i => i.sortOrder))
    : 0
  draftItems.value.push({
    _uid: genUid(),
    category: CATEGORIES.value[0] || '',
    pitfall: '',
    solution: '',
    remark: '',
    sortOrder: maxSort + 1,
  })
}

function hasDraftChanges(): boolean {
  if (draftItems.value.length !== originalItems.value.length) return true
  const ds = [...draftItems.value].sort((a, b) => a.sortOrder - b.sortOrder)
  const os = [...originalItems.value].sort((a, b) => a.sortOrder - b.sortOrder)
  for (let i = 0; i < ds.length; i++) {
    const d = ds[i], o = os[i]
    if (!o || d.category !== o.category || d.pitfall !== o.pitfall || d.solution !== o.solution || d.remark !== o.remark) return true
  }
  return false
}

function handleCancel() {
  if (hasDraftChanges()) {
    showConfirm({
      title: '确认放弃修改？',
      content: '所有未保存的改动将丢失',
      okText: '放弃',
      cancelText: '继续编辑',
      danger: true,
      onOk: () => {
        draftItems.value = deepCloneItems(originalItems.value)
        editing.value = false
      },
    })
  } else {
    editing.value = false
  }
}

async function handleSave() {
  // 前端校验
  for (let i = 0; i < draftItems.value.length; i++) {
    const item = draftItems.value[i]
    if (!item.category || !CATEGORIES.value.includes(item.category)) {
      showError(`第 ${i + 1} 条「类型」必须选择有效分类`)
      return
    }
    if (!item.pitfall.trim() && !item.solution.trim()) {
      showError(`第 ${i + 1} 条「踩坑点」和「解决方案」不能同时为空`)
      return
    }
    if (item.pitfall.length > 5000) { showError(`第 ${i + 1} 条「踩坑点」不能超过 5000 字`); return }
    if (item.solution.length > 5000) { showError(`第 ${i + 1} 条「解决方案」不能超过 5000 字`); return }
    if (item.remark.length > 5000) { showError(`第 ${i + 1} 条「备注」不能超过 5000 字`); return }
  }

  saving.value = true
  try {
    const creates: any[] = []
    const updates: any[] = []
    const deletes: any[] = []

    const origMap = new Map<number, DraftItem>()
    for (const item of originalItems.value) {
      if (item.id) origMap.set(item.id, item)
    }

    const draftIdSet = new Set<number>()
    for (let idx = 0; idx < draftItems.value.length; idx++) {
      const item = draftItems.value[idx]
      if (!item.id) {
        creates.push({
          category: item.category,
          pitfall: item.pitfall,
          solution: item.solution,
          remark: item.remark,
          sort_order: idx + 1,
        })
      } else {
        draftIdSet.add(item.id)
        const orig = origMap.get(item.id)
        if (orig && (
          orig.category !== item.category ||
          orig.pitfall !== item.pitfall ||
          orig.solution !== item.solution ||
          orig.remark !== item.remark ||
          orig.sortOrder !== idx + 1
        )) {
          updates.push({
            id: item.id,
            category: item.category,
            pitfall: item.pitfall,
            solution: item.solution,
            remark: item.remark,
            sort_order: idx + 1,
          })
        }
      }
    }

    for (const item of originalItems.value) {
      if (item.id && !draftIdSet.has(item.id)) {
        deletes.push({ id: item.id })
      }
    }

    const body = { creates, updates, deletes }
    const res = await $fetch<{ code: number; msg: string; data: { items: PitfallItemAPI[] } }>(
      '/api/renovation/pitfall/items/batch',
      { method: 'POST', body },
    )
    if (res.data?.items) {
      originalItems.value = apiItemsToDraft(res.data.items)
    }
    draftItems.value = deepCloneItems(originalItems.value)
    editing.value = false
    showSuccess('保存成功')
  } catch (err: any) {
    showError(err?.data?.statusMessage || err?.data?.msg || err?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

// ====== 分类样式 ======
function categoryBadgeClass(category: string): string {
  const map: Record<string, string> = {
    '局改': 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300',
    '全屋定制': 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
    '家电': 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300',
    '家私': 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300',
    '油漆': 'bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300',
    '美缝': 'bg-teal-100 text-teal-700 dark:bg-teal-900/30 dark:text-teal-300',
  }
  return map[category] || 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-300'
}

// ====== SEO ======
useSeoMeta({
  title: '踩坑日记 - 装修 - fatwill 的小屋',
  description: '记录装修过程中的坑与教训',
})
</script>
