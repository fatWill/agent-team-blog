<template>
  <Teleport to="body">
    <Transition name="media-viewer-fade">
      <div
        v-if="visible"
        class="fixed inset-0 z-[9999] bg-black"
        @touchstart.passive="onTouchStart"
        @touchmove.passive="onTouchMove"
        @touchend="onTouchEnd"
      >
        <!-- 顶部栏 -->
        <div class="absolute left-0 right-0 top-0 z-10 flex items-center justify-between px-4 py-3">
          <span v-if="items.length > 1" class="text-sm font-medium text-white/80">{{ currentIndex + 1 }}/{{ items.length }}</span>
          <span v-else />
          <button class="flex h-8 w-8 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur" @click="close">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>

        <!-- 图片预览 -->
        <div v-if="currentItem?.type === 'image'" class="flex h-full w-full items-center justify-center">
          <img
            ref="imageRef"
            :src="currentItem.url"
            :alt="currentItem.name"
            class="max-h-full max-w-full select-none"
            :style="imageTransformStyle"
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
          <!-- 视频控制栏 -->
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

        <!-- 右下角操作按钮（图片：旋转 + 下载；视频：下载） -->
        <div class="absolute bottom-8 right-4 z-20 flex flex-col items-center gap-3">
          <!-- 旋转按钮（仅图片显示） -->
          <button
            v-if="currentItem?.type === 'image'"
            class="flex h-10 w-10 items-center justify-center rounded-full bg-black/40 text-white backdrop-blur-sm transition-opacity active:opacity-70"
            @click.stop="rotateImage"
          >
            <!-- 顺时针旋转图标 -->
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
            </svg>
          </button>
          <!-- 下载按钮 -->
          <button
            class="flex h-10 w-10 items-center justify-center rounded-full bg-black/40 text-white backdrop-blur-sm transition-opacity active:opacity-70"
            :class="downloading ? 'opacity-50' : ''"
            :disabled="downloading"
            @click.stop="handleDownload"
          >
            <svg v-if="!downloading" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
            </svg>
            <!-- 下载中 loading -->
            <svg v-else class="h-5 w-5 animate-spin" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
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
const imageRef = ref<HTMLImageElement | null>(null)
const videoRef = ref<HTMLVideoElement | null>(null)

// 图片变换状态
const scale = ref(1)
const translateX = ref(0)
const translateY = ref(0)
const rotation = ref(0)

// 下载状态
const downloading = ref(false)

// 视频控制
const videoPaused = ref(true)
const videoCurrentTime = ref(0)
const videoDuration = ref(0)
const showVideoControls = ref(true)
let controlsTimer: ReturnType<typeof setTimeout> | null = null

// 触摸状态
let touchStartTime = 0
let touchStartX = 0
let touchStartY = 0
let isPinching = false
let initialPinchDistance = 0
let initialScale = 1
let isSwiping = false
let isDragging = false

// 双击检测
let tapCount = 0
let tapTimer: ReturnType<typeof setTimeout> | null = null
let lastTapX = 0
let lastTapY = 0

const currentItem = computed(() => props.items[currentIndex.value])

const imageTransformStyle = computed(() => ({
  transform: `scale(${scale.value}) translate(${translateX.value}px, ${translateY.value}px) rotate(${rotation.value}deg)`,
  transition: isPinching || isDragging ? 'none' : 'transform 0.2s ease',
}))

watch(() => props.visible, (v) => {
  if (v) {
    currentIndex.value = props.initialIndex
    resetImageState()
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
    resetVideoState()
  }
})

watch(currentIndex, () => {
  resetImageState()
  resetVideoState()
})

function resetImageState() {
  scale.value = 1
  translateX.value = 0
  translateY.value = 0
  rotation.value = 0
}

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
  rotation.value = (rotation.value + 90) % 360
}

// ====== 触摸事件 ======
function onTouchStart(e: TouchEvent) {
  touchStartTime = Date.now()

  if (e.touches.length === 2) {
    isPinching = true
    initialPinchDistance = getPinchDistance(e.touches)
    initialScale = scale.value
    // 双指时清除单击计时
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
    const dist = getPinchDistance(e.touches)
    const newScale = Math.max(0.5, Math.min(5, initialScale * (dist / initialPinchDistance)))
    scale.value = newScale
    return
  }

  if (e.touches.length === 1) {
    const dx = e.touches[0].clientX - touchStartX
    const dy = e.touches[0].clientY - touchStartY

    if (scale.value > 1) {
      isDragging = true
      translateX.value += dx / scale.value
      translateY.value += dy / scale.value
      touchStartX = e.touches[0].clientX
      touchStartY = e.touches[0].clientY
    } else if (Math.abs(dx) > 8) {
      isSwiping = true
    }
  }
}

function onTouchEnd(e: TouchEvent) {
  isPinching = false

  const elapsed = Date.now() - touchStartTime

  // 滑动切换图片/视频
  if (isSwiping && scale.value <= 1 && e.changedTouches.length > 0) {
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

  // 只处理短促点击（非拖动）
  if (elapsed > 300) return

  // 双击检测
  tapCount++
  if (tapCount === 1) {
    tapTimer = setTimeout(() => {
      // 单击：关闭
      if (currentItem.value?.type === 'image' && scale.value <= 1) {
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
  if (scale.value > 1) {
    scale.value = 1
    translateX.value = 0
    translateY.value = 0
  } else {
    scale.value = 2.5
  }
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
  } catch {
    // fallback 静默
  }
}

// ====== 下载（fetch blob 方式，避免跨域直接打开） ======
async function handleDownload() {
  if (!currentItem.value || downloading.value) return
  downloading.value = true
  try {
    const response = await fetch(currentItem.value.url)
    if (!response.ok) throw new Error('Download failed')
    const blob = await response.blob()
    const blobUrl = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = blobUrl
    a.download = currentItem.value.name || 'download'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    setTimeout(() => URL.revokeObjectURL(blobUrl), 1000)
  } catch {
    // 降级：直接打开新标签页
    window.open(currentItem.value.url, '_blank')
  } finally {
    downloading.value = false
  }
}

onBeforeUnmount(() => {
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
