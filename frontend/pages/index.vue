<template>
  <div class="landing-page">
    <!-- 全屏 Canvas 背景 -->
    <canvas ref="canvasRef" class="bg-canvas" />

    <!-- 进入博客按钮 -->
    <div
      ref="btnRef"
      class="enter-btn"
      :class="{ dragging: isDragging }"
      :style="{ left: btnLeft + 'px', top: btnTop + 'px' }"
      @mousedown.prevent="onPointerDown"
      @touchstart.prevent="onPointerDown"
    >
      <span class="enter-text">进入博客</span>
      <svg class="enter-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
      </svg>
    </div>

    <!-- 底部署名 -->
    <div class="footer">fatwill.cloud</div>
  </div>
</template>

<script setup lang="ts">
import { prepareWithSegments, layoutNextLine } from '@chenglou/pretext'
import type { PreparedTextWithSegments } from '@chenglou/pretext'

useHead({
  title: 'fatwill.cloud',
  meta: [{ name: 'description', content: '全栈开发者 fatwillzeng 的个人博客' }],
})

// ====== 代码片段 ======
const snippets = [
  '<div class="container mx-auto px-4 py-8">',
  'const handleClick = (e: MouseEvent) => { e.preventDefault(); router.push("/home") }',
  '.flex { display: flex; align-items: center; gap: 16px; justify-content: space-between; }',
  'import { ref, computed, onMounted, watch } from "vue"',
  'function fibonacci(n: number): number { return n <= 1 ? n : fib(n-1) + fib(n-2) }',
  'background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px;',
  'export default defineNuxtConfig({ ssr: true, modules: ["@pinia/nuxt", "@nuxtjs/tailwindcss"] })',
  'const [state, dispatch] = useReducer(reducer, initialState)',
  '@keyframes fadeIn { from { opacity: 0; transform: translateY(10px) } to { opacity: 1; transform: none } }',
  'SELECT id, title, created_at FROM articles WHERE published = 1 ORDER BY created_at DESC LIMIT 10',
  '<template><NuxtPage /><ClientOnly><AppLoading /></ClientOnly></template>',
  'border-radius: 12px; box-shadow: 0 8px 32px rgba(0,0,0,.12); backdrop-filter: blur(10px);',
  'async function fetchData<T>(url: string, options?: RequestInit): Promise<T> { const res = await fetch(url, options) }',
  'git commit -m "feat: integrate @chenglou/pretext for text layout engine"',
  'docker-compose up -d --build && docker ps -a | grep blog',
  'npm install @chenglou/pretext @pinia/nuxt @nuxtjs/tailwindcss --save',
  'interface Article { id: string; title: string; content: string; tags: string[]; createdAt: Date }',
  'transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); will-change: transform, opacity;',
  'const { data, pending, error } = await useAsyncData("articles", () => $fetch("/api/articles"))',
  'CREATE TABLE articles (id INT PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255) NOT NULL, content LONGTEXT);',
  'pm2 start ecosystem.config.cjs --env production && pm2 save && pm2 logs --lines 20',
  'type Props = { title: string; description?: string; onClick?: () => void; className?: string }',
  'nginx -t && systemctl reload nginx  # test config and graceful reload',
  'const router = useRouter(); const route = useRoute(); watch(() => route.path, () => { scrollTo(0, 0) })',
]

const colors = [
  '#7dd3fc', '#4ade80', '#fb923c', '#a78bfa',
  '#f472b6', '#22d3ee', '#facc15', '#34d399',
  '#c084fc', '#f87171', '#38bdf8', '#86efac',
]

// ====== 行数据结构 ======
interface TextLine {
  text: string
  y: number
  offsetX: number
  speed: number
  color: string
  font: string
  opacity: number
  prepared: PreparedTextWithSegments | null
  measuredWidth: number
}

// ====== Refs ======
const canvasRef = ref<HTMLCanvasElement | null>(null)
const btnRef = ref<HTMLElement | null>(null)
const isDragging = ref(false)
const btnLeft = ref(0)
const btnTop = ref(0)

// ====== 动画数据 ======
let lines: TextLine[] = []
let animId = 0
let ctx: CanvasRenderingContext2D | null = null
let canvasW = 0
let canvasH = 0
const LINE_COUNT = 18
const BTN_PAD = 24

// ====== 构建每行文字（随机 2-3 个片段拼接，重复 4 次）======
function buildLineText(index: number): string {
  const count = 2 + (index % 2) // 2 或 3 个片段
  const parts: string[] = []
  for (let i = 0; i < count; i++) {
    parts.push(snippets[(index * 3 + i * 7) % snippets.length])
  }
  const single = parts.join('  //  ')
  return Array(4).fill(single).join('      ')
}

// ====== 初始化行数据 ======
function initLines() {
  lines = []
  const lineHeight = canvasH / LINE_COUNT
  for (let i = 0; i < LINE_COUNT; i++) {
    const fontSize = 13 + (i % 4) // 13-16px
    const speed = (0.4 + ((i * 37) % 100) / 100 * 0.8) * (i % 2 === 0 ? 1 : -1)
    const text = buildLineText(i)

    lines.push({
      text,
      y: lineHeight * (i + 0.5) + fontSize * 0.35, // 垂直居中于行带
      offsetX: (i * 137) % 500, // 随机初始偏移
      speed,
      color: colors[i % colors.length],
      font: `${fontSize}px "Fira Code", "Cascadia Code", "JetBrains Mono", Consolas, monospace`,
      opacity: 0.15 + (i % 5) * 0.05, // 0.15 - 0.35
      prepared: null,
      measuredWidth: 0,
    })
  }
}

// ====== 用 pretext 预处理并测量每行宽度 ======
function measureLines() {
  for (const line of lines) {
    try {
      const prepared = prepareWithSegments(line.text, line.font)
      line.prepared = prepared
      const result = layoutNextLine(prepared, { segmentIndex: 0, graphemeIndex: 0 }, 99999)
      if (result) {
        line.measuredWidth = result.width + 60 // 加间距
      }
    } catch {
      // fallback: 用 canvas measureText
      line.prepared = null
      if (ctx) {
        ctx.font = line.font
        line.measuredWidth = ctx.measureText(line.text).width + 60
      }
    }
  }
}

// ====== 绘制滚动文字（无限循环）======
function drawScrollingText(line: TextLine) {
  if (!ctx) return
  const w = line.measuredWidth
  if (w <= 0) return

  // 计算循环起始 x
  let x = ((line.offsetX % w) + w) % w - w
  while (x < canvasW) {
    ctx.fillText(line.text, x, line.y)
    x += w
  }
}

// ====== 每帧绘制 ======
function drawFrame() {
  if (!ctx || !canvasRef.value) return

  ctx.clearRect(0, 0, canvasW, canvasH)

  // 获取按钮位置（每帧实时获取，支持拖拽后位置变化）
  const btnEl = btnRef.value
  let btnRect = { left: -9999, top: -9999, right: -9999, bottom: -9999, width: 0, height: 0 }
  if (btnEl) {
    btnRect = btnEl.getBoundingClientRect()
  }

  for (const line of lines) {
    line.offsetX += line.speed

    ctx.save()
    ctx.globalAlpha = line.opacity
    ctx.font = line.font
    ctx.fillStyle = line.color

    // 检查这行文字是否与按钮在 y 方向重叠
    // 字体高度大约等于 fontSize，baseline 在下方约 0.8 处
    const fontSize = parseInt(line.font)
    const lineTop = line.y - fontSize
    const lineBottom = line.y + 4

    const overlapsBtn =
      lineBottom >= btnRect.top - BTN_PAD &&
      lineTop <= btnRect.bottom + BTN_PAD

    if (!overlapsBtn) {
      // 不与按钮重叠：正常绘制整行
      drawScrollingText(line)
    } else {
      // 与按钮重叠：分左右两段绘制，跳过按钮区域
      const leftEdge = btnRect.left - BTN_PAD
      const rightEdge = btnRect.right + BTN_PAD

      // 左段
      if (leftEdge > 0) {
        ctx.save()
        ctx.beginPath()
        ctx.rect(0, 0, leftEdge, canvasH)
        ctx.clip()
        drawScrollingText(line)
        ctx.restore()
      }

      // 右段
      if (rightEdge < canvasW) {
        ctx.save()
        ctx.beginPath()
        ctx.rect(rightEdge, 0, canvasW - rightEdge, canvasH)
        ctx.clip()
        drawScrollingText(line)
        ctx.restore()
      }
    }

    ctx.restore()
  }

  animId = requestAnimationFrame(drawFrame)
}

// ====== 按钮拖拽 ======
let pointerDownTimer: ReturnType<typeof setTimeout> | null = null
let isLongPress = false
let hasMoved = false
let startPointerX = 0
let startPointerY = 0

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

function onPointerDown(e: MouseEvent | TouchEvent) {
  isLongPress = false
  hasMoved = false
  const pos = getPointerPos(e)
  startPointerX = pos.x
  startPointerY = pos.y

  const el = btnRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
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
    if (!isLongPress && !hasMoved) {
      navigateTo('/home')
    }
    document.removeEventListener('mouseup', onEarlyUp)
    document.removeEventListener('touchend', onEarlyUp)
  }

  document.addEventListener('mouseup', onEarlyUp)
  document.addEventListener('touchend', onEarlyUp)
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
  const lineHeight = canvasH / LINE_COUNT
  for (let i = 0; i < lines.length; i++) {
    const fontSize = parseInt(lines[i].font)
    lines[i].y = lineHeight * (i + 0.5) + fontSize * 0.35
  }

  // 重新居中按钮
  if (btnRef.value) {
    const rect = btnRef.value.getBoundingClientRect()
    btnLeft.value = (canvasW - rect.width) / 2
    btnTop.value = (canvasH - rect.height) / 2
  }
}

// ====== 生命周期 ======
onMounted(() => {
  // 先初始化 canvas 尺寸
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

  initLines()
  measureLines()

  // 初始化按钮位置（居中）
  nextTick(() => {
    if (btnRef.value) {
      const rect = btnRef.value.getBoundingClientRect()
      btnLeft.value = (canvasW - rect.width) / 2
      btnTop.value = (canvasH - rect.height) / 2
    }
  })

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
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 44px;
  background: rgba(255, 255, 255, 0.07);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 60px;
  color: rgba(255, 255, 255, 0.88);
  font-size: 17px;
  font-weight: 500;
  letter-spacing: 2px;
  cursor: grab;
  user-select: none;
  box-shadow: 0 0 40px rgba(99, 102, 241, 0.12), 0 8px 32px rgba(0, 0, 0, 0.25);
  transition: box-shadow 0.3s ease, border-color 0.3s ease, background 0.3s ease;
  white-space: nowrap;
}

.enter-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.3);
  box-shadow: 0 0 60px rgba(99, 102, 241, 0.25), 0 8px 32px rgba(0, 0, 0, 0.3);
}

.enter-btn.dragging {
  cursor: grabbing;
  transition: none;
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
.footer {
  position: fixed;
  bottom: 20px;
  left: 0;
  right: 0;
  text-align: center;
  color: rgba(255, 255, 255, 0.18);
  font-size: 12px;
  font-family: 'Fira Code', monospace;
  letter-spacing: 3px;
  pointer-events: none;
  z-index: 5;
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
