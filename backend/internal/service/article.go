package service

import (
	"log"
	"time"
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"github.com/google/uuid"

	"github.com/romeokeita231/Article_Generator/internal/model"
	"github.com/romeokeita231/Article_Generator/internal/store"
	"github.com/romeokeita231/Article_Generator/internal/common"
	

)

type ArticleService struct {
	store        *store.ArticleStore
	agentSvc     *ArticleAgentService
	quotaSvc     *QuotaService
	sseManager   *common.SSEManager
}



// NewArticleService 创建文章服务
func NewArticleService(st *store.ArticleStore, agentSvc *ArticleAgentService, quotaSvc *QuotaService, sseManager *common.SSEManager) *ArticleService {
	return &ArticleService{
		store:      st,
		agentSvc:   agentSvc,
		quotaSvc:   quotaSvc,
		sseManager: sseManager,
	}
}


// Create 创建文章任务
func (s *ArticleService) Create(user *model.User, req *model.CreateArticleRequest) (string, error) {
    if req.Topic == "" {
        return "", common.ErrParams.WithMessage("选题不能为空")
    }

    // 检查并消耗配额（原子操作）
    if err := s.quotaSvc.CheckAndConsumeQuota(user); err != nil {
        return "", err
    }

    // 生成任务ID
    taskID := uuid.NewString()

    // 将 enabledImageMethods 转为 JSON（为空时设置为 nil）
	var methodsJSON *string
	if len(req.EnabledImageMethods) > 0 {
		methodsBytes, _ := json.Marshal(req.EnabledImageMethods)
		methodsStr := string(methodsBytes)
		methodsJSON = &methodsStr
	}

    // 创建文章记录
	article := &model.Article{
		TaskID:              taskID,
		UserID:              user.ID,
		Topic:               req.Topic,
		Style:               req.Style,
		EnabledImageMethods: methodsJSON,
		Status:              model.StatusPending,
		Phase:               model.PhasePending,
		CreateTime:          time.Now(),
	}

    if err := s.store.Create(article); err != nil {
        return "", common.ErrOperation
    }

    // 异步执行阶段1：生成标题方案
	go s.ExecutePhase1Async(taskID, req.Topic, req.Style)

	log.Printf("文章任务已创建, taskId=%s, userId=%d, style=%s", taskID, user.ID, req.Style)
    return taskID, nil
}

// ExecutePhase1Async 阶段1：异步生成标题方案
func (s *ArticleService) ExecutePhase1Async(taskID, topic, style string) {
    _ = s.store.UpdateStatus(taskID, model.StatusProcessing, nil)
    _ = s.UpdatePhase(taskID, model.PhaseTitleGenerating)

    state := &model.ArticleState{TaskID: taskID, Topic: topic, Style: style}

    ctx := context.Background()
    err := s.agentSvc.ExecutePhase1(ctx, state)
    if err != nil {
        errMsg := err.Error()
        _ = s.store.UpdateStatus(taskID, model.StatusFailed, &errMsg)
        s.sseManager.Send(taskID, map[string]interface{}{
            "type": common.SSEMsgError, "message": errMsg,
        })
        s.sseManager.Complete(taskID)
        return
    }

    // 保存标题方案，更新阶段为等待选择
    _ = s.SaveTitleOptions(taskID, state.TitleOptions)
    _ = s.UpdatePhase(taskID, model.PhaseTitleSelecting)

    // 推送标题方案生成完成消息（不关闭 SSE）
    s.sseManager.Send(taskID, map[string]interface{}{
        "type":         common.SSEMsgTitlesGenerated,
        "titleOptions": state.TitleOptions,
    })
}

// ExecutePhase2Async 阶段2：异步生成大纲（用户确认标题后调用）
func (s *ArticleService) ExecutePhase2Async(taskID string) {
    article, err := s.store.GetByTaskID(taskID)
    if err != nil { return }

    _ = s.UpdatePhase(taskID, model.PhaseOutlineGenerating)

    state := &model.ArticleState{
        TaskID:          taskID,
        Style:           article.Style,
        UserDescription: "",
    }
    if article.UserDescription != nil {
        state.UserDescription = *article.UserDescription
    }
    state.Title = &model.TitleResult{
        MainTitle: *article.MainTitle,
        SubTitle:  *article.SubTitle,
    }

    ctx := context.Background()
    err = s.agentSvc.ExecutePhase2(ctx, state)
    if err != nil {
        // ... 错误处理 ...
        return
    }

    // 保存大纲，更新阶段为等待编辑
    outlineJSON, _ := json.Marshal(state.Outline.Sections)
    outlineStr := string(outlineJSON)
    article.Outline = &outlineStr
    _ = s.store.Update(article)
    _ = s.UpdatePhase(taskID, model.PhaseOutlineEditing)

    // 推送大纲生成完成消息（不关闭 SSE）
    s.sseManager.Send(taskID, map[string]interface{}{
        "type":    common.SSEMsgOutlineGenerated,
        "outline": state.Outline.Sections,
    })
}

// ExecutePhase3Async 阶段3：异步生成正文+配图（用户确认大纲后调用）
func (s *ArticleService) ExecutePhase3Async(taskID string) {
    article, err := s.store.GetByTaskID(taskID)
    if err != nil { return }

    _ = s.UpdatePhase(taskID, model.PhaseContentGenerating)

    state := &model.ArticleState{TaskID: taskID, Style: article.Style}

    // 从数据库恢复配图方式、标题、大纲
    if article.EnabledImageMethods != nil && *article.EnabledImageMethods != "" {
        _ = json.Unmarshal([]byte(*article.EnabledImageMethods), &state.EnabledImageMethods)
    }
    state.Title = &model.TitleResult{
        MainTitle: *article.MainTitle, SubTitle: *article.SubTitle,
    }
    var outlineSections []model.OutlineSection
    if article.Outline != nil {
        _ = json.Unmarshal([]byte(*article.Outline), &outlineSections)
    }
    state.Outline = &model.OutlineResult{Sections: outlineSections}

    ctx := context.Background()
    err = s.agentSvc.ExecutePhase3(ctx, state)
    if err != nil {
        // ... 错误处理 ...
        return
    }

    _ = s.saveArticle(taskID, state)
    _ = s.store.UpdateStatus(taskID, model.StatusCompleted, nil)

    s.sseManager.Send(taskID, map[string]interface{}{
        "type": common.SSEMsgAllComplete, "taskId": taskID,
    })
    s.sseManager.Complete(taskID) // 阶段3结束才关闭 SSE
}

// GetByTaskID 根据任务ID获取文章详情
func (s *ArticleService) GetByTaskID(taskID string, userID int64, isAdmin bool) (*model.ArticleInfo, error) {
    article, err := s.store.GetByTaskID(taskID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, common.ErrNotFound.WithMessage("文章不存在")
        }
        return nil, common.ErrSystem
    }

    // 权限校验：只能查看自己的文章（管理员除外）
    if !isAdmin && article.UserID != userID {
        return nil, common.ErrNoAuth
    }

    return article.ToArticleInfo(), nil
}

// ListByPage 分页查询文章列表
func (s *ArticleService) ListByPage(req *model.QueryArticleRequest,
    userID int64, isAdmin bool) (*model.PageResult, error) {
    if req.PageNum <= 0 { req.PageNum = common.DefaultPageNum }
    if req.PageSize <= 0 { req.PageSize = common.DefaultPageSize }
    if req.PageSize > common.MaxPageSize { req.PageSize = common.MaxPageSize }

    queryUserID := &userID
    if isAdmin && req.UserID != nil {
        queryUserID = req.UserID
    }

    articles, total, err := s.store.List(queryUserID, req.Status, isAdmin, req.PageNum, req.PageSize)
    if err != nil {
        return nil, common.ErrSystem
    }

    articleInfos := make([]model.ArticleInfo, 0, len(articles))
    for i := range articles {
        if info := articles[i].ToArticleInfo(); info != nil {
            articleInfos = append(articleInfos, *info)
        }
    }

    return &model.PageResult{
        Total: total, Records: articleInfos,
        PageNum: req.PageNum, PageSize: req.PageSize,
    }, nil
}

// Delete 删除文章
func (s *ArticleService) Delete(id int64, userID int64, isAdmin bool) error {
    article, err := s.store.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return common.ErrNotFound
        }
        return common.ErrSystem
    }

    if !isAdmin && article.UserID != userID {
        return common.ErrNoAuth
    }

    if err := s.store.Delete(id); err != nil {
        return common.ErrOperation
    }
    return nil
}

// saveArticle 保存文章到数据库
func (s *ArticleService) saveArticle(taskID string, state *model.ArticleState) error {
    article, err := s.store.GetByTaskID(taskID)
    if err != nil {
        return err
    }

    outlineJSON, _ := json.Marshal(state.Outline.Sections)
    imagesJSON, _ := json.Marshal(state.Images)
    outlineStr := string(outlineJSON)
    imagesStr := string(imagesJSON)
    now := time.Now()

    article.MainTitle = &state.Title.MainTitle
    article.SubTitle = &state.Title.SubTitle
    article.Outline = &outlineStr
    article.Content = &state.Content
    article.FullContent = &state.FullContent
    article.Images = &imagesStr
    article.CompletedTime = &now

    return s.store.Update(article)
}

// ConfirmTitle 确认标题并输入补充描述
func (s *ArticleService) ConfirmTitle(taskID, mainTitle, subTitle string, userDescription *string, userID int64, isAdmin bool) error {
	// 获取文章
	article, err := s.store.GetByTaskID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNotFound.WithMessage("文章不存在")
		}
		return common.ErrSystem
	}

	// 权限校验
	if !isAdmin && article.UserID != userID {
		return common.ErrNoAuth
	}

	// 校验当前阶段
	if article.Phase != model.PhaseTitleSelecting {
		return common.ErrParams.WithMessage("当前阶段不允许确认标题")
	}

	// 更新标题和用户补充描述
	article.MainTitle = &mainTitle
	article.SubTitle = &subTitle
	article.UserDescription = userDescription

	if err := s.store.Update(article); err != nil {
		return common.ErrOperation
	}

	// 异步执行阶段2：生成大纲
	go s.ExecutePhase2Async(taskID)

	return nil
}

// ConfirmOutline 确认大纲
func (s *ArticleService) ConfirmOutline(taskID string, outline []model.OutlineSection, userID int64, isAdmin bool) error {
	// 获取文章
	article, err := s.store.GetByTaskID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrNotFound.WithMessage("文章不存在")
		}
		return common.ErrSystem
	}

	// 权限校验
	if !isAdmin && article.UserID != userID {
		return common.ErrNoAuth
	}

	// 校验当前阶段
	if article.Phase != model.PhaseOutlineEditing {
		return common.ErrParams.WithMessage("当前阶段不允许确认大纲")
	}

	// 更新大纲
	outlineJSON, _ := json.Marshal(outline)
	outlineStr := string(outlineJSON)
	article.Outline = &outlineStr

	if err := s.store.Update(article); err != nil {
		return common.ErrOperation
	}

	// 异步执行阶段3：生成正文+配图
	go s.ExecutePhase3Async(taskID)

	return nil
}

// UpdatePhase 更新阶段
func (s *ArticleService) UpdatePhase(taskID, phase string) error {
	return s.store.UpdatePhase(taskID, phase)
}

// SaveTitleOptions 保存标题方案
func (s *ArticleService) SaveTitleOptions(taskID string, titleOptions []model.TitleOption) error {
	optionsJSON, _ := json.Marshal(titleOptions)
	optionsStr := string(optionsJSON)
	return s.store.UpdateTitleOptions(taskID, optionsStr)
}

// AiModifyOutline AI 修改大纲
func (s *ArticleService) AiModifyOutline(taskID, modifySuggestion string, userID int64, isAdmin bool) ([]model.OutlineSection, error) {
    article, err := s.store.GetByTaskID(taskID)
    if err != nil {
        return nil, common.ErrNotFound.WithMessage("文章不存在")
    }

    if !isAdmin && article.UserID != userID {
        return nil, common.ErrNoAuth
    }

    if article.Phase != model.PhaseOutlineEditing {
        return nil, common.ErrParams.WithMessage("当前阶段不允许修改大纲")
    }

    var currentOutline []model.OutlineSection
    if article.Outline != nil {
        _ = json.Unmarshal([]byte(*article.Outline), &currentOutline)
    }

    ctx := context.Background()
    return s.agentSvc.AiModifyOutline(ctx, *article.MainTitle, *article.SubTitle, currentOutline, modifySuggestion)
}
