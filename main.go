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

	go startTcpServer()

	// tcpServer.Stop()

	ctx := context.Background()
	if err := envconfig.Process(ctx, &config.ApplicationConfig); err != nil {
		log.Fatal(err)
	}

	startHttpServer()

}

func startTcpServer() {
	addr := network.Addr{
		Host: "localhost",
		Port: "8081",
	}
	tcpServer, _ := network.NewTcpServer(&network.TcpServerConfig{
		Network: "tcp",
		Addr:    addr,
	})
	err := tcpServer.Init()
	if err != nil {
		log.Printf("TCP server init failure. %s", err)
		return
	}

	err = tcpServer.Start()
	if err != nil {
		log.Printf("TCP server init failure. %s", err)
		return
	}

	log.Printf("TCP server startup")
}

func startHttpServer() {
	// to set gin Mode, either you can use env or code
	// - using env:    export GIN_MODE=release
	// - using code:    gin.SetMode(gin.ReleaseMode)
	// if envValue, isExisting := os.LookupEnv("GIN_MODE"); isExisting {
	//     gin.SetMode(envValue)
	// } else {
	//     gin.SetMode(gin.DebugMode)
	// }
	gin.SetMode(config.GetHttpServerConfig().GinMode)

	r := gin.Default()

	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    time.Duration(config.GetHttpServerConfig().ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.GetHttpServerConfig().WriteTimeout) * time.Second,
		MaxHeaderBytes: config.GetHttpServerConfig().MaxHeaderBytes,
	}

	router.InitRouter(r)

	server.ListenAndServe()
}
