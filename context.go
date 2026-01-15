//go:build tinygo || wasm
// +build tinygo wasm

package certm

import "fmt"

// Context 上下文
type Context struct {
	Language  string `json:"language"`
	ProjectID int    `json:"project_id"`
}

// GetCertContainerList 获取证书容器列表
func (c *Context) GetCertContainerList() ([]*CertContainerInfo, error) {
	list, err := call[[]*CertContainerInfo]("db_get_cert_container_list", c.ProjectID)
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// GetCertAssetListOfContainer 获取证书资产列表
func (c *Context) GetCertAssetListOfContainer(containerID int) ([]*CertAssetInfo, error) {
	list, err := call[[]*CertAssetInfo]("db_get_cert_asset_list_of_container", containerID)
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// GetCertAssetDetail 获取证书资产详情
func (c *Context) GetCertAssetDetail(assetID int) (*CertAssetDetail, error) {
	asset, err := call[*CertAssetDetail]("db_get_cert_asset_detail", assetID)
	if err != nil {
		return nil, err
	}
	return *asset, nil
}

// GetDeployerList 获取部署器列表
func (c *Context) GetDeployerList() ([]*DeployerInfo, error) {
	list, err := call[[]*DeployerInfo]("db_get_deployer_list", c.ProjectID)
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// GetDeployerDetail 获取部署器详情
func (c *Context) GetDeployerDetail(deployerID int) (*DeployerDetail, error) {
	deployer, err := call[*DeployerDetail]("db_get_deployer_detail", deployerID)
	if err != nil {
		return nil, err
	}
	return *deployer, nil
}

// GetNoticeRuleList 获取告警规则列表
func (c *Context) GetNoticeRuleList() ([]*NoticeRuleInfo, error) {
	list, err := call[[]*NoticeRuleInfo]("db_get_notice_rule_list", c.ProjectID)
	if err != nil {
		return nil, err
	}
	return *list, nil
}

// sprintf 简化的格式化字符串（兼容TinyGo）
func sprintf(format string, args ...interface{}) string {
	// 简化版本：如果有参数就用fmt.Sprintf，否则直接返回
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

// Debug 输出调试日志
func (c *Context) Debug(format string, args ...interface{}) {
	hostLogImpl("debug", sprintf(format, args...))
}

// Error 输出错误日志
func (c *Context) Error(format string, args ...interface{}) {
	hostLogImpl("error", sprintf(format, args...))
}

// Info 输出信息日志
func (c *Context) Info(format string, args ...interface{}) {
	hostLogImpl("info", sprintf(format, args...))
}
