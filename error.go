package certm

import "fmt"

// PluginError represents a plugin-specific error
type PluginError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *PluginError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common error constructors
var (
	ErrInvalidConfig = func(msg string) error {
		return &PluginError{Code: "INVALID_CONFIG", Message: msg}
	}

	ErrExecutionFailed = func(msg string) error {
		return &PluginError{Code: "EXECUTION_FAILED", Message: msg}
	}

	ErrDependencyMissing = func(dep string) error {
		return &PluginError{Code: "DEPENDENCY_MISSING", Message: fmt.Sprintf("required dependency missing: %s", dep)}
	}
)
