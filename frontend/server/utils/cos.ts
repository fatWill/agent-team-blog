import COS from 'cos-nodejs-sdk-v5'

let _cosClient: COS | null = null

function getCOSClient(): COS {
  if (!_cosClient) {
    _cosClient = new COS({
      SecretId: process.env.COS_ID || '',
      SecretKey: process.env.COS_KEY || '',
    })
  }
  return _cosClient
}

/**
 * 生成 COS 存储路径
 * 格式：upload/YYYYMMDD/时间戳-随机串.ext
 */
export function generateCOSKey(ext: string): string {
  const now = new Date()
  const date = `${now.getFullYear()}${String(now.getMonth() + 1).padStart(2, '0')}${String(now.getDate()).padStart(2, '0')}`
  const randomStr = Math.random().toString(36).slice(2, 10)
  return `upload/${date}/${Date.now()}-${randomStr}${ext}`
}

/**
 * 上传 Buffer 到 COS
 * @returns 完整的访问 URL
 */
export async function uploadToCOS(key: string, buffer: Buffer, contentType?: string): Promise<string> {
  const client = getCOSClient()
  const bucket = process.env.COS_BUCKET || 'fatwill-cloud-1253664788'
  const region = process.env.COS_REGION || 'ap-guangzhou'
  const baseURL = process.env.COS_BASE_URL || 'https://assets.fatwill.cloud'

  await new Promise<void>((resolve, reject) => {
    client.putObject({
      Bucket: bucket,
      Region: region,
      Key: key,
      Body: buffer,
      ContentType: contentType || 'application/octet-stream',
    }, (err) => {
      if (err) reject(err)
      else resolve()
    })
  })

  return `${baseURL}/${key}`
}
