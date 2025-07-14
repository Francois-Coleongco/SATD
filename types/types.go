package types

import "time"

type PacketMeta struct {
	SrcIP     string
	DstIP     string
	SrcPort   string
	DstPort   string
	Protocol  string
	Timestamp time.Time
}

type SynAckRatio struct {
	Syns    uint32
	SynAcks uint32
}
