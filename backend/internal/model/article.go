package model

import (
    "time"
)

// Article 文章实体
type Article struct {
    ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
    TaskID        string     `gorm:"column:taskId;uniqueIndex:uk_taskId" json:"taskId"`
    UserID        int64      `gorm:"column:userId;index:idx_userId" json:"userId"`
    Topic         string     `gorm:"column:topic" json:"topic"`
    MainTitle     *string    `gorm:"column:mainTitle" json:"mainTitle"`
    SubTitle      *string    `gorm:"column:subTitle" json:"subTitle"`
    Outline       *string    `gorm:"column:outline;type:json" json:"outline"`
    Content       *string    `gorm:"column:content;type:text" json:"content"`
    FullContent   *string    `gorm:"column:fullContent;type:text" json:"fullContent"`
    Images        *string    `gorm:"column:images;type:json" json:"images"`
    Status        string     `gorm:"column:status;default:PENDING;index:idx_status" json:"status"`
    ErrorMessage  *string    `gorm:"column:errorMessage;type:text" json:"errorMessage"`
    CreateTime    time.Time  `gorm:"column:createTime;autoCreateTime;index:idx_createTime" json:"createTime"`
    CompletedTime *time.Time `gorm:"column:completedTime" json:"completedTime"`
    UpdateTime    time.Time  `gorm:"column:updateTime;autoUpdateTime" json:"updateTime"`
    IsDelete      int        `gorm:"column:isDelete;default:0" json:"-"`
}

func (Article) TableName() string {
    return "article"
}

const (
    StatusPending    = "PENDING"
    StatusProcessing = "PROCESSING"
    StatusCompleted  = "COMPLETED"
    StatusFailed     = "FAILED"
)

// ArticleInfo 文章信息（响应）
type ArticleInfo struct {
    ID            int64            `json:"id"`
    TaskID        string           `json:"taskId"`
    UserID        int64            `json:"userId"`
    Topic         string           `json:"topic"`
    MainTitle     *string          `json:"mainTitle"`
    SubTitle      *string          `json:"subTitle"`
    Outline       []OutlineSection `json:"outline"`
    Content       *string          `json:"content"`
    FullContent   *string          `json:"fullContent"`
    Images        []ImageResult    `json:"images"`
    Status        string           `json:"status"`
    ErrorMessage  *string          `json:"errorMessage"`
    CreateTime    time.Time        `json:"createTime"`
    CompletedTime *time.Time       `json:"completedTime"`
}

func (a *Article) ToArticleInfo() *ArticleInfo {
    if a == nil {
        return nil
    }

    info := &ArticleInfo{
        ID: a.ID, TaskID: a.TaskID, UserID: a.UserID, Topic: a.Topic,
        MainTitle: a.MainTitle, SubTitle: a.SubTitle,
        Content: a.Content, FullContent: a.FullContent,
        Status: a.Status, ErrorMessage: a.ErrorMessage,
        CreateTime: a.CreateTime, CompletedTime: a.CompletedTime,
    }

    if a.Outline != nil {
        parseJSON(*a.Outline, &info.Outline)
    }
    if a.Images != nil {
        parseJSON(*a.Images, &info.Images)
    }

    return info
}

// ArticleState 文章生成状态（智能体间共享）
type ArticleState struct {
    TaskID            string             `json:"taskId"`
    Topic             string             `json:"topic"`
    Title             *TitleResult       `json:"title"`
    Outline           *OutlineResult     `json:"outline"`
    Content           string             `json:"content"`
    FullContent       string             `json:"fullContent"`
    ImageRequirements []ImageRequirement `json:"imageRequirements"`
    Images            []ImageResult      `json:"images"`
}

// TitleResult 标题结果（智能体1输出）
type TitleResult struct {
    MainTitle string `json:"mainTitle"`
    SubTitle  string `json:"subTitle"`
}

// OutlineResult 大纲结果（智能体2输出）
type OutlineResult struct {
    Sections []OutlineSection `json:"sections"`
}

// OutlineSection 大纲章节
type OutlineSection struct {
    Section int      `json:"section"`
    Title   string   `json:"title"`
    Points  []string `json:"points"`
}

// ImageRequirement 配图需求（智能体4输出）
type ImageRequirement struct {
    Position     int    `json:"position"`
    Type         string `json:"type"`
    SectionTitle string `json:"sectionTitle"`
    Keywords     string `json:"keywords"`
}

// ImageResult 配图结果（智能体5输出）
type ImageResult struct {
    Position     int    `json:"position"`
    URL          string `json:"url"`
    Method       string `json:"method"`
    Keywords     string `json:"keywords"`
    SectionTitle string `json:"sectionTitle"`
    Description  string `json:"description"`
}

// ArticlePage 文章分页结果
type ArticlePage struct {
    PageNumber int64          `json:"pageNumber"`
    PageSize   int64          `json:"pageSize"`
    TotalRow   int64          `json:"totalRow"`
    TotalPage  int64          `json:"totalPage"`
    Records    []*ArticleInfo `json:"records"`
}