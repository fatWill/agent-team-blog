<template>
  <div>
    <!-- Hero 区域 -->
    <BusinessHeroSection />

    <!-- 主内容区 -->
    <div class="mx-auto max-w-6xl px-4 py-10 sm:px-6 lg:px-8">
      <!-- 分类筛选 -->
      <div class="mb-8">
        <h2 class="mb-4 text-xl font-bold text-gray-900">探索文章</h2>
        <BusinessCategoryFilter
          :categories="allCategories"
          :active-slug="activeCategory"
          @select="handleCategorySelect"
        />
      </div>

      <!-- 文章列表 + 侧边栏 -->
      <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
        <!-- 文章列表 -->
        <div class="lg:col-span-2">
          <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <BusinessArticleCard
              v-for="article in filteredArticles"
              :key="article.id"
              :article="article"
            />
          </div>

          <!-- 空状态 -->
          <div v-if="filteredArticles.length === 0" class="py-20 text-center">
            <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100">
              <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m5.231 13.481L15 17.25m-4.5-15H5.625c-.621 0-1.125.504-1.125 1.125v16.5c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9zm3.75 11.625a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
              </svg>
            </div>
            <p class="text-sm text-gray-500">该分类下暂无文章</p>
          </div>
        </div>

        <!-- 侧边栏 -->
        <div class="hidden lg:block">
          <BusinessSideBar
            :tags="tags"
            :recommended-articles="recommendedArticles"
            :stats="blogStats"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { articles, categories, tags } = useMockData()

const activeCategory = ref('all')

const allCategories = computed(() => [
  { id: 0, name: '全部', slug: 'all', articleCount: articles.length },
  ...categories,
])

const filteredArticles = computed(() => {
  if (activeCategory.value === 'all') return articles
  return articles.filter((a) => a.category.slug === activeCategory.value)
})

const recommendedArticles = computed(() =>
  [...articles].sort((a, b) => b.viewCount - a.viewCount).slice(0, 5),
)

const blogStats = computed(() => ({
  articles: articles.length,
  categories: categories.length,
  tags: tags.length,
}))

function handleCategorySelect(slug: string) {
  activeCategory.value = slug
}

useHead({
  title: 'Agent Team Blog - 技术探索，无限可能',
  meta: [
    { name: 'description', content: '一个现代化的技术博客平台，分享前沿技术趋势、深度实践经验与独到的工程思考。' },
  ],
})
</script>
