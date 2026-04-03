<template>
  <div class="index-page">
    <canvas ref="canvasRef" class="bg-canvas" />
    <button
      ref="btnRef"
      class="enter-btn"
      :class="{ dragging: isDragging }"
      :style="{ left: btnLeft + 'px', top: btnTop + 'px' }"
      @mousedown="onPointerDown"
      @touchstart.prevent="onPointerDown"
      @click="onBtnClick"
    >
      <span>进入博客</span>
      <span class="arrow">→</span>
    </button>
    <div class="footer">fatwill.cloud</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import {
  prepareWithSegments,
  layoutNextLine,
  type LayoutCursor,
  type PreparedTextWithSegments,
} from '@chenglou/pretext'

useHead({
  title: 'fatwill.cloud',
  meta: [{ name: 'description', content: '全栈开发者 fatwillzeng 的个人博客' }],
})

// ---- 代码片段 ----
const CODE_SNIPPETS = [
  '<div class="container mx-auto px-4 py-8">',
  'const handleClick = (e: MouseEvent) => { e.preventDefault(); router.push("/home") }',
  '.flex { display: flex; align-items: center; gap: 16px; justify-content: space-between; }',
  'import { ref, computed, onMounted, watch } from "vue"',
  'function fibonacci(n: number): number { return n <= 1 ? n : fib(n-1) + fib(n-2) }',
  'background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px;',
  'export default defineNuxtConfig({ ssr: true, modules: ["@pinia/nuxt", "@nuxtjs/tailwindcss"] })',
  'const [state, dispatch] = useReducer(reducer, initialState)',
  '@keyframes fadeIn { from { opacity: 0; transform: translateY(10px) } to { opacity: 1 } }',
  'SELECT id, title, created_at FROM articles WHERE published = 1 ORDER BY created_at DESC LIMIT 10',
  '<template><NuxtPage /><ClientOnly><AppLoading /></ClientOnly></template>',
  'border-radius: 12px; box-shadow: 0 8px 32px rgba(0,0,0,.12); backdrop-filter: blur(10px);',
  'async function fetchData<T>(url: string): Promise<T> { const res = await fetch(url); return res.json() }',
  'git commit -m "feat: integrate @chenglou/pretext for inline text layout"',
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
  'z-index: 100; position: sticky; top: 0; backdrop-filter: saturate(180%) blur(20px);',
  'Promise.all([fetchUser(), fetchPosts(), fetchComments()]).then(([user, posts, comments]) => {})',
  'rsync -az --delete .output/ root@server:/root/blog-frontend/.output/ -e "ssh -i ~/.ssh/key"',
]

const COLORS = [
  '#7dd3fc', '#4ade80', '#fb923c', '#a78bfa', '#f472b6',
  '#22d3ee', '#facc15', '#34d399', '#c084fc', '#f87171',
  '#38bdf8', '#86efac', '#fbbf24', '#60a5fa',
]

const CODE_FONT = '13px "Fira Code", "Courier New", monospace'
const LINE_HEIGHT = 26
const PADDING_X = 16
const TEXT_OPACITY = 0.22
const BTN_GAP = 12 // 按钮与文字之间的间距

// ---- refs ----
const canvasRef = ref<HTMLCanvasElement | null>(null)
const btnRef = ref<HTMLButtonElement | null>(null)
const btnLeft = ref(0)
const btnTop = ref(0)
const isDragging = ref(false)

// ---- 状态 ----
let ctx: CanvasRenderingContext2D | null = null
let animFrameId = 0
let offsetY = 0
const scrollSpeed = 0.5

// 按钮尺寸（从 DOM 读取）
let btnW = 160
let btnH = 52

// ---- 预处理文字 ----
// 把所有代码片段拼接成一段长文字，用 prepareWithSegments 预处理
// 然后用 layoutNextLine 逐行排版
let preparedText: PreparedTextWithSegments | null = null
let fullText = ''
let colorRanges: Array<{ start: number; end: number; color: string }> = []

// 预排版结果（全宽排版，用于计算总高度和循环）
interface PreLayoutLine {
  text: string
  width: number
  startCursor: LayoutCursor
  endCursor: LayoutCursor
  y: number // 内容坐标
  colorIdx: number // 用于着色
}

let preLayoutLines: PreLayoutLine[] = []
let totalContentHeight = 0

/**
 * 构建预处理文本和颜色映射
 */
function buildPreparedText() {
  // 拼接所有代码片段，用分隔符连接
  const separator = '   ·   '
  const parts: string[] = []
  // 重复多次确保内容足够
  for (let rep = 0; rep < 4; rep++) {
    for (const snippet of CODE_SNIPPETS) {
      parts.push(snippet)
    }
  }
  fullText = parts.join(separator)

  // 构建颜色映射：每个代码片段对应一个颜色
  colorRanges = []
  let pos = 0
  for (let rep = 0; rep < 4; rep++) {
    for (let i = 0; i < CODE_SNIPPETS.length; i++) {
      const start = pos
      pos += CODE_SNIPPETS[i]!.length
      colorRanges.push({ start, end: pos, color: COLORS[i % COLORS.length]! })
      pos += separator.length // 跳过分隔符
    }
  }

  // 用 pretext 预处理（正确的 API：text, font）
  preparedText = prepareWithSegments(fullText, CODE_FONT)
}

/**
 * 用全宽对文字进行预排版，记录每行的 cursor 和 y 坐标
 * 这些行在绘制时会根据按钮位置动态重新排版
 */
function buildPreLayout() {
  if (!preparedText) return

  const cssW = canvasRef.value ? canvasRef.value.width / (window.devicePixelRatio || 1) : window.innerWidth
  const contentWidth = cssW - PADDING_X * 2

  preLayoutLines = []
  let cursor: LayoutCursor = { segmentIndex: 0, graphemeIndex: 0 }
  let y = 0
  let lineIdx = 0
  const maxHeight = window.innerHeight * 4 // 足够循环

  while (y < maxHeight) {
    const line = layoutNextLine(preparedText, cursor, contentWidth)
    if (!line) break

    preLayoutLines.push({
      text: line.text,
      width: line.width,
      startCursor: { ...cursor },
      endCursor: { ...line.end },
      y,
      colorIdx: lineIdx,
    })

    cursor = line.end
    y += LINE_HEIGHT
    lineIdx++
  }

  totalContentHeight = y
}

/**
 * 根据文本在 fullText 中的字符位置，获取对应的颜色
 * 简化方案：用行索引取颜色
 */
function getLineColor(lineIdx: number): string {
  return COLORS[lineIdx % COLORS.length]!
}

/**
 * 核心绘制函数
 *
 * 关键逻辑：对于每一行，计算其屏幕 y 坐标后，
 * 判断是否与按钮区域重叠。如果重叠，则用 layoutNextLine
 * 传入受限宽度重新排版这行文字，文字从排版源头就绕开了按钮。
 */
function drawFrame() {
  if (!ctx || !canvasRef.value || !preparedText) {
    animFrameId = requestAnimationFrame(drawFrame)
    return
  }

  const dpr = window.devicePixelRatio || 1
  const w = canvasRef.value.width
  const h = canvasRef.value.height
  const cssW = w / dpr
  const cssH = h / dpr

  ctx.clearRect(0, 0, w, h)
  ctx.save()
  ctx.scale(dpr, dpr)

  // 更新滚动
  offsetY -= scrollSpeed
  if (totalContentHeight > 0 && -offsetY >= totalContentHeight) {
    offsetY += totalContentHeight
  }

  // 按钮屏幕坐标
  const bLeft = btnLeft.value
  const bTop = btnTop.value
  const bRight = bLeft + btnW
  const bBottom = bTop + btnH + 4

  ctx.font = CODE_FONT
  ctx.globalAlpha = TEXT_OPACITY

  // 绘制两遍实现无缝循环
  for (let pass = 0; pass < 2; pass++) {
    const passOffset = offsetY + pass * totalContentHeight

    for (let i = 0; i < preLayoutLines.length; i++) {
      const pLine = preLayoutLines[i]!
      const screenY = pLine.y + passOffset + LINE_HEIGHT // baseline

      // 屏幕外的行跳过
      if (screenY < -LINE_HEIGHT || screenY > cssH + LINE_HEIGHT) continue

      // 判断这一行是否与按钮重叠
      const lineTop = screenY - LINE_HEIGHT
      const lineBottom = screenY
      const overlapsWithBtn = lineBottom > bTop && lineTop < bBottom

      const color = getLineColor(pLine.colorIdx)
      ctx.fillStyle = color

      if (!overlapsWithBtn) {
        // ===== 普通行：全宽绘制 =====
        ctx.fillText(pLine.text, PADDING_X, screenY)
      } else {
        // ===== 按钮重叠行：用 layoutNextLine 动态宽度重新排版 =====
        // 核心用法：传入受限宽度，文字在 layout 阶段就只排在按钮左侧
        const availableWidth = bLeft - PADDING_X - BTN_GAP

        if (availableWidth > 40) {
          // 对这行文字重新 prepare + layoutNextLine
          // 用 pLine.text（已经是全宽排版的这行文字）重新排版
          const linePrepared = prepareWithSegments(pLine.text, CODE_FONT)
          const shorterLine = layoutNextLine(
            linePrepared,
            { segmentIndex: 0, graphemeIndex: 0 },
            availableWidth,
          )

          if (shorterLine && shorterLine.text) {
            ctx.fillText(shorterLine.text, PADDING_X, screenY)
          }
        }

        // 按钮右侧也可以绘制文字（如果有空间）
        const rightStart = bRight + BTN_GAP
        const rightAvailableWidth = cssW - rightStart - PADDING_X

        if (rightAvailableWidth > 40) {
          // 从 pLine.text 中取出按钮右侧应该显示的文字
          // 用全宽排版的行文字，跳过左侧已显示的部分
          const linePrepared = prepareWithSegments(pLine.text, CODE_FONT)

          // 先用左侧宽度排版获取消耗掉的文字
          const leftLine = layoutNextLine(
            linePrepared,
            { segmentIndex: 0, graphemeIndex: 0 },
            availableWidth > 40 ? availableWidth : 0.001,
          )

          if (leftLine) {
            // 从左侧消耗后的 cursor 继续排版右侧文字
            const rightLine = layoutNextLine(linePrepared, leftLine.end, rightAvailableWidth)
            if (rightLine && rightLine.text) {
              ctx.fillText(rightLine.text, rightStart, screenY)
            }
          }
        }
      }
    }
  }

  ctx.restore()
  animFrameId = requestAnimationFrame(drawFrame)
}

// ---- Canvas 初始化 ----
function initCanvas() {
  if (!canvasRef.value) return
  ctx = canvasRef.value.getContext('2d')
  resizeCanvas()
}

function resizeCanvas() {
  if (!canvasRef.value) return
  const dpr = window.devicePixelRatio || 1
  canvasRef.value.width = window.innerWidth * dpr
  canvasRef.value.height = window.innerHeight * dpr

  if (btnRef.value) {
    btnW = btnRef.value.offsetWidth || 160
    btnH = btnRef.value.offsetHeight || 52
  }

  buildPreLayout()
}

// ---- 长按拖拽 ----
let dragStartX = 0
let dragStartY = 0
let btnStartX = 0
let btnStartY = 0
let pressTimer: ReturnType<typeof setTimeout> | null = null
let pressMoved = false
let pressStartTime = 0

function onPointerDown(e: MouseEvent | TouchEvent) {
  e.stopPropagation()
  pressMoved = false
  pressStartTime = Date.now()

  const point = 'touches' in e ? e.touches[0]! : e
  dragStartX = point.clientX
  dragStartY = point.clientY
  btnStartX = btnLeft.value
  btnStartY = btnTop.value

  pressTimer = setTimeout(() => {
    isDragging.value = true
  }, 300)

  const onMove = (e2: MouseEvent | TouchEvent) => {
    const p = 'touches' in e2 ? e2.touches[0]! : (e2 as MouseEvent)
    const dx = p.clientX - dragStartX
    const dy = p.clientY - dragStartY

    if (Math.abs(dx) > 5 || Math.abs(dy) > 5) {
      pressMoved = true
    }

    if (isDragging.value) {
      btnLeft.value = btnStartX + dx
      btnTop.value = btnStartY + dy
    }
  }

  const onUp = () => {
    if (pressTimer) {
      clearTimeout(pressTimer)
      pressTimer = null
    }
    isDragging.value = false
    document.removeEventListener('mousemove', onMove as any)
    document.removeEventListener('mouseup', onUp)
    document.removeEventListener('touchmove', onMove as any)
    document.removeEventListener('touchend', onUp)
  }

  document.addEventListener('mousemove', onMove as any)
  document.addEventListener('mouseup', onUp)
  document.addEventListener('touchmove', onMove as any, { passive: true })
  document.addEventListener('touchend', onUp)
}

function onBtnClick() {
  if (!pressMoved && Date.now() - pressStartTime < 300) {
    navigateTo('/home')
  }
}

// ---- 生命周期 ----
onMounted(async () => {
  await nextTick()

  // 读取按钮实际尺寸
  if (btnRef.value) {
    btnW = btnRef.value.offsetWidth || 160
    btnH = btnRef.value.offsetHeight || 52
  }

  // 初始化按钮位置（屏幕右侧中央偏右）
  btnLeft.value = window.innerWidth - btnW - 60
  btnTop.value = window.innerHeight / 2 - btnH / 2

  // 构建预处理文本
  buildPreparedText()

  // 初始化 canvas
  initCanvas()

  // 启动动画
  animFrameId = requestAnimationFrame(drawFrame)
  window.addEventListener('resize', onResize)
})

function onResize() {
  if (btnRef.value) {
    btnW = btnRef.value.offsetWidth || 160
    btnH = btnRef.value.offsetHeight || 52
  }
  btnLeft.value = window.innerWidth - btnW - 60
  btnTop.value = window.innerHeight / 2 - btnH / 2
  resizeCanvas()
}

onUnmounted(() => {
  cancelAnimationFrame(animFrameId)
  animFrameId = 0
  window.removeEventListener('resize', onResize)
  if (pressTimer) clearTimeout(pressTimer)
})
</script>

<style scoped>
.index-page {
  position: fixed;
  inset: 0;
  background: #0d1117;
  overflow: hidden;
}

.bg-canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.enter-btn {
  position: fixed;
  z-index: 10;
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
  display: flex;
  align-items: center;
  gap: 10px;
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
}

.arrow {
  font-size: 18px;
  opacity: 0.7;
}

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

@media (max-width: 768px) {
  .enter-btn {
    padding: 12px 30px;
    font-size: 15px;
    letter-spacing: 1px;
  }
}
</style>
