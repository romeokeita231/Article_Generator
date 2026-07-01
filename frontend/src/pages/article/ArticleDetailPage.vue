<template>
  <!-- 大纲 -->
  <div v-if="article?.outline && article.outline.length > 0" class="outline-section">
    <h2 class="section-title">
      <OrderedListOutlined class="section-icon" />
      文章大纲
    </h2>

    <div class="outline-list">
      <div v-for="item in article.outline" :key="item.section" class="outline-item">
        <div class="outline-title">{{ item.section }}. {{ item.title }}</div>

        <ul class="outline-points">
          <li v-for="(point, idx) in item.points" :key="idx">{{ point }}</li>

        </ul>

      </div>

    </div>

  </div>

  <!-- 完整图文（优先展示） -->
  <div v-if="article?.fullContent" class="content-section">
    <h2 class="section-title">
      <FileTextOutlined class="section-icon" />
      完整图文
    </h2>

    <div v-html="markdownToHtml(article.fullContent)" class="markdown-content"></div>

  </div>

  <!-- 普通正文（无 fullContent 时展示） -->
  <div v-else-if="article?.content" class="content-section">
    <h2 class="section-title">
      <FileTextOutlined class="section-icon" />
      文章正文
    </h2>

    <div v-html="markdownToHtml(article.content)" class="markdown-content"></div>

  </div>

</template>
<script lang="ts" setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { getArticleTaskId } from '@/api/articleHandler'
import { marked } from 'marked'




const route = useRoute()

const loading = ref(false)
const article = ref<API.ArticleInfo | null>(null)

// Markdown 转 HTML
const markdownToHtml = (markdown: string) => {
  return marked(markdown)
}

// 加载文章
const loadArticle = async () => {
  const taskId = route.params.taskId as string
  if (!taskId) {
    message.error('文章ID不存在')
    return
  }

  loading.value = true
  try {
    const res = await getArticleTaskId({ taskId })
    article.value = res.data.data
  } catch (error: any) {
    message.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// 导出 Markdown
const exportMarkdown = () => {
  if (!article.value) return

  let markdown = `# ${article.value.mainTitle}\n\n`
  markdown += `> ${article.value.subTitle}\n\n`

  // 优先使用完整图文
  if (article.value.fullContent) {
    markdown += article.value.fullContent
  } else {
    if (article.value.outline && article.value.outline.length > 0) {
      markdown += `## 目录\n\n`
      article.value.outline.forEach(item => {
        markdown += `${item.section}. ${item.title}\n`
      })
      markdown += `\n---\n\n`
    }

    markdown += article.value.content || ''

    if (article.value.images && article.value.images.length > 0) {
      markdown += `\n\n## 配图\n\n`
      article.value.images.forEach(image => {
        markdown += `![${image.description}](${image.url})\n\n`
      })
    }
  }

  const blob = new Blob([markdown], { type: 'text/markdown' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${article.value.mainTitle}.md`
  a.click()
  URL.revokeObjectURL(url)

  message.success('导出成功')
}


</script>
