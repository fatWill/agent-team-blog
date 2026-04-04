<template>
  <div class="landing-page">
    <canvas ref="canvasRef" class="bg-canvas" />
    <div
      ref="btnRef"
      class="enter-btn"
      :class="{ dragging: isDragging }"
      :style="{ left: btnLeft + 'px', top: btnTop + 'px', opacity: btnVisible ? '' : '0', pointerEvents: btnVisible ? '' : 'none' }"
      @mousedown.prevent="onPointerDown"
      @touchstart.prevent="onPointerDown"
    >
      <span class="enter-text">进入博客</span>
      <svg class="enter-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import {
  prepareWithSegments,
  layoutNextLine,
  type PreparedTextWithSegments,
  type LayoutCursor,
} from '@chenglou/pretext'

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
  'async function fetchData<T>(url: string): Promise<T> { const res = await fetch(url); return res.json() }',
  'git commit -m "feat: integrate @chenglou/pretext for text layout engine"',
  'docker-compose up -d --build && docker ps -a | grep blog',
  'npm install @chenglou/pretext @pinia/nuxt @nuxtjs/tailwindcss --save',
  'interface Article { id: string; title: string; content: string; tags: string[]; createdAt: Date }',
  'transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); will-change: transform, opacity;',
  'const { data, pending } = await useAsyncData("articles", () => $fetch("/api/articles"))',
  'CREATE TABLE articles (id INT PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255) NOT NULL);',
  'pm2 start ecosystem.config.cjs --env production && pm2 save && pm2 logs --lines 20',
  'type Props = { title: string; description?: string; onClick?: () => void; className?: string }',
  'nginx -t && systemctl reload nginx  # test config and graceful reload',
  'const router = useRouter(); watch(() => route.path, () => { scrollTo(0, 0) })',
]

const colors = [
  '#7dd3fc', '#4ade80', '#fb923c', '#a78bfa',
  '#f472b6', '#22d3ee', '#facc15', '#34d399',
  '#c084fc', '#f87171', '#38bdf8', '#86efac',
]

// ====== 每行数据结构 ======
interface TextLine {
  text: string          // 完整文字（重复多次用于无缝循环）
  y: number             // 行 baseline 的 y 坐标
  offsetX: number       // 当前水平偏移
  speed: number         // 每帧移动的像素（正=向右，负=向左）
  color: string
  font: string
  opacity: number
  prepared: PreparedTextWithSegments | null
  measuredWidth: number // 整段文字宽度（用于循环）
}

// ====== Refs ======
const canvasRef = ref<HTMLCanvasElement | null>(null)
const btnRef = ref<HTMLElement | null>(null)
const isDragging = ref(false)
const btnLeft = ref(-9999)
const btnTop = ref(-9999)
const btnVisible = ref(false)

// ====== 动画数据 ======
let lines: TextLine[] = []
let animId = 0
let ctx: CanvasRenderingContext2D | null = null
let canvasW = 0
let canvasH = 0
const BTN_GAP = 16 // 按钮与文字之间的间距
const LINE_START: LayoutCursor = { segmentIndex: 0, graphemeIndex: 0 }
const DRAG_THRESHOLD = 8 // 拖拽触发阈值（px）
let lineCount = 18 // 动态计算的行数

// 按钮位置缓存（避免在 rAF 中调用 getBoundingClientRect 触发回流）
let cachedBtnL = -9999
let cachedBtnT = -9999
let cachedBtnR = -9999
let cachedBtnB = -9999

function updateBtnCache() {
  if (!btnRef.value) return
  const rect = btnRef.value.getBoundingClientRect()
  cachedBtnL = rect.left
  cachedBtnT = rect.top
  cachedBtnR = rect.right
  cachedBtnB = rect.bottom
}

// 估算按钮尺寸（避免依赖 getBoundingClientRect，初始居中用）
function estimateBtnSize(): { w: number; h: number } {
  if (typeof window !== 'undefined' && window.innerWidth < 768) {
    // 移动端: padding 14px 32px, font-size 15px, gap 10px, arrow 20px
    // 宽度 ≈ 32 + text(~75px) + 10 + 20 + 32 ≈ 169, 高度 ≈ 14 + 15*1.2 + 14 ≈ 46
    return { w: 169, h: 46 }
  }
  // 桌面端: padding 16px 44px, font-size 17px, gap 10px, arrow 20px
  // 宽度 ≈ 44 + text(~85px) + 10 + 20 + 44 ≈ 203, 高度 ≈ 16 + 17*1.2 + 16 ≈ 52
  return { w: 203, h: 52 }
}

// 根据屏幕宽度判断是否为移动端，动态调整参数
function isMobile(): boolean {
  return canvasW < 768
}
function getMinLineGap(): number {
  return isMobile() ? 36 : 26
}
function getFontSize(index: number): number {
  // 移动端字号 10-11px，桌面端 13-16px
  return isMobile() ? 10 + (index % 2) : 13 + (index % 4)
}

// ====== 构建每行文字（随机 2-3 个片段拼接，重复多次）======
function buildLineText(index: number): string {
  const count = 2 + (index % 2) // 2 或 3 个片段
  const parts: string[] = []
  for (let i = 0; i < count; i++) {
    parts.push(snippets[(index * 3 + i * 7) % snippets.length]!)
  }
  const single = parts.join('  //  ')
  return Array(4).fill(single).join('      ')
}

// ====== 初始化行数据 ======
function initLines() {
  lines = []
  const minGap = getMinLineGap()
  // 动态计算行数：确保行间距 >= minGap
  lineCount = Math.max(6, Math.floor(canvasH / minGap))
  const lineHeight = canvasH / lineCount

  for (let i = 0; i < lineCount; i++) {
    const fontSize = getFontSize(i)
    // 速度 0.4~1.2，奇偶行方向相反
    const speed = (0.4 + ((i * 37) % 100) / 100 * 0.8) * (i % 2 === 0 ? 1 : -1)
    const text = buildLineText(i)

    lines.push({
      text,
      y: lineHeight * (i + 0.5),
      offsetX: (i * 137) % 500, // 随机初始偏移
      speed,
      color: colors[i % colors.length]!,
      font: `${fontSize}px "Fira Code", "Cascadia Code", "JetBrains Mono", Consolas, monospace`,
      opacity: 0.15 + (i % 5) * 0.05, // 0.15 ~ 0.35
      prepared: null,
      measuredWidth: 0,
    })
  }
}

// ====== 用 pretext 预处理（绕排用）+ canvas measureText 测量真实宽度 ======
function measureLines() {
  if (!ctx) return
  for (const line of lines) {
    // 用 pretext 预处理（绕排时需要）
    try {
      line.prepared = prepareWithSegments(line.text, line.font)
    } catch {
      line.prepared = null
    }

    // 用 canvas measureText 测量实际像素宽度（比 layoutNextLine 更可靠）
    ctx.font = line.font
    line.measuredWidth = ctx.measureText(line.text).width + 80
  }
}

/**
 * 绘制一行横向滚动文字（无限循环）
 *
 * 核心绕排逻辑：
 * - 对于每个重复绘制的文字副本，检查其每个字符是否落在按钮区域
 * - 如果行与按钮在 y 方向不重叠 → 正常绘制
 * - 如果行与按钮在 y 方向重叠 → 用 layoutNextLine 动态宽度把文字截断到按钮左侧，
 *   然后从按钮右侧继续绘制剩余文字
 */
function drawLine(
  line: TextLine,
  btnL: number,
  btnT: number,
  btnR: number,
  btnB: number,
) {
  if (!ctx) return
  const w = line.measuredWidth
  if (w <= 0) return

  ctx.save()
  ctx.globalAlpha = line.opacity
  ctx.font = line.font
  ctx.fillStyle = line.color

  // 判断这行是否与按钮在 y 方向重叠
  const fontSize = parseInt(line.font)
  const lineTop = line.y - fontSize
  const lineBottom = line.y + 4
  const overlapsBtn = lineBottom >= btnT - BTN_GAP && lineTop <= btnB + BTN_GAP

  // 计算循环起始 x
  let x = ((line.offsetX % w) + w) % w - w

  if (!overlapsBtn) {
    // ===== 不与按钮重叠：正常水平循环绘制 =====
    while (x < canvasW) {
      ctx.fillText(line.text, x, line.y)
      x += w
    }
  } else {
    // ===== 与按钮重叠：用 pretext layoutNextLine 动态宽度绕排 =====
    // 每个重复副本都需要检查是否与按钮水平重叠
    while (x < canvasW) {
      const textRight = x + w

      // 这个副本完全在按钮左侧或右侧，不需要绕排
      if (textRight <= btnL - BTN_GAP || x >= btnR + BTN_GAP) {
        ctx.fillText(line.text, x, line.y)
        x += w
        continue
      }

      // 这个副本与按钮水平重叠，需要绕排
      if (line.prepared) {
        // 左侧文字：从文字起点 x 到按钮左边界
        const leftAvail = btnL - BTN_GAP - x
        if (leftAvail > 10) {
          const leftLine = layoutNextLine(line.prepared, LINE_START, leftAvail)
          if (leftLine && leftLine.text) {
            ctx.fillText(leftLine.text, x, line.y)
          }
        }

        // 右侧文字：从按钮右边界到文字副本结束
        const rightStart = btnR + BTN_GAP
        if (rightStart < textRight && rightStart < canvasW) {
          // 计算右侧文字应该从原文的哪个位置开始
          // 方法：用 layoutNextLine 先消耗掉「从文字起点到按钮右边界」的宽度
          const consumeWidth = rightStart - x
          if (consumeWidth > 0) {
            const consumed = layoutNextLine(line.prepared, LINE_START, consumeWidth)
            if (consumed) {
              // 从消耗后的 cursor 继续排版剩余文字
              const rightAvail = Math.min(textRight, canvasW) - rightStart
              if (rightAvail > 10) {
                const rightLine = layoutNextLine(line.prepared, consumed.end, rightAvail)
                if (rightLine && rightLine.text) {
                  ctx.fillText(rightLine.text, rightStart, line.y)
                }
              }
            }
          }
        }
      } else {
        // fallback：没有 prepared 数据，用 clip 裁切
        const leftEdge = btnL - BTN_GAP
        const rightEdge = btnR + BTN_GAP

        if (leftEdge > 0) {
          ctx.save()
          ctx.beginPath()
          ctx.rect(0, lineTop - 4, leftEdge, lineBottom - lineTop + 8)
          ctx.clip()
          ctx.fillText(line.text, x, line.y)
          ctx.restore()
        }
        if (rightEdge < canvasW) {
          ctx.save()
          ctx.beginPath()
          ctx.rect(rightEdge, lineTop - 4, canvasW - rightEdge, lineBottom - lineTop + 8)
          ctx.clip()
          ctx.fillText(line.text, x, line.y)
          ctx.restore()
        }
      }

      x += w
    }
  }

  ctx.restore()
}

// ====== 每帧绘制 ======
function drawFrame() {
  if (!ctx || !canvasRef.value) {
    animId = requestAnimationFrame(drawFrame)
    return
  }

  // 每帧重置变换矩阵，避免 scale 累积
  const dpr = window.devicePixelRatio || 1
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, canvasW, canvasH)

  for (const line of lines) {
    line.offsetX += line.speed
    drawLine(line, cachedBtnL, cachedBtnT, cachedBtnR, cachedBtnB)
  }

  animId = requestAnimationFrame(drawFrame)
}

// ====== 按钮拖拽 ======
let hasDragged = false
let startPointerX = 0
let startPointerY = 0
let dragOffsetX = 0
let dragOffsetY = 0

function getPointerPos(e: MouseEvent | TouchEvent): { x: number; y: number } {
  if ('touches' in e && e.touches.length > 0) {
    return { x: e.touches[0]!.clientX, y: e.touches[0]!.clientY }
  }
  if ('changedTouches' in e && e.changedTouches.length > 0) {
    return { x: e.changedTouches[0]!.clientX, y: e.changedTouches[0]!.clientY }
  }
  return { x: (e as MouseEvent).clientX, y: (e as MouseEvent).clientY }
}

function onPointerDown(e: MouseEvent | TouchEvent) {
  hasDragged = false
  const pos = getPointerPos(e)
  startPointerX = pos.x
  startPointerY = pos.y

  const el = btnRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  dragOffsetX = pos.x - rect.left
  dragOffsetY = pos.y - rect.top

  const onMove = (ev: MouseEvent | TouchEvent) => {
    const p = getPointerPos(ev)
    const dx = p.x - startPointerX
    const dy = p.y - startPointerY
    const dist = Math.sqrt(dx * dx + dy * dy)

    if (!hasDragged && dist < DRAG_THRESHOLD) return

    ev.preventDefault()
    hasDragged = true
    isDragging.value = true
    btnLeft.value = p.x - dragOffsetX
    btnTop.value = p.y - dragOffsetY
    // 拖拽时同步更新按钮位置缓存
    updateBtnCache()
  }

  const onUp = () => {
    if (!hasDragged) {
      // 没有拖动 → 视为点击跳转
      navigateTo('/home')
    }
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
  // 不在这里 scale！drawFrame 每帧用 setTransform 重置

  // 完全重建行数据（行数/字号可能因屏幕尺寸变化而不同）
  initLines()
  measureLines()

  // 重新居中按钮并更新缓存
  if (btnRef.value) {
    const { w, h } = estimateBtnSize()
    btnLeft.value = (canvasW - w) / 2
    btnTop.value = (canvasH - h) / 2
  }
  nextTick(() => updateBtnCache())
}

// ====== 生命周期 ======
onMounted(() => {
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
  // 不在这里 scale！drawFrame 每帧用 setTransform 重置

  // 先用估算值设置初始居中位置（此时按钮还是 opacity:0 不可见）
  const { w, h } = estimateBtnSize()
  btnLeft.value = (canvasW - w) / 2
  btnTop.value = (canvasH - h) / 2

  initLines()
  measureLines()

  // 初始化按钮位置（居中）+ 缓存 + 启动动画
  nextTick(() => {
    updateBtnCache()
    animId = requestAnimationFrame(drawFrame)
    // 下一帧再显示，避免左上角闪烁
    requestAnimationFrame(() => {
      btnVisible.value = true
    })
    // 动画结束后（0.6s）用真实尺寸更新缓存
    setTimeout(() => {
      updateBtnCache()
    }, 700)
  })

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  cancelAnimationFrame(animId)
  window.removeEventListener('resize', handleResize)
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
  animation: btnEntrance 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) both;
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

@keyframes btnEntrance {
  from {
    opacity: 0;
    transform: scale(0.3);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@media (max-width: 768px) {
  .enter-btn {
    padding: 14px 32px;
    font-size: 15px;
    letter-spacing: 1px;
  }
}
</style>