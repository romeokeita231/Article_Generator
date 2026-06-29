package store

import (
	"gorm.io/gorm"

	"github.com/romeokeita231/Article_Generator/internal/model"
)

// ArticleStore 文章存储
type ArticleStore struct {
    db *gorm.DB
}

// NewArticleStore 创建文章存储
func NewArticleStore(db *gorm.DB) *ArticleStore {
    return &ArticleStore{db: db}
}

// Create 创建文章
func (s *ArticleStore) Create(article *model.Article) error {
    return s.db.Create(article).Error
}

// GetByTaskID 根据任务ID获取文章
func (s *ArticleStore) GetByTaskID(taskID string) (*model.Article, error) {
    var article model.Article
    err := s.db.Scopes(NotDeleted).Where("taskId = ?", taskID).First(&article).Error
    if err != nil {
        return nil, err
    }
    return &article, nil
}

// GetByID 根据ID获取文章
func (s *ArticleStore) GetByID(id int64) (*model.Article, error) {
	var article model.Article
	err := s.db.Scopes(NotDeleted).Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Update 更新文章
func (s *ArticleStore) Update(article *model.Article) error {
	return s.db.Scopes(NotDeleted).Where("taskId = ?", article.TaskID).Updates(article).Error
}

// UpdateStatus 更新文章状态
func (s *ArticleStore) UpdateStatus(taskID, status string, errorMsg *string) error {
    updates := map[string]interface{}{"status": status}
    if errorMsg != nil {
        updates["errorMessage"] = *errorMsg
    }
    return s.db.Model(&model.Article{}).Where("taskId = ?", taskID).Updates(updates).Error
}

// List 分页查询文章列表
func (s *ArticleStore) List(userID *int64, status *string, isAdmin bool,
    pageNum, pageSize int64) ([]model.Article, int64, error) {

    var articles []model.Article
    var total int64
    query := s.db.Scopes(NotDeleted)

    if !isAdmin && userID != nil {
        query = query.Where("userId = ?", *userID)
    } else if userID != nil {
        query = query.Where("userId = ?", *userID)
    }

    if status != nil && *status != "" {
        query = query.Where("status = ?", *status)
    }

    query.Model(&model.Article{}).Count(&total)

    offset := (pageNum - 1) * pageSize
    query.Order("createTime DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&articles)

    return articles, total, nil
}

// Delete 删除文章（逻辑删除）
func (s *ArticleStore) Delete(id int64) error {
	return s.db.Model(&model.Article{}).Where("id = ?", id).Update("isDelete", 1).Error
}