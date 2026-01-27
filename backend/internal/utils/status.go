package utils

import "gost-panel/internal/model"

// GostStateToRuleStatus 将 Gost 服务状态转换为规则状态
func GostStateToRuleStatus(state string) model.RuleStatus {
	switch state {
	case "ready", "running":
		return model.RuleStatusRunning
	case "failed":
		return model.RuleStatusError
	default:
		return model.RuleStatusStopped
	}
}

// GostStateToTunnelStatus 将 Gost 服务状态转换为隧道状态
func GostStateToTunnelStatus(state string) model.TunnelStatus {
	switch state {
	case "ready", "running":
		return model.TunnelStatusRunning
	case "failed":
		return model.TunnelStatusError
	default:
		return model.TunnelStatusStopped
	}
}
