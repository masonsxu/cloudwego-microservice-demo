package permission

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/enum"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl implements the Converter interface.
type ConverterImpl struct {
	enumConverter enum.Converter
}

// NewConverter creates a new ConverterImpl.
func NewConverter(enumConverter enum.Converter) Converter {
	return &ConverterImpl{enumConverter: enumConverter}
}

// ModelToThrift converts a models.Permission to an identity_srv.Permission.
func (c *ConverterImpl) ModelToThrift(
	model *models.Permission,
) *identity_srv.Permission {
	if model == nil {
		return nil
	}

	return &identity_srv.Permission{
		Resource:    &model.Resource,
		Action:      &model.Action,
		Description: &model.Description,
	}
}

// ThriftToModel converts an identity_srv.Permission to a models.Permission.
func (c *ConverterImpl) ThriftToModel(
	thrift *identity_srv.Permission,
) *models.Permission {
	if thrift == nil {
		return nil
	}

	description := ""
	if thrift.Description != nil {
		description = *thrift.Description
	}

	return &models.Permission{
		Resource:    *thrift.Resource,
		Action:      *thrift.Action,
		Description: description,
	}
}

// ModelSliceToThrift converts a slice of models.Permission to a slice of identity_srv.Permission.
func (c *ConverterImpl) ModelSliceToThrift(
	modelSlice []*models.Permission,
) []*identity_srv.Permission {
	if modelSlice == nil {
		return nil
	}

	thriftSlice := make([]*identity_srv.Permission, len(modelSlice))
	for i, m := range modelSlice {
		thriftSlice[i] = c.ModelToThrift(m)
	}

	return thriftSlice
}

// ThriftSliceToModel converts a slice of identity_srv.Permission to a slice of models.Permission.
func (c *ConverterImpl) ThriftSliceToModel(
	thriftSlice []*identity_srv.Permission,
) []*models.Permission {
	if thriftSlice == nil {
		return nil
	}

	modelSlice := make([]*models.Permission, len(thriftSlice))
	for i, t := range thriftSlice {
		modelSlice[i] = c.ThriftToModel(t)
	}

	return modelSlice
}
