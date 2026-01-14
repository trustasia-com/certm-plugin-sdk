package certm

import (
	"context"

	"github.com/trustasia-com/certm-plugin-sdk/helper"
)

// BaseComponent 基础组件
type BaseComponent struct{}

// Info 提供组件元数据基础实现，建议重写
func (b *BaseComponent) Info() ComponentInfo {
	return ComponentInfo{
		Type:        ComponentTypeDeploy,
		ID:          "base",
		Name:        "Base Component",
		Description: "Base Component",
		InputTypes:  []DataType{DataTypeAny},
		OutputType:  DataTypeAny,
	}
}

// GetConfigSchema 获取组件输出Schema
func (b *BaseComponent) GetConfigSchema(ctx context.Context) ([]helper.Field, error) {
	return []helper.Field{
		{
			Type:     helper.FieldTypeString,
			Name:     "Name",
			Key:      "name",
			Required: true,
		},
	}, nil
}

// ValidateConfig 验证组件配置是否合法
func (b *BaseComponent) ValidateConfig(ctx context.Context, config helper.FieldConfig) error {
	return nil
}

// Execute 执行组件
func (b *BaseComponent) Execute(ctx context.Context, config helper.FieldConfig, input []*StepOutput) (*StepOutput, error) {
	output, err := NewStepOutput(true, nil, DataTypeDeployResult, "Success")
	if err != nil {
		return nil, err
	}
	return output, nil
}
