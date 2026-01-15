//go:build tinygo || wasm
// +build tinygo wasm

package certm

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"unsafe"
)

// ptrFrom 从字节切片获取指针
func ptrFrom(data []byte) uint32 {
	if len(data) == 0 {
		return 0
	}
	return uint32(uintptr(unsafe.Pointer(&data[0])))
}

// ptrToSlice 从指针和长度重建切片
func ptrToSlice(ptr, length uint32) []byte {
	if ptr == 0 || length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

// JSON json & write to memory
func (r *Result) writeToMemory(ptr *uint32) {
	data, _ := json.Marshal(r)

	size := uint32(len(data))

	// 创建buffer: 4字节长度 + 数据
	buf := make([]byte, 4+size)

	// 写入长度
	binary.LittleEndian.PutUint32(buf[0:4], size)

	// 写入数据
	copy(buf[4:], data)

	// 返回buffer起始指针
	*ptr = ptrFrom(buf)
}

// readFromMemory 从WASM内存读取数据
func readFromMemory(ptr uint32) []byte {
	if ptr == 0 {
		return nil
	}

	// 读取长度
	lengthBytes := ptrToSlice(ptr, 4)
	length := binary.LittleEndian.Uint32(lengthBytes)

	// 读取数据
	return ptrToSlice(ptr+4, length)
}

// call 统一的主机函数调用（使用泛型）
func call[T any](fnName string, args ...interface{}) (*T, error) {
	// 1. 序列化参数
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	// 2. 调用主机函数
	fnNameBytes := []byte(fnName)
	resultPtrAndLen := hostCall(
		ptrFrom(fnNameBytes), uint32(len(fnNameBytes)),
		ptrFrom(argsJSON), uint32(len(argsJSON)),
	)

	// 3. 解析返回值（高32位是指针，低32位是长度）
	resultPtr := uint32(resultPtrAndLen >> 32)
	resultLen := uint32(resultPtrAndLen & 0xFFFFFFFF)

	resultJSON := ptrToSlice(resultPtr, resultLen)

	// 4. 检查是否是错误
	var errResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(resultJSON, &errResp); err == nil && errResp.Error != "" {
		return nil, fmt.Errorf("%s", errResp.Error)
	}

	// 5. 反序列化结果
	var result T
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		return nil, fmt.Errorf("unmarshal result: %w", err)
	}

	return &result, nil
}

// hostLogImpl 日志输出（封装主机函数）
func hostLogImpl(level, msg string) {
	levelBytes := []byte(level)
	msgBytes := []byte(msg)

	hostLog(
		ptrFrom(levelBytes), uint32(len(levelBytes)),
		ptrFrom(msgBytes), uint32(len(msgBytes)),
	)
}
