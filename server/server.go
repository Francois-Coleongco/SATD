package main

import (
	"crypto/tls"
	"fmt"
	"io"

	"log"
	"net"

	pb "server-agent-threat-detection/satd/v1"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// need to add ip whitelist for server to look through before accepting a connection

type serverFeederServer struct {
	pb.UnimplementedServerFeederServer
}

func (s *serverFeederServer) Feed(stream pb.ServerFeeder_FeedServer) error {

	totalBytes := 0

	lt_dat, err := stream.Recv()

	linkType := layers.LinkType(lt_dat.GetPayload()[0])

	if err != nil {
		log.Printf("error reading link type %s\n", err)
	}

	for {
		netDat, err := stream.Recv()
		if err == io.EOF {
			log.Println("end of stream")
			break
		}

		if err != nil {
			log.Printf("error receiving netDat, error thrown: %s", err)
			return err
		}
		log_chunk := netDat.GetPayload()

		packet := gopacket.NewPacket(log_chunk, linkType, gopacket.Default)
		// need to index packet meta data into elastic search (use docker to spawn it)
		fmt.Println("received: ", packet)
		totalBytes += len(log_chunk)
	}

	return stream.SendAndClose(&pb.RecConf{Success: true, BytesReceived: int64(totalBytes)})
}

func main() {

	// spawn the elastic search query and analysis go routine at the beginning.

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
