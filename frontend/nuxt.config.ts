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
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
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
