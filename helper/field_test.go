package helper

import (
	"encoding/json"
	"testing"
)

// TestFieldConfig_String 测试字符串类型获取
func TestFieldConfig_String(t *testing.T) {
	config := FieldConfig{
		"name":   "test",
		"region": "cn-hangzhou",
	}

	t.Run("existing key", func(t *testing.T) {
		name := config.String("name")
		if name != "test" {
			t.Errorf("Expected 'test', got '%s'", name)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		missing := config.String("missing")
		if missing != "" {
			t.Errorf("Expected empty string for missing key, got '%s'", missing)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		config := FieldConfig{"port": 80}
		result := config.String("port")
		if result != "" {
			t.Errorf("Expected empty string for wrong type, got '%s'", result)
		}
	})
}

// TestFieldConfig_Int 测试整数类型获取及转换
func TestFieldConfig_Int(t *testing.T) {
	t.Run("direct int", func(t *testing.T) {
		config := FieldConfig{"count": 10}
		count := config.Int("count")
		if count != 10 {
			t.Errorf("Expected 10, got %d", count)
		}
	})

	t.Run("int64 conversion", func(t *testing.T) {
		config := FieldConfig{"size": int64(1024)}
		size := config.Int("size")
		if size != 1024 {
			t.Errorf("Expected 1024, got %d", size)
		}
	})

	t.Run("float64 conversion (JSON)", func(t *testing.T) {
		config := FieldConfig{"port": float64(443)}
		port := config.Int("port")
		if port != 443 {
			t.Errorf("Expected 443, got %d", port)
		}
	})

	t.Run("float32 conversion", func(t *testing.T) {
		config := FieldConfig{"timeout": float32(30)}
		timeout := config.Int("timeout")
		if timeout != 30 {
			t.Errorf("Expected 30, got %d", timeout)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		config := FieldConfig{}
		result := config.Int("missing")
		if result != 0 {
			t.Errorf("Expected 0 for missing key, got %d", result)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		config := FieldConfig{"key": "not-a-number"}
		result := config.Int("key")
		if result != 0 {
			t.Errorf("Expected 0 for wrong type, got %d", result)
		}
	})
}

// TestFieldConfig_Float 测试浮点数类型获取及转换
func TestFieldConfig_Float(t *testing.T) {
	t.Run("direct float64", func(t *testing.T) {
		config := FieldConfig{"ratio": 0.95}
		ratio := config.Float("ratio")
		if ratio != 0.95 {
			t.Errorf("Expected 0.95, got %f", ratio)
		}
	})

	t.Run("float32 conversion", func(t *testing.T) {
		config := FieldConfig{"factor": float32(1.5)}
		factor := config.Float("factor")
		if factor != 1.5 {
			t.Errorf("Expected 1.5, got %f", factor)
		}
	})

	t.Run("int conversion", func(t *testing.T) {
		config := FieldConfig{"percent": 95}
		percent := config.Float("percent")
		if percent != 95.0 {
			t.Errorf("Expected 95.0, got %f", percent)
		}
	})

	t.Run("int64 conversion", func(t *testing.T) {
		config := FieldConfig{"value": int64(100)}
		value := config.Float("value")
		if value != 100.0 {
			t.Errorf("Expected 100.0, got %f", value)
		}
	})
}

// TestFieldConfig_Boolean 测试布尔类型获取
func TestFieldConfig_Boolean(t *testing.T) {
	config := FieldConfig{
		"enabled": true,
		"debug":   false,
	}

	t.Run("true value", func(t *testing.T) {
		enabled := config.Boolean("enabled")
		if !enabled {
			t.Error("Expected true")
		}
	})

	t.Run("false value", func(t *testing.T) {
		debug := config.Boolean("debug")
		if debug {
			t.Error("Expected false")
		}
	})

	t.Run("missing key", func(t *testing.T) {
		missing := config.Boolean("missing")
		if missing {
			t.Error("Expected false for missing key")
		}
	})
}

// TestFieldConfig_Map 测试 map 类型获取
func TestFieldConfig_Map(t *testing.T) {
	config := FieldConfig{
		"metadata": map[string]any{
			"env":     "production",
			"version": "1.0.0",
		},
	}

	t.Run("existing map", func(t *testing.T) {
		metadata := config.Map("metadata")
		if metadata == nil {
			t.Fatal("Expected metadata map, got nil")
		}
		if metadata["env"] != "production" {
			t.Errorf("Expected env='production', got '%v'", metadata["env"])
		}
	})

	t.Run("missing key", func(t *testing.T) {
		missing := config.Map("missing")
		if missing != nil {
			t.Errorf("Expected nil for missing key, got %v", missing)
		}
	})
}

// TestFieldConfig_StringSlice 测试字符串切片获取
func TestFieldConfig_StringSlice(t *testing.T) {
	t.Run("direct []string", func(t *testing.T) {
		config := FieldConfig{
			"domains": []string{"example.com", "test.com"},
		}
		domains := config.StringSlice("domains")
		if len(domains) != 2 {
			t.Errorf("Expected 2 domains, got %d", len(domains))
		}
		if domains[0] != "example.com" {
			t.Errorf("Expected 'example.com', got '%s'", domains[0])
		}
	})

	t.Run("[]any from JSON", func(t *testing.T) {
		config := FieldConfig{
			"tags": []any{"tag1", "tag2", "tag3"},
		}
		tags := config.StringSlice("tags")
		if len(tags) != 3 {
			t.Errorf("Expected 3 tags, got %d", len(tags))
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		config := FieldConfig{"empty": []string{}}
		empty := config.StringSlice("empty")
		if empty == nil {
			t.Error("Expected empty slice, got nil")
		}
		if len(empty) != 0 {
			t.Errorf("Expected length 0, got %d", len(empty))
		}
	})

	t.Run("missing key", func(t *testing.T) {
		config := FieldConfig{}
		missing := config.StringSlice("missing")
		if missing != nil {
			t.Errorf("Expected nil for missing key, got %v", missing)
		}
	})
}

// TestFieldConfig_IntSlice 测试整数切片获取及转换
func TestFieldConfig_IntSlice(t *testing.T) {
	t.Run("direct []int", func(t *testing.T) {
		config := FieldConfig{
			"ports": []int{80, 443, 8080},
		}
		ports := config.IntSlice("ports")
		expected := []int{80, 443, 8080}
		if len(ports) != len(expected) {
			t.Errorf("Expected %d ports, got %d", len(expected), len(ports))
		}
		for i, port := range ports {
			if port != expected[i] {
				t.Errorf("ports[%d]: expected %d, got %d", i, expected[i], port)
			}
		}
	})

	t.Run("[]any with float64 from JSON", func(t *testing.T) {
		config := FieldConfig{
			"numbers": []any{float64(1), float64(2), float64(3)},
		}
		numbers := config.IntSlice("numbers")
		if len(numbers) != 3 {
			t.Errorf("Expected 3 numbers, got %d", len(numbers))
		}
		if numbers[0] != 1 {
			t.Errorf("Expected 1, got %d", numbers[0])
		}
	})

	t.Run("mixed number types", func(t *testing.T) {
		config := FieldConfig{
			"mixed": []any{float64(10), int(20), int64(30), float32(40)},
		}
		mixed := config.IntSlice("mixed")
		expected := []int{10, 20, 30, 40}
		if len(mixed) != len(expected) {
			t.Errorf("Expected %d items, got %d", len(expected), len(mixed))
		}
		for i, val := range mixed {
			if val != expected[i] {
				t.Errorf("mixed[%d]: expected %d, got %d", i, expected[i], val)
			}
		}
	})
}

// TestFieldConfig_FloatSlice 测试浮点数切片获取及转换
func TestFieldConfig_FloatSlice(t *testing.T) {
	t.Run("direct []float64", func(t *testing.T) {
		config := FieldConfig{
			"rates": []float64{0.1, 0.5, 0.9},
		}
		rates := config.FloatSlice("rates")
		if len(rates) != 3 {
			t.Errorf("Expected 3 rates, got %d", len(rates))
		}
	})

	t.Run("mixed number types", func(t *testing.T) {
		config := FieldConfig{
			"weights": []any{float64(0.3), float32(0.5), int(1), int64(2)},
		}
		weights := config.FloatSlice("weights")
		if len(weights) != 4 {
			t.Errorf("Expected 4 weights, got %d", len(weights))
		}
		if weights[2] != 1.0 {
			t.Errorf("Expected weights[2]=1.0, got %f", weights[2])
		}
	})
}

// TestFieldConfig_BooleanSlice 测试布尔切片获取
func TestFieldConfig_BooleanSlice(t *testing.T) {
	t.Run("direct []bool", func(t *testing.T) {
		config := FieldConfig{
			"flags": []bool{true, false, true},
		}
		flags := config.BooleanSlice("flags")
		expected := []bool{true, false, true}
		if len(flags) != len(expected) {
			t.Errorf("Expected %d flags, got %d", len(expected), len(flags))
		}
		for i, flag := range flags {
			if flag != expected[i] {
				t.Errorf("flags[%d]: expected %v, got %v", i, expected[i], flag)
			}
		}
	})

	t.Run("[]any from JSON", func(t *testing.T) {
		config := FieldConfig{
			"features": []any{true, false, true, false},
		}
		features := config.BooleanSlice("features")
		if len(features) != 4 {
			t.Errorf("Expected 4 features, got %d", len(features))
		}
	})
}

// TestFieldConfig_JSONUnmarshal 测试 JSON 反序列化兼容性
func TestFieldConfig_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "test-app",
		"port": 8080,
		"ratio": 0.95,
		"enabled": true,
		"tags": ["api", "web"],
		"numbers": [1, 2, 3],
		"weights": [0.1, 0.5, 0.9],
		"flags": [true, false],
		"metadata": {
			"env": "prod",
			"region": "cn-east"
		}
	}`

	var config FieldConfig
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// 测试各种类型
	if config.String("name") != "test-app" {
		t.Error("String type failed")
	}

	if config.Int("port") != 8080 {
		t.Error("Int type (from JSON float64) failed")
	}

	if config.Float("ratio") != 0.95 {
		t.Error("Float type failed")
	}

	if !config.Boolean("enabled") {
		t.Error("Boolean type failed")
	}

	tags := config.StringSlice("tags")
	if len(tags) != 2 || tags[0] != "api" {
		t.Error("StringSlice type failed")
	}

	numbers := config.IntSlice("numbers")
	if len(numbers) != 3 || numbers[0] != 1 {
		t.Error("IntSlice type (from JSON []any) failed")
	}

	weights := config.FloatSlice("weights")
	if len(weights) != 3 || weights[0] != 0.1 {
		t.Error("FloatSlice type failed")
	}

	flags := config.BooleanSlice("flags")
	if len(flags) != 2 || !flags[0] {
		t.Error("BooleanSlice type failed")
	}

	metadata := config.Map("metadata")
	if metadata == nil || metadata["env"] != "prod" {
		t.Error("Map type failed")
	}
}

// TestFieldConfig_Validate 测试配置验证
func TestFieldConfig_Validate(t *testing.T) {
	schema := []Field{
		{
			Type:     FieldTypeString,
			Key:      "name",
			Name:     "Name",
			Required: true,
		},
		{
			Type:     FieldTypeInt,
			Key:      "port",
			Name:     "Port",
			Required: true,
		},
		{
			Type:     FieldTypeString,
			Key:      "description",
			Name:     "Description",
			Required: false,
		},
	}

	t.Run("valid config", func(t *testing.T) {
		config := FieldConfig{
			"name": "test",
			"port": 8080,
		}
		err := config.Validate(schema)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("nil config", func(t *testing.T) {
		var config FieldConfig
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for nil config")
		}
		if !IsInputNilError(err) {
			t.Errorf("Expected InputNilError, got: %v", err)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		config := FieldConfig{
			"name": "test",
		}
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for missing required field")
		}
		if !IsRequiredFieldError(err) {
			t.Errorf("Expected RequiredFieldError, got: %v", err)
		}
		if GetErrorField(err) != "port" {
			t.Errorf("Expected error field 'port', got: %s", GetErrorField(err))
		}
	})

	t.Run("empty required field", func(t *testing.T) {
		config := FieldConfig{
			"name": "",
			"port": 8080,
		}
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for empty required field")
		}
		if !IsEmptyFieldError(err) {
			t.Errorf("Expected EmptyFieldError, got: %v", err)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		config := FieldConfig{
			"name": "test",
			"port": "not-a-number",
		}
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for invalid type")
		}
		if !IsInvalidValueError(err) && GetErrorCode(err) != ValidationErrorInvalidType {
			t.Errorf("Expected type error, got: %v", err)
		}
	})

	t.Run("optional field missing", func(t *testing.T) {
		config := FieldConfig{
			"name": "test",
			"port": 8080,
		}
		err := config.Validate(schema)
		if err != nil {
			t.Errorf("Expected no error for missing optional field, got: %v", err)
		}
	})
}

// TestFieldConfig_ValidateWithOptions 测试带选项的验证
func TestFieldConfig_ValidateWithOptions(t *testing.T) {
	schema := []Field{
		{
			Type:     FieldTypeString,
			Key:      "env",
			Name:     "Environment",
			Required: true,
			Options: []FieldOption{
				{Value: "dev", Name: "Development"},
				{Value: "staging", Name: "Staging"},
				{Value: "prod", Name: "Production"},
			},
		},
	}

	t.Run("valid option", func(t *testing.T) {
		config := FieldConfig{"env": "prod"}
		err := config.Validate(schema)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("invalid option", func(t *testing.T) {
		config := FieldConfig{"env": "invalid"}
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for invalid option")
		}
		if !IsInvalidValueError(err) {
			t.Errorf("Expected InvalidValueError, got: %v", err)
		}
	})
}

// TestFieldConfig_ValidateWithShowCondition 测试条件显示
func TestFieldConfig_ValidateWithShowCondition(t *testing.T) {
	schema := []Field{
		{
			Type:     FieldTypeString,
			Key:      "mode",
			Name:     "Mode",
			Required: true,
		},
		{
			Type:     FieldTypeInt,
			Key:      "ssl_port",
			Name:     "SSL Port",
			Required: true,
			ShowWhen: &ShowCondition{
				Key:   "mode",
				Value: "ssl",
			},
		},
	}

	t.Run("condition met - field required", func(t *testing.T) {
		config := FieldConfig{
			"mode": "ssl",
		}
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for missing ssl_port")
		}
	})

	t.Run("condition met - field provided", func(t *testing.T) {
		config := FieldConfig{
			"mode":     "ssl",
			"ssl_port": 443,
		}
		err := config.Validate(schema)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("condition not met - field not required", func(t *testing.T) {
		config := FieldConfig{
			"mode": "http",
			// ssl_port 不是必需的，因为 mode != "ssl"
		}
		err := config.Validate(schema)
		if err != nil {
			t.Errorf("Expected no error when condition not met, got: %v", err)
		}
	})
}

// TestFieldConfig_ValidateFormat 测试格式验证
func TestFieldConfig_ValidateFormat(t *testing.T) {
	t.Run("valid email", func(t *testing.T) {
		schema := []Field{
			{
				Type:   FieldTypeString,
				Format: FieldFormatEmail,
				Key:    "email",
				Name:   "Email",
			},
		}
		config := FieldConfig{"email": "test@example.com"}
		err := config.Validate(schema)
		if err != nil {
			t.Errorf("Expected no error for valid email, got: %v", err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		schema := []Field{
			{
				Type:   FieldTypeString,
				Format: FieldFormatEmail,
				Key:    "email",
				Name:   "Email",
			},
		}
		config := FieldConfig{"email": "invalid-email"}
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for invalid email")
		}
	})

	t.Run("valid port", func(t *testing.T) {
		schema := []Field{
			{
				Type:   FieldTypeString,
				Format: FieldFormatPort,
				Key:    "port",
				Name:   "Port",
			},
		}
		config := FieldConfig{"port": "8080"}
		err := config.Validate(schema)
		if err != nil {
			t.Errorf("Expected no error for valid port, got: %v", err)
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		schema := []Field{
			{
				Type:   FieldTypeString,
				Format: FieldFormatPort,
				Key:    "port",
				Name:   "Port",
			},
		}
		config := FieldConfig{"port": "99999"}
		err := config.Validate(schema)
		if err == nil {
			t.Error("Expected error for invalid port")
		}
	})
}

// TestValidationErrorHelpers 测试错误辅助函数
func TestValidationErrorHelpers(t *testing.T) {
	t.Run("IsRequiredFieldError", func(t *testing.T) {
		err := &ValidationError{Code: ValidationErrorRequired}
		if !IsRequiredFieldError(err) {
			t.Error("Expected true for RequiredFieldError")
		}
	})

	t.Run("IsEmptyFieldError", func(t *testing.T) {
		err := &ValidationError{Code: ValidationErrorEmpty}
		if !IsEmptyFieldError(err) {
			t.Error("Expected true for EmptyFieldError")
		}
	})

	t.Run("IsInvalidValueError", func(t *testing.T) {
		err := &ValidationError{Code: ValidationErrorInvalidValue}
		if !IsInvalidValueError(err) {
			t.Error("Expected true for InvalidValueError")
		}
	})

	t.Run("IsInputNilError", func(t *testing.T) {
		err := &ValidationError{Code: ValidationErrorInputNil}
		if !IsInputNilError(err) {
			t.Error("Expected true for InputNilError")
		}
	})

	t.Run("GetErrorCode", func(t *testing.T) {
		err := &ValidationError{Code: ValidationErrorRequired, Field: "test"}
		code := GetErrorCode(err)
		if code != ValidationErrorRequired {
			t.Errorf("Expected ValidationErrorRequired, got: %v", code)
		}
	})

	t.Run("GetErrorField", func(t *testing.T) {
		err := &ValidationError{Field: "username"}
		field := GetErrorField(err)
		if field != "username" {
			t.Errorf("Expected 'username', got: %s", field)
		}
	})

	t.Run("GetErrorFieldName", func(t *testing.T) {
		err := &ValidationError{FieldName: "User Name"}
		fieldName := GetErrorFieldName(err)
		if fieldName != "User Name" {
			t.Errorf("Expected 'User Name', got: %s", fieldName)
		}
	})
}
