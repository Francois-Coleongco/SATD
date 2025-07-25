package types

import "time"

type AgentBeatData struct {
	LastBeat time.Time
}

type PacketMeta struct {
	AgentID   string
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
