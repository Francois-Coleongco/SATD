import { useState } from 'react';
import './App.css'

function App() {

	const [authed, setAuthed] = useState(false);
	const [badAuth, setBadAuthed] = useState(false);
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");

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
			updateAgents()
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
		}


	}

	const agents = [
		{
			agentID: 'AGENT001',
			agentIP: '192.168.1.1',
			uniqueIPs: ['192.168.1.2', '192.168.1.3'],
			threatSummary: 'High risk - 3 threats detected',
			health: 'Healthy',
			lastCheckin: '2023-08-15 10:30:00',
		},
		{
			agentID: 'AGENT002',
			agentIP: '192.168.1.4',
			uniqueIPs: ['192.168.1.5'],
			threatSummary: 'Medium risk - 1 threat detected',
			health: 'Warning',
			lastCheckin: '2023-08-15 09:50:00',
		},
		{
			agentID: 'AGENT003',
			agentIP: '192.168.1.6',
			uniqueIPs: ['192.168.1.7', '192.168.1.8', '192.168.1.9'],
			threatSummary: 'Low risk - 0 threats detected',
			health: 'Healthy',
			lastCheckin: '2023-08-15 11:00:00',
		},
	];

	if (authed) {
		return (
			<div className="min-h-screen flex flex-col items-center p-8 bg-gray-100">
				<h1 className="text-4xl font-bold text-gray-800 mb-8">Agent Dashboard</h1>
				<div className="w-full max-w-6xl bg-black rounded-lg shadow-lg overflow-hidden">
					<div className="grid grid-cols-6 bg-gray-200 text-gray-700 font-semibold py-4 px-6">
						<div>Agent ID</div>
						<div>Agent IP</div>
						<div>Unique IPs</div>
						<div>Threat Summary</div>
						<div>Health</div>
						<div>Last Check-In</div>
					</div>
					{agents.map((agent) => (
						<div
							className="grid grid-cols-6 py-4 px-6 border-b border-gray-200 hover:bg-gray-800"
							key={agent.agentID}
						>
							<div>{agent.agentID}</div>
							<div>{agent.agentIP}</div>
							<div>{agent.uniqueIPs.join(', ')}</div>
							<div>{agent.threatSummary}</div>
							<div
								className={`font-semibold text-center rounded-lg py-1 px-2 ${agent.health === 'Healthy'
									? 'bg-green-100 text-green-700'
									: agent.health === 'Warning'
										? 'bg-yellow-100 text-yellow-700'
										: 'bg-red-100 text-red-700'
									}`}
							>
								{agent.health}
							</div>
							<div>{agent.lastCheckin}</div>
						</div>
					))}
				</div>
			</div>
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

