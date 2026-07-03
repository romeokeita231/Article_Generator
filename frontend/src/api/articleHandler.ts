// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 获取文章 GET /article/${param0} */
export async function getArticleTaskId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getArticleTaskIdParams,
  options?: { [key: string]: any }
) {
  const { taskId: param0, ...queryParams } = params
  return request<API.BaseResponse & { data?: API.ArticleInfo }>(`/article/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  })
}

/** 使用 AI 修改大纲 POST /article/aiModifyOutline */
export async function postArticleAiModifyOutline(
  body: API.AiModifyOutlineRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: API.OutlineSection }>('/article/aiModifyOutline', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 确认大纲 POST /article/confirmOutline */
export async function postArticleConfirmOutline(
  body: API.ConfirmOutlineRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse>('/article/confirmOutline', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 确认标题并输入补充描述 POST /article/confirmTitle */
export async function postArticleConfirmTitle(
  body: API.ConfirmTitleRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse>('/article/confirmTitle', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 创建文章 POST /article/create */
export async function postArticleCreate(
  body: API.CreateArticleRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: number }>('/article/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 删除文章 POST /article/delete */
export async function postArticleOpenApiDelete(
  body: API.DeleteRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: boolean }>('/article/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取任务执行日志 GET /article/execution-logs/${param0} */
export async function getArticleExecutionLogsTaskId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getArticleExecutionLogsTaskIdParams,
  options?: { [key: string]: any }
) {
  const { taskId: param0, ...queryParams } = params
  return request<API.BaseResponse & { data?: API.AgentExecutionStats }>(
    `/article/execution-logs/${param0}`,
    {
      method: 'GET',
      params: { ...queryParams },
      ...(options || {}),
    }
  )
}

/** 分页查询文章列表 POST /article/list */
export async function postArticleList(
  body: API.QueryArticleRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: API.ArticlePage }>('/article/list', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取文章进度 GET /article/progress/${param0} */
export async function getArticleProgressTaskId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getArticleProgressTaskIdParams,
  options?: { [key: string]: any }
) {
  const { taskId: param0, ...queryParams } = params
  return request<string>(`/article/progress/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  })
}
