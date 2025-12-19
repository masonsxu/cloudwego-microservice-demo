package menu

import (
	"context"

	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"gorm.io/gorm"
)

// menuRepository implements the MenuRepository interface.
type menuRepository struct {
	// 嵌入基础仓储接口
	base.BaseRepository[*models.Menu]

	db *gorm.DB
}

// NewMenuRepository creates a new menu repository.
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{
		BaseRepository: base.NewBaseRepository[*models.Menu](db),
		db:             db,
	}
}

// CreateMenuTree saves a new menu tree (as a flat list) to the database in batches within a transaction.
func (r *menuRepository) CreateMenuTree(ctx context.Context, menus []*models.Menu) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Using CreateInBatches for efficient bulk insertion.
		if err := tx.CreateInBatches(menus, 100).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAllVersions retrieves a list of all unique menu version identifiers, sorted descending.
func (r *menuRepository) GetAllVersions(ctx context.Context) ([]string, error) {
	var versions []string

	err := r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Select("DISTINCT version").
		Order("version DESC").
		Pluck("version", &versions).
		Error

	return versions, err
}

// GetLatestMenuTree retrieves the full menu tree for the most recent version.
func (r *menuRepository) GetLatestMenuTree(ctx context.Context) ([]*models.Menu, error) {
	// Step 1: Find the latest version string based on creation time.
	var latestVersion string

	err := r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Order("created_at DESC").
		Limit(1).
		Pluck("version", &latestVersion).
		Error
	if err != nil {
		return nil, err
	}

	if latestVersion == "" {
		// No menus found in the database.
		return []*models.Menu{}, nil
	}

	// Step 2: Fetch all nodes for the latest version.
	var menus []*models.Menu

	err = r.db.WithContext(ctx).
		Where("version = ?", latestVersion).
		Order("sort ASC").
		Find(&menus).
		Error
	if err != nil {
		return nil, err
	}

	// Step 3: Build the tree from the flat list.
	menuMap := make(map[uuid.UUID]*models.Menu)

	var rootNodes []*models.Menu

	for _, menu := range menus {
		// Initialize children slice to avoid nil pointer issues later.
		menu.Children = []*models.Menu{}
		menuMap[menu.ID] = menu
	}

	for _, menu := range menus {
		if menu.ParentID == nil {
			rootNodes = append(rootNodes, menu)
		} else {
			if parent, ok := menuMap[*menu.ParentID]; ok {
				parent.Children = append(parent.Children, menu)
			}
			// Note: If a parent is not found, the node becomes an orphan and won't appear in the tree.
			// This indicates data integrity issues if it happens.
		}
	}

	return rootNodes, nil
}

// GetBySemanticID 根据语义ID和版本查询菜单
func (r *menuRepository) GetBySemanticID(
	ctx context.Context,
	semanticID string,
	version string,
) (*models.Menu, error) {
	if semanticID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	// 如果版本为空，获取最新版本
	if version == "" {
		latestVersion, err := r.getLatestVersion(ctx)
		if err != nil {
			return nil, err
		}

		version = latestVersion
	}

	var menu models.Menu

	err := r.db.WithContext(ctx).
		Where("semantic_id = ? AND version = ?", semanticID, version).
		First(&menu).
		Error
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// GetBySemanticIDs 批量根据语义ID查询菜单（使用最新版本）
func (r *menuRepository) GetBySemanticIDs(
	ctx context.Context,
	semanticIDs []string,
) ([]*models.Menu, error) {
	if len(semanticIDs) == 0 {
		return []*models.Menu{}, nil
	}

	// 获取最新版本
	latestVersion, err := r.getLatestVersion(ctx)
	if err != nil {
		return nil, err
	}

	var menus []*models.Menu

	err = r.db.WithContext(ctx).
		Where("semantic_id IN ? AND version = ?", semanticIDs, latestVersion).
		Find(&menus).
		Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// GetLatestSemanticIDMapping 获取最新版本的语义ID到UUID的映射
func (r *menuRepository) GetLatestSemanticIDMapping(
	ctx context.Context,
) (map[string]uuid.UUID, error) {
	// 获取最新版本
	latestVersion, err := r.getLatestVersion(ctx)
	if err != nil {
		return nil, err
	}

	// 查询最新版本的所有菜单
	var menus []struct {
		ID         uuid.UUID `gorm:"column:id"`
		SemanticID string    `gorm:"column:semantic_id"`
	}

	err = r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Select("id, semantic_id").
		Where("version = ?", latestVersion).
		Find(&menus).
		Error
	if err != nil {
		return nil, err
	}

	// 构建映射
	mapping := make(map[string]uuid.UUID)
	for _, menu := range menus {
		mapping[menu.SemanticID] = menu.ID
	}

	return mapping, nil
}

// getLatestVersion 获取最新版本号的辅助方法
func (r *menuRepository) getLatestVersion(ctx context.Context) (string, error) {
	var latestVersion string

	err := r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Order("created_at DESC").
		Limit(1).
		Pluck("version", &latestVersion).
		Error
	if err != nil {
		return "", err
	}

	if latestVersion == "" {
		return "", gorm.ErrRecordNotFound
	}

	return latestVersion, nil
}
