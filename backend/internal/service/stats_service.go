package service

import (
	"gost-panel/internal/config"
	"gost-panel/internal/model"
	"gost-panel/internal/repository"

	"gorm.io/gorm"
)

// StatsService 统计服务
type StatsService struct {
	nodeRepo   *repository.NodeRepository
	ruleRepo   *repository.RuleRepository
	tunnelRepo *repository.TunnelRepository
	logRepo    *repository.OperationLogRepository
}

// NewStatsService 创建统计服务
func NewStatsService(db *gorm.DB) *StatsService {
	return &StatsService{
		nodeRepo:   repository.NewNodeRepository(db),
		ruleRepo:   repository.NewRuleRepository(db),
		tunnelRepo: repository.NewTunnelRepository(db),
		logRepo:    repository.NewOperationLogRepository(db),
	}
}

// DashboardStats 仪表盘统计
type DashboardStats struct {
	Nodes   NodeStats   `json:"nodes"`
	Rules   RuleStats   `json:"rules"`
	Tunnels TunnelStats `json:"tunnels"`
	Version string      `json:"version"`
}

// NodeStats 节点统计
type NodeStats struct {
	Total   int64 `json:"total"`
	Online  int64 `json:"online"`
	Offline int64 `json:"offline"`
}

// RuleStats 规则统计
type RuleStats struct {
	Total       int64 `json:"total"`
	Running     int64 `json:"running"`
	Stopped     int64 `json:"stopped"`
	ForwardType int64 `json:"forward_type"` // 端口转发类型数量
	TunnelType  int64 `json:"tunnel_type"`  // 隧道转发类型数量
}

// TunnelStats 隧道统计
type TunnelStats struct {
	Total   int64 `json:"total"`
	Running int64 `json:"running"`
	Stopped int64 `json:"stopped"`
}

// GetDashboardStats 获取仪表盘统计
func (s *StatsService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 节点统计
	nodeTotal, err := s.nodeRepo.CountAll()
	if err != nil {
		return nil, err
	}
	nodeOnline, err := s.nodeRepo.CountByStatus("online")
	if err != nil {
		return nil, err
	}
	stats.Nodes = NodeStats{
		Total:   nodeTotal,
		Online:  nodeOnline,
		Offline: nodeTotal - nodeOnline,
	}

	// 规则统计
	ruleTotal, err := s.ruleRepo.CountAll()
	if err != nil {
		return nil, err
	}
	ruleRunning, err := s.ruleRepo.CountByStatus(model.RuleStatusRunning)
	if err != nil {
		return nil, err
	}
	forwardType, err := s.ruleRepo.CountByType(model.RuleTypeForward)
	if err != nil {
		return nil, err
	}
	tunnelType, err := s.ruleRepo.CountByType(model.RuleTypeTunnel)
	if err != nil {
		return nil, err
	}
	stats.Rules = RuleStats{
		Total:       ruleTotal,
		Running:     ruleRunning,
		Stopped:     ruleTotal - ruleRunning,
		ForwardType: forwardType,
		TunnelType:  tunnelType,
	}

	// 隧道统计
	tunnelTotal, err := s.tunnelRepo.CountAll()
	if err != nil {
		return nil, err
	}
	tunnelRunning, err := s.tunnelRepo.CountByStatus(model.TunnelStatusRunning)
	if err != nil {
		return nil, err
	}
	stats.Tunnels = TunnelStats{
		Total:   tunnelTotal,
		Running: tunnelRunning,
		Stopped: tunnelTotal - tunnelRunning,
	}

	stats.Version = config.Version

	return stats, nil
}
