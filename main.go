package main

import (
	"context"
	"go-networking/config"
	"go-networking/docs"
	"go-networking/log"
	"go-networking/network"
	"go-networking/router"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title http server
// @version 1.0
// @description http server
// @termsOfService http://swagger.io/terms/

// @contact.name errorfatal89@gmail.com
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:8080
// @BasePath /api/v1
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.InitLogger()

	go startTcpServer()

	// tcpServer.Stop()

	ctx := context.Background()
	if err := envconfig.Process(ctx, &config.ApplicationConfig); err != nil {
		log.ErrorErr(err)
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
		log.ErrorErrMsg(err, "TCP server init failure.")
		return
	}

	err = tcpServer.Start()
	if err != nil {
		log.ErrorErrMsg(err, "TCP server init failure.")
		return
	}

	log.Info("TCP server startup")
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
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.InitRouter(r)

	server.ListenAndServe()
}
