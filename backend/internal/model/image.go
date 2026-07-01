package model

// ImageData 图片数据（可以是字节数据或URL）
type ImageData struct {
    Bytes    []byte // 图片字节数据（适用于生成类服务：Mermaid、NanoBanana）
    URL      string // 图片URL（适用于检索类服务：Pexels、Iconify）
    MimeType string // MIME类型（如 image/png, image/svg+xml）
}

func (d *ImageData) IsValid() bool {
    return (len(d.Bytes) > 0) || (d.URL != "")
}

// FromBytes 从字节数据创建 ImageData
func FromBytes(bytes []byte, mimeType string) *ImageData {
	return &ImageData{
		Bytes:    bytes,
		MimeType: mimeType,
	}
}

// FromURL 从 URL 创建 ImageData
func FromURL(url string) *ImageData {
	return &ImageData{
		URL: url,
	}
}

// ImageRequest 图片请求（统一封装不同服务所需的参数）
type ImageRequest struct {
    Keywords string // 关键词（用于图片检索类服务）
    Prompt   string // 提示词（用于 AI 生成类服务）
    Position int    // 位置
    Type     string // 类型：cover/section
}

// GetEffectiveParam 根据服务类型获取有效参数
func (r *ImageRequest) GetEffectiveParam(usePromptFirst bool) string {
    if usePromptFirst {
        if r.Prompt != "" {
            return r.Prompt
        }
        return r.Keywords
    }
    if r.Keywords != "" {
        return r.Keywords
    }
    return r.Prompt
}

// ImageStrategyResult 图片策略结果
type ImageStrategyResult struct {
    URL    string // COS URL
    Method string // 图片方法
}
