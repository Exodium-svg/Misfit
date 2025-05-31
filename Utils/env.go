package Utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type EnvEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Env struct {
	Path   string                 `json:"path"`
	Values map[string]interface{} `json:"values"`
	lock   sync.Mutex
}

func detectType(val interface{}) string {
	switch val.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64:
		return "int"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}

func EnvSerialize(env *Env) []EnvEntry {
	env.lock.Lock()
	defer env.lock.Unlock()

	var entries []EnvEntry
	for key, val := range env.Values {
		entry := EnvEntry{
			Key:   key,
			Value: fmt.Sprintf("%v", val),
			Type:  detectType(val),
		}
		entries = append(entries, entry)
	}

	return entries
}

func EnvGet[T any](env *Env, key string, fallback T) T {
	env.lock.Lock()
	defer env.lock.Unlock()

	value, ok := env.Values[key]
	if !ok {
		return fallback
	}

	typed, ok := value.(T)
	if !ok {
		return fallback
	}

	return typed
}

func EnvSet[T any](env *Env, key string, value T) {
	env.lock.Lock()
	defer env.lock.Unlock()

	env.Values[key] = value
}

func EnvSave(env *Env) error {
	env.lock.Lock()
	defer env.lock.Unlock()

	var lines []string

	for key, value := range env.Values {
		var line string

		switch v := value.(type) {
		case int:
			line = fmt.Sprintf("int:%s:%d", key, v)
		case bool:
			line = fmt.Sprintf("bool:%s:%t", key, v)
		case string:
			line = fmt.Sprintf("string:%s:%s", key, v)
		case float64:
			line = fmt.Sprintf("float:%s:%f", key, v)
		default:
			// Unsupported type, skip
			continue
		}

		lines = append(lines, line)
	}

	data := strings.Join(lines, "\n")
	err := os.WriteFile(env.Path, []byte(data), 0644)

	return err
}

func NewEnv(path string) (Env, error) {
	env := Env{
		Path:   path,
		Values: make(map[string]interface{}),
	}

	data, err := ReadFile(path)

	if nil != err {
		return env, err
	}

	var lines []string = strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if "" == line {
			continue // skip empty lines
		}

		var lineInfo []string = strings.Split(line, ":")

		if 3 != len(lineInfo) {
			continue
		}

		var key string = lineInfo[1]
		var value interface{}

		switch lineInfo[0] {
		case "int":
			value, _ = strconv.Atoi(lineInfo[2])
			break
		case "bool":
			value, _ = strconv.ParseBool(lineInfo[2])
			break
		case "string":
			value = lineInfo[2]
			break
		case "float":
			value, _ = strconv.ParseFloat(lineInfo[2], 64)
			break
		}

		env.Values[key] = value
	}

	return env, nil
}
