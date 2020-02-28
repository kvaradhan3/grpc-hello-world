package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	hw "kannan.ieee.org/proto/helloWorld"

	ts "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

const (
	port_default  = "65432"
	host_default  = "localhost"
	tag_default   = "H"
	count_default = 1
)

var (
	port  string
	host  string
	tag   string
	count int
)

func main() {
	flag.StringVar(&port, "port", port_default, "port to connect to")
	flag.StringVar(&host, "host", host_default, "host to connect to")
	flag.StringVar(&tag, "tag", tag_default, "Tag to log with")
	flag.IntVar(&count, "count", count_default, "#msg iterations (default 1)")
	flag.Parse()

	rand.Seed(time.Now().Unix())
	cname, err := net.LookupCNAME(host)
	if err != nil {
		cname = host
	}
	addrs, err := net.LookupHost(cname)
	if err != nil {
		log.Fatalf("LookupHost(%s): no addresses returned\n", cname)
	}
	addr := addrs[rand.Intn(len(addrs))]
	name := fmt.Sprintf("%s%06d", tag, rand.Intn(999999))

	conn, err := grpc.Dial(net.JoinHostPort(addr, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	h := hw.NewHelloWorldClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sum := int64(0)
	for cnt := count; cnt > 0; cnt-- {
		t := time.Now()
		r, err := h.HelloWorld(ctx,
			&hw.Hello{
				SendTime: &ts.Timestamp{
					Seconds: t.Unix(),
					Nanos:   int32(t.UnixNano() % 1e9),
				},
				Rqst: name,
			})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		rtime := int64(r.SendTime.Seconds)*1e9 + int64(r.SendTime.Nanos)
		diff := rtime - t.UnixNano()
		sum += diff
		log.Printf(">msg: %s <cookie: %d diff: %dns., sum: %d",
			r.Resp, r.Cookie, diff, sum)
	}

	if count > 1 {
		ave := sum / int64(count)
		msec, usec := ave/1e6, (ave/1e3)%1e3

		log.Printf("rtt estimate over %d iterations: %d (total), %d.%dms per iter",
			count, sum, msec, usec)
	}
	return
}
