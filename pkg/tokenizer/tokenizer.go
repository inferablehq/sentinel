package tokenizer

import (
	"fmt"
	"log"
	"slices"
)

func Tokenizer(data interface{}, path string, except []string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return tokenizerObject(v, path, except)
	case []interface{}:
		return tokenizerArray(v, path, except)
	default:
		return maskValue(v)
	}
}

func tokenizerObject(obj map[string]interface{}, path string, except []string) (map[string]interface{}, error) {
	for key, value := range obj {
		currentPath := fmt.Sprintf("%s.%s", path, key)

		if slices.Contains(except, currentPath) {
			fmt.Println("except", currentPath)
			continue
		}

		log.Println(">>> tokenizing", except, currentPath)

		tokenized, err := Tokenizer(value, currentPath, except)
		if err != nil {
			return nil, err
		}
		obj[key] = tokenized
	}

	return obj, nil
}

func tokenizerArray(arr []interface{}, path string, except []string) ([]interface{}, error) {
	for i, item := range arr {
		currentPath := fmt.Sprintf("%s.%d", path, i)
		tokenized, err := Tokenizer(item, currentPath, except)
		if err != nil {
			return nil, err
		}
		arr[i] = tokenized
	}
	return arr, nil
}

func Detokenizer(data interface{}, path string, except []string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return detokenizerObject(v, path, except)
	case []interface{}:
		return detokenizerArray(v, path, except)
	case string:
		if isMask(v) {
			return unmaskValue(v)
		}
		return v, nil
	default:
		return v, nil
	}
}

func detokenizerObject(obj map[string]interface{}, path string, except []string) (map[string]interface{}, error) {
	for key, value := range obj {
		currentPath := fmt.Sprintf("%s.%s", path, key)

		if slices.Contains(except, currentPath) {
			continue
		}

		detokenized, err := Detokenizer(value, currentPath, except)
		if err != nil {
			return nil, err
		}
		obj[key] = detokenized
	}

	return obj, nil
}

func detokenizerArray(arr []interface{}, path string, except []string) ([]interface{}, error) {
	for i, item := range arr {
		currentPath := fmt.Sprintf("%s.%d", path, i)
		detokenized, err := Detokenizer(item, currentPath, except)
		if err != nil {
			return nil, err
		}
		arr[i] = detokenized
	}
	return arr, nil
}
