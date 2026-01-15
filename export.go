//go:build tinygo || wasm
// +build tinygo wasm

package certm

import (
	"encoding/json"

	"github.com/trustasia-com/certm-plugin-sdk/helper"
)

// component_info 获取组件信息
//
//export component_info
func componentInfo() (ptr uint32) {
	result := checkComponentRegistered()
	defer result.writeToMemory(&ptr)

	if !result.Success {
		return
	}

	result.Data = component.Info()
	return
}

// get_config_schema 获取配置Schema
//
//export get_config_schema
func getConfigSchema(ctxPtr uint32) (ptr uint32) {
	result := checkComponentRegistered()
	defer result.writeToMemory(&ptr)

	if !result.Success {
		return
	}

	// 1. 读取并解析Context
	ctx := parseContext(ctxPtr)

	// 2. 调用组件方法获取Schema
	fields, err := component.GetConfigSchema(ctx)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	result.Data = fields
	return
}

// get_dynamic_options 获取动态选项
//
//export get_dynamic_options
func getDynamicOptions(ctxPtr, configPtr, keyPtr uint32) (ptr uint32) {
	result := checkComponentRegistered()
	defer result.writeToMemory(&ptr)

	if !result.Success {
		return
	}

	// 1. 解析参数
	ctx := parseContext(ctxPtr)
	configData := readFromMemory(configPtr)
	keyData := readFromMemory(keyPtr)

	var config helper.FieldConfig
	json.Unmarshal(configData, &config)
	key := string(keyData)

	// 2. 调用组件方法
	options, err := component.GetDynamicOptions(ctx, config, key)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	result.Data = options
	return
}

// validate_config 验证配置
//
//export validate_config
func validateConfig(ctxPtr, configPtr uint32) (ptr uint32) {
	result := checkComponentRegistered()
	defer result.writeToMemory(&ptr)

	if !result.Success {
		return
	}

	// 1. 解析参数
	ctx := parseContext(ctxPtr)
	configData := readFromMemory(configPtr)

	var config helper.FieldConfig
	json.Unmarshal(configData, &config)

	// 2. 调用验证方法
	err := component.ValidateConfig(ctx, config)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	return
}

// execute 执行组件
//
//export execute
func execute(ctxPtr, configPtr, inputPtr uint32) (ptr uint32) {
	result := checkComponentRegistered()
	defer result.writeToMemory(&ptr)

	if !result.Success {
		return
	}

	// 1. 解析参数
	ctx := parseContext(ctxPtr)
	configData := readFromMemory(configPtr)
	inputData := readFromMemory(inputPtr)

	var config helper.FieldConfig
	var input []*StepOutput

	err := json.Unmarshal(configData, &config)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	err = json.Unmarshal(inputData, &input)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}

	// 2. 执行
	output, err := component.Execute(ctx, config, input)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	result.Data = output
	return
}

// parseContext 从内存指针解析Context
func parseContext(ptr uint32) *Context {
	ctx := &Context{}

	if ptr == 0 {
		return ctx
	}

	// 从内存读取Context数据
	data := readFromMemory(ptr)

	// Context数据格式: {"language": "zh-CN", "project_id": 123}
	if err := json.Unmarshal(data, ctx); err != nil {
		return ctx
	}

	return ctx
}

// checkComponentRegistered 检查组件是否已注册
func checkComponentRegistered() *Result {
	if component == nil {
		return &Result{Success: false, Error: "component not registered"}
	}
	return &Result{Success: true}
}
