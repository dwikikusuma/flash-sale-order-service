package config

type Config struct {
	App      App           `mapstructure:"app" validate:"required"`
	DB       DB            `mapstructure:"db" validate:"required"`
	Redis    Redis         `mapstructure:"redis" validate:"required"`
	Secret   SecreteConfig `mapstructure:"secret" validate:"required"`
	Services Services      `mapstructure:"services" validate:"required"`
	Kafka    Kafka         `mapstructure:"kafka" validate:"required"`
}

type App struct {
	Port string `mapstructure:"port" validate:"required"`
}

type DB struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
	NameS1   string `mapstructure:"nameS1" validate:"required"` // For sharding, e.g., db_name-s1
	NameS2   string `mapstructure:"nameS2" validate:"required"` // For sharding, e.g., db_name-s2
}

type SecreteConfig struct {
	JWTSecret string `mapstructure:"jwt_secret" validate:"required"`
}

type Redis struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Password string `mapstructure:"password"`
}

type Services struct {
	Product string `mapstructure:"product" validate:"required"`
	Pricing string `mapstructure:"pricing" validate:"required"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers" validate:"required"`
	Topic   string   `mapstructure:"topic" validate:"required"`
}
