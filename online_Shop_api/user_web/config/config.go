package config

type UserSrvConfig struct {
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type EmailConfig struct {
	FromEmail string `mapstructure:"fromEmail"`
	SecretKey string `mapstructure:"secretKey"`
}

type RedisConfig struct {
	IP     string `mapstructure:"ip"`
	Port   int    `mapstructure:"port"`
	Expire int    `mapstructure:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name            string        `mapstructure:"name"`
	Port            int           `mapstructure:"port"`
	UserSrvInfo     UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo         JWTConfig     `mapstructure:"jwt"`
	RedisConfigInfo RedisConfig   `mapstructure:"redis"`
	EmailConfigInfo EmailConfig   `mapstructure:"email"`
	ConsulInfo      ConsulConfig  `mapstructure:"consul"`
}
