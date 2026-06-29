declare namespace API {
  type AddUserRequest = {
    userAccount: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }

  type BaseResponse = {
    code?: number
    data?: any
    message?: string
  }

  type DeleteRequest = {
    id: number
  }

  type getUserGetParams = {
    /** 用户 ID */
    id: number
  }

  type LoginRequest = {
    userAccount: string
    userPassword: string
  }

  type LoginUser = {
    createTime?: string
    editTime?: string
    id?: number
    quota?: number
    updateTime?: string
    userAccount?: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
    vipTime?: string
  }

  type PageResult = {
    pageNum?: number
    pageSize?: number
    records?: any
    total?: number
  }

  type QueryUserRequest = {
    id?: number
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    userAccount?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }

  type RegisterRequest = {
    checkPassword: string
    userAccount: string
    userPassword: string
  }

  type UpdateUserRequest = {
    id: number
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }

  type User = {
    createTime?: string
    editTime?: string
    id?: number
    updateTime?: string
    userAccount?: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }

  type UserInfo = {
    createTime?: string
    editTime?: string
    id?: number
    updateTime?: string
    userAccount?: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
    vipTime?: string
  }
}
