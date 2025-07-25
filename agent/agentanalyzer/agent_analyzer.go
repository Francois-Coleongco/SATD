package agentanalyzer

import (
	"SATD/types"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func Tcp_Packet_Analyzer(packet gopacket.Packet, ip string, ratio map[string]*types.SynAckRatio, mut *sync.Mutex) {
	// for syn flood / tcp half scan, but for tcp half scan, we need to make sure that we track syn/synacks per port on the remote ip
	mut.Lock()
	defer mut.Unlock()

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

	for _, i := range ratio {
		fmt.Fprintf(os.Stderr, "i.SynAcks: %v / %v\n", i.SynAcks, i.Syns)
	}

}
