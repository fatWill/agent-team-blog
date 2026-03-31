import { promises as fs } from 'node:fs'
import { resolve } from 'node:path'

const DATA_FILE = resolve(process.cwd(), 'server/data/changelog.json')

export default defineEventHandler(async () => {
  const raw = await fs.readFile(DATA_FILE, 'utf-8')
  return JSON.parse(raw)
})
