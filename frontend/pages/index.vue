<template>
  <div ref="pageRef" class="landing-page" @click="handlePageClick">
    <!-- 代码滚动背景 -->
    <div class="code-bg">
      <div
        v-for="(line, i) in codeLines"
        :key="i"
        class="code-line"
        :style="{
          top: `${line.top}%`,
          animationDuration: `${line.speed}s`,
          animationDirection: line.direction,
          color: line.color,
          fontSize: `${line.fontSize}px`,
          opacity: line.opacity,
        }"
      >
        <span class="code-text">{{ line.text }}&nbsp;&nbsp;&nbsp;&nbsp;{{ line.text }}&nbsp;&nbsp;&nbsp;&nbsp;{{ line.text }}&nbsp;&nbsp;&nbsp;&nbsp;{{ line.text }}</span>
      </div>
    </div>

    <!-- 花盆+花 -->
    <div
      ref="flowerRef"
      class="flower-container"
      :class="{
        'fly-in': flyInStarted,
        'idle': idleStarted && !isDragging,
        'dragging': isDragging,
      }"
      :style="flowerStyle"
      @mousedown.stop="onPointerDown"
      @touchstart.stop="onPointerDown"
      @click.stop="onFlowerClick"
    >
      <!-- 花朵部分 -->
      <div class="flower-head-group" :class="{ 'wind-blow': flyInStarted && !idleStarted }">
        <!-- 花瓣 -->
        <div
          v-for="n in 5"
          :key="'petal-' + n"
          class="petal"
          :class="[`petal-${n}`, { 'wave-left': waving && n % 2 === 1, 'wave-right': waving && n % 2 === 0 }]"
          :style="{ transform: `rotate(${(n - 1) * 72}deg) translateY(-28px)` }"
        />
        <!-- 花脸 -->
        <div class="flower-face">
          <!-- 眯眼（飞入时） -->
          <template v-if="!showSmile">
            <div class="squint-eye left" />
            <div class="squint-eye right" />
          </template>
          <!-- 笑脸（停下后） -->
          <template v-else>
            <div class="smile-eye left" />
            <div class="smile-eye right" />
            <div class="smile-mouth" />
          </template>
          <!-- 腮红 -->
          <div class="blush left" />
          <div class="blush right" />
        </div>
      </div>

      <!-- 花茎 -->
      <div class="stem" :class="{ 'wind-blow-stem': flyInStarted && !idleStarted }">
        <!-- 左叶 -->
        <div class="leaf leaf-left" :class="{ 'wave-leaf-left': waving }" />
        <!-- 右叶 -->
        <div class="leaf leaf-right" :class="{ 'wave-leaf-right': waving }" />
      </div>

      <!-- 花盆 -->
      <div class="pot">
        <div class="pot-rim" />
        <div class="pot-body">
          <div class="pot-stripe" />
        </div>
        <div class="pot-base" />
      </div>
    </div>

    <!-- 进入主页按钮 -->
    <NuxtLink to="/home" class="enter-btn" :class="{ 'show': idleStarted }">
      <span class="enter-text">进入博客</span>
      <svg class="enter-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
      </svg>
    </NuxtLink>

    <!-- 底部署名 -->
    <div class="landing-footer" :class="{ 'show': idleStarted }">
      <span>fatwill.cloud</span>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'fatwill.cloud',
  meta: [
    { name: 'description', content: '全栈开发者 fatwillzeng 的个人博客' },
  ],
})

// ====== 代码滚动行 ======
const codeSnippets = [
  '<div class="container mx-auto">',
  'const handleClick = (e: Event) => {',
  '.flex { display: flex; align-items: center; }',
  'import React from "react"',
  'function fibonacci(n: number): number {',
  'background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);',
  'export default defineNuxtConfig({',
  'const [state, setState] = useState(false)',
  '@keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }',
  'SELECT * FROM articles ORDER BY created_at DESC',
  '<template><NuxtPage /></template>',
  'border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,.1);',
  'async function fetchData(url: string) {',
  'git commit -m "feat: add new feature"',
  'docker-compose up -d --build',
  'npm install -D tailwindcss postcss autoprefixer',
  'interface User { id: string; name: string; }',
  'transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);',
  'const router = useRouter()',
  'CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT);',
]

const codeColors = [
  '#7dd3fc', '#4ade80', '#fb923c', '#a78bfa',
  '#f472b6', '#22d3ee', '#facc15', '#34d399',
  '#c084fc', '#f87171', '#38bdf8', '#86efac',
]

interface CodeLine {
  text: string
  top: number
  speed: number
  direction: string
  color: string
  fontSize: number
  opacity: number
}

const codeLines = ref<CodeLine[]>([])

function generateCodeLines() {
  const lines: CodeLine[] = []
  const lineCount = 18
  for (let i = 0; i < lineCount; i++) {
    lines.push({
      text: codeSnippets[i % codeSnippets.length],
      top: (i / lineCount) * 100,
      speed: 20 + Math.random() * 40,
      direction: i % 2 === 0 ? 'normal' : 'reverse',
      color: codeColors[i % codeColors.length],
      fontSize: 13 + Math.random() * 5,
      opacity: 0.12 + Math.random() * 0.18,
    })
  }
  codeLines.value = lines
}

generateCodeLines()

// ====== 花盆动画状态 ======
const pageRef = ref<HTMLElement | null>(null)
const flowerRef = ref<HTMLElement | null>(null)
const flyInStarted = ref(false)
const idleStarted = ref(false)
const showSmile = ref(false)
const waving = ref(false)
const isDragging = ref(false)

// 花盆位置
const flowerX = ref(0)
const flowerY = ref(0)
const restX = ref(0)
const restY = ref(0)

const flowerStyle = computed(() => {
  if (isDragging.value) {
    return {
      left: `${flowerX.value}px`,
      top: `${flowerY.value}px`,
      transition: 'none',
    }
  }
  return {}
})

// ====== 入场动画 ======
onMounted(() => {
  // 计算停留位置
  const vw = window.innerWidth
  const vh = window.innerHeight
  restX.value = vw * 0.08
  restY.value = vh * 0.42

  // 延迟触发飞入
  requestAnimationFrame(() => {
    flyInStarted.value = true
  })

  // 飞入完成后切换到 idle
  setTimeout(() => {
    idleStarted.value = true
  }, 500)

  // 切换笑脸
  setTimeout(() => {
    showSmile.value = true
  }, 700)
})

// ====== 点击打招呼 ======
let clickTimer: ReturnType<typeof setTimeout> | null = null
let isLongPress = false

function onFlowerClick(e: Event) {
  e.stopPropagation()
  if (isLongPress || isDragging.value) return
  if (waving.value) return
  waving.value = true
  setTimeout(() => {
    waving.value = false
  }, 600)
}

// ====== 长按拖拽 ======
function onPointerDown(e: MouseEvent | TouchEvent) {
  e.stopPropagation()
  isLongPress = false
  const startPos = getPointerPos(e)
  const el = flowerRef.value
  if (!el) return

  const rect = el.getBoundingClientRect()
  const offsetX = startPos.x - rect.left
  const offsetY = startPos.y - rect.top

  clickTimer = setTimeout(() => {
    isLongPress = true
    isDragging.value = true
    flowerX.value = startPos.x - offsetX
    flowerY.value = startPos.y - offsetY

    const onMove = (ev: MouseEvent | TouchEvent) => {
      const pos = getPointerPos(ev)
      flowerX.value = pos.x - offsetX
      flowerY.value = pos.y - offsetY
    }

    const onUp = () => {
      isDragging.value = false
      document.removeEventListener('mousemove', onMove)
      document.removeEventListener('mouseup', onUp)
      document.removeEventListener('touchmove', onMove)
      document.removeEventListener('touchend', onUp)
    }

    document.addEventListener('mousemove', onMove)
    document.addEventListener('mouseup', onUp)
    document.addEventListener('touchmove', onMove, { passive: true })
    document.addEventListener('touchend', onUp)
  }, 300)

  const onEarlyUp = () => {
    if (clickTimer) {
      clearTimeout(clickTimer)
      clickTimer = null
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
  const me = e as MouseEvent
  return { x: me.clientX, y: me.clientY }
}

function handlePageClick() {
  // 点击页面空白处无特殊行为
}
</script>

<style scoped>
/* ====== 页面基础 ====== */
.landing-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: #0d1117;
  cursor: default;
  user-select: none;
}

/* ====== 代码滚动背景 ====== */
.code-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  z-index: 0;
}

.code-line {
  position: absolute;
  white-space: nowrap;
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', 'Consolas', monospace;
  animation: scrollLine 30s linear infinite;
  pointer-events: none;
}

.code-text {
  display: inline-block;
}

@keyframes scrollLine {
  0% { transform: translateX(-25%); }
  100% { transform: translateX(0%); }
}

/* ====== 花盆容器 ====== */
.flower-container {
  position: fixed;
  left: 8%;
  top: 42%;
  transform: translate(-300px, 0) scale(1);
  z-index: 10;
  cursor: grab;
  width: 150px;
  display: flex;
  flex-direction: column;
  align-items: center;
  will-change: transform;
  transition: none;
}

.flower-container.fly-in {
  transform: translate(0, 0) scale(1);
  transition: transform 400ms cubic-bezier(0.34, 1.56, 0.64, 1);
}

.flower-container.idle {
  animation: idleSway 3s ease-in-out infinite;
}

.flower-container.dragging {
  position: fixed;
  cursor: grabbing;
  animation: none !important;
}

@keyframes idleSway {
  0%, 100% { transform: translate(0, 0) rotate(-2deg); }
  50% { transform: translate(0, 0) rotate(2deg); }
}

/* ====== 花头组 ====== */
.flower-head-group {
  position: relative;
  width: 90px;
  height: 90px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.flower-head-group.wind-blow {
  transform: rotate(15deg);
  transition: transform 0.1s ease;
}

/* ====== 花瓣 ====== */
.petal {
  position: absolute;
  width: 30px;
  height: 42px;
  border-radius: 50% 50% 50% 50% / 60% 60% 40% 40%;
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 50%, #d63031 100%);
  box-shadow: inset -2px -2px 4px rgba(0,0,0,0.15), inset 2px 2px 4px rgba(255,255,255,0.2);
  transform-origin: center 42px;
  transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.flower-head-group.wind-blow .petal {
  transform-origin: center 42px;
}

.flower-head-group.wind-blow .petal-1 { transform: rotate(0deg) translateY(-28px) rotate(35deg) !important; }
.flower-head-group.wind-blow .petal-2 { transform: rotate(72deg) translateY(-28px) rotate(30deg) !important; }
.flower-head-group.wind-blow .petal-3 { transform: rotate(144deg) translateY(-28px) rotate(25deg) !important; }
.flower-head-group.wind-blow .petal-4 { transform: rotate(216deg) translateY(-28px) rotate(30deg) !important; }
.flower-head-group.wind-blow .petal-5 { transform: rotate(288deg) translateY(-28px) rotate(35deg) !important; }

/* 花瓣呼吸动画 */
.idle .petal {
  animation: petalBreath 2.5s ease-in-out infinite;
}

.petal-1 { animation-delay: 0s; }
.petal-2 { animation-delay: 0.2s; }
.petal-3 { animation-delay: 0.4s; }
.petal-4 { animation-delay: 0.6s; }
.petal-5 { animation-delay: 0.8s; }

@keyframes petalBreath {
  0%, 100% { scale: 1; }
  50% { scale: 1.06; }
}

/* 挥手时花瓣动画 */
.wave-left {
  animation: wavePetalLeft 0.6s ease-in-out !important;
}
.wave-right {
  animation: wavePetalRight 0.6s ease-in-out !important;
}

@keyframes wavePetalLeft {
  0%, 100% { transform: var(--petal-transform); }
  50% { scale: 1.12; }
}
@keyframes wavePetalRight {
  0%, 100% { transform: var(--petal-transform); }
  50% { scale: 0.9; }
}

/* ====== 花脸 ====== */
.flower-face {
  position: relative;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 35%, #fff3a0, #ffd93d 40%, #f0c420 100%);
  box-shadow: 0 2px 8px rgba(0,0,0,0.15), inset 0 -2px 4px rgba(0,0,0,0.08);
  z-index: 2;
}

/* 眯眼 */
.squint-eye {
  position: absolute;
  top: 42%;
  width: 10px;
  height: 3px;
  border-radius: 2px;
  background: #5a3e28;
}
.squint-eye.left { left: 25%; transform: rotate(-5deg); }
.squint-eye.right { right: 25%; transform: rotate(5deg); }

/* 笑脸眼睛 */
.smile-eye {
  position: absolute;
  top: 38%;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #3d2b1f;
  animation: eyeBlink 4s ease-in-out infinite;
}
.smile-eye.left { left: 28%; }
.smile-eye.right { right: 28%; }

@keyframes eyeBlink {
  0%, 90%, 100% { transform: scaleY(1); }
  95% { transform: scaleY(0.1); }
}

/* 笑嘴 */
.smile-mouth {
  position: absolute;
  top: 56%;
  left: 50%;
  transform: translateX(-50%);
  width: 14px;
  height: 7px;
  border-bottom: 2.5px solid #5a3e28;
  border-radius: 0 0 50% 50%;
}

/* 腮红 */
.blush {
  position: absolute;
  top: 52%;
  width: 10px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 130, 130, 0.45);
}
.blush.left { left: 10%; }
.blush.right { right: 10%; }

/* ====== 花茎 ====== */
.stem {
  width: 6px;
  height: 70px;
  background: linear-gradient(180deg, #27ae60, #2ecc71);
  border-radius: 3px;
  position: relative;
  margin-top: -4px;
  transition: transform 0.3s ease;
}

.stem.wind-blow-stem {
  transform: skewX(-8deg);
}

/* ====== 叶子 ====== */
.leaf {
  position: absolute;
  width: 28px;
  height: 14px;
  border-radius: 50% 0 50% 0;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  box-shadow: inset -1px -1px 2px rgba(0,0,0,0.1);
  transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.leaf-left {
  left: -24px;
  top: 20px;
  transform: rotate(-20deg);
  border-radius: 0 50% 0 50%;
}

.leaf-right {
  right: -24px;
  top: 38px;
  transform: rotate(20deg);
}

/* idle 叶子随风 */
.idle .leaf-left {
  animation: leafSwayLeft 3.5s ease-in-out infinite;
}
.idle .leaf-right {
  animation: leafSwayRight 3.2s ease-in-out infinite 0.3s;
}

@keyframes leafSwayLeft {
  0%, 100% { transform: rotate(-20deg); }
  50% { transform: rotate(-30deg) translateY(-2px); }
}

@keyframes leafSwayRight {
  0%, 100% { transform: rotate(20deg); }
  50% { transform: rotate(30deg) translateY(-2px); }
}

/* 风吹叶子偏转 */
.wind-blow-stem .leaf-left {
  transform: rotate(-45deg) !important;
  animation: none !important;
}
.wind-blow-stem .leaf-right {
  transform: rotate(40deg) translateY(-3px) !important;
  animation: none !important;
}

/* 挥手叶子动画 */
.wave-leaf-left {
  animation: waveLeafL 0.6s ease-in-out !important;
}
.wave-leaf-right {
  animation: waveLeafR 0.6s ease-in-out !important;
}

@keyframes waveLeafL {
  0%, 100% { transform: rotate(-20deg); }
  40% { transform: rotate(-60deg) translateY(-4px); }
}
@keyframes waveLeafR {
  0%, 100% { transform: rotate(20deg); }
  50% { transform: rotate(60deg) translateY(-4px); }
}

/* ====== 花盆 ====== */
.pot {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: -2px;
}

.pot-rim {
  width: 62px;
  height: 10px;
  background: linear-gradient(180deg, #e67e22, #d35400);
  border-radius: 4px 4px 0 0;
  box-shadow: 0 1px 3px rgba(0,0,0,0.15);
}

.pot-body {
  width: 56px;
  height: 42px;
  background: linear-gradient(180deg, #e67e22 0%, #ca6f1e 60%, #b8600f 100%);
  clip-path: polygon(4% 0%, 96% 0%, 85% 100%, 15% 100%);
  position: relative;
  overflow: hidden;
}

.pot-stripe {
  position: absolute;
  top: 40%;
  left: 10%;
  right: 10%;
  height: 3px;
  background: rgba(255,255,255,0.2);
  border-radius: 2px;
}

.pot-base {
  width: 36px;
  height: 6px;
  background: linear-gradient(180deg, #b8600f, #a0510d);
  border-radius: 0 0 3px 3px;
}

/* ====== 进入按钮 ====== */
.enter-btn {
  position: fixed;
  right: 8%;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 32px;
  border: 1px solid rgba(255,255,255,0.15);
  border-radius: 60px;
  background: rgba(255,255,255,0.05);
  backdrop-filter: blur(10px);
  color: rgba(255,255,255,0.8);
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 1px;
  text-decoration: none;
  opacity: 0;
  transition: all 0.5s ease;
  pointer-events: none;
  z-index: 20;
}

.enter-btn.show {
  opacity: 1;
  pointer-events: auto;
}

.enter-btn:hover {
  background: rgba(255,255,255,0.12);
  border-color: rgba(255,255,255,0.3);
  color: #fff;
  box-shadow: 0 0 30px rgba(99, 102, 241, 0.2);
}

.enter-btn:hover .enter-arrow {
  transform: translateX(4px);
}

.enter-text {
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.enter-arrow {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease;
}

/* ====== 底部署名 ====== */
.landing-footer {
  position: fixed;
  bottom: 24px;
  left: 0;
  right: 0;
  text-align: center;
  color: rgba(255,255,255,0.2);
  font-size: 13px;
  font-family: 'Fira Code', monospace;
  letter-spacing: 2px;
  opacity: 0;
  transition: opacity 0.8s ease 0.3s;
  z-index: 20;
  pointer-events: none;
}

.landing-footer.show {
  opacity: 1;
}

/* ====== 移动端适配 ====== */
@media (max-width: 768px) {
  .flower-container {
    left: 5%;
    top: 35%;
    transform: translate(-200px, 0) scale(0.85);
  }

  .flower-container.fly-in {
    transform: translate(0, 0) scale(0.85);
  }

  .flower-container.idle {
    animation: idleSwayMobile 3s ease-in-out infinite;
  }

  @keyframes idleSwayMobile {
    0%, 100% { transform: translate(0, 0) rotate(-2deg) scale(0.85); }
    50% { transform: translate(0, 0) rotate(2deg) scale(0.85); }
  }

  .enter-btn {
    right: auto;
    left: 50%;
    top: auto;
    bottom: 80px;
    transform: translateX(-50%);
    padding: 12px 28px;
    font-size: 15px;
  }

  .enter-btn.show {
    transform: translateX(-50%);
  }

  .code-line {
    font-size: 11px !important;
  }
}
</style>
