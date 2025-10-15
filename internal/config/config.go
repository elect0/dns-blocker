package config

import "github.com/spf13/viper"

type Config struct {
	ListenAddress  string            `mapstructure:"listen_address"`
	UpstreamServer string            `mapstructure:"upstream_server"`
	BlocklistPath  string            `mapstructure:"blocklist_path"`
	CustomRecords  map[string]string `mapstructure:"custom_records"`
	Logging        struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
