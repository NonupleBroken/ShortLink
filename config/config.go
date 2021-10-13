package config

import "github.com/BurntSushi/toml"

type Config struct {
	Server ServerConfig `toml:"server"`
	Redis  RedisConfig  `toml:"redis"`
	Log    LogConfig    `toml:"log"`
}

type ServerConfig struct {
	Addr     string `toml:"addr"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DB       int    `toml:"db"`
	Password string `toml:"password"`
}

type LogConfig struct {
	LogPath       string `toml:"log_path"`
	ServerLogName string `toml:"server_log_name"`
	WebLogName    string `toml:"web_log_name"`
}

func InitConfig(configFile string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(configFile, &config)
	return config, err
}
