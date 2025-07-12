package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/google/gopacket/pcap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"

	pb "server-agent-threat-detection/satd/v1"
)

var (
	id = 0111

	interfaceName = flag.String("interface_name", "enp0s3", "give a valid network interface (hint: run `ip a`)")
	MTU           = flag.Int64("MTU", 1500, "give the MTU of your network interface")                                                   // 1500 default ethernet
	isPromiscuous = flag.Bool("promiscuous_mode", false, "set promiscuous mode to true if you wish to see packets not for your device") // 1500 default ethernet

)

func start_data_stream(client pb.ServerFeederClient) {
	ctx := context.Background() // should be ample time to send a chunk

	handle, err := pcap.OpenLive(*interfaceName, int32(*MTU), *isPromiscuous, pcap.BlockForever)

	if err != nil {
		log.Fatalf("error starting packet capture in pcap.OpenLive, error thrown: %s\n", err)
	}

	stream, err := client.Feed(ctx)

	if err != nil {
		log.Fatalf("error creating stream from client in start_data_stream, error thrown: %s\n", err)
	}

	lt := []byte{byte(handle.LinkType())}

	stream.Send(&pb.NetDat{Payload: lt})

	for {
		data, ci, err := handle.ReadPacketData()

		if err != nil {
			log.Printf("couldn't read  this packet %s\n", err)
		}

		fmt.Println(ci)

		err = stream.Send(&pb.NetDat{Payload: data})
		if err != nil {
			log.Printf("error occurred during stream of to_send in start_data_stream, error thrown: %s\n", err)
			break
		}
	}

	r, err := stream.CloseAndRecv()

	if err != nil {
		log.Printf("failure in stream.CloseAndRecv(), error thrown: %s\n", err)
	}

	log.Printf("stream summary %v\n", r)
}

func main() {
	flag.Parse()
	caCert, err := os.ReadFile("cert.pem")

	if err != nil {
		log.Fatalf("failed to read CA cert: %s", err)
	}

	caCertPool := x509.NewCertPool()

	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS((&tls.Config{
		RootCAs: caCertPool,
	}))

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("couldn't connect to server, threw error: %s", err)
	}

	defer conn.Close()

	c := pb.NewServerFeederClient(conn)

	start_data_stream(c)

}
