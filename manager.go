package mflag

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// mapManager holds configuration values.
// It supports nested structures, which can be accessed using dot notation (e.g., "database.host").
type mapManager struct {
	data map[string]interface{}
}

// newManager creates and returns a new, empty mapManager.
func newManager() *mapManager {
	return &mapManager{
		data: make(map[string]interface{}),
	}
}

// Clone creates a deep copy of the mapManager.
func (m *mapManager) Clone() *mapManager {
	return &mapManager{
		data: deepCopyMap(m.data),
	}
}

// Merge merges another mapManager into this one. Values in the other manager
// take precedence by overwriting existing keys.
func (m *mapManager) Merge(other *mapManager) {
	m.data = mergeMaps(m.data, other.data)
}

// LoadFile reads a YAML configuration file from the specified path and populates the config.
func (m *mapManager) LoadFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		// It's not an error if the file doesn't exist; we just won't load it.
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("%w: failed to read config file %s: %w", ErrInitFailed, filename, err)
	}

	var parsedData map[string]interface{}
	if err := yaml.Unmarshal(content, &parsedData); err != nil {
		return fmt.Errorf("%w: failed to parse yaml: %w", ErrInitFailed, err)
	}

	// The YAML library can create map[any]any, which we need to convert.
	m.data = convertMap(parsedData)
	return nil
}

// SetValue sets a value for a given key. The key can be a dot-separated path to create nested maps.
func (m *mapManager) SetValue(key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := m.data

	for i, k := range keys {
		if i == len(keys)-1 {
			// This is the last key, so set the value.
			current[k] = value
		} else {
			// This is a key in the path.
			if _, exists := current[k]; !exists {
				// Create a new map if the key doesn't exist.
				current[k] = make(map[string]interface{})
			}
			// Move to the next level.
			if nested, ok := current[k].(map[string]interface{}); ok {
				current = nested
			} else {
				// A value already exists at this path but it's not a map,
				// so we cannot create a nested key. We'll overwrite it.
				newMap := make(map[string]interface{})
				current[k] = newMap
				current = newMap
			}
		}
	}
}

// Get retrieves a configuration value by key.
func (m *mapManager) Get(key string) interface{} {
	keys := strings.Split(key, ".")
	var current interface{} = m.data

	for _, k := range keys {
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil // Cannot traverse further down a non-map value.
		}

		value, exists := currentMap[k]
		if !exists {
			return nil
		}
		current = value
	}
	return current
}

// GetString returns the value associated with the key as a string.
func (m *mapManager) GetString(key string) string {
	val := m.Get(key)
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

// GetInt returns the value associated with the key as an integer.
func (m *mapManager) GetInt(key string) int {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64: // YAML can unmarshal to int64
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// GetInt8 returns the value associated with the key as an int8.
func (m *mapManager) GetInt8(key string) int8 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int8:
		return v
	case int:
		return int8(v)
	case int16:
		return int8(v)
	case int32:
		return int8(v)
	case int64:
		return int8(v)
	case uint:
		return int8(v)
	case uint8:
		return int8(v)
	case uint16:
		return int8(v)
	case uint32:
		return int8(v)
	case uint64:
		return int8(v)
	case float64:
		return int8(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 8); err == nil {
			return int8(i)
		}
	}
	return 0
}

// GetInt16 returns the value associated with the key as an int16.
func (m *mapManager) GetInt16(key string) int16 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int16:
		return v
	case int:
		return int16(v)
	case int8:
		return int16(v)
	case int32:
		return int16(v)
	case int64:
		return int16(v)
	case uint:
		return int16(v)
	case uint8:
		return int16(v)
	case uint16:
		return int16(v)
	case uint32:
		return int16(v)
	case uint64:
		return int16(v)
	case float64:
		return int16(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 16); err == nil {
			return int16(i)
		}
	}
	return 0
}

// GetInt32 returns the value associated with the key as an int32.
func (m *mapManager) GetInt32(key string) int32 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int32:
		return v
	case int:
		return int32(v)
	case int8:
		return int32(v)
	case int16:
		return int32(v)
	case int64:
		return int32(v)
	case uint:
		return int32(v)
	case uint8:
		return int32(v)
	case uint16:
		return int32(v)
	case uint32:
		return int32(v)
	case uint64:
		return int32(v)
	case float64:
		return int32(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			return int32(i)
		}
	}
	return 0
}

// GetInt64 returns the value associated with the key as an int64.
func (m *mapManager) GetInt64(key string) int64 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// GetUint returns the value associated with the key as a uint.
func (m *mapManager) GetUint(key string) uint {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case uint:
		return v
	case uint8:
		return uint(v)
	case uint16:
		return uint(v)
	case uint32:
		return uint(v)
	case uint64:
		return uint(v)
	case int:
		if v < 0 {
			return 0
		}
		return uint(v)
	case int8:
		if v < 0 {
			return 0
		}
		return uint(v)
	case int16:
		if v < 0 {
			return 0
		}
		return uint(v)
	case int32:
		if v < 0 {
			return 0
		}
		return uint(v)
	case int64:
		if v < 0 {
			return 0
		}
		return uint(v)
	case float64:
		if v < 0 {
			return 0
		}
		return uint(v)
	case string:
		if i, err := strconv.ParseUint(v, 10, 0); err == nil {
			return uint(i)
		}
	}
	return 0
}

// GetUint8 returns the value associated with the key as a uint8.
func (m *mapManager) GetUint8(key string) uint8 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case uint8:
		return v
	case uint:
		return uint8(v)
	case uint16:
		return uint8(v)
	case uint32:
		return uint8(v)
	case uint64:
		return uint8(v)
	case int:
		if v < 0 {
			return 0
		}
		return uint8(v)
	case int8:
		if v < 0 {
			return 0
		}
		return uint8(v)
	case int16:
		if v < 0 {
			return 0
		}
		return uint8(v)
	case int32:
		if v < 0 {
			return 0
		}
		return uint8(v)
	case int64:
		if v < 0 {
			return 0
		}
		return uint8(v)
	case float64:
		if v < 0 {
			return 0
		}
		return uint8(v)
	case string:
		if i, err := strconv.ParseUint(v, 10, 8); err == nil {
			return uint8(i)
		}
	}
	return 0
}

// GetUint16 returns the value associated with the key as a uint16.
func (m *mapManager) GetUint16(key string) uint16 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case uint16:
		return v
	case uint:
		return uint16(v)
	case uint8:
		return uint16(v)
	case uint32:
		return uint16(v)
	case uint64:
		return uint16(v)
	case int:
		if v < 0 {
			return 0
		}
		return uint16(v)
	case int8:
		if v < 0 {
			return 0
		}
		return uint16(v)
	case int16:
		if v < 0 {
			return 0
		}
		return uint16(v)
	case int32:
		if v < 0 {
			return 0
		}
		return uint16(v)
	case int64:
		if v < 0 {
			return 0
		}
		return uint16(v)
	case float64:
		if v < 0 {
			return 0
		}
		return uint16(v)
	case string:
		if i, err := strconv.ParseUint(v, 10, 16); err == nil {
			return uint16(i)
		}
	}
	return 0
}

// GetUint32 returns the value associated with the key as a uint32.
func (m *mapManager) GetUint32(key string) uint32 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case uint32:
		return v
	case uint:
		return uint32(v)
	case uint8:
		return uint32(v)
	case uint16:
		return uint32(v)
	case uint64:
		return uint32(v)
	case int:
		if v < 0 {
			return 0
		}
		return uint32(v)
	case int8:
		if v < 0 {
			return 0
		}
		return uint32(v)
	case int16:
		if v < 0 {
			return 0
		}
		return uint32(v)
	case int32:
		if v < 0 {
			return 0
		}
		return uint32(v)
	case int64:
		if v < 0 {
			return 0
		}
		return uint32(v)
	case float64:
		if v < 0 {
			return 0
		}
		return uint32(v)
	case string:
		if i, err := strconv.ParseUint(v, 10, 32); err == nil {
			return uint32(i)
		}
	}
	return 0
}

// GetUint64 returns the value associated with the key as a uint64.
func (m *mapManager) GetUint64(key string) uint64 {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case uint64:
		return v
	case uint:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case int:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int8:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int16:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int32:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case int64:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case float64:
		if v < 0 {
			return 0
		}
		return uint64(v)
	case string:
		if i, err := strconv.ParseUint(v, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// GetBool returns the value associated with the key as a boolean.
func (m *mapManager) GetBool(key string) bool {
	val := m.Get(key)
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return false
}

// GetFloat64 returns the value associated with the key as a float64.
func (m *mapManager) GetFloat64(key string) float64 {
	val := m.Get(key)
	if val == nil {
		return 0.0
	}
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

// GetDuration returns the value associated with the key as a time.Duration.
// It can parse duration strings (e.g., "10s", "5m").
// If the value is a number, it's treated as nanoseconds.
func (m *mapManager) GetDuration(key string) time.Duration {
	val := m.Get(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case time.Duration:
		return v
	case string:
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	case int:
		return time.Duration(v)
	case int64:
		return time.Duration(v)
	case float64:
		return time.Duration(v)
	}
	return 0
}

// GetStringMapString returns the value associated with the key as a map of strings.
// If the value is not a map, it returns an empty map. All values in the map
// are converted to strings.
func (m *mapManager) GetStringMapString(key string) map[string]string {
	val := m.Get(key)
	if val == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	if v, ok := val.(map[string]interface{}); ok {
		for k, item := range v {
			result[k] = fmt.Sprintf("%v", item)
		}
	}
	return result
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (m *mapManager) GetStringSlice(key string) []string {
	val := m.Get(key)
	if val == nil {
		return []string{}
	}

	switch v := val.(type) {
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case []string:
		return v
	case string:
		if strings.Contains(v, ",") {
			parts := strings.Split(v, ",")
			result := make([]string, len(parts))
			for i, part := range parts {
				result[i] = strings.TrimSpace(part)
			}
			return result
		}
		return []string{v}
	}
	return []string{}
}

// IsSet checks if a key is set in the configuration.
func (m *mapManager) IsSet(key string) bool {
	return m.Get(key) != nil
}

// AllKeys returns all keys in the config, flattened with dot notation.
func (m *mapManager) AllKeys() []string {
	var keys []string
	collectKeys("", m.data, &keys)
	sort.Strings(keys)
	return keys
}

// collectKeys is a recursive helper for AllKeys.
func collectKeys(prefix string, data map[string]interface{}, keys *[]string) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]interface{}); ok {
			collectKeys(fullKey, nested, keys)
		} else {
			*keys = append(*keys, fullKey)
		}
	}
}

// Debug prints all configuration values to standard output.
func (m *mapManager) Debug() {
	fmt.Println("--- mflag configuration ---")
	keys := m.AllKeys()
	if len(keys) == 0 {
		fmt.Println("  (empty)")
		return
	}
	for _, key := range keys {
		value := m.Get(key)
		fmt.Printf("  %s: %v (%T)\n", key, value, value)
	}
	fmt.Println("---------------------------")
}

// convertMap recursively converts map[interface{}]interface{} to map[string]interface{}.
// The standard YAML library can unmarshal into the former, but we need the latter for
// structured access.
func convertMap(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			// The key is interface{}, convert it to string
			strKeyMap := make(map[string]interface{})
			for k3, v3 := range v2 {
				strKeyMap[fmt.Sprintf("%v", k3)] = v3
			}
			res[k] = convertMap(strKeyMap)
		case map[string]interface{}:
			res[k] = convertMap(v2)
		case []interface{}:
			// Handle slices that might contain maps
			res[k] = convertSlice(v2)
		default:
			res[k] = v
		}
	}
	return res
}

// convertSlice recursively converts slices containing maps
func convertSlice(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			strKeyMap := make(map[string]interface{})
			for k3, v3 := range v2 {
				strKeyMap[fmt.Sprintf("%v", k3)] = v3
			}
			result[i] = convertMap(strKeyMap)
		case map[string]interface{}:
			result[i] = convertMap(v2)
		case []interface{}:
			result[i] = convertSlice(v2)
		default:
			result[i] = v
		}
	}
	return result
}

// deepCopyMap creates a deep copy of a map.
func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	if original == nil {
		return nil
	}
	clone := make(map[string]interface{})
	for k, v := range original {
		clone[k] = deepCopyValue(v)
	}
	return clone
}

// deepCopyValue creates a deep copy of any value, handling nested maps and slices.
func deepCopyValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		return deepCopyMap(val)
	case []interface{}:
		clone := make([]interface{}, len(val))
		for i, item := range val {
			clone[i] = deepCopyValue(item)
		}
		return clone
	case []string:
		clone := make([]string, len(val))
		copy(clone, val)
		return clone
	default:
		// For basic types (string, int, bool, etc.), direct assignment is fine
		// as they are copied by value in Go
		return val
	}
}

// mergeMaps recursively merges two maps. Values in src overwrite values in dst.
func mergeMaps(dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = make(map[string]interface{})
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcOk := srcVal.(map[string]interface{})
			dstMap, dstOk := dstVal.(map[string]interface{})
			if srcOk && dstOk {
				dst[key] = mergeMaps(dstMap, srcMap)
				continue
			}
		}
		dst[key] = srcVal
	}
	return dst
}
