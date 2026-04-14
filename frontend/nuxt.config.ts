// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-03-25',

  runtimeConfig: {
    // Go 后端地址（同一台服务器内网直连）
    backendUrl: process.env.BACKEND_URL || 'http://127.0.0.1:8080',
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
    cdnURL: 'https://cdn.fatwill.cloud',
    head: {
      title: 'fatwill - 个人博客',
      htmlAttrs: {
        lang: 'zh-CN',
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '全栈开发者 fatwill 的个人博客，分享技术与生活。' },
        { name: 'author', content: 'fatwill' },
        { property: 'og:site_name', content: 'fatwill' },
        { property: 'og:type', content: 'website' },
        { name: 'twitter:card', content: 'summary' },
        { name: 'robots', content: 'index, follow' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'canonical', href: 'https://fatwill.cloud' },
        { rel: 'sitemap', type: 'application/xml', href: '/sitemap.xml' },
        { rel: 'preconnect', href: 'https://assets.fatwill.cloud', crossorigin: '' },
        { rel: 'dns-prefetch', href: 'https://assets.fatwill.cloud' },
        { rel: 'preconnect', href: 'https://cdn.fatwill.cloud', crossorigin: '' },
        { rel: 'dns-prefetch', href: 'https://cdn.fatwill.cloud' },
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
