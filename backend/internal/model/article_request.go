package model

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
    Topic string `json:"topic" binding:"required"`
    Style               string   `json:"style"`               // 文章风格，允许为空
    EnabledImageMethods []string `json:"enabledImageMethods"` // 允许的配图方式，为空表示支持所有

}

// QueryArticleRequest 查询文章请求
type QueryArticleRequest struct {
    UserID   *int64  `json:"userId"`
    Status   *string `json:"status"`
    PageNum  int64   `json:"pageNum"`
    PageSize int64   `json:"pageSize"`
}

// ConfirmTitleRequest 确认标题请求
type ConfirmTitleRequest struct {
    TaskID            string  `json:"taskId" binding:"required"`
    SelectedMainTitle string  `json:"selectedMainTitle" binding:"required"`
    SelectedSubTitle  string  `json:"selectedSubTitle" binding:"required"`
    UserDescription   *string `json:"userDescription"` // 用户补充描述（可选）
}

// ConfirmOutlineRequest 确认大纲请求
type ConfirmOutlineRequest struct {
    TaskID  string           `json:"taskId" binding:"required"`
    Outline []OutlineSection `json:"outline" binding:"required"`
}

// AiModifyOutlineRequest AI 修改大纲请求
type AiModifyOutlineRequest struct {
    TaskID           string `json:"taskId" binding:"required"`
    ModifySuggestion string `json:"modifySuggestion" binding:"required"`
}
