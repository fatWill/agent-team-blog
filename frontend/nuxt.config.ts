// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-03-25',

  runtimeConfig: {
    // Go 后端地址（同一台服务器内网直连）
    backendUrl: process.env.BACKEND_URL || 'http://127.0.0.1:8080',

    public: {
      // 站点主域名（用于 canonical / OG / sitemap / robots 等绝对 URL 拼接）
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'https://fatwill.cloud',
      // 图片/静态素材 CDN（腾讯云 COS 自定义域名）
      assetsUrl: process.env.NUXT_PUBLIC_ASSETS_URL || 'https://assets.fatwill.cloud',
      // Nuxt 构建产物 CDN
      cdnUrl: process.env.NUXT_PUBLIC_CDN_URL || 'https://cdn.fatwill.cloud',
    },
  },

  nitro: {
    maxRequestBodySize: 500 * 1024 * 1024, // 500MB
  },

  hooks: {
    'pages:extend'(pages) {
      const { resolve } = require('path')
      // 为每个 Tab 路径注册独立路由，指向 home.vue（SSR 可正确识别）
      const tabPaths = [
        { name: 'tab-articles', path: '/articles' },
        { name: 'tab-life', path: '/life' },
        { name: 'tab-renovation', path: '/renovation' },
        { name: 'tab-toys', path: '/toys' },
        { name: 'tab-agent-team', path: '/agent-team' },
        { name: 'tab-guestbook', path: '/guestbook' },
        { name: 'tab-changelog', path: '/changelog' },
      ]
      for (const tab of tabPaths) {
        pages.push({
          name: tab.name,
          path: tab.path,
          file: resolve(__dirname, 'pages/home.vue'),
        })
      }
    },
  },

  devtools: { enabled: false },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
  ],

  app: {
    // ⚠️ Nuxt 构建期常量：cdnURL 在 build 时被内联进产物路径，Nuxt 官方不支持运行时注入。
    // 如需换域名，必须在「构建时」设置 NUXT_PUBLIC_CDN_URL 并重新 build，
    // 仅改运行时环境变量对本项无效。同步参考 .env.example。
    cdnURL: process.env.NUXT_PUBLIC_CDN_URL || 'https://cdn.fatwill.cloud',
    head: {
      title: 'fatwill 的小屋',
      htmlAttrs: {
        lang: 'zh-CN',
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'fatwill 的小屋 — 分享技术与生活' },
        { name: 'author', content: 'fatwill' },
        { property: 'og:site_name', content: 'fatwill 的小屋' },
        { property: 'og:type', content: 'website' },
        { name: 'twitter:card', content: 'summary' },
        { name: 'robots', content: 'index, follow' },
      ],
      // 注意：canonical / preconnect / dns-prefetch 已迁移到 app.vue，
      // 通过 useHead + useRuntimeConfig 运行时动态注入（支持 NUXT_PUBLIC_SITE_URL /
      // NUXT_PUBLIC_ASSETS_URL / NUXT_PUBLIC_CDN_URL 换域名）。此处只保留与域名无关的静态项。
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
        { rel: 'manifest', href: '/site.webmanifest' },
        { rel: 'sitemap', type: 'application/xml', href: '/sitemap.xml' },
      ],
    },
  },

  tailwindcss: {
    cssPath: '~/assets/styles/tailwind.css',
    configPath: 'tailwind.config.ts',
  },

  typescript: {
    strict: true,
  },
})
