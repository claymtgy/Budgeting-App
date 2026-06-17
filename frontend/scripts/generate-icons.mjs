import sharp from 'sharp'
import { readFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const iconsDir = join(__dirname, '../public/icons')
const svg = readFileSync(join(iconsDir, 'icon.svg'))

async function writePng(name, size, paddingScale = 1) {
  const inner = Math.round(size * paddingScale)
  const offset = Math.round((size - inner) / 2)
  const resized = await sharp(svg).resize(inner, inner).png().toBuffer()

  await sharp({
    create: {
      width: size,
      height: size,
      channels: 4,
      background: { r: 26, g: 26, b: 46, alpha: 1 }
    }
  })
    .composite([{ input: resized, left: offset, top: offset }])
    .png()
    .toFile(join(iconsDir, name))

  console.log(`wrote ${name} (${size}x${size})`)
}

await writePng('icon-192.png', 192)
await writePng('icon-512.png', 512)
await writePng('icon-maskable-512.png', 512, 0.72)
await writePng('apple-touch-icon.png', 180)
