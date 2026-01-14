# CertM Plugin SDK

CertM 工作流组件插件开发工具包 (SDK)，助力开发者快速扩展系统功能。

---

## 📂 项目结构

一个标准的插件项目应当包含以下文件：

```text
my-plugin/
├── plugin.yaml       # 插件元数据（必须）
├── main.go           # 插件入口（必须）
├── component.go      # 业务逻辑实现（必须）
└── go.mod            # 依赖管理
```

---

## 🚀 最小实现 (Minimal Implementation)

### 1. `plugin.yaml`
定义插件的基础信息，用于安装前的展示与校验。
```yaml
name: "my-plugin"
display_name: "我的插件"
version: "1.0.0"
component:
  type: "deploy"  # 可选: cert, deploy, check, notice
```

### 2. `main.go`
只需导出一个 `NewComponent` 函数，作为主程序加载插件的切入点：
```go
package main

import "git.trustasia.cn/certcloud/certm/pkg/plugin-sdk"

func NewComponent() sdk.Component {
	return &MyComponent{}
}
```

### 3. `component.go`
继承 `sdk.BaseComponent` 并实现核心逻辑。
```go
type MyComponent struct {
    sdk.BaseComponent
}

// Info 返回组件基本信息（ID 需与 plugin.yaml 保持一致）
func (c *MyComponent) Info() sdk.ComponentInfo {
    return sdk.ComponentInfo{
        ID:   "my-plugin",
        Name: "我的插件",
        Type: sdk.ComponentTypeDeploy,
    }
}

// Execute 核心业务逻辑
func (c *MyComponent) Execute(ctx context.Context, config helper.FieldConfig, input []*sdk.StepOutput) (*sdk.StepOutput, error) {
    // 逻辑实现...
    return sdk.NewStepOutput(true, nil, sdk.DataTypeNone, "Success"), nil
}
```

---

## 🧩 进阶功能

### 1. 生命周期与元数据 (Optional)
如果你需要初始化资源或在运行时自证身份，可以在 `main.go` 中添加：
```go
// 二进制层面的身份声明
var PluginMeta = sdk.PluginMetadata{
    ID:      "my-plugin",
    Version: "1.0.0",
}

// 生命周期钩子
func OnLoad() error { /* 初始化逻辑 */ return nil }
func OnUnload() error { /* 清理逻辑 */ return nil }
```

### 2. 数据访问与上下文
通过 `ctx` 获取系统资源：
- **数据查询**: `sdk.GetDataAccess(ctx)` 提供证书、部署器等只读查询接口。
- **环境信息**: `sdk.GetLang(ctx)`, `sdk.GetProjectID(ctx)`。

### 3. 配置 UI 定义
通过 `GetConfigSchema()` 返回 `helper.Field` 列表，前端会自动渲染对应的配置表单。

---

## 💡 开发建议 (Best Practices)

- **推荐组件**: 日志打印推荐使用 `logx`，HTTP 请求推荐 `httpx` (来自 `go-van` 框架)。
- **常见误区**: `plugin.yaml` 是安装包的说明书，`PluginMetadata` 是二进制的身份证。两者虽有重叠，但生效阶段不同（安装前 vs 加载后）。推荐两者保持一致。

---

## 📦 打包与发布

为了保证安全性，CertM 要求插件包必须进行签名：

### 1. 打包结构
```text
my-plugin.zip
├── plugin.so          # 编译后的产物
├── plugin.yaml        # 元数据
├── manifest.json      # 文件哈希清单
└── signature          # Ed25519 签名文件
```

### 2. 编译命令
```bash
go build -buildmode=plugin -o plugin.so .
```

*参考示例：[pkg/plugin-sdk/example/](./example/)*
