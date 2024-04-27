package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"online_Shop_api/user_web/config"
	"online_Shop_api/user_web/global"
)

// GetEnvInfo 获取系统环境变量值
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig()  {
	debug := GetEnvInfo("ONLINE_SHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user_web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user_web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.NacosConfig); err != nil{
		panic(err)
	}

	zap.S().Infof("配置信息：%v", global.NacosConfig)

	//viper 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置信息产生变化：%s", e.Name)
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := v.Unmarshal(global.NacosConfig); err != nil{
			panic(err)
		}

		zap.S().Infof("配置信息：%v", global.NacosConfig)
	})

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port: uint64(global.NacosConfig.Port),

		},
	}

	cc := []constant.ClientConfig{
		{
			NamespaceId: global.NacosConfig.NameSpace,
			TimeoutMs: 5000,
			NotLoadCacheAtStart: true,
			LogDir: "../tmp/nacos/log",
			CacheDir: "../tmp/nacos/cache",
			LogLevel: "debug",
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig": cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	if err != nil {
		panic(err)
	}

	serverConfig := config.ServerConfig{}

	err = json.Unmarshal([]byte(content), &serverConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败：%s", err)
	}

	fmt.Println(serverConfig)
}
