package certm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/trustasia-com/certm-plugin-sdk/helper"
)

// Component 组件接口
type Component interface {
	Info() ComponentInfo

	// GetConfigSchema 获取组件配置Schema
	GetConfigSchema(ctx context.Context) ([]helper.Field, error)
	// GetDynamicOptions 获取动态选项
	// config: 当前配置，可以携带前置值
	// key: 需要的那个字段的值
	GetDynamicOptions(ctx context.Context, config helper.FieldConfig, key string) ([]helper.FieldOption, error)
	// ValidateConfig 验证组件配置是否合法
	ValidateConfig(ctx context.Context, config helper.FieldConfig) error

	// ctx: 上下文，包含超时控制
	// config: 组件配置（从 WorkflowStep.ComponentConfig 反序列化）
	// input: 上一步骤的输出（首个步骤为空，DAG模式下可能是多个输入的合并）
	Execute(ctx context.Context, config helper.FieldConfig, input []*StepOutput) (*StepOutput, error)
}

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

// PluginMetadata 插件元数据
type PluginMetadata struct {
	ID          string        `json:"id"`          // 插件ID
	Name        string        `json:"name"`        // 插件名称
	Version     string        `json:"version"`     // 插件版本
	Type        ComponentType `json:"type"`        // 插件类型
	Author      string        `json:"author"`      // 插件作者
	Description string        `json:"description"` // 插件描述
}

// StepOutput 步骤输出
type StepOutput struct {
	Success  bool            `json:"success"`
	Data     json.RawMessage `json:"data,omitempty"`
	DataType DataType        `json:"data_type,omitempty"` // 数据类型标识
	Message  string          `json:"message,omitempty"`
}

// CertOutputData 证书组件输出数据
type CertOutputData struct {
	CommonName  string    `json:"common_name"`  // 通用名称
	KeyPEM      string    `json:"key_pem"`      // 私钥PEM
	ChainPEM    []string  `json:"chain_pem"`    // 证书链PEM，包含叶子证书
	NotAfter    time.Time `json:"not_after"`    // 过期时间
	SHA1        string    `json:"sha1"`         // 证书SHA1
	HistorySHA1 []string  `json:"history_sha1"` // 历史证书SHA1

	EncKeyPEM  string `json:"enc_key_pem"`  // 加密私钥PEM
	EncCertPEM string `json:"enc_cert_pem"` // 加密证书PEM
}

// DeployOutputData 部署组件输出数据
type DeployOutputData struct {
	// 部署目标信息
	TargetType string `json:"target_type"` // 目标类型：CDN/LB...
	TargetName string `json:"target_name"` // 目标名称: 域名/IP/节点ID...
	// 部署结果
	Deployed bool   `json:"deployed"` // 是否部署成功: true/false
	Error    string `json:"error"`    // 错误信息: 部署失败原因
	// 证书信息用于后续check判断
	SHA1       string    `json:"sha1"`        // 证书SHA1
	CommonName string    `json:"common_name"` // 通用名称
	NotAfter   time.Time `json:"not_after"`   // 过期时间
	DeployedAt time.Time `json:"deployed_at"` // 部署时间
}

// CheckEndpointResult 检测端点结果
type CheckEndpointResult struct {
	Endpoint string `json:"endpoint"` // 检测地址: example.com:443
	IP       string `json:"ip"`       // 检测IP
	// 检测结果
	TrustStatus  TrustStatus `json:"trust_status"`  // SSL信任状态
	CertMatch    bool        `json:"cert_match"`    // 证书是否匹配
	ResponseTime int         `json:"response_time"` // 响应时间(ms)
	Error        string      `json:"error"`         // 错误信息: 检测失败原因
	// 证书信息
	SHA1       string    `json:"sha1"`        // 证书SHA1
	CommonName string    `json:"common_name"` // 通用名称
	NotAfter   time.Time `json:"not_after"`   // 过期时间
	CheckedAt  time.Time `json:"checked_at"`  // 检测时间
}

// CheckOutputData 检测组件输出数据
type CheckOutputData struct {
	Endpoints []*CheckEndpointResult `json:"endpoints"` // 检测端点结果
}

// NotifyOutputData 通知组件输出数据
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

	return &StepOutput{
		Success:  success,
		Data:     dataBytes,
		DataType: dataType,
		Message:  message,
	}, nil
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
