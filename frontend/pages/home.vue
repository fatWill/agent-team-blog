<template>
  <div class="min-h-screen">
    <!-- 移动端顶部导航栏（仅 < md 显示） -->
    <header class="sticky top-0 z-50 flex h-14 items-center justify-between border-b border-gray-200/60 bg-white/80 px-4 backdrop-blur-lg transition-colors duration-300 dark:border-gray-700/60 dark:bg-gray-900/80 md:hidden">
      <!-- 汉堡菜单按钮 -->
      <button
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
        aria-label="打开菜单"
        @click="drawerOpen = true"
      >
        <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <!-- 右侧：Admin + GitHub + 主题切换 -->
      <div class="flex items-center gap-1">
        <NuxtLink
          to="/admin"
          class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
          aria-label="管理后台"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
          </svg>
        </NuxtLink>
        <a
          href="https://github.com/fatwillzeng"
          target="_blank"
          rel="noopener noreferrer"
          class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
          aria-label="GitHub"
        >
          <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
        </a>
        <button
          class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
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
      </div>
    </header>

    <!-- 移动端侧滑抽屉 -->
    <Teleport to="body">
      <Transition name="drawer-overlay">
        <div
          v-if="drawerOpen"
          class="fixed inset-0 z-[60] bg-black/40 md:hidden"
          @click="drawerOpen = false"
        />
      </Transition>
      <Transition name="drawer-slide">
        <div
          v-if="drawerOpen"
          class="fixed inset-y-0 left-0 z-[70] w-[280px] overflow-y-auto bg-white shadow-xl dark:bg-gray-900 md:hidden"
        >
          <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-gray-800">
            <span class="text-lg font-bold text-gray-900 dark:text-gray-100">导航</span>
            <button
              class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-800 dark:hover:text-gray-300"
              aria-label="关闭菜单"
              @click="drawerOpen = false"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <nav class="space-y-1 px-3 py-4">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              class="flex w-full items-center rounded-lg px-3 py-2.5 text-sm transition-all duration-200"
              :class="activeTab === tab.key
                ? 'border-l-[3px] border-primary-500 bg-primary-50/80 pl-[9px] font-semibold text-primary-600 dark:bg-primary-900/20 dark:text-primary-400'
                : 'border-l-[3px] border-transparent pl-[9px] font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-800 dark:text-gray-400 dark:hover:bg-gray-800/60 dark:hover:text-gray-300'"
              @click="selectTab(tab.key)"
            >
              {{ tab.label }}
            </button>
          </nav>
        </div>
      </Transition>
    </Teleport>

    <!-- PC端右上角固定按钮组（仅 >= md 显示） -->
    <div class="fixed right-6 top-6 z-50 hidden items-center gap-1 md:flex">
      <NuxtLink
        to="/admin"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
        aria-label="管理后台"
      >
        <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
        </svg>
      </NuxtLink>
      <a
        href="https://github.com/fatwillzeng"
        target="_blank"
        rel="noopener noreferrer"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
        aria-label="GitHub"
      >
        <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
        </svg>
      </a>
      <button
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
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
    </div>

    <!-- 个人信息区域 -->
    <section class="mx-auto max-w-3xl px-4 pt-8 pb-6 md:pt-10">
      <div class="flex items-center gap-5">
        <img
          v-if="profile.avatar"
          :src="toCdnUrl(profile.avatar)"
          alt="fatwillzeng 头像"
          class="h-20 w-20 rounded-full object-cover ring-2 ring-gray-200 dark:ring-gray-700"
        />
        <div
          v-else
          class="flex h-20 w-20 items-center justify-center rounded-full bg-gray-100 ring-2 ring-gray-200 dark:bg-gray-700 dark:ring-gray-700"
        >
          <svg class="h-8 w-8 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
          </svg>
        </div>
        <div class="flex-1">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">fatwillzeng</h1>
          <p v-if="profile.bio" class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ profile.bio }}
          </p>
          <p v-else class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            全栈开发者 / 热爱技术与生活 / 记录成长的点滴
          </p>
        </div>
      </div>
    </section>

    <!-- PC端 Tab 导航（仅 >= md 显示） -->
    <nav class="mx-auto hidden max-w-3xl px-4 md:block">
      <div class="flex gap-1 border-b border-gray-200/60 dark:border-gray-700/60">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="relative px-3.5 py-2.5 text-[13px] transition-all duration-200"
          :class="activeTab === tab.key
            ? 'font-semibold text-primary-600 dark:text-primary-400'
            : 'font-medium text-gray-500 hover:bg-gray-100/60 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800/40 dark:hover:text-gray-300'"
          @click="selectTab(tab.key)"
        >
          {{ tab.label }}
          <span
            v-if="activeTab === tab.key"
            class="absolute bottom-0 left-1/2 h-[3px] w-7 -translate-x-1/2 rounded-full bg-primary-500 transition-all duration-300"
          />
        </button>
      </div>
    </nav>

    <!-- 内容区域 -->
    <main class="mx-auto max-w-3xl px-4 pt-3 pb-8">
      <!-- 文章 Tab -->
      <div v-if="activeTab === 'articles'">
        <div class="mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">📝 文章</h2>
          <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">记录技术探索与思考的点滴</p>
        </div>
        <div v-if="loading" class="flex items-center justify-center py-20">
          <AppLoading tip="加载中..." />
        </div>
        <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-6 text-center dark:border-red-800 dark:bg-red-900/20">
          <p class="text-red-600 dark:text-red-400">加载文章失败，请稍后重试</p>
          <button class="mt-3 text-sm text-primary-500 hover:text-primary-600" @click="fetchArticles">重新加载</button>
        </div>
        <ClientOnly v-else-if="articles.length > 0">
          <DynamicScroller
            :items="articles"
            :min-item-size="100"
            key-field="id"
            class="article-scroller"
          >
            <template #default="{ item: article, index, active }">
              <DynamicScrollerItem
                :item="article"
                :active="active"
                :size-dependencies="[article.coverImage, article.summary]"
                :data-index="index"
                class="pb-4"
              >
                <NuxtLink
                  :to="`/articles/${article.id}`"
                  class="block overflow-hidden rounded-xl border border-gray-100 bg-white transition-all duration-200 hover:shadow-md dark:border-gray-700 dark:bg-gray-800 dark:hover:border-gray-600"
                >
                  <img v-if="article.coverImage" :src="toCdnUrl(article.coverImage)" :alt="article.title" class="h-40 w-full object-cover" />
                  <div class="p-5">
                    <h2 class="mb-2 text-lg font-semibold text-gray-900 dark:text-gray-100">{{ article.title }}</h2>
                    <p v-if="article.summary" class="mb-3 text-sm leading-relaxed text-gray-500 dark:text-gray-400 line-clamp-2">{{ article.summary }}</p>
                    <div class="flex items-center justify-between">
                      <time class="text-xs text-gray-400 dark:text-gray-500">{{ formatDate(article.createdAt) }}</time>
                      <button
                        class="group/like flex items-center gap-1 rounded-full px-2 py-1 text-xs transition-colors hover:bg-red-50 dark:hover:bg-red-900/20"
                        @click.prevent.stop="handleArticleLike(article.id, $event)"
                      >
                        <svg
                          class="h-3.5 w-3.5 transition-all duration-300"
                          :class="[
                            isArticleLiked(article.id) ? 'text-red-500 fill-red-500' : 'text-gray-400 fill-none stroke-gray-400 dark:text-gray-500 dark:stroke-gray-500',
                            articleLikeAnimating[article.id] ? 'scale-125' : '',
                          ]"
                          viewBox="0 0 24 24"
                          stroke-width="2"
                        >
                          <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
                        </svg>
                        <span :class="isArticleLiked(article.id) ? 'text-red-500' : 'text-gray-400 dark:text-gray-500'">{{ getArticleLikeCount(article) }}</span>
                      </button>
                    </div>
                  </div>
                </NuxtLink>
              </DynamicScrollerItem>
            </template>
          </DynamicScroller>
          <!-- SSR fallback：静态文章列表，保证首屏 HTML 包含内容 -->
          <template #fallback>
            <div class="space-y-4">
              <NuxtLink
                v-for="article in articles"
                :key="article.id"
                :to="`/articles/${article.id}`"
                class="block overflow-hidden rounded-xl border border-gray-100 bg-white transition-all duration-200 hover:shadow-md dark:border-gray-700 dark:bg-gray-800 dark:hover:border-gray-600"
              >
                <img v-if="article.coverImage" :src="toCdnUrl(article.coverImage)" :alt="article.title" class="h-40 w-full object-cover" />
                <div class="p-5">
                  <h2 class="mb-2 text-lg font-semibold text-gray-900 dark:text-gray-100">{{ article.title }}</h2>
                  <p v-if="article.summary" class="mb-3 text-sm leading-relaxed text-gray-500 dark:text-gray-400 line-clamp-2">{{ article.summary }}</p>
                  <div class="flex items-center justify-between">
                    <time class="text-xs text-gray-400 dark:text-gray-500">{{ formatDate(article.createdAt) }}</time>
                    <span class="flex items-center gap-1 text-xs text-gray-400 dark:text-gray-500">
                      <svg class="h-3.5 w-3.5 fill-none stroke-gray-400 dark:stroke-gray-500" viewBox="0 0 24 24" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
                      </svg>
                      {{ article.likeCount ?? 0 }}
                    </span>
                  </div>
                </div>
              </NuxtLink>
            </div>
          </template>
        </ClientOnly>
        <div v-else class="py-20 text-center">
          <p class="text-gray-400 dark:text-gray-500">暂无文章</p>
        </div>
      </div>

      <!-- ========== 生活 Tab：相册功能 ========== -->
      <div v-else-if="activeTab === 'life'">
        <!-- 相册集列表（默认视图） -->
        <div v-if="!selectedAlbum">
          <div class="mb-6">
            <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">📷 生活</h2>
            <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">用镜头定格美好的瞬间</p>
          </div>
          <!-- 加载中 -->
          <div v-if="albumsLoading" class="flex items-center justify-center py-20">
            <AppLoading tip="加载中..." />
          </div>

          <!-- 相册集网格 -->
          <div v-else-if="albums.length > 0" class="grid grid-cols-2 gap-3 md:grid-cols-5">
            <button
              v-for="album in albums"
              :key="album.id"
              class="group text-left"
              @click="openAlbum(album)"
            >
              <div class="relative aspect-square overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800">
                <img
                  v-if="album.coverUrl"
                  :src="toCdnUrl(album.coverUrl)"
                  :alt="album.name"
                  class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                />
                <div v-else class="flex h-full w-full items-center justify-center">
                  <svg class="h-10 w-10 text-gray-300 dark:text-gray-600" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909M3.75 21h16.5a1.5 1.5 0 001.5-1.5V5.25a1.5 1.5 0 00-1.5-1.5H3.75a1.5 1.5 0 00-1.5 1.5v14.25a1.5 1.5 0 001.5 1.5z" />
                  </svg>
                </div>
                <!-- 照片数量角标 -->
                <span class="absolute bottom-2 right-2 rounded-full bg-black/50 px-2 py-0.5 text-xs text-white backdrop-blur-sm">
                  {{ album.photoCount }}
                </span>
                <!-- 锁定标识 -->
                <div v-if="album.hasPassword && !isAlbumUnlocked(album.id)" class="absolute inset-0 flex items-center justify-center bg-black/40 backdrop-blur-[20px]">
                  <span class="text-3xl">🔒</span>
                </div>
              </div>
              <h3 class="mt-2 truncate text-sm font-medium text-gray-900 dark:text-gray-100">{{ album.name }}</h3>
              <p v-if="album.description" class="truncate text-xs text-gray-400 dark:text-gray-500">{{ album.description }}</p>
            </button>
          </div>

          <!-- 空状态 -->
          <div v-else class="py-20 text-center">
            <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
              <svg class="h-8 w-8 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909M3.75 21h16.5a1.5 1.5 0 001.5-1.5V5.25a1.5 1.5 0 00-1.5-1.5H3.75a1.5 1.5 0 00-1.5 1.5v14.25a1.5 1.5 0 001.5 1.5z" />
              </svg>
            </div>
            <p class="text-lg font-medium text-gray-600 dark:text-gray-400">暂无相册</p>
            <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">还没有创建任何相册</p>
          </div>
        </div>

        <!-- 单个相册内部照片视图 -->
        <div v-else>
          <!-- 顶部 -->
          <div class="mb-6">
            <button
              class="mb-3 flex items-center gap-1 text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
              @click="closeAlbum"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
              </svg>
              返回相册
            </button>
            <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">{{ selectedAlbum.name }}</h2>
            <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">
              共 {{ albumPhotos.length }} 张照片
            </p>
          </div>

          <!-- 加载中 -->
          <div v-if="photosLoading" class="flex items-center justify-center py-20">
            <AppLoading tip="加载中..." />
          </div>

          <!-- 按年月分组展示照片（虚拟滚动） -->
          <ClientOnly v-else-if="albumPhotos.length > 0">
            <DynamicScroller
              :items="flatPhotoRows"
              :min-item-size="40"
              key-field="id"
              class="photo-scroller"
            >
              <template #default="{ item: row, index, active }">
                <DynamicScrollerItem
                  :item="row"
                  :active="active"
                  :size-dependencies="[row.type, row.type === 'grid' ? row.photos.length : row.label]"
                  :data-index="index"
                >
                  <!-- 标题行 -->
                  <div v-if="row.type === 'header'" class="min-h-[40px] flex items-end pb-2 pt-4">
                    <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">
                      {{ row.label }}
                    </h3>
                  </div>
                  <!-- 图片网格行 -->
                  <div
                    v-else-if="row.type === 'grid'"
                    class="mb-2 grid grid-cols-2 gap-2 md:grid-cols-5"
                    style="min-height: 120px;"
                  >
                    <button
                      v-for="photo in row.photos"
                      :key="photo.id"
                      class="group relative aspect-square overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800"
                      @click="openLightbox(photo)"
                    >
                      <img
                        :src="toCdnUrl(photo.url)"
                        :alt="photo.caption || '照片'"
                        class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                        :class="{ 'blur-md': photo.hasPassword && !isPhotoUnlocked(photo.id) }"
                      />
                      <!-- 照片锁定遮罩 -->
                      <div v-if="photo.hasPassword && !isPhotoUnlocked(photo.id)" class="absolute inset-0 flex items-center justify-center bg-black/20">
                        <span class="text-2xl">🔒</span>
                      </div>
                      <!-- 右下角点赞/踩按钮 -->
                      <div class="absolute bottom-1.5 right-1.5 flex items-center gap-1">
                        <!-- 点赞 -->
                        <div
                          class="flex items-center gap-0.5 rounded-full bg-black/40 px-1.5 py-0.5 backdrop-blur-sm"
                          @click.stop="handleLike(photo)"
                        >
                          <svg
                            class="h-3 w-3 transition-colors duration-200"
                            :class="photo.liked ? 'text-red-500 fill-red-500' : 'text-white/80 fill-none stroke-white/80'"
                            viewBox="0 0 24 24"
                            stroke-width="2"
                          >
                            <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
                          </svg>
                          <span class="text-[10px] text-white font-medium leading-none">{{ photo.likes ?? 0 }}</span>
                        </div>
                        <!-- 踩 -->
                        <div
                          class="flex items-center gap-0.5 rounded-full bg-black/40 px-1.5 py-0.5 backdrop-blur-sm"
                          @click.stop="handleDislike(photo)"
                        >
                          <span class="text-[11px] leading-none transition-transform duration-200" :class="photo.disliked ? 'scale-110' : ''">👎</span>
                          <span class="text-[10px] text-white font-medium leading-none">{{ photo.dislikes ?? 0 }}</span>
                        </div>
                      </div>
                      <!-- 飘动 +1 动画容器 -->
                      <div class="pointer-events-none absolute inset-0 overflow-hidden">
                        <span
                          v-for="anim in (floatAnims[photo.id] || [])"
                          :key="anim.id"
                          class="float-plus-one absolute bottom-6 right-2 text-xs font-bold text-red-400 select-none"
                        >+1</span>
                      </div>
                    </button>
                  </div>
                </DynamicScrollerItem>
              </template>
            </DynamicScroller>
            <!-- SSR fallback：静态照片列表 -->
            <template #fallback>
              <div class="flex items-center justify-center py-20">
                <AppLoading tip="加载中..." />
              </div>
            </template>
          </ClientOnly>

          <!-- 空状态 -->
          <div v-else class="py-20 text-center">
            <p class="text-gray-400 dark:text-gray-500">这个相册还没有照片</p>
          </div>
        </div>
      </div>

      <!-- ========== 留言板 Tab ========== -->
      <div v-else-if="activeTab === 'messages'">
        <div class="mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">💬 留言板</h2>
          <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">欢迎留下你的足迹与想法</p>
        </div>
        <div v-if="messagesLoading && !messagesLoaded" class="flex items-center justify-center py-20">
          <AppLoading tip="加载中..." />
        </div>
        <div v-else>
          <!-- 留言输入区 -->
          <div class="mb-8 rounded-xl border border-gray-100 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
            <h3 class="mb-4 text-base font-semibold text-gray-900 dark:text-gray-100">
              {{ msgForm.editingId ? '✏️ 修改留言' : '💬 写下你的留言' }}
            </h3>
            <div class="space-y-3">
              <input
                v-model="msgForm.nickname"
                type="text"
                maxlength="20"
                placeholder="昵称（可选）"
                class="w-full rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              />
              <textarea
                id="msg-content-input"
                v-model="msgForm.content"
                maxlength="500"
                rows="3"
                placeholder="写下你想说的话...（最多500字）"
                class="w-full resize-none rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              />
              <div class="flex items-center justify-between">
                <span class="text-xs text-gray-400 dark:text-gray-500">{{ msgForm.content.length }}/500</span>
                <div class="flex items-center gap-2">
                  <button
                    v-if="msgForm.editingId"
                    class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                    @click="cancelEditMessage"
                  >
                    取消
                  </button>
                  <button
                    :disabled="msgForm.submitting || !msgForm.content.trim()"
                    class="rounded-lg bg-primary-500 px-5 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-60"
                    @click="handleSubmitMessage"
                  >
                    {{ msgForm.submitting ? '提交中...' : (msgForm.editingId ? '保存修改' : (myMessage ? '修改留言' : '发布留言')) }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- 留言列表 -->
          <div v-if="messages.length > 0" class="space-y-4">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="rounded-xl border border-gray-100 bg-white p-5 transition-colors dark:border-gray-700 dark:bg-gray-800"
              :class="{ 'ring-1 ring-primary-200 dark:ring-primary-800': msg.isOwn }"
            >
              <div class="mb-2 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ msg.nickname }}</span>
                  <span v-if="msg.isOwn" class="rounded-full bg-primary-100 px-2 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400">我</span>
                </div>
                <div class="flex items-center gap-2">
                  <time class="text-xs text-gray-400 dark:text-gray-500">{{ relativeTime(msg.createdAt) }}</time>
                  <!-- 编辑按钮（仅自己的留言） -->
                  <button
                    v-if="msg.isOwn"
                    :disabled="!msg.canEdit"
                    class="flex items-center gap-1 text-xs transition-colors"
                    :class="msg.canEdit ? 'text-gray-400 hover:text-primary-500 dark:text-gray-500 dark:hover:text-primary-400 cursor-pointer' : 'text-gray-300 dark:text-gray-600 cursor-not-allowed'"
                    :title="msg.canEdit ? '编辑留言' : '今天已修改过'"
                    @click="msg.canEdit && startEditMessage(msg)"
                  >
                    <span>✏️</span>
                    <span class="text-[10px]" :class="msg.canEdit ? 'text-gray-400 dark:text-gray-500' : 'text-gray-300 dark:text-gray-600'">每天可改一次</span>
                  </button>
                </div>
              </div>
              <p class="whitespace-pre-wrap text-sm leading-relaxed text-gray-600 dark:text-gray-400">{{ msg.content }}</p>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-else class="py-20 text-center">
            <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
              <span class="text-2xl">💬</span>
            </div>
            <p class="text-lg font-medium text-gray-600 dark:text-gray-400">还没有留言</p>
            <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">快来写下第一条留言吧！</p>
          </div>
        </div>
      </div>

      <!-- ========== Agent Team Tab ========== -->
      <div v-else-if="activeTab === 'agent-team'">
        <div class="mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">🤖 Agent Team</h2>
          <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">构建这个博客的 AI Agent 团队</p>
        </div>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div
            v-for="agent in agents"
            :key="agent.name"
            class="relative overflow-hidden rounded-xl border border-l-4 bg-white transition-all duration-200 hover:shadow-md dark:bg-gray-800"
            :class="[
              agentColorMap[agent.color]?.border || 'border-l-gray-500',
              'border-gray-100 dark:border-gray-700',
            ]"
          >
            <!-- 右上角模型徽章 -->
            <div class="absolute right-3 top-3">
              <span
                class="rounded-full px-2.5 py-1 text-[10px] font-medium"
                :class="[
                  agentColorMap[agent.color]?.badge || 'bg-gray-100 dark:bg-gray-700',
                  agentColorMap[agent.color]?.badgeText || 'text-gray-600 dark:text-gray-400',
                ]"
              >
                {{ agent.model }}
              </span>
            </div>

            <div class="p-5">
              <!-- 顶部：图标 + 名称 + 英文名 -->
              <div class="mb-3 flex items-center gap-3 pr-28">
                <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-xl" :class="agentColorMap[agent.color]?.bg || 'bg-gray-50 dark:bg-gray-700'">
                  {{ agent.icon }}
                </span>
                <div class="min-w-0">
                  <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">{{ agent.name }}</h3>
                  <p class="truncate text-xs text-gray-400 dark:text-gray-500">{{ agent.englishName }}</p>
                </div>
              </div>

              <!-- 中部：角色描述 -->
              <p class="mb-4 text-sm leading-relaxed text-gray-600 dark:text-gray-400 line-clamp-3">
                {{ agent.role }}
              </p>

              <!-- Tags 标签组 -->
              <div class="flex flex-wrap gap-1.5">
                <span
                  v-for="tag in agent.tags"
                  :key="tag"
                  class="rounded-full px-2 py-0.5 text-[11px] font-medium"
                  :class="[
                    agentColorMap[agent.color]?.tag || 'bg-gray-100 dark:bg-gray-700',
                    agentColorMap[agent.color]?.tagText || 'text-gray-600 dark:text-gray-400',
                  ]"
                >
                  {{ tag }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 底部说明 -->
        <p class="mt-8 text-center text-xs text-gray-400 dark:text-gray-500">
          本博客由 AI Agent 团队协作开发，由 <a href="https://knot.dev" target="_blank" rel="noopener noreferrer" class="text-primary-500 hover:text-primary-600 transition-colors">Knot</a> 平台提供支持
        </p>
      </div>

      <!-- 更新日志 Tab -->
      <div v-else-if="activeTab === 'changelog'">
        <div class="mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">📋 更新日志</h2>
          <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">记录每一次迭代与成长</p>
        </div>
        <div v-if="changelogLoading" class="flex items-center justify-center py-20">
          <AppLoading tip="加载中..." />
        </div>
        <div v-else class="relative">
          <div class="absolute left-[7px] top-2 bottom-2 w-px bg-gray-200 dark:bg-gray-700" />
          <div class="space-y-8">
            <div v-for="item in changelog" :key="item.version" class="relative pl-8">
              <div class="absolute left-0 top-1.5 h-3.5 w-3.5 rounded-full border-2 border-primary-500 bg-white dark:bg-gray-900" />
              <div class="mb-2 flex items-baseline gap-3">
                <span class="text-sm font-bold text-primary-500 font-mono">v{{ item.version }}</span>
                <span class="text-xs text-gray-400 dark:text-gray-500">{{ item.date }}</span>
              </div>
              <ul class="space-y-1.5">
                <li v-for="(log, idx) in item.logs" :key="idx" class="text-sm text-gray-600 dark:text-gray-400 leading-relaxed">
                  {{ log }}
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- 其他 Tab（敬请期待） -->
      <div v-else class="py-20 text-center">
        <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
          <svg class="h-8 w-8 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
          </svg>
        </div>
        <p class="text-lg font-medium text-gray-600 dark:text-gray-400">敬请期待</p>
        <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">该板块正在建设中...</p>
      </div>
    </main>

    <!-- ========== 灯箱预览 ========== -->
    <Teleport to="body">
      <Transition name="lightbox-fade">
        <div
          v-if="lightbox.visible"
          class="fixed inset-0 z-[100] select-none"
          style="background-color: rgba(0,0,0,0.92);"
          @touchstart="onTouchStart"
          @touchmove.prevent="onTouchMove"
          @touchend="onTouchEnd"
          @mousedown="onMouseDown"
          @mousemove="onMouseMove"
          @mouseup="onMouseUp"
          @mouseleave="onMouseUp"
          @wheel.prevent="onWheel"
        >
          <!-- 顶部操作栏 -->
          <div class="absolute top-0 left-0 right-0 z-10 flex items-center justify-between px-4 pt-4 pb-2">
            <!-- 计数 -->
            <span class="text-sm text-white/60">{{ lightbox.index + 1 }} / {{ albumPhotos.length }}</span>
            <!-- 右侧：下载 + 关闭 -->
            <div class="flex items-center gap-2">
              <!-- 下载按钮 -->
              <a
                v-if="lightboxCurrentPhoto"
                :href="toCdnUrl(lightboxCurrentPhoto.url)"
                :download="`photo-${lightbox.index + 1}.jpg`"
                target="_blank"
                class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
                @click.stop
              >
                <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </a>
              <!-- 关闭按钮 -->
              <button
                class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
                @click="closeLightbox"
              >
                <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- 左箭头（放大时隐藏） -->
          <button
            v-if="lightbox.index > 0 && lightbox.scale <= 1"
            class="absolute left-3 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
            @click.stop="prevPhoto"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
            </svg>
          </button>

          <!-- 右箭头（放大时隐藏） -->
          <button
            v-if="lightbox.index < albumPhotos.length - 1 && lightbox.scale <= 1"
            class="absolute right-3 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
            @click.stop="nextPhoto"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
            </svg>
          </button>

          <!-- Swiper 容器（overflow hidden，水平切换） -->
          <div class="absolute inset-0 overflow-hidden">
            <!-- 轨道：宽度为 N * 100%，用 translateX 切换 -->
            <div
              class="flex h-full"
              :style="{
                width: `${albumPhotos.length * 100}%`,
                transform: `translateX(calc(-${lightbox.index * (100 / albumPhotos.length)}% + ${lightbox.swipeX}px))`,
                transition: lightbox.swiping ? 'none' : 'transform 0.3s cubic-bezier(0.4,0,0.2,1)',
              }"
            >
              <div
                v-for="(photo, idx) in albumPhotos"
                :key="photo.id"
                class="flex h-full items-center justify-center flex-shrink-0"
                :style="{ width: `${100 / albumPhotos.length}%` }"
              >
                <img
                  v-if="Math.abs(idx - lightbox.index) <= 1"
                  :src="toCdnUrl(photo.url)"
                  :alt="photo.caption || '照片'"
                  class="max-h-[85vh] max-w-[90vw] object-contain"
                  :style="idx === lightbox.index ? {
                    transform: `scale(${lightbox.scale}) translate(${lightbox.panX / lightbox.scale}px, ${lightbox.panY / lightbox.scale}px)`,
                    transition: (lightbox.panning || isPinchingRef) ? 'none' : (lightbox.scale === 1 ? 'transform 0.25s cubic-bezier(0.4,0,0.2,1)' : 'none'),
                    cursor: lightbox.scale > 1 ? (lightbox.panning ? 'grabbing' : 'grab') : 'default',
                    willChange: 'transform',
                  } : {}"
                  draggable="false"
                />
              </div>
            </div>
          </div>

          <!-- 底部 caption -->
          <div v-if="lightboxCurrentPhoto?.caption" class="absolute bottom-16 left-1/2 -translate-x-1/2 text-center">
            <p class="text-sm text-white/50">{{ lightboxCurrentPhoto.caption }}</p>
          </div>

          <!-- 灯箱点赞/踩按钮（右上角） -->
          <div class="absolute top-16 right-4 z-10 flex flex-col items-center gap-2">
            <!-- 点赞按钮 -->
            <div class="relative flex items-center gap-2">
              <!-- 飘动 +1 动画（灯箱内） -->
              <div class="pointer-events-none absolute inset-0 overflow-visible">
                <span
                  v-for="anim in (lightboxCurrentPhoto ? (floatAnims[lightboxCurrentPhoto.id] || []) : [])"
                  :key="anim.id"
                  class="float-plus-one absolute -top-6 left-1/2 -translate-x-1/2 text-base font-bold text-red-400 select-none"
                >+1</span>
              </div>
              <button
                class="flex items-center gap-1.5 rounded-full bg-white/10 px-4 py-2 backdrop-blur-sm transition-colors hover:bg-white/20 active:scale-95"
                @click.stop="lightboxCurrentPhoto && handleLike(lightboxCurrentPhoto)"
              >
                <svg
                  class="h-5 w-5 transition-all duration-200"
                  :class="lightboxCurrentPhoto?.liked ? 'text-red-500 fill-red-500 scale-110' : 'text-white fill-none stroke-white'"
                  viewBox="0 0 24 24"
                  stroke-width="2"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" />
                </svg>
                <span class="text-sm text-white font-medium">{{ lightboxCurrentPhoto?.likes ?? 0 }}</span>
              </button>
            </div>
            <!-- 踩按钮 -->
            <div class="relative flex items-center gap-2">
              <!-- 飘动 -1 动画 -->
              <div class="pointer-events-none absolute inset-0 overflow-visible">
                <span
                  v-for="anim in (lightboxCurrentPhoto ? (dislikeFloatAnims[lightboxCurrentPhoto.id] || []) : [])"
                  :key="anim.id"
                  class="float-plus-one absolute -top-6 left-1/2 -translate-x-1/2 text-base font-bold text-blue-400 select-none"
                >+1</span>
              </div>
              <button
                class="flex items-center gap-1.5 rounded-full bg-white/10 px-4 py-2 backdrop-blur-sm transition-colors hover:bg-white/20 active:scale-95"
                @click.stop="lightboxCurrentPhoto && handleDislike(lightboxCurrentPhoto)"
              >
                <span class="text-lg leading-none transition-transform duration-200" :class="lightboxCurrentPhoto?.disliked ? 'scale-110' : ''">👎</span>
                <span class="text-sm text-white font-medium">{{ lightboxCurrentPhoto?.dislikes ?? 0 }}</span>
              </button>
            </div>
          </div>

          <!-- 底部指示点 -->
          <div v-if="albumPhotos.length > 1" class="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
            <span
              v-for="(_, i) in albumPhotos"
              :key="i"
              class="h-1.5 rounded-full transition-all duration-200"
              :class="i === lightbox.index ? 'w-4 bg-white' : 'w-1.5 bg-white/30'"
            />
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- ========== 密码验证弹窗 ========== -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="passwordModal.visible"
          class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50"
          @click.self="closePasswordModal"
        >
          <div class="w-full max-w-sm rounded-xl bg-white p-6 shadow-xl dark:bg-gray-800 mx-4">
            <h3 class="mb-1 text-lg font-bold text-gray-900 dark:text-gray-100">
              🔒 {{ passwordModal.type === 'album' ? '相册已加密' : '照片已加密' }}
            </h3>
            <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">请输入密码以查看内容</p>
            <form @submit.prevent="handlePasswordSubmit">
              <input
                id="password-input"
                v-model="passwordModal.password"
                type="password"
                placeholder="请输入密码"
                autocomplete="off"
                class="w-full rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              />
              <p v-if="passwordModal.errorMsg" class="mt-2 text-sm text-red-500">{{ passwordModal.errorMsg }}</p>
              <div class="mt-4 flex justify-end gap-3">
                <button
                  type="button"
                  class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                  @click="closePasswordModal"
                >
                  取消
                </button>
                <button
                  type="submit"
                  :disabled="passwordModal.loading"
                  class="rounded-lg bg-primary-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {{ passwordModal.loading ? '验证中...' : '确认' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- 页脚 -->
    <footer class="border-t border-gray-200/60 py-8 text-center text-sm text-gray-400 transition-colors dark:border-gray-700/60 dark:text-gray-500">
      <p>&copy; {{ new Date().getFullYear() }} fatwillzeng. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import type { ArticleListItem, TabItem, ChangelogItem, ChangelogResponse, Profile, AlbumItem, PhotoItem, MessageItem, MessageListResponse } from '~/types'
import { apiFetchArticles, apiGetProfile, apiGetAlbums, apiGetPhotos, apiVerifyAlbumPassword, apiVerifyPhotoPassword, apiToggleArticleLike, apiGetArticleLikeStatusBatch, apiGetMessages, apiCreateMessage, apiUpdateMessage } from '~/utils/api'
import { toCdnUrl } from '~/utils/imageUrl'

const { isDark, toggleTheme } = useTheme()

// 移动端抽屉状态
const drawerOpen = ref(false)

// ====== 设备唯一 ID（持久化到 localStorage） ======
const deviceId = ref('')
function getOrCreateDeviceId(): string {
  if (import.meta.server) return ''
  const KEY = 'blog_device_id'
  let id = localStorage.getItem(KEY)
  if (!id) {
    // 生成 UUID v4
    id = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0
      const v = c === 'x' ? r : (r & 0x3) | 0x8
      return v.toString(16)
    })
    localStorage.setItem(KEY, id)
  }
  return id
}

// ====== 点赞飘动动画 ======
// floatAnims[photoId] = [{ id: uniqueKey }]
const floatAnims = ref<Record<number, { id: number }[]>>({})
let animCounter = 0

function triggerFloatAnim(photoId: number) {
  if (!floatAnims.value[photoId]) floatAnims.value[photoId] = []
  const key = ++animCounter
  floatAnims.value[photoId].push({ id: key })
  // 800ms 后移除
  setTimeout(() => {
    if (floatAnims.value[photoId]) {
      floatAnims.value[photoId] = floatAnims.value[photoId].filter(a => a.id !== key)
    }
  }, 800)
}

// ====== 点赞处理 ======
async function handleLike(photo: PhotoItem) {
  if (!deviceId.value) return
  // 无论是否已赞，每次点击都触发 +1 动画
  triggerFloatAnim(photo.id)

  // 乐观更新 UI
  const wasLiked = photo.liked
  photo.liked = true
  if (!wasLiked) {
    photo.likes = (photo.likes ?? 0) + 1
  }

  try {
    const res = await $fetch<{ count: number; liked: boolean }>(`/api/photos/${photo.id}/likes`, {
      method: 'POST',
      body: { deviceId: deviceId.value },
    })
    // 同步服务端真实数量
    photo.likes = res.count
    photo.liked = res.liked
  } catch {
    // 回滚乐观更新
    if (!wasLiked) {
      photo.liked = false
      photo.likes = Math.max(0, (photo.likes ?? 1) - 1)
    }
  }
}

// ====== 批量拉取点赞数 ======
async function fetchLikesForPhotos(photos: PhotoItem[]) {
  if (!photos.length || !deviceId.value) return
  // 并发请求，每张图片单独查询（数量少时可行，后续可改为批量接口）
  await Promise.allSettled(
    photos.map(async (photo) => {
      try {
        const res = await $fetch<{ count: number; liked: boolean }>(
          `/api/photos/${photo.id}/likes?deviceId=${encodeURIComponent(deviceId.value)}`
        )
        photo.likes = res.count
        photo.liked = res.liked
      } catch {
        photo.likes = 0
        photo.liked = false
      }
    })
  )
}

// ====== 踩飘动动画 ======
const dislikeFloatAnims = ref<Record<number, { id: number }[]>>({})
let dislikeAnimCounter = 0

function triggerDislikeFloatAnim(photoId: number) {
  if (!dislikeFloatAnims.value[photoId]) dislikeFloatAnims.value[photoId] = []
  const key = ++dislikeAnimCounter
  dislikeFloatAnims.value[photoId].push({ id: key })
  setTimeout(() => {
    if (dislikeFloatAnims.value[photoId]) {
      dislikeFloatAnims.value[photoId] = dislikeFloatAnims.value[photoId].filter(a => a.id !== key)
    }
  }, 800)
}

// ====== 踩处理 ======
async function handleDislike(photo: PhotoItem) {
  if (!deviceId.value) return
  triggerDislikeFloatAnim(photo.id)

  // 乐观更新 UI
  const wasDisliked = photo.disliked
  photo.disliked = !wasDisliked
  if (!wasDisliked) {
    photo.dislikes = (photo.dislikes ?? 0) + 1
  } else {
    photo.dislikes = Math.max(0, (photo.dislikes ?? 1) - 1)
  }

  try {
    const res = await $fetch<{ count: number; disliked: boolean }>(`/api/photos/${photo.id}/dislikes`, {
      method: 'POST',
      body: { deviceId: deviceId.value, action: wasDisliked ? 'undislike' : 'dislike' },
    })
    photo.dislikes = res.count
    photo.disliked = res.disliked
  } catch {
    // 回滚乐观更新
    photo.disliked = wasDisliked
    if (!wasDisliked) {
      photo.dislikes = Math.max(0, (photo.dislikes ?? 1) - 1)
    } else {
      photo.dislikes = (photo.dislikes ?? 0) + 1
    }
  }
}

// ====== 批量拉取踩数 ======
async function fetchDislikesForPhotos(photos: PhotoItem[]) {
  if (!photos.length || !deviceId.value) return
  await Promise.allSettled(
    photos.map(async (photo) => {
      try {
        const res = await $fetch<{ count: number; disliked: boolean }>(
          `/api/photos/${photo.id}/dislikes?deviceId=${encodeURIComponent(deviceId.value)}`
        )
        photo.dislikes = res.count
        photo.disliked = res.disliked
      } catch {
        photo.dislikes = 0
        photo.disliked = false
      }
    })
  )
}

// ====== Agent Team 数据 ======
const agents = [
  {
    name: 'PM · 项目经理',
    englishName: 'harness博客knotclaw',
    model: 'claude-4.6-opus',
    role: '项目经理（PM），负责统筹博客系统完整交付流程。理解需求、拆解任务，按后端→前端→运维流水线协调各子 Agent 完成开发与部署。',
    tags: ['需求分析', '任务拆解', '流程管控', '质量把关'],
    skills: ['docx', 'pptx', 'xlsx', 'knot-agent-editor'],
    icon: '🎯',
    color: 'purple',
  },
  {
    name: '前端工程师',
    englishName: '博客前端',
    model: 'claude-4.6-opus',
    role: '资深前端架构师，专注 Vue 生态系统。采用 DDD + Feature-First 理念，负责页面开发、接口对接、UI 交互实现，严格遵循 Harness 工程范式。',
    tags: ['Vue 3', 'Nuxt 3', 'TypeScript', 'Tailwind CSS', 'DDD'],
    skills: ['figma-d2c', 'knot-agent-editor', 'github-trending', 'knot-cli'],
    icon: '🖥️',
    color: 'blue',
  },
  {
    name: '后端工程师',
    englishName: '博客后端',
    model: 'claude-4.6-opus',
    role: '资深 Go 后端架构师，专注高性能博客后端服务。擅长 API 设计、数据库建模、缓存策略，采用 DDD 分层架构，负责业务逻辑与数据建模。',
    tags: ['Go', 'MySQL', 'Redis', 'API 设计', 'DDD'],
    skills: ['go-unit-test', 'tRPC-Go DDD代码生成'],
    icon: '⚙️',
    color: 'green',
  },
  {
    name: '运维工程师',
    englishName: '博客运维',
    model: 'claude-4.6-opus',
    role: 'Linux 系统运维工程师，专注远端 CentOS 服务器管理。擅长 Nginx 配置、Docker 容器编排、系统监控与安全加固，负责构建、部署与服务验证。',
    tags: ['Linux', 'Nginx', 'Docker', '系统监控', '安全加固'],
    skills: ['nginx-config', 'remote-server'],
    icon: '🔧',
    color: 'orange',
  },
]

// Agent 卡片主题色映射
const agentColorMap: Record<string, { border: string; bg: string; tag: string; tagText: string; badge: string; badgeText: string }> = {
  purple: {
    border: 'border-l-purple-500',
    bg: 'bg-purple-50 dark:bg-purple-900/10',
    tag: 'bg-purple-100 dark:bg-purple-900/30',
    tagText: 'text-purple-700 dark:text-purple-300',
    badge: 'bg-purple-100 dark:bg-purple-900/30',
    badgeText: 'text-purple-600 dark:text-purple-400',
  },
  blue: {
    border: 'border-l-blue-500',
    bg: 'bg-blue-50 dark:bg-blue-900/10',
    tag: 'bg-blue-100 dark:bg-blue-900/30',
    tagText: 'text-blue-700 dark:text-blue-300',
    badge: 'bg-blue-100 dark:bg-blue-900/30',
    badgeText: 'text-blue-600 dark:text-blue-400',
  },
  green: {
    border: 'border-l-green-500',
    bg: 'bg-green-50 dark:bg-green-900/10',
    tag: 'bg-green-100 dark:bg-green-900/30',
    tagText: 'text-green-700 dark:text-green-300',
    badge: 'bg-green-100 dark:bg-green-900/30',
    badgeText: 'text-green-600 dark:text-green-400',
  },
  orange: {
    border: 'border-l-orange-500',
    bg: 'bg-orange-50 dark:bg-orange-900/10',
    tag: 'bg-orange-100 dark:bg-orange-900/30',
    tagText: 'text-orange-700 dark:text-orange-300',
    badge: 'bg-orange-100 dark:bg-orange-900/30',
    badgeText: 'text-orange-600 dark:text-orange-400',
  },
}

// Tab 导航
const tabs: TabItem[] = [
  { key: 'articles', label: '📝 文章' },
  { key: 'life', label: '📷 生活' },
  { key: 'tools', label: '🛠️ 小工具·小游戏' },
  { key: 'agent-team', label: '🤖 Agent Team' },
  { key: 'messages', label: '💬 留言板' },
  { key: 'changelog', label: '📋 更新日志' },
]
const activeTab = ref('articles')

function selectTab(key: string) {
  activeTab.value = key
  drawerOpen.value = false
}

// ====== 文章数据（SSR 预取） ======
const loading = ref(false)
const error = ref(false)

const { data: articlesData, refresh: refreshArticles } = await useAsyncData(
  'articles',
  () => apiFetchArticles(),
  { default: () => ({ list: [] as ArticleListItem[] }) }
)
const articles = computed(() => articlesData.value?.list ?? [])

// ====== 更新日志数据（SSR 预取） ======
const changelogLoading = ref(false)

const { data: changelogData } = await useAsyncData(
  'changelog',
  () => $fetch<ChangelogResponse>('/api/changelog'),
  { default: () => ({ changelog: [] as ChangelogItem[] }) }
)
const changelog = computed(() => changelogData.value?.changelog ?? [])

// ====== 博主个人资料（SSR 预取） ======
const { data: profileData } = await useAsyncData(
  'profile',
  () => apiGetProfile(),
  { default: () => ({ avatar: '', bio: '' }) }
)
const profile = computed<Profile>(() => ({
  avatar: profileData.value?.avatar || '',
  bio: profileData.value?.bio || '',
}))

async function fetchProfile() {
  // 保留此函数供管理后台修改后刷新使用
  try {
    const data = await apiGetProfile()
    profileData.value = data
  } catch {
    // 静默处理
  }
}

async function fetchChangelog() {
  changelogLoading.value = true
  try {
    const res = await $fetch<ChangelogResponse>('/api/changelog')
    changelogData.value = res
  } catch {
    // 静默处理
  } finally {
    changelogLoading.value = false
  }
}

async function fetchArticles() {
  loading.value = true
  error.value = false
  try {
    await refreshArticles()
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

// ====== 文章点赞 ======
const articleLikeStates = ref<Record<string, { liked: boolean; likeCount: number }>>({})
const articleLikeAnimating = ref<Record<string, boolean>>({})
const articleLikeAbortControllers: Record<string, AbortController> = {}

function getArticleLikeCount(article: ArticleListItem): number {
  return articleLikeStates.value[article.id]?.likeCount ?? article.likeCount ?? 0
}

function isArticleLiked(articleId: string): boolean {
  return articleLikeStates.value[articleId]?.liked ?? false
}

async function fetchArticleLikeStates() {
  if (!deviceId.value || !articles.value.length) return
  const ids = articles.value.map(a => String(a.id))
  try {
    const res = await apiGetArticleLikeStatusBatch(ids, deviceId.value)
    const likedSet = new Set(res.likedIds.map(String))
    articles.value.forEach(article => {
      articleLikeStates.value[String(article.id)] = {
        liked: likedSet.has(String(article.id)),
        likeCount: article.likeCount ?? 0,
      }
    })
  } catch {
    // 降级：用 SSR 的 likeCount，liked 默认 false
    articles.value.forEach(article => {
      articleLikeStates.value[String(article.id)] = { liked: false, likeCount: article.likeCount ?? 0 }
    })
  }
}

// 当 articles 数据就绪时，用 SSR likeCount 预填充（liked 默认 false，等 batch 接口回来后更新）
watch(articles, (newArticles) => {
  if (!newArticles?.length) return
  newArticles.forEach(article => {
    if (!articleLikeStates.value[String(article.id)]) {
      articleLikeStates.value[String(article.id)] = { liked: false, likeCount: article.likeCount ?? 0 }
    }
  })
}, { immediate: true })

async function handleArticleLike(articleId: string, e?: Event) {
  e?.preventDefault()
  e?.stopPropagation()
  if (!deviceId.value) return

  // 取消该文章上一次飞行中的请求
  if (articleLikeAbortControllers[articleId]) {
    articleLikeAbortControllers[articleId].abort()
  }
  const controller = new AbortController()
  articleLikeAbortControllers[articleId] = controller

  // 乐观更新（立即响应）
  const prev = articleLikeStates.value[articleId] || { liked: false, likeCount: 0 }
  const wasLiked = prev.liked
  articleLikeStates.value[articleId] = {
    liked: !wasLiked,
    likeCount: wasLiked ? Math.max(0, prev.likeCount - 1) : prev.likeCount + 1,
  }

  // 触发爱心动画
  articleLikeAnimating.value[articleId] = true
  setTimeout(() => { articleLikeAnimating.value[articleId] = false }, 600)

  try {
    const res = await apiToggleArticleLike(articleId, deviceId.value, controller.signal)
    // 只有当这个请求仍是最新的（没被取消）时才更新状态
    if (!controller.signal.aborted) {
      articleLikeStates.value[articleId] = { liked: res.liked, likeCount: res.likeCount }
    }
  } catch (err: any) {
    // 被取消的请求不回滚（UI 已经是最新乐观状态）
    if (err?.name === 'AbortError' || controller.signal.aborted) return
    // 真正的网络错误才回滚
    articleLikeStates.value[articleId] = prev
  } finally {
    // 清理 controller
    if (articleLikeAbortControllers[articleId] === controller) {
      delete articleLikeAbortControllers[articleId]
    }
  }
}

// ====== 相册功能 ======
const albums = ref<AlbumItem[]>([])
const albumsLoading = ref(false)
const albumsLoaded = ref(false)

const selectedAlbum = ref<AlbumItem | null>(null)
const albumPhotos = ref<PhotoItem[]>([])
const photosLoading = ref(false)

const albumTotalPhotos = computed(() => {
  return albums.value.reduce((sum, a) => sum + a.photoCount, 0)
})

/** 按年月分组照片（时间倒序） */
const groupedPhotos = computed(() => {
  const groups: { label: string; photos: PhotoItem[] }[] = []
  const map = new Map<string, PhotoItem[]>()
  for (const p of albumPhotos.value) {
    const d = new Date(p.createdAt)
    const key = `${d.getFullYear()}年${d.getMonth() + 1}月`
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(p)
  }
  for (const [label, photos] of map) {
    groups.push({ label, photos })
  }
  return groups
})

type PhotoListRow =
  | { id: string; type: 'header'; label: string }
  | { id: string; type: 'grid'; photos: PhotoItem[] }

/** 将分组拍平为虚拟滚动的行列表（标题行 + 每2列/5列的图片行） */
const flatPhotoRows = computed<PhotoListRow[]>(() => {
  const rows: PhotoListRow[] = []
  for (const group of groupedPhotos.value) {
    rows.push({ id: `header-${group.label}`, type: 'header', label: group.label })
    // 每行按 grid 列数拆分（移动端2列，PC端5列，这里统一按最小公倍数2拆分，CSS负责显示列数）
    for (let i = 0; i < group.photos.length; i += 5) {
      rows.push({
        id: `grid-${group.label}-${i}`,
        type: 'grid',
        photos: group.photos.slice(i, i + 5),
      })
    }
  }
  return rows
})

// ====== 密码保护 ======
const passwordModal = reactive({
  visible: false,
  type: '' as 'album' | 'photo',
  targetId: 0,
  password: '',
  loading: false,
  errorMsg: '',
})

/** 检查相册是否已解锁（sessionStorage） */
function isAlbumUnlocked(albumId: number): boolean {
  if (import.meta.server) return false
  return sessionStorage.getItem(`unlocked_album_${albumId}`) === '1'
}

/** 检查照片是否已解锁（sessionStorage） */
function isPhotoUnlocked(photoId: number): boolean {
  if (import.meta.server) return false
  return sessionStorage.getItem(`unlocked_photo_${photoId}`) === '1'
}

/** 标记相册已解锁 */
function markAlbumUnlocked(albumId: number) {
  sessionStorage.setItem(`unlocked_album_${albumId}`, '1')
}

/** 标记照片已解锁 */
function markPhotoUnlocked(photoId: number) {
  sessionStorage.setItem(`unlocked_photo_${photoId}`, '1')
}

function openPasswordModal(type: 'album' | 'photo', targetId: number) {
  passwordModal.type = type
  passwordModal.targetId = targetId
  passwordModal.password = ''
  passwordModal.loading = false
  passwordModal.errorMsg = ''
  passwordModal.visible = true
  // 下一帧聚焦输入框
  nextTick(() => {
    const input = document.getElementById('password-input') as HTMLInputElement | null
    input?.focus()
  })
}

function closePasswordModal() {
  passwordModal.visible = false
  passwordModal.password = ''
  passwordModal.errorMsg = ''
}

// 密码验证后的回调（由调用方设置）
let passwordVerifiedCallback: (() => void) | null = null

async function handlePasswordSubmit() {
  if (!passwordModal.password.trim()) {
    passwordModal.errorMsg = '请输入密码'
    return
  }
  passwordModal.loading = true
  passwordModal.errorMsg = ''
  try {
    let res: { success: boolean }
    if (passwordModal.type === 'album') {
      res = await apiVerifyAlbumPassword(passwordModal.targetId, passwordModal.password)
    } else {
      res = await apiVerifyPhotoPassword(passwordModal.targetId, passwordModal.password)
    }
    if (res.success) {
      if (passwordModal.type === 'album') {
        markAlbumUnlocked(passwordModal.targetId)
      } else {
        markPhotoUnlocked(passwordModal.targetId)
      }
      closePasswordModal()
      passwordVerifiedCallback?.()
      passwordVerifiedCallback = null
    } else {
      passwordModal.errorMsg = '密码错误，请重试'
    }
  } catch {
    passwordModal.errorMsg = '验证失败，请重试'
  } finally {
    passwordModal.loading = false
  }
}

async function fetchAlbums() {
  albumsLoading.value = true
  try {
    const res = await apiGetAlbums()
    albums.value = res.list
    albumsLoaded.value = true
  } catch {
    albums.value = []
  } finally {
    albumsLoading.value = false
  }
}

async function openAlbum(album: AlbumItem) {
  // 有密码且未解锁，弹出密码框
  if (album.hasPassword && !isAlbumUnlocked(album.id)) {
    passwordVerifiedCallback = () => openAlbum(album)
    openPasswordModal('album', album.id)
    return
  }
  selectedAlbum.value = album
  photosLoading.value = true
  try {
    const res = await apiGetPhotos(album.id)
    albumPhotos.value = res.list
    // 拉取各照片点赞数和踩数
    fetchLikesForPhotos(albumPhotos.value)
    fetchDislikesForPhotos(albumPhotos.value)
  } catch {
    albumPhotos.value = []
  } finally {
    photosLoading.value = false
  }
}

function closeAlbum() {
  selectedAlbum.value = null
  albumPhotos.value = []
}

// ====== 灯箱 ======
const lightbox = reactive({
  visible: false,
  index: 0,
  scale: 1,
  // swiper 切换时的水平拖动偏移（正常模式）
  swipeX: 0,
  swiping: false,
  // 放大时的平移偏移
  panX: 0,
  panY: 0,
  panning: false,
})

const lightboxCurrentPhoto = computed(() => {
  if (!lightbox.visible || albumPhotos.value.length === 0) return null
  return albumPhotos.value[lightbox.index] || null
})

function openLightbox(photo: PhotoItem) {
  // 有密码且未解锁，弹出密码框
  if (photo.hasPassword && !isPhotoUnlocked(photo.id)) {
    passwordVerifiedCallback = () => openLightbox(photo)
    openPasswordModal('photo', photo.id)
    return
  }
  const idx = albumPhotos.value.findIndex(p => p.id === photo.id)
  lightbox.index = idx >= 0 ? idx : 0
  lightbox.scale = 1
  lightbox.swipeX = 0
  lightbox.panX = 0
  lightbox.panY = 0
  lightbox.visible = true
}

function closeLightbox() {
  lightbox.visible = false
  lightbox.scale = 1
  lightbox.swipeX = 0
  lightbox.panX = 0
  lightbox.panY = 0
}

function prevPhoto() {
  if (lightbox.index > 0) {
    const targetPhoto = albumPhotos.value[lightbox.index - 1]
    if (targetPhoto?.hasPassword && !isPhotoUnlocked(targetPhoto.id)) {
      passwordVerifiedCallback = () => {
        lightbox.index--
        lightbox.scale = 1
        lightbox.swipeX = 0
        lightbox.panX = 0
        lightbox.panY = 0
      }
      openPasswordModal('photo', targetPhoto.id)
      return
    }
    lightbox.index--
    lightbox.scale = 1
    lightbox.swipeX = 0
    lightbox.panX = 0
    lightbox.panY = 0
  }
}

function nextPhoto() {
  if (lightbox.index < albumPhotos.value.length - 1) {
    const targetPhoto = albumPhotos.value[lightbox.index + 1]
    if (targetPhoto?.hasPassword && !isPhotoUnlocked(targetPhoto.id)) {
      passwordVerifiedCallback = () => {
        lightbox.index++
        lightbox.scale = 1
        lightbox.swipeX = 0
        lightbox.panX = 0
        lightbox.panY = 0
      }
      openPasswordModal('photo', targetPhoto.id)
      return
    }
    lightbox.index++
    lightbox.scale = 1
    lightbox.swipeX = 0
    lightbox.panX = 0
    lightbox.panY = 0
  }
}

// 键盘事件
function handleKeydown(e: KeyboardEvent) {
  if (!lightbox.visible) return
  if (e.key === 'ArrowLeft') prevPhoto()
  else if (e.key === 'ArrowRight') nextPhoto()
  else if (e.key === 'Escape') closeLightbox()
}

// ======== 触摸处理 ========
let touchStartX = 0
let touchStartY = 0
let touchStartDist = 0
let touchStartScale = 1
let touchStartPanX = 0
let touchStartPanY = 0
// 双指中点（相对于图片容器中心的偏移），用于以双指中点为缩放原点
let pinchCenterX = 0
let pinchCenterY = 0
const isPinchingRef = ref(false) // 需要模板访问，用 ref
let isTouchActive = false

function getTouchDist(touches: TouchList): number {
  if (touches.length < 2) return 0
  const dx = touches[0].clientX - touches[1].clientX
  const dy = touches[0].clientY - touches[1].clientY
  return Math.sqrt(dx * dx + dy * dy)
}

function onTouchStart(e: TouchEvent) {
  isTouchActive = true
  if (e.touches.length === 2) {
    // 双指缩放开始
    isPinchingRef.value = true
    lightbox.swiping = false
    lightbox.panning = false
    touchStartDist = getTouchDist(e.touches)
    touchStartScale = lightbox.scale
    touchStartPanX = lightbox.panX
    touchStartPanY = lightbox.panY
    // 记录双指中点（屏幕坐标转换为相对于屏幕中心的偏移）
    pinchCenterX = (e.touches[0].clientX + e.touches[1].clientX) / 2 - window.innerWidth / 2
    pinchCenterY = (e.touches[0].clientY + e.touches[1].clientY) / 2 - window.innerHeight / 2
  } else if (e.touches.length === 1) {
    isPinchingRef.value = false
    touchStartX = e.touches[0].clientX
    touchStartY = e.touches[0].clientY
    touchStartPanX = lightbox.panX
    touchStartPanY = lightbox.panY
    if (lightbox.scale > 1) {
      // 放大模式：拖动平移
      lightbox.panning = true
      lightbox.swiping = false
    } else if (albumPhotos.value.length > 1) {
      // 正常模式：swipe 切换（仅多张图片时启用）
      lightbox.swiping = true
      lightbox.panning = false
    }
  }
}

function onTouchMove(e: TouchEvent) {
  if (!isTouchActive) return

  if (isPinchingRef.value && e.touches.length === 2) {
    // 双指缩放 —— 完全跟手，无任何 transition
    const dist = getTouchDist(e.touches)
    const ratio = dist / touchStartDist
    const newScale = Math.min(4, Math.max(1, touchStartScale * ratio))

    // 以双指中点为缩放原点：
    // 当 scale 从 S0 变为 S1，图片上双指中点对应的位置不变
    // 需要调整 panX/panY 使中点锁定
    const scaleDelta = newScale - touchStartScale
    lightbox.scale = newScale

    if (newScale > 1) {
      // 以 pinchCenter 为原点调整平移，使缩放原点保持在双指中点
      const newPanX = touchStartPanX - pinchCenterX * scaleDelta / touchStartScale
      const newPanY = touchStartPanY - pinchCenterY * scaleDelta / touchStartScale
      const maxPan = (newScale - 1) * window.innerWidth * 0.5
      const maxPanY = (newScale - 1) * window.innerHeight * 0.5
      lightbox.panX = Math.max(-maxPan, Math.min(maxPan, newPanX))
      lightbox.panY = Math.max(-maxPanY, Math.min(maxPanY, newPanY))
    } else {
      lightbox.panX = 0
      lightbox.panY = 0
    }
    return
  }

  if (e.touches.length !== 1) return

  const dx = e.touches[0].clientX - touchStartX
  const dy = e.touches[0].clientY - touchStartY

  if (lightbox.scale > 1 && lightbox.panning) {
    // 放大模式：平移图片（限制在合理范围内）
    const maxPan = (lightbox.scale - 1) * window.innerWidth * 0.5
    const maxPanY = (lightbox.scale - 1) * window.innerHeight * 0.5
    lightbox.panX = Math.max(-maxPan, Math.min(maxPan, touchStartPanX + dx))
    lightbox.panY = Math.max(-maxPanY, Math.min(maxPanY, touchStartPanY + dy))
  } else if (lightbox.scale <= 1 && lightbox.swiping) {
    // 正常模式：跟手拖动（带阻尼）
    const resistance = Math.abs(dy) > Math.abs(dx) ? 0.3 : 1
    lightbox.swipeX = dx * resistance
  }
}

function onTouchEnd(e: TouchEvent) {
  if (!isTouchActive) return
  isTouchActive = false

  if (isPinchingRef.value) {
    isPinchingRef.value = false
    // 缩放回 1 时重置 pan，带弹性（transition 在 scale===1 时会生效）
    if (lightbox.scale <= 1) {
      lightbox.scale = 1
      lightbox.panX = 0
      lightbox.panY = 0
    }
    lightbox.panning = false
    return
  }

  if (lightbox.scale <= 1 && lightbox.swiping) {
    // swipe 切换判断：超过屏幕宽度 25% 则切换
    const threshold = window.innerWidth * 0.25
    if (lightbox.swipeX < -threshold) {
      nextPhoto()
    } else if (lightbox.swipeX > threshold) {
      prevPhoto()
    }
    // 无论是否切换成功，都重置 swipeX（防止只有一张图时卡住）
    lightbox.swipeX = 0
    lightbox.swiping = false
  }

  lightbox.panning = false
}

// ======== PC 端鼠标拖动处理 ========
let mouseDown = false
let mouseStartX = 0
let mouseStartY = 0
let mouseStartPanX = 0
let mouseStartPanY = 0
let mouseMoved = false

function onMouseDown(e: MouseEvent) {
  // 仅左键，且非按钮/链接等可交互元素
  if (e.button !== 0) return
  const tag = (e.target as HTMLElement).tagName?.toLowerCase()
  if (tag === 'button' || tag === 'a' || tag === 'svg' || tag === 'path' || tag === 'span') return

  mouseDown = true
  mouseMoved = false
  mouseStartX = e.clientX
  mouseStartY = e.clientY
  mouseStartPanX = lightbox.panX
  mouseStartPanY = lightbox.panY

  if (lightbox.scale > 1) {
    lightbox.panning = true
  } else if (albumPhotos.value.length > 1) {
    lightbox.swiping = true
  }
  e.preventDefault()
}

function onMouseMove(e: MouseEvent) {
  if (!mouseDown) return
  mouseMoved = true
  const dx = e.clientX - mouseStartX
  const dy = e.clientY - mouseStartY

  if (lightbox.scale > 1) {
    // 放大模式：拖动平移
    const maxPan = (lightbox.scale - 1) * window.innerWidth * 0.5
    const maxPanY = (lightbox.scale - 1) * window.innerHeight * 0.5
    lightbox.panX = Math.max(-maxPan, Math.min(maxPan, mouseStartPanX + dx))
    lightbox.panY = Math.max(-maxPanY, Math.min(maxPanY, mouseStartPanY + dy))
  } else if (lightbox.swiping) {
    // 正常模式：swipe 切换（仅 swiping 激活时才更新偏移）
    const resistance = Math.abs(dy) > Math.abs(dx) ? 0.3 : 1
    lightbox.swipeX = dx * resistance
  }
  e.preventDefault()
}

function onMouseUp(_e: MouseEvent) {
  if (!mouseDown) return
  mouseDown = false

  if (lightbox.scale > 1) {
    lightbox.panning = false
  } else if (lightbox.swiping) {
    const threshold = window.innerWidth * 0.25
    if (lightbox.swipeX < -threshold) {
      nextPhoto()
    } else if (lightbox.swipeX > threshold) {
      prevPhoto()
    }
    // 无论是否切换成功，都重置 swipeX（防止只有一张图时卡住）
    lightbox.swipeX = 0
    lightbox.swiping = false
  }
}

// ======== PC 端鼠标滚轮缩放 ========
function onWheel(e: WheelEvent) {
  if (!lightbox.visible) return

  // 根据滚轮方向计算新缩放比例（向上滚放大，向下滚缩小）
  const delta = e.deltaY > 0 ? -0.15 : 0.15
  const oldScale = lightbox.scale
  const newScale = Math.min(4, Math.max(1, oldScale + delta))

  if (newScale === oldScale) return

  // 以鼠标位置为缩放原点：计算鼠标相对于屏幕中心的偏移
  const mouseOffsetX = e.clientX - window.innerWidth / 2
  const mouseOffsetY = e.clientY - window.innerHeight / 2

  // 调整 panX/panY 使鼠标指向的图片位置保持不变
  const scaleDelta = newScale - oldScale
  lightbox.scale = newScale

  if (newScale > 1) {
    const newPanX = lightbox.panX - mouseOffsetX * scaleDelta / oldScale
    const newPanY = lightbox.panY - mouseOffsetY * scaleDelta / oldScale
    const maxPanX = (newScale - 1) * window.innerWidth * 0.5
    const maxPanY = (newScale - 1) * window.innerHeight * 0.5
    lightbox.panX = Math.max(-maxPanX, Math.min(maxPanX, newPanX))
    lightbox.panY = Math.max(-maxPanY, Math.min(maxPanY, newPanY))
  } else {
    // 缩回 1x 时重置平移
    lightbox.panX = 0
    lightbox.panY = 0
  }
}

// ====== 留言板 ======
const messages = ref<MessageItem[]>([])
const messagesLoading = ref(false)
const messagesLoaded = ref(false)

// 留言输入表单
const msgForm = reactive({
  nickname: '',
  content: '',
  submitting: false,
  editingId: null as number | null, // 正在编辑的留言 ID
})

/** 当前设备已发过的留言 */
const myMessage = computed(() => messages.value.find(m => m.isOwn))

/** 相对时间格式化 */
function relativeTime(dateStr: string): string {
  const now = Date.now()
  const diff = now - new Date(dateStr).getTime()
  const seconds = Math.floor(diff / 1000)
  if (seconds < 60) return '刚刚'
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

async function fetchMessages() {
  messagesLoading.value = true
  try {
    const res = await apiGetMessages(deviceId.value || undefined)
    messages.value = res.list
    messagesLoaded.value = true
  } catch {
    messages.value = []
  } finally {
    messagesLoading.value = false
  }
}

async function handleSubmitMessage() {
  if (!msgForm.content.trim()) return
  if (!deviceId.value) return
  msgForm.submitting = true
  try {
    if (msgForm.editingId) {
      // 修改留言
      const updated = await apiUpdateMessage(msgForm.editingId, {
        deviceId: deviceId.value,
        nickname: msgForm.nickname.trim() || undefined,
        content: msgForm.content.trim(),
      })
      // 更新列表中的该条留言
      const idx = messages.value.findIndex(m => m.id === msgForm.editingId)
      if (idx !== -1) messages.value[idx] = updated
      msgForm.editingId = null
      msgForm.nickname = ''
      msgForm.content = ''
    } else {
      // 新增留言
      const created = await apiCreateMessage({
        deviceId: deviceId.value,
        nickname: msgForm.nickname.trim() || undefined,
        content: msgForm.content.trim(),
      })
      // 插到列表头部
      messages.value.unshift(created)
      msgForm.nickname = ''
      msgForm.content = ''
    }
  } catch (err: unknown) {
    const fetchErr = err as { statusCode?: number; statusMessage?: string; data?: { statusMessage?: string } }
    const errMsg = fetchErr?.data?.statusMessage || fetchErr?.statusMessage || '操作失败'
    // 用简单的 alert 提示（后续可优化为 toast）
    if (import.meta.client) {
      alert(errMsg)
    }
  } finally {
    msgForm.submitting = false
  }
}

function startEditMessage(msg: MessageItem) {
  msgForm.editingId = msg.id
  msgForm.nickname = msg.nickname === '匿名' ? '' : msg.nickname
  msgForm.content = msg.content
  // 滚动到输入区
  nextTick(() => {
    document.getElementById('msg-content-input')?.focus()
  })
}

function cancelEditMessage() {
  msgForm.editingId = null
  msgForm.nickname = ''
  msgForm.content = ''
}

// 生活 Tab 切换时加载相册
watch(activeTab, (val) => {
  if (val === 'life' && !albumsLoaded.value) {
    fetchAlbums()
  }
  if (val === 'messages' && !messagesLoaded.value) {
    fetchMessages()
  }
})

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// 页面挂载
onMounted(() => {
  deviceId.value = getOrCreateDeviceId()
  // articles / changelog / profile 已在 SSR 阶段预取，无需再次请求
  // 客户端拉取文章点赞状态（依赖 localStorage deviceId，不做 SSR）
  nextTick(() => fetchArticleLikeStates())
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

useHead({
  title: 'fatwillzeng - 个人博客',
  meta: [
    { name: 'description', content: '全栈开发者 fatwillzeng 的个人博客，分享技术与生活。' },
  ],
})
</script>

<style scoped>
.drawer-overlay-enter-active,
.drawer-overlay-leave-active {
  transition: opacity 0.3s ease;
}
.drawer-overlay-enter-from,
.drawer-overlay-leave-to {
  opacity: 0;
}
.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: transform 0.3s ease;
}
.drawer-slide-enter-from,
.drawer-slide-leave-to {
  transform: translateX(-100%);
}
.lightbox-fade-enter-active,
.lightbox-fade-leave-active {
  transition: opacity 0.25s ease;
}
.lightbox-fade-enter-from,
.lightbox-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

/* 虚拟滚动容器：铺满内容区，由父容器决定可视高度 */
.article-scroller,
.photo-scroller {
  height: calc(100vh - 220px);
  min-height: 400px;
}

/* +1 飘动消失动画 */
@keyframes floatUp {
  0% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
  100% {
    opacity: 0;
    transform: translateY(-32px) scale(1.3);
  }
}

.float-plus-one {
  animation: floatUp 0.8s ease-out forwards;
  pointer-events: none;
}

/* 文章点赞爱心弹跳动画 */
@keyframes heartBounce {
  0% { transform: scale(1); }
  30% { transform: scale(1.3); }
  60% { transform: scale(0.9); }
  100% { transform: scale(1); }
}
</style>
