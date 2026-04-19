<template>
  <Teleport to="body">
    <Transition name="media-viewer-fade">
      <div
        v-if="visible"
        ref="containerRef"
        class="fixed inset-0 z-[9999] bg-black touch-none"
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

        <!-- 图片预览 -->
        <div v-if="currentItem?.type === 'image'" class="flex h-full w-full items-center justify-center">
          <img
            ref="imageRef"
            :src="currentItem.url"
            :alt="currentItem.name"
            class="max-h-full max-w-full select-none media-img"
            draggable="false"
          />
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
const containerRef = ref<HTMLDivElement | null>(null)
const imageRef = ref<HTMLImageElement | null>(null)
const videoRef = ref<HTMLVideoElement | null>(null)

// 下载状态
const downloading = ref(false)

// 视频控制
const videoPaused = ref(true)
const videoCurrentTime = ref(0)
const videoDuration = ref(0)
const showVideoControls = ref(true)
let controlsTimer: ReturnType<typeof setTimeout> | null = null

// ====== 图片变换：全部用普通变量，完全脱离 Vue 响应式 ======
let liveScale = 1
let liveTranslateX = 0
let liveTranslateY = 0
let liveRotation = 0
let rafPending = false

// 触摸状态
let touchStartTime = 0
let touchStartX = 0
let touchStartY = 0
let isPinching = false
let initialPinchDistance = 0
let initialScale = 1
let isSwiping = false
let isDragging = false

// 双指缩放中心点
let pinchCenterX = 0
let pinchCenterY = 0

// 双击检测
let tapCount = 0
let tapTimer: ReturnType<typeof setTimeout> | null = null
let lastTapX = 0
let lastTapY = 0

const currentItem = computed(() => props.items[currentIndex.value])

// 直接写 DOM，animate=true 时加过渡动画（旋转、双击缩放用）
function applyTransform(animate = false) {
  if (!imageRef.value) return
  imageRef.value.style.transition = animate ? 'transform 0.25s ease' : 'none'
  imageRef.value.style.transform = `scale(${liveScale}) translate(${liveTranslateX}px, ${liveTranslateY}px) rotate(${liveRotation}deg)`
}

// rAF 节流版本，用于 pinch/drag 过程中
function applyTransformRaf() {
  if (rafPending) return
  rafPending = true
  requestAnimationFrame(() => {
    rafPending = false
    applyTransform(false)
  })
}

function resetImageState() {
  liveScale = 1
  liveTranslateX = 0
  liveTranslateY = 0
  liveRotation = 0
  if (imageRef.value) {
    imageRef.value.style.transition = 'none'
    imageRef.value.style.transform = 'scale(1) translate(0px, 0px) rotate(0deg)'
  }
}

// ====== 手动绑定 touch 事件（passive: false，可 preventDefault） ======
function bindTouchEvents() {
  if (!containerRef.value) return
  containerRef.value.addEventListener('touchstart', onTouchStart, { passive: false })
  containerRef.value.addEventListener('touchmove', onTouchMove, { passive: false })
  containerRef.value.addEventListener('touchend', onTouchEnd, { passive: true })
}

function unbindTouchEvents() {
  if (!containerRef.value) return
  containerRef.value.removeEventListener('touchstart', onTouchStart)
  containerRef.value.removeEventListener('touchmove', onTouchMove)
  containerRef.value.removeEventListener('touchend', onTouchEnd)
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
watch(() => props.visible, async (v) => {
  if (v) {
    currentIndex.value = props.initialIndex
    resetImageState()
    document.body.style.overflow = 'hidden'
    pushHistoryState()
    window.addEventListener('popstate', onPopState)
    await nextTick()
    bindTouchEvents()
  } else {
    unbindTouchEvents()
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
  liveRotation += 90
  liveTranslateX = 0
  liveTranslateY = 0
  applyTransform(true)
}

// ====== 触摸事件（手动绑定，passive: false） ======
function onTouchStart(e: TouchEvent) {
  touchStartTime = Date.now()

  if (e.touches.length === 2) {
    e.preventDefault() // 阻止浏览器默认双指缩放
    isPinching = true
    initialPinchDistance = getPinchDistance(e.touches)
    initialScale = liveScale
    // 记录双指中心点
    pinchCenterX = (e.touches[0].clientX + e.touches[1].clientX) / 2
    pinchCenterY = (e.touches[0].clientY + e.touches[1].clientY) / 2
    if (tapTimer) { clearTimeout(tapTimer); tapTimer = null; tapCount = 0 }
    return
  }

  if (e.touches.length === 1) {
    touchStartX = e.touches[0].clientX
    touchStartY = e.touches[0].clientY
    lastTapX = touchStartX
    lastTapY = touchStartY
    isSwiping = false
    isDragging = false
  }
}

function onTouchMove(e: TouchEvent) {
  if (e.touches.length === 2 && isPinching) {
    e.preventDefault() // 关键：阻止浏览器整页缩放！
    const dist = getPinchDistance(e.touches)
    liveScale = Math.max(0.5, Math.min(5, initialScale * (dist / initialPinchDistance)))
    applyTransformRaf()
    return
  }

  if (e.touches.length === 1) {
    const dx = e.touches[0].clientX - touchStartX
    const dy = e.touches[0].clientY - touchStartY

    if (liveScale > 1) {
      e.preventDefault() // 放大状态下拖动也阻止默认滚动
      isDragging = true
      liveTranslateX += dx / liveScale
      liveTranslateY += dy / liveScale
      touchStartX = e.touches[0].clientX
      touchStartY = e.touches[0].clientY
      applyTransformRaf()
    } else if (Math.abs(dx) > 8) {
      isSwiping = true
    }
  }
}

function onTouchEnd(e: TouchEvent) {
  const wasPinching = isPinching
  isPinching = false

  // pinch/drag 结束，应用最终状态（无动画，保持当前位置）
  if (wasPinching || isDragging) {
    applyTransform(false)
  }

  const elapsed = Date.now() - touchStartTime

  if (isSwiping && liveScale <= 1 && e.changedTouches.length > 0) {
    const totalDx = e.changedTouches[0].clientX - lastTapX
    if (Math.abs(totalDx) > 60) {
      if (totalDx < 0 && currentIndex.value < props.items.length - 1) {
        currentIndex.value++
      } else if (totalDx > 0 && currentIndex.value > 0) {
        currentIndex.value--
      }
    }
    isSwiping = false
    isDragging = false
    return
  }

  isSwiping = false
  isDragging = false

  if (elapsed > 300) return

  tapCount++
  if (tapCount === 1) {
    tapTimer = setTimeout(() => {
      if (currentItem.value?.type === 'image' && liveScale <= 1) {
        close()
      } else if (currentItem.value?.type === 'video') {
        toggleControls()
      }
      tapCount = 0
      tapTimer = null
    }, 250)
  } else if (tapCount === 2) {
    if (tapTimer) { clearTimeout(tapTimer); tapTimer = null }
    tapCount = 0
    handleDoubleTap()
  }
}

function handleDoubleTap() {
  if (currentItem.value?.type !== 'image') return
  if (liveScale > 1) {
    liveScale = 1
    liveTranslateX = 0
    liveTranslateY = 0
  } else {
    liveScale = 2.5
  }
  applyTransform(true)
}

function getPinchDistance(touches: TouchList) {
  const dx = touches[0].clientX - touches[1].clientX
  const dy = touches[0].clientY - touches[1].clientY
  return Math.sqrt(dx * dx + dy * dy)
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
  unbindTouchEvents()
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

/* GPU 加速，提升 transform 性能 */
.media-img {
  will-change: transform;
  transform: scale(1) translate(0px, 0px) rotate(0deg);
  backface-visibility: hidden;
  -webkit-backface-visibility: hidden;
}

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
