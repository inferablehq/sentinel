package tokenizer

import (
	"encoding/json"
	"math/rand"
	"os"
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

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func setValue(key string, value Stored) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return os.WriteFile(key, jsonValue, 0644)
}

func getMaskedValue(value interface{}) (string, error) {
	key := randomString(8)
	err := setValue(key, Stored{Value: value})
	if err != nil {
		return "", err
	}
	return key, nil
}
