package menu

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

const (
	// DefaultProductLine 默认产品线标识
	DefaultProductLine = "default"
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

// GetMaxVersion 获取指定产品线的最大版本号
func (r *menuRepository) GetMaxVersion(ctx context.Context, productLine string) (int, error) {
	var maxVersion int

	err := r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Where("product_line = ?", productLine).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).
		Error
	if err != nil {
		return 0, err
	}

	return maxVersion, nil
}

// GetLatestContentHash 获取指定产品线最新版本的内容哈希
func (r *menuRepository) GetLatestContentHash(ctx context.Context, productLine string) (string, error) {
	// 先获取最新版本号
	maxVersion, err := r.GetMaxVersion(ctx, productLine)
	if err != nil {
		return "", err
	}

	if maxVersion == 0 {
		// 没有任何版本
		return "", nil
	}

	// 获取该版本的第一条记录的内容哈希（同一版本所有记录哈希相同）
	var contentHash string
	err = r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Where("product_line = ? AND version = ?", productLine, maxVersion).
		Limit(1).
		Pluck("content_hash", &contentHash).
		Error
	if err != nil {
		return "", err
	}

	return contentHash, nil
}

// GetLatestMenuTree retrieves the full menu tree for the most recent version of a product line.
func (r *menuRepository) GetLatestMenuTree(ctx context.Context, productLine string) ([]*models.Menu, error) {
	// Step 1: Find the latest version for this product line.
	maxVersion, err := r.GetMaxVersion(ctx, productLine)
	if err != nil {
		return nil, err
	}

	if maxVersion == 0 {
		// No menus found for this product line.
		return []*models.Menu{}, nil
	}

	// Step 2: Fetch all nodes for the latest version.
	var menus []*models.Menu

	err = r.db.WithContext(ctx).
		Where("product_line = ? AND version = ?", productLine, maxVersion).
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
		}
	}

	return rootNodes, nil
}

// GetBySemanticID 根据语义ID和版本查询菜单
func (r *menuRepository) GetBySemanticID(
	ctx context.Context,
	productLine string,
	semanticID string,
	version int,
) (*models.Menu, error) {
	if semanticID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	// 如果版本为0，获取最新版本
	if version == 0 {
		maxVersion, err := r.GetMaxVersion(ctx, productLine)
		if err != nil {
			return nil, err
		}

		version = maxVersion
	}

	var menu models.Menu

	err := r.db.WithContext(ctx).
		Where("product_line = ? AND semantic_id = ? AND version = ?", productLine, semanticID, version).
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
	productLine string,
	semanticIDs []string,
) ([]*models.Menu, error) {
	if len(semanticIDs) == 0 {
		return []*models.Menu{}, nil
	}

	// 获取最新版本
	maxVersion, err := r.GetMaxVersion(ctx, productLine)
	if err != nil {
		return nil, err
	}

	if maxVersion == 0 {
		return []*models.Menu{}, nil
	}

	var menus []*models.Menu

	err = r.db.WithContext(ctx).
		Where("product_line = ? AND semantic_id IN ? AND version = ?", productLine, semanticIDs, maxVersion).
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
	productLine string,
) (map[string]uuid.UUID, error) {
	// 获取最新版本
	maxVersion, err := r.GetMaxVersion(ctx, productLine)
	if err != nil {
		return nil, err
	}

	if maxVersion == 0 {
		return make(map[string]uuid.UUID), nil
	}

	// 查询最新版本的所有菜单
	var menus []struct {
		ID         uuid.UUID `gorm:"column:id"`
		SemanticID string    `gorm:"column:semantic_id"`
	}

	err = r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Select("id, semantic_id").
		Where("product_line = ? AND version = ?", productLine, maxVersion).
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
