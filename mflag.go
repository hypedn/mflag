// Package mflag provides a simple configuration library for Go applications,
// inspired by the standard `flag` package. It merges configuration from a YAML
// file with command-line flags, providing a clear precedence order.
//
// # Usage Pattern
//
// The intended usage follows a strict three-step process in your main function:
//
//  1. mflag.Init("path/to/configmap.yaml"): Load configuration from a file. This
//     step is optional; if skipped, configuration will only come from flags and
//     defaults.
//
//  2. Define flags using mflag.String(), mflag.Int(), etc. These functions
//     wrap the standard library's flag functions. If a value is present in the
//     loaded config file, it will be used as the default for the flag.
//
//  3. mflag.Parse(): Parses the command-line flags. Values provided on the
//     command line will override values from the config file.
//
// # Precedence Order
//
// 1. Command-line flags (highest precedence)
// 2. Values from the YAML configuration file
// 3. Default values defined in code (lowest precedence)
//
// All Get* functions must be called after Parse().
package mflag

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

var (
	ErrNotParsed   = errors.New("mflag: Parse() must be called before using Get* functions")
	ErrKeyNotFound = errors.New("mflag: key not found")
)

var (
	defaults    = newManager()
	config      = newManager()
	finalConfig = newManager()
	parsed      = false
)

// SetDefault sets a default value for a key.
// Defaults have the lowest precedence and are overridden by config files and flags.
// It should be called before Init and Parse.
func SetDefault(key string, value interface{}) {
	defaults.SetValue(key, value)
}

// Init loads configuration from a YAML file at the given path. It should be
// called after setting defaults and before parsing flags.
func Init(filename string) error {
	return config.LoadFile(filename)
}

// mustBeParsed checks if Parse() has been called and panics if not.
// This follows the same pattern as the standard flag package.
func mustBeParsed() {
	if !parsed {
		panic("mflag: Parse() must be called before using Get* functions")
	}
}

// GetString returns the value associated with the key as a string.
// It returns the final value after merging defaults, config file, and flags.
// Must be called after Parse.
func GetString(key string) string {
	mustBeParsed()
	return finalConfig.GetString(key)
}

// GetInt returns the value associated with the key as an integer.
// Must be called after Parse.
func GetInt(key string) int {
	mustBeParsed()
	return finalConfig.GetInt(key)
}

// GetBool returns the value associated with the key as a boolean.
// Must be called after Parse.
func GetBool(key string) bool {
	mustBeParsed()
	return finalConfig.GetBool(key)
}

// GetFloat64 returns the value associated with the key as a float64.
// Must be called after Parse.
func GetFloat64(key string) float64 {
	mustBeParsed()
	return finalConfig.GetFloat64(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
// Must be called after Parse.
func GetStringSlice(key string) []string {
	mustBeParsed()
	return finalConfig.GetStringSlice(key)
}

// IsSet checks if a key is set in the configuration.
// Must be called after Parse.
func IsSet(key string) bool {
	mustBeParsed()
	return finalConfig.IsSet(key)
}

// AllKeys returns all keys in the config, flattened with dot notation.
// Must be called after Parse.
func AllKeys() []string {
	mustBeParsed()
	return finalConfig.AllKeys()
}

// Debug prints all configuration values to standard output.
// Must be called after Parse.
func Debug() {
	mustBeParsed()
	finalConfig.Debug()
}

// --- Flag integration: https://pkg.go.dev/flag#pkg-index ---

// String defines a string flag. The default value is overridden by a value from the
// configuration file if one exists for the given key.
func String(name string, value string, usage string) *string {
	if val := getFlagDefault(name); val != nil {
		value = fmt.Sprintf("%v", val)
	}
	return flag.String(name, value, usage)
}

// Int defines an int flag. The default value is overridden by a value from the
// configuration file if one exists for the given key.
func Int(name string, value int, usage string) *int {
	if val := getFlagDefault(name); val != nil {
		newValue, err := castToInt(val)
		if err != nil {
			panic(fmt.Sprintf("mflag: invalid value for key %q in config/defaults: %v", name, err))
		}
		value = newValue
	}
	return flag.Int(name, value, usage)
}

// Bool defines a bool flag. The default value is overridden by a value from the
// configuration file if one exists for the given key.
func Bool(name string, value bool, usage string) *bool {
	if val := getFlagDefault(name); val != nil {
		newValue, err := castToBool(val)
		if err != nil {
			panic(fmt.Sprintf("mflag: invalid value for key %q in config/defaults: %v", name, err))
		}
		value = newValue
	}
	return flag.Bool(name, value, usage)
}

// Float64 defines a float64 flag. The default value is overridden by a value from the
// configuration file if one exists for the given key.
func Float64(name string, value float64, usage string) *float64 {
	if val := getFlagDefault(name); val != nil {
		newValue, err := castToFloat64(val)
		if err != nil {
			panic(fmt.Sprintf("mflag: invalid value for key %q in config/defaults: %v", name, err))
		}
		value = newValue
	}
	return flag.Float64(name, value, usage)
}

// getFlagDefault determines the default value for a flag based on the precedence:
// config file > code defaults (from SetDefault).
func getFlagDefault(name string) interface{} {
	if config.IsSet(name) {
		return config.Get(name)
	}
	return defaults.Get(name) // returns nil if not set
}

// Parse parses command-line arguments and merges all configuration sources.
// It MUST be called after all flags are defined. After Parse is called, the Get*
// functions will return the final, merged configuration values.
// Precedence: Flags > Config File > Defaults.
func Parse() {
	// 1. Start with a copy of the defaults.
	finalConfig = defaults.Clone()

	// 2. Merge config file values on top of defaults.
	finalConfig.Merge(config)

	// 3. Parse flags. The flag package will use values from the config file as
	//    defaults if they were present when the flag was defined.
	flag.Parse()

	// 4. Merge all flag values into the final config. This gives flags the
	//    highest precedence. We use VisitAll to get all flags, not just those
	//    set on the command line, to capture their default values as well.
	flag.VisitAll(func(f *flag.Flag) {
		// Use the flag.Getter interface to get typed values from flags
		if getter, ok := f.Value.(flag.Getter); ok {
			finalConfig.SetValue(f.Name, getter.Get())
		}
	})

	parsed = true
}

// Parsed reports whether the command-line flags have been parsed.
// This mirrors flag.Parsed() for consistency.
func Parsed() bool {
	return parsed
}

// castToInt converts an interface{} to an int, handling common numeric types.
func castToInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int64: // YAML can unmarshal to int64
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("cannot cast string %q to int: %w", val, err)
		}
		return i, nil
	}
	return 0, fmt.Errorf("cannot cast type %T to int", v)
}

// castToBool converts an interface{} to a bool.
func castToBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("cannot cast string %q to bool: %w", val, err)
		}
		return b, nil
	}
	return false, fmt.Errorf("cannot cast type %T to bool", v)
}

// castToFloat64 converts an interface{} to a float64.
func castToFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot cast string %q to float64: %w", val, err)
		}
		return f, nil
	}
	return 0.0, fmt.Errorf("cannot cast type %T to float64", v)
}
