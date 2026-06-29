package common


// BaseResponse 统一响应结构
type BaseResponse struct {
    Code    int         `json:"code"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message"`
}

// Success 成功响应
func Success(data interface{}) BaseResponse {
    return BaseResponse{
        Code:    0,
        Data:    data,
        Message: "ok",
    }
}

// Error 错误响应
func Error(err *AppError) BaseResponse {
    return BaseResponse{
        Code:    err.Code,
        Data:    nil,
        Message: err.Message,
    }
}
