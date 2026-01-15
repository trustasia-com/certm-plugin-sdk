//go:build tinygo || wasm
// +build tinygo wasm

package certm

// host_call 统一的主机函数调用入口
//
//go:wasm-module env
//export host_call
func hostCall(fnNamePtr, fnNameLen, argsPtr, argsLen uint32) uint64

// host_log 日志输出
//
//go:wasm-module env
//export host_log
func hostLog(levelPtr, levelLen, msgPtr, msgLen uint32)
