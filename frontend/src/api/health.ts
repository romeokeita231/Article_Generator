// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 健康检查 检查服务是否可用 GET /health */
export async function getHealth(options?: { [key: string]: any }) {
  return request<string>('/health', {
    method: 'GET',
    ...(options || {}),
  })
}
