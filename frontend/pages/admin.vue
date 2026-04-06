<template>
  <div class="flex h-screen flex-col overflow-hidden bg-gray-50 transition-colors dark:bg-gray-900">
    <!-- 顶部导航 -->
    <header class="shrink-0 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg transition-colors duration-300 dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex h-14 max-w-4xl items-center justify-between px-4">
        <div class="flex items-center gap-3">
          <NuxtLink to="/home" class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            ← 首页
          </NuxtLink>
          <span class="text-gray-300 dark:text-gray-600">/</span>
          <span class="text-sm font-medium text-gray-900 dark:text-gray-100">管理后台</span>
        </div>
        <div class="flex items-center gap-2">
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
          <button
            class="rounded-lg px-3 py-1.5 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200"
            @click="handleLogout"
          >
            退出登录
          </button>
        </div>
      </div>
    </header>

    <!-- Tab 切换 -->
    <nav class="shrink-0 border-b border-gray-200/60 bg-white/80 backdrop-blur-lg dark:border-gray-700/60 dark:bg-gray-900/80">
      <div class="mx-auto flex max-w-4xl gap-0 overflow-x-auto px-4">
        <button
          v-for="tab in adminTabs"
          :key="tab.key"
          class="relative shrink-0 px-5 py-3 text-sm font-medium transition-colors"
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

    <!-- 内容区域 -->
    <main class="flex-1 overflow-y-auto">
      <div class="mx-auto max-w-4xl px-4 py-8">

      <!-- ========== Tab 1: 写文章 / 编辑文章 ========== -->
      <div v-if="activeTab === 'write'">
        <div class="mb-6 flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">
            {{ editingArticleId ? '编辑文章' : '写文章' }}
          </h2>
          <button
            v-if="editingArticleId"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
            @click="cancelEdit"
          >
            取消编辑
          </button>
        </div>

        <form class="space-y-5" @submit.prevent="handlePublish">
          <div>
            <input
              v-model="form.title"
              type="text"
              placeholder="文章标题"
              required
              class="w-full rounded-lg border border-gray-200 bg-white px-4 py-3 text-lg font-medium text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400"
            />
          </div>

          <!-- 封面图上传 -->
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">封面图（可选）</label>
            <input ref="coverFileInput" type="file" accept="image/*" class="hidden" @change="handleCoverUpload" />
            <div
              v-if="!form.coverImage"
              class="flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-300 bg-gray-50 px-6 py-8 transition-colors hover:border-primary-400 hover:bg-primary-50/30 dark:border-gray-600 dark:bg-gray-800/50 dark:hover:border-primary-500 dark:hover:bg-primary-900/10"
              :class="{ 'pointer-events-none opacity-60': coverUploading }"
              @click="coverFileInput?.click()"
              @dragover.prevent
              @drop.prevent="handleCoverDrop"
            >
              <svg v-if="!coverUploading" class="mb-2 h-10 w-10 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6.827 6.175A2.31 2.31 0 015.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 00-1.134-.175 2.31 2.31 0 01-1.64-1.055l-.822-1.316a2.192 2.192 0 00-1.736-1.039 48.774 48.774 0 00-5.232 0 2.192 2.192 0 00-1.736 1.039l-.821 1.316z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 12.75a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0z" />
              </svg>
              <div v-if="coverUploading" class="flex flex-col items-center gap-2">
                <svg class="h-8 w-8 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
                <div class="w-full max-w-xs bg-gray-200 rounded-full h-1.5 dark:bg-gray-600">
                  <div class="bg-primary-500 h-1.5 rounded-full transition-all duration-300" :style="{ width: coverUploadPercent + '%' }" />
                </div>
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ coverUploadPercent }}%</span>
              </div>
              <span v-if="!coverUploading" class="text-sm font-medium text-gray-600 dark:text-gray-400">点击上传封面图</span>
              <span v-if="!coverUploading" class="mt-1 text-xs text-gray-400 dark:text-gray-500">支持 JPG/PNG/GIF/WebP</span>
            </div>
            <div v-else class="relative inline-block">
              <img :src="toCdnUrl(form.coverImage)" alt="封面图预览" class="h-40 w-auto rounded-lg border border-gray-200 object-cover dark:border-gray-600" />
              <button type="button" class="absolute -right-2 -top-2 flex h-6 w-6 items-center justify-center rounded-full bg-red-500 text-white shadow-md transition-colors hover:bg-red-600" @click="form.coverImage = ''">
                <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
          </div>

          <div>
            <textarea v-model="form.summary" placeholder="文章摘要（可选）" rows="2" class="w-full resize-none rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400" />
          </div>

          <!-- Tiptap 编辑器 -->
          <div class="tiptap-editor overflow-hidden rounded-lg border border-gray-200 bg-white transition-colors dark:border-gray-600 dark:bg-gray-800">
            <div v-if="editor" class="flex flex-wrap gap-1 border-b border-gray-200 px-3 py-2 dark:border-gray-600">
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('bold') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleBold().run()"><strong>B</strong></button>
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('italic') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleItalic().run()"><em>I</em></button>
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('strike') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleStrike().run()"><s>S</s></button>
              <span class="mx-1 w-px bg-gray-200 dark:bg-gray-600" />
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('heading', { level: 1 }) ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleHeading({ level: 1 }).run()">H1</button>
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('heading', { level: 2 }) ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()">H2</button>
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('heading', { level: 3 }) ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleHeading({ level: 3 }).run()">H3</button>
              <span class="mx-1 w-px bg-gray-200 dark:bg-gray-600" />
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('bulletList') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleBulletList().run()">• 列表</button>
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('orderedList') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleOrderedList().run()">1. 列表</button>
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('blockquote') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleBlockquote().run()">引用</button>
              <button type="button" class="rounded px-2 py-1 text-sm transition-colors" :class="editor.isActive('codeBlock') ? 'bg-gray-200 text-gray-900 dark:bg-gray-600 dark:text-gray-100' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'" @click="editor.chain().focus().toggleCodeBlock().run()">代码块</button>
              <span class="mx-1 w-px bg-gray-200 dark:bg-gray-600" />
              <button type="button" class="rounded px-2 py-1 text-sm text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700" @click="addImage">图片</button>
              <button type="button" class="rounded px-2 py-1 text-sm text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700" @click="addLink">链接</button>
              <button type="button" class="rounded px-2 py-1 text-sm text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700" @click="editor.chain().focus().setHorizontalRule().run()">分割线</button>
            </div>
            <div class="text-gray-900 dark:text-gray-100">
              <EditorContent :editor="editor" />
            </div>
          </div>

          <div class="flex items-center justify-end gap-3">
            <span v-if="publishSuccess" class="text-sm text-green-600 dark:text-green-400">✓ {{ editingArticleId ? '保存成功！' : '发布成功！' }}</span>
            <span v-if="publishError" class="text-sm text-red-500">{{ publishError }}</span>
            <button type="submit" :disabled="publishing || !form.title.trim()" class="rounded-lg bg-primary-500 px-6 py-2.5 text-sm font-medium text-white transition-all duration-200 hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-300 disabled:cursor-not-allowed disabled:opacity-60">
              {{ publishing ? (editingArticleId ? '保存中...' : '发布中...') : (editingArticleId ? '保存修改' : '发布文章') }}
            </button>
          </div>
        </form>
      </div>

      <!-- ========== Tab 2: 文章管理 ========== -->
      <div v-if="activeTab === 'manage'">
        <h2 class="mb-6 text-xl font-bold text-gray-900 dark:text-gray-100">文章管理</h2>
        <div class="mb-6">
          <input v-model="searchKeyword" type="text" placeholder="按标题搜索文章..." class="w-full rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400" />
        </div>
        <div v-if="manageLoading" class="flex items-center justify-center py-20">
          <AppLoading tip="加载中..." />
        </div>
        <div v-else-if="manageArticles.length > 0" class="space-y-3">
          <div v-for="article in manageArticles" :key="article.id" class="rounded-lg border border-gray-100 bg-white p-4 transition-colors dark:border-gray-700 dark:bg-gray-800">
            <div class="flex items-start gap-4">
              <img v-if="article.coverImage" :src="toCdnUrl(article.coverImage)" :alt="article.title" class="h-16 w-24 flex-shrink-0 rounded-md object-cover" />
              <div v-else class="flex h-16 w-24 flex-shrink-0 items-center justify-center rounded-md bg-gray-100 dark:bg-gray-700">
                <svg class="h-6 w-6 text-gray-300 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909M3.75 21h16.5a1.5 1.5 0 001.5-1.5V5.25a1.5 1.5 0 00-1.5-1.5H3.75a1.5 1.5 0 00-1.5 1.5v14.25a1.5 1.5 0 001.5 1.5z" /></svg>
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="truncate text-sm font-semibold text-gray-900 dark:text-gray-100">{{ article.title }}</h3>
                <p v-if="article.summary" class="mt-1 truncate text-xs text-gray-500 dark:text-gray-400">{{ article.summary }}</p>
                <time class="mt-1 block text-xs text-gray-400 dark:text-gray-500">{{ formatDate(article.createdAt) }}</time>
              </div>
              <div class="flex flex-shrink-0 items-center gap-2">
                <button class="rounded-md px-3 py-1.5 text-xs font-medium text-primary-500 transition-colors hover:bg-primary-50 dark:hover:bg-primary-900/20" @click="startEdit(article)">编辑</button>
                <button class="rounded-md px-3 py-1.5 text-xs font-medium text-red-500 transition-colors hover:bg-red-50 dark:hover:bg-red-900/20" :disabled="deletingId === article.id" @click="handleDelete(article)">{{ deletingId === article.id ? '删除中...' : '删除' }}</button>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="py-20 text-center">
          <p class="text-gray-400 dark:text-gray-500">{{ searchKeyword ? '没有找到匹配的文章' : '暂无文章' }}</p>
        </div>
      </div>

      <!-- ========== Tab 3: 个人资料 ========== -->
      <div v-if="activeTab === 'profile'">
        <h2 class="mb-6 text-xl font-bold text-gray-900 dark:text-gray-100">个人资料</h2>
        <div class="space-y-8">
          <div class="flex flex-col items-center gap-4">
            <input ref="avatarFileInput" type="file" accept="image/*" class="hidden" @change="handleAvatarFileChange" />
            <div class="relative">
              <img v-if="profileForm.avatar" :src="toCdnUrl(profileForm.avatar)" alt="头像" class="h-24 w-24 rounded-full border-2 border-gray-200 object-cover dark:border-gray-600" />
              <div v-else class="flex h-24 w-24 items-center justify-center rounded-full border-2 border-gray-200 bg-gray-100 dark:border-gray-600 dark:bg-gray-700">
                <svg class="h-10 w-10 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>
              </div>
            </div>
            <button type="button" :disabled="avatarUploading" class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-60 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800" @click="avatarFileInput?.click()">{{ avatarUploading ? '上传中...' : '更换头像' }}</button>
            <div v-if="avatarUploading" class="w-32 mt-2">
              <div class="bg-gray-200 rounded-full h-1.5 dark:bg-gray-600">
                <div class="bg-primary-500 h-1.5 rounded-full transition-all duration-300" :style="{ width: avatarUploadPercent + '%' }" />
              </div>
              <p class="text-xs text-center text-gray-500 mt-1 dark:text-gray-400">{{ avatarUploadPercent }}%</p>
            </div>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">个人简介</label>
            <textarea v-model="profileForm.bio" rows="4" maxlength="200" placeholder="写点什么介绍自己吧..." class="w-full resize-none rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-primary-400" />
            <p class="mt-1 text-right text-xs text-gray-400 dark:text-gray-500">{{ profileForm.bio.length }} / 200</p>
          </div>
          <div class="flex items-center justify-end gap-3">
            <span v-if="profileSuccess" class="text-sm text-green-600 dark:text-green-400">✓ 保存成功</span>
            <span v-if="profileError" class="text-sm text-red-500">{{ profileError }}</span>
            <button type="button" :disabled="profileSaving" class="rounded-lg bg-primary-500 px-6 py-2.5 text-sm font-medium text-white transition-all duration-200 hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-300 disabled:cursor-not-allowed disabled:opacity-60" @click="handleSaveProfile">{{ profileSaving ? '保存中...' : '保存资料' }}</button>
          </div>
        </div>
      </div>

      <!-- ========== Tab 4: 相册管理 ========== -->
      <div v-if="activeTab === 'albums'">
        <div class="mb-6 flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">相册管理</h2>
          <button
            class="rounded-lg bg-primary-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-600"
            @click="openAlbumModal()"
          >
            + 新建相册集
          </button>
        </div>

        <!-- 加载中 -->
        <div v-if="adminAlbumsLoading" class="flex items-center justify-center py-20">
          <AppLoading tip="加载中..." />
        </div>

        <!-- 相册列表 -->
        <div v-else-if="adminAlbums.length > 0" class="space-y-3">
          <div
            v-for="album in adminAlbums"
            :key="album.id"
            class="rounded-lg border bg-white p-4 transition-colors dark:bg-gray-800"
            :class="adminSelectedAlbumId === album.id ? 'border-primary-300 dark:border-primary-600' : 'border-gray-100 dark:border-gray-700'"
          >
            <div class="flex items-center gap-4">
              <!-- 封面缩略图 -->
              <div class="h-10 w-10 flex-shrink-0 overflow-hidden rounded bg-gray-100 dark:bg-gray-700">
                <img v-if="album.coverUrl" :src="toCdnUrl(album.coverUrl)" :alt="album.name" class="h-full w-full object-cover" />
                <div v-else class="flex h-full w-full items-center justify-center">
                  <svg class="h-5 w-5 text-gray-300 dark:text-gray-500" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909M3.75 21h16.5a1.5 1.5 0 001.5-1.5V5.25a1.5 1.5 0 00-1.5-1.5H3.75a1.5 1.5 0 00-1.5 1.5v14.25a1.5 1.5 0 001.5 1.5z" />
                  </svg>
                </div>
              </div>
              <!-- 信息 -->
              <div class="min-w-0 flex-1">
                <h3 class="truncate text-sm font-semibold text-gray-900 dark:text-gray-100">
                  {{ album.name }}
                  <span v-if="album.hasPassword" class="ml-1 text-xs">🔒</span>
                </h3>
                <p class="text-xs text-gray-400 dark:text-gray-500">{{ album.photoCount }} 张照片</p>
              </div>
              <!-- 操作按钮 -->
              <div class="flex flex-shrink-0 items-center gap-2">
                <button
                  class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
                  :class="adminSelectedAlbumId === album.id ? 'bg-primary-500 text-white' : 'text-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20'"
                  @click="selectAdminAlbum(album)"
                >
                  {{ adminSelectedAlbumId === album.id ? '管理中' : '管理照片' }}
                </button>
                <button
                  class="rounded-md px-3 py-1.5 text-xs font-medium text-gray-500 transition-colors hover:bg-gray-50 dark:hover:bg-gray-700"
                  @click="openAlbumModal(album)"
                >
                  编辑
                </button>
                <button
                  class="rounded-md px-3 py-1.5 text-xs font-medium text-red-500 transition-colors hover:bg-red-50 dark:hover:bg-red-900/20"
                  :disabled="adminDeletingAlbumId === album.id"
                  @click="handleDeleteAlbum(album)"
                >
                  {{ adminDeletingAlbumId === album.id ? '删除中...' : '删除' }}
                </button>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="py-12 text-center">
          <p class="text-gray-400 dark:text-gray-500">暂无相册，点击上方按钮创建</p>
        </div>

        <!-- 照片管理区 -->
        <div v-if="adminSelectedAlbumId" class="mt-8">
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-lg font-bold text-gray-900 dark:text-gray-100">
              {{ adminSelectedAlbumName }} · 照片管理
            </h3>
            <div class="flex items-center gap-2">
              <input ref="photoFileInput" type="file" accept="image/*" multiple class="hidden" @change="handlePhotoUpload" />
              <button
                class="rounded-lg bg-primary-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-600"
                @click="photoFileInput?.click()"
              >
                上传照片
              </button>
            </div>
          </div>

          <!-- 加载中 -->
          <div v-if="adminPhotosLoading" class="flex items-center justify-center py-20">
            <AppLoading tip="加载中..." />
          </div>

          <!-- 照片网格 -->
          <div v-else-if="adminPhotos.length > 0" class="grid grid-cols-3 gap-3">
            <div
              v-for="photo in adminPhotos"
              :key="photo.id"
              class="group relative aspect-square overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800"
            >
              <!-- 上传中：loading 占位 -->
              <template v-if="photo.status === 'uploading'">
                <div class="absolute inset-0 flex flex-col items-center justify-center gap-2 bg-gray-100 dark:bg-gray-800">
                  <svg class="h-8 w-8 animate-spin text-primary-400" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  <!-- 进度条 -->
                  <div class="w-3/4">
                    <div class="h-1 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-gray-600">
                      <div
                        class="h-full rounded-full bg-primary-400 transition-all duration-200"
                        :style="{ width: photo.uploadPercent + '%' }"
                      />
                    </div>
                    <p class="mt-1 text-center text-xs text-gray-400">{{ photo.uploadPercent }}%</p>
                  </div>
                </div>
              </template>

              <!-- 上传失败：裂图 -->
              <template v-else-if="photo.status === 'error'">
                <div class="flex h-full w-full flex-col items-center justify-center gap-1 bg-gray-100 dark:bg-gray-800">
                  <!-- 裂图 SVG -->
                  <svg class="h-12 w-12 text-gray-300 dark:text-gray-600" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect width="64" height="64" rx="4" fill="currentColor" fill-opacity="0.15"/>
                    <!-- 山和太阳 -->
                    <path d="M8 48 L20 28 L28 38 L36 24 L56 48 Z" fill="currentColor" fill-opacity="0.25"/>
                    <circle cx="46" cy="18" r="7" fill="currentColor" fill-opacity="0.25"/>
                    <!-- 裂缝 -->
                    <path d="M22 10 L30 28 L24 36 L34 56" stroke="#EF4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M30 28 L38 22 L44 34" stroke="#EF4444" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <span class="text-xs text-red-400">上传失败</span>
                </div>
                <!-- 失败时也可以点 × 移除 -->
                <button
                  class="absolute right-1 top-1 flex h-6 w-6 items-center justify-center rounded-full bg-black/50 text-white opacity-0 transition-opacity group-hover:opacity-100"
                  @click="adminPhotos = adminPhotos.filter(p => p.id !== photo.id)"
                >
                  <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
              </template>

              <!-- 上传成功：正常显示图片 -->
              <template v-else>
                <img :src="toCdnUrl(photo.url)" :alt="photo.caption || '照片'" class="h-full w-full object-cover" />
                <!-- 锁图标（有密码时显示） -->
                <span v-if="photo.hasPassword" class="absolute top-1 left-1 rounded-full bg-black/50 px-1.5 py-0.5 text-xs backdrop-blur-sm">🔒</span>
                <!-- 设置密码按钮 -->
                <button
                  class="absolute left-1 bottom-1 flex h-6 items-center gap-0.5 rounded-full bg-black/50 px-2 text-xs text-white opacity-0 transition-opacity group-hover:opacity-100 backdrop-blur-sm"
                  @click="openPhotoPasswordModal(photo)"
                >
                  🔒 密码
                </button>
                <!-- 删除按钮 -->
                <button
                  class="absolute right-1 top-1 flex h-6 w-6 items-center justify-center rounded-full bg-black/50 text-white opacity-0 transition-opacity group-hover:opacity-100"
                  :disabled="adminDeletingPhotoId === photo.id"
                  @click="handleDeletePhoto(photo)"
                >
                  <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
              </template>
            </div>
          </div>
          <div v-else class="py-12 text-center">
            <p class="text-gray-400 dark:text-gray-500">暂无照片，点击上方按钮上传</p>
          </div>
        </div>
      </div>

      <!-- ========== Tab 5: 留言管理 ========== -->
      <div v-if="activeTab === 'messages'">
        <h2 class="mb-6 text-xl font-bold text-gray-900 dark:text-gray-100">留言管理</h2>

        <!-- 加载中 -->
        <div v-if="adminMsgLoading" class="flex items-center justify-center py-20">
          <AppLoading tip="加载中..." />
        </div>

        <!-- 留言列表 -->
        <div v-else-if="adminMessages.length > 0" class="space-y-4">
          <div
            v-for="msg in adminMessages"
            :key="msg.id"
            class="rounded-xl border border-gray-100 bg-white p-5 transition-colors dark:border-gray-700 dark:bg-gray-800"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0 flex-1">
                <div class="mb-2 flex items-center gap-2">
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ msg.nickname }}</span>
                  <time class="text-xs text-gray-400 dark:text-gray-500">{{ formatTime(msg.createdAt) }}</time>
                </div>
                <p class="whitespace-pre-wrap text-sm leading-relaxed text-gray-600 dark:text-gray-400">{{ msg.content }}</p>
              </div>
              <button
                class="shrink-0 rounded-lg px-3 py-1.5 text-xs font-medium text-red-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:text-red-400 dark:hover:bg-red-900/20 dark:hover:text-red-300"
                @click="handleDeleteMessage(msg)"
              >
                删除
              </button>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="py-20 text-center">
          <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
            <span class="text-2xl">💬</span>
          </div>
          <p class="text-lg font-medium text-gray-600 dark:text-gray-400">暂无留言</p>
        </div>
      </div>

      <!-- ========== Tab 6: 数据统计 ========== -->
      <div v-if="activeTab === 'analytics'">
        <h2 class="mb-6 text-xl font-bold text-gray-900 dark:text-gray-100">数据统计</h2>

        <!-- 概览卡片 -->
        <div class="mb-8 grid grid-cols-2 gap-3 md:grid-cols-4">
          <div class="rounded-xl border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">今日 PV</p>
            <p class="mt-1 text-2xl font-bold text-blue-500">{{ analyticsOverview ? analyticsOverview.today_pv : '--' }}</p>
          </div>
          <div class="rounded-xl border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">今日 UV</p>
            <p class="mt-1 text-2xl font-bold text-emerald-500">{{ analyticsOverview ? analyticsOverview.today_uv : '--' }}</p>
          </div>
          <div class="rounded-xl border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">总 PV</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-gray-100">{{ analyticsOverview ? analyticsOverview.total_pv.toLocaleString() : '--' }}</p>
          </div>
          <div class="rounded-xl border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">总 UV</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-gray-100">{{ analyticsOverview ? analyticsOverview.total_uv.toLocaleString() : '--' }}</p>
          </div>
        </div>

        <!-- PV/UV 趋势折线图 -->
        <div class="mb-8 rounded-xl border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">PV/UV 趋势</h3>
            <div class="flex gap-1">
              <button
                class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
                :class="trendDays === 7 ? 'bg-primary-500 text-white' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="trendDays = 7"
              >7 天</button>
              <button
                class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
                :class="trendDays === 30 ? 'bg-primary-500 text-white' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="trendDays = 30"
              >30 天</button>
            </div>
          </div>
          <!-- 图例 -->
          <div class="mb-3 flex items-center gap-4">
            <span class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400">
              <span class="inline-block h-2 w-4 rounded-sm bg-[#3b82f6]" /> PV
            </span>
            <span class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400">
              <span class="inline-block h-2 w-4 rounded-sm bg-[#10b981]" /> UV
            </span>
          </div>
          <!-- SVG 图表 -->
          <div class="overflow-x-auto">
          <div ref="chartContainerRef" class="relative" :style="{ width: chartWidth + 'px', minWidth: '320px', height: '280px' }">
            <svg v-if="trendData.length > 0" :width="chartWidth" :height="chartHeight">
              <!-- 纵轴网格线 + 标签 -->
              <template v-for="(val, i) in yAxisTicks" :key="'y' + i">
                <line
                  :x1="chartPadding.left" :y1="yScale(val)"
                  :x2="chartWidth - chartPadding.right" :y2="yScale(val)"
                  stroke="currentColor" class="text-gray-100 dark:text-gray-700" stroke-width="1"
                />
                <text
                  :x="chartPadding.left - 8" :y="yScale(val) + 4"
                  text-anchor="end" class="fill-gray-400 text-[10px] dark:fill-gray-500"
                >{{ val }}</text>
              </template>
              <!-- PV 折线 -->
              <polyline
                :points="pvPoints"
                fill="none" stroke="#3b82f6" stroke-width="2" stroke-linejoin="round" stroke-linecap="round"
              />
              <!-- UV 折线 -->
              <polyline
                :points="uvPoints"
                fill="none" stroke="#10b981" stroke-width="2" stroke-linejoin="round" stroke-linecap="round"
              />
              <!-- PV 数据点 -->
              <circle
                v-for="(pt, i) in pvCircles" :key="'pvc' + i"
                :cx="pt.x" :cy="pt.y" r="3.5"
                fill="#3b82f6" stroke="white" stroke-width="1.5"
                class="cursor-pointer dark:stroke-gray-800"
                @mouseenter="showTooltip($event, i)"
                @mouseleave="hideTooltip"
              />
              <!-- UV 数据点 -->
              <circle
                v-for="(pt, i) in uvCircles" :key="'uvc' + i"
                :cx="pt.x" :cy="pt.y" r="3.5"
                fill="#10b981" stroke="white" stroke-width="1.5"
                class="cursor-pointer dark:stroke-gray-800"
                @mouseenter="showTooltip($event, i)"
                @mouseleave="hideTooltip"
              />
              <!-- 横轴标签 -->
              <text
                v-for="(label, i) in xLabels" :key="'xl' + i"
                :x="xScale(i)"
                :y="chartHeight - 8"
                text-anchor="middle"
                class="fill-gray-400 text-[10px] dark:fill-gray-500"
              >{{ label }}</text>
            </svg>
            <div v-else class="flex h-full items-center justify-center">
              <p class="text-sm text-gray-400 dark:text-gray-500">暂无数据</p>
            </div>
            <!-- Tooltip -->
            <div
              v-if="tooltip.visible"
              class="pointer-events-none absolute z-10 rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs shadow-lg dark:border-gray-600 dark:bg-gray-700"
              :style="{ left: tooltip.x + 'px', top: tooltip.y + 'px' }"
            >
              <p class="font-medium text-gray-700 dark:text-gray-200">{{ tooltip.date }}</p>
              <p class="text-blue-500">PV: {{ tooltip.pv }}</p>
              <p class="text-emerald-500">UV: {{ tooltip.uv }}</p>
            </div>
          </div>
          </div>
        </div>

        <!-- Top5 页面报表 -->
        <div class="mb-8 rounded-xl border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">Top5 页面</h3>
            <div class="flex gap-1">
              <button
                class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
                :class="topPagesDays === 7 ? 'bg-primary-500 text-white' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="topPagesDays = 7"
              >7 天</button>
              <button
                class="rounded-md px-3 py-1 text-xs font-medium transition-colors"
                :class="topPagesDays === 30 ? 'bg-primary-500 text-white' : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
                @click="topPagesDays = 30"
              >30 天</button>
            </div>
          </div>
          <!-- PC 端表格 -->
          <div class="hidden md:block">
            <table class="w-full text-left text-sm">
              <thead>
                <tr class="border-b border-gray-100 dark:border-gray-700">
                  <th class="pb-2 font-medium text-gray-500 dark:text-gray-400">排名</th>
                  <th class="pb-2 font-medium text-gray-500 dark:text-gray-400">页面路径</th>
                  <th class="pb-2 text-right font-medium text-gray-500 dark:text-gray-400">PV</th>
                  <th class="pb-2 text-right font-medium text-gray-500 dark:text-gray-400">UV</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(page, i) in topPagesData" :key="page.path" class="border-b border-gray-50 dark:border-gray-700/50">
                  <td class="py-2.5 text-gray-900 dark:text-gray-100">
                    <span
                      class="inline-flex h-5 w-5 items-center justify-center rounded-full text-xs font-bold"
                      :class="i < 3 ? 'bg-primary-100 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400' : 'text-gray-400 dark:text-gray-500'"
                    >{{ i + 1 }}</span>
                  </td>
                  <td class="py-2.5 font-mono text-xs text-gray-700 dark:text-gray-300">{{ page.path }}</td>
                  <td class="py-2.5 text-right font-medium text-gray-900 dark:text-gray-100">{{ page.pv }}</td>
                  <td class="py-2.5 text-right font-medium text-gray-900 dark:text-gray-100">{{ page.uv }}</td>
                </tr>
                <tr v-if="topPagesData.length === 0">
                  <td colspan="4" class="py-8 text-center text-gray-400 dark:text-gray-500">暂无数据</td>
                </tr>
              </tbody>
            </table>
          </div>
          <!-- 移动端卡片 -->
          <div class="space-y-2 md:hidden">
            <div
              v-for="(page, i) in topPagesData" :key="page.path"
              class="flex items-center gap-3 rounded-lg bg-gray-50 p-3 dark:bg-gray-700/50"
            >
              <span
                class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-bold"
                :class="i < 3 ? 'bg-primary-100 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400' : 'bg-gray-100 text-gray-400 dark:bg-gray-600 dark:text-gray-500'"
              >{{ i + 1 }}</span>
              <div class="min-w-0 flex-1">
                <p class="truncate font-mono text-xs text-gray-700 dark:text-gray-300">{{ page.path }}</p>
                <p class="mt-0.5 text-xs text-gray-400 dark:text-gray-500">PV {{ page.pv }} · UV {{ page.uv }}</p>
              </div>
            </div>
            <div v-if="topPagesData.length === 0" class="py-8 text-center text-sm text-gray-400 dark:text-gray-500">暂无数据</div>
          </div>
        </div>

        <!-- 访问日志 -->
        <div class="rounded-xl border border-gray-100 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
          <h3 class="mb-4 text-sm font-semibold text-gray-900 dark:text-gray-100">访问日志</h3>
          <!-- 筛选条件 -->
          <div class="mb-4 flex flex-wrap items-end gap-3">
            <div class="min-w-0 flex-1">
              <label class="mb-1 block text-xs text-gray-500 dark:text-gray-400">页面路径</label>
              <input
                v-model="logFilter.path"
                type="text"
                placeholder="/home"
                class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              />
            </div>
            <div class="w-36">
              <label class="mb-1 block text-xs text-gray-500 dark:text-gray-400">日期</label>
              <input
                v-model="logFilter.date"
                type="date"
                class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              />
            </div>
            <div class="flex gap-2">
              <button
                class="rounded-lg bg-primary-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-600"
                @click="fetchLogs(1)"
              >查询</button>
              <button
                class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                @click="resetLogFilter"
              >重置</button>
            </div>
          </div>
          <!-- PC 端表格 -->
          <div class="hidden md:block">
            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm">
                <thead>
                  <tr class="border-b border-gray-100 dark:border-gray-700">
                    <th class="whitespace-nowrap pb-2 pr-4 font-medium text-gray-500 dark:text-gray-400">时间</th>
                    <th class="whitespace-nowrap pb-2 pr-4 font-medium text-gray-500 dark:text-gray-400">路径</th>
                    <th class="whitespace-nowrap pb-2 pr-4 font-medium text-gray-500 dark:text-gray-400">设备</th>
                    <th class="whitespace-nowrap pb-2 pr-4 font-medium text-gray-500 dark:text-gray-400">浏览器</th>
                    <th class="whitespace-nowrap pb-2 pr-4 font-medium text-gray-500 dark:text-gray-400">系统</th>
                    <th class="whitespace-nowrap pb-2 pr-4 font-medium text-gray-500 dark:text-gray-400">IP</th>
                    <th class="whitespace-nowrap pb-2 font-medium text-gray-500 dark:text-gray-400">来源</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="log in logsData.list" :key="log.id" class="border-b border-gray-50 dark:border-gray-700/50">
                    <td class="whitespace-nowrap py-2 pr-4 text-xs text-gray-500 dark:text-gray-400">{{ formatTime(log.created_at) }}</td>
                    <td class="max-w-[160px] truncate py-2 pr-4 font-mono text-xs text-gray-700 dark:text-gray-300">{{ log.path }}</td>
                    <td class="whitespace-nowrap py-2 pr-4 text-xs text-gray-500 dark:text-gray-400">{{ log.device_type || '-' }}</td>
                    <td class="whitespace-nowrap py-2 pr-4 text-xs text-gray-500 dark:text-gray-400">{{ log.browser || '-' }}</td>
                    <td class="whitespace-nowrap py-2 pr-4 text-xs text-gray-500 dark:text-gray-400">{{ log.os || '-' }}</td>
                    <td class="whitespace-nowrap py-2 pr-4 font-mono text-xs text-gray-500 dark:text-gray-400">{{ log.ip || '-' }}</td>
                    <td class="max-w-[120px] truncate py-2 text-xs text-gray-400 dark:text-gray-500">{{ log.referer || '-' }}</td>
                  </tr>
                  <tr v-if="logsData.list.length === 0">
                    <td colspan="7" class="py-8 text-center text-gray-400 dark:text-gray-500">暂无日志</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <!-- 移动端卡片列表 -->
          <div class="space-y-2 md:hidden">
            <div
              v-for="log in logsData.list" :key="log.id"
              class="rounded-lg bg-gray-50 p-3 dark:bg-gray-700/50"
            >
              <div class="mb-1.5 flex items-center justify-between">
                <span class="font-mono text-xs font-medium text-gray-700 dark:text-gray-300">{{ log.path }}</span>
                <span class="text-xs text-gray-400 dark:text-gray-500">{{ formatTime(log.created_at) }}</span>
              </div>
              <div class="flex flex-wrap gap-x-3 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
                <span v-if="log.device_type">{{ log.device_type }}</span>
                <span v-if="log.browser">{{ log.browser }}</span>
                <span v-if="log.os">{{ log.os }}</span>
                <span v-if="log.ip" class="font-mono">{{ log.ip }}</span>
              </div>
            </div>
            <div v-if="logsData.list.length === 0" class="py-8 text-center text-sm text-gray-400 dark:text-gray-500">暂无日志</div>
          </div>
          <!-- 分页 -->
          <div v-if="logsData.total > 0" class="mt-4 flex items-center justify-between">
            <span class="text-xs text-gray-400 dark:text-gray-500">
              第 {{ logsData.page }} / {{ Math.ceil(logsData.total / logsData.page_size) }} 页，共 {{ logsData.total }} 条
            </span>
            <div class="flex gap-2">
              <button
                :disabled="logsData.page <= 1"
                class="rounded-md border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                @click="fetchLogs(logsData.page - 1)"
              >上一页</button>
              <button
                :disabled="logsData.page >= Math.ceil(logsData.total / logsData.page_size)"
                class="rounded-md border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                @click="fetchLogs(logsData.page + 1)"
              >下一页</button>
            </div>
          </div>
        </div>
      </div>
      </div>
    </main>

    <!-- ========== 新建/编辑相册 Modal ========== -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="albumModal.visible"
          class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50"
          @click.self="closeAlbumModal"
        >
          <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-gray-800">
            <h3 class="mb-4 text-lg font-bold text-gray-900 dark:text-gray-100">
              {{ albumModal.editingId ? '编辑相册' : '新建相册' }}
            </h3>
            <form @submit.prevent="handleSaveAlbum">
              <div class="space-y-4">
                <div>
                  <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">相册名称</label>
                  <input
                    v-model="albumModal.name"
                    type="text"
                    required
                    placeholder="输入相册名称"
                    class="w-full rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                  />
                </div>
                <div>
                  <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">描述（可选）</label>
                  <textarea
                    v-model="albumModal.description"
                    rows="3"
                    placeholder="输入相册描述"
                    class="w-full resize-none rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                  />
                </div>
                <div>
                  <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">密码保护（可选）</label>
                  <input
                    v-model="albumModal.password"
                    type="password"
                    :placeholder="albumModal.editingId && albumModal.hasPassword ? '输入新密码可修改，留空不变' : '留空则不加密'"
                    autocomplete="new-password"
                    class="w-full rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                  />
                  <p v-if="albumModal.editingId && albumModal.hasPassword" class="mt-1 text-xs text-amber-500 dark:text-amber-400">
                    🔒 已设置密码。输入新密码可修改，留空则保持不变。
                  </p>
                </div>
              </div>
              <div class="mt-6 flex justify-end gap-3">
                <button
                  type="button"
                  class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                  @click="closeAlbumModal"
                >
                  取消
                </button>
                <button
                  type="submit"
                  :disabled="albumModal.saving || !albumModal.name.trim()"
                  class="rounded-lg bg-primary-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {{ albumModal.saving ? '保存中...' : '确定' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- ========== 照片密码设置 Modal ========== -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="photoPasswordModal.visible"
          class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50"
          @click.self="closePhotoPasswordModal"
        >
          <div class="w-full max-w-sm rounded-xl bg-white p-6 shadow-xl dark:bg-gray-800 mx-4">
            <h3 class="mb-1 text-lg font-bold text-gray-900 dark:text-gray-100">🔒 设置照片密码</h3>
            <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">
              {{ photoPasswordModal.hasPassword ? '该照片已设置密码，可修改或清除' : '设置密码后需验证才能查看' }}
            </p>
            <form @submit.prevent="handleSavePhotoPassword">
              <input
                v-model="photoPasswordModal.password"
                type="password"
                :placeholder="photoPasswordModal.hasPassword ? '输入新密码，留空则清除密码' : '输入密码，留空则不加密'"
                autocomplete="new-password"
                class="w-full rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
              />
              <div class="mt-4 flex justify-end gap-3">
                <button
                  type="button"
                  class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                  @click="closePhotoPasswordModal"
                >
                  取消
                </button>
                <button
                  type="submit"
                  :disabled="photoPasswordModal.saving"
                  class="rounded-lg bg-primary-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {{ photoPasswordModal.saving ? '保存中...' : '确定' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import { showSuccess, showError, showConfirm } from '~/utils/ui'
import type { ArticleListItem, AlbumItem, PhotoItem, MessageItem } from '~/types'
import { useAuthStore } from '~/stores/auth'
import { chunkedUpload } from '~/utils/chunkedUpload'
import {
  apiCreateArticle,
  apiUpdateArticle,
  apiDeleteArticle,
  apiGetArticles,
  apiFetchArticle,
  apiUploadImage,
  apiGetProfile,
  apiUpdateProfile,
  apiGetAlbums,
  apiCreateAlbum,
  apiUpdateAlbum,
  apiDeleteAlbum,
  apiGetPhotos,
  apiAddPhoto,
  apiDeletePhoto,
  apiUpdatePhoto,
  apiGetMessages,
  apiDeleteMessage,
} from '~/utils/api'
import { toCdnUrl } from '~/utils/imageUrl'
interface AdminPhotoItem {
  id: number           // 正式照片为正数 DB id；占位项为负数临时 id
  albumId: number
  url: string          // 占位时为空字符串，成功后为真实 url
  caption: string | null
  hasPassword: boolean
  createdAt: string
  updatedAt: string
  // 扩展字段（仅本地使用）
  status: 'done' | 'uploading' | 'error'
  uploadPercent: number  // 0-100，仅 uploading 时有效
}

const authStore = useAuthStore()
const router = useRouter()
const { isDark, toggleTheme } = useTheme()

// ====== Tab 管理 ======
const adminTabs = [
  { key: 'write', label: '写文章' },
  { key: 'manage', label: '文章管理' },
  { key: 'profile', label: '个人资料' },
  { key: 'albums', label: '相册管理' },
  { key: 'messages', label: '留言管理' },
  { key: 'analytics', label: '数据统计' },
]
const activeTab = ref('write')

// ====== 写文章 / 编辑文章 ======
const editingArticleId = ref<string | null>(null)
const coverFileInput = ref<HTMLInputElement | null>(null)
const coverUploading = ref(false)
const coverUploadPercent = ref(0)

const form = reactive({
  title: '',
  summary: '',
  coverImage: '',
})

const publishing = ref(false)
const publishSuccess = ref(false)
const publishError = ref('')

const editor = useEditor({
  extensions: [
    StarterKit,
    Image,
    Link.configure({
      openOnClick: false,
      HTMLAttributes: { target: '_blank', rel: 'noopener noreferrer' },
    }),
  ],
  content: '<p>开始写作...</p>',
  editorProps: {
    attributes: { class: 'prose prose-sm max-w-none' },
  },
})

function addImage() {
  const url = window.prompt('请输入图片 URL')
  if (url && editor.value) editor.value.chain().focus().setImage({ src: url }).run()
}

function addLink() {
  const url = window.prompt('请输入链接 URL')
  if (url && editor.value) editor.value.chain().focus().setLink({ href: url }).run()
}

async function handleCoverUpload(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  await uploadCoverFile(file)
  input.value = ''
}

async function handleCoverDrop(e: DragEvent) {
  const file = e.dataTransfer?.files?.[0]
  if (!file || !file.type.startsWith('image/')) return
  await uploadCoverFile(file)
}

async function uploadCoverFile(file: File) {
  coverUploading.value = true
  coverUploadPercent.value = 0
  try {
    const url = await chunkedUpload(file, (p) => { coverUploadPercent.value = p.percent })
    form.coverImage = url
  } catch { showError('封面图上传失败，请重试') }
  finally { coverUploading.value = false; coverUploadPercent.value = 0 }
}

function resetForm() {
  editingArticleId.value = null
  form.title = ''
  form.summary = ''
  form.coverImage = ''
  editor.value?.commands.setContent('<p>开始写作...</p>')
}

function cancelEdit() { resetForm() }

async function handlePublish() {
  if (!editor.value || !form.title.trim()) return
  publishing.value = true
  publishSuccess.value = false
  publishError.value = ''
  try {
    const content = editor.value.getJSON()
    const payload = {
      title: form.title.trim(),
      summary: form.summary.trim() || undefined,
      coverImage: form.coverImage.trim() || undefined,
      content,
    }
    if (editingArticleId.value) {
      await apiUpdateArticle(editingArticleId.value, payload)
      publishSuccess.value = true
      showSuccess('文章保存成功')
      setTimeout(() => { resetForm(); activeTab.value = 'manage'; publishSuccess.value = false; fetchManageArticles() }, 1500)
    } else {
      await apiCreateArticle(payload)
      publishSuccess.value = true
      showSuccess('文章发布成功')
      resetForm()
      setTimeout(() => { publishSuccess.value = false }, 3000)
    }
  } catch (err: unknown) {
    const fetchErr = err as { data?: { statusMessage?: string }; statusCode?: number }
    if (fetchErr?.statusCode === 401 || (fetchErr?.data as Record<string, unknown>)?.statusCode === 401) {
      publishError.value = '登录已过期，请重新登录'
      showError('登录已过期，请重新登录')
      authStore.setLoggedIn(false)
      setTimeout(() => router.push('/login'), 1500)
    } else {
      const errMsg = editingArticleId.value ? '保存失败，请重试' : '发布失败，请重试'
      publishError.value = errMsg
      showError(errMsg)
    }
  } finally { publishing.value = false }
}

// ====== 文章管理 ======
const searchKeyword = ref('')
const manageArticles = ref<ArticleListItem[]>([])
const manageLoading = ref(false)
const deletingId = ref<string | null>(null)
let searchTimer: ReturnType<typeof setTimeout> | null = null

async function fetchManageArticles() {
  manageLoading.value = true
  try {
    const params = searchKeyword.value.trim() ? { title: searchKeyword.value.trim() } : undefined
    const res = await apiGetArticles(params)
    manageArticles.value = res.list
  } catch { manageArticles.value = [] }
  finally { manageLoading.value = false }
}

watch(searchKeyword, () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchManageArticles(), 300)
})

watch(activeTab, (val) => {
  if (val === 'manage') fetchManageArticles()
  if (val === 'profile') fetchProfile()
  if (val === 'albums') fetchAdminAlbums()
  if (val === 'messages') fetchAdminMessages()
  if (val === 'analytics') fetchAnalyticsData()
})

async function startEdit(article: ArticleListItem) {
  try {
    const detail = await apiFetchArticle(article.id)
    editingArticleId.value = detail.id
    form.title = detail.title
    form.summary = detail.summary || ''
    form.coverImage = detail.coverImage || ''
    editor.value?.commands.setContent(detail.content)
    activeTab.value = 'write'
  } catch { showError('加载文章详情失败，请重试') }
}

async function handleDelete(article: ArticleListItem) {
  showConfirm({
    title: '确认删除',
    content: `确定要删除「${article.title}」吗？此操作不可撤销。`,
    okText: '确定删除',
    danger: true,
    cancelText: '取消',
    async onOk() {
      deletingId.value = article.id
      try {
        await apiDeleteArticle(article.id)
        showSuccess('文章删除成功')
        await fetchManageArticles()
      } catch (err: unknown) {
        const fetchErr = err as { statusCode?: number }
        if (fetchErr?.statusCode === 401) { showError('登录已过期，请重新登录'); authStore.setLoggedIn(false); router.push('/login') }
        else showError('删除失败，请重试')
      } finally { deletingId.value = null }
    },
  })
}

// ====== 个人资料 ======
const avatarFileInput = ref<HTMLInputElement | null>(null)
const profileForm = reactive({ avatar: '', bio: '' })
const profileSaving = ref(false)
const profileSuccess = ref(false)
const profileError = ref('')
const avatarUploading = ref(false)
const avatarUploadPercent = ref(0)

async function fetchProfile() {
  try {
    const data = await apiGetProfile()
    profileForm.avatar = data.avatar || ''
    profileForm.bio = data.bio || ''
  } catch { /* 静默处理 */ }
}

async function handleAvatarFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  avatarUploading.value = true
  avatarUploadPercent.value = 0
  try {
    const url = await chunkedUpload(file, (p) => { avatarUploadPercent.value = p.percent })
    profileForm.avatar = url
    showSuccess('头像上传成功')
  } catch { showError('头像上传失败，请重试') }
  finally { avatarUploading.value = false; avatarUploadPercent.value = 0; input.value = '' }
}

async function handleSaveProfile() {
  profileSaving.value = true
  profileSuccess.value = false
  profileError.value = ''
  try {
    await apiUpdateProfile({ avatar: profileForm.avatar, bio: profileForm.bio })
    profileSuccess.value = true
    showSuccess('个人信息保存成功')
    setTimeout(() => { profileSuccess.value = false }, 3000)
  } catch (err: unknown) {
    const fetchErr = err as { statusCode?: number }
    if (fetchErr?.statusCode === 401) {
      profileError.value = '登录已过期，请重新登录'
      showError('登录已过期，请重新登录')
      authStore.setLoggedIn(false)
      setTimeout(() => router.push('/login'), 1500)
    } else {
      profileError.value = '保存失败，请重试'
      showError('保存失败，请重试')
    }
  } finally { profileSaving.value = false }
}

// ====== 相册管理 ======
const adminAlbums = ref<AlbumItem[]>([])
const adminAlbumsLoading = ref(false)
const adminDeletingAlbumId = ref<number | null>(null)

const adminSelectedAlbumId = ref<number | null>(null)
const adminPhotos = ref<AdminPhotoItem[]>([])
const adminPhotosLoading = ref(false)
const adminDeletingPhotoId = ref<number | null>(null)

const photoFileInput = ref<HTMLInputElement | null>(null)

const adminSelectedAlbumName = computed(() => {
  return adminAlbums.value.find(a => a.id === adminSelectedAlbumId.value)?.name || ''
})

// 相册 Modal
const albumModal = reactive({
  visible: false,
  editingId: null as number | null,
  name: '',
  description: '',
  password: '',
  hasPassword: false, // 编辑时原相册是否已设密码
  saving: false,
})

async function fetchAdminAlbums() {
  adminAlbumsLoading.value = true
  try {
    const res = await apiGetAlbums()
    adminAlbums.value = res.list
  } catch { adminAlbums.value = [] }
  finally { adminAlbumsLoading.value = false }
}

async function selectAdminAlbum(album: AlbumItem) {
  if (adminSelectedAlbumId.value === album.id) {
    adminSelectedAlbumId.value = null
    adminPhotos.value = []
    return
  }
  adminSelectedAlbumId.value = album.id
  adminPhotosLoading.value = true
  try {
    const res = await apiGetPhotos(album.id)
    adminPhotos.value = res.list.map(p => ({ ...p, status: 'done' as const, uploadPercent: 100, hasPassword: p.hasPassword ?? false }))
  } catch { adminPhotos.value = [] }
  finally { adminPhotosLoading.value = false }
}

function openAlbumModal(album?: AlbumItem) {
  albumModal.editingId = album?.id ?? null
  albumModal.name = album?.name ?? ''
  albumModal.description = album?.description ?? ''
  albumModal.password = ''
  albumModal.hasPassword = album?.hasPassword ?? false
  albumModal.saving = false
  albumModal.visible = true
}

function closeAlbumModal() {
  albumModal.visible = false
}

async function handleSaveAlbum() {
  if (!albumModal.name.trim()) return
  albumModal.saving = true
  try {
    if (albumModal.editingId) {
      await apiUpdateAlbum(albumModal.editingId, {
        name: albumModal.name.trim(),
        description: albumModal.description.trim() || undefined,
        password: albumModal.password,
      })
    } else {
      await apiCreateAlbum({
        name: albumModal.name.trim(),
        description: albumModal.description.trim() || undefined,
        password: albumModal.password || undefined,
      })
    }
    closeAlbumModal()
    showSuccess(albumModal.editingId ? '相册更新成功' : '相册创建成功')
    await fetchAdminAlbums()
  } catch (err: unknown) {
    const fetchErr = err as { statusCode?: number }
    if (fetchErr?.statusCode === 401) {
      showError('登录已过期，请重新登录')
      authStore.setLoggedIn(false)
      router.push('/login')
    } else {
      showError('保存失败，请重试')
    }
  } finally { albumModal.saving = false }
}

async function handleDeleteAlbum(album: AlbumItem) {
  showConfirm({
    title: '确认删除',
    content: `确定要删除相册「${album.name}」吗？其中所有照片也会被删除，此操作不可撤销。`,
    okText: '确定删除',
    danger: true,
    cancelText: '取消',
    async onOk() {
      adminDeletingAlbumId.value = album.id
      try {
        await apiDeleteAlbum(album.id)
        if (adminSelectedAlbumId.value === album.id) {
          adminSelectedAlbumId.value = null
          adminPhotos.value = []
        }
        showSuccess('相册删除成功')
        await fetchAdminAlbums()
      } catch (err: unknown) {
        const fetchErr = err as { statusCode?: number }
        if (fetchErr?.statusCode === 401) {
          showError('登录已过期，请重新登录')
          authStore.setLoggedIn(false)
          router.push('/login')
        } else { showError('删除失败，请重试') }
      } finally { adminDeletingAlbumId.value = null }
    },
  })
}

/** 批量上传照片（乐观更新：先插入占位卡，成功后替换，失败后显示裂图） */
async function handlePhotoUpload(e: Event) {
  const input = e.target as HTMLInputElement
  const files = input.files
  if (!files || files.length === 0 || !adminSelectedAlbumId.value) return

  const albumId = adminSelectedAlbumId.value

  // 为每张图片生成临时 id，并插入占位卡（插到列表最前面）
  const tempItems: AdminPhotoItem[] = Array.from(files).map((_, i) => ({
    id: -(Date.now() + i),   // 负数 id，避免与真实 id 冲突
    albumId,
    url: '',
    caption: null,
    hasPassword: false,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    status: 'uploading' as const,
    uploadPercent: 0,
  }))

  // 插入到列表最前面
  adminPhotos.value = [...tempItems, ...adminPhotos.value]

  // 并发上传所有图片（每张独立处理，互不影响）
  const uploadTasks = Array.from(files).map(async (file, i) => {
    const tempId = tempItems[i].id

    try {
      // 上传文件，获取 url
      const url = await chunkedUpload(file, (p) => {
        const idx = adminPhotos.value.findIndex(ph => ph.id === tempId)
        if (idx !== -1) adminPhotos.value[idx].uploadPercent = p.percent
      })

      // 调接口存入 DB
      const photo = await apiAddPhoto(albumId, { url })

      // 替换占位项为真实数据
      const idx = adminPhotos.value.findIndex(ph => ph.id === tempId)
      if (idx !== -1) {
        adminPhotos.value[idx] = { ...photo, status: 'done' as const, uploadPercent: 100, hasPassword: photo.hasPassword ?? false }
      }
    } catch (err: unknown) {
      const fetchErr = err as { statusCode?: number }
      if (fetchErr?.statusCode === 401) {
        showError('登录已过期，请重新登录')
        authStore.setLoggedIn(false)
        router.push('/login')
        return
      }
      // 标记为失败
      const idx = adminPhotos.value.findIndex(ph => ph.id === tempId)
      if (idx !== -1) adminPhotos.value[idx].status = 'error'
    }
  })

  await Promise.all(uploadTasks)

  input.value = ''
  await fetchAdminAlbums() // 刷新 photoCount 和封面
}

// ====== 照片密码设置 Modal ======
const photoPasswordModal = reactive({
  visible: false,
  photoId: 0,
  password: '',
  hasPassword: false,
  saving: false,
})

function openPhotoPasswordModal(photo: AdminPhotoItem) {
  photoPasswordModal.photoId = photo.id
  photoPasswordModal.password = ''
  photoPasswordModal.hasPassword = photo.hasPassword
  photoPasswordModal.saving = false
  photoPasswordModal.visible = true
}

function closePhotoPasswordModal() {
  photoPasswordModal.visible = false
  photoPasswordModal.password = ''
}

async function handleSavePhotoPassword() {
  photoPasswordModal.saving = true
  try {
    const updated = await apiUpdatePhoto(photoPasswordModal.photoId, {
      password: photoPasswordModal.password,
    })
    // 更新本地列表中的 hasPassword 状态
    const idx = adminPhotos.value.findIndex(p => p.id === photoPasswordModal.photoId)
    if (idx !== -1) {
      adminPhotos.value[idx].hasPassword = updated.hasPassword
    }
    closePhotoPasswordModal()
    showSuccess(photoPasswordModal.password ? '密码设置成功' : '密码已清除')
  } catch (err: unknown) {
    const fetchErr = err as { statusCode?: number }
    if (fetchErr?.statusCode === 401) {
      showError('登录已过期，请重新登录')
      authStore.setLoggedIn(false)
      router.push('/login')
    } else {
      showError('保存失败，请重试')
    }
  } finally {
    photoPasswordModal.saving = false
  }
}

async function handleDeletePhoto(photo: AdminPhotoItem) {
  showConfirm({
    title: '确认删除',
    content: '确定要删除这张照片吗？',
    okText: '确定删除',
    danger: true,
    cancelText: '取消',
    async onOk() {
      adminDeletingPhotoId.value = photo.id
      try {
        await apiDeletePhoto(photo.id)
        adminPhotos.value = adminPhotos.value.filter(p => p.id !== photo.id)
        showSuccess('照片删除成功')
        await fetchAdminAlbums() // 刷新 photoCount
      } catch (err: unknown) {
        const fetchErr = err as { statusCode?: number }
        if (fetchErr?.statusCode === 401) {
          showError('登录已过期，请重新登录')
          authStore.setLoggedIn(false)
          router.push('/login')
        } else { showError('删除失败，请重试') }
      } finally { adminDeletingPhotoId.value = null }
    },
  })
}

// ====== 留言管理 ======
const adminMessages = ref<MessageItem[]>([])
const adminMsgLoading = ref(false)

async function fetchAdminMessages() {
  adminMsgLoading.value = true
  try {
    const res = await apiGetMessages()
    adminMessages.value = res.list
  } catch {
    showError('获取留言列表失败')
  } finally {
    adminMsgLoading.value = false
  }
}

function handleDeleteMessage(msg: MessageItem) {
  showConfirm({
    title: '删除留言',
    content: `确定要删除「${msg.nickname}」的这条留言吗？`,
    danger: true,
    okText: '删除',
    async onOk() {
      try {
        await apiDeleteMessage(msg.id)
        adminMessages.value = adminMessages.value.filter(m => m.id !== msg.id)
        showSuccess('删除成功')
      } catch (err: any) {
        if (err?.response?.status === 401) {
          showError('登录已过期，请重新登录')
          authStore.setLoggedIn(false)
          router.push('/login')
        } else if (err?.response?.status === 404) {
          showError('留言不存在')
          adminMessages.value = adminMessages.value.filter(m => m.id !== msg.id)
        } else {
          showError('删除失败，请重试')
        }
      }
    },
  })
}

function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

async function handleLogout() {
  try {
    await $fetch('/api/auth/logout', { method: 'POST' })
  } catch {}
  authStore.logout()
  router.push('/home')
}

// ====== 数据统计 ======
interface AnalyticsOverview {
  today_pv: number
  today_uv: number
  total_pv: number
  total_uv: number
}
interface TrendItem {
  date: string
  pv: number
  uv: number
}
interface TopPageItem {
  path: string
  pv: number
  uv: number
}
interface LogItem {
  id: number
  path: string
  device_id: string
  ip: string
  device_type: string
  browser: string
  os: string
  referer: string
  created_at: string
}

const analyticsOverview = ref<AnalyticsOverview | null>(null)
const trendData = ref<TrendItem[]>([])
const trendDays = ref(7)
const topPagesData = ref<TopPageItem[]>([])
const topPagesDays = ref(7)

const logFilter = reactive({ path: '', date: '' })
const logsData = reactive({
  list: [] as LogItem[],
  total: 0,
  page: 1,
  page_size: 5,
})

// 折线图相关
const chartContainerRef = ref<HTMLElement | null>(null)
const chartWidth = ref(600)
const chartHeight = 280
const chartPadding = { top: 20, right: 20, bottom: 50, left: 50 }

const tooltip = reactive({ visible: false, x: 0, y: 0, date: '', pv: 0, uv: 0 })

let chartResizeObserver: ResizeObserver | null = null

onMounted(() => {
  if (chartContainerRef.value) {
    const updateWidth = () => {
      const containerEl = chartContainerRef.value?.parentElement
      if (!containerEl) return
      const containerW = containerEl.clientWidth
      chartWidth.value = trendDays.value <= 7 ? Math.max(containerW, 320) : Math.max(containerW, 600)
    }
    updateWidth()
    chartResizeObserver = new ResizeObserver(updateWidth)
    chartResizeObserver.observe(chartContainerRef.value.parentElement!)
  }
})

// 纵轴刻度
const yMax = computed(() => {
  if (trendData.value.length === 0) return 5
  const max = Math.max(...trendData.value.map(d => Math.max(d.pv, d.uv)))
  return max < 5 ? 5 : Math.ceil(max * 1.2)
})
const yAxisTicks = computed(() => {
  const ticks: number[] = []
  const step = Math.ceil(yMax.value / 5)
  for (let i = 0; i <= 5; i++) ticks.push(step * i)
  return ticks
})

// 坐标映射
function xScale(i: number): number {
  const total = trendData.value.length
  if (total <= 1) return chartPadding.left + (chartWidth.value - chartPadding.left - chartPadding.right) / 2
  const drawW = chartWidth.value - chartPadding.left - chartPadding.right
  return chartPadding.left + (i / (total - 1)) * drawW
}
function yScale(val: number): number {
  const drawH = chartHeight - chartPadding.top - chartPadding.bottom
  return chartPadding.top + drawH - (val / yMax.value) * drawH
}

const pvCircles = computed(() => trendData.value.map((d, i) => ({ x: xScale(i), y: yScale(d.pv) })))
const uvCircles = computed(() => trendData.value.map((d, i) => ({ x: xScale(i), y: yScale(d.uv) })))
const pvPoints = computed(() => pvCircles.value.map(p => `${p.x},${p.y}`).join(' '))
const uvPoints = computed(() => uvCircles.value.map(p => `${p.x},${p.y}`).join(' '))

const xLabels = computed(() => {
  return trendData.value.map((d, i) => {
    const short = d.date.slice(5) // "04-01"
    if (trendDays.value <= 7) return short
    return i % 5 === 0 || i === trendData.value.length - 1 ? short : ''
  })
})

function showTooltip(e: MouseEvent, idx: number) {
  const d = trendData.value[idx]
  if (!d || !chartContainerRef.value) return
  const rect = chartContainerRef.value.getBoundingClientRect()
  tooltip.x = e.clientX - rect.left + 12
  tooltip.y = e.clientY - rect.top - 60
  tooltip.date = d.date
  tooltip.pv = d.pv
  tooltip.uv = d.uv
  tooltip.visible = true
}
function hideTooltip() {
  tooltip.visible = false
}

// 数据请求
async function fetchOverview() {
  try {
    const res = await $fetch<{ ok: boolean; data: AnalyticsOverview }>('/api/pv/overview')
    analyticsOverview.value = res.data
  } catch { /* 静默 */ }
}

async function fetchTrend() {
  try {
    const res = await $fetch<{ ok: boolean; data: TrendItem[] }>('/api/pv/trend', {
      params: { days: trendDays.value },
    })
    trendData.value = res.data || []
  } catch { trendData.value = [] }
}

async function fetchTopPages() {
  try {
    const res = await $fetch<{ ok: boolean; data: TopPageItem[] }>('/api/pv/top-pages', {
      params: { days: topPagesDays.value },
    })
    topPagesData.value = res.data || []
  } catch { topPagesData.value = [] }
}

async function fetchLogs(page: number = 1) {
  try {
    const res = await $fetch<{ ok: boolean; data: { list: LogItem[]; total: number; page: number; page_size: number } }>('/api/pv/logs', {
      params: {
        page,
        page_size: logsData.page_size,
        path: logFilter.path || undefined,
        date: logFilter.date || undefined,
      },
    })
    logsData.list = res.data.list || []
    logsData.total = res.data.total
    logsData.page = res.data.page
    logsData.page_size = res.data.page_size
  } catch {
    logsData.list = []
    logsData.total = 0
  }
}

function resetLogFilter() {
  logFilter.path = ''
  logFilter.date = ''
  fetchLogs(1)
}

async function fetchAnalyticsData() {
  await Promise.all([fetchOverview(), fetchTrend(), fetchTopPages(), fetchLogs(1)])
}

watch(trendDays, () => {
  const containerEl = chartContainerRef.value?.parentElement
  if (containerEl) {
    chartWidth.value = trendDays.value <= 7 ? Math.max(containerEl.clientWidth, 320) : Math.max(containerEl.clientWidth, 600)
  }
  fetchTrend()
})
watch(topPagesDays, () => fetchTopPages())

onUnmounted(() => {
  editor.value?.destroy()
  if (searchTimer) clearTimeout(searchTimer)
  chartResizeObserver?.disconnect()
})

useHead({ title: '管理后台 - fatwill' })
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
