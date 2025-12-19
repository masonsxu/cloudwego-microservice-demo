package parser

import (
	"testing"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	version := "v20250928"

	// 测试解析和扁平化
	menus, err := ParseAndFlattenMenu(yamlContent, version)
	require.NoError(t, err)
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
    icon: "IconPatientData"
    component: "Layout"
`

	version := "v20250928"

	// 测试空的语义ID应该报错
	_, err := ParseAndFlattenMenu(yamlContent, version)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "菜单节点缺少语义化ID")
}

func TestParseAndFlattenMenu_DuplicateSemanticID(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据1"
    id: "patient_data"
    path: "/patient-data-1"
    icon: "IconPatientData"
    component: "Layout"
  - name: "患者数据2"
    id: "patient_data"
    path: "/patient-data-2"
    icon: "IconPatientData"
    component: "Layout"
`

	version := "v20250928"

	// 测试重复的语义ID应该报错
	_, err := ParseAndFlattenMenu(yamlContent, version)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "检测到重复的语义化ID")
}

func TestParseAndFlattenMenu_EmptyVersion(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据"
    id: "patient_data"
    path: "/patient-data"
    icon: "IconPatientData"
    component: "Layout"
`

	// 测试空版本应该报错
	_, err := ParseAndFlattenMenu(yamlContent, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "版本号不能为空")
}

func TestParseAndFlattenMenu_InvalidYAML(t *testing.T) {
	yamlContent := `
menu:
  - name: "患者数据
    id: "patient_data"
    path: "/patient-data"
`

	version := "v20250928"

	// 测试无效的YAML应该报错
	_, err := ParseAndFlattenMenu(yamlContent, version)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "解析菜单YAML失败")
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
