package menu

import (
	"context"

	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// MenuRepository defines the data access layer for menu-related operations.
type MenuRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[*models.Menu]

	// CreateMenuTree saves a new menu tree (provided as a flat list of nodes) to the database.
	// It's recommended to perform this operation within a transaction.
	CreateMenuTree(ctx context.Context, menus []*models.Menu) error

	// GetMaxVersion 获取指定产品线的最大版本号
	GetMaxVersion(ctx context.Context, productLine string) (int, error)

	// GetLatestContentHash 获取指定产品线最新版本的内容哈希
	GetLatestContentHash(ctx context.Context, productLine string) (string, error)

	// GetLatestMenuTree retrieves the full menu tree for the most recent version of a product line.
	GetLatestMenuTree(ctx context.Context, productLine string) ([]*models.Menu, error)

	// GetBySemanticID 根据语义ID和版本查询菜单
	// productLine: 产品线标识
	// semanticID: 语义化菜单ID (来自menu.yaml)
	// version: 菜单版本，如果为0则使用最新版本
	GetBySemanticID(ctx context.Context, productLine string, semanticID string, version int) (*models.Menu, error)

	// GetBySemanticIDs 批量根据语义ID查询菜单（使用最新版本）
	// productLine: 产品线标识
	// semanticIDs: 语义化菜单ID列表
	GetBySemanticIDs(ctx context.Context, productLine string, semanticIDs []string) ([]*models.Menu, error)

	// GetLatestSemanticIDMapping 获取最新版本的语义ID到UUID的映射
	// productLine: 产品线标识
	GetLatestSemanticIDMapping(ctx context.Context, productLine string) (map[string]uuid.UUID, error)
}
