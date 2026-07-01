package service

import (
	"fmt"
	"context"
	"strings"
	"encoding/json"

	"github.com/tmc/langchaingo/llms"
		
	"github.com/romeokeita231/Article_Generator/internal/common"
	"github.com/romeokeita231/Article_Generator/internal/model"
)

func (s *ArticleAgentService) agent2GenerateOutlineStream(ctx context.Context, state *model.ArticleState) error {
    // 根据是否有用户补充描述决定是否插入
    descriptionSection := ""
    if state.UserDescription != "" {
        descriptionSection = strings.ReplaceAll(
            common.Agent2DescriptionSection, "{userDescription}", state.UserDescription,
        )
    }

    prompt := strings.ReplaceAll(common.Agent2OutlinePrompt, "{mainTitle}", state.Title.MainTitle)
    prompt = strings.ReplaceAll(prompt, "{subTitle}", state.Title.SubTitle)
    prompt = strings.ReplaceAll(prompt, "{descriptionSection}", descriptionSection)
    prompt += s.getStylePrompt(state.Style)


    var contentBuilder strings.Builder

    _, err := s.llm.GenerateContent(ctx, []llms.MessageContent{
        llms.TextParts(llms.ChatMessageTypeHuman, prompt),
    }, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
        text := string(chunk)
        contentBuilder.WriteString(text)
        s.sendMessage(state.TaskID, map[string]interface{}{
            "type": "AGENT2_STREAMING", "content": text,
        })
        return nil
    }))
    if err != nil {
        return err
    }

    var outline model.OutlineResult
    if err := json.Unmarshal([]byte(contentBuilder.String()), &outline); err != nil {
        return fmt.Errorf("parse outline failed: %w", err)
    }

    state.Outline = &outline
    return nil
}
