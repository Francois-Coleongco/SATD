package main

import (
	"crypto/tls"
	"fmt"
	"io"

	"log"
	"net"

	pb "server-agent-threat-detection/satd/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// need to add ip whitelist for server to look through before accepting a connection

type serverFeederServer struct {
	pb.UnimplementedServerFeederServer
}

func (s *serverFeederServer) Feed(stream pb.ServerFeeder_FeedServer) error {

	totalBytes := 0

	for {
		netDat, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("end of stream")
			break
		}

		if err != nil {
			log.Printf("error receiving netDat, error thrown: %s", err)
			return err
		}
		log_chunk := netDat.GetPayload()
		totalBytes += len(log_chunk)
		fmt.Println(string(log_chunk))
	}

	msg := fmt.Sprintln("once processed, this message should tell the client whether there is a potential attack and take subsequent measures to quarantine/kill the connection / find the source etcc etc")

	return stream.SendAndClose(&pb.RecConf{Success: true, Message: msg, BytesReceived: int64(totalBytes)})

}

func main() {

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")

	if err != nil {
		log.Fatalf("failed to load keypair, error thrown: %s", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	})

	s := grpc.NewServer(grpc.Creds(creds))

	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("failed to listen, error thrown: %s", err)
	}

	pb.RegisterServerFeederServer(s, &serverFeederServer{})

	log.Println("Server Booted")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve, error thrown: %s", err)
	}

}
