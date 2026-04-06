/**
 * 分片上传工具 re-export（向后兼容层）
 * 实际实现已移至 shared/utils/chunkedUpload.ts
 */
export type { UploadProgress, UploadResult } from '~/shared/utils/chunkedUpload'
export { chunkedUpload } from '~/shared/utils/chunkedUpload'
