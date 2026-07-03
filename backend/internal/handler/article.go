package handler

import (
	"net/http"
	"io"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

    "github.com/romeokeita231/Article_Generator/internal/service"
    "github.com/romeokeita231/Article_Generator/internal/common"
    "github.com/romeokeita231/Article_Generator/internal/model"
	
)

type ArticleHandler struct {
    svc        *service.ArticleService
    userSvc    *service.UserService
    sseManager *common.SSEManager
    agentLogService *service.AgentLogService
}

func NewArticleHandler(svc *service.ArticleService, userSvc *service.UserService, agentLogService *service.AgentLogService, sseManager *common.SSEManager) *ArticleHandler {
    return &ArticleHandler{
        svc: svc, userSvc: userSvc, sseManager: sseManager,
    }
}


// Create
// @Summary 创建文章
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param request body model.CreateArticleRequest true "创建文章请求体"
// @Success 200 {object} common.BaseResponse{data=int64}
// @Router /article/create [post]
func (h *ArticleHandler) Create(c *gin.Context) {
    var req model.CreateArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrParams))
        return
    }

    // 校验风格参数（允许为空）
    if !common.IsValidArticleStyle(req.Style) {
        c.JSON(http.StatusOK, common.Error(common.ErrParams.WithMessage("无效的文章风格")))
        return
    }

    session := sessions.Default(c)
    user, err := h.userSvc.GetLoginUser(session)
    if err != nil {
        handleError(c, err)
        return
    }

    taskID, err := h.svc.Create(user, &req)
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, common.Success(taskID))
}

// GetProgress
// @Summary 获取文章进度
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param taskId path string true "文章 ID"
// @Success 200 {string} string "文章进度"
// @Router /article/progress/{taskId} [get]
func (h *ArticleHandler) GetProgress(c *gin.Context) {
    taskID := c.Param("taskId")
    // ... 权限校验 ...

    // 设置 SSE 响应头
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("X-Accel-Buffering", "no")

    messageChan := h.sseManager.Register(taskID)
    defer h.sseManager.Unregister(taskID)

    c.Stream(func(w io.Writer) bool {
        select {
        case msg, ok := <-messageChan:
            if !ok {
                return false
            }
            c.SSEvent("message", msg)
            c.Writer.Flush()
            return true
        case <-c.Request.Context().Done():
            return false
        case <-time.After(30 * time.Minute):
            return false
        }
    })
}

// Get
// @Summary 获取文章
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param taskId path string true "文章 ID"
// @Success 200 {object} common.BaseResponse{data=model.ArticleInfo}
// @Router /article/{taskId} [get]
func (h *ArticleHandler) Get(c *gin.Context) {
    taskID := c.Param("taskId")
    session := sessions.Default(c)
    user, err := h.userSvc.GetLoginUser(session)
    if err != nil {
        handleError(c, err)
        return
    }
    isAdmin := user.UserRole == common.AdminRole
    article, err := h.svc.GetByTaskID(taskID, user.ID, isAdmin)
    if err != nil {
        handleError(c, err)
        return
    }
    c.JSON(http.StatusOK, common.Success(article))
}

// List
// @Summary 分页查询文章列表
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param request body model.QueryArticleRequest true "查询文章列表请求体"
// @Success 200 {object} common.BaseResponse{data=model.ArticlePage}
// @Router /article/list [post]
func (h *ArticleHandler) List(c *gin.Context) {
    var req model.QueryArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrParams))
        return
    }
    session := sessions.Default(c)
    user, err := h.userSvc.GetLoginUser(session)
    if err != nil {
        handleError(c, err)
        return
    }
    isAdmin := user.UserRole == common.AdminRole
    page, err := h.svc.ListByPage(&req, user.ID, isAdmin)
    if err != nil {
        handleError(c, err)
        return
    }
    c.JSON(http.StatusOK, common.Success(page))
}

// Delete
// @Summary 删除文章
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param request body model.DeleteRequest true "删除文章请求体"
// @Success 200 {object} common.BaseResponse{data=bool}
// @Router /article/delete [post]
func (h *ArticleHandler) Delete(c *gin.Context) {
    var req model.DeleteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrParams))
        return
    }
    session := sessions.Default(c)
    user, err := h.userSvc.GetLoginUser(session)
    if err != nil {
        handleError(c, err)
        return
    }
    isAdmin := user.UserRole == common.AdminRole
    if err := h.svc.Delete(req.ID, user.ID, isAdmin); err != nil {
        handleError(c, err)
        return
    }
    c.JSON(http.StatusOK, common.Success(true))
}

// ConfirmTitle 
// @Summary 确认标题并输入补充描述
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param request body model.ConfirmTitleRequest true "确认标题并输入补充描述请求体"
// @Success 200 {object} common.BaseResponse
// @Router /article/confirmTitle [post]
func (h *ArticleHandler) ConfirmTitle(c *gin.Context) {
    var req model.ConfirmTitleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrParams))
        return
    }

    session := sessions.Default(c)
    user, err := h.userSvc.GetLoginUser(session)
    if err != nil {
        handleError(c, err)
        return
    }

    isAdmin := user.UserRole == common.AdminRole
    if err := h.svc.ConfirmTitle(req.TaskID, req.SelectedMainTitle, req.SelectedSubTitle,
        req.UserDescription, user.ID, isAdmin); err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, common.Success(nil))
}

// ConfirmOutline 
// @Summary 确认大纲
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param request body model.ConfirmOutlineRequest true "确认大纲请求体"
// @Success 200 {object} common.BaseResponse
// @Router /article/confirmOutline [post]
func (h *ArticleHandler) ConfirmOutline(c *gin.Context) {
    var req model.ConfirmOutlineRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrParams))
        return
    }

    session := sessions.Default(c)
    user, err := h.userSvc.GetLoginUser(session)
    if err != nil {
        handleError(c, err)
        return
    }

    isAdmin := user.UserRole == common.AdminRole
    if err := h.svc.ConfirmOutline(req.TaskID, req.Outline, user.ID, isAdmin); err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, common.Success(nil))
}

// AiModifyOutline 
// @Summary 使用 AI 修改大纲
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param request body model.AiModifyOutlineRequest true "AI 修改大纲请求体"
// @Success 200 {object} common.BaseResponse{data=model.OutlineSection}
// @Router /article/aiModifyOutline [post]
func (h *ArticleHandler) AiModifyOutline(c *gin.Context) {
    var req model.AiModifyOutlineRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrParams))
        return
    }

    session := sessions.Default(c)
    user, err := h.userSvc.GetLoginUser(session)
    if err != nil {
        handleError(c, err)
        return
    }

    isAdmin := user.UserRole == common.AdminRole
    modifiedOutline, err := h.svc.AiModifyOutline(req.TaskID, req.ModifySuggestion, user.ID, isAdmin)
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, common.Success(modifiedOutline))
}

// GetExecutionLogs 
// @Summary 获取任务执行日志
// @Tags articleHandler
// @Accept json
// @Produce json
// @Param taskId path string true "任务ID"
// @Success 200 {object} common.BaseResponse{data=model.AgentExecutionStats}
// @Router /article/execution-logs/{taskId} [get]
func (h *ArticleHandler) GetExecutionLogs(c *gin.Context) {
    taskID := c.Param("taskId")
    if taskID == "" {
        c.JSON(http.StatusOK, common.Error(common.ErrParams.WithMessage("任务ID不能为空")))
        return
    }

    stats, err := h.agentLogService.GetExecutionStats(taskID)
    if err != nil {
        c.JSON(http.StatusOK, common.Error(common.ErrSystem.WithMessage("获取执行日志失败: "+err.Error())))
        return
    }

    c.JSON(http.StatusOK, common.Success(stats))
}
