package service

import (
	"github.com/romeokeita231/Article_Generator/internal/model"
)

// ImageService 图片服务接口
type ImageService interface {
    GetMethod() string
    IsAvailable() bool

    // SearchImage 适用于网络检索类服务（Pexels、Iconify、EmojiPack）
    SearchImage(keywords string) (string, error)

    // GetImageData 适用于生成类服务（Mermaid、NanoBanana、SVG_DIAGRAM）
    GetImageData(req *model.ImageRequest) (*model.ImageData, error)
}
