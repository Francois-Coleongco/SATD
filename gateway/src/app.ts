import express from 'express'
import fs from 'fs'
import path from 'path'
import https from 'https'
import jsonwebtoken from 'jsonwebtoken'
import { userExists } from './auth'
import { AgentInfo } from './types'

interface JwtPayload {
	username: string;
	iat?: number;
	exp?: number;
}

const app = express()

app.use(express.json())


app.post('/login', async (req, res) => {

	const username = req.body.username
	const password = req.body.password

	console.log("this was the username: ", username)
	console.log("this was the password: ", password)
	const exists = await userExists(username, password)

	if (!exists) {
		res.status(401).send()
		return
	}

	const secretKey = process.env.SECRET_JWT_KEY

	if (secretKey === undefined) {
		res.status(500).send("HUHHHH?")
		console.log("no SECRET KE WHAHTHA")
		return
	}

	const token = jsonwebtoken.sign({ username }, secretKey, { expiresIn: '1h' })
	res.json({ token })
	res.status(200)
	res.send()
})


app.get('/fetch-dashboard-info', (req, res) => {

	const authHeader = req.headers.authorization

	if (!authHeader) {
		return res.status(401).send("no auth header provided")
	}

	const token = authHeader.split(' ')[1]

	try {

		const decoded = jsonwebtoken.verify(token, String(process.env.SECRET_JWT_KEY)) as JwtPayload

		res.locals.user = { username: decoded.username }

	} catch {
		return res.status(401).send("WHO DA HELL IS THIS GUY???? INVALID TOKEN")
	}

	// server.go can now be called upon to give data to this endpoint

	// get the id of the agent
	// threat summary (being hacked, being scanned, healthy)
	// health ["Critical", "High", "Medium", "Low"]
	// last check in (when agent last communicated it's status)
	// cpu and memory usage of each agent
	// total connections over time graph for each agent
	//
	// const info = AgentInfo	{
	// 	AgentId: agentID,
	// 	ThreatSummary: threatSummary,
	// 	Health: health,
	// 	LastCheckIn: lastCheckIn,
	// 	CPUUsage: cpuUsage,
	// 	RAMUsage: ramUsage,
	// }

	const info: AgentInfo = {
		AgentId: ,
		ThreatSummary: threatSummary,
		Health: health,
		LastCheckIn: lastCheckIn,
	}



	res.status(200)


	return res.send(`welcome to dashboard endpoint`)

})


const options = {
	key: fs.readFileSync(path.join(__dirname, "../key.pem")),
	cert: fs.readFileSync(path.join(__dirname, "../cert.pem")),
}


const server = https.createServer(options, app)

export default server
