package main

import (
	"log"

	"github.com/hypedn/mflag"
)

const (
	debug   = "debug"
	appPort = "app_port"

	dbKey  = "database"
	dbHost = "host"
	dbPort = "port"
	dbUser = "user"

	features    = "features"
	darkMode    = "dark_mode"
	betaTesting = "beta_testing"
)

type FeatureFlags struct {
	UseDarkMode    bool
	UseBetaTesting bool
}

type AppSettings struct {
	Debug    bool
	AppPort  int
	Database struct {
		Host string
		Port string
		User string
	}
	Flags FeatureFlags
}

func defaults() {
	mflag.SetDefault(debug, true)
	mflag.SetDefault(appPort, 8080)

	mflag.SetDefault(dbKey, map[string]interface{}{
		dbHost: "localhost",
		dbPort: 5432,
		dbUser: "default_user",
	})

	mflag.SetDefault(features, []string{darkMode, betaTesting})
}

func GetSettings() AppSettings {
	dbSettings := mflag.GetStringMapString(dbKey)
	featureFlags := mflag.GetStringSet(features)

	return AppSettings{
		Debug:   mflag.GetBool(debug),
		AppPort: mflag.GetInt(appPort),
		Database: struct {
			Host string
			Port string
			User string
		}{
			Host: dbSettings[dbHost],
			Port: dbSettings[dbPort],
			User: dbSettings[dbUser],
		},
		Flags: FeatureFlags{
			UseDarkMode:    featureFlags[darkMode],
			UseBetaTesting: featureFlags[betaTesting],
		},
	}
}

func main() {
	defaults()
	if err := mflag.Init("configmap.yaml"); err != nil {
		log.Fatal(err)
	}
	mflag.Parse()

	config := GetSettings()
	if config.Debug {
		mflag.Debug()
	}

	if config.Flags.UseBetaTesting {
		log.Println("✅ Beta testing is enabled.")
	}
	if config.Flags.UseDarkMode {
		log.Println("✅ Dark mode is enabled.")
	}
}
