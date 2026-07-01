package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"github.com/romeokeita231/Article_Generator/internal/common"
	"github.com/romeokeita231/Article_Generator/internal/config"
	"github.com/romeokeita231/Article_Generator/internal/model"
)

// PexelsService Pexels 图片检索服务
type PexelsService struct {
    apiKey string
    client *http.Client
}

// GetMethod 返回方法名
func (s *PexelsService) GetMethod() string {
	return common.ImageMethodPexels
}

// IsAvailable 是否可用
func (s *PexelsService) IsAvailable() bool {
	return s.apiKey != ""
}

func NewPexelsService(cfg *config.Config) *PexelsService {
    return &PexelsService{apiKey: cfg.Pexels.APIKey, client: &http.Client{}}
}

func (s *PexelsService) SearchImage(keywords string) (string, error) {
    apiURL := fmt.Sprintf("https://api.pexels.com/v1/search?query=%s&per_page=1",
        url.QueryEscape(keywords))

    req, _ := http.NewRequest("GET", apiURL, nil)
    req.Header.Set("Authorization", s.apiKey)

    resp, err := s.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var result struct {
        Photos []struct {
            Src struct{ Large string `json:"large"` } `json:"src"`
        } `json:"photos"`
    }
    json.Unmarshal(body, &result)

    if len(result.Photos) == 0 {
        return "", fmt.Errorf("no image found")
    }
    return result.Photos[0].Src.Large, nil
}

// GetImageData 获取图片数据（Pexels 是检索类服务，使用 SearchImage）
func (s *PexelsService) GetImageData(req *model.ImageRequest) (*model.ImageData, error) {
	url, err := s.SearchImage(req.Keywords)
	if err != nil {
		return nil, err
	}
	return model.FromURL(url), nil
}

// GetFallbackImage 获取降级图片（Picsum）
func (s *PexelsService) GetFallbackImage(position int) string {
    return fmt.Sprintf("https://picsum.photos/seed/%d/800/600", position)
}
