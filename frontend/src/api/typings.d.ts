declare namespace API {
  type AddUserRequest = {
    userAccount: string
    userAvatar?: string
    userName?: string
    userProfile?: string
    userRole?: string
  }

  type AgentExecutionStats = {
    agentCount?: number
    /** key: agentName, value: durationMs */
    agentDurations?: Record<string, any>
    logs?: AgentLog[]
    /** SUCCESS/FAILED/RUNNING */
    overallStatus?: string
    taskId?: string
    totalDurationMs?: number
  }

  type AgentLog = {
    agentName?: string
    createTime?: string
    durationMs?: number
    endTime?: string
    errorMessage?: string
    id?: number
    inputData?: string
    isDelete?: number
    outputData?: string
    prompt?: string
    startTime?: string
    /** RUNNING/SUCCESS/FAILED */
    status?: string
    taskId?: string
    updateTime?: string
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

  type getArticleExecutionLogsTaskIdParams = {
    /** 任务ID */
    taskId: string
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

  type StatisticsVO = {
    /** 活跃用户数（本周） */
    activeUserCount?: number
    /** 平均耗时（毫秒） */
    avgDurationMs?: number
    /** 本月创作数量 */
    monthCount?: number
    /** 配额总使用量 */
    quotaUsed?: number
    /** 成功率（百分比） */
    successRate?: number
    /** 今日创作数量 */
    todayCount?: number
    /** 总创作数量 */
    totalCount?: number
    /** 总用户数 */
    totalUserCount?: number
    /** VIP 用户数 */
    vipUserCount?: number
    /** 本周创作数量 */
    weekCount?: number
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
