package main

import (
	"context"
	"go-networking/config"
	"go-networking/network"
	"go-networking/router"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	tcpServer, _ := network.NewTcpServer()
	tcpServer.Start()
	connProcessor := network.ConnProcessor{}
	tcpServer.AddProcessor(network.CONN, connProcessor)
	tcpServer.Stop()

	ctx := context.Background()
	if err := envconfig.Process(ctx, &config.ApplicationConfig); err != nil {
		log.Fatal(err)
	}
	// to set gin Mode, either you can use env or code
	// - using env:	export GIN_MODE=release
	// - using code:	gin.SetMode(gin.ReleaseMode)
	// if envValue, isExisting := os.LookupEnv("GIN_MODE"); isExisting {
	// 	gin.SetMode(envValue)
	// } else {
	// 	gin.SetMode(gin.DebugMode)
	// }

	gin.SetMode(config.GetHttpServerConfig().GinMode)

	r := gin.Default()

	server := &http.Server{
		Addr:           config.GetHttpServerConfig().Host + ":" + config.GetHttpServerConfig().Port,
		Handler:        r,
		ReadTimeout:    time.Duration(config.GetHttpServerConfig().ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.GetHttpServerConfig().WriteTimeout) * time.Second,
		MaxHeaderBytes: config.GetHttpServerConfig().MaxHeaderBytes,
	}

	router.InitRouter(r)

	server.ListenAndServe()
}
