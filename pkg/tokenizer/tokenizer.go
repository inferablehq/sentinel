package tokenizer

import (
	"fmt"
	"slices"
)

func Tokenizer(obj map[string]interface{}, path string, except []string) (map[string]interface{}, error) {
	for key, value := range obj {
		currentPath := fmt.Sprintf("%s.%s", path, key)

		if slices.Contains(except, currentPath) {
			fmt.Println("except", currentPath)
			continue
		}

		switch v := value.(type) {
		case string, int, float64:
			fmt.Println(currentPath, v)
			m, err := getMaskedValue(v)
			if err != nil {
				return nil, err
			}
			obj[key] = m
		case map[string]interface{}:
			subMap, err := Tokenizer(v, currentPath, except)
			if err != nil {
				return nil, err
			}
			obj[key] = subMap
		case []interface{}:
			for i, item := range v {
				if subMap, ok := item.(map[string]interface{}); ok {
					t, err := Tokenizer(subMap, fmt.Sprintf("%s.%d", currentPath, i), except)
					v[i] = t
					if err != nil {
						return nil, err
					}
				} else if _, ok := item.(string); ok {
					m, err := getMaskedValue(item)
					if err != nil {
						return nil, err
					}
					v[i] = m
				} else if _, ok := item.(int); ok {
					m, err := getMaskedValue(item)
					if err != nil {
						return nil, err
					}
					v[i] = m
				}
			}
		}
	}

	return obj, nil
}
