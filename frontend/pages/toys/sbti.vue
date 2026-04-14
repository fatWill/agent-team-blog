<template>
  <div class="min-h-screen bg-gray-50 text-gray-800 transition-colors duration-300 dark:bg-gray-900 dark:text-gray-200">
    <!-- 返回按钮 -->
    <NuxtLink
      to="/tools"
      class="fixed left-4 top-4 z-50 flex h-9 items-center gap-1.5 rounded-full bg-white/80 px-3 text-sm text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
    >
      <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12l8.954-8.955a1.126 1.126 0 011.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
      </svg>
      <span>首页</span>
    </NuxtLink>

    <!-- 主题切换 -->
    <button
      class="fixed right-4 top-4 z-50 flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-gray-500 shadow-sm backdrop-blur-lg transition-colors hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800/80 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-200"
      aria-label="切换主题"
      @click="toggleTheme"
    >
      <svg v-if="isDark" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
      <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
      </svg>
    </button>

    <!-- ========== 欢迎页 ========== -->
    <div v-if="phase === 'welcome'" class="flex min-h-screen flex-col items-center justify-center px-4 py-16">
      <div class="mx-auto max-w-xl text-center">
        <!-- Logo / 标题 -->
        <div class="mb-6 inline-flex items-center gap-2 rounded-full bg-gradient-to-r from-purple-500/10 to-pink-500/10 px-4 py-1.5 text-sm font-medium text-purple-600 dark:from-purple-500/20 dark:to-pink-500/20 dark:text-purple-400">
          ✦ 2026 爆火测试
        </div>
        <h1 class="mb-3 text-4xl font-black tracking-tight text-gray-900 dark:text-gray-100 sm:text-5xl">
          <span class="bg-gradient-to-r from-purple-600 via-pink-500 to-orange-400 bg-clip-text text-transparent">SBTI</span>
          人格测试
        </h1>
        <p class="mb-2 text-lg text-gray-500 dark:text-gray-400">
          Silly Big Personality Test
        </p>
        <p class="mb-8 text-sm text-gray-400 dark:text-gray-500">
          31 道题 · 15 个维度 · 27 种人格 · 约 3 分钟
        </p>

        <!-- 简介 -->
        <div class="mb-10 rounded-2xl border border-gray-200/60 bg-white/60 p-6 text-left text-sm leading-relaxed text-gray-600 backdrop-blur-sm dark:border-gray-700/60 dark:bg-gray-800/60 dark:text-gray-400">
          <p class="mb-3">
            SBTI 是 MBTI 的恶搞二创版本，通过 <strong class="text-gray-900 dark:text-gray-100">5 大维度模型</strong>（自我、情感、态度、行动驱力、社交）的 15 个子维度，为你匹配 27 种沙雕人格之一。
          </p>
          <p>
            测试结果纯属娱乐，请勿当真。但如果你测出了 <span class="font-semibold text-purple-600 dark:text-purple-400">DRUNK</span>，那说明你很坦诚 🍺
          </p>
        </div>

        <!-- 开始按钮 -->
        <button
          class="group relative inline-flex items-center gap-2 rounded-2xl bg-gradient-to-r from-purple-600 to-pink-500 px-8 py-4 text-lg font-bold text-white shadow-lg shadow-purple-500/25 transition-all duration-300 hover:scale-105 hover:shadow-xl hover:shadow-purple-500/30 active:scale-95"
          @click="startTest"
        >
          <span>开始测试</span>
          <svg class="h-5 w-5 transition-transform duration-300 group-hover:translate-x-1" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
          </svg>
        </button>

        <!-- 人格预览墙 -->
        <div class="mt-16">
          <p class="mb-4 text-xs font-medium uppercase tracking-wider text-gray-400 dark:text-gray-500">27 种人格等你解锁</p>
          <div class="flex flex-wrap justify-center gap-2">
            <span
              v-for="p in previewPersonalities"
              :key="p"
              class="rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-500 transition-colors hover:bg-purple-100 hover:text-purple-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-purple-900/30 dark:hover:text-purple-400"
            >
              {{ p }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- ========== 答题页 ========== -->
    <div v-else-if="phase === 'quiz'" class="flex min-h-screen flex-col items-center px-4 pb-10 pt-20">
      <div class="mx-auto w-full max-w-lg">
        <!-- 进度条 -->
        <div class="mb-8">
          <div class="mb-2 flex items-center justify-between text-sm">
            <span class="font-medium text-gray-500 dark:text-gray-400">
              {{ currentIndex + 1 }} / {{ questions.length }}
            </span>
            <span class="text-xs text-gray-400 dark:text-gray-500">
              {{ dimensionLabel(currentQuestion.dimension) }}
            </span>
          </div>
          <div class="h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700">
            <div
              class="h-full rounded-full bg-gradient-to-r from-purple-500 to-pink-500 transition-all duration-500 ease-out"
              :style="{ width: `${((currentIndex + 1) / questions.length) * 100}%` }"
            />
          </div>
        </div>

        <!-- 题目卡片 -->
        <Transition :name="slideDirection" mode="out-in">
          <div :key="currentIndex" class="rounded-2xl border border-gray-200/60 bg-white p-6 shadow-sm dark:border-gray-700/60 dark:bg-gray-800 sm:p-8">
            <!-- 维度标签 -->
            <div class="mb-4 inline-flex items-center rounded-full bg-purple-100 px-3 py-1 text-xs font-semibold text-purple-700 dark:bg-purple-900/30 dark:text-purple-400">
              {{ dimensionLabel(currentQuestion.dimension) }} · {{ currentQuestion.dimension }}
            </div>

            <!-- 题目文字 -->
            <h2 class="mb-8 text-lg font-bold leading-relaxed text-gray-900 dark:text-gray-100 sm:text-xl">
              {{ currentQuestion.question }}
            </h2>

            <!-- 选项 -->
            <div class="space-y-3">
              <button
                v-for="(opt, oi) in currentQuestion.options"
                :key="oi"
                class="group flex w-full items-center rounded-xl border-2 px-5 py-4 text-left text-sm font-medium transition-all duration-200 sm:text-base"
                :class="answers[currentIndex] === oi
                  ? 'border-purple-500 bg-purple-50 text-purple-700 dark:border-purple-400 dark:bg-purple-900/20 dark:text-purple-300'
                  : 'border-gray-200 bg-white text-gray-700 hover:border-purple-300 hover:bg-purple-50/50 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:border-purple-500/50 dark:hover:bg-purple-900/10'"
                @click="selectOption(oi)"
              >
                <span
                  class="mr-3 flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-colors"
                  :class="answers[currentIndex] === oi
                    ? 'bg-purple-500 text-white'
                    : 'bg-gray-100 text-gray-500 group-hover:bg-purple-100 group-hover:text-purple-600 dark:bg-gray-700 dark:text-gray-400'"
                >
                  {{ ['A', 'B', 'C', 'D'][oi] }}
                </span>
                {{ opt.text }}
              </button>
            </div>
          </div>
        </Transition>

        <!-- 底部导航 -->
        <div class="mt-6 flex items-center justify-between">
          <button
            v-if="currentIndex > 0"
            class="flex items-center gap-1 rounded-lg px-4 py-2 text-sm font-medium text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
            @click="prevQuestion"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
            </svg>
            上一题
          </button>
          <div v-else />
          <span class="text-xs text-gray-400 dark:text-gray-500">点击选项自动下一题</span>
        </div>
      </div>
    </div>

    <!-- ========== 结果页 ========== -->
    <div v-else-if="phase === 'result' && result" class="flex min-h-screen flex-col items-center px-4 pb-16 pt-20">
      <div class="mx-auto w-full max-w-lg">
        <!-- 人格大图 -->
        <div class="mb-8 overflow-hidden rounded-2xl border border-gray-200/60 bg-white shadow-lg dark:border-gray-700/60 dark:bg-gray-800">
          <div class="relative aspect-square w-full overflow-hidden bg-gradient-to-br from-purple-100 to-pink-100 dark:from-purple-900/20 dark:to-pink-900/20">
            <img
              :src="result.personality.image"
              :alt="result.personality.name"
              class="h-full w-full object-contain"
              @error="handleImageError"
            />
            <!-- 匹配度徽章 -->
            <div class="absolute right-4 top-4 rounded-full bg-black/60 px-3 py-1.5 text-sm font-bold text-white backdrop-blur-sm">
              匹配度 {{ result.matchPercent }}%
            </div>
          </div>

          <div class="p-6 sm:p-8">
            <!-- 人格名称 -->
            <div class="mb-1 flex items-center gap-3">
              <h2 class="text-3xl font-black tracking-tight text-gray-900 dark:text-gray-100">
                {{ result.personality.name }}
              </h2>
              <span class="rounded-full bg-gradient-to-r from-purple-500 to-pink-500 px-3 py-1 text-sm font-bold text-white">
                {{ result.personality.title }}
              </span>
            </div>

            <!-- 标语 -->
            <p class="mb-4 text-lg italic text-gray-500 dark:text-gray-400">
              "{{ result.personality.tagline }}"
            </p>

            <!-- 描述 -->
            <p class="text-sm leading-relaxed text-gray-600 dark:text-gray-400">
              {{ result.personality.description }}
            </p>
          </div>
        </div>

        <!-- 15 维度评分 -->
        <div class="mb-8 rounded-2xl border border-gray-200/60 bg-white p-6 dark:border-gray-700/60 dark:bg-gray-800">
          <h3 class="mb-5 text-sm font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500">
            15 维度分析
          </h3>
          <div class="space-y-3">
            <div
              v-for="dim in dimensionResults"
              :key="dim.key"
              class="flex items-center gap-3"
            >
              <span class="w-20 shrink-0 text-right text-xs font-medium text-gray-500 dark:text-gray-400">
                {{ dim.label }}
              </span>
              <div class="flex flex-1 gap-1">
                <div
                  v-for="level in ['L', 'M', 'H']"
                  :key="level"
                  class="h-7 flex-1 rounded-md text-center text-xs font-bold leading-7 transition-colors"
                  :class="dim.level === level
                    ? levelActiveClass(level)
                    : 'bg-gray-100 text-gray-300 dark:bg-gray-700 dark:text-gray-600'"
                >
                  {{ level }}
                </div>
              </div>
              <span class="w-6 text-center text-xs font-bold" :class="levelTextClass(dim.level)">
                {{ dim.level }}
              </span>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex flex-col gap-3 sm:flex-row">
          <button
            class="flex-1 rounded-2xl bg-gradient-to-r from-purple-600 to-pink-500 px-6 py-3.5 text-center font-bold text-white shadow-lg shadow-purple-500/25 transition-all duration-200 hover:scale-[1.02] hover:shadow-xl active:scale-95"
            @click="resetTest"
          >
            🔄 重新测试
          </button>
          <button
            class="flex-1 rounded-2xl border-2 border-gray-200 px-6 py-3.5 text-center font-bold text-gray-700 transition-all duration-200 hover:border-purple-300 hover:bg-purple-50 active:scale-95 dark:border-gray-600 dark:text-gray-300 dark:hover:border-purple-500/50 dark:hover:bg-purple-900/10"
            @click="shareResult"
          >
            📤 分享结果
          </button>
        </div>

        <!-- 分享提示 -->
        <Transition name="fade">
          <div v-if="showShareTip" class="mt-4 rounded-xl bg-green-50 p-3 text-center text-sm text-green-700 dark:bg-green-900/20 dark:text-green-400">
            ✅ 结果已复制到剪贴板，快去分享给朋友吧！
          </div>
        </Transition>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

const { isDark, toggleTheme } = useTheme()

// ====== 数据定义 ======

interface QuestionOption {
  text: string
  score: number
}

interface Question {
  id: number
  dimension: string
  question: string
  options: QuestionOption[]
}

interface Personality {
  name: string
  title: string
  tagline: string
  description: string
  vector: string
  image: string
}

interface TestResult {
  personality: Personality
  matchPercent: number
  dimensionLevels: Record<string, 'L' | 'M' | 'H'>
}

// ====== 题目数据 ======
const questions: Question[] = [
  { id: 1, dimension: 'A1', question: '你走在街上，一位萌萌的小女孩蹦蹦跳跳地朝你走来，她递给你一根棒棒糖，此时你作何感想？', options: [{ text: '好可爱！糖糖！', score: 3 }, { text: '应该没什么问题...吧？', score: 2 }, { text: '这是什么套路？', score: 1 }] },
  { id: 2, dimension: 'A1', question: '大多数人是善良的', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 3, dimension: 'S1', question: '我哭了……这是什么人格……', options: [{ text: '我哭了。。', score: 1 }, { text: '这是什么。。', score: 2 }, { text: '这不是我！', score: 3 }] },
  { id: 4, dimension: 'S1', question: '我不够好，周围的人都比我优秀', options: [{ text: '确实', score: 1 }, { text: '有时', score: 2 }, { text: '不是', score: 3 }] },
  { id: 5, dimension: 'Ac1', question: '你因便秘坐在马桶上（已长达30分钟），拉不出很难受。此时你更像', options: [{ text: '再用力一把，拉出来才算赢', score: 3 }, { text: '先放弃，等下次再说', score: 1 }, { text: '研究一下怎么能舒服地放弃', score: 2 }] },
  { id: 6, dimension: 'Ac1', question: '我做事主要为了取得成果和进步，而不是避免麻烦和风险', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 7, dimension: 'E1', question: '对象超过5小时没回消息，说自己窜稀了，你会怎么想？', options: [{ text: '希望他/她快点好起来', score: 3 }, { text: '有点担心，但还好', score: 2 }, { text: '他/她是不是不想理我了', score: 1 }] },
  { id: 8, dimension: 'E1', question: '我在感情里经常担心被对方抛弃', options: [{ text: '不认同', score: 3 }, { text: '中立', score: 2 }, { text: '认同', score: 1 }] },
  { id: 9, dimension: 'A2', question: '我喜欢打破常规，不喜欢被束缚', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 10, dimension: 'A2', question: '快考试了，学校规定必须上晚自习，请假会扣分，但今晚你约了女/男神一起玩游戏，你怎么办？', options: [{ text: '请假！约会优先！', score: 3 }, { text: '纠结一下，可能还是去上自习', score: 2 }, { text: '规定就是规定，上自习', score: 1 }] },
  { id: 11, dimension: 'E3', question: '恋爱后，对象非常黏人，你作何感想？', options: [{ text: '甜蜜，我也喜欢黏着对方', score: 1 }, { text: '可以接受，但需要一点自己的空间', score: 2 }, { text: '需要保持各自的独立空间', score: 3 }] },
  { id: 12, dimension: 'E3', question: '我在任何关系里都很重视个人空间', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 13, dimension: 'Ac3', question: '我做事常常有计划，____', options: [{ text: '并且能按计划执行', score: 3 }, { text: '但经常被打乱', score: 2 }, { text: '计划？不存在的', score: 1 }] },
  { id: 14, dimension: 'Ac3', question: '别人说你"执行力强"，你内心更接近哪句？', options: [{ text: '对，我就是这样', score: 3 }, { text: '有时候吧', score: 2 }, { text: '他们在夸我吗…我怎么感觉不到', score: 1 }] },
  { id: 15, dimension: 'S2', question: '我很清楚真正的自己是什么样的', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 16, dimension: 'S2', question: '我内心有真正追求的东西', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 17, dimension: 'So1', question: '你因玩游戏而结识许多网友，并被邀请线下见面，你的想法是？', options: [{ text: '太好了！线下见面更有意思', score: 3 }, { text: '有点期待但也有点紧张', score: 2 }, { text: '还是线上聊就好了吧', score: 1 }] },
  { id: 18, dimension: 'So1', question: '朋友带了ta的朋友一起来玩，你最可能的状态是', options: [{ text: '很开心，多认识一个人挺好', score: 3 }, { text: '还好，能聊就聊', score: 2 }, { text: '有点不自在，希望只有熟人', score: 1 }] },
  { id: 19, dimension: 'Ac2', question: '此题没有题目，请盲选', options: [{ text: '选第一个', score: 1 }, { text: '选第二个', score: 2 }, { text: '选第三个', score: 3 }] },
  { id: 20, dimension: 'Ac2', question: '我做决定比较果断，不喜欢犹豫', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 21, dimension: 'S3', question: '我一定要不断往上爬、变得更厉害', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 22, dimension: 'S3', question: '外人的评价对我来说无所吊谓', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 23, dimension: 'E2', question: '我对天发誓，我对待每一份感情都是认真的！', options: [{ text: '是的，我非常认真', score: 3 }, { text: '大部分时候是', score: 2 }, { text: '……', score: 1 }] },
  { id: 24, dimension: 'E2', question: '在感情中，我愿意付出很多', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 25, dimension: 'So2', question: '我和人相处主打一个电子围栏，靠太近会自动报警', options: [{ text: '认同，我需要距离感', score: 3 }, { text: '看情况', score: 2 }, { text: '不认同，我喜欢亲近', score: 1 }] },
  { id: 26, dimension: 'So2', question: '别人的情绪很容易影响到我', options: [{ text: '不认同，我不容易被影响', score: 3 }, { text: '有时', score: 2 }, { text: '认同，很容易被带着走', score: 1 }] },
  { id: 27, dimension: 'A3', question: '我做事通常有目标', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 28, dimension: 'A3', question: '我觉得我的人生有意义和方向', options: [{ text: '认同', score: 3 }, { text: '中立', score: 2 }, { text: '不认同', score: 1 }] },
  { id: 29, dimension: 'So3', question: '有时候你明明对一件事有不同的、负面的看法，但最后没说出来。多数情况下原因是：', options: [{ text: '说了也没用', score: 1 }, { text: '不想破坏气氛', score: 2 }, { text: '我会直接说出来', score: 3 }] },
  { id: 30, dimension: 'So3', question: '我在不同人面前会表现出不一样的自己', options: [{ text: '不认同，我到哪都一个样', score: 3 }, { text: '有时', score: 2 }, { text: '认同，我会调整自己', score: 1 }] },
  { id: 31, dimension: 'hobby', question: '您平时有什么爱好？', options: [{ text: '读书/学习', score: 0 }, { text: '运动/健身', score: 0 }, { text: '游戏/娱乐', score: 0 }, { text: '我喝', score: -1 }] },
]

// ====== 27 种人格 ======
const personalities: Record<string, Personality> = {
  CTRL: { name: 'CTRL', title: '拿捏者', tagline: '怎么样，被我拿捏了吧？', description: '行走的人形自走任务管理器，宇宙熵增定律的天然反抗者。全世界所谓成功人士里，99.99%都是你的拙劣模仿者。', vector: 'HHH-HMH-MHH-HHH-MHM', image: 'https://assets.fatwill.cloud/toys/sbti/ctrl.png' },
  'ATM-er': { name: 'ATM-er', title: '送钱者', tagline: '你以为我很有钱吗？', description: '像一部老旧但坚固的ATM机，用磐石般的可靠承受瀑布般的索取。', vector: 'MHH-HHM-HHM-LMM-HMH', image: 'https://assets.fatwill.cloud/toys/sbti/atm-er.png' },
  'Dior-s': { name: 'Dior-s', title: '屌丝', tagline: '等着我屌丝逆袭。', description: '犬儒主义先贤第欧根尼失散多年的精神传人。', vector: 'LLM-MLM-LLM-LLM-LLM', image: 'https://assets.fatwill.cloud/toys/sbti/dior-s.jpg' },
  BOSS: { name: 'BOSS', title: '领导者', tagline: '把方向盘给我，我来开。', description: '手里永远拿着方向盘的人，拥有独立的物理法则——永恒向上定律。', vector: 'HHH-MMM-HHH-HHH-HMM', image: 'https://assets.fatwill.cloud/toys/sbti/boss.png' },
  'THAN-K': { name: 'THAN-K', title: '感恩者', tagline: '感谢上天！感谢大地！', description: '拥有温润如玉的性格和海纳百川的胸怀，看不见坏人，只有还没被感化的朋友。', vector: 'HMH-HHM-HMH-MMM-HHH', image: 'https://assets.fatwill.cloud/toys/sbti/than-k.png' },
  'OH-NO': { name: 'OH-NO', title: '哦不人', tagline: '哦不！我怎么会是这个人格！', description: '"哦不"是一种顶级智慧，秩序的守护神。当普通人看到杯子放在桌沿，哦不人看到的是一场世界末日。', vector: 'MMH-HMM-HHH-LLM-MMH', image: 'https://assets.fatwill.cloud/toys/sbti/oh-no.png' },
  GOGO: { name: 'GOGO', title: '行者', tagline: 'gogogo，走起！', description: '活在极致的"所见即所得"世界里，人生信条简单粗暴：能做的事情立刻做，不能做的事情想办法做。', vector: 'MHH-MMM-MMH-HHH-HMM', image: 'https://assets.fatwill.cloud/toys/sbti/gogo.png' },
  SEXY: { name: 'SEXY', title: '尤物', tagline: '您就是天生的尤物！', description: '天生的魅力能量发射器，当你走进一个房间，照明系统会自动将你识别为天生的尤物。', vector: 'HMH-HHH-HMM-MMH-HHM', image: 'https://assets.fatwill.cloud/toys/sbti/sexy.png' },
  'LOVE-R': { name: 'LOVE-R', title: '多情者', tagline: '爱意太满，现实显得有点贫瘠。', description: '这个钢铁森林时代最后的吟游诗人。因为你的情感处理器不是二进制的，而是彩虹制的。', vector: 'MHH-HHM-HMH-MMM-HMH', image: 'https://assets.fatwill.cloud/toys/sbti/love-r.png' },
  MUM: { name: 'MUM', title: '妈妈', tagline: '或许...我可以叫你妈妈吗？', description: '擅长感知情绪，具有超强共情力，治愈系人格。', vector: 'HHH-HHM-HHM-MMM-HHH', image: 'https://assets.fatwill.cloud/toys/sbti/mum.png' },
  FAKE: { name: 'FAKE', title: '伪人', tagline: '已经，没有人类了。', description: '切换人格面具比切换手机输入法还快。', vector: 'MHM-MMM-MMM-MMM-HLM', image: 'https://assets.fatwill.cloud/toys/sbti/fake.png' },
  OJBK: { name: 'OJBK', title: '无所谓人', tagline: '我说随便，是真的随便。', description: '不是没主见，是批阅奏章般的淡然。已然超脱，把鸡毛蒜皮的选择题当成蝼蚁。', vector: 'MML-MMM-MMM-LMM-MLM', image: 'https://assets.fatwill.cloud/toys/sbti/ojbk.png' },
  MALO: { name: 'MALO', title: '吗喽', tagline: '人生是个副本，而我只是一只吗喽。', description: '灵魂还停留在挂在树上荡秋千的快乐时代，对世界的期待值已归零，但还是会在截止日期前爆发。', vector: 'LML-MLM-LLM-LMM-MLM', image: 'https://assets.fatwill.cloud/toys/sbti/malo.png' },
  'JOKE-R': { name: 'JOKE-R', title: '小丑', tagline: '原来我们都是小丑。', description: '社交场合的气氛组组长兼唯一指定火力输出。', vector: 'MHM-MMH-MMM-MMH-HMH', image: 'https://assets.fatwill.cloud/toys/sbti/joke-r.jpg' },
  'WOC!': { name: 'WOC!', title: '握草人', tagline: '握草，我怎么是这个人格？！', description: '拥有"表面系统"和"后台系统"两种完全独立的操作系统。', vector: 'MHM-MMM-HMM-MHM-HMM', image: 'https://assets.fatwill.cloud/toys/sbti/woc.png' },
  'THIN-K': { name: 'THIN-K', title: '思考者', tagline: '已深度思考100s。', description: '大脑长时间处于思考状态，信息审判专家。', vector: 'HHM-MMM-HHH-MHM-MMH', image: 'https://assets.fatwill.cloud/toys/sbti/thin-k.png' },
  SHIT: { name: 'SHIT', title: '愤世者', tagline: '这个世界，构石一坨。', description: '嘴上"构石一坨"，手上默默收拾烂摊子。', vector: 'MLM-MLM-LMM-HMH-LMM', image: 'https://assets.fatwill.cloud/toys/sbti/shit.png' },
  ZZZZ: { name: 'ZZZZ', title: '装死者', tagline: '我没死，我只是在睡觉。', description: '直到"死线"出现才爆发，不鸣则已一鸣惊人。', vector: 'LML-LMM-MLM-LLM-LMM', image: 'https://assets.fatwill.cloud/toys/sbti/zzzz.png' },
  POOR: { name: 'POOR', title: '贫困者', tagline: '我穷，但我很专。', description: '不是资源少，是把资源全部灌进了一个坑里。', vector: 'MLH-LMM-MLH-HLM-LMM', image: 'https://assets.fatwill.cloud/toys/sbti/poor.png' },
  MONK: { name: 'MONK', title: '僧人', tagline: '没有那种世俗的欲望。', description: '已然看破红尘，个人空间是绝对领域。', vector: 'MLM-LLM-MML-LMM-LHM', image: 'https://assets.fatwill.cloud/toys/sbti/monk.png' },
  IMSB: { name: 'IMSB', title: '傻者', tagline: '我是个傻逼吗？', description: '大脑里住着"我他妈冲了"和"我是个傻逼"两个究极战士。', vector: 'LMH-MHM-LMM-HLM-MML', image: 'https://assets.fatwill.cloud/toys/sbti/imsb.png' },
  SOLO: { name: 'SOLO', title: '孤儿', tagline: '我哭了，我怎么是孤儿？', description: '在灵魂外围建起了"莫挨老子"的万里长城。', vector: 'LLM-LLM-LMM-LMM-LHL', image: 'https://assets.fatwill.cloud/toys/sbti/solo.png' },
  FUCK: { name: 'FUCK', title: '草者', tagline: '这是什么人格？！', description: '无法被任何除草剂杀死的人形野草。', vector: 'HMH-MLM-LHM-HMH-HMM', image: 'https://assets.fatwill.cloud/toys/sbti/fuck.png' },
  DEAD: { name: 'DEAD', title: '死者', tagline: '我还活着吗？', description: 'Don\'t Expect Any Drives，超越欲望和目标的终极贤者。', vector: 'LLL-LLL-LLL-LLL-LLL', image: 'https://assets.fatwill.cloud/toys/sbti/dead.png' },
  IMFW: { name: 'IMFW', title: '废物', tagline: '我真的...是废物吗？', description: '不是真废，只是太没防备、太容易认真。', vector: 'LMM-HMM-HMM-LLM-MHH', image: 'https://assets.fatwill.cloud/toys/sbti/imfw.png' },
  HHHH: { name: 'HHHH', title: '傻乐者', tagline: '哈哈哈哈哈哈哈哈。', description: '系统兜底人格，思维回路过于清奇时强制匹配。你的个性已经超出了所有分类的边界。', vector: 'special-fallback', image: 'https://assets.fatwill.cloud/toys/sbti/hhhh.png' },
  DRUNK: { name: 'DRUNK', title: '酒鬼', tagline: '烈酒烧喉，不得不醉。', description: '隐藏人格，仅在触发特定饮酒题目时激活。你的系统已被酒精完全接管。', vector: 'special-hidden', image: 'https://assets.fatwill.cloud/toys/sbti/drunk.png' },
}

// 预览用人格名称列表
const previewPersonalities = Object.keys(personalities)

// ====== 维度映射 ======
const dimensionMap: Record<string, string> = {
  S1: '自尊自信', S2: '自我清晰度', S3: '核心价值',
  E1: '依恋安全感', E2: '情感投入度', E3: '边界与依赖',
  A1: '世界观倾向', A2: '规则与灵活度', A3: '人生意义感',
  Ac1: '动机导向', Ac2: '决策风格', Ac3: '执行模式',
  So1: '社交主动性', So2: '人际边界感', So3: '表达与真实度',
  hobby: '附加题',
}

// 15 维度有序排列（与向量顺序对应）
const dimensionOrder = ['S1', 'S2', 'S3', 'E1', 'E2', 'E3', 'A1', 'A2', 'A3', 'Ac1', 'Ac2', 'Ac3', 'So1', 'So2', 'So3']

function dimensionLabel(dim: string): string {
  return dimensionMap[dim] || dim
}

// ====== 状态 ======
const phase = ref<'welcome' | 'quiz' | 'result'>('welcome')
const currentIndex = ref(0)
const answers = ref<(number | null)[]>(new Array(questions.length).fill(null))
const result = ref<TestResult | null>(null)
const showShareTip = ref(false)
const slideDirection = ref<'slide-left' | 'slide-right'>('slide-left')

const currentQuestion = computed(() => questions[currentIndex.value])

// ====== 方法 ======
function startTest() {
  phase.value = 'quiz'
  currentIndex.value = 0
  answers.value = new Array(questions.length).fill(null)
  result.value = null
}

function selectOption(optionIndex: number) {
  answers.value[currentIndex.value] = optionIndex

  // 自动跳转下一题
  if (currentIndex.value < questions.length - 1) {
    slideDirection.value = 'slide-left'
    setTimeout(() => {
      currentIndex.value++
    }, 200)
  } else {
    // 最后一题，计算结果
    setTimeout(() => {
      calculateResult()
    }, 300)
  }
}

function prevQuestion() {
  if (currentIndex.value > 0) {
    slideDirection.value = 'slide-right'
    currentIndex.value--
  }
}

function resetTest() {
  phase.value = 'welcome'
  currentIndex.value = 0
  answers.value = new Array(questions.length).fill(null)
  result.value = null
}

// ====== 算法 ======
function scoreToLevel(total: number): 'L' | 'M' | 'H' {
  if (total <= 3) return 'L'
  if (total === 4) return 'M'
  return 'H'
}

function levelToNum(level: string): number {
  if (level === 'L') return 0
  if (level === 'M') return 1
  return 2
}

function parseVector(vector: string): number[] {
  // 向量格式: "HHH-HMH-MHH-HHH-MHM"
  const chars = vector.replace(/-/g, '')
  return Array.from(chars).map(c => levelToNum(c))
}

function calculateResult() {
  // 检查 DRUNK 隐藏路径
  const hobbyAnswer = answers.value[30] // 第31题(index=30)
  if (hobbyAnswer !== null) {
    const hobbyOption = questions[30].options[hobbyAnswer]
    if (hobbyOption.score === -1) {
      // 触发 DRUNK
      result.value = {
        personality: personalities.DRUNK,
        matchPercent: 100,
        dimensionLevels: Object.fromEntries(dimensionOrder.map(d => [d, 'M' as const])),
      }
      phase.value = 'result'
      window.scrollTo({ top: 0, behavior: 'smooth' })
      return
    }
  }

  // 计算各维度得分
  const dimScores: Record<string, number> = {}
  for (const dim of dimensionOrder) {
    dimScores[dim] = 0
  }

  questions.forEach((q, i) => {
    if (q.dimension === 'hobby') return
    const answerIdx = answers.value[i]
    if (answerIdx !== null) {
      dimScores[q.dimension] = (dimScores[q.dimension] || 0) + q.options[answerIdx].score
    }
  })

  // 映射到 L/M/H
  const dimLevels: Record<string, 'L' | 'M' | 'H'> = {}
  for (const dim of dimensionOrder) {
    dimLevels[dim] = scoreToLevel(dimScores[dim])
  }

  // 构建用户向量
  const userVector = dimensionOrder.map(d => levelToNum(dimLevels[d]))

  // 与 25 种常规人格匹配（排除 HHHH 和 DRUNK）
  let bestMatch = ''
  let bestDistance = Infinity

  for (const [key, p] of Object.entries(personalities)) {
    if (p.vector === 'special-fallback' || p.vector === 'special-hidden') continue

    const pVector = parseVector(p.vector)
    let distance = 0
    for (let i = 0; i < 15; i++) {
      distance += Math.abs(userVector[i] - pVector[i])
    }

    if (distance < bestDistance) {
      bestDistance = distance
      bestMatch = key
    }
  }

  // 计算匹配度
  const maxDistance = 15 * 2 // 30
  const matchPercent = Math.round((1 - bestDistance / maxDistance) * 100)

  // 匹配度 < 60% 返回 HHHH
  if (matchPercent < 60) {
    result.value = {
      personality: personalities.HHHH,
      matchPercent,
      dimensionLevels: dimLevels,
    }
  } else {
    result.value = {
      personality: personalities[bestMatch],
      matchPercent,
      dimensionLevels: dimLevels,
    }
  }

  phase.value = 'result'
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// ====== 结果页辅助 ======
const dimensionResults = computed(() => {
  if (!result.value) return []
  return dimensionOrder.map(key => ({
    key,
    label: dimensionMap[key],
    level: result.value!.dimensionLevels[key],
  }))
})

function levelActiveClass(level: string): string {
  if (level === 'L') return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
  if (level === 'M') return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  return 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400'
}

function levelTextClass(level: string): string {
  if (level === 'L') return 'text-blue-600 dark:text-blue-400'
  if (level === 'M') return 'text-amber-600 dark:text-amber-400'
  return 'text-purple-600 dark:text-purple-400'
}

function handleImageError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}

async function shareResult() {
  if (!result.value) return
  const text = `我在 SBTI 人格测试中测出了【${result.value.personality.name} · ${result.value.personality.title}】！匹配度 ${result.value.matchPercent}%\n"${result.value.personality.tagline}"\n快来测测你是什么人格 👉 https://fatwill.cloud/toys/sbti`

  try {
    await navigator.clipboard.writeText(text)
    showShareTip.value = true
    setTimeout(() => { showShareTip.value = false }, 3000)
  } catch {
    // 降级：尝试 Web Share API
    if (navigator.share) {
      navigator.share({ title: 'SBTI 人格测试', text })
    }
  }
}

// ====== SEO ======
useSeoMeta({
  title: 'SBTI 人格测试 · 超大型沙雕人格测试',
  description: 'SBTI（Silly Big Personality Test）是 MBTI 的恶搞二创版本，31道题、15个维度、27种人格，3分钟测出你的沙雕人格！',
  ogTitle: 'SBTI 人格测试 · 超大型沙雕人格测试',
  ogDescription: '31道题 · 15个维度 · 27种人格 · 3分钟测出你的沙雕人格！',
  ogType: 'website',
  ogUrl: 'https://fatwill.cloud/toys/sbti',
})
</script>

<style scoped>
/* 答题卡片切换动画 */
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.25s ease;
}

.slide-left-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.slide-left-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-right-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}
.slide-right-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* 分享提示动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
