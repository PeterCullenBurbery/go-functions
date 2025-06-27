package gofunctions

import (
	"strings"
)

func GetCaseInsensitiveMap(m map[string]interface{}, key string) map[string]interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if result, ok := v.(map[string]interface{}); ok {
				return result
			}
		}
	}
	return nil
}

func GetCaseInsensitiveList(m map[string]interface{}, key string) []string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			raw, ok := v.([]interface{})
			if !ok {
				return nil
			}
			var result []string
			for _, val := range raw {
				if s, ok := val.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

func GetCaseInsensitiveString(m map[string]interface{}, key string) string {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			if str, ok := v.(string); ok {
				return str
			}
		}
	}
	return ""
}

func GetNestedString(m map[string]interface{}, key string) string {
	if val := GetCaseInsensitiveString(m, key); val != "" {
		return val
	}
	if sub := GetCaseInsensitiveMap(m, key); sub != nil {
		for _, v := range sub {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

func GetNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	return GetCaseInsensitiveMap(m, key)
}