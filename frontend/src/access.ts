import { useLoginUserStore } from '@/stores/loginUser'
import { message } from 'ant-design-vue'
import router from '@/router'
import { USER_ROLE_ADMIN } from '@/constants/user'

// 是否为首次获取登录用户
let firstFetchLoginUser = true

/**
 * 全局权限校验
 * 使用新写法：返回 true/false 或路径字符串，不使用 next()
 */
router.beforeEach(async (to, from) => {
  const loginUserStore = useLoginUserStore()
  let loginUser = loginUserStore.loginUser

  // 首次加载时，等后端返回用户信息后再校验权限
  if (firstFetchLoginUser) {
    await loginUserStore.fetchLoginUser()
    loginUser = loginUserStore.loginUser
    firstFetchLoginUser = false
  }

  const toUrl = to.fullPath

  // 1. 如果访问的是登录页或注册页，直接放行（不需要登录）
  if (toUrl.startsWith('/user/login') || toUrl.startsWith('/user/register')) {
    return true
  }

  // 2. 如果访问的是首页，直接放行
  if (toUrl === '/' || toUrl === '') {
    return true
  }

  // 3. 管理员页面权限校验
  if (toUrl.startsWith('/admin')) {
    if (!loginUser || loginUser.userRole !== USER_ROLE_ADMIN) {
      message.error('没有权限')
      return `/user/login?redirect=${to.fullPath}`
    }
  }

  // 4. 其他页面需要登录
  // 如果未登录，跳转到登录页
  if (!loginUser || !loginUser.id) {
    return `/user/login?redirect=${to.fullPath}`
  }

  return true
})
