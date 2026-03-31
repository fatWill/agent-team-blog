/** 分片上传工具 —— 支持进度回调、断点续传（分片）、自动清理 */

const CHUNK_SIZE = 1.5 * 1024 * 1024 // 1.5MB 每片，预留 FormData 字段开销

export interface UploadProgress {
  loaded: number   // 已上传字节
  total: number    // 总字节
  percent: number  // 0-100
}

/**
 * 小文件直传（使用 fetch，确保 cookie 鉴权正常）
 */
async function uploadSmallFile(
  file: File,
  onProgress?: (progress: UploadProgress) => void,
): Promise<string> {
  onProgress?.({ loaded: 0, total: file.size, percent: 0 })

  const formData = new FormData()
  formData.append('file', file)

  const res = await fetch('/api/upload', {
    method: 'POST',
    credentials: 'include',
    body: formData,
  })

  if (!res.ok) {
    throw new Error(`上传失败，状态码 ${res.status}`)
  }

  const data = (await res.json()) as { url: string }
  onProgress?.({ loaded: file.size, total: file.size, percent: 100 })
  return data.url
}

/**
 * 清理临时分片文件
 */
async function cleanupChunks(uploadId: string): Promise<void> {
  try {
    await fetch('/api/upload/chunk', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ uploadId }),
    })
  } catch {
    // 清理失败静默处理
  }
}

/**
 * 分片上传大文件
 */
async function uploadLargeFile(
  file: File,
  onProgress?: (progress: UploadProgress) => void,
): Promise<string> {
  const uploadId = Date.now().toString(36) + Math.random().toString(36).slice(2)
  const totalChunks = Math.ceil(file.size / CHUNK_SIZE)
  const totalSize = file.size

  try {
    for (let i = 0; i < totalChunks; i++) {
      const start = i * CHUNK_SIZE
      const end = Math.min(start + CHUNK_SIZE, file.size)
      const chunk = file.slice(start, end)

      const formData = new FormData()
      formData.append('file', chunk)
      formData.append('uploadId', uploadId)
      formData.append('chunkIndex', String(i))
      formData.append('totalChunks', String(totalChunks))
      formData.append('filename', file.name)

      const res = await fetch('/api/upload/chunk', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      })

      if (!res.ok) {
        throw new Error(`分片 ${i} 上传失败，状态码 ${res.status}`)
      }

      // 更新进度：每完成一个分片更新
      const loaded = i === totalChunks - 1
        ? totalSize
        : (i + 1) * CHUNK_SIZE
      onProgress?.({
        loaded,
        total: totalSize,
        percent: Math.round((loaded / totalSize) * 100),
      })
    }

    // 所有分片上传完成，请求合并
    const mergeRes = await fetch('/api/upload/merge', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ uploadId, totalChunks, filename: file.name }),
    })

    if (!mergeRes.ok) {
      throw new Error(`合并分片失败，状态码 ${mergeRes.status}`)
    }

    const mergeData = (await mergeRes.json()) as { url: string }
    return mergeData.url
  } catch (err) {
    // 出错时清理临时文件
    await cleanupChunks(uploadId)
    throw err
  }
}

/**
 * 分片上传主函数
 * - 文件 ≤ 1.5MB：直接用 fetch 上传（开始/结束两段式进度）
 * - 文件 > 1.5MB：分片上传
 * @param file 要上传的文件
 * @param onProgress 进度回调
 * @returns 最终图片 URL
 */
export async function chunkedUpload(
  file: File,
  onProgress?: (progress: UploadProgress) => void,
): Promise<string> {
  if (file.size <= CHUNK_SIZE) {
    return uploadSmallFile(file, onProgress)
  }
  return uploadLargeFile(file, onProgress)
}
