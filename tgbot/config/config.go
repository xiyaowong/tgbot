package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	ServerPort     = 8889
	HookPostURL    = ""
	ShowRawMessage = false
)

func init() {
	if cfg, err := ini.ShadowLoad("config.ini"); err == nil {
		sec, _ := cfg.GetSection(ini.DEFAULT_SECTION)
		ServerPort = sec.Key("ServerPort").MustInt(8889)
		HookPostURL = sec.Key("HookPostURL").MustString("")
		ShowRawMessage = sec.Key("ShowRawMessage").MustBool(false)
	}

	fmt.Println("当前配置: ")
	fmt.Printf("HTTP监听端口: %d\n", ServerPort)
	if HookPostURL != "" {
		fmt.Printf("数据上报地址: %s\n\n", HookPostURL)
	} else {
		fmt.Println("数据上报未开启")
	}
	if ShowRawMessage {
		fmt.Println("打印消息的原始数据")
	} else {
		fmt.Println("不打印消息的原始数据")
	}
	fmt.Println(" ")
}
