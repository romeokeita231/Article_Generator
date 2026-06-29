package model

import (
    "encoding/json"
)

// parseJSON 解析 JSON 字符串到目标对象
func parseJSON(jsonStr string, target interface{}) {
    if jsonStr == "" {
        return
    }
    _ = json.Unmarshal([]byte(jsonStr), target)
}
