package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "server-agent-threat-detection/satd/v1"
	"server-agent-threat-detection/types"
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

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {

		if err != nil {
			log.Printf("couldn't read  this packet %s\n", err)
		}

		fmt.Println("fucking hell")

		netLayer := packet.NetworkLayer()

		var dstIP string = "--"
		var srcIP string = "--"
		var dstPort string = "--"
		var srcPort string = "--"
		var protocol string = "--"

		if netLayer != nil {
			dstIP = netLayer.NetworkFlow().Dst().String()
			srcIP = netLayer.NetworkFlow().Dst().String()

			transLayer := packet.TransportLayer()

			if transLayer != nil {
				dstPort = transLayer.TransportFlow().Dst().String()
				srcPort = transLayer.TransportFlow().Src().String()
				protocol = transLayer.LayerType().String()
			}
		}

		dat := types.PacketMeta{
			SrcIP:     srcIP,
			DstIP:     dstIP,
			SrcPort:   srcPort,
			DstPort:   dstPort,
			Protocol:  protocol,
			Timestamp: packet.Metadata().Timestamp,
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		enc.Encode(dat)

		fmt.Println("error???")

		err = stream.Send(&pb.NetDat{Payload: buf.Bytes()})
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
