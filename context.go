package certm

import (
	"context"
	"encoding/json"
	"time"
)

var (
	DataAccessCtxKey = "dataAccess"
	LangCtxKey       = "lang"
	ProjectCtxKey    = "project"
)

// DataAccess 组件访问数据接口
type DataAccess interface {
	// 证书组件访问接口
	SelectCertContainerList(projectID int) ([]*CertContainerInfo, error)
	SelectCertAssetListOfContainer(projectID, containerID int) ([]*CertAssetInfo, error)
	SelectCertAssetDetail(projectID, assetID int) (*CertAssetDetail, error)

	// 部署组件访问接口
	SelectDeployerList(projectID int, targetID string) ([]*DeployerInfo, error)
	SelectDeployerDetail(projectID, deployerID int) (*DeployerDetail, error)

	//
	// 检测组件访问接口

	// 通知组件访问接口
	SelectNoticeRuleList() ([]*NoticeRuleInfo, error)
}

// GetDataAccess 获取组件访问数据
// nolint:errcheck
func GetDataAccess(ctx context.Context) DataAccess {
	return ctx.Value(DataAccessCtxKey).(DataAccess)
}

// GetLang 获取语言
// nolint:errcheck
func GetLang(ctx context.Context) string {
	return ctx.Value(LangCtxKey).(string)
}

// GetProjectID 获取项目ID
// nolint:errcheck
func GetProjectID(ctx context.Context) int {
	return ctx.Value(ProjectCtxKey).(int)
}

// CertContainerInfo 证书容器信息
type CertContainerInfo struct {
	ID int `json:"id"`

	Status     string `json:"status"`
	CommonName string `json:"common_name"`
	KeyAlgo    string `json:"key_algo"`
	ExistKey   bool   `json:"exist_key"`
}

// CertAssetInfo 证书资产信息
type CertAssetInfo struct {
	ID int `json:"id"`

	SHA1       string    `json:"sha1"`        // 证书SHA1
	CommonName string    `json:"common_name"` // 通用名称
	NotAfter   time.Time `json:"not_after"`   // 过期时间
}

// CertAssetDetail 证书资产详情
type CertAssetDetail struct {
	CertAssetInfo

	KeyPEM   string   `json:"key_pem"`   // 私钥PEM
	ChainPEM []string `json:"chain_pem"` // 证书链PEM，包含叶子证书
}

// NoticeRuleInfo 告警规则信息
type DeployerInfo struct {
	ID int `json:"id"`

	Name   string `json:"name"`
	Status string `json:"status"`
	Remark string `json:"remark"`
}

// DeployerDetail 部署器详情
type DeployerDetail struct {
	DeployerInfo

	Credentials json.RawMessage `json:"credentials"`
	Config      json.RawMessage `json:"config"`
}

// WorkflowStepInfo 工作流步骤信息
type WorkflowStepInfo struct {
	ID int `json:"id"`

	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

// NoticeRuleInfo 告警规则信息
type NoticeRuleInfo struct {
	ID int `json:"id"`

	Name string `json:"name"`
}
