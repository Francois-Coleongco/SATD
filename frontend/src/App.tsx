import { useEffect, useState } from 'react';
import './App.css'

interface AgentData {
	AgentID: string
	AgentIP: string
	UniqueIPs: Map<string, number> // ips, AbuseIPDB score. these ips are by the day
	LastCheckIn: Date
}

function App() {

	const [authed, setAuthed] = useState(false);
	const [badAuth, setBadAuthed] = useState(false);
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [agents, setAgents] = useState<Map<string, AgentData>>(new Map());

	const authenticate = async () => {
		console.log("sending cresd", username, password)

		const csrfResp = await fetch('https://localhost:3000/csrf', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
			}
		});
		const csrfData = await csrfResp.json()

		const authResp = await fetch('https://localhost:3000/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'X-CSRF-Token': csrfData.csrfToken
			},
			body: JSON.stringify({ username: username, password: password }),
			credentials: "include"
		});

		if (!authResp.ok) {
			console.log(authResp)
			setAuthed(false);
			setBadAuthed(true);
			return "";
		} else {
			setAuthed(true)
		}

		const authData = await authResp.json()

		return authData.token
	}

	const updateAgents = async () => {
		const evtSource = new EventSource("https://localhost:3000/dashboard-stream", {
			withCredentials: true
		});
		evtSource.onmessage = (event) => {
			const parsedData = JSON.parse(event.data)
			console.log("parsedData", parsedData)

			let lastCheckIn: Date;
			try {
				lastCheckIn = new Date(parsedData.LastCheckIn);
				if (isNaN(lastCheckIn.getTime())) {
					console.warn(`Invalid date for AgentID ${parsedData.AgentID}: ${parsedData.LastCheckIn}`);
					lastCheckIn = new Date(); // Fallback to current time
				}
			} catch (error) {
				console.warn(`Error parsing date for AgentID ${parsedData.AgentID}: ${error}`);
				lastCheckIn = new Date();
			}

			const agentData: AgentData = {
				AgentID: parsedData.AgentID || 'Unknown',
				AgentIP: parsedData.AgentIP || 'Unknown',
				UniqueIPs: new Map(Object.entries(parsedData.UniqueIPs || {})),
				LastCheckIn: lastCheckIn,
			};

			setAgents((prevAgents) => {
				const newAgents = new Map(prevAgents);
				newAgents.set(agentData.AgentID, agentData);
				return newAgents;
			});
		};

		evtSource.onerror = () => {
			console.error('EventSource failed');
			evtSource.close();
		};

		return () => {
			evtSource.close();
		};
	};




	useEffect(() => {
		if (authed) {
			updateAgents()
		}
	}, [authed])

	if (authed) {
		return (
			< div className="min-h-screen flex flex-col items-center p-8 bg-gray-100" >
				<h1 className="text-4xl font-bold text-gray-800 mb-8">Agent Dashboard</h1>
				test
				<div className="w-full max-w-6xl bg-white rounded-lg shadow-lg overflow-hidden">
					{/* Header */}
					<div className="grid grid-cols-4 py-4 px-6 bg-gray-900 text-white font-bold text-sm uppercase tracking-wider">
						<div>Agent ID</div>
						<div>Agent IP</div>
						<div>Unique IPs</div>
						<div>Last Check-In</div>
					</div>

					{/* Data Rows */}
					{[...agents.values()].map((agent) => (
						<div
							className="grid grid-cols-4 gap-4 py-4 px-6 text-sm border-b border-gray-200 hover:bg-gray-100 transition-colors duration-200"
							key={agent.AgentID}
						>
							<div className="text-gray-800 font-medium break-all">{agent.AgentID}</div>
							<div className="text-gray-700 break-all">{agent.AgentIP}</div>

							{/* Unique IPs */}
							<div className="space-y-1 max-h-32 overflow-y-auto pr-2">
								{[...agent.UniqueIPs.entries()].map(([ip, score]) => (
									<div key={ip} className="flex justify-between text-gray-600 text-xs border-b border-gray-100 pb-1">
										<span className="font-mono">{ip}</span>
										<span className="text-right font-semibold">{score}</span>
									</div>
								))}
							</div>

							<div className="text-gray-600">{agent.LastCheckIn.toLocaleString()}</div>
						</div>
					))}
				</div>
			</div >
		);
	} else {
		return (
			<div className="min-h-screen flex items-center justify-center p-6">
				<div className="w-full max-w-sm bg-white p-8 rounded-lg shadow-lg">
					<h2 className="text-2xl font-semibold text-gray-800 mb-6 text-center">Login</h2>
					<form className="space-y-6">
						<div>
							<label htmlFor="username" className="block text-gray-700 font-medium mb-2">Username</label>
							<input
								type="text"
								id="username"
								value={username}
								onChange={(e) => setUsername(e.target.value)}
								className="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-700"
								placeholder="Enter your username"
							/>
						</div>

						<div>
							<label htmlFor="password" className="block text-gray-700 font-medium mb-2">Password</label>
							<input
								type="password"
								id="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								className="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-700"
								placeholder="Enter your password"
							/>
						</div>
						<div>{badAuth}</div>

						<button
							type='button'
							className="w-full bg-blue-600 text-white p-3 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
							onClick={authenticate}
						>
							Login
						</button>
					</form>
				</div>
			</div>
		);
	}

}

export default App;

