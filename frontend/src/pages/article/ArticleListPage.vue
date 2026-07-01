<template>

</template>
<script lang="ts">
import { getArticleTaskId, postArticleList, postArticleOpenApiDelete } from '@/api/articleHandler'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import dayjs, { type Dayjs } from 'dayjs'
import { message } from 'ant-design-vue'

const router = useRouter()

// 搜索筛选
const searchKeyword = ref('')
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const statusFilter = ref<string>('')

const columns = [
  {
    title: '选题',
    dataIndex: 'topic',
    key: 'topic',
    width: 180,
    ellipsis: true,
  },
  {
    title: '标题',
    key: 'title',
    width: 280,
  },
  {
    title: '状态',
    key: 'status',
    width: 110,
  },
  {
    title: '创建时间',
    key: 'createTime',
    width: 160,
  },
  {
    title: '操作',
    key: 'action',
    width: 200,
  },
]

const loading = ref(false)
const dataSource = ref<API.ArticleInfo[]>([])
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条`,
  pageSizeOptions: ['10', '20', '50', '100']
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await postArticleList({
      pageNum: pagination.value.current,
      pageSize: pagination.value.pageSize,
    })
    const pageData = res.data.data
    let records = pageData?.records || []

    // 前端过滤（如果后端不支持）
    if (searchKeyword.value) {
      const keyword = searchKeyword.value.toLowerCase()
      records = records.filter((item: API.ArticleInfo) =>
        item.mainTitle?.toLowerCase().includes(keyword) ||
        item.topic?.toLowerCase().includes(keyword)
      )
    }

    if (statusFilter.value) {
      records = records.filter((item: API.ArticleInfo) => item.status === statusFilter.value)
    }

    if (dateRange.value) {
      const [start, end] = dateRange.value
      records = records.filter((item: API.ArticleInfo) => {
        const createTime = dayjs(item.createTime)
        return createTime.isAfter(start.startOf('day')) && createTime.isBefore(end.endOf('day'))
      })
    }

    dataSource.value = records
    pagination.value.total = pageData?.totalRow || 0
  } catch (error: any) {
    message.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// 导出文章
const exportArticle = async (record: API.ArticleInfo) => {
  try {
    const res = await getArticleTaskId({ taskId: record.taskId! })
    const article = res.data.data
    if (!article) {
      message.error('文章数据不存在')
      return
    }

    let markdown = `# ${article.mainTitle}\n\n`
    markdown += `> ${article.subTitle}\n\n`

    if (article.fullContent) {
      markdown += article.fullContent
    } else {
      markdown += article.content || ''
    }

    const blob = new Blob([markdown], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${article.mainTitle || '文章'}.md`
    a.click()
    URL.revokeObjectURL(url)

    message.success('导出成功')
  } catch (error: any) {
    message.error(error.message || '导出失败')
  }
}

// 删除文章
const deleteArticle = async (record: API.ArticleInfo) => {
  try {
    await postArticleOpenApiDelete({ id: record.id! })
    message.success('删除成功')
    loadData()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}


</script>
