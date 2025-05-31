package Utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var stringLocales = map[string]*KeyString{}
var stringLock sync.Mutex
var filepath = "translations.json"

type KeyString struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GetString(key string, fallback string) string {
	stringLock.Lock()
	defer stringLock.Unlock()

	keyString, exists := stringLocales[key]

	if !exists {
		stringLocales[key] = &KeyString{Key: key, Value: fallback}
		return fallback
	}

	return keyString.Value
}

func SetString(key string, value string) {
	stringLock.Lock()
	defer stringLock.Unlock()

	stringLocales[key] = &KeyString{key, value}
}

func LoadStrings() {
	stringLock.Lock()
	defer stringLock.Unlock()

	jsonData, err := ReadFile(filepath)

	if nil != err {
		return
	}

	err = json.Unmarshal(jsonData, &stringLocales)

	if nil != err {
		fmt.Printf("Failed to load strings %s\n", err.Error())
		return
	}
}

func SaveStrings() {
	stringLock.Lock()
	defer stringLock.Unlock()

	jsonData, err := json.Marshal(stringLocales)
	if nil != err {
		fmt.Printf("Failed to marshal strings %s\n", err.Error())
		return
	}

	err = os.WriteFile(filepath, jsonData, 0644)

	if nil != err {
		fmt.Printf("Failed to save strings %s\n", err.Error())
		return
	}
}

func SerializeStrings() ([]byte, error) {
	stringLock.Lock()
	defer stringLock.Unlock()

	return json.Marshal(stringLocales)
}
