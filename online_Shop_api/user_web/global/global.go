package global

import (
	ut "github.com/go-playground/universal-translator"
	"online_Shop_api/user_web/config"
)

var (
	Trans ut.Translator

	ServerConfig = &config.ServerConfig{}
)
