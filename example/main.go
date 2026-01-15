//go:build tinygo || wasm
// +build tinygo wasm

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	certm "github.com/trustasia-com/certm-plugin-sdk"
	"github.com/trustasia-com/certm-plugin-sdk/helper"
)

// MyDeployer 实现Component接口
type MyDeployer struct{}

// Info 返回组件信息
func (d *MyDeployer) Info() certm.ComponentInfo {
	return certm.ComponentInfo{
		Type:        certm.ComponentTypeDeploy,
		ID:          "example-deployer",
		Name:        "示例部署器",
		Description: "这是一个示例WASM插件",
		InputTypes:  []certm.DataType{certm.DataTypeCertificate},
		OutputType:  certm.DataTypeDeployResult,
	}
}

// GetConfigSchema 返回配置字段
func (d *MyDeployer) GetConfigSchema(ctx context.Context) ([]helper.Field, error) {
	return []helper.Field{
		{
			Type:     helper.FieldTypeString,
			Format:   helper.FieldFormatText,
			Key:      "target_url",
			Name:     "目标地址",
			Required: true,
		},
		{
			Type:    helper.FieldTypeInt,
			Format:  helper.FieldFormatNumber,
			Key:     "timeout",
			Name:    "超时时间（秒）",
			Default: 60,
		},
		{
			Type:    helper.FieldTypeBoolean,
			Format:  helper.FieldFormatCheckbox,
			Key:     "verify_ssl",
			Name:    "验证SSL证书",
			Default: true,
		},
		{
			Type:   helper.FieldTypeString,
			Format: helper.FieldFormatSelect,
			Key:    "region",
			Name:   "地区",
			OptionsSource: &helper.OptionsSource{
				DependsOn: []string{},
			},
		},
	}, nil
}

// ValidateConfig 验证配置
func (d *MyDeployer) ValidateConfig(ctx context.Context, config helper.FieldConfig) error {
	// 示例：验证URL格式
	if url, ok := config["target_url"].(string); ok {
		if url == "" {
			return fmt.Errorf("字段 'target_url': 目标地址不能为空")
		}
	}
	return nil // 验证通过
}

// GetDynamicOptions 获取动态选项
func (d *MyDeployer) GetDynamicOptions(ctx context.Context, config helper.FieldConfig, key string) ([]helper.FieldOption, error) {
	// 示例：根据不同字段返回不同选项
	switch key {
	case "region":
		return []helper.FieldOption{
			{Value: "cn-beijing", Name: "北京"},
			{Value: "cn-shanghai", Name: "上海"},
			{Value: "us-west", Name: "美国西部"},
		}, nil
	default:
		return nil, nil
	}
}

// Execute 执行部署逻辑
func (d *MyDeployer) Execute(ctx context.Context, config helper.FieldConfig, input []*certm.StepOutput) (*certm.StepOutput, error) {
	// 1. 获取配置
	targetURL := config.String("target_url")
	timeout := config.Float("timeout")
	verifySSL := config.Boolean("verify_ssl")

	// 2. 获取DataAccess用于数据查询（如果需要）
	// dataAccess := certm.GetDataAccess(ctx)
	// projectID := certm.GetProjectID(ctx)
	// lang := certm.GetLang(ctx)

	fmt.Printf("开始部署: target=%s, timeout=%.0f, verifySSL=%v\n", targetURL, timeout, verifySSL)

	// 3. 获取输入证书数据
	if len(input) == 0 || input[0] == nil {
		return certm.NewStepOutput(false, nil, certm.DataTypeCertificate, "未找到证书数据")
	}

	// 解析证书数据
	var certData certm.CertOutputData
	if err := json.Unmarshal(input[0].Data, &certData); err != nil {
		return certm.NewStepOutput(false, nil, certm.DataTypeCertificate, fmt.Sprintf("解析证书数据失败: %v", err))
	}

	fmt.Printf("证书: CN=%s, SHA1=%s\n", certData.CommonName, certData.SHA1)

	// 4. 执行部署（示例：可以调用DataAccess提供的方法或HTTP）
	// 这里演示简单的成功返回

	// 5. 返回成功结果
	deployResult := certm.DeployOutputData{
		TargetType: "cdn",
		TargetName: targetURL,
		Deployed:   true,
		Error:      "",
		SHA1:       certData.SHA1,
		CommonName: certData.CommonName,
		NotAfter:   certData.NotAfter,
		DeployedAt: time.Now(),
	}

	return certm.NewStepOutput(true, deployResult, certm.DataTypeDeployResult, fmt.Sprintf("已部署到 %s", targetURL))
}

func main() {
	// 注册组件实现
	certm.Register(&MyDeployer{})
}
