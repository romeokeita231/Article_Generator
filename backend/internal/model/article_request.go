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
