package common

import (
	"encoding/json"
	"time"
    "sync"
)

// SSEManager SSE 连接管理器
type SSEManager struct {
    clients map[string]chan string
    mu      sync.RWMutex
}

func NewSSEManager() *SSEManager {
    return &SSEManager{clients: make(map[string]chan string)}
}

func (m *SSEManager) Register(taskID string) chan string {
    m.mu.Lock()
    defer m.mu.Unlock()

    ch := make(chan string, 100) // 缓冲通道，避免阻塞
    m.clients[taskID] = ch
    return ch
}

func (m *SSEManager) Send(taskID string, data interface{}) {
    m.mu.RLock()
    ch, ok := m.clients[taskID]
    m.mu.RUnlock()
    if !ok {
        return
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return
    }

    select {
    case ch <- string(jsonData):
    case <-time.After(5 * time.Second):
        // 超时则放弃
    }
}

func (m *SSEManager) Unregister(taskID string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if ch, ok := m.clients[taskID]; ok {
        close(ch)
        delete(m.clients, taskID)
    }
}

func (m *SSEManager) Complete(taskID string) {
    m.Unregister(taskID)
}
