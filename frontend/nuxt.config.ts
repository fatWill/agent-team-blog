// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-03-25',

  runtimeConfig: {
    // SQLite 数据库路径（从 .env 自动读取）
    dbPath: process.env.DB_PATH || '/root/blog-data/blog.db',
    // Redis 配置
    redisHost: process.env.REDIS_HOST || '127.0.0.1',
    redisPort: process.env.REDIS_PORT || '6379',
    redisPassword: process.env.REDIS_PASSWORD || '',
  },

  nitro: {
    maxRequestBodySize: 500 * 1024 * 1024, // 500MB
    // better-sqlite3 是原生模块，不打包进 bundle，作为外部依赖在服务器上安装
    externals: {
      external: ['better-sqlite3'],
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
