package main

import (
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
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

// internals
var (
	latencyLogger  *log.Logger                    = nil
	agentsMap      map[string]types.AgentBeatData // key is the id of the agent
	agentsMapMutex sync.Mutex
)

var (
	elasticApiKey string
	esClient      *elasticsearch.Client
)

type serverFeederServer struct {
	pb.UnimplementedServerFeederServer
}

func createESIndex(utcTime string) {
	_, err := esClient.Indices.Create(utcTime)

	if err != nil {
		log.Printf("continuing without ES index, error thrown: %s\n", err)
	}
}

func (s *serverFeederServer) Feed(stream pb.ServerFeeder_FeedServer) error {

	fmt.Println("called Feed")
	totalBytes := 0

	utcTime := time.Now().UTC().Format("2006-01-02")

	createESIndex(utcTime)

	for {

		fmt.Println("supposed to read here")
		netDat, err := stream.Recv()

		if err == io.EOF {
			log.Println("END OF STREAM")
			break
		}

		if err != nil {
			log.Printf("error receiving netDat, error thrown: %s", err)
			return err
		}

		log_chunk := netDat.GetPayload()

		if netDat.GetIsHeartbeat() {
			agentsMapMutex.Lock()
			agentID := string(netDat.GetPayload())
			agentsMap[agentID] = time.Now()
			agentsMapMutex.Unlock()
			continue
		}

		currUtcTime := time.Now().UTC().Format("2006-01-02")
		createESIndex(currUtcTime)
		utcTime = currUtcTime

		var packetMetaData types.PacketMeta

		buf := bytes.NewBuffer(log_chunk)

		dec := gob.NewDecoder(buf)

		err = dec.Decode(&packetMetaData)

		if err != nil {
			log.Printf("couldn't decode netData in Feed loop, error thrown: %s\n", err)
			continue
		}

		latency := time.Now().UTC().Sub(packetMetaData.Timestamp)

		latencyLogger.Printf("%d ms", latency.Milliseconds()) // do the math of averaging after running to not interfere with latency calculations

		fmt.Printf("%s %s, %s, %s, %s, %s, %s\n", packetMetaData.AgentID, packetMetaData.SrcIP, packetMetaData.DstIP, packetMetaData.SrcPort, packetMetaData.DstPort, packetMetaData.Protocol, packetMetaData.Timestamp)

		data, err := json.Marshal(packetMetaData)

		if err != nil {
			fmt.Printf("error serializing to json, error thrown: %s\n", err)
		}

		esClient.Index(utcTime, bytes.NewReader(data))

		// need to index packet meta data into elastic search (use docker to spawn it)
		totalBytes += len(log_chunk)
	}

	return stream.SendAndClose(&pb.RecConf{Success: true, BytesReceived: int64(totalBytes)})
}

func heartbeatCheck() {
	for {
		currTime := time.Now()
		fmt.Println("checking beat")
		agentsMapMutex.Lock()
		fmt.Println("I GOT A LOCK YAYYYYY")

		for agent, agentData := range agentsMap {
			if agentData.LastBeat.IsZero() {
				continue
			}
			diff := currTime.Sub(agentData.LastBeat)
			if diff.Seconds() >= 4 {
				log.Println("agent", agent, " is dead")
				agentsMap[agent] = types.AgentBeatData{} // zero time. aka in 1970 or something i think. indicating that this agent is dead;
			} else {
				log.Println(agent, "IS STILL HERE!")
			}
		}
		agentsMapMutex.Unlock()

		time.Sleep(time.Second * 4) // same interval length as heartbeat sends
	}

}
func initLogger() {
	file, err := os.OpenFile("latency.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Printf("couldn't open latency.log, error thrown: %s\n", err)
		// not fatal, if cant start this it's fine
	}

	latencyLogger = log.New(file, "", log.Ltime|log.LUTC|log.Lmicroseconds)

}

func initTLS() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")

	if err != nil {
		log.Fatalf("failed to load keypair, error thrown: %s", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	})

	return creds

}

func main() {

	agentsMap = make(map[string]time.Time)

	godotenv.Load(".server_env")

	initLogger()

	elasticApiKey = os.Getenv("ELASTIC_API_KEY")

	var err error

	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		APIKey: elasticApiKey,
	})

	go heartbeatCheck()

	fmt.Println("made it after")
	s := grpc.NewServer(grpc.Creds(initTLS()))

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

	log.Println("Graceful Exit :)")

}
