<template>
  <div class="flex min-h-screen flex-col bg-gray-50 transition-colors dark:bg-gray-900">
    <!-- 顶部导航 -->
    <header class="sticky top-0 z-30 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex h-14 max-w-4xl items-center px-4">
        <NuxtLink
          to="/life"
          class="flex items-center gap-2 text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
          </svg>
          返回生活
        </NuxtLink>
        <span class="mx-3 text-gray-300 dark:text-gray-600">/</span>
        <span class="text-sm font-medium text-gray-900 dark:text-gray-100">阳阳的成长日记</span>
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
      <!-- ====== 编辑模式 ====== -->
      <div v-if="editing" class="space-y-4">
        <div
          v-for="(item, idx) in draftItems"
          :key="item._uid"
          class="relative rounded-xl border border-blue-200 bg-blue-50/50 p-5 shadow-sm dark:border-blue-800/50 dark:bg-blue-900/10"
        >
          <!-- 顶部：序号 + 删除 -->
          <div class="mb-3 flex items-center justify-between">
            <span class="text-xs font-medium text-gray-400">#{{ idx + 1 }}</span>
            <button
              class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-gray-300 transition-colors hover:bg-red-50 hover:text-red-500 dark:text-gray-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
              title="删除此条"
              @click="removeDraftRow(item._uid)"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
              </svg>
            </button>
          </div>

          <!-- 内容 -->
          <div class="mb-3">
            <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-gray-500">内容</label>
            <textarea
              v-model="item.content"
              rows="3"
              class="w-full resize-none rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              :style="{ fieldSizing: 'content' }"
              placeholder="记录阳阳的成长瞬间..."
            />
          </div>

          <!-- 图片上传 -->
          <div class="mb-3">
            <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-gray-500">图片</label>
            <div class="flex flex-wrap gap-2">
              <div v-for="(img, imgIdx) in item.images" :key="imgIdx" class="group relative h-20 w-20 overflow-hidden rounded-lg border border-gray-200 dark:border-gray-600">
                <img :src="img" class="h-full w-full object-cover" />
                <button
                  class="absolute right-0.5 top-0.5 flex h-5 w-5 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity group-hover:opacity-100"
                  @click="item.images.splice(imgIdx, 1)"
                >
                  <svg class="h-3 w-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
              </div>
              <label class="flex h-20 w-20 cursor-pointer items-center justify-center rounded-lg border-2 border-dashed border-gray-300 text-gray-400 transition-colors hover:border-indigo-400 hover:text-indigo-500 dark:border-gray-600 dark:hover:border-indigo-500">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
                <input type="file" accept="image/*" multiple class="hidden" @change="(e) => handleImageUpload(e, item)" />
              </label>
            </div>
          </div>

          <!-- 视频上传 -->
          <div class="mb-3">
            <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-gray-500">视频</label>
            <div class="space-y-2">
              <div v-for="(vid, vidIdx) in item.videos" :key="vidIdx" class="group relative overflow-hidden rounded-lg border border-gray-200 dark:border-gray-600">
                <video :src="vid" class="h-32 w-full object-cover" />
                <button
                  class="absolute right-1 top-1 flex h-6 w-6 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity group-hover:opacity-100"
                  @click="item.videos.splice(vidIdx, 1)"
                >
                  <svg class="h-3 w-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
              </div>
              <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg border-2 border-dashed border-gray-300 py-3 text-sm text-gray-400 transition-colors hover:border-indigo-400 hover:text-indigo-500 dark:border-gray-600">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 002.25-2.25v-9a2.25 2.25 0 00-2.25-2.25h-9A2.25 2.25 0 002.25 7.5v9a2.25 2.25 0 002.25 2.25z" /></svg>
                <span>添加视频</span>
                <input type="file" accept="video/*" multiple class="hidden" @change="(e) => handleVideoUpload(e, item)" />
              </label>
            </div>
          </div>

          <!-- 发生时间 -->
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-400 dark:text-gray-500">发生时间</label>
            <input
              v-model="item.happenedAt"
              type="datetime-local"
              class="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
            />
          </div>
        </div>

        <!-- 新增一行按钮 -->
        <button
          class="flex w-full items-center justify-center gap-2 rounded-xl border-2 border-dashed border-gray-300 py-4 text-sm font-medium text-gray-400 transition-colors hover:border-indigo-400 hover:text-indigo-500 dark:border-gray-600 dark:text-gray-500 dark:hover:border-indigo-500 dark:hover:text-indigo-400"
          @click="addDraftRow"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
          新增一条
        </button>
      </div>

      <!-- ====== 前台展示模式 ====== -->
      <template v-else>
        <!-- 筛选区：时间 Tab + 搜索 -->
        <div class="mb-5 space-y-3">
          <!-- 时间 Tab -->
          <div v-if="months.length > 0" class="flex flex-wrap gap-2">
            <button
              class="rounded-full px-3 py-1.5 text-xs font-medium transition-all duration-200"
              :class="activeMonth === '全部'
                ? 'bg-indigo-500 text-white shadow-sm'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700'"
              @click="activeMonth = '全部'"
            >
              全部
            </button>
            <button
              v-for="m in months"
              :key="m"
              class="rounded-full px-3 py-1.5 text-xs font-medium transition-all duration-200"
              :class="activeMonth === m
                ? 'bg-indigo-500 text-white shadow-sm'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700'"
              @click="activeMonth = m"
            >
              {{ formatMonthLabel(m) }}
            </button>
          </div>
          <!-- 搜索框 -->
          <div class="relative">
            <svg class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索内容..."
              class="w-full rounded-lg border border-gray-200 bg-white py-2 pl-10 pr-4 text-sm text-gray-900 outline-none transition-colors focus:border-indigo-400 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            />
          </div>
        </div>

        <!-- 加载态 -->
        <div v-if="loading" class="space-y-4">
          <div v-for="i in 3" :key="i" class="animate-pulse rounded-xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
            <div class="mb-3 h-4 w-24 rounded bg-gray-200 dark:bg-gray-700" />
            <div class="mb-2 h-4 w-3/4 rounded bg-gray-200 dark:bg-gray-700" />
            <div class="h-32 w-full rounded bg-gray-200 dark:bg-gray-700" />
          </div>
        </div>

        <!-- 空态 -->
        <div v-else-if="filteredGroups.length === 0" class="py-20 text-center">
          <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-pink-50 dark:bg-pink-900/20">
            <span class="text-3xl">👶</span>
          </div>
          <h2 class="text-lg font-bold text-gray-900 dark:text-gray-100">暂无成长记录</h2>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">阳阳的精彩瞬间即将开始记录 🌟</p>
          <button
            v-if="isLoggedIn"
            class="mt-6 inline-flex items-center gap-1.5 rounded-lg bg-indigo-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-indigo-600"
            @click="enterEditModeWithNewRow"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
            记录第一条
          </button>
        </div>

        <!-- 按月分组的信息流 -->
        <div v-else class="space-y-8">
          <div v-for="group in filteredGroups" :key="group.month">
            <!-- 月份分割线 -->
            <div class="mb-4 flex items-center gap-3">
              <span class="text-sm font-bold text-gray-900 dark:text-gray-100">{{ formatMonthLabel(group.month) }}</span>
              <div class="h-px flex-1 bg-gray-200 dark:bg-gray-700" />
            </div>

            <!-- 动态卡片 -->
            <div class="space-y-4">
              <div
                v-for="item in group.items"
                :key="item.id"
                class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all dark:border-gray-700 dark:bg-gray-800"
              >
                <!-- 日期时间 -->
                <div class="mb-3 text-xs text-gray-400 dark:text-gray-500">
                  {{ formatDateTime(item.happenedAt) }}
                </div>

                <!-- 文本内容 -->
                <p v-if="item.content" class="mb-3 whitespace-pre-wrap break-words text-sm leading-relaxed text-gray-800 dark:text-gray-200">
                  {{ item.content }}
                </p>

                <!-- 图片区 -->
                <div v-if="item.images.length > 0" class="mb-3" :class="imageGridClass(item.images.length)">
                  <div
                    v-for="(img, imgIdx) in item.images"
                    :key="imgIdx"
                    class="cursor-pointer overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-700"
                    :class="imageItemClass(item.images.length, imgIdx)"
                    @click="openLightbox(item.images, imgIdx)"
                  >
                    <img :src="img" :alt="`图片 ${imgIdx + 1}`" class="h-full w-full object-cover transition-transform duration-200 hover:scale-105" loading="lazy" />
                  </div>
                </div>

                <!-- 视频区 -->
                <div v-if="item.videos.length > 0" class="space-y-3">
                  <div v-for="(vid, vidIdx) in item.videos" :key="vidIdx" class="overflow-hidden rounded-lg">
                    <video :src="vid" controls preload="metadata" class="w-full rounded-lg" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </main>

    <!-- 图片灯箱 -->
    <ClientOnly>
      <Teleport to="body">
        <Transition name="fade">
          <div
            v-if="lightboxVisible"
            class="fixed inset-0 z-[100] flex items-center justify-center bg-black/90"
            @click.self="lightboxVisible = false"
          >
            <button
              class="absolute right-4 top-4 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition-colors hover:bg-white/20"
              @click="lightboxVisible = false"
            >
              <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
            </button>
            <!-- 左箭头 -->
            <button
              v-if="lightboxImages.length > 1"
              class="absolute left-4 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition-colors hover:bg-white/20"
              @click="lightboxIndex = (lightboxIndex - 1 + lightboxImages.length) % lightboxImages.length"
            >
              <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" /></svg>
            </button>
            <!-- 图片 -->
            <img :src="lightboxImages[lightboxIndex]" class="max-h-[90vh] max-w-[90vw] rounded-lg object-contain" />
            <!-- 右箭头 -->
            <button
              v-if="lightboxImages.length > 1"
              class="absolute right-4 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition-colors hover:bg-white/20"
              @click="lightboxIndex = (lightboxIndex + 1) % lightboxImages.length"
            >
              <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" /></svg>
            </button>
            <!-- 计数 -->
            <div v-if="lightboxImages.length > 1" class="absolute bottom-6 left-1/2 -translate-x-1/2 rounded-full bg-black/50 px-3 py-1 text-sm text-white">
              {{ lightboxIndex + 1 }} / {{ lightboxImages.length }}
            </div>
          </div>
        </Transition>
      </Teleport>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import { showSuccess, showError, showConfirm } from '~/shared/utils/ui'

// ====== 类型 ======
interface DiaryItemAPI {
  id: number
  content: string
  images: string  // JSON string
  videos: string  // JSON string
  happenedAt: string
  createdAt: string
  updatedAt: string
}

interface DiaryItem {
  id: number
  content: string
  images: string[]
  videos: string[]
  happenedAt: string
}

interface DraftItem {
  _uid: string
  id?: number
  content: string
  images: string[]
  videos: string[]
  happenedAt: string
}

interface MonthGroup {
  month: string  // "2026-08"
  items: DiaryItem[]
}

// ====== SSR 数据 ======
const { data: itemsData } = await useAsyncData('growth-diary-items', () =>
  $fetch<{ items: DiaryItemAPI[] }>('/api/growth-diary/items'),
)
const { data: monthsData } = await useAsyncData('growth-diary-months', () =>
  $fetch<{ months: string[] }>('/api/growth-diary/months'),
)

const months = computed(() => monthsData.value?.months ?? [])

// ====== 状态 ======
const { isLoggedIn } = useAuth()
const editing = ref(false)
const saving = ref(false)
const loading = ref(false)
const activeMonth = ref('全部')
const searchQuery = ref('')
const draftItems = ref<DraftItem[]>([])

// 灯箱
const lightboxVisible = ref(false)
const lightboxImages = ref<string[]>([])
const lightboxIndex = ref(0)

let uidCounter = 0
function genUid(): string { return `uid_${++uidCounter}_${Date.now()}` }

// ====== 解析 API 数据 ======
function parseApiItem(item: DiaryItemAPI): DiaryItem {
  let images: string[] = []
  let videos: string[] = []
  try { images = JSON.parse(item.images || '[]') } catch { images = [] }
  try { videos = JSON.parse(item.videos || '[]') } catch { videos = [] }
  return {
    id: item.id,
    content: item.content || '',
    images,
    videos,
    happenedAt: item.happenedAt,
  }
}

const parsedItems = computed<DiaryItem[]>(() => {
  const raw = itemsData.value?.items ?? []
  return raw.map(parseApiItem).sort((a, b) => new Date(b.happenedAt).getTime() - new Date(a.happenedAt).getTime())
})

// ====== 前端过滤 ======
const filteredItems = computed<DiaryItem[]>(() => {
  let items = parsedItems.value
  // 时间筛选
  if (activeMonth.value !== '全部') {
    items = items.filter(i => i.happenedAt.startsWith(activeMonth.value))
  }
  // 搜索过滤
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    items = items.filter(i => i.content.toLowerCase().includes(q))
  }
  return items
})

// 按月分组
const filteredGroups = computed<MonthGroup[]>(() => {
  const map = new Map<string, DiaryItem[]>()
  for (const item of filteredItems.value) {
    const month = item.happenedAt.slice(0, 7) // "2026-08"
    if (!map.has(month)) map.set(month, [])
    map.get(month)!.push(item)
  }
  // 按月份降序
  return Array.from(map.entries())
    .sort(([a], [b]) => b.localeCompare(a))
    .map(([month, items]) => ({ month, items }))
})

// ====== 格式化 ======
function formatMonthLabel(month: string): string {
  const [year, m] = month.split('-')
  return `${year} 年 ${parseInt(m)} 月`
}

function formatDateTime(iso: string): string {
  const d = new Date(iso)
  return `${d.getMonth() + 1} 月 ${d.getDate()} 日 ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

// ====== 图片网格样式 ======
function imageGridClass(count: number): string {
  if (count === 1) return ''
  if (count === 2) return 'grid grid-cols-2 gap-2'
  if (count === 3) return 'grid grid-cols-3 gap-2'
  if (count === 4) return 'grid grid-cols-2 gap-2'
  return 'grid grid-cols-3 gap-2'
}

function imageItemClass(count: number, _idx: number): string {
  if (count === 1) return 'max-w-md aspect-auto max-h-80'
  if (count <= 4) return 'aspect-square'
  return 'aspect-square'
}

// ====== 灯箱 ======
function openLightbox(images: string[], index: number) {
  lightboxImages.value = images
  lightboxIndex.value = index
  lightboxVisible.value = true
}

// ====== 编辑模式 ======
function enterEditMode() {
  draftItems.value = parsedItems.value.map(item => ({
    _uid: genUid(),
    id: item.id,
    content: item.content,
    images: [...item.images],
    videos: [...item.videos],
    happenedAt: toLocalDatetime(item.happenedAt),
  }))
  editing.value = true
}

function enterEditModeWithNewRow() {
  draftItems.value = parsedItems.value.map(item => ({
    _uid: genUid(),
    id: item.id,
    content: item.content,
    images: [...item.images],
    videos: [...item.videos],
    happenedAt: toLocalDatetime(item.happenedAt),
  }))
  editing.value = true
  addDraftRow()
}

function toLocalDatetime(iso: string): string {
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function removeDraftRow(uid: string) {
  draftItems.value = draftItems.value.filter(i => i._uid !== uid)
}

function addDraftRow() {
  const now = new Date()
  const pad = (n: number) => String(n).padStart(2, '0')
  const defaultTime = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}T${pad(now.getHours())}:${pad(now.getMinutes())}`
  draftItems.value.push({
    _uid: genUid(),
    content: '',
    images: [],
    videos: [],
    happenedAt: defaultTime,
  })
}

function handleCancel() {
  const hasChanges = draftItems.value.length !== parsedItems.value.length ||
    draftItems.value.some((d, i) => {
      const o = parsedItems.value[i]
      if (!o) return true
      return d.content !== o.content ||
        JSON.stringify(d.images) !== JSON.stringify(o.images) ||
        JSON.stringify(d.videos) !== JSON.stringify(o.videos)
    })
  if (hasChanges) {
    showConfirm({
      title: '确认放弃修改？',
      content: '所有未保存的改动将丢失',
      okText: '放弃',
      cancelText: '继续编辑',
      danger: true,
      onOk: () => { editing.value = false },
    })
  } else {
    editing.value = false
  }
}

// ====== 上传 ======
async function handleImageUpload(e: Event, item: DraftItem) {
  const input = e.target as HTMLInputElement
  if (!input.files) return
  for (const file of Array.from(input.files)) {
    const formData = new FormData()
    formData.append('file', file)
    try {
      const res = await $fetch<{ url: string }>('/api/upload', { method: 'POST', body: formData })
      if (res.url) item.images.push(res.url)
    } catch (err: any) {
      showError(`图片上传失败: ${err?.message || '未知错误'}`)
    }
  }
  input.value = ''
}

async function handleVideoUpload(e: Event, item: DraftItem) {
  const input = e.target as HTMLInputElement
  if (!input.files) return
  for (const file of Array.from(input.files)) {
    if (file.size > 100 * 1024 * 1024) {
      showError('视频文件不能超过 100MB')
      continue
    }
    const formData = new FormData()
    formData.append('file', file)
    try {
      const res = await $fetch<{ url: string }>('/api/upload', { method: 'POST', body: formData })
      if (res.url) item.videos.push(res.url)
    } catch (err: any) {
      showError(`视频上传失败: ${err?.message || '未知错误'}`)
    }
  }
  input.value = ''
}

// ====== 保存 ======
async function handleSave() {
  // 校验
  for (let i = 0; i < draftItems.value.length; i++) {
    const item = draftItems.value[i]
    if (!item.content.trim() && item.images.length === 0 && item.videos.length === 0) {
      showError(`第 ${i + 1} 条：内容、图片、视频至少填写一项`)
      return
    }
    if (!item.happenedAt) {
      showError(`第 ${i + 1} 条：发生时间不能为空`)
      return
    }
  }

  saving.value = true
  try {
    const creates: any[] = []
    const updates: any[] = []
    const deletes: number[] = []

    const origMap = new Map<number, DiaryItem>()
    for (const item of parsedItems.value) {
      origMap.set(item.id, item)
    }

    const draftIdSet = new Set<number>()
    for (const item of draftItems.value) {
      const happenedAtISO = new Date(item.happenedAt).toISOString()
      if (!item.id) {
        creates.push({
          content: item.content,
          images: JSON.stringify(item.images),
          videos: JSON.stringify(item.videos),
          happened_at: happenedAtISO,
        })
      } else {
        draftIdSet.add(item.id)
        const orig = origMap.get(item.id)
        if (orig && (
          orig.content !== item.content ||
          JSON.stringify(orig.images) !== JSON.stringify(item.images) ||
          JSON.stringify(orig.videos) !== JSON.stringify(item.videos) ||
          orig.happenedAt !== happenedAtISO
        )) {
          updates.push({
            id: item.id,
            content: item.content,
            images: JSON.stringify(item.images),
            videos: JSON.stringify(item.videos),
            happened_at: happenedAtISO,
          })
        }
      }
    }

    for (const item of parsedItems.value) {
      if (!draftIdSet.has(item.id)) {
        deletes.push(item.id)
      }
    }

    const body = { creates, updates, deletes }
    await $fetch('/api/growth-diary/items/batch', { method: 'POST', body })

    // 刷新数据
    await refreshNuxtData('growth-diary-items')
    await refreshNuxtData('growth-diary-months')
    editing.value = false
    showSuccess('保存成功')
  } catch (err: any) {
    showError(err?.data?.statusMessage || err?.data?.msg || err?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

// ====== SEO ======
useSeoMeta({
  title: '阳阳的成长日记 - fatwill 的小屋',
  description: '记录阳阳成长的每一个精彩瞬间',
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
