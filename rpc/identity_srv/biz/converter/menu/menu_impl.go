package menu

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl implements the Converter interface.
type ConverterImpl struct{}

// NewConverter creates a new Converter.
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ModelToThrift converts a menu model to its thrift representation, including its children recursively.
// 使用 SemanticID 而非 UUID 作为对外暴露的菜单标识符，确保在版本间的稳定性。
func (c *ConverterImpl) ModelToThrift(model *models.Menu) *identity_srv.MenuNode {
	if model == nil {
		return nil
	}

	return &identity_srv.MenuNode{
		Id:        &model.SemanticID, // 使用语义化ID而非UUID
		Name:      &model.Name,
		Path:      &model.Path,
		Icon:      &model.Icon,
		Component: &model.Component,
		Children:  c.ModelsToThrift(model.Children), // Recursively convert children
	}
}

// ModelsToThrift converts a slice of menu models to a slice of thrift representations.
func (c *ConverterImpl) ModelsToThrift(models []*models.Menu) []*identity_srv.MenuNode {
	if len(models) == 0 {
		return nil
	}

	thrifts := make([]*identity_srv.MenuNode, 0, len(models))
	for _, model := range models {
		thrifts = append(thrifts, c.ModelToThrift(model))
	}

	return thrifts
}
