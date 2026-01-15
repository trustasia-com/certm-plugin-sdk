package helper

import (
	"fmt"
	"reflect"
)

// FieldType 字段类型
type FieldType string

// String 字段类型字符串
func (f FieldType) String() string {
	return string(f)
}

// FieldType 字段类型
const (
	FieldTypeString       FieldType = "string"        // 字符串
	FieldTypeInt          FieldType = "int"           // 整数
	FieldTypeFloat        FieldType = "float"         // 浮点数
	FieldTypeBoolean      FieldType = "boolean"       // 布尔值
	FieldTypeStringArray  FieldType = "string_array"  // 字符串数组
	FieldTypeIntArray     FieldType = "int_array"     // 整数数组
	FieldTypeFloatArray   FieldType = "float_array"   // 浮点数数组
	FieldTypeBooleanArray FieldType = "boolean_array" // 布尔值数组
	FieldTypeObject       FieldType = "object"        // 对象
)

// FieldFormat 字段格式
type FieldFormat string

// String 字段格式字符串
func (f FieldFormat) String() string {
	return string(f)
}

// FieldFormat 字段格式
const (
	FieldFormatTextarea FieldFormat = "textarea" // 多行文本：textarea
	FieldFormatText     FieldFormat = "text"     // 单行文本：input[type="text"]
	FieldFormatNumber   FieldFormat = "number"   // 数字：input[type="number"]

	FieldFormatRadio         FieldFormat = "radio"          // 单选：radio
	FieldFormatCheckbox      FieldFormat = "checkbox"       // 多选：input[type="checkbox"]
	FieldFormatCheckboxGroup FieldFormat = "checkbox-group" // 多选组：input[type="checkbox"]
	FieldFormatSelect        FieldFormat = "select"         // 下拉：select
	FieldFormatSelectGroup   FieldFormat = "select-group"   // 下拉组：select multiple
	FieldFormatFile          FieldFormat = "file"           // 文件：input[type="file"]
	FieldFormatFilePaste     FieldFormat = "file_paste"     // 文件+粘贴：input[type="file"]

	// 常见格式
	FieldFormatEmail    FieldFormat = "email"    // 邮箱：input[type="email"]
	FieldFormatURL      FieldFormat = "url"      // URL：input[type="url"]
	FieldFormatPassword FieldFormat = "password" // 密码：input[type="password"]
	FieldFormatDate     FieldFormat = "date"     // 日期：input[type="date"]
	FieldFormatTime     FieldFormat = "time"     // 时间：input[type="time"]
	FieldFormatTel      FieldFormat = "tel"      // 电话：input[type="tel"]
	FieldFormatIP       FieldFormat = "ip"       // IP地址：input[type="text"]
	FieldFormatCIDR     FieldFormat = "cidr"     // CIDR地址范围：input[type="text"]
	FieldFormatPort     FieldFormat = "port"     // 端口：input[type="number"]
)

// FieldOption 字段选项
type FieldOption struct {
	Value any    `json:"value"` // 具体值
	Name  string `json:"name"`  // 显示值
}

// ShowCondition 显示条件
type ShowCondition struct {
	Key   string `json:"key"`   // 字段key
	Value any    `json:"value"` // 具体值
}

// OptionsSource 选项来源
type OptionsSource struct {
	Endpoint  string   `json:"endpoint,omitempty"`   // API端点
	Method    string   `json:"method,omitempty"`     // API方法
	Params    []string `json:"params,omitempty"`     // API参数字段列表
	DependsOn []string `json:"depends_on,omitempty"` // 依赖字段（当这些字段变化时重新加载）
}

// Field 字段描述
type Field struct {
	Type     FieldType   `json:"type"`     // 字段类型，用于传入后端值类型
	Format   FieldFormat `json:"format"`   // 字段格式，用于限制输入值的格式
	Name     string      `json:"name"`     // 字段显示名称
	Key      string      `json:"key"`      // 字段key
	Default  any         `json:"default"`  // 字段默认值
	Required bool        `json:"required"` // 是否必须

	Description   string         `json:"description,omitempty"`    // 解释该字段
	DocURL        string         `json:"doc_url,omitempty"`        // 文档链接
	Options       []FieldOption  `json:"options,omitempty"`        // 静态选项
	OptionsSource *OptionsSource `json:"options_source,omitempty"` // 动态选项
	ShowWhen      *ShowCondition `json:"show_when,omitempty"`      // 显示条件
}

// Error 字段错误
func (f *Field) Error(code ValidationErrorCode, value any) error {
	return &ValidationError{Field: f.Key, FieldName: f.Name, Value: value, Code: code}
}

// ValidationErrorCode 验证错误码
type ValidationErrorCode string

const (
	ValidationErrorRequired      ValidationErrorCode = "FIELD_REQUIRED" // 必填字段错误
	ValidationErrorEmpty         ValidationErrorCode = "FIELD_EMPTY"    // 字段为空错误
	ValidationErrorInvalidValue  ValidationErrorCode = "INVALID_VALUE"  // 无效值错误
	ValidationErrorInputNil      ValidationErrorCode = "INPUT_NIL"      // 输入为空错误
	ValidationErrorInvalidFormat ValidationErrorCode = "INVALID_FORMAT" // 无效格式错误
	ValidationErrorInvalidType   ValidationErrorCode = "INVALID_TYPE"   // 无效类型错误
)

// ValidationError 验证错误类型
type ValidationError struct {
	Field     string              // 字段Key
	FieldName string              // 字段显示名称
	Value     any                 // 字段值
	Code      ValidationErrorCode // 错误码
}

func (e *ValidationError) Error() string {
	// 提供默认英文消息，实际使用时调用方应该根据Code生成本地化消息
	switch e.Code {
	case ValidationErrorRequired:
		return fmt.Sprintf("field '%s' is required", e.FieldName)
	case ValidationErrorEmpty:
		return fmt.Sprintf("field '%s' cannot be empty", e.FieldName)
	case ValidationErrorInvalidValue:
		return fmt.Sprintf("field '%s' has invalid value: %v", e.FieldName, e.Value)
	case ValidationErrorInputNil:
		return "input is nil"
	case ValidationErrorInvalidFormat:
		return fmt.Sprintf("field '%s' has invalid format: %v", e.FieldName, e.Value)
	case ValidationErrorInvalidType:
		return fmt.Sprintf("field '%s' has invalid type: %T", e.FieldName, e.Value)
	default:
		return "validation error"
	}
}

// IsRequiredFieldError 判断是否为必填字段错误
func IsRequiredFieldError(err error) bool {
	if ve, ok := err.(*ValidationError); ok {
		return ve.Code == ValidationErrorRequired
	}
	return false
}

// IsEmptyFieldError 判断是否为字段为空错误
func IsEmptyFieldError(err error) bool {
	if ve, ok := err.(*ValidationError); ok {
		return ve.Code == ValidationErrorEmpty
	}
	return false
}

// IsInvalidValueError 判断是否为无效值错误
func IsInvalidValueError(err error) bool {
	if ve, ok := err.(*ValidationError); ok {
		return ve.Code == ValidationErrorInvalidValue
	}
	return false
}

// IsInputNilError 判断是否为输入为空错误
func IsInputNilError(err error) bool {
	if ve, ok := err.(*ValidationError); ok {
		return ve.Code == ValidationErrorInputNil
	}
	return false
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) ValidationErrorCode {
	if ve, ok := err.(*ValidationError); ok {
		return ve.Code
	}
	return ""
}

// GetErrorField 获取错误相关的字段Key
func GetErrorField(err error) string {
	if ve, ok := err.(*ValidationError); ok {
		return ve.Field
	}
	return ""
}

// GetErrorFieldName 获取错误相关的字段显示名称
func GetErrorFieldName(err error) string {
	if ve, ok := err.(*ValidationError); ok {
		return ve.FieldName
	}
	return ""
}

// FieldConfig 字段配置
type FieldConfig map[string]any

// getValue 泛型辅助函数，统一处理字段获取和类型转换
func getValue[T any](f FieldConfig, key string, typeName string) T {
	var zero T
	val, ok := f[key]
	if !ok {
		fmt.Printf("[Helper] field %s not found\n", key)
		return zero
	}

	result, ok := val.(T)
	if !ok {
		fmt.Printf("[Helper] field %s has invalid type, expected %s, got %T\n", key, typeName, val)
		return zero
	}
	return result
}

// getSliceValue 泛型切片辅助函数，支持 JSON 反序列化的 []any 转换
func getSliceValue[T any](f FieldConfig, key string, typeName string, converter func(any) (T, bool)) []T {
	val, ok := f[key]
	if !ok {
		fmt.Printf("[Helper] field %s not found\n", key)
		return nil
	}

	// 直接是目标类型的切片
	if v, ok := val.([]T); ok {
		return v
	}

	// 处理 JSON 反序列化的 []any
	if arr, ok := val.([]any); ok {
		result := make([]T, 0, len(arr))
		for i, item := range arr {
			if converted, ok := converter(item); ok {
				result = append(result, converted)
			} else {
				fmt.Printf("[Helper] field %s[%d] has invalid type, expected %s, got %T\n", key, i, typeName, item)
			}
		}
		return result
	}

	fmt.Printf("[Helper] field %s has invalid type, expected []%s or []any, got %T\n", key, typeName, val)
	return nil
}

// String 字段配置字符串
func (f FieldConfig) String(key string) string {
	return getValue[string](f, key, "string")
}

// toInt 数字类型转换为 int
func toInt(val any) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case float32:
		return int(v), true
	default:
		return 0, false
	}
}

// toFloat 数字类型转换为 float64
func toFloat(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

// Int 字段配置整数
func (f FieldConfig) Int(key string) int {
	val := f[key]
	if result, ok := toInt(val); ok {
		return result
	}

	fmt.Printf("[Helper] field %s has invalid type, expected numeric, got %T\n", key, val)
	return 0
}

// Float 字段配置浮点数
func (f FieldConfig) Float(key string) float64 {
	val := f[key]
	if result, ok := toFloat(val); ok {
		return result
	}
	fmt.Printf("[Helper] field %s has invalid type, expected numeric, got %T\n", key, val)
	return 0
}

// Boolean 字段配置布尔值
func (f FieldConfig) Boolean(key string) bool {
	return getValue[bool](f, key, "bool")
}

// StringSlice 字段配置字符串切片
func (f FieldConfig) StringSlice(key string) []string {
	return getSliceValue(f, key, "string", func(item any) (string, bool) {
		str, ok := item.(string)
		return str, ok
	})
}

// IntSlice 字段配置整数切片
func (f FieldConfig) IntSlice(key string) []int {
	return getSliceValue(f, key, "int", func(item any) (int, bool) {
		return toInt(item)
	})
}

// FloatSlice 字段配置浮点数切片
func (f FieldConfig) FloatSlice(key string) []float64 {
	return getSliceValue(f, key, "float64", toFloat)
}

// BooleanSlice 字段配置布尔值切片
func (f FieldConfig) BooleanSlice(key string) []bool {
	return getSliceValue(f, key, "bool", func(item any) (bool, bool) {
		b, ok := item.(bool)
		return b, ok
	})
}

// Map 字段配置map
func (f FieldConfig) Map(key string) map[string]any {
	return getValue[map[string]any](f, key, "map[string]any")
}

func (f FieldConfig) Validate(fields []Field) error {
	// 检查输入是否为 nil
	if f == nil {
		return &ValidationError{Code: ValidationErrorInputNil}
	}

	for _, field := range fields {
		// 检查显示条件
		if !shouldShowField(field, f) {
			continue // 跳过不满足显示条件的字段
		}

		value, exists := f[field.Key]
		if !exists { // 处理缺失字段
			if field.Required {
				return field.Error(ValidationErrorRequired, value)
			}
			continue
		}

		// 检查值是否为空
		if isEmpty(value) {
			if field.Required {
				return field.Error(ValidationErrorEmpty, value)
			}
			continue // 跳过非必填的空字段
		}

		// 类型验证
		newValue, err := validateFieldType(field, value)
		if err != nil {
			return field.Error(ValidationErrorInvalidType, value)
		}

		// 选项验证
		if err := validateFieldOptions(field, newValue); err != nil {
			return field.Error(ValidationErrorInvalidValue, value)
		}

		// 格式验证
		if err := validateFieldFormat(field, newValue); err != nil {
			return field.Error(ValidationErrorInvalidFormat, value)
		}
	}
	return nil
}

// shouldShowField 检查字段是否应该显示
func shouldShowField(field Field, input FieldConfig) bool {
	// 如果没有显示条件，总是显示
	if field.ShowWhen == nil {
		return true
	}

	// 检查依赖字段的值
	dependValue, exists := input[field.ShowWhen.Key]
	if !exists {
		return false // 依赖字段不存在，不显示
	}

	// 检查值是否匹配
	return reflect.DeepEqual(dependValue, field.ShowWhen.Value)
}

// isEmpty 检查值是否为空
func isEmpty(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

// validateFieldType 验证字段类型
func validateFieldType(field Field, value any) (any, error) {
	switch field.Type {
	case FieldTypeString:
		if val, ok := value.(string); ok {
			return val, nil
		}
	case FieldTypeInt:
		// 复用 toInt 转换器来验证
		if val, ok := toInt(value); ok {
			return val, nil
		}
	case FieldTypeFloat:
		if val, ok := toFloat(value); ok {
			return val, nil
		}
	case FieldTypeBoolean:
		if val, ok := value.(bool); ok {
			return val, nil
		}
	case FieldTypeStringArray:
		// 直接是 []string
		if val, ok := value.([]string); ok {
			return val, nil
		}
		// 处理 JSON 反序列化的 []any
		if arr, ok := value.([]any); ok {
			for _, item := range arr {
				if _, ok := item.(string); !ok {
					return nil, field.Error(ValidationErrorInvalidType, value)
				}
			}
			return arr, nil
		}
	case FieldTypeIntArray:
		// 直接是 []int
		if val, ok := value.([]int); ok {
			return val, nil
		}
		// 处理 JSON 反序列化的 []any
		if arr, ok := value.([]any); ok {
			for _, item := range arr {
				if _, ok := toInt(item); !ok {
					return nil, field.Error(ValidationErrorInvalidType, value)
				}
			}
			return arr, nil
		}
	case FieldTypeFloatArray:
		// 直接是 []float64
		if val, ok := value.([]float64); ok {
			return val, nil
		}
		// 处理 JSON 反序列化的 []any
		if arr, ok := value.([]any); ok {
			for _, item := range arr {
				if _, ok := toFloat(item); !ok {
					return nil, field.Error(ValidationErrorInvalidType, value)
				}
			}
			return arr, nil
		}
	case FieldTypeBooleanArray:
		// 直接是 []bool
		if val, ok := value.([]bool); ok {
			return val, nil
		}
		// 处理 JSON 反序列化的 []any
		if arr, ok := value.([]any); ok {
			for _, item := range arr {
				if _, ok := item.(bool); !ok {
					return nil, field.Error(ValidationErrorInvalidType, value)
				}
			}
			return arr, nil
		}
	case FieldTypeObject:
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Map || v.Kind() == reflect.Struct {
			return value, nil
		}
	}

	return nil, field.Error(ValidationErrorInvalidType, value)
}

// validateFieldOptions 验证字段选项
func validateFieldOptions(field Field, value any) error {
	if len(field.Options) == 0 {
		return nil // 没有选项限制
	}

	// 检查值是否在选项中
	for _, option := range field.Options {
		if reflect.DeepEqual(option.Value, value) {
			return nil // 找到匹配选项
		}
	}

	return field.Error(ValidationErrorInvalidValue, value)
}

// validateFieldFormat 验证字段格式
func validateFieldFormat(field Field, value any) error {
	// 只对字符串类型进行格式验证
	str, ok := value.(string)
	if !ok {
		return nil // 非字符串类型跳过格式验证
	}

	// 查找并执行对应的验证器
	if validator, ok := formatValidators[field.Format]; ok {
		return validator(field, str)
	}

	return nil
}
