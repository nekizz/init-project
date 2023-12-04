package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Load loads configuration from local .env file
func Load(out interface{}, stage string) error {
	if err := PreloadLocalENV(stage); err != nil {
		return err
	}

	if err := env.Parse(out); err != nil {
		return err
	}

	return nil
}

// PreloadLocalENV reads .env* files and sets the values to os ENV
func PreloadLocalENV(stage string) error {
	basePath := ""
	if stage == "test" {
		basePath = "testdata/"
	}
	// // local config per stage
	// if stage != "" {
	// 	godotenv.Load(basePath + ".env." + stage + ".local")
	// }

	// local config
	godotenv.Load(basePath + ".env.local")

	// // per stage config
	// if stage != "" {
	// 	godotenv.Load(basePath + ".env." + stage)
	// }

	// default config
	return godotenv.Load(basePath + ".env")
}
