//go:build tinygo || wasm
// +build tinygo wasm

package certm

import (
	"context"
	"encoding/json"
	"fmt"

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
	certmCtx := parseCertmContext(ctxPtr)
	ctx := SetContextKey(context.Background(), certmCtx, certmCtx.Language, certmCtx.ProjectID)

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

	// 1. 读取并解析Context
	certmCtx := parseCertmContext(ctxPtr)
	ctx := SetContextKey(context.Background(), certmCtx, certmCtx.Language, certmCtx.ProjectID)

	// 2. 解析参数
	configData := readFromMemory(configPtr)
	keyData := readFromMemory(keyPtr)

	var config helper.FieldConfig
	err := json.Unmarshal(configData, &config)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	key := string(keyData)

	// 3. 调用组件方法
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

	// 1. 读取并解析Context
	certmCtx := parseCertmContext(ctxPtr)
	ctx := SetContextKey(context.Background(), certmCtx, certmCtx.Language, certmCtx.ProjectID)

	// 2. 解析参数
	configData := readFromMemory(configPtr)

	var config helper.FieldConfig
	err := json.Unmarshal(configData, &config)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}

	// 3. 调用验证方法
	err = component.ValidateConfig(ctx, config)
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

	// 1. 读取并解析Context
	certmCtx := parseCertmContext(ctxPtr)
	ctx := SetContextKey(context.Background(), certmCtx, certmCtx.Language, certmCtx.ProjectID)

	// 2. 解析参数
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

	// 3. 执行
	output, err := component.Execute(ctx, config, input)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return
	}
	result.Data = output
	return
}

// GetCertContainerList 获取证书容器列表
func (c *CertmContext) GetCertContainerList(projectID int) ([]*CertContainerInfo, error) {
	list, err := call[[]*CertContainerInfo]("db_get_cert_container_list", projectID)
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// GetCertAssetListOfContainer 获取证书资产列表
func (c *CertmContext) GetCertAssetListOfContainer(projectID, containerID int) ([]*CertAssetInfo, error) {
	list, err := call[[]*CertAssetInfo]("db_get_cert_asset_list_of_container", projectID, containerID)
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// GetCertAssetDetail 获取证书资产详情
func (c *CertmContext) GetCertAssetDetail(projectID, assetID int) (*CertAssetDetail, error) {
	asset, err := call[*CertAssetDetail]("db_get_cert_asset_detail", projectID, assetID)
	if err != nil {
		return nil, err
	}
	return *asset, nil
}

// GetDeployerList 获取部署器列表
func (c *CertmContext) GetDeployerList(projectID int, targetID string) ([]*DeployerInfo, error) {
	list, err := call[[]*DeployerInfo]("db_get_deployer_list", projectID, targetID)
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// GetDeployerDetail 获取部署器详情
func (c *CertmContext) GetDeployerDetail(projectID, deployerID int) (*DeployerDetail, error) {
	deployer, err := call[*DeployerDetail]("db_get_deployer_detail", projectID, deployerID)
	if err != nil {
		return nil, err
	}
	return *deployer, nil
}

// GetNoticeRuleList 获取告警规则列表
func (c *CertmContext) GetNoticeRuleList() ([]*NoticeRuleInfo, error) {
	list, err := call[[]*NoticeRuleInfo]("db_get_notice_rule_list")
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// sprintf 简化的格式化字符串（兼容TinyGo）
func sprintf(format string, args ...any) string {
	// 简化版本：如果有参数就用fmt.Sprintf，否则直接返回
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

// Debug 输出调试日志
func (c *CertmContext) Debug(format string, args ...any) {
	hostLogImpl("debug", sprintf(format, args...))
}

// Error 输出错误日志
func (c *CertmContext) Error(format string, args ...any) {
	hostLogImpl("error", sprintf(format, args...))
}

// Info 输出信息日志
func (c *CertmContext) Info(format string, args ...any) {
	hostLogImpl("info", sprintf(format, args...))
}

// parseCertmContext 从内存指针解析Context
func parseCertmContext(ptr uint32) *CertmContext {
	ctx := &CertmContext{}

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
