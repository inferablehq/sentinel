package tokenizer

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestTokenizerAndDetokenizer(t *testing.T) {
	SetRandomSeed(1)

	// Set up a temporary directory for masked values
	tempDir, err := os.MkdirTemp("", "tokenizer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	os.Setenv("INFERABLE_DATA_DIR", tempDir)

	tests := []struct {
		name   string
		input  string
		except []string
	}{
		{
			name:   "Simple object",
			input:  `{"name": "John", "age": 30}`,
			except: []string{},
		},
		{
			name:   "Nested object",
			input:  `{"user": {"name": "Alice", "email": "alice@example.com"}, "score": 95}`,
			except: []string{},
		},
		{
			name:   "Array of objects",
			input:  `{"users": [{"name": "Bob", "age": 25}, {"name": "Charlie", "age": 35}]}`,
			except: []string{},
		},
		{
			name:   "Mixed types",
			input:  `{"name": "Dave", "active": true, "hobbies": ["reading", "cycling"]}`,
			except: []string{},
		},
		{
			name:   "With exceptions",
			input:  `{"user": {"name": "Eve", "email": "eve@example.com"}, "message": "Hello"}`,
			except: []string{".user.name", ".message"},
		},
		{
			name:   "Deep nesting",
			input:  `{"level1": {"level2": {"level3": {"value": "secret"}}}}`,
			except: []string{},
		},
		{
			name:   "Array with mixed types",
			input:  `{"data": [42, "text", {"key": "value"}, true]}`,
			except: []string{},
		},
		{
			name:   "Except with different types",
			input:  `{"name": "Except", "age": 30, "active": true}`,
			except: []string{".name", ".age"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input map[string]interface{}
			err := json.Unmarshal([]byte(tt.input), &input)
			if err != nil {
				t.Fatalf("Failed to unmarshal input JSON: %v", err)
				return
			}

			// Tokenize
			tokenized, err := Tokenizer(input, "", tt.except)
			if err != nil {
				t.Fatalf("Tokenizer() error = %v", err)
				return
			}

			// Detokenize
			detokenized, err := Detokenizer(tokenized, "", tt.except)
			if err != nil {
				t.Fatalf("Detokenizer() error = %v", err)
				return
			}

			// Compare original input with detokenized result
			if !reflect.DeepEqual(input, detokenized) {
				t.Errorf("Detokenized result does not match original input.\nOriginal: %v\nDetokenized: %v", input, detokenized)
				return
			}

			// Check if excepted fields were not tokenized
			for _, path := range tt.except {
				originalValue := getNestedValue(input, path)
				tokenizedValue := getNestedValue(tokenized.(map[string]interface{}), path)
				if !reflect.DeepEqual(originalValue, tokenizedValue) {
					t.Errorf("Excepted field %s was tokenized. Original: %v, Tokenized: %v", path, originalValue, tokenizedValue)
					return
				}
			}
		})
	}
}

// Helper function to get nested value from a map using a dot-separated path
func getNestedValue(m map[string]interface{}, path string) interface{} {
	keys := strings.Split(path, ".")
	current := m
	for i, key := range keys {
		if i == len(keys)-1 {
			return current[key]
		}
		if v, ok := current[key].(map[string]interface{}); ok {
			current = v
		} else {
			return nil
		}
	}
	return nil
}
