package model

import (
	"time"

	"gorm.io/gorm"
)

// RuleStatus 规则状态
type RuleStatus string

const (
	RuleStatusRunning RuleStatus = "running" // 运行中
	RuleStatusStopped RuleStatus = "stopped" // 已停止
	RuleStatusError   RuleStatus = "error"   // 错误
)

// RuleProtocol 规则协议
type RuleProtocol string

const (
	RuleProtocolTCP RuleProtocol = "tcp" // TCP 协议
	RuleProtocolUDP RuleProtocol = "udp" // UDP 协议
)

// RuleType 规则类型
type RuleType string

const (
	RuleTypeForward RuleType = "forward" // 端口转发（直连目标）
	RuleTypeTunnel  RuleType = "tunnel"  // 隧道转发（通过隧道链路）
)

// GostRule 转发规则模型
// 入口选择：NodeID 或 TunnelID 二选一
// - 端口转发 (forward)：选择 NodeID，直接在该节点上创建转发服务
// - 隧道转发 (tunnel)：选择 TunnelID，在隧道的入口节点上创建转发服务，使用隧道的 Chain
type GostRule struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	NodeID     *uint        `gorm:"index" json:"node_id"`                         // 入口节点 ID（端口转发时使用）
	Name       string       `gorm:"size:100;not null" json:"name"`                // 规则名称
	Type       RuleType     `gorm:"size:20;not null;default:forward" json:"type"` // 规则类型
	TunnelID   *uint        `gorm:"index" json:"tunnel_id"`                       // 隧道 ID（隧道转发时使用）
	Protocol   RuleProtocol `gorm:"size:10;not null;default:tcp" json:"protocol"` // 协议
	ListenPort int          `gorm:"not null" json:"listen_port"`                  // 监听端口

	Targets   []string   `gorm:"type:json;serializer:json" json:"targets"` // 多目标列表 (host:port)
	Strategy  string     `gorm:"size:20;default:round" json:"strategy"`    // 负载均衡策略 (round, random, fifo)
	EnableTLS bool       `gorm:"default:false" json:"enable_tls"`          // 是否启用 TLS
	Status    RuleStatus `gorm:"size:20;default:stopped" json:"status"`    // 状态
	ServiceID string     `gorm:"size:100" json:"service_id"`               // Gost 服务 ID

	// 流量监控配置
	ObserverID string `gorm:"size:100" json:"observer_id"` // 观察器 ID

	// 流量统计 (由观察器更新)
	InputBytes    int64 `gorm:"default:0" json:"input_bytes"`    // 入站总流量 (bytes)
	OutputBytes   int64 `gorm:"default:0" json:"output_bytes"`   // 出站总流量 (bytes)
	TotalBytes    int64 `gorm:"default:0" json:"total_bytes"`    // 总流量 (Input + Output)
	TotalRequests int64 `gorm:"default:0" json:"total_requests"` // 总请求数

	Remark    string         `gorm:"type:text" json:"remark"` // 备注
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联 - 入口节点
	Node *GostNode `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	// 关联 - 隧道
	Tunnel *GostTunnel `gorm:"foreignKey:TunnelID" json:"tunnel,omitempty"`
}

// TableName 指定表名
func (GostRule) TableName() string {
	return "rules"
}
