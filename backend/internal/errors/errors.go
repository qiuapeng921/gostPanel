// Package errors 定义统一的业务错误
// 所有业务错误都应该在此包中定义，便于统一管理和维护
package errors

import "net/http"

// BizError 业务错误结构
// 用于封装业务逻辑中的错误，包含错误码、错误消息和 HTTP 状态码
type BizError struct {
	Code     int    `json:"code"`    // 业务错误码
	Message  string `json:"message"` // 错误消息
	HTTPCode int    `json:"-"`       // HTTP 状态码（不序列化）
}

// Error 实现 error 接口
func (e *BizError) Error() string {
	return e.Message
}

// New 创建业务错误
func New(code int, message string, httpCode int) *BizError {
	return &BizError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
	}
}

// ==================== 节点相关错误 (100xx) ====================

var (
	// ErrNodeNotFound 节点不存在
	ErrNodeNotFound = New(10001, "节点不存在", http.StatusNotFound)
	// ErrNodeNameExists 节点名称已存在
	ErrNodeNameExists = New(10002, "节点名称已存在", http.StatusBadRequest)
	// ErrNodeHasRules 节点下存在规则
	ErrNodeHasRules = New(10003, "节点下存在规则，无法删除", http.StatusBadRequest)
	// ErrNodeHasTunnels 节点下存在隧道配置
	ErrNodeHasTunnels = New(10004, "节点下存在隧道配置，无法删除", http.StatusBadRequest)
	// ErrNodeHasObservers 节点下存在流量监控
	ErrNodeHasObservers = New(10005, "节点下存在流量监控，无法删除", http.StatusBadRequest)
	// ErrNodeOffline 节点已离线
	ErrNodeOffline = New(10006, "节点已离线", http.StatusBadRequest)
)

// ==================== 规则相关错误 (101xx) ====================

var (
	// ErrRuleNotFound 规则不存在
	ErrRuleNotFound = New(10101, "规则不存在", http.StatusNotFound)
	// ErrRulePortExists 端口已被使用
	ErrRulePortExists = New(10102, "端口已被使用", http.StatusBadRequest)
	// ErrRuleRunning 规则正在运行中
	ErrRuleRunning = New(10103, "规则正在运行中，请先停止", http.StatusBadRequest)
	// ErrRuleStartFailed 启动规则失败
	ErrRuleStartFailed = New(10104, "启动规则失败", http.StatusInternalServerError)
	// ErrTunnelRequired 隧道转发需要选择隧道
	ErrTunnelRequired = New(10105, "隧道转发类型需要选择隧道", http.StatusBadRequest)
	// ErrRuleChainCreateFailed 创建规则链失败
	ErrRuleChainCreateFailed = New(10106, "创建规则链失败", http.StatusInternalServerError)
	// ErrNodeRequired 端口转发需要选择节点
	ErrNodeRequired = New(10107, "端口转发类型需要选择节点", http.StatusBadRequest)
	// ErrRuleTypeInvalid 无效的规则类型
	ErrRuleTypeInvalid = New(10108, "无效的规则类型", http.StatusBadRequest)
	// ErrTunnelChainNotFound 隧道链不存在
	ErrTunnelChainNotFound = New(10109, "隧道未启动或链路不存在", http.StatusBadRequest)
)

// ==================== 隧道相关错误 (102xx) ====================

var (
	// ErrTunnelNotFound 隧道不存在
	ErrTunnelNotFound = New(10201, "隧道不存在", http.StatusNotFound)
	// ErrTunnelRunning 隧道正在运行中
	ErrTunnelRunning = New(10202, "隧道正在运行中，请先停止", http.StatusBadRequest)
	// ErrTunnelNameExists 隧道名称已存在
	ErrTunnelNameExists = New(10203, "隧道名称已存在", http.StatusBadRequest)
	// ErrTunnelInUse 隧道正在被规则使用
	ErrTunnelInUse = New(10209, "隧道正在被规则使用，无法删除", http.StatusBadRequest)
	// ErrTunnelRelayCreateFailed 创建隧道 Relay 服务失败
	ErrTunnelRelayCreateFailed = New(10214, "创建隧道 Relay 服务失败", http.StatusInternalServerError)
	// ErrTunnelHasRules 隧道被规则引用
	ErrTunnelHasRules = New(10215, "隧道被规则引用，无法删除", http.StatusBadRequest)
	// ErrEntryNodeOffline 入口节点离线
	ErrEntryNodeOffline = New(10216, "入口节点已离线", http.StatusBadRequest)
	// ErrExitNodeOffline 出口节点离线
	ErrExitNodeOffline = New(10217, "出口节点已离线", http.StatusBadRequest)
	// ErrTunnelChainCreateFailed 创建隧道 Chain 失败
	ErrTunnelChainCreateFailed = New(10218, "创建隧道Chain失败", http.StatusInternalServerError)
	// ErrTunnelObserverCreateFailed 创建观察器失败
	ErrTunnelObserverCreateFailed = New(10213, "创建观察器失败", http.StatusInternalServerError)
)

// ==================== 认证相关错误 (103xx) ====================

var (
	// ErrInvalidCredentials 用户名或密码错误
	ErrInvalidCredentials = New(10301, "用户名或密码错误", http.StatusUnauthorized)
	// ErrTokenExpired 登录已过期
	ErrTokenExpired = New(10302, "登录已过期，请重新登录", http.StatusUnauthorized)
	// ErrTokenInvalid Token 无效
	ErrTokenInvalid = New(10303, "Token 无效", http.StatusUnauthorized)
	// ErrPasswordMismatch 原密码错误
	ErrPasswordMismatch = New(10304, "原密码错误", http.StatusBadRequest)
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = New(10305, "用户不存在", http.StatusNotFound)
)

// ==================== 通用错误 (500xx) ====================

var (
	// ErrInternal 服务器内部错误
	ErrInternal = New(50000, "服务器内部错误", http.StatusInternalServerError)
	// ErrBadRequest 请求参数错误
	ErrBadRequest = New(40000, "请求参数错误", http.StatusBadRequest)
	// ErrOperationFailed 操作失败
	ErrOperationFailed = New(50001, "操作失败", http.StatusInternalServerError)
)

// ==================== 认证相关错误 (103xx) ====================

var (
	// ErrTokenGenerationFailed Token 生成失败
	ErrTokenGenerationFailed = New(10306, "Token 生成失败", http.StatusInternalServerError)
)

// ==================== 系统/配置相关错误 (104xx) ====================

var (
	// ErrPanelURLNotFound 未配置面板地址
	ErrPanelURLNotFound = New(10401, "未配置面板地址，无法启动流量监控，请先在[系统配置]中设置面板URL", http.StatusBadRequest)
	// ErrSMTPConfigIncomplete SMTP配置不完整
	ErrSMTPConfigIncomplete = New(10402, "SMTP配置不完整", http.StatusBadRequest)
	// ErrSMTPConnectFailed 连接SMTP服务器失败
	ErrSMTPConnectFailed = New(10403, "连接SMTP服务器失败", http.StatusInternalServerError)
	// ErrSMTPClientFailed 创建SMTP客户端失败
	ErrSMTPClientFailed = New(10404, "创建SMTP客户端失败", http.StatusInternalServerError)
	// ErrSMTPAuthFailed SMTP认证失败
	ErrSMTPAuthFailed = New(10405, "SMTP认证失败", http.StatusUnauthorized)
	// ErrSMTPSenderFailed 设置发件人失败
	ErrSMTPSenderFailed = New(10406, "设置发件人失败", http.StatusInternalServerError)
	// ErrSMTPRecipientFailed 设置收件人失败
	ErrSMTPRecipientFailed = New(10407, "设置收件人失败", http.StatusInternalServerError)
	// ErrSMTPDataFailed 创建邮件数据流失败
	ErrSMTPDataFailed = New(10408, "创建邮件数据流失败", http.StatusInternalServerError)
	// ErrSMTPWriteFailed 写入邮件内容失败
	ErrSMTPWriteFailed = New(10409, "写入邮件内容失败", http.StatusInternalServerError)
	// ErrSMTPCloseFailed 关闭邮件数据流失败
	ErrSMTPCloseFailed = New(10410, "关闭邮件数据流失败", http.StatusInternalServerError)

	// ErrDBPathNotConfigured 未配置数据库路径
	ErrDBPathNotConfigured = New(10411, "未配置数据库路径", http.StatusInternalServerError)
	// ErrBackupDirCreateFailed 创建备份目录失败
	ErrBackupDirCreateFailed = New(10412, "创建备份目录失败", http.StatusInternalServerError)
	// ErrBackupFailed 备份失败
	ErrBackupFailed = New(10413, "备份失败", http.StatusInternalServerError)

	// ErrObserverCreateFailed 创建观察器失败
	ErrObserverCreateFailed = New(10414, "创建流量监控失败", http.StatusInternalServerError)
	// ErrExtractHostFailed 提取主机IP失败
	ErrExtractHostFailed = New(10415, "无法从API地址提取主机IP", http.StatusInternalServerError)
)

// ==================== 隧道相关补全 (102xx) ====================

var (
	// ErrTunnelNodeSame 入口和出口节点不能相同
	ErrTunnelNodeSame = New(10204, "入口和出口节点不能相同", http.StatusBadRequest)
	// ErrEntryNodeNotFound 入口节点不存在
	ErrEntryNodeNotFound = New(10205, "入口节点不存在", http.StatusNotFound)
	// ErrExitNodeNotFound 出口节点不存在
	ErrExitNodeNotFound = New(10206, "出口节点不存在", http.StatusNotFound)
	// ErrTunnelNotRunning 隧道未运行
	ErrTunnelNotRunning = New(10207, "隧道未运行", http.StatusBadRequest)
)
