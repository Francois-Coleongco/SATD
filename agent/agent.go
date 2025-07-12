package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "server-agent-threat-detection/satd/v1"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
	id   = 0111
)

func start_data_stream(client pb.ServerFeederClient, to_send [][]byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20) // should be ample time to send a chunk
	defer cancel()

	stream, err := client.Feed(ctx)

	if err != nil {
		log.Fatalf("error creating stream from client in start_data_stream, error thrown: %s", err)
	}

	for _, dat := range to_send {
		err := stream.Send(&pb.NetDat{Payload: dat})
		if err != nil {
			log.Fatalf("error occurred during stream of to_send in start_data_stream")
		}
	}

	r, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("failure in stream.CloseAndRecv(), error thrown: ", err)
	}

	log.Printf("stream summary %v\n", r)
}

func main() {
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

	var dummy_data = [][]byte{
		[]byte("i bought a cat"),
		[]byte("i have a cat"),
		[]byte("i had a cat"),
		[]byte("cry"),
	}

	start_data_stream(c, dummy_data)

}
