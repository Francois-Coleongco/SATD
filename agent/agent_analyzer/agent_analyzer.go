package agent_analyzer

import (
	"SATD/types"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func Tcp_Packet_Analyzer(packet gopacket.Packet, ip string, ratio map[string]*types.SynAckRatio) {

	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	if tcpLayer == nil {
		log.Println("error accessing tcp layer interface")
		return
	}

	tcp, ok := tcpLayer.(*layers.TCP)

	if !ok {
		log.Println("error accessing tcp layer struct")
		return
	}

	if ratio[ip] == nil {
		ratio[ip] = &types.SynAckRatio{}
	}

	if tcp.SYN && !tcp.ACK {
		ratio[ip].Syns++
	} else if tcp.SYN && tcp.ACK {
		ratio[ip].SynAcks++
	}

}
