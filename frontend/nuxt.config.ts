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
      // 为 home 页面添加 Tab 路径别名
      const homePage = pages.find(p => p.path === '/home')
      if (homePage) {
        homePage.alias = ['/articles', '/life', '/tools', '/agent-team', '/guestbook', '/changelog']
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
