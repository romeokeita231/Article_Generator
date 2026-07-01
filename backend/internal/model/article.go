package model

import (
	"time"
)

// Article 文章实体
type Article struct {
	ID                  int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID              string     `gorm:"column:taskId;uniqueIndex:uk_taskId" json:"taskId"`
	UserID              int64      `gorm:"column:userId;index:idx_userId" json:"userId"`
	Topic               string     `gorm:"column:topic" json:"topic"`
	MainTitle           *string    `gorm:"column:mainTitle" json:"mainTitle"`
	SubTitle            *string    `gorm:"column:subTitle" json:"subTitle"`
	Outline             *string    `gorm:"column:outline;type:json" json:"outline"`
	Content             *string    `gorm:"column:content;type:text" json:"content"`
	FullContent         *string    `gorm:"column:fullContent;type:text" json:"fullContent"`
	Images              *string    `gorm:"column:images;type:json" json:"images"`
	Status              string     `gorm:"column:status;default:PENDING;index:idx_status" json:"status"`
	ErrorMessage        *string    `gorm:"column:errorMessage;type:text" json:"errorMessage"`
	CreateTime          time.Time  `gorm:"column:createTime;autoCreateTime;index:idx_createTime" json:"createTime"`
	CompletedTime       *time.Time `gorm:"column:completedTime" json:"completedTime"`
	UpdateTime          time.Time  `gorm:"column:updateTime;autoUpdateTime" json:"updateTime"`
	IsDelete            int        `gorm:"column:isDelete;default:0" json:"-"`
	Style               string     `gorm:"column:style" json:"style"`
	EnabledImageMethods *string    `gorm:"column:enabledImageMethods;type:json" json:"enabledImageMethods"`
	UserDescription     *string    `gorm:"column:userDescription;type:text" json:"userDescription"`
	TitleOptions        *string    `gorm:"column:titleOptions;type:json" json:"titleOptions"`
	Phase               string     `gorm:"column:phase;default:PENDING" json:"phase"`
}

func (Article) TableName() string {
	return "article"
}

// ArticleStatus 文章状态
const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
)

// ArticlePhase 文章阶段
const (
	PhasePending           = "PENDING"            // 等待处理
	PhaseTitleGenerating   = "TITLE_GENERATING"   // 生成标题中
	PhaseTitleSelecting    = "TITLE_SELECTING"    // 等待选择标题
	PhaseOutlineGenerating = "OUTLINE_GENERATING" // 生成大纲中
	PhaseOutlineEditing    = "OUTLINE_EDITING"    // 等待编辑大纲
	PhaseContentGenerating = "CONTENT_GENERATING" // 生成正文中
)

// ArticleInfo 文章信息（响应）
type ArticleInfo struct {
	ID                  int64            `json:"id"`
	TaskID              string           `json:"taskId"`
	UserID              int64            `json:"userId"`
	Topic               string           `json:"topic"`
	MainTitle           *string          `json:"mainTitle"`
	SubTitle            *string          `json:"subTitle"`
	Outline             []OutlineSection `json:"outline"`
	Content             *string          `json:"content"`
	FullContent         *string          `json:"fullContent"`
	Images              []ImageResult    `json:"images"`
	Status              string           `json:"status"`
	ErrorMessage        *string          `json:"errorMessage"`
	CreateTime          time.Time        `json:"createTime"`
	CompletedTime       *time.Time       `json:"completedTime"`
	Style               string           `json:"style"`               // 文章风格
	EnabledImageMethods []string         `json:"enabledImageMethods"` // 允许的配图方式列表
	UserDescription     *string          `json:"userDescription"`     // 用户补充描述
	TitleOptions        []TitleOption    `json:"titleOptions"`        // 标题方案列表
	Phase               string           `json:"phase"`               // 当前阶段
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

	if a.TitleOptions != nil {
		parseJSON(*a.TitleOptions, &info.TitleOptions)
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
	TaskID                  string             `json:"taskId"`
	Topic                   string             `json:"topic"`
	Title                   *TitleResult       `json:"title"`
	Outline                 *OutlineResult     `json:"outline"`
	Content                 string             `json:"content"`
	FullContent             string             `json:"fullContent"`
	ImageRequirements       []ImageRequirement `json:"imageRequirements"`
	Images                  []ImageResult      `json:"images"`
	ContentWithPlaceholders string             `json:"contentWithPlaceholders"`
	UserDescription         string             `json:"userDescription"` // 用户补充描述
	Phase                   string             `json:"phase"`           // 当前阶段
	TitleOptions            []TitleOption      `json:"titleOptions"`    // 标题方案列表
	Style                   string             `json:"style"`               // 文章风格
	EnabledImageMethods     []string           `json:"enabledImageMethods"` // 允许的配图方式列表
}

// TitleOption 标题方案
type TitleOption struct {
	MainTitle string `json:"mainTitle"`
	SubTitle  string `json:"subTitle"`
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
	Position      int    `json:"position"`
	Type          string `json:"type"`
	SectionTitle  string `json:"sectionTitle"`
	ImageSource   string `json:"imageSource"`
	Keywords      string `json:"keywords"`
	Prompt        string `json:"prompt"`
	PlaceholderID string `json:"placeholderId"`
}

// ImageResult 配图结果（智能体5输出）
type ImageResult struct {
	Position      int    `json:"position"`
	URL           string `json:"url"`
	Method        string `json:"method"`
	Keywords      string `json:"keywords"`
	SectionTitle  string `json:"sectionTitle"`
	Description   string `json:"description"`
	PlaceholderID string `json:"placeholderID"`
}

// ArticlePage 文章分页结果
type ArticlePage struct {
	PageNumber int64          `json:"pageNumber"`
	PageSize   int64          `json:"pageSize"`
	TotalRow   int64          `json:"totalRow"`
	TotalPage  int64          `json:"totalPage"`
	Records    []*ArticleInfo `json:"records"`
}
