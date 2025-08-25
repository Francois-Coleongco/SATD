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
import cors from 'cors'
import { doubleCsrf } from 'csrf-csrf'

interface JwtPayload {
	username: string;
	iat?: number;
	exp?: number;
}

const app = express()

app.use(cors({
	origin: "http://localhost:5173",
	credentials: true
}));

app.use(session({
	secret: crypto.randomBytes(32).toString('hex'),
	resave: false,
	saveUninitialized: true,
	cookie: {
		secure: true,
		httpOnly: true,
		sameSite: 'lax',
		maxAge: 24 * 60 * 60 * 1000
	}
}))

app.use(express.json())
app.use(cookieParser())

let clients: express.Response[] = []
let agentsMap: Map<string, AgentInfo> = new Map(); // id, AgentData

const doubleCsrfUtilities = doubleCsrf({
	getSecret: () => "Secret", // A function that optionally takes the request and returns a secret
	getSessionIdentifier: (req) => req.session.id, // A function that returns the unique identifier for the request
	cookieName: "__Host-psifi.x-csrf-token", // The name of the cookie to be used, recommend using Host prefix.
	cookieOptions: {
		sameSite: "strict",
		secure: true,
		httpOnly: true,
	},
	size: 32, // The size of the random value used to construct the message used for hmac generation
	ignoredMethods: ["GET", "HEAD", "OPTIONS"], // A list of request methods that will not be protected.
	getCsrfTokenFromRequest: (req) => req.headers["x-csrf-token"], // A function that returns the token from the request
	skipCsrfProtection: undefined
});


const csrf = (req: express.Request, res: express.Response, next: express.NextFunction) => {
	try {
		doubleCsrfUtilities.validateRequest(req);
		next();
	} catch (err) {
		res.status(403).json({ error: 'Invalid CSRF token' });
	}
};


app.get('/csrf', (req, res) => {
	const csrfToken = doubleCsrfUtilities.generateCsrfToken(req, res)
	return res.json({ csrfToken: csrfToken })
})


app.post('/login', csrf, async (req, res) => {

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

	return res.json({ token: token })
})


app.post('/add-dashboard-info', authMiddleware, async (req, res) => {
	console.log(req.body)
	return res.status(200).send("SUCCESSFUL SEND DATA");
})

app.get('/dashboard-stream', authMiddleware, (req, res) => {
	res.setHeader('Content-Type', 'text/event-stream');
	res.setHeader('Cache-Control', 'no-cache');
	res.setHeader('Connection', 'keep-alive');

	clients.push(res);

	req.on('close', () => {
		clients = clients.filter(c => c !== res);
	});
});


const options = {
	key: fs.readFileSync(path.join(__dirname, "../key.pem")),
	cert: fs.readFileSync(path.join(__dirname, "../cert.pem")),
}


const server = https.createServer(options, app)

export default server
