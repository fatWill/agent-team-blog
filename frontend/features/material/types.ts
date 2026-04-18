export interface MaterialAttachment {
  id: string
  type: 'pdf' | 'image' | 'video'
  url: string
  name: string
  size?: number
  thumbnailUrl?: string
  duration?: number
}

export interface MaterialItem {
  id: number
  title: string
  tags: string[]
  attachments: MaterialAttachment[]
  sortOrder: number
  createdAt: string
  updatedAt: string
}
