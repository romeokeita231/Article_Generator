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
