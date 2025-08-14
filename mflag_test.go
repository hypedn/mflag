package mflag

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func testReset() {
	defaults = newManager()
	config = newManager()
	finalConfig = newManager()
	parsed = false

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
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

func TestPrecedenceOrder(t *testing.T) {
	testReset()

	SetDefault("port", 1111)
	SetDefault("db.user", "default_user")
	SetDefault("feature.new", false)

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

	oldArgs := os.Args
	os.Args = []string{
		"test_app",
		"-port=3333",
		"--db.host=flag.host",
		"-enabled",
	}
	defer func() { os.Args = oldArgs }()

	_ = Int("port", 9999, "Server port")
	_ = String("db.host", "fallback.host", "DB host")
	_ = String("db.user", "fallback_user", "DB user")
	_ = Bool("enabled", false, "Enable flag")
	_ = Bool("feature.new", true, "A new feature flag")
	_ = String("log.level", "info", "Log level")

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
		{"log.level", "info", GetString("log.level"), "Flag's hardcoded default"},
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
	testReset()
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
	testReset()
	err := Init("non-existent-file-for-test.yaml")
	if err != nil {
		t.Errorf("Init() with non-existent file should not return an error, but got: %v", err)
	}
}

func TestStringSliceHandling(t *testing.T) {
	testReset()
	oldArgs := os.Args
	os.Args = []string{"test"}
	defer func() { os.Args = oldArgs }()

	SetDefault("features_csv", "three, four")
	Parse()

	expected := []string{"three", "four"}
	actual := GetStringSlice("features_csv")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestTypeConversions(t *testing.T) {
	testReset()
	oldArgs := os.Args
	os.Args = []string{"test"}
	defer func() { os.Args = oldArgs }()

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

func Example() {
	defer func(oldArgs []string) {
		os.Args = oldArgs
		testReset()
	}(os.Args)
	testReset()
	os.Args = []string{"cmd", "--host=flag.host", "--debug=false"}

	SetDefault("host", "default.host")
	SetDefault("port", 8080)
	SetDefault("debug", true)

	config.SetValue("host", "config.host")
	config.SetValue("port", 9090)

	_ = String("host", "fallback.host", "Hostname")
	_ = Int("port", 3000, "Port number")
	_ = Bool("debug", false, "Enable debug mode")
	_ = Int("timeout", 5, "Request timeout")

	Parse()

	fmt.Printf("Host: %s (from flag)\n", GetString("host"))
	fmt.Printf("Port: %d (from config)\n", GetInt("port"))
	fmt.Printf("Debug: %t (from flag)\n", GetBool("debug"))
	fmt.Printf("Timeout: %d (from flag's hardcoded default)\n", GetInt("timeout"))

	// Output:
	// Host: flag.host (from flag)
	// Port: 9090 (from config)
	// Debug: false (from flag)
	// Timeout: 5 (from flag's hardcoded default)
}
