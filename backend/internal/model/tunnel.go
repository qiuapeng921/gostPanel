package model

import (
	"time"

	"gorm.io/gorm"
)

// TunnelStatus 隧道状态
type TunnelStatus string

const (
	TunnelStatusStopped TunnelStatus = "stopped" // 已停止
	TunnelStatusRunning TunnelStatus = "running" // 运行中
	TunnelStatusError   TunnelStatus = "error"   // 错误
)

// GostTunnel 隧道模型 - 管理入口节点与出口节点的链路关系
// 启动隧道时：在出口节点创建 Relay 服务，在入口节点创建 Chain 连接到出口节点
type GostTunnel struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"size:100;not null" json:"name"`       // 隧道名称
	EntryNodeID uint         `gorm:"not null;index" json:"entry_node_id"` // 入口节点 ID
	ExitNodeID  uint         `gorm:"not null;index" json:"exit_node_id"`  // 出口节点 ID
	Protocol    string       `gorm:"size:10;default:tcp" json:"protocol"` // 协议 (tcp/udp)
	RelayPort   int          `gorm:"default:8443" json:"relay_port"`      // 出口节点 Relay 服务端口
	Status      TunnelStatus `gorm:"size:20;default:stopped" json:"status"`

	// Gost 服务相关 ID（启动时创建）
	ServiceID string `gorm:"size:100" json:"service_id"` // 出口节点 Relay 服务 ID
	ChainID   string `gorm:"size:100" json:"chain_id"`   // 入口节点 Chain ID

	Remark    string         `gorm:"type:text" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联 - 入口节点
	EntryNode *GostNode `gorm:"foreignKey:EntryNodeID" json:"entry_node,omitempty"`
	// 关联 - 出口节点
	ExitNode *GostNode `gorm:"foreignKey:ExitNodeID" json:"exit_node,omitempty"`
	// 关联 - 使用该隧道的规则
	Rules []GostRule `gorm:"foreignKey:TunnelID" json:"rules,omitempty"`
}

// TableName 指定表名
func (GostTunnel) TableName() string {
	return "tunnels"
}
