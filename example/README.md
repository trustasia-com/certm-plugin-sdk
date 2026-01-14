# 示例部署器插件 (Example Deployer)

这是一个 CertM 工作流组件插件的开发示例，展示如何创建自定义部署组件。

## 快速开始

### 1. 构建插件

```bash
make build
```

构建产物将生成在 `dist/` 目录下：`dist/example-deployer.so`

### 2. 运行测试

```bash
make test
```

## 开发指南

### 核心接口实现

参考 `component.go`：

1.  **Info()**: 返回组件元数据（类型、ID、名称等）。
2.  **GetConfigSchema()**: 定义组件的配置字段（前端将根据此动态生成表单）。
3.  **Execute()**: 实现具体的部署逻辑。接收前置步骤的输出列表并执行。
4.  **ValidateConfig()**: 在保存或执行前验证用户配置。

### 日志记录与 HTTP 调用

建议使用 `go-van` 框架中的组件：

- **日志**: 使用 `logx` 打印日志。
    ```go
    logx.WithContext(ctx).Infof("Deploying certificate for %s", certData.CommonName)
    ```
- **HTTP**: 使用 `httpx` 进行外部 API 调用。

### 数据解析

始终使用 SDK 提供的辅助方法解析输入数据：

```go
certData, err := input[0].ParseCertificate()
```

## 参考文档

- [Plugin SDK 主文档](../README.md)
- [工作流系统设计](../../../docs/工作流组件插件系统设计.md)
