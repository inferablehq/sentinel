package tokenizer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Stored struct {
	Value interface{}
}

var (
	// Create a new random number generator with a seed
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func SetRandomSeed(seed int64) {
	rng = rand.New(rand.NewSource(seed))
}

func valueId(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return fmt.Sprintf("$$masked_%s", string(b))
}

func setValue(key string, value Stored) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return os.WriteFile(key, jsonValue, 0644)
}

func maskValue(value interface{}) (string, error) {
	key := valueId(8)

	dataPath := filepath.Join(os.Getenv("INFERABLE_DATA_DIR"), key)

	err := setValue(dataPath, Stored{Value: value})
	if err != nil {
		return "", err
	}
	return key, nil
}

func isMask(key string) bool {
	return strings.HasPrefix(key, "$$masked_")
}

func unmaskValue(key string) (interface{}, error) {
	dataPath := filepath.Join(os.Getenv("INFERABLE_DATA_DIR"), key)

	if value, err := os.ReadFile(dataPath); err != nil {
		return nil, err
	} else {
		var stored Stored
		if err := json.Unmarshal(value, &stored); err != nil {
			return nil, err
		}
		return stored.Value, nil
	}
}
