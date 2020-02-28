package main

import (
	"flag"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	"math/rand"

	hw "kannan.ieee.org/proto/helloWorld"

	ts "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

const (
	port_default = "65432"
)

var (
	port string
)

type helloWorldService struct{}

func (h *helloWorldService) HelloWorld(ctx context.Context,
	in *hw.Hello) (*hw.World, error) {
	t := time.Now()
	cookie := rand.Int31()
	log.Printf("<msg: %s, >cookie: %d", in.Rqst, cookie)
	return &hw.World{
		SendTime: &ts.Timestamp{
			Seconds: t.Unix(),
			Nanos:   int32(t.UnixNano() % 1e9),
		},
		Resp:   in.Rqst,
		Cookie: cookie,
	}, nil
}

func main() {
	flag.StringVar(&port, "port", port_default, 
		"Specify port to use (default 65432) (1024, 65535]")
	flag.Parse()
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil || p <= 1024 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalf("Failed listen(%s): %v", port, err)
	}

	server := grpc.NewServer()

	hw.RegisterHelloWorldServer(server, &helloWorldService{})
	if err := server.Serve(lis); err != nil {
		log.Fatal("Failed to start server!: %v", err)
	}
}
