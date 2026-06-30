// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** 创建用户（管理员） POST /user/add */
export async function postUserAdd(body: API.AddUserRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponse & { data?: number }>('/user/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 删除用户（管理员） POST /user/delete */
export async function postUserOpenApiDelete(
  body: API.DeleteRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: boolean }>('/user/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 根据 ID 获取用户信息 GET /user/get */
export async function getUserGet(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getUserGetParams,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: API.UserInfo }>('/user/get', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  })
}

/** 分页查询用户列表（管理员） POST /user/list */
export async function postUserList(body: API.QueryUserRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponse & { data?: API.PageResult }>('/user/list', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取当前登录用户 GET /user/login */
export async function getUserLogin(options?: { [key: string]: any }) {
  return request<API.BaseResponse & { data?: API.LoginUser }>('/user/login', {
    method: 'GET',
    ...(options || {}),
  })
}

/** 用户登录 POST /user/login */
export async function postUserLogin(body: API.LoginRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponse & { data?: API.LoginUser }>('/user/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 用户注销 POST /user/logout */
export async function postUserLogout(options?: { [key: string]: any }) {
  return request<API.BaseResponse & { data?: boolean }>('/user/logout', {
    method: 'POST',
    ...(options || {}),
  })
}

/** 用户注册 POST /user/register */
export async function postUserRegister(
  body: API.RegisterRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: number }>('/user/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 更新用户（管理员） POST /user/update */
export async function postUserUpdate(
  body: API.UpdateUserRequest,
  options?: { [key: string]: any }
) {
  return request<API.BaseResponse & { data?: boolean }>('/user/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}
