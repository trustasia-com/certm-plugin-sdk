package certm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/trustasia-com/certm-plugin-sdk/helper"
)

// Context 上下文
type Context struct {
	Language  string `json:"language"`
	ProjectID int    `json:"project_id"`
}

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

// ComponentType 组件类型
type ComponentType string

// String 实现 Stringer 接口
func (c ComponentType) String() string {
	return string(c)
}

const (
	ComponentTypeCert   ComponentType = "cert"   // 证书组件（证书来源）
	ComponentTypeDeploy ComponentType = "deploy" // 部署组件
	ComponentTypeCheck  ComponentType = "check"  // 检测组件
	ComponentTypeNotice ComponentType = "notice" // 通知组件
)

// DataType 数据类型
type DataType string

// String 实现 Stringer 接口
func (d DataType) String() string {
	return string(d)
}

const (
	DataTypeNone         DataType = "none"          // 无输入（起始组件）
	DataTypeCertificate  DataType = "certificate"   // 证书数据
	DataTypeDeployResult DataType = "deploy_result" // 部署结果
	DataTypeCheckResult  DataType = "check_result"  // 检测结果
	DataTypeNoticeResult DataType = "notice_result" // 通知结果
	DataTypeAny          DataType = "any"           // 任意类型（通用接收器）
)

// TrustStatus 信任状态
type TrustStatus int

// Int 实现 Inter 接口
func (t TrustStatus) Int() int {
	return int(t)
}

const (
	TrustStatusUnspecified    TrustStatus = 0  // 未指定
	TrustStatusTrusted        TrustStatus = 1  // 信任
	TrustStatusCertExpired    TrustStatus = 2  // 证书过期
	TrustStatusCertRevoked    TrustStatus = 4  // 证书被吊销
	TrustStatusCADisabled     TrustStatus = 5  // CA被禁用
	TrustStatusCAExpired      TrustStatus = 6  // CA过期
	TrustStatusCARevoked      TrustStatus = 7  // CA被吊销
	TrustStatusCANotFind      TrustStatus = 8  // 无法组链
	TrustStatusCARemoved      TrustStatus = 9  // CA被移除
	TrustStatusChainErr       TrustStatus = 10 // 证书链配置错误
	TrustStatusSelfSigned     TrustStatus = 11 // 自签名证书
	TrustStatusDomainNotMatch TrustStatus = 12 // 域名不匹配
)

// ComponentInfo 组件信息
type ComponentInfo struct {
	Type ComponentType `json:"type"`

	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// 组件输入类型
	InputTypes []DataType `json:"input_types"`
	// 组件输出类型
	OutputType DataType `json:"output_type"`
}

// StepOutput 步骤输出
type StepOutput struct {
	Success  bool            `json:"success"`           // 是否成功
	Data     json.RawMessage `json:"data,omitempty"`    // 输出数据
	DataType DataType        `json:"data_type"`         // 数据类型
	Message  string          `json:"message,omitempty"` // 消息
}

// CertOutputData 证书数据
type CertOutputData struct {
	SHA1        string    `json:"sha1"`         // 证书指纹
	CommonName  string    `json:"common_name"`  // 主域名
	KeyPEM      string    `json:"key_pem"`      // 私钥PEM
	ChainPEM    []string  `json:"chain_pem"`    // 证书链PEM
	NotAfter    time.Time `json:"not_after"`    // 过期时间
	HistorySHA1 []string  `json:"history_sha1"` // 历史证书SHA1

	EncKeyPEM  string `json:"enc_key_pem,omitempty"`  // 加密私钥PEM
	EncCertPEM string `json:"enc_cert_pem,omitempty"` // 加密证书PEM
}

// DeployOutputData 部署结果
type DeployOutputData struct {
	// 部署目标
	TargetType string `json:"target_type"` // 目标类型: CDN/LB
	TargetName string `json:"target_name"` // 目标名称: 域名/IP/节点ID

	// 部署结果
	Deployed bool   `json:"deployed"` // 是否部署成功
	Error    string `json:"error"`    // 错误信息

	// 证书信息（用于后续check）
	SHA1       string    `json:"sha1"`        // 证书指纹
	CommonName string    `json:"common_name"` // 主域名
	NotAfter   time.Time `json:"not_after"`   // 过期时间
	DeployedAt time.Time `json:"deployed_at"` // 部署时间
}

// CheckEndpointResult 单个检测端点结果
type CheckEndpointResult struct {
	Endpoint string `json:"endpoint"` // 检测地址: example.com:443
	IP       string `json:"ip"`       // 检测IP

	// 检测结果
	TrustStatus  TrustStatus `json:"trust_status"`  // SSL信任状态
	CertMatch    bool        `json:"cert_match"`    // 证书是否匹配
	ResponseTime int         `json:"response_time"` // 响应时间(ms)
	Error        string      `json:"error"`         // 错误信息

	// 证书信息
	SHA1       string    `json:"sha1"`        // 证书指纹
	CommonName string    `json:"common_name"` // 主域名
	NotAfter   time.Time `json:"not_after"`   // 过期时间
	CheckedAt  time.Time `json:"checked_at"`  // 检测时间
}

// CheckOutputData 检测结果
type CheckOutputData struct {
	Endpoints []*CheckEndpointResult `json:"endpoints"` // 检测端点列表
}

// NotifyOutputData 通知结果
type NotifyOutputData struct {
	Channel string `json:"channel"` // 通知渠道
	Sent    bool   `json:"sent"`    // 是否发送成功
}

// NewStepOutput 创建新的步骤输出
func NewStepOutput(success bool, data any, dataType DataType, message string) (*StepOutput, error) {
	var dataBytes json.RawMessage
	if data != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("marshal data: %w", err)
		}
		dataBytes = bytes
	}

	return &StepOutput{Success: success, Data: dataBytes, DataType: dataType, Message: message}, nil
}

// ParseData 解析数据
func (s *StepOutput) ParseData(v any) error {
	if s == nil || len(s.Data) == 0 {
		return fmt.Errorf("no data to parse")
	}
	return json.Unmarshal(s.Data, v)
}

// GetDataType 获取数据类型
func (s *StepOutput) GetDataType() DataType {
	if s == nil {
		return DataTypeNone
	}
	return s.DataType
}

// IsCertificate 检查输出是否包含证书数据
func (s *StepOutput) IsCertificate() bool {
	return s != nil && s.DataType == DataTypeCertificate
}

// IsDeployResult 检查输出是否包含部署结果
func (s *StepOutput) IsDeployResult() bool {
	return s != nil && s.DataType == DataTypeDeployResult
}

// IsCheckResult 检查输出是否包含检测结果
func (s *StepOutput) IsCheckResult() bool {
	return s != nil && s.DataType == DataTypeCheckResult
}

// ParseCertificate 解析证书数据
func (s *StepOutput) ParseCertificate() (*CertOutputData, error) {
	if !s.IsCertificate() {
		return nil, fmt.Errorf("output is not certificate data, got: %s", s.DataType)
	}
	var data CertOutputData
	if err := s.ParseData(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// ParseDeployResult 解析部署结果
func (s *StepOutput) ParseDeployResult() (*DeployOutputData, error) {
	if !s.IsDeployResult() {
		return nil, fmt.Errorf("output is not deploy result, got: %s", s.DataType)
	}
	var data DeployOutputData
	if err := s.ParseData(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// ParseCheckResult 解析检测结果
func (s *StepOutput) ParseCheckResult() (*CheckOutputData, error) {
	if !s.IsCheckResult() {
		return nil, fmt.Errorf("output is not check result, got: %s", s.DataType)
	}
	var data CheckOutputData
	if err := s.ParseData(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

/////////////////////  Host  /////////////////////

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

// Result 执行结果
type Result struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

/////////////////////  Plugin  /////////////////////

// PluginYaml 插件信息
type PluginYaml struct {
	ID string `json:"id" yaml:"id"`

	Name        string        `json:"name" yaml:"name"`
	Version     string        `json:"version" yaml:"version"`
	Description string        `json:"description" yaml:"description"`
	Type        ComponentType `json:"type" yaml:"type"`

	Author        *AuthorInfo        `json:"author" yaml:"author"`
	Tags          []string           `json:"tags" yaml:"tags"`
	Compatibility *CompatibilityInfo `json:"compatibility" yaml:"compatibility"`
}

// AuthorInfo 作者信息
type AuthorInfo struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"email" yaml:"email"`
	URL   string `json:"url" yaml:"url"`
}

// CompatibilityInfo 兼容性信息
type CompatibilityInfo struct {
	MinVersion string `json:"min_version" yaml:"min_version"`
	MaxVersion string `json:"max_version" yaml:"max_version"`
}
