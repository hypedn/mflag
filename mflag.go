// Package mflag provides integrated configuration management for Go applications,
// merging settings from default values, YAML files, and command-line flags.
package mflag

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	ErrInitFailed = errors.New("mflag: Init failed")
)

var (
	defaults    = newManager()
	config      = newManager()
	finalConfig = newManager()
	parsed      = false
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}
}

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

// GetInt8 returns the value associated with the key as an int8.
// Must be called after Parse.
func GetInt8(key string) int8 {
	mustBeParsed()
	return finalConfig.GetInt8(key)
}

// GetInt16 returns the value associated with the key as an int16.
// Must be called after Parse.
func GetInt16(key string) int16 {
	mustBeParsed()
	return finalConfig.GetInt16(key)
}

// GetInt32 returns the value associated with the key as an int32.
// Must be called after Parse.
func GetInt32(key string) int32 {
	mustBeParsed()
	return finalConfig.GetInt32(key)
}

// GetInt64 returns the value associated with the key as an int64.
// Must be called after Parse.
func GetInt64(key string) int64 {
	mustBeParsed()
	return finalConfig.GetInt64(key)
}

// GetUint returns the value associated with the key as a uint.
// Must be called after Parse.
func GetUint(key string) uint {
	mustBeParsed()
	return finalConfig.GetUint(key)
}

// GetUint8 returns the value associated with the key as a uint8.
// Must be called after Parse.
func GetUint8(key string) uint8 {
	mustBeParsed()
	return finalConfig.GetUint8(key)
}

// GetUint16 returns the value associated with the key as a uint16.
// Must be called after Parse.
func GetUint16(key string) uint16 {
	mustBeParsed()
	return finalConfig.GetUint16(key)
}

// GetUint32 returns the value associated with the key as a uint32.
// Must be called after Parse.
func GetUint32(key string) uint32 {
	mustBeParsed()
	return finalConfig.GetUint32(key)
}

// GetUint64 returns the value associated with the key as a uint64.
// Must be called after Parse.
func GetUint64(key string) uint64 {
	mustBeParsed()
	return finalConfig.GetUint64(key)
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

// GetDuration returns the value associated with the key as a time.Duration.
// Must be called after Parse.
func GetDuration(key string) time.Duration {
	mustBeParsed()
	return finalConfig.GetDuration(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
// Must be called after Parse.
func GetStringMapString(key string) map[string]string {
	mustBeParsed()
	return finalConfig.GetStringMapString(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
// Must be called after Parse.
func GetStringSlice(key string) []string {
	mustBeParsed()
	return finalConfig.GetStringSlice(key)
}

// GetStringSet returns the string slice value associated with a key as a map[string]bool (a set).
// This is useful for efficiently checking for the existence of an item in a list, like a feature flag.
// Must be called after Parse.
func GetStringSet(key string) map[string]bool {
	mustBeParsed()
	l := finalConfig.GetStringSlice(key)
	m := make(map[string]bool, len(l))
	for _, item := range l {
		m[item] = true
	}
	return m
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

// populateFlagSet dynamically creates flags for all known keys on a given flag set.
// It returns a slice of errors for any invalid default values encountered.
func populateFlagSet(fs *flag.FlagSet) []error {
	allKeys := finalConfig.AllKeys()
	var errs []error
	for _, key := range allKeys {
		value := finalConfig.Get(key)
		usage := fmt.Sprintf("override configuration for '%s'", key)

		switch v := value.(type) {
		case bool:
			fs.Bool(key, v, usage)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			isUint := false
			if dv := defaults.Get(key); dv != nil {
				switch dv.(type) {
				case uint, uint8, uint16, uint32, uint64:
					isUint = true
				}
			}

			if isUint {
				val, err := castToUint64(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("invalid value for uint flag %q: %w", key, err))
					continue
				}
				fs.Uint64(key, val, usage)
			} else {
				val, err := castToInt(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("invalid default for flag %q: %w", key, err))
					continue
				}
				fs.Int(key, val, usage)
			}
		case float64:
			val, err := castToFloat64(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid default for flag %q: %w", key, err))
				continue
			}
			fs.Float64(key, val, usage)
		case time.Duration:
			val, err := castToDuration(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid default for flag %q: %w", key, err))
				continue
			}
			fs.Duration(key, val, usage)
		default: // string, slices, maps, etc.
			fs.String(key, finalConfig.GetString(key), usage)
		}
	}
	return errs
}

// Parse parses command-line arguments and merges all configuration sources.
// It MUST be called after setting defaults and calling Init. It dynamically creates
// command-line flags for all known configuration keys.
// Precedence: Flags > Config File > Defaults.
func Parse() {
	// 1. Start with a copy of the defaults.
	finalConfig = defaults.Clone()

	// 2. Merge config file values on top of defaults.
	finalConfig.Merge(config)

	// 3. Populate the global command-line flag set.
	errs := populateFlagSet(flag.CommandLine)

	if len(errs) > 0 {
		// Mimic the behavior of the standard flag package on error.
		fmt.Fprintln(flag.CommandLine.Output(), errors.Join(errs...))
		os.Exit(1)
	}

	flag.Parse()

	// 4. Overwrite finalConfig with values from flags that were explicitly set
	//    on the command line. This gives them the highest precedence.
	flag.Visit(func(f *flag.Flag) {
		getter := f.Value.(flag.Getter)
		finalConfig.SetValue(f.Name, getter.Get())
	})
	parsed = true
}

// ParseWithError is similar to Parse but returns an error on failure.
// This allows for more granular error handling.
// Note: This function creates its own temporary flag set and does not parse
// flags defined globally via the standard `flag` package.
func ParseWithError() error {
	// 1. Start with a copy of the defaults.
	finalConfig = defaults.Clone()

	// 2. Merge config file values on top of defaults.
	finalConfig.Merge(config)

	// 3. Dynamically create flags for all known keys on a temporary flag set.
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// 4. Populate the temporary flag set.
	if errs := populateFlagSet(fs); len(errs) > 0 {
		return errors.Join(errs...)
	}

	// 5. Parse the command-line arguments.
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	fs.Visit(func(f *flag.Flag) {
		getter := f.Value.(flag.Getter)
		finalConfig.SetValue(f.Name, getter.Get())
	})
	parsed = true
	return nil
}

// castToInt converts an interface{} to an int, handling common numeric types.
func castToInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
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

// castToUint64 converts an interface{} to a uint64.
func castToUint64(v interface{}) (uint64, error) {
	switch val := v.(type) {
	case uint64:
		return val, nil
	case uint:
		return uint64(val), nil
	case uint8:
		return uint64(val), nil
	case uint16:
		return uint64(val), nil
	case uint32:
		return uint64(val), nil
	case int:
		if val < 0 {
			return 0, fmt.Errorf("cannot cast negative int %d to uint64", val)
		}
		return uint64(val), nil
	case int64:
		if val < 0 {
			return 0, fmt.Errorf("cannot cast negative int64 %d to uint64", val)
		}
		return uint64(val), nil
	case float64:
		if val < 0 {
			return 0, fmt.Errorf("cannot cast negative float64 %f to uint64", val)
		}
		return uint64(val), nil
	case string:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot cast string %q to uint64: %w", val, err)
		}
		return u, nil
	}
	return 0, fmt.Errorf("cannot cast type %T to uint64", v)
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

// castToDuration converts an interface{} to a time.Duration.
func castToDuration(v interface{}) (time.Duration, error) {
	switch val := v.(type) {
	case time.Duration:
		return val, nil
	case string:
		d, err := time.ParseDuration(val)
		if err != nil {
			return 0, fmt.Errorf("cannot cast string %q to time.Duration: %w", val, err)
		}
		return d, nil
	case int:
		return time.Duration(val), nil
	case int64:
		return time.Duration(val), nil
	case float64:
		return time.Duration(val), nil
	}
	return 0, fmt.Errorf("cannot cast type %T to time.Duration", v)
}
