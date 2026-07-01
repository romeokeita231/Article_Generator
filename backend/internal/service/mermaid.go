package service

import (
	"os"
	"fmt"
	"time"
	"context"
	"strings"
	"os/exec"

	"github.com/romeokeita231/Article_Generator/internal/common"
	"github.com/romeokeita231/Article_Generator/internal/config"
	"github.com/romeokeita231/Article_Generator/internal/model"
)

// MermaidService Mermaid 流程图生成服务
// 使用 mermaid-cli 将 Mermaid 代码转换为图片
type MermaidService struct {
	config config.MermaidConfig
}

// NewMermaidService 创建 Mermaid 服务
func NewMermaidService(cfg config.MermaidConfig) *MermaidService {
	return &MermaidService{
		config: cfg,
	}
}

// GetMethod 返回方法名
func (s *MermaidService) GetMethod() string {
	return common.ImageMethodMermaid
}

// IsAvailable 是否可用（检查 mmdc 命令是否存在）
func (s *MermaidService) IsAvailable() bool {
	if s.config.CLI == "" {
		return false
	}
	// 检查 mmdc 命令是否可用
	_, err := exec.LookPath(s.config.CLI)
	return err == nil
}

// SearchImage Mermaid 是生成类服务，不实现此方法
func (s *MermaidService) SearchImage(keywords string) (string, error) {
	return "", fmt.Errorf("Mermaid 是生成类服务，请使用 GetImageData")
}

// GetImageData 生成 Mermaid 图表数据
func (s *MermaidService) GetImageData(req *model.ImageRequest) (*model.ImageData, error) {
	// 优先使用 Prompt（Mermaid 代码），否则使用 Keywords
	mermaidCode := req.GetEffectiveParam(true)
	return s.GenerateDiagramData(mermaidCode)
}

func (s *MermaidService) GenerateDiagramData(mermaidCode string) (*model.ImageData, error) {
    // 写入临时文件
    tmpInput, _ := os.CreateTemp("", "mermaid_input_*.mmd")
    defer os.Remove(tmpInput.Name())
    tmpInput.WriteString(mermaidCode)
    tmpInput.Close()

    tmpOutput, _ := os.CreateTemp("", "mermaid_output_*."+s.config.OutputFormat)
    tmpOutput.Close()
    defer os.Remove(tmpOutput.Name())

    // 调用 mmdc 命令
    if err := s.convertMermaidToImage(tmpInput.Name(), tmpOutput.Name()); err != nil {
        return nil, err
    }

    imageBytes, _ := os.ReadFile(tmpOutput.Name())
    mimeType := s.getMimeType(s.config.OutputFormat)
    return model.FromBytes(imageBytes, mimeType), nil
}

func (s *MermaidService) convertMermaidToImage(inputFile, outputFile string) error {
    args := []string{
        "-i", inputFile, "-o", outputFile,
        "-t", s.config.Theme,
        "-w", fmt.Sprintf("%d", s.config.Width),
        "-H", fmt.Sprintf("%d", s.config.Height),
    }
    ctx, cancel := context.WithTimeout(context.Background(),
        time.Duration(s.config.Timeout)*time.Millisecond)
    defer cancel()

    cmd := exec.CommandContext(ctx, s.config.CLI, args...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("Mermaid CLI 执行失败: %w, output: %s", err, string(output))
    }
    return nil
}



// getMimeType 根据输出格式获取 MIME 类型
func (s *MermaidService) getMimeType(format string) string {
	switch strings.ToLower(format) {
	case "png":
		return "image/png"
	case "svg":
		return "image/svg+xml"
	case "pdf":
		return "application/pdf"
	default:
		return "image/png"
	}
}