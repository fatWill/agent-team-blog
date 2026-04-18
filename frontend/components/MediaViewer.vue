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
        <div v-if="currentItem?.type === 'image'" class="flex h-full w-full items-center justify-center" @click="handleSingleTap">
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

        <!-- 长按操作栏 -->
        <Transition name="action-sheet">
          <div v-if="showActionSheet" class="absolute bottom-0 left-0 right-0 z-20" @click.self="showActionSheet = false">
            <div class="mx-3 mb-2">
              <div class="overflow-hidden rounded-xl bg-white dark:bg-gray-800">
                <button class="w-full px-4 py-3.5 text-center text-base font-medium text-gray-900 active:bg-gray-100 dark:text-gray-100 dark:active:bg-gray-700" @click="handleSave">
                  💾 保存{{ currentItem?.type === 'video' ? '视频' : '图片' }}
                </button>
              </div>
            </div>
            <div class="mx-3 mb-6">
              <button class="w-full rounded-xl bg-white px-4 py-3.5 text-center text-base font-medium text-gray-900 active:bg-gray-100 dark:bg-gray-800 dark:text-gray-100 dark:active:bg-gray-700" @click="showActionSheet = false">
                取消
              </button>
            </div>
          </div>
        </Transition>

        <!-- 长按操作栏遮罩 -->
        <Transition name="fade">
          <div v-if="showActionSheet" class="absolute inset-0 z-[19] bg-black/40" @click="showActionSheet = false" />
        </Transition>
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

// 图片缩放 & 移动
const scale = ref(1)
const translateX = ref(0)
const translateY = ref(0)

// 视频控制
const videoPaused = ref(true)
const videoCurrentTime = ref(0)
const videoDuration = ref(0)
const showVideoControls = ref(true)
let controlsTimer: ReturnType<typeof setTimeout> | null = null

// 操作栏
const showActionSheet = ref(false)

// 触摸状态
let touchStartTime = 0
let touchStartX = 0
let touchStartY = 0
let lastTapTime = 0
let longPressTimer: ReturnType<typeof setTimeout> | null = null
let isPinching = false
let initialPinchDistance = 0
let initialScale = 1
let isSwiping = false
let isDragging = false

const currentItem = computed(() => props.items[currentIndex.value])

const imageTransformStyle = computed(() => ({
  transform: `scale(${scale.value}) translate(${translateX.value}px, ${translateY.value}px)`,
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
  showActionSheet.value = false
  emit('close')
}

// ====== 触摸事件 ======
function onTouchStart(e: TouchEvent) {
  if (showActionSheet.value) return
  touchStartTime = Date.now()
  clearLongPress()

  if (e.touches.length === 2) {
    // 双指缩放开始
    isPinching = true
    initialPinchDistance = getPinchDistance(e.touches)
    initialScale = scale.value
    return
  }

  if (e.touches.length === 1) {
    touchStartX = e.touches[0].clientX
    touchStartY = e.touches[0].clientY
    isSwiping = false
    isDragging = false

    // 长按计时
    longPressTimer = setTimeout(() => {
      showActionSheet.value = true
    }, 600)
  }
}

function onTouchMove(e: TouchEvent) {
  if (showActionSheet.value) return
  clearLongPress()

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
      // 缩放状态下拖动
      isDragging = true
      translateX.value += dx / scale.value
      translateY.value += dy / scale.value
      touchStartX = e.touches[0].clientX
      touchStartY = e.touches[0].clientY
    } else if (Math.abs(dx) > 10) {
      isSwiping = true
    }
  }
}

function onTouchEnd(e: TouchEvent) {
  if (showActionSheet.value) return
  clearLongPress()
  isPinching = false

  const elapsed = Date.now() - touchStartTime
  if (elapsed < 300 && !isSwiping && !isDragging) {
    const now = Date.now()
    if (now - lastTapTime < 300) {
      // 双击
      handleDoubleTap(e)
      lastTapTime = 0
    } else {
      lastTapTime = now
      // 延迟判断是否为单击
      setTimeout(() => {
        if (lastTapTime !== 0 && Date.now() - lastTapTime >= 280) {
          if (currentItem.value?.type === 'image') {
            close()
          }
          lastTapTime = 0
        }
      }, 300)
    }
    return
  }

  // 滑动切换
  if (isSwiping && scale.value <= 1 && e.changedTouches.length > 0) {
    const dx = e.changedTouches[0].clientX - touchStartX
    if (Math.abs(dx) > 60) {
      if (dx < 0 && currentIndex.value < props.items.length - 1) {
        currentIndex.value++
      } else if (dx > 0 && currentIndex.value > 0) {
        currentIndex.value--
      }
    }
  }

  isSwiping = false
  isDragging = false
}

function handleSingleTap() {
  // PC 端单击关闭（移动端由 touch 事件处理）
  if (currentItem.value?.type === 'image' && scale.value <= 1) {
    close()
  }
}

function handleDoubleTap(_e: TouchEvent) {
  if (currentItem.value?.type !== 'image') return
  if (scale.value > 1) {
    scale.value = 1
    translateX.value = 0
    translateY.value = 0
  } else {
    scale.value = 2.5
  }
}

function clearLongPress() {
  if (longPressTimer) {
    clearTimeout(longPressTimer)
    longPressTimer = null
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

// ====== 保存 ======
function handleSave() {
  if (!currentItem.value) return
  const a = document.createElement('a')
  a.href = currentItem.value.url
  a.download = currentItem.value.name
  a.target = '_blank'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  showActionSheet.value = false
}

onBeforeUnmount(() => {
  clearLongPress()
  if (controlsTimer) clearTimeout(controlsTimer)
  document.body.style.overflow = ''
})
</script>

<style scoped>
/* 进入/离开过渡 */
.media-viewer-fade-enter-active,
.media-viewer-fade-leave-active { transition: opacity 0.25s ease; }
.media-viewer-fade-enter-from,
.media-viewer-fade-leave-to { opacity: 0; }

.fade-enter-active,
.fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from,
.fade-leave-to { opacity: 0; }

/* 操作栏从底部滑入 */
.action-sheet-enter-active { transition: transform 0.3s ease; }
.action-sheet-leave-active { transition: transform 0.2s ease; }
.action-sheet-enter-from,
.action-sheet-leave-to { transform: translateY(100%); }

/* 视频进度条样式 */
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
