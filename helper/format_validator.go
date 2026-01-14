package helper

import (
	"net"
	"strconv"
	"strings"
	"time"

	"regexp"
)

// formatValidators 格式验证器映射表
var formatValidators = map[FieldFormat]func(Field, string) error{
	FieldFormatEmail:    validateEmail,
	FieldFormatURL:      validateURL,
	FieldFormatPort:     validatePort,
	FieldFormatIP:       validateIP,
	FieldFormatCIDR:     validateCIDR,
	FieldFormatDate:     validateDate,
	FieldFormatTime:     validateTime,
	FieldFormatTel:      validateTel,
	FieldFormatPassword: validatePassword,
}

// emailRegex 邮箱正则表达式（简化版，符合 RFC 5322 的基本要求）
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// validateEmail 验证邮箱格式
func validateEmail(field Field, value string) error {
	if !emailRegex.MatchString(value) {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}

// urlRegex URL 正则表达式（支持 http、https、ftp 等协议）
var urlRegex = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)

// validateURL 验证 URL 格式
func validateURL(field Field, value string) error {
	if !urlRegex.MatchString(value) {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}

// validatePort 验证端口号格式
func validatePort(field Field, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	if port < 1 || port > 65535 {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}

// validateIP 验证 IP 地址格式
func validateIP(field Field, value string) error {
	if ip := net.ParseIP(value); ip == nil {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}

// validateCIDR 验证 CIDR 地址范围格式
func validateCIDR(field Field, value string) error {
	if _, _, err := net.ParseCIDR(value); err != nil {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}

// validateDate 验证日期格式 (YYYY-MM-DD)
func validateDate(field Field, value string) error {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}

// validateTime 验证时间格式 (HH:MM:SS 或 HH:MM)
func validateTime(field Field, value string) error {
	// 尝试解析 HH:MM:SS 格式
	_, err := time.Parse("15:04:05", value)
	if err == nil {
		return nil
	}
	// 尝试解析 HH:MM 格式
	_, err = time.Parse("15:04", value)
	if err != nil {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}

// telRegex 电话正则表达式（支持国际格式和国内格式）
var telRegex = regexp.MustCompile(`^(\+?\d{1,3}[-.\s]?)?\(?\d{1,4}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}$`)

// validateTel 验证电话格式
func validateTel(field Field, value string) error {
	// 移除空格和特殊字符后验证
	cleaned := strings.ReplaceAll(value, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// 基本验证：至少包含数字，长度在 7-15 位之间
	if len(cleaned) < 7 || len(cleaned) > 15 {
		return field.Error(ValidationErrorInvalidFormat, value)
	}

	// 检查是否全是数字（可能包含 + 号）
	hasDigit := false
	for _, r := range cleaned {
		if r >= '0' && r <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return field.Error(ValidationErrorInvalidFormat, value)
	}

	return nil
}

// validatePassword 验证密码格式（至少 6 位，可根据需求调整）
func validatePassword(field Field, value string) error {
	if len(value) < 6 {
		return field.Error(ValidationErrorInvalidFormat, value)
	}
	return nil
}
