package config

// Configuration holds data necessery for configuring application
type Configuration struct {
	Stage               string   `env:"UP_STAGE"`
	Host                string   `env:"HOST"`
	Port                int      `env:"PORT"`
	ReadTimeout         int      `env:"READ_TIMEOUT"`
	WriteTimeout        int      `env:"WRITE_TIMEOUT"`
	AllowOrigins        []string `env:"ALLOW_ORIGINS"`
	Debug               bool     `env:"DEBUG"`
	DbLog               bool     `env:"DB_LOG"`
	DbPsn               string   `env:"DB_PSN"`
	JwtSecret           string   `env:"JWT_SECRET"`
	JwtDuration         int      `env:"JWT_DURATION"`
	JwtSigningAlgorithm string   `env:"JWT_SIGNINGALGORITHM"`
}

// Load returns Configuration struct
func Load() (*Configuration, error) {

	return nil, nil
}
