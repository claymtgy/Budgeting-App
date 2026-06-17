import sharp from 'sharp'
import { readFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const iconsDir = join(__dirname, '../public/icons')

async function writePng(name, size, svgFile = 'icon.svg') {
  const svg = readFileSync(join(iconsDir, svgFile))
  await sharp(svg, { density: 300 })
    .resize(size, size, { fit: 'contain', background: '#1a1a2e' })
    .png()
    .toFile(join(iconsDir, name))

  console.log(`wrote ${name} (${size}x${size}) from ${svgFile}`)
}

await writePng('icon-192.png', 192)
await writePng('icon-512.png', 512)
await writePng('icon-maskable-512.png', 512, 'icon-maskable.svg')
await writePng('apple-touch-icon.png', 180)
