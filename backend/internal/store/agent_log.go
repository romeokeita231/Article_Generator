package store

import (
    
    "gorm.io/gorm"

	"github.com/romeokeita231/Article_Generator/internal/model"
)
// AgentLogStore 智能体日志数据访问
type AgentLogStore struct {
    db *gorm.DB
}

// NewAgentLogStore 创建智能体日志 Store
func NewAgentLogStore(db *gorm.DB) *AgentLogStore {
	return &AgentLogStore{db: db}
}

// Create 创建日志记录
func (s *AgentLogStore) Create(log *model.AgentLog) error {
    return s.db.Create(log).Error
}

// GetByTaskID 根据任务ID查询所有日志（按创建时间升序）
func (s *AgentLogStore) GetByTaskID(taskID string) ([]*model.AgentLog, error) {
    var logs []*model.AgentLog
    err := s.db.Where("taskId = ?", taskID).
        Where("isDelete = ?", 0).
        Order("createTime ASC").
        Find(&logs).Error
    return logs, err
}
