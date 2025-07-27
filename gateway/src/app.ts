import express from 'express'
import fs from 'fs'
import path from 'path'
import https from 'https'
import jsonwebtoken from 'jsonwebtoken'
import { authMiddleware, userExists } from './auth'
import { AgentInfo } from './types'
import { csrfSync } from 'csrf-sync'
import cookieParser from 'cookie-parser'
import session, { Session } from 'express-session'
import crypto from 'crypto'

interface JwtPayload {
	username: string;
	iat?: number;
	exp?: number;
}

const app = express()

let agentsMap: Map<string, AgentInfo> = new Map(); // id, AgentData

const {
	invalidCsrfTokenError, // This is just for convenience if you plan on making your own middleware.
	generateToken, // Use this in your routes to generate, store, and get a CSRF token.
	getTokenFromRequest, // use this to retrieve the token submitted by a user
	getTokenFromState, // The default method for retrieving a token from state.
	storeTokenInState, // The default method for storing a token in state.
	revokeToken, // Revokes/deletes a token by calling storeTokenInState(undefined)
	csrfSynchronisedProtection, // This is the default CSRF protection middleware.
} = csrfSync({
	getTokenFromRequest: (req) => {
		return req.cookies['XSRF-TOKEN']
	}
});

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

app.use(csrfSynchronisedProtection)

app.use((req, res, next) => {
	console.log('Session:', req.session);
	console.log('Cookies:', req.cookies);
	next();
});

app.get('/get-csrf-token', (req, res) => {
	const token = generateToken(req)
	res.cookie('XSRF-TOKEN', token, {
		httpOnly: false,
		secure: true,
		sameSite: 'strict'

	})
	res.json({ csrfToken: token })
})

app.post('/login', async (req, res) => {

	// need to implement csrf protection here
	console.log("this was my csrf token: ", req.csrfToken)


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

	// need to remove this once ui is complete

	res.cookie('jwt', token)

	return res.json({ token })
})


app.post('/add-dashboard-info', agentAuthMiddleware, async (req, res) => {
	return status(200).send("SUCCESSFUL SEND DATA");
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
