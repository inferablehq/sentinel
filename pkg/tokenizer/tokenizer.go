package tokenizer

import (
	"fmt"
	"slices"
)

func Tokenizer(obj map[string]interface{}, path string, except []string) map[string]interface{} {
	for key, value := range obj {
		currentPath := fmt.Sprintf("%s.%s", path, key)

		if slices.Contains(except, currentPath) {
			fmt.Println("except", currentPath)
			continue
		}

		switch v := value.(type) {
		case string, int, float64:
			fmt.Println(currentPath, v)
			obj[key] = "MASKED"
		case map[string]interface{}:
			obj[key] = Tokenizer(v, currentPath, except)
		case []interface{}:
			for i, item := range v {
				if subMap, ok := item.(map[string]interface{}); ok {
					v[i] = Tokenizer(subMap, fmt.Sprintf("%s.%d", currentPath, i), except)
				} else if _, ok := item.(string); ok {
					v[i] = "MASKED"
				} else if _, ok := item.(int); ok {
					v[i] = "MASKED"
				} else if _, ok := item.(float64); ok {
					v[i] = "MASKED"
				}
			}
		}
	}

	return obj
}
