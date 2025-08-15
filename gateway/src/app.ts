import express from 'express'
import fs from 'fs'
import path from 'path'
import https from 'https'
import jsonwebtoken from 'jsonwebtoken'
import { authMiddleware, userExists } from './auth'
import { AgentInfo } from './types'
import cookieParser from 'cookie-parser'
import session from 'express-session'
import crypto from 'crypto'

interface JwtPayload {
	username: string;
	iat?: number;
	exp?: number;
}

const app = express()

let agentsMap: Map<string, AgentInfo> = new Map(); // id, AgentData


app.use(session({
	secret: crypto.randomBytes(32).toString('hex'),
	resave: false,
	saveUninitialized: true,
	cookie: {
		secure: true,
		httpOnly: true,
		sameSite: 'strict',
		maxAge: 24 * 60 * 60 * 1000
	}
}))

app.use(express.json())
app.use(cookieParser())

app.post('/login', async (req, res) => {

	// no need for csrf because we're using jwts that are httponly

	const username = req.body.username
	const password = req.body.password

	const exists = await userExists(username, password)

	if (!exists) {
		res.status(401).send()
		return
	}

	const secretKey = process.env.SECRET_JWT_KEY

	if (secretKey === undefined) {
		res.status(500).send("HUHHHH?")
		console.log("no SECRET KEY WHAHTHA")
		return
	}

	const token = jsonwebtoken.sign({ username }, secretKey, { expiresIn: '1h' })

	res.cookie('jwt', token)

	return res.json({ token })
})


app.post('/add-dashboard-info', authMiddleware, async (req, res) => {
	console.log(req.body)
	return res.status(200).send("SUCCESSFUL SEND DATA");
})


app.get('/fetch-dashboard-info', authMiddleware, async (req, res) => {

	// need to repopulate if there was a change, how do we know if there was a change?
	for (const [agentID, agentInfo] of agentsMap) {
	}

	// return only the updated ones

	return res.status(200).send(`welcome to dashboard endpoint`)

})


const options = {
	key: fs.readFileSync(path.join(__dirname, "../key.pem")),
	cert: fs.readFileSync(path.join(__dirname, "../cert.pem")),
}


const server = https.createServer(options, app)

export default server
