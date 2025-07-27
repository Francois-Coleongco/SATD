package types

import "time"

type AgentInfo struct { // heartbeat data
	AgentIP       string
	ThreatSummary string
	Health        string
	LastCheckIn   time.Time
}

type PacketMeta struct {
	AgentID   string
	AgentIP   string
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
