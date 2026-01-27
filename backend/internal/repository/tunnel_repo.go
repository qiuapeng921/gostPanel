package repository

import (
	"gost-panel/internal/model"

	"gorm.io/gorm"
)

// TunnelRepository 隧道仓库
type TunnelRepository struct {
	*BaseRepository
}

// NewTunnelRepository 创建隧道仓库
func NewTunnelRepository(db *gorm.DB) *TunnelRepository {
	return &TunnelRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建隧道
func (r *TunnelRepository) Create(tunnel *model.GostTunnel) error {
	return r.DB.Create(tunnel).Error
}

// Update 更新隧道
func (r *TunnelRepository) Update(tunnel *model.GostTunnel) error {
	return r.DB.Save(tunnel).Error
}

// Delete 删除隧道
func (r *TunnelRepository) Delete(id uint) error {
	return r.DB.Delete(&model.GostTunnel{}, id).Error
}

// FindByID 根据 ID 查询隧道（包含关联节点）
func (r *TunnelRepository) FindByID(id uint) (*model.GostTunnel, error) {
	var tunnel model.GostTunnel
	err := r.DB.Preload("EntryNode").Preload("ExitNode").First(&tunnel, id).Error
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}

// List 查询隧道列表
func (r *TunnelRepository) List(opt *QueryOption) ([]model.GostTunnel, int64, error) {
	var tunnels []model.GostTunnel
	var total int64

	db := r.DB.Model(&model.GostTunnel{})

	// 应用条件过滤
	db = ApplyConditions(db, opt)

	// 统计总数（包含过滤条件）
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 预加载节点和关联规则
	db = db.Preload("EntryNode").Preload("ExitNode")

	// 默认按创建时间倒序
	if opt == nil || len(opt.Orders) == 0 {
		db = db.Order("created_at DESC")
	}

	// 应用分页
	db = ApplyPagination(db, opt)

	if err := db.Find(&tunnels).Error; err != nil {
		return nil, 0, err
	}

	return tunnels, total, nil
}

// UpdateStatus 更新隧道状态
func (r *TunnelRepository) UpdateStatus(id uint, status model.TunnelStatus) error {
	return r.UpdateField(&model.GostTunnel{}, id, "status", status)
}

// CountAll 统计总数
func (r *TunnelRepository) CountAll() (int64, error) {
	var count int64
	err := r.DB.Model(&model.GostTunnel{}).Count(&count).Error
	return count, err
}

// FindByNodeID 查找节点相关的隧道
func (r *TunnelRepository) FindByNodeID(nodeID uint) ([]model.GostTunnel, error) {
	var tunnels []model.GostTunnel
	err := r.DB.Where("entry_node_id = ? OR exit_node_id = ?", nodeID, nodeID).Find(&tunnels).Error
	return tunnels, err
}

// StopByNodeID 停止与该节点相关的所有隧道
func (r *TunnelRepository) StopByNodeID(nodeID uint) error {
	return r.DB.Model(&model.GostTunnel{}).
		Where("(entry_node_id = ? OR exit_node_id = ?) AND status = ?", nodeID, nodeID, model.TunnelStatusRunning).
		Update("status", model.TunnelStatusStopped).Error
}

// HasRulesUsingTunnel 检查是否有规则正在使用该隧道
func (r *TunnelRepository) HasRulesUsingTunnel(tunnelID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&model.GostRule{}).Where("tunnel_id = ?", tunnelID).Count(&count).Error
	return count > 0, err
}

// HasRules 检查隧道是否被规则使用（别名方法）
func (r *TunnelRepository) HasRules(tunnelID uint) (bool, error) {
	return r.HasRulesUsingTunnel(tunnelID)
}

// CountByStatus 按状态统计
func (r *TunnelRepository) CountByStatus(status model.TunnelStatus) (int64, error) {
	var count int64
	err := r.DB.Model(&model.GostTunnel{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// UpdateServiceInfo 更新隧道的服务 ID 和 Chain ID
func (r *TunnelRepository) UpdateServiceInfo(id uint, serviceID, chainID string) error {
	return r.DB.Model(&model.GostTunnel{}).Where("id = ?", id).
		Updates(map[string]any{
			"service_id": serviceID,
			"chain_id":   chainID,
		}).Error
}
