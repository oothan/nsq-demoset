package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/cmd/market/server"
	marketpb "nsq-demoset/app/app-services/internal/proto/market/v1/pb"
)

func main() {
	port := flag.String("port", "8081", "Default port is 8081.")
	flag.Parse()

	addr := fmt.Sprintf(":%s", *port)

	gs := grpc.NewServer()

	// create the instance of market
	market := server.NewMarketServer()

	// register instance of market
	marketpb.RegisterMarketServer(gs, market)

	// register the reflection service which allows clients to determine the methods
	reflection.Register(gs)

	// create a TCP socket for inbound connection
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	logger.Sugar.Debug("Server is listening on port: ", *port)

	// listen fot request
	gs.Serve(l)
}
