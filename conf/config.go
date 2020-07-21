package conf

import (
	"strings"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config {
		Name: cfg,
	}

	if err := c.initConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) initConfig() error  {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("conf")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO_DISTRIBUTE_STORAGE")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
