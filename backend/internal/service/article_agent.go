package service

import (
	"context"
	"fmt"
    "log"
    "strings"
    "encoding/json"

	"github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/llms/openai"
	"github.com/romeokeita231/Article_Generator/internal/common"
	"github.com/romeokeita231/Article_Generator/internal/config"
	"github.com/romeokeita231/Article_Generator/internal/model"

)

// ArticleAgentService 文章智能体编排服务
type ArticleAgentService struct {
    llm        llms.Model
    imageStrategy *ImageServiceStrategy // 替换原来的 pexels + cos
    sseManager *common.SSEManager
}

func NewArticleAgentService(cfg *config.Config, imageStrategy *ImageServiceStrategy, sseManager *common.SSEManager) (*ArticleAgentService, error) {

    llm, err := openai.New(
        openai.WithToken(cfg.AI.DashScope.APIKey),
        openai.WithModel(cfg.AI.DashScope.Model),
        openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
    )
    if err != nil {
        return nil, fmt.Errorf("create dashscope client: %w", err)
    }

    return &ArticleAgentService{
        llm: llm, 
        imageStrategy: imageStrategy, 
        sseManager: sseManager,
    }, nil
}

func (s *ArticleAgentService) Execute(ctx context.Context, state *model.ArticleState) error {
    // 智能体1：生成标题
    if err := s.agent1GenerateTitle(ctx, state); err != nil {
        return fmt.Errorf("agent1 failed: %w", err)
    }
    s.sendMessage(state.TaskID, map[string]interface{}{
        "type": "AGENT1_COMPLETE", "title": state.Title,
    })

    // 智能体2：生成大纲（流式）
    if err := s.agent2GenerateOutlineStream(ctx, state); err != nil {
        return fmt.Errorf("agent2 failed: %w", err)
    }
    s.sendMessage(state.TaskID, map[string]interface{}{
        "type": "AGENT2_COMPLETE", "outline": state.Outline.Sections,
    })

    // 智能体3：生成正文（流式）
    if err := s.agent3GenerateContent(ctx, state); err != nil {
        return fmt.Errorf("agent3 failed: %w", err)
    }
    s.sendMessage(state.TaskID, map[string]interface{}{"type": "AGENT3_COMPLETE"})

    // 智能体4：分析配图需求
    if err := s.agent4AnalyzeImageRequirements(ctx, state); err != nil {
        return fmt.Errorf("agent4 failed: %w", err)
    }
    s.sendMessage(state.TaskID, map[string]interface{}{
        "type": "AGENT4_COMPLETE", "imageRequirements": state.ImageRequirements,
    })

    // 智能体5：生成配图
    if err := s.agent5GenerateImages(ctx, state); err != nil {
        return fmt.Errorf("agent5 failed: %w", err)
    }
    s.sendMessage(state.TaskID, map[string]interface{}{
        "type": "AGENT5_COMPLETE", "images": state.Images,
    })

    // 图文合成
    s.mergeImagesIntoContent(state)
    s.sendMessage(state.TaskID, map[string]interface{}{
        "type": "MERGE_COMPLETE", "fullContent": state.FullContent,
    })

    return nil
}

func (s *ArticleAgentService) agent1GenerateTitle(ctx context.Context, state *model.ArticleState) error {
    prompt := strings.ReplaceAll(common.Agent1TitlePrompt, "{topic}", state.Topic)

    content, err := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt)
    if err != nil {
        return fmt.Errorf("LLM call failed: %w", err)
    }

    var title model.TitleResult
    if err := json.Unmarshal([]byte(content), &title); err != nil {
        return fmt.Errorf("parse title failed: %w", err)
    }

    state.Title = &title
    log.Printf("智能体1：标题生成成功, mainTitle=%s", title.MainTitle)
    return nil
}



func (s *ArticleAgentService) agent3GenerateContent(ctx context.Context, state *model.ArticleState) error {
    outlineJSON, _ := json.Marshal(state.Outline.Sections)
    prompt := strings.ReplaceAll(common.Agent3ContentPrompt, "{mainTitle}", state.Title.MainTitle)
    prompt = strings.ReplaceAll(prompt, "{subTitle}", state.Title.SubTitle)
    prompt = strings.ReplaceAll(prompt, "{outline}", string(outlineJSON))

    var contentBuilder strings.Builder

    _, err := s.llm.GenerateContent(ctx, []llms.MessageContent{
        llms.TextParts(llms.ChatMessageTypeHuman, prompt),
    }, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
        text := string(chunk)
        contentBuilder.WriteString(text)
        s.sendMessage(state.TaskID, map[string]interface{}{
            "type": "AGENT3_STREAMING", "content": text,
        })
        return nil
    }))
    if err != nil {
        return err
    }

    state.Content = contentBuilder.String()
    log.Printf("智能体3：正文生成成功, length=%d", len(state.Content))
    return nil
}

func (s *ArticleAgentService) agent4AnalyzeImageRequirements(ctx context.Context, state *model.ArticleState) error {
    prompt := strings.ReplaceAll(common.Agent4ImageRequirementsPrompt, "{mainTitle}", state.Title.MainTitle)
    prompt = strings.ReplaceAll(prompt, "{content}", state.Content)

    content, err := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt)
    if err != nil {
        return err
    }

    var requirements []model.ImageRequirement
    if err := json.Unmarshal([]byte(content), &requirements); err != nil {
        return fmt.Errorf("parse image requirements failed: %w", err)
    }

    state.ImageRequirements = requirements
    log.Printf("智能体4：配图需求分析成功, count=%d", len(requirements))
    return nil
}

func (s *ArticleAgentService) agent5GenerateImages(ctx context.Context, state *model.ArticleState) error {
    var imageResults []model.ImageResult

    for _, req := range state.ImageRequirements {
        imageRequest := &model.ImageRequest{
            Keywords: req.Keywords,
            Prompt:   req.Prompt,
            Position: req.Position,
            Type:     req.Type,
        }

        // 通过策略选择器获取图片并上传到 COS
        result, err := s.imageStrategy.GetImageAndUpload(req.ImageSource, imageRequest)
        if err != nil {
            log.Printf("智能体5：获取图片失败, position=%d, error=%v", req.Position, err)
            continue // 失败时跳过，不中断整个流程
        }

        imageResult := s.buildImageResult(&req, result.URL, result.Method)
        imageResults = append(imageResults, imageResult)

        s.sendMessage(state.TaskID, map[string]interface{}{
            "type": "IMAGE_COMPLETE", "image": imageResult,
        })
    }

    state.Images = imageResults
    return nil
}

// buildImageResult 构建配图结果对象
func (s *ArticleAgentService) buildImageResult(req *model.ImageRequirement, cosURL, method string) model.ImageResult {
	return model.ImageResult{
		Position:      req.Position,
		URL:           cosURL,
		Method:        method,
		Keywords:      req.Keywords,
		SectionTitle:  req.SectionTitle,
		Description:   req.Type,
		PlaceholderID: req.PlaceholderID,
	}
}

func (s *ArticleAgentService) mergeImagesIntoContent(state *model.ArticleState) {
    // 使用包含占位符的正文（Agent4 已在正文中预埋好占位符）
    content := state.ContentWithPlaceholders

    if len(state.Images) == 0 {
        state.FullContent = content
        return
    }

    fullContent := content

    
    // 遍历所有配图，根据占位符替换为实际图片
    for _, image := range state.Images {
        if image.PlaceholderID != "" && strings.Contains(fullContent, image.PlaceholderID) {
            imageMarkdown := fmt.Sprintf("![%s](%s)", image.Description, image.URL)
            fullContent = strings.ReplaceAll(fullContent, image.PlaceholderID, imageMarkdown)
        }
    }

    state.FullContent = fullContent
}

// sendMessage 发送 SSE 消息
func (s *ArticleAgentService) sendMessage(taskID string, data interface{}) {
	s.sseManager.Send(taskID, data)
}

// ExecutePhase1 阶段1：生成标题方案
func (s *ArticleAgentService) ExecutePhase1(ctx context.Context, state *model.ArticleState) error {
    if err := s.agent1GenerateTitleOptions(ctx, state); err != nil {
        return fmt.Errorf("agent1 failed: %w", err)
    }
    s.sendMessage(state.TaskID, map[string]interface{}{
        "type": common.SSEMsgAgent1Complete, "titleOptions": state.TitleOptions,
    })
    return nil
}

// ExecutePhase2 阶段2：生成大纲
func (s *ArticleAgentService) ExecutePhase2(ctx context.Context, state *model.ArticleState) error {
    if err := s.agent2GenerateOutlineStream(ctx, state); err != nil {
        return fmt.Errorf("agent2 failed: %w", err)
    }
    s.sendMessage(state.TaskID, map[string]interface{}{
        "type": common.SSEMsgAgent2Complete, "outline": state.Outline.Sections,
    })
    return nil
}

// ExecutePhase3 阶段3：生成正文+配图
func (s *ArticleAgentService) ExecutePhase3(ctx context.Context, state *model.ArticleState) error {
    // 智能体3-5 + 图文合成（逻辑与原 Execute 方法相同）
    if err := s.agent3GenerateContent(ctx, state); err != nil {
        return fmt.Errorf("agent3 failed: %w", err)
    }
    // ... 后续智能体调用不变 ...
    return nil
}

// agent1GenerateTitleOptions 智能体1：生成标题方案（3-5个）
func (s *ArticleAgentService) agent1GenerateTitleOptions(ctx context.Context, state *model.ArticleState) error {
    prompt := strings.ReplaceAll(common.Agent1TitlePrompt, "{topic}", state.Topic)
    prompt += s.getStylePrompt(state.Style)

    content, err := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt)
    if err != nil {
        return fmt.Errorf("LLM call failed: %w", err)
    }

    // 解析标题方案列表（JSON 数组）
    var titleOptions []model.TitleOption
    if err := json.Unmarshal([]byte(content), &titleOptions); err != nil {
        return fmt.Errorf("parse title options: %w", err)
    }

    state.TitleOptions = titleOptions
    return nil
}


// getStylePrompt 根据风格获取对应的 Prompt 附加内容
func (s *ArticleAgentService) getStylePrompt(style string) string {
	if style == "" {
		return ""
	}

	switch style {
	case common.ArticleStyleTech:
		return common.StyleTechPrompt
	case common.ArticleStyleEmotional:
		return common.StyleEmotionalPrompt
	case common.ArticleStyleEducational:
		return common.StyleEducationalPrompt
	case common.ArticleStyleHumorous:
		return common.StyleHumorousPrompt
	default:
		return ""
	}
}

// AiModifyOutline AI 修改大纲
func (s *ArticleAgentService) AiModifyOutline(ctx context.Context, mainTitle, subTitle string,
    currentOutline []model.OutlineSection, modifySuggestion string) ([]model.OutlineSection, error) {

    currentOutlineJSON, _ := json.Marshal(currentOutline)

    prompt := common.AiModifyOutlinePrompt
    prompt = strings.ReplaceAll(prompt, "{mainTitle}", mainTitle)
    prompt = strings.ReplaceAll(prompt, "{subTitle}", subTitle)
    prompt = strings.ReplaceAll(prompt, "{currentOutline}", string(currentOutlineJSON))
    prompt = strings.ReplaceAll(prompt, "{modifySuggestion}", modifySuggestion)

    content, err := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt)
    if err != nil {
        return nil, fmt.Errorf("LLM call failed: %w", err)
    }

    var outlineResult model.OutlineResult
    if err := json.Unmarshal([]byte(content), &outlineResult); err != nil {
        return nil, fmt.Errorf("parse outline: %w", err)
    }

    return outlineResult.Sections, nil
}
