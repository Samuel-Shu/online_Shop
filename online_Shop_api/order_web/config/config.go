package config

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}

type ServerConfig struct {
	Name             string         `mapstructure:"name" json:"name"`
	Host             string         `mapstructure:"host" json:"host"`
	Port             int            `mapstructure:"port" json:"port"`
	Tags             []string       `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo     GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	OrderSrvInfo     GoodsSrvConfig `mapstructure:"order_srv" json:"order_srv"`
	InventorySrvInfo GoodsSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	JWTInfo          JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo       ConsulConfig   `mapstructure:"consul" json:"consul"`
	OssInfo          OssConfig      `mapstructure:"oss" json:"oss"`
	AlipayInfo       AlipayConfig   `mapstructure:"alipay" json:"alipay"`
}

type OssConfig struct {
	OSSConfig     string `mapstructure:"oss_config" json:"oss_config"`
	AccessKey     string `mapstructure:"access_key" json:"access_key"`
	SecretKey     string `mapstructure:"secret_key" json:"secret_key"`
	VideoBucket   string `mapstructure:"video_bucket" json:"video_bucket"`
	PictureBucket string `mapstructure:"picture_bucket" json:"picture_bucket"`
	DomainVideo   string `mapstructure:"domain_video" json:"domain_video"`
	DomainPicture string `mapstructure:"domain_picture" json:"domain_picture"`
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
