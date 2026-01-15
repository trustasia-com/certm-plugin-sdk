# CertM WASM Plugin SDK

基于Component接口的WASM插件开发SDK，为证书管理系统提供类型安全的插件扩展能力。

## 功能特性

- ✅ **类型安全** - 通过Go接口确保编译时类型检查
- ✅ **标准化接口** - 统一的Component接口，支持证书、部署、检测、通知等组件
- ✅ **主机能力** - 日志输出、数据库查询、HTTP请求等主机功能
- ✅ **签名验证** - 基于Ed25519的插件包完整性校验
- ✅ **动态配置** - 支持静态/动态选项、条件显示等高级表单功能

## 安装

### 前置要求

- Go 1.21+
- TinyGo 0.40.1+

### 添加依赖

```bash
go get github.com/trustasia-com/certm-plugin-sdk
```

## 快速开始

### 1. 创建项目

```bash
mkdir my-plugin && cd my-plugin
go mod init github.com/myorg/my-plugin
go get github.com/trustasia-com/certm-plugin-sdk
```

### 2. 实现组件

```go
package main

import (
    "encoding/json"
    "fmt"
    certm "github.com/trustasia-com/certm-plugin-sdk"
    "github.com/trustasia-com/certm-plugin-sdk/helper"
)

type MyDeployer struct{}

func (d *MyDeployer) Info() certm.ComponentInfo {
    return certm.ComponentInfo{
        Type:        certm.ComponentTypeDeploy,
        ID:          "my-deployer",
        Name:        "我的部署器",
        Description: "部署证书到目标服务器",
        InputTypes:  []certm.DataType{certm.DataTypeCertificate},
        OutputType:  certm.DataTypeDeployResult,
    }
}

func (d *MyDeployer) GetConfigSchema(ctx *certm.Context) ([]helper.Field, error) {
    return []helper.Field{
        {
            Type:     helper.FieldTypeString,
            Format:   helper.FieldFormatText,
            Key:      "target_url",
            Name:     "目标地址",
            Required: true,
        },
    }, nil
}

func (d *MyDeployer) ValidateConfig(ctx *certm.Context, config helper.FieldConfig) error {
    if url := config.String("target_url"); url == "" {
        return fmt.Errorf("目标地址不能为空")
    }
    return nil
}

func (d *MyDeployer) GetDynamicOptions(ctx *certm.Context, config helper.FieldConfig, key string) ([]helper.FieldOption, error) {
    return nil, nil
}

func (d *MyDeployer) Execute(ctx *certm.Context, config helper.FieldConfig, input []*certm.StepOutput) (*certm.StepOutput, error) {
    targetURL := config.String("target_url")
    ctx.Info("开始部署到: %s", targetURL)
    
    // 解析证书数据
    var certData certm.CertOutputData
    json.Unmarshal(input[0].Data, &certData)
    
    // 执行部署逻辑...
    
    return certm.NewStepOutput(true, certm.DeployOutputData{
        TargetType: "server",
        TargetName: targetURL,
        Deployed:   true,
        SHA1:       certData.SHA1,
        CommonName: certData.CommonName,
    }, certm.DataTypeDeployResult, "部署成功")
}

func main() {
    certm.Register(&MyDeployer{})
}
```

### 3. 编译

```bash
tinygo build -o plugin.wasm -target=wasi main.go
```

### 4. 优化（可选）

```bash
wasm-opt -Oz plugin.wasm -o plugin.optimized.wasm
```

## API文档

### Component接口

所有插件必须实现以下接口：

```go
type Component interface {
    // 返回组件信息
    Info() ComponentInfo
    
    // 获取配置字段定义
    GetConfigSchema(ctx *Context) ([]Field, error)
    
    // 验证配置
    ValidateConfig(ctx *Context, config FieldConfig) error
    
    // 获取动态选项（用于下拉框等）
    GetDynamicOptions(ctx *Context, config FieldConfig, key string) ([]FieldOption, error)
    
    // 执行组件逻辑
    Execute(ctx *Context, config FieldConfig, input []*StepOutput) (*StepOutput, error)
}
```

### Context 上下文

#### 1. 日志输出

```go
ctx.Info("信息: %s", value)    // 信息日志
ctx.Error("错误: %v", err)     // 错误日志
ctx.Debug("调试信息")          // 调试日志
```

#### 2. 数据查询

```go
// 证书相关
containers, err := ctx.GetCertContainerList()
asset, err := ctx.GetCertAssetDetail(assetID)

// 部署器相关
deployers, err := ctx.GetDeployerList()
deployer, err := ctx.GetDeployerDetail(deployerID)

// 告警规则
rules, err := ctx.GetNoticeRuleList()
```

#### 3. 上下文信息

```go
language := ctx.Language    // 当前语言（如 "zh-CN"）
projectID := ctx.ProjectID  // 当前项目ID
```

### 组件类型

```go
const (
    ComponentTypeCert   ComponentType = "cert"   // 证书组件（证书来源）
    ComponentTypeDeploy ComponentType = "deploy" // 部署组件
    ComponentTypeCheck  ComponentType = "check"  // 检测组件
    ComponentTypeNotice ComponentType = "notice" // 通知组件
)
```

### 数据类型

```go
const (
    DataTypeNone         DataType = "none"          // 无输入（起始组件）
    DataTypeCertificate  DataType = "certificate"   // 证书数据
    DataTypeDeployResult DataType = "deploy_result" // 部署结果
    DataTypeCheckResult  DataType = "check_result"  // 检测结果
    DataTypeNoticeResult DataType = "notice_result" // 通知结果
    DataTypeAny          DataType = "any"           // 任意类型
)
```

### 字段定义

#### 字段类型

```go
const (
    FieldTypeString      FieldType = "string"       // 字符串
    FieldTypeInt         FieldType = "int"          // 整数
    FieldTypeBoolean     FieldType = "boolean"      // 布尔值
    FieldTypeStringArray FieldType = "string_array" // 字符串数组
)
```

#### 字段格式

```go
const (
    FieldFormatText     FieldFormat = "text"     // 单行文本
    FieldFormatTextarea FieldFormat = "textarea" // 多行文本
    FieldFormatPassword FieldFormat = "password" // 密码
    FieldFormatNumber   FieldFormat = "number"   // 数字
    FieldFormatSelect   FieldFormat = "select"   // 下拉选择
    FieldFormatCheckbox FieldFormat = "checkbox" // 复选框
)
```

#### 示例

```go
// 基础字段
{
    Type:     helper.FieldTypeString,
    Format:   helper.FieldFormatText,
    Key:      "api_key",
    Name:     "API密钥",
    Required: true,
}

// 带默认值
{
    Type:    helper.FieldTypeInt,
    Key:     "timeout",
    Name:    "超时时间",
    Default: 60,
}

// 下拉选择（静态）
{
    Type:   helper.FieldTypeString,
    Format: helper.FieldFormatSelect,
    Key:    "region",
    Name:   "地区",
    Options: []helper.FieldOption{
        {Value: "cn", Name: "中国"},
        {Value: "us", Name: "美国"},
    },
}

// 下拉选择（动态）
{
    Type:   helper.FieldTypeString,
    Format: helper.FieldFormatSelect,
    Key:    "server",
    Name:   "服务器",
    OptionsSource: &helper.OptionsSource{
        DependsOn: []string{"region"}, // 依赖region字段
    },
}

// 条件显示
{
    Type: helper.FieldTypeString,
    Key:  "ssl_port",
    Name: "SSL端口",
    ShowWhen: &helper.ShowWhen{
        Key:   "use_ssl",
        Value: true,
    },
}
```

### 工具函数

#### FieldConfig辅助方法

```go
config := helper.FieldConfig{
    "timeout": 30,
    "verify":  true,
    "url":     "https://example.com",
}

timeout := config.Float("timeout")  // 30.0
verify := config.Boolean("verify")  // true
url := config.String("url")         // "https://example.com"
```

#### StepOutput创建

```go
// 成功
return certm.NewStepOutput(true, data, dataType, "操作成功")

// 失败
return certm.NewStepOutput(false, nil, dataType, "操作失败: " + err.Error())
```

## 签名验证

SDK提供基于Ed25519的插件签名验证功能，确保插件包未被篡改。

### 生成密钥对

```go
package main

import (
    "crypto/ed25519"
    "encoding/hex"
    "fmt"
)

func main() {
    publicKey, privateKey, _ := ed25519.GenerateKey(nil)
    fmt.Printf("公钥: %s\n", hex.EncodeToString(publicKey))
    fmt.Printf("私钥: %s\n", hex.EncodeToString(privateKey))
}
```

### 签名插件包

```go
package main

import (
    "crypto/ed25519"
    "encoding/hex"
    "encoding/json"
    "os"
    certm "github.com/trustasia-com/certm-plugin-sdk"
)

func main() {
    // 1. 创建清单
    manifest, _ := certm.CreateManifest([]string{
        "plugin.wasm",
        "plugin.yml",
    }, ".")
    
    // 2. 签名
    privateKeyHex := "YOUR_PRIVATE_KEY"
    privateKey, _ := hex.DecodeString(privateKeyHex)
    
    manifestJSON, _ := json.Marshal(manifest)
    signature := certm.Sign(manifestJSON, privateKey)
    
    // 3. 保存文件
    os.WriteFile("manifest.json", manifestJSON, 0644)
    os.WriteFile("signature", signature, 0644)
}
```

### 验证插件包

```go
package main

import (
    "crypto/ed25519"
    "encoding/hex"
    certm "github.com/trustasia-com/certm-plugin-sdk"
)

func main() {
    publicKeyHex := "DEVELOPER_PUBLIC_KEY"
    publicKey, _ := hex.DecodeString(publicKeyHex)
    
    err := certm.VerifyZip("plugin.zip", publicKey)
    if err != nil {
        panic("签名验证失败: " + err.Error())
    }
}
```

## 项目结构

```
certm-plugin-sdk/
├── auth.go           # 签名验证
├── context.go        # Context实现
├── export.go         # WASM导出函数
├── host.go           # 主机函数声明
├── memory.go         # 内存管理
├── sdk.go            # SDK核心
├── types.go          # 类型定义
├── helper/           # 辅助工具
│   ├── field.go      # 字段定义
│   └── config.go     # 配置解析
└── example/          # 示例插件
    ├── main.go
    ├── plugin.yml
    └── Makefile
```

## 完整示例

查看 [`example/`](./example) 目录获取完整示例，包括：

- 完整的部署器实现
- 动态字段选项
- 配置验证
- Makefile构建脚本

## 常见问题

### Q: 如何调试WASM插件？

A: 使用日志输出：

```go
ctx.Debug("调试信息: %+v", data)
ctx.Info("当前状态: %s", state)
ctx.Error("错误详情: %v", err)
```

### Q: FieldConfig如何处理类型转换？

A: 使用辅助方法安全转换：

```go
timeout := config.Float("timeout")     // 自动转换为float64
verify := config.Boolean("verify")     // 自动转换为bool
url := config.String("url")            // 自动转换为string
tags := config.StringArray("tags")     // 自动转换为[]string
```

### Q: 如何处理动态字段依赖？

A: 在`GetDynamicOptions`中根据配置返回选项：

```go
func (d *MyDeployer) GetDynamicOptions(ctx *certm.Context, config helper.FieldConfig, key string) ([]helper.FieldOption, error) {
    if key == "instance" {
        region := config.String("region")
        // 根据region查询实例列表
        instances := queryInstancesByRegion(ctx, region)
        return instances, nil
    }
    return nil, nil
}
```

### Q: Execute方法可以执行HTTP请求吗？

A: TinyGo的网络支持有限，建议通过Context提供的主机函数实现，或在主机端添加自定义主机函数。

### Q: 如何发布插件？

A: 
1. 编译WASM文件
2. 创建`plugin.yml`描述文件
3. 生成签名（manifest.json + signature）
4. 打包为ZIP文件
5. 通过平台上传

## 设计原则

- ✅ **类型安全** - 通过Go接口和类型系统确保编译时检查
- ✅ **统一返回** - 所有导出函数统一返回Result格式，便于错误处理
- ✅ **Context传递** - 主机信息通过Context传递，避免全局变量
- ✅ **无状态组件** - 组件必须无状态，所有状态通过参数传递
- ✅ **错误透明** - 主机能够感知和处理所有错误情况

## 许可证

Apache 2.0
