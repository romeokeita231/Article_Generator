package service

import (
	"log"
	"fmt"
	"math/rand"


	"github.com/romeokeita231/Article_Generator/internal/common"
	"github.com/romeokeita231/Article_Generator/internal/model"
)

// ImageServiceStrategy 图片服务策略
type ImageServiceStrategy struct {
    services   map[string]ImageService
    cos        *CosService
    cosEnabled bool
}

// NewImageServiceStrategy 创建图片服务策略
func NewImageServiceStrategy(cos *CosService, cosEnabled bool) *ImageServiceStrategy {
	return &ImageServiceStrategy{
		services:   make(map[string]ImageService),
		cos:        cos,
		cosEnabled: cosEnabled,
	}
}

// RegisterService 注册图片服务
func (s *ImageServiceStrategy) RegisterService(service ImageService) {
    if service == nil {
        return
    }
    method := service.GetMethod()
    s.services[method] = service
}

func (s *ImageServiceStrategy) GetImageAndUpload(imageSource string, req *model.ImageRequest) (*model.ImageStrategyResult, error) {
    method := s.resolveMethod(imageSource)
    service, exists := s.services[method]

    if !exists || !service.IsAvailable() {
        return s.handleFallback(req.Position)
    }

    imageData, err := s.getImage(service, req)
    if err != nil || imageData == nil || !imageData.IsValid() {
        return s.handleFallback(req.Position)
    }

    // COS 可用则上传，否则直接返回原始 URL
    if s.cosEnabled && s.cos != nil {
        folder := common.GetFolderName(method)
        cosURL, err := s.cos.UploadImageData(imageData, folder)
        if err == nil && cosURL != "" {
            return &model.ImageStrategyResult{URL: cosURL, Method: method}, nil
        }
    }

    if imageData.URL != "" {
        return &model.ImageStrategyResult{URL: imageData.URL, Method: method}, nil
    }

    return s.handleFallback(req.Position)
}

// resolveMethod 解析方法名（转大写，支持多种格式）
func (s *ImageServiceStrategy) resolveMethod(imageSource string) string {
	if imageSource == "" {
		return common.GetDefaultSearchMethod()
	}
	// 直接返回，假设已经是标准格式（PEXELS、NANO_BANANA 等）
	return imageSource
}

// getImage 获取图片（根据服务类型选择不同的获取方式）
func (s *ImageServiceStrategy) getImage(service ImageService, req *model.ImageRequest) (*model.ImageData, error) {
	method := service.GetMethod()

	// AI 生成类服务使用 GetImageData
	if common.IsAIGenerated(method) {
		return service.GetImageData(req)
	}

	// 网络检索类服务使用 SearchImage
	url, err := service.SearchImage(req.Keywords)
	if err != nil || url == "" {
		return nil, err
	}

	return model.FromURL(url), nil
}

// handleFallback 降级处理（支持 COS 和非 COS 模式）
func (s *ImageServiceStrategy) handleFallback(position int) (*model.ImageStrategyResult, error) {
	fallbackMethod := common.GetDefaultFallbackMethod()
	fallbackService, exists := s.services[fallbackMethod]

	if !exists || !fallbackService.IsAvailable() {
		return nil, fmt.Errorf("降级服务不可用: %s", fallbackMethod)
	}

	// 使用 Picsum 随机图片
	randomNum := rand.Intn(1000) + position
	url := fmt.Sprintf(common.PicsumURLTemplate, randomNum)

	// 如果 COS 可用，上传到 COS
	if s.cosEnabled && s.cos != nil {
		imageData := model.FromURL(url)
		folder := common.GetFolderName(fallbackMethod)
		cosURL, err := s.cos.UploadImageData(imageData, folder)
		if err != nil {
			log.Printf("降级图片上传 COS 失败, 使用原始 URL, err=%v", err)
			// COS 失败时返回原始 URL
			return &model.ImageStrategyResult{
				URL:    url,
				Method: fallbackMethod,
			}, nil
		}
		log.Printf("使用降级图片, method=%s, cosURL=%s", fallbackMethod, cosURL)
		return &model.ImageStrategyResult{
			URL:    cosURL,
			Method: fallbackMethod,
		}, nil
	}

	// COS 不可用，直接返回 Picsum URL
	log.Printf("COS 未启用, 使用降级图片原始 URL, method=%s, url=%s", fallbackMethod, url)
	return &model.ImageStrategyResult{
		URL:    url,
		Method: fallbackMethod,
	}, nil
}

