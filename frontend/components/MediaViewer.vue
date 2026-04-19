<template>
  <Teleport to="body">
    <Transition name="media-viewer-fade">
      <div
        v-if="visible"
        class="fixed inset-0 z-[9999] bg-black touch-none select-none"
        @touchstart="onTouchStart"
        @touchmove.prevent="onTouchMove"
        @touchend="onTouchEnd"
      >
        <!-- 顶部栏 -->
        <div class="absolute left-0 right-0 top-0 z-10 flex items-center justify-between px-4 py-3">
          <span v-if="items.length > 1" class="text-sm font-medium text-white/80">{{ currentIndex + 1 }}/{{ items.length }}</span>
          <span v-else />
          <button
            class="flex h-8 w-8 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur"
            @touchstart.stop
            @click="close"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>

        <!-- 图片预览（Swiper 容器） -->
        <div v-if="currentItem?.type === 'image'" class="absolute inset-0 overflow-hidden">
          <div
            class="flex h-full"
            :style="{
              width: `${imageItems.length * 100}%`,
              transform: `translateX(calc(-${currentImageIdx * (100 / imageItems.length)}% + ${swipeX / imageItems.length}px))`,
              transition: isSwiping ? 'none' : 'transform 0.3s cubic-bezier(0.4,0,0.2,1)',
            }"
          >
            <div
              v-for="(item, idx) in imageItems"
              :key="item.url"
              class="flex items-center justify-center flex-shrink-0"
              :style="{ width: `${100 / imageItems.length}%`, height: '100vh' }"
            >
              <img
                v-if="Math.abs(idx - currentImageIdx) <= 1"
                :src="item.url"
                :alt="item.name"
                class="max-h-full max-w-full select-none"
                :style="idx === currentImageIdx ? {
                  transform: `scale(${imgScale}) translate(${imgPanX / imgScale}px, ${imgPanY / imgScale}px) rotate(${imgRotation}deg)`,
                  transition: (isPanning || isPinching) ? 'none' : 'transform 0.25s cubic-bezier(0.4,0,0.2,1)',
                  willChange: 'transform',
                } : {}"
                draggable="false"
              />
            </div>
          </div>
        </div>

        <!-- 视频预览 -->
        <div v-else-if="currentItem?.type === 'video'" class="flex h-full w-full flex-col">
          <div class="relative flex flex-1 items-center justify-center" @click="toggleControls">
            <video
              ref="videoRef"
              :src="currentItem.url"
              class="max-h-full max-w-full object-contain"
              playsinline
              preload="metadata"
              @timeupdate="onTimeUpdate"
              @loadedmetadata="onVideoLoaded"
              @ended="videoPaused = true"
              @play="videoPaused = false"
              @pause="videoPaused = true"
            />
          </div>
          <Transition name="fade">
            <div v-show="showVideoControls" class="absolute bottom-0 left-0 right-0 z-10 bg-gradient-to-t from-black/70 to-transparent px-4 pb-6 pt-10">
              <div class="flex items-center gap-3">
                <button class="flex h-8 w-8 shrink-0 items-center justify-center text-white" @click.stop="togglePlay">
                  <svg v-if="videoPaused" class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
                  <svg v-else class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24"><path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z" /></svg>
                </button>
                <span class="shrink-0 text-xs text-white/70">{{ formatTime(videoCurrentTime) }}</span>
                <input
                  type="range"
                  class="video-progress flex-1"
                  :value="videoCurrentTime"
                  :max="videoDuration"
                  step="0.1"
                  @input="onSeek"
                  @click.stop
                />
                <span class="shrink-0 text-xs text-white/70">{{ formatTime(videoDuration) }}</span>
                <button class="flex h-8 w-8 shrink-0 items-center justify-center text-white" @click.stop="enterLandscape">
                  <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 3.75v4.5m0-4.5h4.5m-4.5 0L9 9M3.75 20.25v-4.5m0 4.5h4.5m-4.5 0L9 15M20.25 3.75h-4.5m4.5 0v4.5m0-4.5L15 9m5.25 11.25h-4.5m4.5 0v-4.5m0 4.5L15 15" /></svg>
                </button>
              </div>
            </div>
          </Transition>
        </div>

        <!-- 右下角操作按钮 -->
        <div class="absolute bottom-8 right-4 z-20 flex flex-row items-center gap-2">
          <button
            v-if="currentItem?.type === 'image'"
            class="flex h-8 w-8 items-center justify-center rounded-full bg-black/40 text-white backdrop-blur-sm active:opacity-70"
            @touchstart.stop
            @click.stop="rotateImage"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
          </button>
          <button
            class="flex h-8 w-8 items-center justify-center rounded-full bg-black/40 text-white backdrop-blur-sm active:opacity-70"
            :class="downloading ? 'opacity-50' : ''"
            :disabled="downloading"
            @touchstart.stop
            @click.stop="handleDownload"
          >
            <svg v-if="!downloading" class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
            </svg>
            <svg v-else class="h-4 w-4 animate-spin" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
interface MediaItem {
  type: 'image' | 'video'
  url: string
  name: string
  thumbnailUrl?: string
  duration?: number
}

const props = withDefaults(defineProps<{
  visible: boolean
  items: MediaItem[]
  initialIndex?: number
}>(), { initialIndex: 0 })

const emit = defineEmits<{ close: [] }>()

const currentIndex = ref(props.initialIndex)
const videoRef = ref<HTMLVideoElement | null>(null)

// 下载状态
const downloading = ref(false)

// 视频控制
const videoPaused = ref(true)
const videoCurrentTime = ref(0)
const videoDuration = ref(0)
const showVideoControls = ref(true)
let controlsTimer: ReturnType<typeof setTimeout> | null = null

// ====== 图片变换状态（Vue 响应式 + :style 绑定，复用 home.vue 灯箱方案） ======
const imgScale = ref(1)
const imgPanX = ref(0)
const imgPanY = ref(0)
const imgRotation = ref(0) // 累加，不取模
const swipeX = ref(0)      // Swiper 容器的水平拖动偏移
const isPinching = ref(false)
const isSwiping = ref(false)
const isPanning = ref(false)

// 触摸状态（普通变量，不需要模板访问）
let touchStartX = 0
let touchStartY = 0
let touchStartDist = 0
let touchStartScale = 1
let touchStartPanX = 0
let touchStartPanY = 0
let pinchCenterX = 0
let pinchCenterY = 0
let isTouchActive = false

// 双击检测
let touchStartTime = 0
let tapCount = 0
let tapTimer: ReturnType<typeof setTimeout> | null = null

// 只取图片类型的 items（用于 Swiper 容器）
const imageItems = computed(() => props.items.filter(i => i.type === 'image'))
// 当前图片在 imageItems 中的索引
const currentImageIdx = computed(() => {
  const item = props.items[currentIndex.value]
  if (!item || item.type !== 'image') return 0
  return imageItems.value.findIndex(i => i.url === item.url)
})

const currentItem = computed(() => props.items[currentIndex.value])

function resetImageState() {
  imgScale.value = 1
  imgPanX.value = 0
  imgPanY.value = 0
  imgRotation.value = 0
  swipeX.value = 0
  isPinching.value = false
  isSwiping.value = false
  isPanning.value = false
}

// ====== History API：浏览器返回关闭预览 ======
let historyPushed = false

function pushHistoryState() {
  history.pushState({ mediaViewer: true }, '')
  historyPushed = true
}

function onPopState() {
  historyPushed = false
  emit('close')
}

// ====== watch visible ======
watch(() => props.visible, (v) => {
  if (v) {
    currentIndex.value = props.initialIndex
    resetImageState()
    document.body.style.overflow = 'hidden'
    pushHistoryState()
    window.addEventListener('popstate', onPopState)
  } else {
    window.removeEventListener('popstate', onPopState)
    document.body.style.overflow = ''
    resetVideoState()
    if (historyPushed) {
      historyPushed = false
      history.back()
    }
  }
})

watch(currentIndex, () => {
  resetImageState()
  resetVideoState()
})

function resetVideoState() {
  if (videoRef.value) {
    videoRef.value.pause()
    videoRef.value.currentTime = 0
  }
  videoPaused.value = true
  videoCurrentTime.value = 0
  showVideoControls.value = true
}

function close() {
  emit('close')
}

function rotateImage() {
  imgRotation.value += 90
  imgPanX.value = 0
  imgPanY.value = 0
}

// ====== 触摸事件（完全复用 home.vue 灯箱逻辑） ======
function getTouchDist(touches: TouchList): number {
  if (touches.length < 2) return 0
  const dx = touches[0].clientX - touches[1].clientX
  const dy = touches[0].clientY - touches[1].clientY
  return Math.sqrt(dx * dx + dy * dy)
}

function onTouchStart(e: TouchEvent) {
  isTouchActive = true
  touchStartTime = Date.now()

  if (e.touches.length === 2) {
    isPinching.value = true
    isSwiping.value = false
    isPanning.value = false
    touchStartDist = getTouchDist(e.touches)
    touchStartScale = imgScale.value
    touchStartPanX = imgPanX.value
    touchStartPanY = imgPanY.value
    pinchCenterX = (e.touches[0].clientX + e.touches[1].clientX) / 2 - window.innerWidth / 2
    pinchCenterY = (e.touches[0].clientY + e.touches[1].clientY) / 2 - window.innerHeight / 2
    // 取消双击定时器
    if (tapTimer) { clearTimeout(tapTimer); tapTimer = null; tapCount = 0 }
  } else if (e.touches.length === 1) {
    isPinching.value = false
    touchStartX = e.touches[0].clientX
    touchStartY = e.touches[0].clientY
    touchStartPanX = imgPanX.value
    touchStartPanY = imgPanY.value
    if (imgScale.value > 1) {
      isPanning.value = true
      isSwiping.value = false
    } else if (imageItems.value.length > 1) {
      isSwiping.value = true
      isPanning.value = false
    }
  }
}

function onTouchMove(e: TouchEvent) {
  if (!isTouchActive) return

  if (isPinching.value && e.touches.length === 2) {
    const dist = getTouchDist(e.touches)
    const ratio = dist / touchStartDist
    const newScale = Math.min(4, Math.max(0.2, touchStartScale * ratio))
    const scaleDelta = newScale - touchStartScale
    imgScale.value = newScale

    if (newScale > 1) {
      const newPanX = touchStartPanX - pinchCenterX * scaleDelta / touchStartScale
      const newPanY = touchStartPanY - pinchCenterY * scaleDelta / touchStartScale
      const maxPan = (newScale - 1) * window.innerWidth * 0.5
      const maxPanY = (newScale - 1) * window.innerHeight * 0.5
      imgPanX.value = Math.max(-maxPan, Math.min(maxPan, newPanX))
      imgPanY.value = Math.max(-maxPanY, Math.min(maxPanY, newPanY))
    } else {
      imgPanX.value = 0
      imgPanY.value = 0
    }
    return
  }

  if (e.touches.length !== 1) return
  const dx = e.touches[0].clientX - touchStartX
  const dy = e.touches[0].clientY - touchStartY

  if (imgScale.value > 1 && isPanning.value) {
    const maxPan = (imgScale.value - 1) * window.innerWidth * 0.5
    const maxPanY = (imgScale.value - 1) * window.innerHeight * 0.5
    imgPanX.value = Math.max(-maxPan, Math.min(maxPan, touchStartPanX + dx))
    imgPanY.value = Math.max(-maxPanY, Math.min(maxPanY, touchStartPanY + dy))
  } else if (imgScale.value <= 1 && isSwiping.value) {
    const resistance = Math.abs(dy) > Math.abs(dx) ? 0.3 : 1
    swipeX.value = dx * resistance
  }
}

function onTouchEnd(e: TouchEvent) {
  if (!isTouchActive) return
  isTouchActive = false

  if (isPinching.value) {
    isPinching.value = false
    if (imgScale.value <= 1) {
      imgScale.value = 1
      imgPanX.value = 0
      imgPanY.value = 0
    }
    isPanning.value = false
    return
  }

  if (imgScale.value <= 1 && isSwiping.value) {
    const threshold = window.innerWidth * 0.25
    if (swipeX.value < -threshold && currentIndex.value < props.items.length - 1) {
      currentIndex.value++
    } else if (swipeX.value > threshold && currentIndex.value > 0) {
      currentIndex.value--
    }
    swipeX.value = 0
    isSwiping.value = false
  }

  isPanning.value = false

  // 双击缩放检测
  const elapsed = Date.now() - touchStartTime
  if (elapsed < 300 && e.changedTouches.length === 1) {
    const movedX = Math.abs(e.changedTouches[0].clientX - touchStartX)
    const movedY = Math.abs(e.changedTouches[0].clientY - touchStartY)
    // 手指移动过多不算点击
    if (movedX > 10 || movedY > 10) return

    tapCount++
    if (tapCount === 1) {
      tapTimer = setTimeout(() => {
        tapCount = 0
        tapTimer = null
        // 单击：缩放为1时关闭
        if (currentItem.value?.type === 'image' && imgScale.value <= 1) {
          close()
        } else if (currentItem.value?.type === 'video') {
          toggleControls()
        }
      }, 250)
    } else if (tapCount === 2) {
      if (tapTimer) { clearTimeout(tapTimer); tapTimer = null }
      tapCount = 0
      handleDoubleTap()
    }
  }
}

function handleDoubleTap() {
  if (currentItem.value?.type !== 'image') return
  if (imgScale.value > 1) {
    imgScale.value = 1
    imgPanX.value = 0
    imgPanY.value = 0
  } else {
    imgScale.value = 2.5
  }
}

// ====== 视频控制 ======
function togglePlay() {
  if (!videoRef.value) return
  if (videoRef.value.paused) {
    videoRef.value.play()
    autoHideControls()
  } else {
    videoRef.value.pause()
    showVideoControls.value = true
    if (controlsTimer) clearTimeout(controlsTimer)
  }
}

function toggleControls() {
  showVideoControls.value = !showVideoControls.value
  if (showVideoControls.value && !videoPaused.value) {
    autoHideControls()
  }
}

function autoHideControls() {
  if (controlsTimer) clearTimeout(controlsTimer)
  controlsTimer = setTimeout(() => {
    if (!videoPaused.value) showVideoControls.value = false
  }, 3000)
}

function onTimeUpdate() {
  if (videoRef.value) videoCurrentTime.value = videoRef.value.currentTime
}

function onVideoLoaded() {
  if (videoRef.value) videoDuration.value = videoRef.value.duration
}

function onSeek(e: Event) {
  const val = Number((e.target as HTMLInputElement).value)
  if (videoRef.value) videoRef.value.currentTime = val
}

function formatTime(s: number): string {
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

async function enterLandscape() {
  try {
    await (screen.orientation as any).lock('landscape')
  } catch { /* fallback 静默 */ }
}

async function handleDownload() {
  if (!currentItem.value || downloading.value) return
  downloading.value = true
  try {
    const url = `/api/download?url=${encodeURIComponent(currentItem.value.url)}&name=${encodeURIComponent(currentItem.value.name || 'download')}`
    const a = document.createElement('a')
    a.href = url
    a.download = currentItem.value.name || 'download'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    await new Promise(resolve => setTimeout(resolve, 1500))
  } catch {
    window.open(currentItem.value?.url, '_blank')
  } finally {
    downloading.value = false
  }
}

onBeforeUnmount(() => {
  window.removeEventListener('popstate', onPopState)
  if (tapTimer) clearTimeout(tapTimer)
  if (controlsTimer) clearTimeout(controlsTimer)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.media-viewer-fade-enter-active,
.media-viewer-fade-leave-active { transition: opacity 0.25s ease; }
.media-viewer-fade-enter-from,
.media-viewer-fade-leave-to { opacity: 0; }

.fade-enter-active,
.fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from,
.fade-leave-to { opacity: 0; }

.video-progress {
  -webkit-appearance: none;
  appearance: none;
  height: 3px;
  background: rgba(255,255,255,0.3);
  border-radius: 2px;
  outline: none;
}
.video-progress::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #fff;
  cursor: pointer;
}
.video-progress::-moz-range-thumb {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #fff;
  border: none;
  cursor: pointer;
}
</style>
