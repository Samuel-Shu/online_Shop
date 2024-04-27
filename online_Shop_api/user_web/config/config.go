package config

type UserSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type EmailConfig struct {
	FromEmail string `mapstructure:"fromEmail" json:"fromEmail"`
	SecretKey string `mapstructure:"secretKey" json:"secretKey"`
}

type RedisConfig struct {
	IP     string `mapstructure:"ip" json:"ip"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name            string        `mapstructure:"name" json:"name"`
	Port            int           `mapstructure:"port" json:"port"`
	UserSrvInfo     UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo         JWTConfig     `mapstructure:"jwt" json:"jwt"`
	RedisConfigInfo RedisConfig   `mapstructure:"redis" json:"redis"`
	EmailConfigInfo EmailConfig   `mapstructure:"email" json:"email"`
	ConsulInfo      ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	NameSpace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
