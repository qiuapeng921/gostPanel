package model

import (
	"time"

	"gorm.io/gorm"
)

// NodeStatus 节点状态
type NodeStatus string

const (
	NodeStatusOnline  NodeStatus = "online"  // 在线
	NodeStatusOffline NodeStatus = "offline" // 离线
	NodeStatusError   NodeStatus = "error"   // 错误
)

// GostNode Gost 节点模型
type GostNode struct {
	ID       uint       `gorm:"primaryKey" json:"id"`
	Name     string     `gorm:"size:100;not null" json:"name"`         // 节点名称
	Address  string     `gorm:"size:255;not null" json:"address"`      // IP 或域名
	Port     int        `gorm:"not null" json:"port"`                  // 端口
	Username string     `gorm:"size:50" json:"username"`               // API 认证用户名
	Password string     `gorm:"size:255" json:"password"`              // API 认证密码
	Status   NodeStatus `gorm:"size:20;default:offline" json:"status"` // 状态

	// 流量统计
	TotalBytes  int64 `gorm:"default:0" json:"total_bytes"`
	InputBytes  int64 `gorm:"default:0" json:"input_bytes"`
	OutputBytes int64 `gorm:"default:0" json:"output_bytes"`

	LastCheckAt *time.Time     `json:"last_check_at"`           // 最后检查时间
	Remark      string         `gorm:"type:text" json:"remark"` // 备注
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联 - 规则
	Rules []GostRule `gorm:"foreignKey:NodeID" json:"rules,omitempty"`

	// 关联 - 隧道（作为入口或出口节点）
	EntryTunnels []GostTunnel `gorm:"foreignKey:EntryNodeID" json:"entry_tunnels,omitempty"`
	ExitTunnels  []GostTunnel `gorm:"foreignKey:ExitNodeID" json:"exit_tunnels,omitempty"`
}

// TableName 指定表名
func (GostNode) TableName() string {
	return "nodes"
}
