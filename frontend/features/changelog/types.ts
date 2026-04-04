/** 单个版本更新日志 */
export interface ChangelogItem {
  version: string
  date: string
  logs: string[]
  createdAt: string
  updatedAt: string
}

/** 更新日志接口响应 */
export interface ChangelogResponse {
  changelog: ChangelogItem[]
}
