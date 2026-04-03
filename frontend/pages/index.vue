<template>
  <div class="landing-page">
    <!-- 全屏 Canvas 背景 -->
    <canvas ref="canvasRef" class="bg-canvas" />

    <!-- 进入博客按钮 -->
    <div
      ref="btnRef"
      class="enter-btn"
      :class="{ dragging: isDragging }"
      :style="btnPositionStyle"
      @mousedown.prevent="onPointerDown"
      @touchstart.prevent="onPointerDown"
    >
      <span class="enter-text">进入博客</span>
      <svg class="enter-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
      </svg>
    </div>

    <!-- 底部署名 -->
    <div class="landing-footer">fatwill.cloud</div>
  </div>
</template>

<script setup lang="ts">
import { prepareWithSegments, layoutNextLine } from '@chenglou/pretext'

useHead({
  title: 'fatwill.cloud',
  meta: [{ name: 'description', content: '全栈开发者 fatwillzeng 的个人博客' }],
})

// ====== 代码片段 ======
const codeSnippets = [
  '<div class="container mx-auto px-4">',
  'const handleClick = (e: MouseEvent) => { e.preventDefault() }',
  '.flex { display: flex; align-items: center; gap: 16px; }',
  'import { ref, computed, onMounted } from "vue"',
  'function fibonacci(n: number): number { return n <= 1 ? n : fib(n-1) + fib(n-2) }',
  'background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);',
  'export default defineNuxtConfig({ ssr: true, modules: ["@pinia/nuxt"] })',
  'const [state, dispatch] = useReducer(reducer, initialState)',
  '@keyframes fadeIn { from { opacity: 0; transform: translateY(10px) } to { opacity: 1 } }',
  'SELECT id, title, created_at FROM articles ORDER BY created_at DESC LIMIT 10',
  '<template><NuxtPage /><ClientOnly><AppLoading /></ClientOnly></template>',
  'border-radius: 12px; box-shadow: 0 8px 32px rgba(0,0,0,.12); backdrop-filter: blur(10px);',
  'async function fetchData<T>(url: string): Promise<T> { const res = await fetch(url) }',
  'git commit -m "feat: add pretext text layout engine integration"',
  'docker run -d -p 3000:3000 --name blog-frontend blog:latest',
  'npm install @chenglou/pretext @pinia/nuxt @nuxtjs/tailwindcss',
  'interface Article { id: string; title: string; content: string; createdAt: Date }',
  'transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); will-change: transform;',
  'const router = useRouter(); await router.push({ path: "/home", query: { tab: "articles" } })',
  'CREATE TABLE articles (id INT PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255) NOT NULL);',
  'pm2 start ecosystem.config.cjs --env production && pm2 save',
  'const { data, pending } = await useAsyncData("articles", () => $fetch("/api/articles"))',
  'nginx -t && nginx -s reload  # test config and graceful reload',
  'type Props = { title: string; onClick?: () => void; className?: string }',
]

const lineColors = [
  '#7dd3fc', '#4ade80', '#fb923c', '#a78bfa',
  '#f472b6', '#22d3ee', '#facc15', '#34d399',
  '#7dd3fc', '#4ade80', '#fb923c', '#a78bfa',
  '#f472b6', '#22d3ee', '#facc15', '#34d399',
  '#7dd3fc', '#4ade80',
]

// ====== 行数据结构 ======
interface CodeLineData {
  text: string
  repeatText: string // 重复拼接后的文字
  font: string
  color: string
  opacity: number
  speed: number // px/frame，正=向右，负=向左
  y: number
  offset: number
  measuredWidth: number // pretext 测量的精确宽度
  gap: number // 重复间距
}

// ====== Refs ======
const canvasRef = ref<HTMLCanvasElement | null>(null)
const btnRef = ref<HTMLElement | null>(null)
const isDragging = ref(false)

// 按钮位置（用 left/top px 值）
const btnLeft = ref(0)
const btnTop = ref(0)
const btnInited = ref(false)

const btnPositionStyle = computed(() => {
  if (!btnInited.value) {
    // 初始居中用 CSS
    return { left: '50%', top: '50%', transform: 'translate(-50%, -50%)' }
  }
  return {
    left: `${btnLeft.value}px`,
    top: `${btnTop.value}px`,
    transform: 'none',
  }
})

// ====== 动画数据 ======
let lines: CodeLineData[] = []
let animId = 0
let ctx: CanvasRenderingContext2D | null = null
let canvasW = 0
let canvasH = 0
const LINE_COUNT = 18
const LINE_HEIGHT = 36 // 行高（含间距）
const BTN_PAD = 28 // 按钮周围文字留白

// ====== 初始化行数据 ======
function initLines() {
  lines = []
  const startY = (canvasH - LINE_COUNT * LINE_HEIGHT) / 2 + 20
  for (let i = 0; i < LINE_COUNT; i++) {
    const fontSize = 13 + Math.random() * 3
    const speed = (0.4 + Math.random() * 0.8) * (i % 2 === 0 ? 1 : -1)
    const snippet = codeSnippets[i % codeSnippets.length]
    const gap = 80
    // 重复文字 5 次，确保填满屏幕
    const repeatText = Array(5).fill(snippet).join('    ')

    lines.push({
      text: snippet,
      repeatText,
      font: `${fontSize}px "Fira Code", "Cascadia Code", "JetBrains Mono", Consolas, monospace`,
      color: lineColors[i % lineColors.length],
      opacity: 0.15 + Math.random() * 0.2,
      speed,
      y: startY + i * LINE_HEIGHT,
      offset: Math.random() * 1000, // 随机初始偏移，避免整齐排列
      measuredWidth: 0,
      gap,
    })
  }
}

// ====== 用 pretext 精确测量每行宽度 ======
function measureLines() {
  for (const line of lines) {
    try {
      const prepared = prepareWithSegments(line.repeatText, line.font)
      const result = layoutNextLine(prepared, { segmentIndex: 0, graphemeIndex: 0 }, 99999)
      if (result) {
        line.measuredWidth = result.width + line.gap
      }
    } catch {
      // fallback: 用 canvas measureText
      if (ctx) {
        ctx.font = line.font
        line.measuredWidth = ctx.measureText(line.repeatText).width + line.gap
      }
    }
  }
}

// ====== Canvas 绘制 ======
function drawFrame() {
  if (!ctx || !canvasRef.value) return

  ctx.clearRect(0, 0, canvasW, canvasH)

  // 获取按钮位置
  const btnEl = btnRef.value
  let btnRect = { left: -9999, top: -9999, right: -9999, bottom: -9999, width: 0, height: 0 }
  if (btnEl) {
    btnRect = btnEl.getBoundingClientRect()
  }

  // 设置 clip 路径排除按钮区域（evenodd 规则）
  ctx.save()
  const region = new Path2D()
  region.rect(0, 0, canvasW, canvasH)
  region.rect(
    btnRect.left - BTN_PAD,
    btnRect.top - BTN_PAD,
    btnRect.width + BTN_PAD * 2,
    btnRect.height + BTN_PAD * 2,
  )
  ctx.clip(region, 'evenodd')

  for (const line of lines) {
    line.offset += line.speed

    ctx.font = line.font
    ctx.globalAlpha = line.opacity
    ctx.fillStyle = line.color

    const mw = line.measuredWidth
    if (mw <= 0) continue

    // 计算起始 x，保证无缝循环
    let startX = line.offset % mw
    if (line.speed > 0) {
      // 向右滚动
      if (startX > 0) startX -= mw
    } else {
      // 向左滚动
      if (startX < 0) startX += mw
      startX -= mw
    }

    // 绘制足够多次以填满屏幕
    let x = startX
    while (x < canvasW + mw) {
      ctx.fillText(line.repeatText, x, line.y)
      x += mw
    }
  }

  ctx.restore()
  ctx.globalAlpha = 1

  animId = requestAnimationFrame(drawFrame)
}

// ====== 按钮拖拽 ======
let pointerDownTimer: ReturnType<typeof setTimeout> | null = null
let isLongPress = false
let hasMoved = false
let startPointerX = 0
let startPointerY = 0

function onPointerDown(e: MouseEvent | TouchEvent) {
  isLongPress = false
  hasMoved = false
  const pos = getPointerPos(e)
  startPointerX = pos.x
  startPointerY = pos.y

  const el = btnRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()

  // 如果还没初始化 px 位置，先设一下
  if (!btnInited.value) {
    btnLeft.value = rect.left
    btnTop.value = rect.top
    btnInited.value = true
  }

  const offsetX = pos.x - rect.left
  const offsetY = pos.y - rect.top

  pointerDownTimer = setTimeout(() => {
    isLongPress = true
    isDragging.value = true

    const onMove = (ev: MouseEvent | TouchEvent) => {
      ev.preventDefault()
      hasMoved = true
      const p = getPointerPos(ev)
      btnLeft.value = p.x - offsetX
      btnTop.value = p.y - offsetY
    }

    const onUp = () => {
      isDragging.value = false
      document.removeEventListener('mousemove', onMove)
      document.removeEventListener('mouseup', onUp)
      document.removeEventListener('touchmove', onMove)
      document.removeEventListener('touchend', onUp)
    }

    document.addEventListener('mousemove', onMove, { passive: false })
    document.addEventListener('mouseup', onUp)
    document.addEventListener('touchmove', onMove, { passive: false })
    document.addEventListener('touchend', onUp)
  }, 300)

  const onEarlyUp = () => {
    if (pointerDownTimer) {
      clearTimeout(pointerDownTimer)
      pointerDownTimer = null
    }
    // 短按 + 没移动 → 跳转
    if (!isLongPress && !hasMoved) {
      const dx = Math.abs(getPointerPos(e).x - startPointerX)
      const dy = Math.abs(getPointerPos(e).y - startPointerY)
      if (dx < 5 && dy < 5) {
        navigateTo('/home')
      }
    }
    document.removeEventListener('mouseup', onEarlyUp)
    document.removeEventListener('touchend', onEarlyUp)
  }

  document.addEventListener('mouseup', onEarlyUp)
  document.addEventListener('touchend', onEarlyUp)
}

function getPointerPos(e: MouseEvent | TouchEvent): { x: number; y: number } {
  if ('touches' in e && e.touches.length > 0) {
    return { x: e.touches[0].clientX, y: e.touches[0].clientY }
  }
  if ('changedTouches' in e && e.changedTouches.length > 0) {
    return { x: e.changedTouches[0].clientX, y: e.changedTouches[0].clientY }
  }
  const me = e as MouseEvent
  return { x: me.clientX, y: me.clientY }
}

// ====== resize ======
function handleResize() {
  const canvas = canvasRef.value
  if (!canvas) return
  const dpr = window.devicePixelRatio || 1
  canvasW = window.innerWidth
  canvasH = window.innerHeight
  canvas.width = canvasW * dpr
  canvas.height = canvasH * dpr
  canvas.style.width = `${canvasW}px`
  canvas.style.height = `${canvasH}px`
  ctx = canvas.getContext('2d')
  if (ctx) {
    ctx.scale(dpr, dpr)
  }
  // 重新计算行 Y 坐标
  const startY = (canvasH - LINE_COUNT * LINE_HEIGHT) / 2 + 20
  for (let i = 0; i < lines.length; i++) {
    lines[i].y = startY + i * LINE_HEIGHT
  }
}

// ====== 生命周期 ======
onMounted(() => {
  initLines()
  handleResize()
  measureLines()
  animId = requestAnimationFrame(drawFrame)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  cancelAnimationFrame(animId)
  window.removeEventListener('resize', handleResize)
  if (pointerDownTimer) clearTimeout(pointerDownTimer)
})
</script>

<style scoped>
.landing-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: #0d1117;
  user-select: none;
}

.bg-canvas {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
}

/* ====== 进入按钮 ====== */
.enter-btn {
  position: fixed;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 40px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 60px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 2px;
  cursor: grab;
  user-select: none;
  box-shadow: 0 0 40px rgba(99, 102, 241, 0.15), 0 8px 32px rgba(0, 0, 0, 0.3);
  transition: box-shadow 0.3s ease, border-color 0.3s ease, background 0.3s ease;
  z-index: 10;
  white-space: nowrap;
}

.enter-btn:hover {
  background: rgba(255, 255, 255, 0.14);
  border-color: rgba(255, 255, 255, 0.35);
  box-shadow: 0 0 60px rgba(99, 102, 241, 0.25), 0 8px 32px rgba(0, 0, 0, 0.3);
}

.enter-btn.dragging {
  cursor: grabbing;
  box-shadow: 0 0 80px rgba(99, 102, 241, 0.35), 0 12px 48px rgba(0, 0, 0, 0.4);
  border-color: rgba(255, 255, 255, 0.4);
}

.enter-text {
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.enter-arrow {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease;
}

.enter-btn:hover .enter-arrow {
  transform: translateX(4px);
}

/* ====== 底部署名 ====== */
.landing-footer {
  position: fixed;
  bottom: 24px;
  left: 0;
  right: 0;
  text-align: center;
  color: rgba(255, 255, 255, 0.15);
  font-size: 13px;
  font-family: 'Fira Code', monospace;
  letter-spacing: 2px;
  z-index: 20;
  pointer-events: none;
}

/* ====== 移动端 ====== */
@media (max-width: 768px) {
  .enter-btn {
    padding: 14px 32px;
    font-size: 15px;
    letter-spacing: 1px;
  }
}
</style>
