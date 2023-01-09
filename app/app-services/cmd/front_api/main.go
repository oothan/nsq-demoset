package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	logger "nsq-demoset/app/_applib"
	libnsq "nsq-demoset/app/_applib/nsq"
	_ "nsq-demoset/app/app-services/cmd/front_api/docs"
	"nsq-demoset/app/app-services/cmd/front_api/handler"
	"nsq-demoset/app/app-services/internal/ds"
	"os"
	"os/signal"
	"syscall"
)

// @title App
// @version 0.1
// @description application description

// @BasePath /

// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func main() {

	port := flag.String("port", "8001", "Default port is 8001")
	mPort := flag.String("market", "8081", "Default port is 8081")
	flag.Parse()

	//addr := fmt.Sprintf(":%s", *port)
	//mAddr := fmt.Sprintf(":%s", *mPort)
	addr := net.JoinHostPort("", *port)
	mAddr := net.JoinHostPort("", *mPort)

	libnsq.InitNSQProducer()

	// load datasource
	ds := ds.NewDataSource()

	//router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	h := handler.NewHandler(&handler.HConfig{
		R:             router,
		DS:            ds,
		MarketRPCAddr: mAddr,
	})
	h.Register()

	server := http.Server{
		Addr:    addr,
		Handler: h.R,
	}

	go func() {
		logger.Sugar.Info("Server started listening on port: ", *port)
		if err := server.ListenAndServe(); err != nil {
			logger.Sugar.Error("Failed to initialized server on port :", *port, " ", err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	// shutdown close
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Sugar.Error("Failed to shutdown server: ", err.Error())
	}
}
