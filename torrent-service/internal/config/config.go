package config

type Config struct {
	App struct {
		Pame string
		Port int
	}
	Postgres struct {
		Host string
		Port string
		User string
		Password string
		DBName string
		maxCon string
	}
	Logger struct {
		Level string 
	}
}


func Load() (*Config, error) {
	return &Config{

	}
}