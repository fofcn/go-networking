package main

import (
	"context"
	"go-networking/config"
	"go-networking/docs"
	"go-networking/infra/model"
	"go-networking/log"
	"go-networking/network"
	"go-networking/router"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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

	insertTestData()

	go startTcpServer()

	// tcpServer.Stop()

	ctx := context.Background()
	if err := envconfig.Process(ctx, &config.ApplicationConfig); err != nil {
		log.ErrorErr(err)
	}

	startHttpServer()

}

func insertTestData() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: "root:example@tcp(127.0.0.1:3306)/trade_order?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		// DefaultStringSize:         65535,                                                                                // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent SQL logging
	})
	if err != nil {
		log.Error("connect db errr")
		return
	}

	err = db.AutoMigrate(&model.UserModel{}, &model.TradeOrderModel{}, &model.OrderProductModel{}, &model.ProductModel{})
	if err != nil {
		log.Errorf("%v", err)
	}

	min := time.Date(1911, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	birthday := time.Unix(sec, 0)

	var wg sync.WaitGroup
	goroutines := make(chan struct{}, 10) // 10 is goroutine count

	batchSize := 100

	// Create 50k users and 50k products
	for i := 0; i < 500000; i++ {
		// Add to wait group
		wg.Add(1)
		goroutines <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-goroutines }()
			var users []model.UserModel = make([]model.UserModel, 0, batchSize)
			var products []model.ProductModel = make([]model.ProductModel, 0, batchSize)
			// Reformat this for your password storage needs
			users = append(users, model.UserModel{Name: "User " + strconv.Itoa(i), Email: "user" + strconv.Itoa(i) + "@test.com", Age: uint8(rand.Intn(100)), Birthday: &birthday, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			products = append(products, model.ProductModel{Name: "Product " + strconv.Itoa(i), Price: uint(rand.Intn(1000)), CreatedAt: time.Now(), UpdatedAt: time.Now()})

			// If users batch reaches 200, insert and reset slice
			db.Create(&users)
			log.Infof("create new 200 users")
			println("[println]create new 200 users")

			// If products batch reaches 200, insert and reset slice
			db.Create(&products)
			log.Infof("create new 200 products")
			println("[println]create new 200 products")
		}(i)
	}

	wg.Wait() // wait for all goroutines to finish

	// // Check if any remaining users/products to insert
	// if len(users) > 0 {
	// 	db.Create(&users)
	// }
	// if len(products) > 0 {
	// 	db.Create(&products)
	// }

	log.Info("Success insert test data to mysql")
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
