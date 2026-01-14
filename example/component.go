package main

import (
	"context"

	certm "github.com/trustasia-com/certm-plugin-sdk"
	"github.com/trustasia-com/certm-plugin-sdk/helper"
)

// ExampleDeployer 示例部署器组件
type ExampleDeployer struct {
	certm.BaseComponent // 嵌入基础实现
}

// Info 返回组件信息
func (d *ExampleDeployer) Info() certm.ComponentInfo {
	return certm.ComponentInfo{
		Type:        PluginMeta.Type,
		ID:          PluginMeta.ID,
		Name:        PluginMeta.Name,
		Description: PluginMeta.Description,
	}
}

// GetConfigSchema 获取组件配置Schema
func (d *ExampleDeployer) GetConfigSchema(ctx context.Context) ([]helper.Field, error) {
	return []helper.Field{
		{
			Type:     helper.FieldTypeString,
			Name:     "部署目标",
			Key:      "target",
			Required: true,
		},
	}, nil
}

// GetDynamicOptions 获取动态选项
// config: 当前配置，可以携带前置值
// key: 需要的那个字段的值
func (d *ExampleDeployer) GetDynamicOptions(ctx context.Context, config helper.FieldConfig, key string) ([]helper.FieldOption, error) {
	return nil, nil
}

// Execute 执行部署逻辑
func (d *ExampleDeployer) Execute(ctx context.Context, config helper.FieldConfig, input []*certm.StepOutput) (*certm.StepOutput, error) {

	return nil, nil
}

// ValidateConfig 验证配置
func (d *ExampleDeployer) ValidateConfig(ctx context.Context, config helper.FieldConfig) error {

	return nil
}
