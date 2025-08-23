package config

import "github.com/spf13/viper"

type option struct {
	ConfigFolder []string
	ConfigType   string
	ConfigFile   string
}

type Option func(*option)

func LoadConfig(opts ...Option) Config {
	opt := &option{
		ConfigFolder: getDefaultConfigFolder(),
		ConfigFile:   getDefaultConfigFile(),
		ConfigType:   getDefaultConfigType(),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	for _, folder := range opt.ConfigFolder {
		viper.AddConfigPath(folder)
	}

	viper.SetConfigName(opt.ConfigFile)
	viper.SetConfigType(opt.ConfigType)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func getDefaultConfigFolder() []string {
	return []string{"./files/config"}
}

func getDefaultConfigFile() string {
	return "config"
}

func getDefaultConfigType() string {
	return "yaml"
}

func WithConfigFolder(folder []string) Option {
	return func(o *option) {
		o.ConfigFolder = folder
	}
}

func WithConfigFile(file string) Option {
	return func(o *option) {
		o.ConfigFile = file
	}
}

func WithConfigType(configType string) Option {
	return func(o *option) {
		o.ConfigType = configType
	}
}
