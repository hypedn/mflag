package mflag

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestPrecedenceOrder(t *testing.T) {
	testReset(t)

	SetDefault("port", 1111)
	SetDefault("db.user", "default_user")
	SetDefault("feature.new", false)
	SetDefault("enabled", false) // Add default for flag-only key
	configFileContent := `
port: 2222
db:
  host: "config.host"
  user: "config_user"
features:
  - dark_mode
`
	configPath := createTempYAML(t, configFileContent)
	if err := Init(configPath); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	os.Args = []string{
		"test_app",
		"--port=3333",
		"--db.host=flag.host",
		"--enabled",
	}

	Parse()

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		source   string
	}{
		{"port", 3333, GetInt("port"), "Command-line flag"},
		{"db.host", "flag.host", GetString("db.host"), "Command-line flag"},
		{"enabled", true, GetBool("enabled"), "Command-line flag"},
		{"db.user", "config_user", GetString("db.user"), "Config file"},
		{"features", []string{"dark_mode"}, GetStringSlice("features"), "Config file"},
		{"feature.new", false, GetBool("feature.new"), "SetDefault value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.actual, tt.expected) {
				t.Errorf("Expected %s to be %v (from %s), but got %v", tt.name, tt.expected, tt.source, tt.actual)
			}
		})
	}

	if !IsSet("port") {
		t.Error("Expected 'port' to be set")
	}
	if IsSet("nonexistent.key") {
		t.Error("Expected 'nonexistent.key' to not be set")
	}
}

func TestGetBeforeParsePanic(t *testing.T) {
	testReset(t)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			if !strings.Contains(fmt.Sprintf("%v", r), "Parse() must be called before using Get* functions") {
				t.Errorf("Unexpected panic message: %v", r)
			}
		}
	}()
	// This should panic
	GetString("any_key")
}

func TestInitNonExistentFile(t *testing.T) {
	testReset(t)

	err := Init("non-existent-file-for-test.yaml")
	if err != nil {
		t.Errorf("Init() with non-existent file should not return an error, but got: %v", err)
	}
}

func TestInit_BadYAML(t *testing.T) {
	testReset(t)

	// This YAML is invalid because of the unclosed quote.
	badYAML := `key: "value`
	path := createTempYAML(t, badYAML)

	err := Init(path)
	if err == nil {
		t.Fatal("Init() should have failed with bad YAML syntax, but it did not")
	}

	if !strings.Contains(err.Error(), "failed to parse yaml") {
		t.Errorf("Expected error message to indicate a YAML parsing failure, but got: %v", err)
	}
}

func TestInit_FilePermissionError(t *testing.T) {
	testReset(t)

	path := createTempYAML(t, "key: value")
	// Make the file unreadable. On Windows, this might not prevent reading,
	// but it's the standard POSIX way to test this.
	if err := os.Chmod(path, 0000); err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0644) })

	err := Init(path)
	if err == nil {
		t.Fatal("Init() should have failed with a file permission error, but it did not")
	}

	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected error message to indicate a file read failure, but got: %v", err)
	}
}

func TestStringSliceHandling(t *testing.T) {
	testReset(t)

	SetDefault("features_csv", "three, four")
	Parse()

	expected := []string{"three", "four"}
	actual := GetStringSlice("features_csv")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestTypeConversions(t *testing.T) {
	testReset(t)

	SetDefault("int_from_string", "123")
	SetDefault("bool_from_string", "true")
	SetDefault("float_from_int", 123)

	Parse()

	if val := GetInt("int_from_string"); val != 123 {
		t.Errorf("Expected GetInt to convert string '123' to 123, got %d", val)
	}
	if val := GetBool("bool_from_string"); val != true {
		t.Errorf("Expected GetBool to convert string 'true' to true, got %t", val)
	}
	if val := GetFloat64("float_from_int"); val != 123.0 {
		t.Errorf("Expected GetFloat64 to convert int 123 to 123.0, got %f", val)
	}
}

func TestIntegerTypeConversions(t *testing.T) {
	testReset(t)

	SetDefault("i8", int8(8))
	SetDefault("i16", int16(16))
	SetDefault("i32", int32(32))
	SetDefault("i64", int64(64))
	SetDefault("i", int(12))
	SetDefault("u8", uint8(8))
	SetDefault("u16", uint16(16))
	SetDefault("u32", uint32(32))
	SetDefault("u64", uint64(64))
	SetDefault("u", uint(12))
	SetDefault("i_from_s", "123")
	SetDefault("u_from_s", "456")
	SetDefault("neg_i", -10)

	Parse()

	if val := GetInt8("i8"); val != 8 {
		t.Errorf("GetInt8 failed, expected 8, got %d", val)
	}
	if val := GetInt16("i16"); val != 16 {
		t.Errorf("GetInt16 failed, expected 16, got %d", val)
	}
	if val := GetInt32("i32"); val != 32 {
		t.Errorf("GetInt32 failed, expected 32, got %d", val)
	}
	if val := GetInt64("i64"); val != 64 {
		t.Errorf("GetInt64 failed, expected 64, got %d", val)
	}
	if val := GetInt("i"); val != 12 {
		t.Errorf("GetInt failed, expected 12, got %d", val)
	}

	if val := GetUint8("u8"); val != 8 {
		t.Errorf("GetUint8 failed, expected 8, got %d", val)
	}
	if val := GetUint16("u16"); val != 16 {
		t.Errorf("GetUint16 failed, expected 16, got %d", val)
	}
	if val := GetUint32("u32"); val != 32 {
		t.Errorf("GetUint32 failed, expected 32, got %d", val)
	}
	if val := GetUint64("u64"); val != 64 {
		t.Errorf("GetUint64 failed, expected 64, got %d", val)
	}
	if val := GetUint("u"); val != 12 {
		t.Errorf("GetUint failed, expected 12, got %d", val)
	}
	if val := GetUint("neg_i"); val != 0 {
		t.Errorf("GetUint from negative int failed, expected 0, got %d", val)
	}
}

func TestGetStringMapString(t *testing.T) {
	testReset(t)

	// 1. Test with YAML config
	configFileContent := `
database:
  host: "db.host.com"
  port: 5432
  enabled: true
not_a_map: "some_string"
`
	configPath := createTempYAML(t, configFileContent)
	if err := Init(configPath); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 2. Test with SetDefault
	SetDefault("user_prefs", map[string]interface{}{
		"theme":         "dark",
		"notifications": true,
	})

	Parse()

	// Test case 1: Get from YAML
	dbSettings := GetStringMapString("database")
	expectedDbSettings := map[string]string{
		"host":    "db.host.com",
		"port":    "5432",
		"enabled": "true",
	}
	if !reflect.DeepEqual(dbSettings, expectedDbSettings) {
		t.Errorf("GetStringMapString('database') failed. Expected %v, got %v", expectedDbSettings, dbSettings)
	}

	// Test case 2: Get from SetDefault
	userPrefs := GetStringMapString("user_prefs")
	expectedUserPrefs := map[string]string{
		"theme":         "dark",
		"notifications": "true",
	}
	if !reflect.DeepEqual(userPrefs, expectedUserPrefs) {
		t.Errorf("GetStringMapString('user_prefs') failed. Expected %v, got %v", expectedUserPrefs, userPrefs)
	}

	// Test case 3: Get non-existent key
	if m := GetStringMapString("non.existent"); len(m) != 0 {
		t.Errorf("Expected empty map for non-existent key, got %v", m)
	}

	// Test case 4: Get key with non-map value
	if m := GetStringMapString("not_a_map"); len(m) != 0 {
		t.Errorf("Expected empty map for non-map key, got %v", m)
	}
}

func Example() {
	defer func(oldArgs []string) {
		os.Args = oldArgs
		resetGlobals()
	}(os.Args)
	resetGlobals()

	os.Args = []string{"cmd", "--host=flag.host", "--debug=false"}

	SetDefault("host", "default.host")
	SetDefault("port", 8080)
	SetDefault("debug", true)
	// Add default for a key that will not be overridden by config or flags
	SetDefault("timeout", 5)

	config.SetValue("host", "config.host")
	config.SetValue("port", 9090)

	Parse()

	fmt.Printf("Host: %s (from flag)\n", GetString("host"))       // Highest precedence
	fmt.Printf("Port: %d (from config)\n", GetInt("port"))        // Middle precedence
	fmt.Printf("Debug: %t (from flag)\n", GetBool("debug"))       // Flag overriding default
	fmt.Printf("Timeout: %d (from default)\n", GetInt("timeout")) // Lowest precedence

	// Output:
	// Host: flag.host (from flag)
	// Port: 9090 (from config)
	// Debug: false (from flag)
	// Timeout: 5 (from default)
}

func TestParseErrorHandling(t *testing.T) {
	// This test checks if ParseWithError returns an error when a config value is incompatible
	// with the type of the default value. Here, the default is a uint, but the
	// config file provides a negative number, which is invalid.
	testReset(t)
	SetDefault("my-uint", uint(10))
	configFileContent := `my-uint: -5`
	configPath := createTempYAML(t, configFileContent)
	if err := Init(configPath); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	err := ParseWithError()
	if err == nil {
		t.Fatal("ParseWithError() should have failed with an incompatible config value, but it did not")
	}
	if !strings.Contains(err.Error(), "invalid value for uint flag") {
		t.Errorf("Expected error about invalid uint value, got: %v", err)
	}
}

// resetGlobals resets all package-level state variables and the default flag set.
// This is the core reset logic, callable from both tests and examples.
func resetGlobals() {
	defaults = newManager()
	config = newManager()
	finalConfig = newManager()
	parsed = false

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// testReset is a helper for Test* functions. It resets global state and
// mocks os.Args to prevent the test runner's flags from being parsed.
// It uses t.Cleanup to restore os.Args automatically.
func testReset(t *testing.T) {
	t.Helper()
	resetGlobals()

	oldArgs := os.Args
	os.Args = []string{"test"}
	t.Cleanup(func() { os.Args = oldArgs })
}

func createTempYAML(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Fatalf("Failed to remove tmpfile %v, error: %v", tmpfile.Name(), err)
		}
	})
	return tmpfile.Name()
}
