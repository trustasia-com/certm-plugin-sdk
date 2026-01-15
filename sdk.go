//go:build tinygo || wasm
// +build tinygo wasm

package certm

import "github.com/trustasia-com/certm-plugin-sdk/helper"

var component Component

// Register 注册组件实现
func Register(c Component) { component = c }

// Component 组件接口，实现的组件必须是无状态的
type Component interface {
	Info() ComponentInfo

	// GetConfigSchema 获取组件配置Schema
	GetConfigSchema(ctx *Context) ([]helper.Field, error)
	// GetDynamicOptions 获取动态选项
	// config: 当前配置，可以携带前置值
	// key: 需要的那个字段的值
	GetDynamicOptions(ctx *Context, config helper.FieldConfig, key string) ([]helper.FieldOption, error)
	// ValidateConfig 验证组件配置是否合法
	ValidateConfig(ctx *Context, config helper.FieldConfig) error

	// ctx: 上下文，包含超时控制
	// config: 组件配置（从 WorkflowStep.ComponentConfig 反序列化）
	// input: 上一步骤的输出（首个步骤为空，DAG模式下可能是多个输入的合并）
	Execute(ctx *Context, config helper.FieldConfig, input []*StepOutput) (*StepOutput, error)
}
