package main

import (
	"code.iirose.cn/soxft/proxy-center/client"
	"code.iirose.cn/soxft/proxy-center/config"
	"code.iirose.cn/soxft/proxy-center/server"
	"flag"
	"github.com/spf13/viper"
	"log"
)

var runMode string

func main() {
	flag.StringVar(&runMode, "mode", "server", "run mode")
	flag.Parse()

	config.Parse()

	switch runMode {
	case "clientTest":
		for i := 0; i < 10; i++ {
			c := client.NewClient(viper.GetString("Server.Address"), 2, 30)
			log.Println(c.GetProxy())
		}
		break
	default:
		server.Run()
	}

}
