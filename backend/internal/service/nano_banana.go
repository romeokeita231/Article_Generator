package service

import (
	"context"
	"fmt"
	
	"google.golang.org/genai"
	"github.com/romeokeita231/Article_Generator/internal/model"
	"github.com/romeokeita231/Article_Generator/internal/config"
	"github.com/romeokeita231/Article_Generator/internal/common"
)

// NanoBananaService Nano Banana AI 生图服务（基于 Gemini API）
type NanoBananaService struct {
	config config.NanoBananaConfig
}

// NewNanoBananaService 创建 Nano Banana 服务
func NewNanoBananaService(cfg config.NanoBananaConfig) *NanoBananaService {
	return &NanoBananaService{
		config: cfg,
	}
}

// GetMethod 返回方法名
func (s *NanoBananaService) GetMethod() string {
	return common.ImageMethodNanoBanana
}

// IsAvailable 是否可用
func (s *NanoBananaService) IsAvailable() bool {
	return s.config.APIKey != "" && s.config.Model != ""
}

// SearchImage NanoBanana 是生成类服务，不实现此方法
func (s *NanoBananaService) SearchImage(keywords string) (string, error) {
	return "", fmt.Errorf("NanoBanana 是生成类服务，请使用 GetImageData")
}

// GetImageData 生成图片数据
func (s *NanoBananaService) GetImageData(req *model.ImageRequest) (*model.ImageData, error) {
	prompt := req.GetEffectiveParam(true)
	return s.GenerateImageData(prompt)
}

func (s *NanoBananaService) GenerateImageData(prompt string) (*model.ImageData, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: s.config.APIKey})
    if err != nil {
        return nil, fmt.Errorf("创建Gemini客户端失败: %w", err)
    }

    contents := []*genai.Content{{
        Role:  genai.RoleUser,
        Parts: []*genai.Part{{Text: prompt}},
    }}

    result, err := client.Models.GenerateContent(ctx, s.config.Model, contents, nil)
    if err != nil {
        return nil, fmt.Errorf("生成图片失败: %w", err)
    }

    // 提取 InlineData（图片字节数据）
    for _, part := range result.Candidates[0].Content.Parts {
        if part.InlineData != nil && len(part.InlineData.Data) > 0 {
            mimeType := part.InlineData.MIMEType
            if mimeType == "" {
                mimeType = "image/png"
            }
            return model.FromBytes(part.InlineData.Data, mimeType), nil
        }
    }

    return nil, fmt.Errorf("Nano Banana 未生成图片")
}
