package server

import (
	"code.iirose.cn/soxft/proxy-center/server/proxy"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gopkg.in/eapache/queue.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runServer() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/getProxy/:timeout/:endtime", func(c *gin.Context) {
		timeout := c.Param("timeout")
		endtime := c.Param("endtime")
		if timeout == "" || endtime == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "param error",
				"data":    gin.H{},
			})
			return
		}

		// 获取代理
		ctx, cancel := context.WithTimeout(c, time.Second*10)
		defer cancel()

		// 获取代理
		prox, err := proxy.Get(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    gin.H{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "pong",
			"data": gin.H{
				"addr":     prox.Addr,
				"city":     prox.City,
				"end_time": prox.EndTime,
			},
		})
	})
	if err := r.Run(viper.GetString("Server.Address")); err != nil {
		log.Fatalf("[ERROR] Server Run Error: %s", err)
	}
}

func Run() {
	proxy.Que = queue.New()

	go proxy.MainCron()

	// clear cron
	c := cron.New()
	_, _ = c.AddFunc("@every 2m", proxy.ClearProxy)
	//_, _ = c.AddFunc("@every 1s", Get)
	c.Start()

	// start server
	go runServer()

	// 监听 ctrl + c 信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigs:
		log.Println("[INFO] Exit")
		os.Exit(0)
	}
}
