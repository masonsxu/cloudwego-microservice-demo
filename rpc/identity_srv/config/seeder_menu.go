package config

import (
	_ "embed"
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/parser"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

//go:embed seed_menu.yaml
var seedMenuYAML string

const (
	// seedProductLine 种子数据使用的产品线标识
	seedProductLine = "default"

	// seedMenuVersion 种子数据使用的初始版本号
	seedMenuVersion = 1
)

// seedDefaultMenus 初始化默认菜单数据
// 仅在数据库中完全没有菜单数据时插入，已有数据一律跳过
func seedDefaultMenus(db *gorm.DB) error {
	log.Println("正在检查默认菜单数据...")

	// 检查是否已有菜单数据
	var count int64
	if err := db.Model(&models.Menu{}).Count(&count).Error; err != nil {
		return fmt.Errorf("查询菜单数据失败: %w", err)
	}

	if count > 0 {
		log.Printf("菜单数据已是最新（共 %d 条），跳过初始化", count)
		return nil
	}

	// 解析嵌入的 YAML 菜单配置
	menus, _, err := parser.ParseAndFlattenMenu(seedMenuYAML, seedProductLine, seedMenuVersion)
	if err != nil {
		return fmt.Errorf("解析菜单配置失败: %w", err)
	}

	// 批量插入菜单数据
	if err := db.Create(&menus).Error; err != nil {
		return fmt.Errorf("插入菜单数据失败: %w", err)
	}

	log.Printf("默认菜单数据初始化成功（共 %d 条）", len(menus))

	return nil
}
