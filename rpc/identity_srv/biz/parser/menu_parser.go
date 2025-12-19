package parser

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"gopkg.in/yaml.v3"
)

// YamlMenuNode 定义了 menu.yaml 文件中单个节点的结构。
type YamlMenuNode struct {
	Name      string          `yaml:"name"`
	ID        string          `yaml:"id"` // 这是来自 ayml 的字符串ID，并非数据库的UUID主键
	Path      string          `yaml:"path"`
	Icon      string          `yaml:"icon"`
	Component string          `yaml:"component"`
	Children  []*YamlMenuNode `yaml:"children"`
}

// YamlMenuContainer 是 menu.yaml 文件的根对象结构。
type YamlMenuContainer struct {
	Menu []*YamlMenuNode `yaml:"menu"`
}

// ParseAndFlattenMenu 解析YAML格式的菜单内容，并将其扁平化为 models.Menu 的列表，
// 以便进行数据库插入。该函数会分配版本号并处理父子关系。
// 注意：此函数假定 models.Menu 结构体允许在代码中预设其 UUID。
func ParseAndFlattenMenu(yamlContent string, version string) ([]*models.Menu, error) {
	var container YamlMenuContainer
	if err := yaml.Unmarshal([]byte(yamlContent), &container); err != nil {
		return nil, fmt.Errorf("解析菜单YAML失败: %w", err)
	}

	if version == "" {
		return nil, fmt.Errorf("版本号不能为空")
	}

	flatList := make([]*models.Menu, 0)
	semanticIDSet := make(map[string]bool) // 用于检测同一版本中的重复语义ID

	if err := flattenNodes(container.Menu, nil, version, &flatList, semanticIDSet); err != nil {
		return nil, err
	}

	return flatList, nil
}

// flattenNodes 是一个递归辅助函数，用于遍历菜单树。
// 它将遍历结果填充到 flatList 中，生成一组可供入库的 models.Menu 对象。
func flattenNodes(
	nodes []*YamlMenuNode,
	parentID *uuid.UUID,
	version string,
	flatList *[]*models.Menu,
	semanticIDSet map[string]bool,
) error {
	for i, node := range nodes {
		if node.Name == "" || node.Path == "" {
			return fmt.Errorf("菜单节点缺少必要字段 (name, path): %+v", node)
		}

		if node.ID == "" {
			return fmt.Errorf("菜单节点缺少语义化ID: %+v", node)
		}

		// 检查语义ID重复
		if semanticIDSet[node.ID] {
			return fmt.Errorf("检测到重复的语义化ID: %s", node.ID)
		}

		semanticIDSet[node.ID] = true

		// 在应用层生成UUID，以便在插入前建立父子关系。
		newID := uuid.New()

		menuModel := &models.Menu{
			BaseModel: models.BaseModel{
				ID: newID,
			},
			SemanticID: node.ID, // 保存YAML中的语义化ID
			Version:    version,
			Name:       node.Name,
			Path:       node.Path,
			Component:  node.Component,
			Icon:       node.Icon,
			ParentID:   parentID,
			Sort:       i, // 使用切片索引作为排序依据
		}

		*flatList = append(*flatList, menuModel)

		if len(node.Children) > 0 {
			// 递归调用子节点，并将当前节点的新ID作为其父ID传入
			if err := flattenNodes(node.Children, &newID, version, flatList, semanticIDSet); err != nil {
				return err
			}
		}
	}

	return nil
}
