package main

import (
	certm "github.com/trustasia-com/certm-plugin-sdk"
)

// 插件元数据（必须导出）
var PluginMeta = certm.PluginMetadata{
	Name:        "example-deployer",
	Version:     "1.0.0",
	Type:        certm.ComponentTypeDeploy,
	Author:      "CertM Team",
	Description: "示例部署器插件",
}

// NewComponent 创建组件实例（必须导出）
func NewComponent() certm.Component {
	return &ExampleDeployer{}
}

// OnLoad 插件加载时的回调（可选）
func OnLoad() error {
	// logger := certm.GetLogger()
	// logger.Info("Example deployer plugin loaded",
	// 	certm.F("version", PluginMeta.Version),
	// )
	return nil
}

// OnUnload 插件卸载时的回调（可选）
func OnUnload() error {
	// logger := certm.GetLogger()
	// logger.Info("Example deployer plugin unloading")
	return nil
}
