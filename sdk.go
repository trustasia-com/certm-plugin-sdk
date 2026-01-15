//go:build tinygo || wasm
// +build tinygo wasm

package certm

var component Component

// Register 注册组件实现
func Register(c Component) { component = c }
