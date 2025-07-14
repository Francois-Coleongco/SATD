package main

import (
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"log"
	"net"

	pb "SATD/network_comms/v1"
	"SATD/types"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// need to add ip whitelist for server to look through before accepting a connection

var (
	elasticApiKey string
	esClient      *elasticsearch.Client
)

type serverFeederServer struct {
	pb.UnimplementedServerFeederServer
}

func (s *serverFeederServer) Feed(stream pb.ServerFeeder_FeedServer) error {

	totalBytes := 0

	currIndex := esClient.Indices.Create(time.Date())

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

		var packetMetaData types.PacketMeta

		buf := bytes.NewBuffer(log_chunk)

		dec := gob.NewDecoder(buf)

		err = dec.Decode(&packetMetaData)

		if err != nil {
			log.Println("couldn't decode netData in Feed loop, error thrown: %s\n", err)
		}

		fmt.Printf("%s, %s, %s, %s, %s, %s\n", packetMetaData.SrcIP, packetMetaData.DstIP, packetMetaData.SrcPort, packetMetaData.DstPort, packetMetaData.Protocol, packetMetaData.Timestamp)

		data, err := json.Marshal(packetMetaData)

		if err != nil {
			fmt.Println("error serializing to json, error thrown: %s\n", err)
		}

		esClient.Index(, bytes.NewReader(data))

		// need to index packet meta data into elastic search (use docker to spawn it)
		totalBytes += len(log_chunk)
	}

	return stream.SendAndClose(&pb.RecConf{Success: true, BytesReceived: int64(totalBytes)})
}

func main() {

	// spawn the elastic search query and analysis go routine at the beginning.

	godotenv.Load(".server_env")
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")

	if err != nil {
		log.Fatalf("failed to load keypair, error thrown: %s", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	})

	elasticApiKey = os.Getenv("ELASTIC_API_KEY")

	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		APIKey: elasticApiKey,
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
