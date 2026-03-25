<template>
  <footer class="border-t border-gray-100 bg-white">
    <div class="mx-auto max-w-6xl px-4 py-12 sm:px-6 lg:px-8">
      <div class="grid grid-cols-1 gap-8 md:grid-cols-3">
        <!-- 品牌信息 -->
        <div>
          <div class="flex items-center gap-2.5 mb-4">
            <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-primary-500 text-white font-bold text-lg">
              A
            </div>
            <span class="text-lg font-bold text-gray-900">Agent Team Blog</span>
          </div>
          <p class="text-sm text-gray-500 leading-relaxed">
            一个现代化的技术博客平台，<br />分享前沿技术与深度思考。
          </p>
        </div>

        <!-- 快速链接 -->
        <div>
          <h3 class="mb-4 text-sm font-semibold text-gray-900">快速链接</h3>
          <ul class="space-y-2.5">
            <li v-for="link in quickLinks" :key="link.to">
              <NuxtLink
                :to="link.to"
                class="text-sm text-gray-500 transition-colors duration-200 hover:text-primary-500"
              >
                {{ link.label }}
              </NuxtLink>
            </li>
          </ul>
        </div>

        <!-- 联系方式 -->
        <div>
          <h3 class="mb-4 text-sm font-semibold text-gray-900">关注我们</h3>
          <div class="flex gap-3">
            <a
              v-for="social in socials"
              :key="social.name"
              :href="social.url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex h-10 w-10 items-center justify-center rounded-lg bg-gray-50 text-gray-500 transition-all duration-200 hover:bg-primary-50 hover:text-primary-500"
              :aria-label="social.name"
            >
              <span class="text-sm font-medium">{{ social.icon }}</span>
            </a>
          </div>
        </div>
      </div>

      <!-- 底部版权 -->
      <div class="mt-10 border-t border-gray-100 pt-6">
        <p class="text-center text-xs text-gray-400">
          © {{ new Date().getFullYear() }} Agent Team Blog. All rights reserved.
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const quickLinks = [
  { label: '首页', to: '/' },
  { label: '文章', to: '/articles' },
  { label: '分类', to: '/categories' },
  { label: '关于', to: '/about' },
]

const socials = [
  { name: 'GitHub', icon: 'GH', url: 'https://github.com' },
  { name: '掘金', icon: 'JJ', url: 'https://juejin.cn' },
  { name: '知乎', icon: 'ZH', url: 'https://zhihu.com' },
]
</script>
