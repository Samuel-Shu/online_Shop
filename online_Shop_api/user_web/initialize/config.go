package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	if err := v.Unmarshal(global.ServerConfig); err != nil{
		panic(err)
	}

	zap.S().Infof("配置信息：%v", global.ServerConfig)

	//viper 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置信息产生变化：%s", e.Name)
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := v.Unmarshal(global.ServerConfig); err != nil{
			panic(err)
		}

		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})
}
