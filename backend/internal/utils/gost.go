package utils

import (
	"fmt"
	"time"

	"gost-panel/internal/model"
	"gost-panel/pkg/gost"
)

// GetGostClient 根据节点配置创建 Gost 客户端
func GetGostClient(node *model.GostNode) *gost.Client {
	return gost.NewClient(&gost.Config{
		APIURL:   fmt.Sprintf("http://%s:%d", node.Address, node.Port),
		Username: node.Username,
		Password: node.Password,
		Timeout:  5 * time.Second, // 统一设置超时时间
	})
}
