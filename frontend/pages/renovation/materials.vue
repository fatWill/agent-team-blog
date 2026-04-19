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
        <span class="text-sm font-medium text-gray-900 dark:text-gray-100">材料清单</span>
      </div>
    </header>

    <main class="mx-auto w-full max-w-4xl flex-1 px-4 py-6">
      <!-- 加载状态 -->
      <div v-if="loading" class="space-y-4">
        <div v-for="i in 3" :key="i" class="animate-pulse rounded-xl bg-white p-5 shadow-sm dark:bg-gray-800">
          <div class="mb-3 h-5 w-1/3 rounded bg-gray-200 dark:bg-gray-700" />
          <div class="mb-2 h-4 w-full rounded bg-gray-100 dark:bg-gray-700/60" />
          <div class="h-4 w-2/3 rounded bg-gray-100 dark:bg-gray-700/60" />
        </div>
      </div>

      <template v-else-if="items.length > 0">
        <!-- 标签过滤栏 -->
        <div class="mb-5 flex gap-2 overflow-x-auto pb-1 scrollbar-hide">
          <button
            class="shrink-0 rounded-full px-3.5 py-1.5 text-sm font-medium transition-all"
            :class="activeTag === null
              ? 'bg-gray-900 text-white shadow-sm dark:bg-gray-100 dark:text-gray-900'
              : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700'"
            @click="activeTag = null"
          >
            全部
          </button>
          <button
            v-for="tag in allTags"
            :key="tag"
            class="shrink-0 rounded-full px-3.5 py-1.5 text-sm font-medium transition-all"
            :style="activeTag === tag ? tagSelectedStyle(tag) : tagNormalStyle(tag)"
            @click="activeTag = tag"
          >
            {{ tag }}
          </button>
        </div>

        <!-- 条目列表 -->
        <div class="space-y-4">
          <div
            v-for="item in filteredItems"
            :key="item.id"
            class="rounded-xl bg-white p-5 shadow-sm transition-colors dark:bg-gray-800"
          >
            <!-- 标题 + 标签 -->
            <div class="mb-3 flex flex-wrap items-center gap-2">
              <h3 class="text-base font-bold text-gray-900 dark:text-gray-100">{{ item.title }}</h3>
              <span
                v-for="tag in item.tags"
                :key="tag"
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :style="tagNormalStyle(tag)"
              >
                {{ tag }}
              </span>
            </div>

            <!-- 附件网格 -->
            <div v-if="item.attachments.length > 0" class="grid grid-cols-3 gap-2.5">
              <template v-for="att in item.attachments" :key="att.id">
                <!-- PDF -->
                <button
                  v-if="att.type === 'pdf'"
                  class="flex flex-col items-center gap-1.5 rounded-lg border border-gray-100 bg-gray-50 p-3 transition-colors hover:bg-red-50 dark:border-gray-700 dark:bg-gray-700/50 dark:hover:bg-red-900/10"
                  @click="openPdf(att)"
                >
                  <svg class="h-8 w-8 text-red-500" fill="currentColor" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 2l5 5h-5V4zM6 20V4h5v7h7v9H6z"/><path d="M8 13h3v1H8zm0 2h6v1H8zm0 2h6v1H8z"/></svg>
                  <span class="line-clamp-2 text-center text-xs text-gray-600 dark:text-gray-400">{{ att.name }}</span>
                </button>

                <!-- 图片 -->
                <button
                  v-else-if="att.type === 'image'"
                  class="group relative aspect-square overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-700"
                  @click="openMediaPreview(item, att)"
                >
                  <img :src="att.url" :alt="att.name" class="h-full w-full object-cover transition-transform group-hover:scale-105" loading="lazy" />
                </button>

                <!-- 视频 -->
                <button
                  v-else-if="att.type === 'video'"
                  class="group relative aspect-square overflow-hidden rounded-lg bg-gray-900"
                  @click="openMediaPreview(item, att)"
                >
                  <img v-if="att.thumbnailUrl" :src="att.thumbnailUrl" :alt="att.name" class="h-full w-full object-cover opacity-80 transition-transform group-hover:scale-105" loading="lazy" />
                  <div v-else class="flex h-full w-full items-center justify-center">
                    <span class="text-3xl">🎬</span>
                  </div>
                  <!-- 播放按钮覆盖 -->
                  <div class="absolute inset-0 flex items-center justify-center">
                    <div class="flex h-10 w-10 items-center justify-center rounded-full bg-white/80 shadow backdrop-blur">
                      <svg class="ml-0.5 h-5 w-5 text-gray-900" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
                    </div>
                  </div>
                  <span v-if="att.duration" class="absolute bottom-1 right-1 rounded bg-black/60 px-1.5 py-0.5 text-[10px] text-white">{{ formatDuration(att.duration) }}</span>
                </button>
              </template>
            </div>
          </div>
        </div>
      </template>

      <!-- 空状态 -->
      <div v-else class="py-20 text-center">
        <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-50 dark:bg-green-900/20">
          <span class="text-3xl">🧱</span>
        </div>
        <h2 class="text-lg font-bold text-gray-900 dark:text-gray-100">暂无材料清单</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">管理员可以在后台添加材料信息</p>
      </div>
    </main>

    <!-- 微信风格预览 -->
    <MediaViewer
      :visible="viewerVisible"
      :items="viewerItems"
      :initial-index="viewerIndex"
      @close="viewerVisible = false"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

import type { MaterialItem, MaterialAttachment } from '~/features/material'
import { apiFetchMaterials } from '~/features/material'

const router = useRouter()

const loading = ref(true)
const items = ref<MaterialItem[]>([])
const activeTag = ref<string | null>(null)

// 媒体预览
const viewerVisible = ref(false)
const viewerItems = ref<Array<{ type: 'image' | 'video'; url: string; name: string; thumbnailUrl?: string; duration?: number }>>([])
const viewerIndex = ref(0)

// 加载数据（客户端）
onMounted(async () => {
  try {
    items.value = await apiFetchMaterials()
  } catch (err: any) {
    // 未登录（401/403）则跳转登录页
    const status = err?.response?.status ?? err?.statusCode ?? err?.status
    if (status === 401 || status === 403) {
      await navigateTo('/login?redirect=' + encodeURIComponent('/renovation/materials'), { replace: true })
      return
    }
  } finally {
    loading.value = false
  }
})

// 提取所有唯一标签
const allTags = computed(() => {
  const set = new Set<string>()
  for (const item of items.value) {
    for (const tag of item.tags) set.add(tag)
  }
  return Array.from(set)
})

// 根据标签过滤
const filteredItems = computed(() => {
  if (activeTag.value === null) return items.value
  return items.value.filter(item => item.tags.includes(activeTag.value!))
})

// ====== 标签颜色 ======
function hashColor(str: string, lightness: number): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  const hue = Math.abs(hash) % 360
  return `hsl(${hue}, 70%, ${lightness}%)`
}

function tagNormalStyle(tag: string) {
  const color = hashColor(tag, 55)
  return {
    backgroundColor: `${hashColor(tag, 95)}`,
    color,
  }
}

function tagSelectedStyle(tag: string) {
  return {
    backgroundColor: hashColor(tag, 55),
    color: '#fff',
    boxShadow: `0 2px 8px ${hashColor(tag, 70)}40`,
  }
}

// ====== PDF 打开 ======
function openPdf(att: MaterialAttachment) {
  localStorage.setItem(`pdf_cache_${att.id}`, JSON.stringify({ id: att.id, url: att.url, name: att.name }))
  router.push(`/pdf/${att.id}`)
}

// ====== 媒体预览 ======
function openMediaPreview(item: MaterialItem, att: MaterialAttachment) {
  const mediaAttachments = item.attachments.filter(a => a.type === 'image' || a.type === 'video')
  viewerItems.value = mediaAttachments.map(a => ({
    type: a.type as 'image' | 'video',
    url: a.url,
    name: a.name,
    thumbnailUrl: a.thumbnailUrl,
    duration: a.duration,
  }))
  viewerIndex.value = mediaAttachments.findIndex(a => a.id === att.id)
  if (viewerIndex.value < 0) viewerIndex.value = 0
  viewerVisible.value = true
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

useSeoMeta({
  title: '材料清单 - 装修',
  description: '装修材料清单',
})
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar { display: none; }
.scrollbar-hide { -ms-overflow-style: none; scrollbar-width: none; }
</style>
