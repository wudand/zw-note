function escapeHtml(text: string): string {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

function splitTableRow(line: string): string[] {
  let trimmed = line.trim()
  if (trimmed.startsWith('|')) trimmed = trimmed.slice(1)
  if (trimmed.endsWith('|')) trimmed = trimmed.slice(0, -1)
  return trimmed.split('|').map((cell) => cell.trim())
}

function isTableSeparator(line: string): boolean {
  const cells = splitTableRow(line)
  if (cells.length === 0) return false
  return cells.every((cell) => /^:?-{3,}:?$/.test(cell))
}

function isTableRow(line: string): boolean {
  const trimmed = line.trim()
  if (!trimmed.includes('|')) return false
  if (trimmed.startsWith('{{')) return false
  return splitTableRow(trimmed).length > 0
}

function parseColumnAlign(sep: string): 'left' | 'center' | 'right' {
  const cell = sep.trim()
  const left = cell.startsWith(':')
  const right = cell.endsWith(':')
  if (left && right) return 'center'
  if (right) return 'right'
  return 'left'
}

function alignAttr(align: 'left' | 'center' | 'right'): string {
  return align === 'left' ? '' : ` style="text-align: ${align}"`
}

/**
 * 将后端返回的站点相对图片路径（如 `/uploads/images/xxx.png`）解析为可访问地址。
 *
 * 内容里始终存相对路径（便于跨环境迁移/备份），只在渲染展示时按需补全：
 * - 前后端同域部署（`VITE_API_URL` 为相对路径，如 `/api`）：原样返回，交由反向代理转发 `/uploads`
 * - 前后端跨域部署（`VITE_API_URL` 为绝对地址）：自动补全为后端源
 */
export function resolveAssetUrl(url: string): string {
  if (!url) return url
  const trimmed = url.trim()
  if (/^([a-z][\w+.-]*:)?\/\//i.test(trimmed) || trimmed.startsWith('data:')) return trimmed
  if (!trimmed.startsWith('/uploads/')) return trimmed

  const apiBase = import.meta.env.VITE_API_URL ?? ''
  if (/^https?:\/\//i.test(apiBase)) {
    try {
      return new URL(apiBase).origin + trimmed
    } catch {
      return trimmed
    }
  }
  return trimmed
}

/** 块级 HTML / 占位符：不应再包进 <p>；行内标签（strong/em/a/code 等）需要包 <p> 才能各自换行 */
function isBlockLevelLine(trimmed: string): boolean {
  if (
    trimmed.startsWith('{{MARKDOWN-CODE-BLOCK-') ||
    trimmed.startsWith('{{ALIGN-OPEN-') ||
    trimmed === '{{ALIGN-CLOSE}}' ||
    trimmed.includes('{{ALIGN-OPEN-') ||
    trimmed.includes('{{ALIGN-CLOSE}}')
  ) {
    return true
  }

  return /^(<\/?(?:h[1-6]|hr|ul|ol|li|table|thead|tbody|tr|blockquote|div|pre)\b|<hr\s*\/?>)/i.test(
    trimmed,
  )
}

/** 若当前行起是 GFM 表格，返回 HTML 与消费到的下一行下标；否则返回 null */
function tryConsumeTable(
  lines: string[],
  startIndex: number,
): { html: string; nextIndex: number } | null {
  const headerLine = lines[startIndex]?.trim()
  const separatorLine = lines[startIndex + 1]?.trim()
  if (!headerLine || !separatorLine) return null
  if (!isTableRow(headerLine) || !isTableSeparator(separatorLine)) return null

  const headers = splitTableRow(headerLine)
  const separators = splitTableRow(separatorLine)
  if (headers.length === 0) return null

  const aligns = headers.map((_, idx) => parseColumnAlign(separators[idx] ?? '---'))

  const bodyRows: string[][] = []
  let i = startIndex + 2
  while (i < lines.length) {
    const rowLine = lines[i]?.trim() ?? ''
    if (!rowLine || !isTableRow(rowLine) || isTableSeparator(rowLine)) break
    bodyRows.push(splitTableRow(rowLine))
    i++
  }

  let html = '<table><thead><tr>'
  headers.forEach((header, idx) => {
    html += `<th${alignAttr(aligns[idx] ?? 'left')}>${header}</th>`
  })
  html += '</tr></thead><tbody>'

  for (const row of bodyRows) {
    html += '<tr>'
    headers.forEach((_, idx) => {
      html += `<td${alignAttr(aligns[idx] ?? 'left')}>${row[idx] ?? ''}</td>`
    })
    html += '</tr>'
  }
  html += '</tbody></table>'

  return { html, nextIndex: i }
}

/** 将 Markdown 转为可安全用于 v-html 的 HTML */
export function renderMarkdown(markdown: string): string {
  if (!markdown) return ''

  let html = markdown

  const codeBlocks: string[] = []
  html = html.replace(/```[\s\S]*?```/g, (match) => {
    const index = codeBlocks.length
    codeBlocks.push(match)
    return `{{MARKDOWN-CODE-BLOCK-${index}}}`
  })

  html = html.replace(
    /<div\s+align="(left|center|right|justify)"\s*>/gi,
    '{{ALIGN-OPEN-$1}}',
  )
  html = html.replace(/<\/div>/gi, '{{ALIGN-CLOSE}}')

  html = html
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  html = html.replace(/^###### (.*$)/gim, '<h6>$1</h6>')
  html = html.replace(/^##### (.*$)/gim, '<h5>$1</h5>')
  html = html.replace(/^#### (.*$)/gim, '<h4>$1</h4>')
  html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>')

  html = html.replace(/^---$/gim, '<hr>')
  html = html.replace(/^\*\*\*$/gim, '<hr>')
  html = html.replace(/^___$/gim, '<hr>')

  html = html.replace(/`([^`]+)`/g, '<code>$1</code>')
  html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/__([^_]+)__/g, '<strong>$1</strong>')
  html = html.replace(/(?<!\*)\*([^*]+)\*(?!\*)/g, '<em>$1</em>')
  html = html.replace(/(?<!_)_([^_]+)_(?!_)/g, '<em>$1</em>')
  html = html.replace(/~~([^~]+)~~/g, '<del>$1</del>')
  // 图片语法 `![alt](url)` 本身也满足链接正则 `[alt](url)`，必须先处理图片、
  // 再处理链接，否则链接正则会先把方括号部分吃掉，图片永远渲染不出来
  html = html.replace(
    /!\[([^\]]*)\]\(([^)]+)\)/g,
    (_match, alt: string, src: string) => `<img src="${resolveAssetUrl(src)}" alt="${alt}">`,
  )
  html = html.replace(
    /(?<!!)\[([^\]]+)\]\(([^)]+)\)/g,
    '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>',
  )

  const lines = html.split('\n')
  const result: string[] = []
  let inList = false
  let listType: 'ul' | 'ol' | 'task' | null = null

  const closeList = () => {
    if (!inList) return
    result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
    inList = false
    listType = null
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (line === undefined) continue
    const trimmed = line.trim()

    if (!trimmed) {
      closeList()
      continue
    }

    if (trimmed === '<hr>') {
      closeList()
      result.push('<hr>')
      continue
    }

    if (trimmed.match(/^<h[1-6]>/)) {
      closeList()
      result.push(trimmed)
      continue
    }

    if (
      trimmed.startsWith('{{ALIGN-OPEN-') ||
      trimmed === '{{ALIGN-CLOSE}}' ||
      trimmed.includes('{{ALIGN-OPEN-') ||
      trimmed.includes('{{ALIGN-CLOSE}}')
    ) {
      closeList()
      result.push(trimmed)
      continue
    }

    const table = tryConsumeTable(lines, i)
    if (table) {
      closeList()
      result.push(table.html)
      i = table.nextIndex - 1
      continue
    }

    if (trimmed.match(/^[-*+] \[[ x]\] /i)) {
      if (!inList || listType !== 'task') {
        closeList()
        result.push('<ul class="task-list">')
        inList = true
        listType = 'task'
      }
      const item = trimmed.replace(/^[-*+] \[([ x])\] /i, '')
      const checked = trimmed.match(/\[[xX]\]/) ? 'checked' : ''
      result.push(
        `<li class="task-list-item"><input type="checkbox" ${checked} disabled> ${item}</li>`,
      )
      continue
    }

    if (trimmed.match(/^[-*+] /)) {
      if (!inList || listType !== 'ul') {
        closeList()
        result.push('<ul>')
        inList = true
        listType = 'ul'
      }
      result.push(`<li>${trimmed.replace(/^[-*+] /, '')}</li>`)
      continue
    }

    if (trimmed.match(/^\d+\. /)) {
      if (!inList || listType !== 'ol') {
        closeList()
        result.push('<ol>')
        inList = true
        listType = 'ol'
      }
      result.push(`<li>${trimmed.replace(/^\d+\. /, '')}</li>`)
      continue
    }

    if (trimmed.startsWith('&gt; ')) {
      closeList()
      result.push(`<blockquote>${trimmed.replace(/^&gt; /, '')}</blockquote>`)
      continue
    }

    if (trimmed.includes('{{MARKDOWN-CODE-BLOCK-')) {
      closeList()
      result.push(trimmed)
      continue
    }

    closeList()

    // 整行是 **粗体** / *斜体* / 链接 等时，会以 <strong>/<em>/<a> 开头；
    // 这些是行内元素，必须包进 <p>，否则预览会挤在同一行
    if (isBlockLevelLine(trimmed)) {
      result.push(trimmed)
    } else {
      result.push(`<p>${trimmed}</p>`)
    }
  }

  closeList()

  html = result.join('\n')

  html = html.replace(
    /{{ALIGN-OPEN-(left|center|right|justify)}}/g,
    '<div style="text-align: $1">',
  )
  html = html.replace(/{{ALIGN-CLOSE}}/g, '</div>')

  codeBlocks.forEach((codeBlock, index) => {
    const code = codeBlock.replace(/```[\w]*\n?/g, '').replace(/```/g, '').trim()
    html = html.replace(
      `{{MARKDOWN-CODE-BLOCK-${index}}}`,
      `<pre><code>${escapeHtml(code)}</code></pre>`,
    )
  })

  return html
}
