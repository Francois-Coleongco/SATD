import dotenv from 'dotenv'
dotenv.config()
import app from './app'
import fs from 'fs'
import path from 'path'
import https from 'https'


const options = {
	key: fs.readFileSync(path.join(__dirname, "../key.pem")),
	cert: fs.readFileSync(path.join(__dirname, "../cert.pem")),
}


const server = https.createServer(options, app)

const PORT = process.env.PORT || 3000;

server.listen(PORT, () => {
	console.log(`live on ${PORT}`)
})
