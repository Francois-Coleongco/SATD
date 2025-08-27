package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"

	// "crypto/x509"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"log"
	"net"

	pb "SATD/network_comms/v1"
	"SATD/server/auth"
	"SATD/server/serveranalyzer"
	"SATD/types"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// need to add ip whitelist for server to look through before accepting a connection

// internals
var (
	latencyLogger  *log.Logger                = nil
	agentsMap      map[string]types.AgentInfo // key is the id of the agent
	agentsMapMutex sync.Mutex
	utcTime        string
)

// auth types
var (
	dashboardServerProtAddr string
	dashboardServerAuthAddr string
	dashboardJWT            string
	dashClient              *http.Client
	dashUserCreds           types.DashCreds
)

// apis
var (
	elasticApiKey string
	esClient      *elasticsearch.Client

	ipdbApiKey string
	ipdbClient *http.Client
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

func healthCheck(agentID string, agentIP string) types.AgentInfo {
	// check ips by querying elastic search for the day's index
	query := fmt.Sprintf(`{
		"query": {
			"bool": {
				"must": [
					{ "match": { "AgentID": "%s" } }
				],
				"must_not": [
					{ "term": { "SrcIP": "%s" } }
				]
			}
		}
	}`, agentID, agentIP)

	res, err := esClient.Search(
		esClient.Search.WithIndex(utcTime),
		esClient.Search.WithBody(strings.NewReader(query)),
		esClient.Search.WithPretty(),
	)

	if err != nil {
		log.Println("bad ES search for SrcIPs")
		return types.AgentInfo{}
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("couldn't ReadAll res.Body")
		return types.AgentInfo{}
	}

	var r types.ESResponse

	err = json.Unmarshal(data, &r)

	if err != nil {
		log.Println("couldn't unmarshal SrcIP results into ESResponse r")
		return types.AgentInfo{}
	}

	inf := types.AgentInfo{
		AgentID:       agentID,
		AgentIP:       agentIP,
		ThreatSummary: "", // include possibly scanned if any
		UniqueIPs:     make(map[string]int),
		LastCheckIn:   time.Now().UTC(),
	}

	for _, hit := range r.Hits.Hits {
		_, exists := inf.UniqueIPs[hit.Source.SrcIP]
		if !exists {
			score, err := serveranalyzer.IpCheckAbuseIPDB(hit.Source.SrcIP, &ipdbApiKey, ipdbClient)
			if err != nil {
				score = -1
			}
			inf.UniqueIPs[hit.Source.SrcIP] = score

			// should instead be a tally average to get the overall score
			// if 0 <= score || score <= 30 {
			// 	inf.ThreatSummary = "Low"
			// } else if 31 <= score || score <= 60 {
			// 	inf.ThreatSummary = "Med"
			// } else if 61 <= score || score <= 100 {
			// 	inf.ThreatSummary = "High"
			// } else {
			// 	inf.ThreatSummary = "Unknown" // (no ipdb score)
			// }

			log.Println("adding unique ip: ", hit.Source.SrcIP)
		}
	}

	return inf

}

func sendBeatToDash(agentID string, inf *types.AgentInfo) {

	// type AgentInfo struct { // heartbeat data
	// 	AgentIP       string
	// 	ThreatSummary string
	// 	Health        string
	// 	UniqueIPs     map[string]int // ips, AbuseIPDB score. these ips are by the day
	// 	LastCheckIn   time.Time
	// }
	// make a request to the protected endpoint
	// if unauthorized, request for a JWT token using a user and password in .server_env
	// retry the protected endpoint request

	inf.AgentID = agentID

	fmt.Println("this is info being sent: ", inf)

	jsonBeat, err := json.Marshal(inf)

	if err != nil {
		log.Printf("error marshalling heartbeat data to json, error thrown: %s\n", err)
		return
	}

	reader := bytes.NewReader(jsonBeat)

	for attempts := 0; attempts < types.MAX_PROT_ATTEMPTS_BEFORE_REAUTH; attempts++ {
		req, err := http.NewRequest("POST", dashboardServerProtAddr, reader)
		reader.Seek(0, io.SeekStart)

		if err != nil {
			log.Printf("could not send heartbeat to protected endpoint, error thrown: %s\n", err)
		}

		req.Header.Set("Authorization", "Bearer "+dashboardJWT)
		req.Header.Set("Content-Type", "application/json")

		res, err := dashClient.Do(req)

		if err != nil {
			log.Printf("unable to get response for dashClient, error thrown: %s", err)
			continue
		}

		if res.StatusCode == 200 {
			fmt.Println("SUPPOSED TO BREAK OUT HEREEE")
			break
		} else if res.StatusCode == 401 {
			log.Printf("was unauthorized when attempting to send heartbeat")
			auth.AuthToDash(dashClient, 4, dashboardServerAuthAddr, dashUserCreds, &dashboardJWT)
		}

		time.Sleep(time.Millisecond * 250)
	}
}

func processHeartbeat(data []byte, agentIP string) { // data is the heartbeat data
	agentID := string(data)
	inf := healthCheck(agentID, agentIP)

	agentsMapMutex.Lock()
	agentsMap[agentID] = inf
	agentsMapMutex.Unlock()

	sendBeatToDash(agentID, &inf)
}

func (s *serverFeederServer) Feed(stream pb.ServerFeeder_FeedServer) error {
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return fmt.Errorf("issue getting peer from context in Feed, error thrown")
	}

	addr, ok := p.Addr.(*net.TCPAddr)

	if !ok {
		return fmt.Errorf("couldn't get TCP address, we only allow TCP 'round these parts. gtfo, ur cringe")
	}

	agentIP := addr.IP.String()

	fmt.Println("called Feed")
	totalBytes := 0

	utcTime := time.Now().UTC().Format("2006-01-02")

	createESIndex(utcTime)

	for {

		log.Println("supposed to read here")
		netDat, err := stream.Recv()

		if err == io.EOF {
			log.Println("END OF STREAM")
			break
		}

		if err != nil {
			log.Printf("error receiving netDat, error thrown: %s", err)
			return err
		}

		data := netDat.GetPayload()

		if netDat.GetIsHeartbeat() {
			fmt.Println("this was a heartbeat")
			go processHeartbeat(data, agentIP)
			continue
		}

		currUtcTime := time.Now().UTC().Format("2006-01-02")
		createESIndex(currUtcTime)
		utcTime = currUtcTime

		var packetMetaData types.PacketMeta

		buf := bytes.NewBuffer(data)

		dec := gob.NewDecoder(buf)

		err = dec.Decode(&packetMetaData)

		if err != nil {
			log.Printf("couldn't decode netData in Feed loop, error thrown: %s\n", err)
			continue
		}

		latency := time.Now().UTC().Sub(packetMetaData.Timestamp)

		latencyLogger.Printf("%d ms", latency.Milliseconds()) // do the math of averaging after running to not interfere with latency calculations

		log.Printf("%s %s, %s, %s, %s, %s, %s\n", packetMetaData.AgentID, packetMetaData.SrcIP, packetMetaData.DstIP, packetMetaData.SrcPort, packetMetaData.DstPort, packetMetaData.Protocol, packetMetaData.Timestamp)

		packetMetaData.AgentIP = agentIP

		pData, err := json.Marshal(packetMetaData)

		if err != nil {
			log.Printf("error serializing to json, error thrown: %s\n", err)
		}

		esClient.Index(utcTime, bytes.NewReader(pData))

		// need to index packet meta data into elastic search (use docker to spawn it)
		totalBytes += len(data)
	}

	return stream.SendAndClose(&pb.RecConf{Success: true, BytesReceived: int64(totalBytes)})
}

func heartbeatCheck() {
	for {
		currTime := time.Now()
		log.Println("checking beat")
		agentsMapMutex.Lock()
		log.Println("I GOT A LOCK YAYYYYY")

		for agent, agentData := range agentsMap {
			if agentData.LastCheckIn.IsZero() {
				continue
			}
			diff := currTime.Sub(agentData.LastCheckIn)
			if diff.Seconds() >= 4 {
				log.Println("agent", agent, " is dead")
				agentsMap[agent] = types.AgentInfo{} // zero time. aka in 1970 or something i think. indicating that this agent is dead;
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

func initializeEnvs() {
	dashboardServerProtAddr = os.Getenv("DASHBOARD_SERVER_PROT_ADDR")
	dashboardServerAuthAddr = os.Getenv("DASHBOARD_SERVER_AUTH_ADDR")
	dashboardServerAuthAddr = os.Getenv("DASHBOARD_SERVER_AUTH_ADDR")
	dashUserCreds.Username = os.Getenv("NODEJS_USER")
	fmt.Println("node user was ", dashUserCreds.Username)
	dashUserCreds.Password = os.Getenv("NODEJS_PASS")
	fmt.Println("node pass was ", dashUserCreds.Password)
	elasticApiKey = os.Getenv("ELASTIC_API_KEY")
}

func main() {

	// fetchWhitelist()

	agentsMap = make(map[string]types.AgentInfo)

	godotenv.Load(".server_env")
	initLogger()
	initTLS()
	initializeEnvs()

	cert, err := os.ReadFile("node_cert.pem")

	if err != nil {
		log.Fatalf("couldn't open node_cert.pem, error thrown: %s", err)
	}

	caCertPool := x509.NewCertPool()

	caCertPool.AppendCertsFromPEM(cert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	dashClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	ipdbClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
		Timeout: 10 * time.Second,
	}

	if auth.AuthToDash(dashClient, 4, dashboardServerAuthAddr, dashUserCreds, &dashboardJWT) != nil {
		return
	}

	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		APIKey: elasticApiKey,
	})

	go heartbeatCheck()

	log.Println("made it after")
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
