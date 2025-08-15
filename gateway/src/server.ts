import dotenv from 'dotenv'
dotenv.config()
import server from './app'

const PORT = process.env.PORT || 3000;

server.listen(PORT, () => {
	console.log(`live on ${PORT}`)
})
