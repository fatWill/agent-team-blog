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
          <nav class="px-3 py-4">
            <button
              v-for="tab in tabs"
              :key="tab.key"
              class="flex w-full items-center rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="activeTab === tab.key
                ? 'bg-primary-50 text-primary-600 dark:bg-primary-900/20 dark:text-primary-400'
                : 'text-gray-600 hover:bg-gray-50 dark:text-gray-400 dark:hover:bg-gray-800'"
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
          :src="profile.avatar"
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
    <nav class="sticky top-0 z-40 hidden border-b border-gray-200/60 bg-white/80 backdrop-blur-lg transition-colors duration-300 dark:border-gray-700/60 dark:bg-gray-900/80 md:block">
      <div class="mx-auto flex max-w-3xl gap-0 px-4">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="relative px-4 py-3 text-sm font-medium transition-colors"
          :class="activeTab === tab.key
            ? 'text-gray-900 dark:text-gray-100'
            : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
          <span
            v-if="activeTab === tab.key"
            class="absolute bottom-0 left-1/2 h-0.5 w-8 -translate-x-1/2 rounded-full bg-primary-500"
          />
        </button>
      </div>
    </nav>

    <!-- 移动端当前 Tab 指示（仅 < md 显示） -->
    <div class="border-b border-gray-200/60 px-4 py-2 dark:border-gray-700/60 md:hidden">
      <span class="text-xs font-medium text-primary-500">{{ currentTabLabel }}</span>
    </div>

    <!-- 内容区域 -->
    <main class="mx-auto max-w-3xl px-4 py-8">
      <!-- 文章 Tab -->
      <div v-if="activeTab === 'articles'">
        <div v-if="loading" class="space-y-4">
          <div v-for="i in 3" :key="i" class="animate-pulse rounded-xl border border-gray-100 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
            <div class="mb-3 h-5 w-3/4 rounded bg-gray-200 dark:bg-gray-700" />
            <div class="mb-2 h-4 w-full rounded bg-gray-100 dark:bg-gray-700/50" />
            <div class="h-4 w-1/2 rounded bg-gray-100 dark:bg-gray-700/50" />
          </div>
        </div>
        <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-6 text-center dark:border-red-800 dark:bg-red-900/20">
          <p class="text-red-600 dark:text-red-400">加载文章失败，请稍后重试</p>
          <button class="mt-3 text-sm text-primary-500 hover:text-primary-600" @click="fetchArticles">重新加载</button>
        </div>
        <div v-else-if="articles.length > 0" class="space-y-4">
          <NuxtLink
            v-for="article in articles"
            :key="article.id"
            :to="`/articles/${article.id}`"
            class="block overflow-hidden rounded-xl border border-gray-100 bg-white transition-all duration-200 hover:shadow-md hover:-translate-y-0.5 dark:border-gray-700 dark:bg-gray-800 dark:hover:border-gray-600"
          >
            <img v-if="article.coverImage" :src="article.coverImage" :alt="article.title" class="h-40 w-full object-cover" />
            <div class="p-5">
              <h2 class="mb-2 text-lg font-semibold text-gray-900 dark:text-gray-100">{{ article.title }}</h2>
              <p v-if="article.summary" class="mb-3 text-sm leading-relaxed text-gray-500 dark:text-gray-400 line-clamp-2">{{ article.summary }}</p>
              <time class="text-xs text-gray-400 dark:text-gray-500">{{ formatDate(article.createdAt) }}</time>
            </div>
          </NuxtLink>
        </div>
        <div v-else class="py-20 text-center">
          <p class="text-gray-400 dark:text-gray-500">暂无文章</p>
        </div>
      </div>

      <!-- ========== 生活 Tab：相册功能 ========== -->
      <div v-else-if="activeTab === 'life'">
        <!-- 相册集列表（默认视图） -->
        <div v-if="!selectedAlbum">
          <!-- 标题区 -->
          <div class="mb-6 flex items-center justify-between">
            <div>
              <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">我的相册</h2>
              <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">
                共 {{ albumTotalPhotos }} 张照片
              </p>
            </div>
          </div>

          <!-- 加载骨架屏 -->
          <div v-if="albumsLoading" class="grid grid-cols-2 gap-3 md:grid-cols-5">
            <div v-for="i in 6" :key="i" class="animate-pulse">
              <div class="aspect-square rounded-lg bg-gray-200 dark:bg-gray-700" />
              <div class="mt-2 h-4 w-2/3 rounded bg-gray-200 dark:bg-gray-700" />
              <div class="mt-1 h-3 w-1/3 rounded bg-gray-100 dark:bg-gray-700/50" />
            </div>
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
                  :src="album.coverUrl"
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

          <!-- 加载骨架屏 -->
          <div v-if="photosLoading" class="grid grid-cols-2 gap-2 md:grid-cols-5">
            <div v-for="i in 8" :key="i" class="animate-pulse">
              <div class="aspect-square rounded-lg bg-gray-200 dark:bg-gray-700" />
            </div>
          </div>

          <!-- 按年月分组展示照片 -->
          <div v-else-if="albumPhotos.length > 0" class="space-y-6">
            <div v-for="group in groupedPhotos" :key="group.label">
              <h3 class="mb-3 text-sm font-medium text-gray-500 dark:text-gray-400">{{ group.label }}</h3>
              <div class="grid grid-cols-2 gap-2 md:grid-cols-5">
                <button
                  v-for="(photo, idx) in group.photos"
                  :key="photo.id"
                  class="group relative aspect-square overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800"
                  @click="openLightbox(photo)"
                >
                  <img
                    :src="photo.url"
                    :alt="photo.caption || '照片'"
                    class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                  />
                </button>
              </div>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-else class="py-20 text-center">
            <p class="text-gray-400 dark:text-gray-500">这个相册还没有照片</p>
          </div>
        </div>
      </div>

      <!-- 更新日志 Tab -->
      <div v-else-if="activeTab === 'changelog'">
        <div v-if="changelogLoading" class="space-y-6">
          <div v-for="i in 2" :key="i" class="animate-pulse">
            <div class="mb-3 h-5 w-32 rounded bg-gray-200 dark:bg-gray-700" />
            <div class="space-y-2 pl-4 border-l-2 border-gray-200 dark:border-gray-700">
              <div class="h-4 w-3/4 rounded bg-gray-100 dark:bg-gray-700/50" />
              <div class="h-4 w-2/3 rounded bg-gray-100 dark:bg-gray-700/50" />
            </div>
          </div>
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
          class="fixed inset-0 z-[100] flex items-center justify-center bg-black/90"
          @click.self="closeLightbox"
          @keydown.left="prevPhoto"
          @keydown.right="nextPhoto"
          @keydown.escape="closeLightbox"
        >
          <!-- 关闭按钮 -->
          <button
            class="absolute right-4 top-4 z-10 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
            @click="closeLightbox"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          <!-- 左箭头 -->
          <button
            v-if="lightbox.index > 0"
            class="absolute left-4 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
            @click.stop="prevPhoto"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
            </svg>
          </button>

          <!-- 右箭头 -->
          <button
            v-if="lightbox.index < albumPhotos.length - 1"
            class="absolute right-4 top-1/2 z-10 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
            @click.stop="nextPhoto"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
            </svg>
          </button>

          <!-- 图片区域（支持触摸滑动 + 双指缩放） -->
          <div
            class="flex max-h-[90vh] max-w-[90vw] items-center justify-center"
            @click.stop
            @touchstart="onTouchStart"
            @touchmove="onTouchMove"
            @touchend="onTouchEnd"
          >
            <img
              v-if="lightboxCurrentPhoto"
              :src="lightboxCurrentPhoto.url"
              :alt="lightboxCurrentPhoto.caption || '照片'"
              class="max-h-[90vh] max-w-[90vw] object-contain transition-transform duration-200"
              :style="{ transform: `scale(${lightbox.scale})` }"
              draggable="false"
            />
          </div>

          <!-- 底部信息 -->
          <div class="absolute bottom-4 left-1/2 -translate-x-1/2 text-center">
            <span class="text-sm text-white/70">{{ lightbox.index + 1 }} / {{ albumPhotos.length }}</span>
            <p v-if="lightboxCurrentPhoto?.caption" class="mt-1 text-sm text-white/50">{{ lightboxCurrentPhoto.caption }}</p>
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
import type { ArticleListItem, TabItem, ChangelogItem, ChangelogResponse, Profile, AlbumItem, PhotoItem } from '~/types'
import { apiFetchArticles, apiGetProfile, apiGetAlbums, apiGetPhotos } from '~/utils/api'

const { isDark, toggleTheme } = useTheme()

// 移动端抽屉状态
const drawerOpen = ref(false)

// Tab 导航
const tabs: TabItem[] = [
  { key: 'articles', label: '文章' },
  { key: 'life', label: '生活' },
  { key: 'tools', label: '小工具·小游戏' },
  { key: 'agent-team', label: 'agent team' },
  { key: 'changelog', label: '更新日志' },
]
const activeTab = ref('articles')

const currentTabLabel = computed(() => {
  return tabs.find(t => t.key === activeTab.value)?.label ?? ''
})

function selectTab(key: string) {
  activeTab.value = key
  drawerOpen.value = false
}

// ====== 文章数据 ======
const articles = ref<ArticleListItem[]>([])
const loading = ref(false)
const error = ref(false)

// ====== 更新日志数据 ======
const changelog = ref<ChangelogItem[]>([])
const changelogLoading = ref(false)

// ====== 博主个人资料 ======
const profile = reactive<Profile>({ avatar: '', bio: '' })

async function fetchProfile() {
  try {
    const data = await apiGetProfile()
    profile.avatar = data.avatar || ''
    profile.bio = data.bio || ''
  } catch {
    // 静默处理
  }
}

async function fetchChangelog() {
  changelogLoading.value = true
  try {
    const res = await $fetch<ChangelogResponse>('/api/changelog')
    changelog.value = res.changelog
  } catch {
    changelog.value = []
  } finally {
    changelogLoading.value = false
  }
}

async function fetchArticles() {
  loading.value = true
  error.value = false
  try {
    const res = await apiFetchArticles()
    articles.value = res.list
  } catch {
    error.value = true
  } finally {
    loading.value = false
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
  selectedAlbum.value = album
  photosLoading.value = true
  try {
    const res = await apiGetPhotos(album.id)
    albumPhotos.value = res.list
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
})

const lightboxCurrentPhoto = computed(() => {
  if (!lightbox.visible || albumPhotos.value.length === 0) return null
  return albumPhotos.value[lightbox.index] || null
})

function openLightbox(photo: PhotoItem) {
  const idx = albumPhotos.value.findIndex(p => p.id === photo.id)
  lightbox.index = idx >= 0 ? idx : 0
  lightbox.scale = 1
  lightbox.visible = true
}

function closeLightbox() {
  lightbox.visible = false
  lightbox.scale = 1
}

function prevPhoto() {
  if (lightbox.index > 0) {
    lightbox.index--
    lightbox.scale = 1
  }
}

function nextPhoto() {
  if (lightbox.index < albumPhotos.value.length - 1) {
    lightbox.index++
    lightbox.scale = 1
  }
}

// 键盘事件
function handleKeydown(e: KeyboardEvent) {
  if (!lightbox.visible) return
  if (e.key === 'ArrowLeft') prevPhoto()
  else if (e.key === 'ArrowRight') nextPhoto()
  else if (e.key === 'Escape') closeLightbox()
}

// 触摸滑动 + 双指缩放
let touchStartX = 0
let touchStartY = 0
let touchStartDist = 0
let touchStartScale = 1
let isSwiping = false
let isPinching = false

function getTouchDist(touches: TouchList): number {
  if (touches.length < 2) return 0
  const dx = touches[0].clientX - touches[1].clientX
  const dy = touches[0].clientY - touches[1].clientY
  return Math.sqrt(dx * dx + dy * dy)
}

function onTouchStart(e: TouchEvent) {
  if (e.touches.length === 2) {
    // 双指缩放开始
    isPinching = true
    isSwiping = false
    touchStartDist = getTouchDist(e.touches)
    touchStartScale = lightbox.scale
  } else if (e.touches.length === 1) {
    isSwiping = true
    isPinching = false
    touchStartX = e.touches[0].clientX
    touchStartY = e.touches[0].clientY
  }
}

function onTouchMove(e: TouchEvent) {
  if (isPinching && e.touches.length === 2) {
    e.preventDefault()
    const dist = getTouchDist(e.touches)
    const ratio = dist / touchStartDist
    lightbox.scale = Math.min(4, Math.max(0.5, touchStartScale * ratio))
  }
}

function onTouchEnd(e: TouchEvent) {
  if (isPinching) {
    isPinching = false
    return
  }
  if (isSwiping && e.changedTouches.length === 1) {
    const dx = e.changedTouches[0].clientX - touchStartX
    const dy = e.changedTouches[0].clientY - touchStartY
    // 水平滑动距离 > 50px 且大于垂直距离
    if (Math.abs(dx) > 50 && Math.abs(dx) > Math.abs(dy)) {
      if (dx > 0) prevPhoto()
      else nextPhoto()
    }
    isSwiping = false
  }
}

// 生活 Tab 切换时加载相册
watch(activeTab, (val) => {
  if (val === 'life' && !albumsLoaded.value) {
    fetchAlbums()
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
  fetchArticles()
  fetchChangelog()
  fetchProfile()
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
</style>
