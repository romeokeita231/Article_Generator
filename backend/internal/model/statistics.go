package model

// AgentExecutionStats 智能体执行统计
type AgentExecutionStats struct {
    TaskID          string         `json:"taskId"`
    TotalDurationMs int            `json:"totalDurationMs"`
    AgentCount      int            `json:"agentCount"`
    AgentDurations  map[string]int `json:"agentDurations"` // key: agentName, value: durationMs
    OverallStatus   string         `json:"overallStatus"`  // SUCCESS/FAILED/RUNNING
    Logs            []*AgentLog    `json:"logs"`
}

// StatisticsVO 统计数据 VO
type StatisticsVO struct {
    TodayCount      int64   `json:"todayCount"`      // 今日创作数量
    WeekCount       int64   `json:"weekCount"`       // 本周创作数量
    MonthCount      int64   `json:"monthCount"`      // 本月创作数量
    TotalCount      int64   `json:"totalCount"`      // 总创作数量
    SuccessRate     float64 `json:"successRate"`     // 成功率（百分比）
    AvgDurationMs   int     `json:"avgDurationMs"`   // 平均耗时（毫秒）
    ActiveUserCount int64   `json:"activeUserCount"` // 活跃用户数（本周）
    TotalUserCount  int64   `json:"totalUserCount"`  // 总用户数
    VipUserCount    int64   `json:"vipUserCount"`    // VIP 用户数
    QuotaUsed       int64   `json:"quotaUsed"`       // 配额总使用量
}
