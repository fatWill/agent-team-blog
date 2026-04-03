<template>
  <div class="index-page">
    <canvas ref="canvasRef" class="bg-canvas" />
    <!-- 按钮由排版引擎定位，position: fixed -->
    <button
      ref="btnRef"
      class="enter-btn"
      :class="{ dragging: isDragging }"
      :style="btnStyle"
      @mousedown.prevent="onPointerDown"
      @touchstart.prevent="onPointerDown"
    >
      <span class="enter-text">进入博客</span>
      <svg class="enter-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
      </svg>
    </button>
    <div class="footer">fatwill.cloud</div>
  </div>
</template>

<script setup lang="ts">
import {
  prepareWithSegments,
  layoutNextLine,
  walkLineRanges,
  type LayoutCursor,
  type PreparedTextWithSegments,
} from '@chenglou/pretext'

useHead({
  title: 'fatwill.cloud',
  meta: [{ name: 'description', content: '全栈开发者 fatwillzeng 的个人博客' }],
})

// ====== 代码片段 ======
const CODE_SNIPPETS = [
  '<div class="container mx-auto px-4 py-8">',
  'const handleClick = (e: MouseEvent) => { e.preventDefault(); router.push("/home") }',
  '.flex { display: flex; align-items: center; gap: 16px; justify-content: space-between; }',
  'import { ref, computed, onMounted, watch } from "vue"',
  'function fibonacci(n: number): number { return n <= 1 ? n : fib(n-1) + fib(n-2) }',
  'background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px;',
  'export default defineNuxtConfig({ ssr: true, modules: ["@pinia/nuxt", "@nuxtjs/tailwindcss"] })',
  'const [state, dispatch] = useReducer(reducer, initialState)',
  '@keyframes fadeIn { from { opacity: 0 } to { opacity: 1; transform: none } }',
  'SELECT id, title, created_at FROM articles WHERE published = 1 ORDER BY created_at DESC',
  '<template><NuxtPage /><ClientOnly><AppLoading /></ClientOnly></template>',
  'border-radius: 12px; box-shadow: 0 8px 32px rgba(0,0,0,.12); backdrop-filter: blur(10px);',
  'async function fetchData<T>(url: string): Promise<T> { return fetch(url).then(r => r.json()) }',
  'git commit -m "feat: integrate @chenglou/pretext for inline text layout"',
  'docker-compose up -d --build && docker ps -a | grep blog',
  'npm install @chenglou/pretext @pinia/nuxt @nuxtjs/tailwindcss --save',
  'interface Article { id: string; title: string; content: string; tags: string[]; createdAt: Date }',
  'transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); will-change: transform, opacity;',
  'const { data, pending } = await useAsyncData("articles", () => $fetch("/api/articles"))',
  'CREATE TABLE articles (id INT PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255) NOT NULL);',
  'pm2 start ecosystem.config.cjs --env production && pm2 save && pm2 logs --lines 20',
  'type Props = { title: string; description?: string; onClick?: () => void }',
  'nginx -t && systemctl reload nginx  # test config and graceful reload',
  'const router = useRouter(); watch(() => route.path, () => { scrollTo(0, 0) })',
]

const COLORS = [
  '#7dd3fc', '#4ade80', '#fb923c', '#a78bfa', '#f472b6',
  '#22d3ee', '#facc15', '#34d399', '#c084fc', '#f87171',
  '#38bdf8', '#86efac',
]

// ====== 字体 & 排版参数 ======
const CODE_FONT = '14px "Fira Code", "Cascadia Code", "JetBrains Mono", Consolas, monospace'
const BTN_FONT = '500 17px "PingFang SC", "Microsoft YaHei", Inter, sans-serif'
const LINE_HEIGHT = 30
const UNBOUNDED_WIDTH = 100_000
const BTN_CHROME_WIDTH = 120 // 按钮 padding + gap + 箭头（不含文字宽度）
const BTN_TEXT = '进入博客 →'
const TEXT_OPACITY = 0.22

// ====== 类型定义（参考 rich-note.ts）======

// 预处理后的内联元素
type TextInlineItem = {
  kind: 'text'
  color: string
  chromeWidth: number
  endCursor: LayoutCursor
  fullText: string
  fullWidth: number
  leadingGap: number
  prepared: PreparedTextWithSegments
}

type ChipInlineItem = {
  kind: 'chip'
  leadingGap: number
  text: string
  width: number
}

type InlineItem = TextInlineItem | ChipInlineItem

// 排版结果：每行的片段
type LineFragment = {
  kind: 'text' | 'chip'
  color: string
  leadingGap: number
  text: string
  width: number
}

type RichLine = {
  fragments: LineFragment[]
}

// ====== Refs ======
const canvasRef = ref<HTMLCanvasElement | null>(null)
const btnRef = ref<HTMLButtonElement | null>(null)
const isDragging = ref(false)
const btnLeft = ref(-9999)
const btnTop = ref(-9999)
const btnVisible = ref(false)
const isDragOverride = ref(false) // 拖拽后按钮位置由用户控制

const btnStyle = computed(() => ({
  left: `${btnLeft.value}px`,
  top: `${btnTop.value}px`,
  opacity: btnVisible.value ? 1 : 0,
  pointerEvents: btnVisible.value ? 'auto' as const : 'none' as const,
}))

// ====== 动画状态 ======
let ctx: CanvasRenderingContext2D | null = null
let animId = 0
let canvasW = 0
let canvasH = 0
let offsetY = 0
let scrollSpeed = 0.5

// 排版结果
let layoutLines: RichLine[] = []
let totalHeight = 0
let btnLineIndex = -1 // 按钮所在行
let btnFragIndex = -1 // 按钮在行内的片段索引
let btnFragX = 0 // 按钮在行内的 x 坐标

// ====== 空格宽度缓存 ======
const spaceWidthCache = new Map<string, number>()
function measureCollapsedSpaceWidth(font: string): number {
  const cached = spaceWidthCache.get(font)
  if (cached !== undefined) return cached
  const joined = prepareWithSegments('A A', font)
  const compact = prepareWithSegments('AA', font)
  let joinedW = 0
  let compactW = 0
  walkLineRanges(joined, UNBOUNDED_WIDTH, (l: any) => { joinedW = l.width })
  walkLineRanges(compact, UNBOUNDED_WIDTH, (l: any) => { compactW = l.width })
  const w = Math.max(0, joinedW - compactW)
  spaceWidthCache.set(font, w)
  return w
}

function measureSingleLineWidth(prepared: PreparedTextWithSegments): number {
  let maxW = 0
  walkLineRanges(prepared, UNBOUNDED_WIDTH, (l: any) => { if (l.width > maxW) maxW = l.width })
  return maxW
}

// ====== 构建 InlineItem 数组 ======
function buildInlineItems(): { items: InlineItem[]; btnIndex: number } {
  const LINE_START: LayoutCursor = { segmentIndex: 0, graphemeIndex: 0 }
  const gap = measureCollapsedSpaceWidth(CODE_FONT)
  const items: InlineItem[] = []

  // 计算按钮应该在第几个 item 后面插入
  // 目标：让按钮大约落在屏幕中间
  const targetLine = Math.floor(canvasH / 2 / LINE_HEIGHT)
  const charsPerLine = Math.max(40, Math.floor(canvasW / 8.5)) // 估算每行字符数
  const charsBeforeBtn = targetLine * charsPerLine
  let totalChars = 0
  let snippetIdx = 0

  // 按钮前的代码文字
  while (totalChars < charsBeforeBtn) {
    const snippet = CODE_SNIPPETS[snippetIdx % CODE_SNIPPETS.length]!
    const trimmed = snippet.trim()
    if (trimmed.length === 0) { snippetIdx++; continue }

    const prepared = prepareWithSegments(trimmed, CODE_FONT)
    const wholeLine = layoutNextLine(prepared, LINE_START, UNBOUNDED_WIDTH)
    if (wholeLine) {
      items.push({
        kind: 'text',
        color: COLORS[snippetIdx % COLORS.length]!,
        chromeWidth: 0,
        endCursor: wholeLine.end,
        fullText: wholeLine.text,
        fullWidth: wholeLine.width,
        leadingGap: gap,
        prepared,
      })
    }
    totalChars += trimmed.length
    snippetIdx++
  }

  // 按钮（chip 类型，原子元素不被截断）
  const btnIndex = items.length
  const btnPrepared = prepareWithSegments(BTN_TEXT, BTN_FONT)
  const btnTextWidth = measureSingleLineWidth(btnPrepared)
  items.push({
    kind: 'chip',
    leadingGap: gap,
    text: BTN_TEXT,
    width: Math.ceil(btnTextWidth) + BTN_CHROME_WIDTH,
  })

  // 按钮后的代码文字（填满剩余屏幕 + 余量）
  const totalTarget = charsBeforeBtn * 2.5
  while (totalChars < totalTarget) {
    const snippet = CODE_SNIPPETS[snippetIdx % CODE_SNIPPETS.length]!
    const trimmed = snippet.trim()
    if (trimmed.length === 0) { snippetIdx++; continue }

    const prepared = prepareWithSegments(trimmed, CODE_FONT)
    const wholeLine = layoutNextLine(prepared, LINE_START, UNBOUNDED_WIDTH)
    if (wholeLine) {
      items.push({
        kind: 'text',
        color: COLORS[snippetIdx % COLORS.length]!,
        chromeWidth: 0,
        endCursor: wholeLine.end,
        fullText: wholeLine.text,
        fullWidth: wholeLine.width,
        leadingGap: gap,
        prepared,
      })
    }
    totalChars += trimmed.length
    snippetIdx++
  }

  return { items, btnIndex }
}

// ====== 排版算法（直接参考 rich-note.ts 的 layoutInlineItems）======
function cursorsMatch(a: LayoutCursor, b: LayoutCursor): boolean {
  return a.segmentIndex === b.segmentIndex && a.graphemeIndex === b.graphemeIndex
}

const LINE_START_CURSOR: LayoutCursor = { segmentIndex: 0, graphemeIndex: 0 }

function layoutInlineItems(items: InlineItem[], maxWidth: number, btnIndex: number): RichLine[] {
  const lines: RichLine[] = []
  const safeWidth = Math.max(1, maxWidth)

  let itemIndex = 0
  let textCursor: LayoutCursor | null = null

  // 重置按钮位置信息
  btnLineIndex = -1
  btnFragIndex = -1
  btnFragX = 0

  while (itemIndex < items.length) {
    const fragments: LineFragment[] = []
    let lineWidth = 0
    let remainingWidth = safeWidth

    lineLoop:
    while (itemIndex < items.length) {
      const item = items[itemIndex]!

      switch (item.kind) {
        case 'chip': {
          const leadingGap = fragments.length === 0 ? 0 : item.leadingGap
          // chip 是原子元素：如果放不下就换行
          if (fragments.length > 0 && leadingGap + item.width > remainingWidth) break lineLoop

          // 记录按钮位置
          if (itemIndex === btnIndex) {
            btnLineIndex = lines.length
            btnFragIndex = fragments.length
            btnFragX = lineWidth + leadingGap
          }

          fragments.push({
            kind: 'chip',
            color: '',
            leadingGap,
            text: item.text,
            width: item.width,
          })
          lineWidth += leadingGap + item.width
          remainingWidth = Math.max(0, safeWidth - lineWidth)
          itemIndex++
          textCursor = null
          continue
        }

        case 'text': {
          // 如果文字已经完全消耗完，跳到下一个 item
          if (textCursor !== null && cursorsMatch(textCursor, item.endCursor)) {
            itemIndex++
            textCursor = null
            continue
          }

          const leadingGap = fragments.length === 0 ? 0 : item.leadingGap
          const reservedWidth = leadingGap + item.chromeWidth
          if (fragments.length > 0 && reservedWidth >= remainingWidth) break lineLoop

          // 快速路径：整段文字放得下
          if (textCursor === null) {
            const fullWidth = leadingGap + item.fullWidth + item.chromeWidth
            if (fullWidth <= remainingWidth) {
              fragments.push({
                kind: 'text',
                color: item.color,
                leadingGap,
                text: item.fullText,
                width: item.fullWidth,
              })
              lineWidth += fullWidth
              remainingWidth = Math.max(0, safeWidth - lineWidth)
              itemIndex++
              continue
            }
          }

          // 慢速路径：用 layoutNextLine 截断文字
          const startCursor = textCursor ?? LINE_START_CURSOR
          const line = layoutNextLine(
            item.prepared,
            startCursor,
            Math.max(1, remainingWidth - reservedWidth),
          )
          if (line === null) {
            itemIndex++
            textCursor = null
            continue
          }
          if (cursorsMatch(startCursor, line.end)) {
            itemIndex++
            textCursor = null
            continue
          }

          fragments.push({
            kind: 'text',
            color: item.color,
            leadingGap,
            text: line.text,
            width: line.width,
          })
          lineWidth += leadingGap + line.width + item.chromeWidth
          remainingWidth = Math.max(0, safeWidth - lineWidth)

          // 检查文字是否已经消耗完
          if (cursorsMatch(line.end, item.endCursor)) {
            itemIndex++
            textCursor = null
            continue
          }

          // 没消耗完，记录光标，下一行继续
          textCursor = line.end
          break lineLoop
        }
      }
    }

    if (fragments.length === 0) break
    lines.push({ fragments })
  }

  return lines
}

// ====== 计算排版 ======
function computeLayout() {
  if (!ctx) return

  const { items, btnIndex } = buildInlineItems()
  layoutLines = layoutInlineItems(items, canvasW, btnIndex)
  totalHeight = layoutLines.length * LINE_HEIGHT + LINE_HEIGHT // 加一行余量

  // 更新按钮 DOM 位置
  if (!isDragOverride.value) {
    updateBtnFromLayout()
  }
}

function updateBtnFromLayout() {
  if (btnLineIndex < 0) return
  const y = btnLineIndex * LINE_HEIGHT
  // 按钮排版位置 = btnFragX (已在 layoutInlineItems 中计算)
  // 存储基准位置，每帧绘制时根据 offsetY 更新
  btnLayoutBaseX = btnFragX
  btnLayoutBaseY = y
}

let btnLayoutBaseX = 0
let btnLayoutBaseY = 0

// ====== 每帧绘制 ======
function drawFrame() {
  if (!ctx || !canvasRef.value) {
    animId = requestAnimationFrame(drawFrame)
    return
  }

  const dpr = window.devicePixelRatio || 1
  ctx.clearRect(0, 0, canvasW * dpr, canvasH * dpr)
  ctx.save()
  ctx.scale(dpr, dpr)

  // 更新滚动偏移
  offsetY -= scrollSpeed
  // 循环：当内容滚完一遍，重置
  if (totalHeight > 0 && Math.abs(offsetY) >= totalHeight) {
    offsetY %= totalHeight
  }

  // 绘制两遍实现无缝循环
  for (let pass = 0; pass < 2; pass++) {
    const passOff = offsetY + pass * totalHeight

    for (let lineIdx = 0; lineIdx < layoutLines.length; lineIdx++) {
      const lineY = lineIdx * LINE_HEIGHT + LINE_HEIGHT + passOff // +LINE_HEIGHT 留出顶部间距
      if (lineY < -LINE_HEIGHT || lineY > canvasH + LINE_HEIGHT) continue

      const line = layoutLines[lineIdx]!
      let x = 0

      for (const frag of line.fragments) {
        x += frag.leadingGap

        if (frag.kind === 'text') {
          ctx.font = CODE_FONT
          ctx.fillStyle = frag.color
          ctx.globalAlpha = TEXT_OPACITY
          ctx.fillText(frag.text, x, lineY)
        }
        // chip（按钮）不在 canvas 绘制，由 DOM 元素覆盖

        x += frag.width
      }
    }
  }

  ctx.restore()

  // 更新按钮 DOM 位置（跟随滚动）
  if (!isDragOverride.value && btnLineIndex >= 0) {
    // 找到按钮当前应该在的 Y 坐标
    const baseY = btnLayoutBaseY + LINE_HEIGHT + offsetY
    // 两遍循环中找到可见的那一遍
    let visY = baseY
    if (visY < -LINE_HEIGHT) visY += totalHeight
    if (visY > canvasH + LINE_HEIGHT) visY -= totalHeight

    // 按钮 DOM 位置（左上角）
    const btnEl = btnRef.value
    if (btnEl) {
      const btnH = btnEl.offsetHeight || 48
      // 文字 baseline 在 lineY，按钮需要垂直居中对齐
      btnLeft.value = btnLayoutBaseX
      btnTop.value = visY - btnH * 0.65
      btnVisible.value = true
    }
  }

  animId = requestAnimationFrame(drawFrame)
}

// ====== Canvas 初始化 & resize ======
function initCanvas() {
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
}

function handleResize() {
  initCanvas()
  offsetY = 0
  isDragOverride.value = false
  computeLayout()
}

// ====== 按钮拖拽 ======
let pointerDownTimer: ReturnType<typeof setTimeout> | null = null
let isLongPress = false
let hasMoved = false
let startPointerX = 0
let startPointerY = 0

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
    isDragOverride.value = true

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

// ====== 生命周期 ======
onMounted(() => {
  initCanvas()
  computeLayout()
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
.index-page {
  position: fixed;
  inset: 0;
  background: #0d1117;
  overflow: hidden;
  user-select: none;
}

.bg-canvas {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

/* ====== 进入按钮 ====== */
.enter-btn {
  position: fixed;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 40px;
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
  transition: box-shadow 0.3s ease, border-color 0.3s ease, background 0.3s ease, opacity 0.3s ease;
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
    padding: 12px 30px;
    font-size: 15px;
    letter-spacing: 1px;
  }
}
</style>
