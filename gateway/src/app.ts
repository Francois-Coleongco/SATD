import express from 'express'
import fs from 'fs'
import path from 'path'
import https from 'https'
import client from '@elastic/elasticsearch'
import jsonwebtoken from 'jsonwebtoken'
import { userExists } from './auth'
import { STATUS_CODES } from 'http'


const app = express()

app.use(express.json())

app.get('/dash', (req, res) => {
	res.send("FFFFFFFFFUUUUUUUCKKKKKKKKKK")
})

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

	const token = jsonwebtoken.sign({ payload: username }, secretKey, { expiresIn: '1h' })
	res.json({ token })
	res.status(200)
	res.send()

})

const options = {
	key: fs.readFileSync(path.join(__dirname, "../key.pem")),
	cert: fs.readFileSync(path.join(__dirname, "../cert.pem")),
}

const server = https.createServer(options, app)

export default server
