export type AgentInfo = { // aka heartbeat data
	ThreatSummary: string,
	Health: string,
	LastCheckIn: Date,
	// to be added // CPUUsage: number, // as percent
	// to be added // RAMUsage: number, // as percent
}

export interface JwtPayload {
	username: string;
	iat?: number;
	exp?: number;
}
