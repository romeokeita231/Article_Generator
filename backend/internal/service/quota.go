package service

import (
	"log"
	"github.com/romeokeita231/Article_Generator/internal/model"
	"github.com/romeokeita231/Article_Generator/internal/common"
	"github.com/romeokeita231/Article_Generator/internal/store"
)

type QuotaService struct {
    userStore *store.UserStore
}

// NewQuotaService 创建配额服务
func NewQuotaService(userStore *store.UserStore) *QuotaService {
	return &QuotaService{userStore: userStore}
}


// CheckAndConsumeQuota 检查并消耗配额（原子操作）
func (s *QuotaService) CheckAndConsumeQuota(user *model.User) error {
    // 管理员跳过检查
    if s.isAdmin(user) {
        return nil
    }

    // 原子更新：检查与消费合并为一个 SQL
    affectedRows, err := s.userStore.DecrementQuota(user.ID)
    if err != nil {
        return common.ErrSystem
    }

    if affectedRows == 0 {
        // 影响行数为0，说明配额不足
        return common.ErrOperation.WithMessage("配额不足，无法创建文章")
    }

    log.Printf("用户配额检查并消耗成功, userId=%d", user.ID)
    return nil
}

// isAdmin 判断是否为管理员
func (s *QuotaService) isAdmin(user *model.User) bool {
	return user.UserRole == common.AdminRole
}