import bcrypt from 'bcrypt'
import { Client } from 'pg'
import pool from './db'
import { Request, Response } from 'express'
import jsonwebtoken from 'jsonwebtoken'
import { JwtPayload } from './types'


export const authMiddleware = (req: Request, res: Response) => {
	const token = req.cookies.jwt

	if (!token) {
		return res.status(401).send("no auth header provided")
	}

	try {
		const decoded = jsonwebtoken.verify(token, String(process.env.SECRET_JWT_KEY)) as JwtPayload

		res.locals.user = { username: decoded.username }
	} catch {
		return res.status(401).send("WHO DA HELL IS THIS GUY???? INVALID TOKEN!!!!")
	}
}

export const userExists = async (username: string, password: string): Promise<boolean> => {

	console.log("this was the username: ", username)
	console.log("this was the password: ", password)
	try {
		const res = await pool.query("SELECT * FROM users WHERE username=$1", [username])

		console.log("after")
		if (res.rows.length === 0) {
			console.log("we have 0 rows")
			return false;
		}

		const user = res.rows[0]

		console.log("did user exist?? ", user)
		const storedHash = user.password_hash

		const isMatch = await bcrypt.compare(String(password), storedHash)

		console.log("YES IT DID")
		return isMatch;
	}
	catch (error) {
		console.log("ERROR this was the username: ", username)
		console.log("ERROR this was the password: ", password)
		console.error("ERROR auth error: ", error)
		return false
	}

}

