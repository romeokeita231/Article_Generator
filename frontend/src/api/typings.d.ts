declare namespace API {
  type AddUserRequest = {
    userAccount: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }

  type AiModifyOutlineRequest = {
    modifySuggestion: string
    taskId: string
  }

  type ArticleInfo = {
    completedTime?: string
    content?: string
    createTime?: string
    /** 允许的配图方式列表 */
    enabledImageMethods?: string[]
    errorMessage?: string
    fullContent?: string
    id?: number
    images?: ImageResult[]
    mainTitle?: string
    outline?: OutlineSection[]
    /** 当前阶段 */
    phase?: string
    status?: string
    /** 文章风格 */
    style?: string
    subTitle?: string
    taskId?: string
    /** 标题方案列表 */
    titleOptions?: TitleOption[]
    topic?: string
    /** 用户补充描述 */
    userDescription?: string
    userId?: number
  }

  type ArticlePage = {
    pageNumber?: number
    pageSize?: number
    records?: ArticleInfo[]
    totalPage?: number
    totalRow?: number
  }

  type BaseResponse = {
    code?: number
    data?: any
    message?: string
  }

  type ConfirmOutlineRequest = {
    outline: OutlineSection[]
    taskId: string
  }

  type ConfirmTitleRequest = {
    selectedMainTitle: string
    selectedSubTitle: string
    taskId: string
    /** 用户补充描述（可选） */
    userDescription?: string
  }

  type CreateArticleRequest = {
    /** 允许的配图方式，为空表示支持所有 */
    enabledImageMethods?: string[]
    /** 文章风格，允许为空 */
    style?: string
    topic: string
  }

  type DeleteRequest = {
    id: number
  }

  type getArticleProgressTaskIdParams = {
    /** 文章 ID */
    taskId: string
  }

  type getArticleTaskIdParams = {
    /** 文章 ID */
    taskId: string
  }

  type getUserGetParams = {
    /** 用户 ID */
    id: number
  }

  type getUserGetVoParams = {
    /** 用户 ID */
    id: number
  }

  type ImageResult = {
    description?: string
    keywords?: string
    method?: string
    placeholderID?: string
    position?: number
    sectionTitle?: string
    url?: string
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

  type OutlineSection = {
    points?: string[]
    section?: number
    title?: string
  }

  type PageResult = {
    pageNum?: number
    pageSize?: number
    records?: any
    total?: number
  }

  type QueryArticleRequest = {
    pageNum?: number
    pageSize?: number
    status?: string
    userId?: number
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

  type TitleOption = {
    mainTitle?: string
    subTitle?: string
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
    quota?: number
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
