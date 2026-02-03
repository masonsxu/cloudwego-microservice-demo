package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestParseAndFlattenMenu_WithSemanticID(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据"
    id: "patient_data"
    path: "/patient-data"
    icon: "IconPatientData"
    component: "Layout"
    children:
      - name: "通用字典"
        id: "common_dictionary"
        path: "common-dictionary"
        icon: "IconCommonDictionary"
        component: "views/patient-data/CommonDictionary"
      - name: "肺癌专病库"
        id: "lung_cancer_special"
        path: "lung-cancer-special"
        icon: "IconLungCancer"
        component: "views/patient-data/LungCancerSpecial"
  - name: "系统设置"
    id: "system_settings"
    path: "/system-settings"
    icon: "IconSystemSettings"
    component: "Layout"
`

	productLine := "default"
	version := 1

	// 测试解析和扁平化
	menus, contentHash, err := ParseAndFlattenMenu(yamlContent, productLine, version)
	require.NoError(t, err)
	require.NotEmpty(t, contentHash)
	require.Len(t, menus, 4) // 2个根菜单 + 2个子菜单

	// 验证根菜单
	patientDataMenu := findMenuBySemanticID(menus, "patient_data")
	require.NotNil(t, patientDataMenu)
	assert.Equal(t, "patient_data", patientDataMenu.SemanticID)
	assert.Equal(t, "患者数据", patientDataMenu.Name)
	assert.Equal(t, "/patient-data", patientDataMenu.Path)
	assert.Equal(t, version, patientDataMenu.Version)
	assert.Nil(t, patientDataMenu.ParentID) // 根菜单没有父ID

	// 验证子菜单
	commonDictMenu := findMenuBySemanticID(menus, "common_dictionary")
	require.NotNil(t, commonDictMenu)
	assert.Equal(t, "common_dictionary", commonDictMenu.SemanticID)
	assert.Equal(t, "通用字典", commonDictMenu.Name)
	assert.Equal(t, "common-dictionary", commonDictMenu.Path)
	assert.Equal(t, version, commonDictMenu.Version)
	assert.NotNil(t, commonDictMenu.ParentID) // 子菜单有父ID
	assert.Equal(t, patientDataMenu.ID, *commonDictMenu.ParentID)

	// 验证另一个子菜单
	lungCancerMenu := findMenuBySemanticID(menus, "lung_cancer_special")
	require.NotNil(t, lungCancerMenu)
	assert.Equal(t, "lung_cancer_special", lungCancerMenu.SemanticID)
	assert.Equal(t, "肺癌专病库", lungCancerMenu.Name)
	assert.Equal(t, patientDataMenu.ID, *lungCancerMenu.ParentID)

	// 验证系统设置菜单
	systemMenu := findMenuBySemanticID(menus, "system_settings")
	require.NotNil(t, systemMenu)
	assert.Equal(t, "system_settings", systemMenu.SemanticID)
	assert.Equal(t, "系统设置", systemMenu.Name)
	assert.Equal(t, "/system-settings", systemMenu.Path)
	assert.Nil(t, systemMenu.ParentID) // 根菜单没有父ID
}

func TestParseAndFlattenMenu_EmptySemanticID(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据"
    id: ""
    path: "/patient-data"
`

	productLine := "default"
	version := 1

	_, _, err := ParseAndFlattenMenu(yamlContent, productLine, version)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "菜单节点缺少语义化ID")
}

func TestParseAndFlattenMenu_DuplicateSemanticID(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据"
    id: "patient_data"
    path: "/patient-data"
  - name: "患者数据2"
    id: "patient_data"
    path: "/patient-data2"
`

	productLine := "default"
	version := 1

	_, _, err := ParseAndFlattenMenu(yamlContent, productLine, version)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "检测到重复的语义化ID")
}

func TestParseAndFlattenMenu_EmptyProductLine(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据"
    id: "patient_data"
    path: "/patient-data"
`

	_, _, err := ParseAndFlattenMenu(yamlContent, "", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "产品线不能为空")
}

func TestParseAndFlattenMenu_InvalidVersion(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据"
    id: "patient_data"
    path: "/patient-data"
`

	_, _, err := ParseAndFlattenMenu(yamlContent, "default", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "版本号必须大于0")
}

func TestParseAndFlattenMenu_InvalidYAML(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据"
    id: "patient_data"
`

	productLine := "default"
	version := 1

	_, _, err := ParseAndFlattenMenu(yamlContent, productLine, version)
	assert.Error(t, err)
}

// findMenuBySemanticID 辅助函数，根据语义ID查找菜单
func findMenuBySemanticID(menus []*models.Menu, semanticID string) *models.Menu {
	for _, menu := range menus {
		if menu.SemanticID == semanticID {
			return menu
		}
	}

	return nil
}
