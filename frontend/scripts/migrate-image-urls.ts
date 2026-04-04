/**
 * 数据库旧图片 URL 迁移脚本
 *
 * 将数据库中 /uploads/xxx 格式的 URL 更新为 https://assets.fatwill.cloud/uploads/xxx
 * 注意：旧图片文件仍在 /root/blog-uploads/，需要单独上传到 COS
 * 这里只更新 URL 前缀，让旧图片通过 assets.fatwill.cloud 访问
 *
 * 使用方式（在服务器上运行）：
 *   npx tsx scripts/migrate-image-urls.ts
 */

import mysql from 'mysql2/promise'

async function migrate() {
  const conn = await mysql.createConnection({
    host: process.env.DB_HOST || 'localhost',
    user: process.env.DB_USER || 'root',
    password: process.env.DB_PASSWORD || '',
    database: process.env.DB_NAME || 'blog',
  })

  const oldPrefix = '/uploads/'
  const newPrefix = 'https://assets.fatwill.cloud/uploads/'

  console.log('开始迁移数据库中的图片 URL ...')
  console.log(`旧前缀: ${oldPrefix}`)
  console.log(`新前缀: ${newPrefix}`)
  console.log('---')

  // 更新 articles 表的 cover_image
  const [r1] = await conn.execute(
    `UPDATE articles SET cover_image = REPLACE(cover_image, ?, ?) WHERE cover_image LIKE ?`,
    [oldPrefix, newPrefix, '/uploads/%'],
  )
  console.log('articles.cover_image:', r1)

  // 更新 photos 表的 url
  const [r2] = await conn.execute(
    `UPDATE photos SET url = REPLACE(url, ?, ?) WHERE url LIKE ?`,
    [oldPrefix, newPrefix, '/uploads/%'],
  )
  console.log('photos.url:', r2)

  // 更新 profiles 表的 avatar
  const [r3] = await conn.execute(
    `UPDATE profiles SET avatar = REPLACE(avatar, ?, ?) WHERE avatar LIKE ?`,
    [oldPrefix, newPrefix, '/uploads/%'],
  )
  console.log('profiles.avatar:', r3)

  await conn.end()
  console.log('---')
  console.log('✅ 迁移完成')
}

migrate().catch(console.error)
