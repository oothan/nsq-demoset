package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/cmd/front_api/handler"
	_ "nsq-demoset/app/app-services/conf"
	_ds "nsq-demoset/app/nsq-services/ds"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// load datasource
	ds := _ds.NewDataSource()

	// server
	router := gin.Default()
	h := handler.NewHandler(&handler.HConfig{
		R:  router,
		DS: ds,
	})
	h.Register()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: h.R,
	}

	go func() {
		logger.Sugar.Info("Server started listening on port: ", os.Getenv("APP_PORT"))
		if err := server.ListenAndServe(); err != nil {
			logger.Sugar.Error("Failed to initialized server: ", err.Error())
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
