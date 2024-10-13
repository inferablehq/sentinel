package tokenizer

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		except []string
		want   string
	}{
		{
			name:   "Simple object",
			input:  `{"name": "John", "age": 30}`,
			except: []string{},
			want:   `{"age": "MASKED", "name": "MASKED"}`,
		},
		{
			name:   "Nested object",
			input:  `{"user": {"name": "Alice", "email": "alice@example.com"}, "score": 95}`,
			except: []string{},
			want:   `{"score": "MASKED", "user": {"email": "MASKED", "name": "MASKED"}}`,
		},
		{
			name:   "Array of objects",
			input:  `{"users": [{"name": "Bob", "age": 25}, {"name": "Charlie", "age": 35}]}`,
			except: []string{},
			want:   `{"users": [{"age": "MASKED", "name": "MASKED"}, {"age": "MASKED", "name": "MASKED"}]}`,
		},
		{
			name:   "Mixed types",
			input:  `{"name": "Dave", "active": true, "hobbies": ["reading", "cycling"]}`,
			except: []string{},
			want:   `{"active": true, "hobbies": ["MASKED", "MASKED"], "name": "MASKED"}`,
		},
		{
			name:   "With exceptions",
			input:  `{"user": {"name": "Eve", "email": "eve@example.com"}, "message": "Hello"}`,
			except: []string{".user.name", ".message"},
			want:   `{"message": "Hello", "user": {"email": "MASKED", "name": "Eve"}}`,
		},
		{
			name:   "Deep nesting",
			input:  `{"level1": {"level2": {"level3": {"value": "secret"}}}}`,
			except: []string{},
			want:   `{"level1": {"level2": {"level3": {"value": "MASKED"}}}}`,
		},
		{
			name:   "Array with mixed types",
			input:  `{"data": [42, "text", {"key": "value"}, true]}`,
			except: []string{},
			want:   `{"data": ["MASKED", "MASKED", {"key": "MASKED"}, true]}`,
		},
		{
			name:   "Except with different types",
			input:  `{"name": "Except", "age": 30, "active": true}`,
			except: []string{".name", ".age"},
			want:   `{"active": true, "age": 30, "name": "Except"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input map[string]interface{}
			err := json.Unmarshal([]byte(tt.input), &input)
			if err != nil {
				t.Fatalf("Failed to unmarshal input JSON: %v", err)
			}

			result := Tokenizer(input, "", tt.except)

			resultJSON, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to marshal result to JSON: %v", err)
			}

			var want, got map[string]interface{}
			json.Unmarshal([]byte(tt.want), &want)
			json.Unmarshal(resultJSON, &got)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("Tokenizer() = %v, want %v", string(resultJSON), tt.want)
			}
		})
	}
}
