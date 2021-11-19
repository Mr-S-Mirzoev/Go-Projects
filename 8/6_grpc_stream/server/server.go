package main

import (
	"fmt"
	"log"
	"net"

	"gitlab.com/mailru-go/lectures-2021-2/8/6_grpc_stream/translit"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	translit.RegisterTransliterationServer(server, NewTr())

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
