<template>
  <div class="flex h-screen flex-col bg-gray-100 dark:bg-gray-900">
    <!-- 顶部栏 -->
    <header class="flex h-12 shrink-0 items-center gap-3 border-b border-gray-200 bg-white px-4 dark:border-gray-700 dark:bg-gray-800">
      <button class="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200" @click="goBack">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" /></svg>
        返回
      </button>
      <span class="truncate text-sm font-medium text-gray-900 dark:text-gray-100">{{ pdfName || 'PDF 预览' }}</span>
      <span v-if="totalPages > 0" class="ml-auto shrink-0 text-xs text-gray-400 dark:text-gray-500">{{ currentPage }} / {{ totalPages }}</span>
    </header>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <div class="text-center">
        <div class="mx-auto mb-3 h-10 w-10 animate-spin rounded-full border-4 border-gray-200 border-t-primary-500" />
        <p class="text-sm text-gray-500 dark:text-gray-400">加载 PDF 中...</p>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="flex flex-1 items-center justify-center px-4">
      <div class="text-center">
        <p class="text-lg font-medium text-gray-600 dark:text-gray-400">{{ error }}</p>
        <button class="mt-4 text-sm text-primary-500 hover:text-primary-600" @click="goBack">返回上一页</button>
      </div>
    </div>

    <!-- PDF 渲染区域 -->
    <div v-else ref="scrollContainer" class="flex-1 overflow-y-auto" @scroll="onScroll">
      <div class="mx-auto max-w-4xl py-4">
        <canvas
          v-for="page in totalPages"
          :key="page"
          :ref="(el) => setCanvasRef(el as HTMLCanvasElement, page)"
          class="mx-auto mb-4 block shadow-md"
        />
      </div>
    </div>

    <!-- 底部翻页栏 -->
    <div v-if="totalPages > 1" class="flex h-12 shrink-0 items-center justify-center gap-4 border-t border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800">
      <button
        class="rounded-lg px-4 py-1.5 text-sm font-medium transition-colors"
        :class="currentPage > 1 ? 'text-primary-600 hover:bg-primary-50 dark:text-primary-400 dark:hover:bg-primary-900/20' : 'cursor-not-allowed text-gray-300 dark:text-gray-600'"
        :disabled="currentPage <= 1"
        @click="scrollToPage(currentPage - 1)"
      >
        上一页
      </button>
      <span class="text-sm text-gray-500 dark:text-gray-400">{{ currentPage }} / {{ totalPages }}</span>
      <button
        class="rounded-lg px-4 py-1.5 text-sm font-medium transition-colors"
        :class="currentPage < totalPages ? 'text-primary-600 hover:bg-primary-50 dark:text-primary-400 dark:hover:bg-primary-900/20' : 'cursor-not-allowed text-gray-300 dark:text-gray-600'"
        :disabled="currentPage >= totalPages"
        @click="scrollToPage(currentPage + 1)"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ ssr: false })

const route = useRoute()
const router = useRouter()
const attachmentId = route.params.id as string

const loading = ref(true)
const error = ref('')
const pdfName = ref('')
const totalPages = ref(0)
const currentPage = ref(1)
const scrollContainer = ref<HTMLElement | null>(null)
const canvasRefs = new Map<number, HTMLCanvasElement>()

function setCanvasRef(el: HTMLCanvasElement | null, page: number) {
  if (el) canvasRefs.set(page, el)
}

function goBack() {
  router.back()
}

// 在客户端挂载后读取 localStorage 并加载 PDF
onMounted(() => {
  const cacheKey = `pdf_cache_${attachmentId}`
  const cached = localStorage.getItem(cacheKey)

  if (!cached) {
    loading.value = false
    error.value = 'PDF 信息未找到，请从材料清单页面打开'
  } else {
    const pdfInfo = JSON.parse(cached) as { id: string; url: string; name: string }
    pdfName.value = pdfInfo.name
    loadPdf(pdfInfo.url)
  }
})

async function loadPdf(url: string) {
  try {
    const pdfjsLib = await import('pdfjs-dist')
    pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/4.4.168/pdf.worker.min.mjs'

    const loadingTask = pdfjsLib.getDocument(url)
    const pdf = await loadingTask.promise
    totalPages.value = pdf.numPages
    loading.value = false

    // 等待 DOM 更新后渲染所有页面
    await nextTick()

    for (let i = 1; i <= pdf.numPages; i++) {
      const page = await pdf.getPage(i)
      const viewport = page.getViewport({ scale: 1.5 })
      const canvas = canvasRefs.get(i)
      if (!canvas) continue

      canvas.width = viewport.width
      canvas.height = viewport.height

      const ctx = canvas.getContext('2d')
      if (!ctx) continue

      await page.render({ canvasContext: ctx, viewport, canvas } as any).promise
    }
  } catch (e) {
    loading.value = false
    error.value = 'PDF 加载失败，请稍后重试'
    console.error('PDF load error:', e)
  }
}

function onScroll() {
  if (!scrollContainer.value) return
  const containerTop = scrollContainer.value.scrollTop
  const containerHeight = scrollContainer.value.clientHeight
  const center = containerTop + containerHeight / 2

  let closestPage = 1
  let closestDist = Infinity

  for (const [page, canvas] of canvasRefs) {
    const top = canvas.offsetTop
    const mid = top + canvas.height / 2
    const dist = Math.abs(mid - center)
    if (dist < closestDist) {
      closestDist = dist
      closestPage = page
    }
  }
  currentPage.value = closestPage
}

function scrollToPage(page: number) {
  const canvas = canvasRefs.get(page)
  if (canvas && scrollContainer.value) {
    canvas.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

useSeoMeta({
  title: () => pdfName.value ? `${pdfName.value} - PDF 预览` : 'PDF 预览',
})
</script>
