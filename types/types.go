package types

import "time"

type AgentInfo struct { // heartbeat data
	AgentIP       string
	ThreatSummary string
	Health        string
	UniqueIPs     map[string]int // ips, AbuseIPDB score. these ips are by the day
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

// they should really provide this in an easy to find place smh

type ESResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source struct {
				AgentID   string `json:"AgentID"`
				SrcIP     string `json:"SrcIP"`
				DstIP     string `json:"DstIP"`
				SrcPort   string `json:"SrcPort"`
				DstPort   string `json:"DstPort"`
				Protocol  string `json:"Protocol"`
				Timestamp string `json:"Timestamp"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
