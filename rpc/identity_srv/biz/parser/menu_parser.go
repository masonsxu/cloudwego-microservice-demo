package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
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

// menuHashNode 用于哈希计算的规范化结构
type menuHashNode struct {
	SemanticID string         `json:"id"`
	Name       string         `json:"name"`
	Path       string         `json:"path"`
	Component  string         `json:"component"`
	Icon       string         `json:"icon"`
	Children   []menuHashNode `json:"children"`
}

// CalculateContentHash 计算菜单内容的 SHA256 哈希
// 用于去重校验，确保相同内容的菜单不会重复创建
func CalculateContentHash(nodes []*YamlMenuNode) string {
	// 1. 转换为规范化结构
	hashNodes := toHashNodes(nodes)

	// 2. JSON 序列化（确保字段顺序一致）
	data, err := json.Marshal(hashNodes)
	if err != nil {
		// 序列化失败时返回空字符串，后续会创建新版本
		return ""
	}

	// 3. 计算 SHA256
	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:])
}

// toHashNodes 将 YamlMenuNode 转换为用于哈希计算的规范化结构
func toHashNodes(nodes []*YamlMenuNode) []menuHashNode {
	if len(nodes) == 0 {
		return []menuHashNode{}
	}

	result := make([]menuHashNode, len(nodes))

	for i, n := range nodes {
		result[i] = menuHashNode{
			SemanticID: n.ID,
			Name:       n.Name,
			Path:       n.Path,
			Component:  n.Component,
			Icon:       n.Icon,
			Children:   toHashNodes(n.Children),
		}
	}

	// 按 SemanticID 排序确保稳定性
	sort.Slice(result, func(i, j int) bool {
		return result[i].SemanticID < result[j].SemanticID
	})

	return result
}

// ParseAndFlattenMenu 解析YAML格式的菜单内容，并将其扁平化为 models.Menu 的列表，
// 以便进行数据库插入。该函数会分配版本号并处理父子关系。
// 注意：此函数假定 models.Menu 结构体允许在代码中预设其 UUID。
func ParseAndFlattenMenu(yamlContent string, productLine string, version int) ([]*models.Menu, string, error) {
	var container YamlMenuContainer
	if err := yaml.Unmarshal([]byte(yamlContent), &container); err != nil {
		return nil, "", fmt.Errorf("解析菜单YAML失败: %w", err)
	}

	if productLine == "" {
		return nil, "", fmt.Errorf("产品线不能为空")
	}

	if version <= 0 {
		return nil, "", fmt.Errorf("版本号必须大于0")
	}

	// 计算内容哈希
	contentHash := CalculateContentHash(container.Menu)

	flatList := make([]*models.Menu, 0)
	semanticIDSet := make(map[string]bool) // 用于检测同一版本中的重复语义ID

	if err := flattenNodes(container.Menu, nil, productLine, version, contentHash, &flatList, semanticIDSet); err != nil {
		return nil, "", err
	}

	return flatList, contentHash, nil
}

// flattenNodes 是一个递归辅助函数，用于遍历菜单树。
// 它将遍历结果填充到 flatList 中，生成一组可供入库的 models.Menu 对象。
func flattenNodes(
	nodes []*YamlMenuNode,
	parentID *uuid.UUID,
	productLine string,
	version int,
	contentHash string,
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
			ProductLine: productLine,
			SemanticID:  node.ID,
			Version:     version,
			ContentHash: contentHash,
			Name:        node.Name,
			Path:        node.Path,
			Component:   node.Component,
			Icon:        node.Icon,
			ParentID:    parentID,
			Sort:        i,
		}

		*flatList = append(*flatList, menuModel)

		if len(node.Children) > 0 {
			// 递归调用子节点，并将当前节点的新ID作为其父ID传入
			if err := flattenNodes(node.Children, &newID, productLine, version, contentHash, flatList, semanticIDSet); err != nil {
				return err
			}
		}
	}

	return nil
}
