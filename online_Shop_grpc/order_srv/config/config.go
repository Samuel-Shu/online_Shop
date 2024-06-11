package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name            string       `mapstructure:"name" json:"name"`
	MysqlConfigInfo MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo      ConsulConfig `mapstructure:"consul" json:"consul"`
	RedisConfigInfo RedisConfig  `mapstructure:"redis" json:"redis"`

	//商品微服务的配置信息
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`

	//库存微服务的配置信息
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
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

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
