import bcrypt from 'bcrypt'
import { Client } from 'pg'
import pool from './db'
import { NextFunction, Request, Response } from 'express'
import jsonwebtoken from 'jsonwebtoken'
import { JwtPayload } from './types'


export const authMiddleware = (req: Request, res: Response, next: NextFunction) => {
	const authHeader = req.headers.authorization
	var token = req.cookies.jwt

	if (!token && authHeader) {
		const parts = authHeader.split(' ')

		if (parts.length === 2 && parts[0] == 'Bearer') {
			token = parts[1]
		}
	}

	try {
		const decoded = jsonwebtoken.verify(token, String(process.env.SECRET_JWT_KEY)) as JwtPayload
		res.locals.user = { username: decoded.username }
		next();
	} catch {
		return res.status(401).send("WHO DA HELL IS THIS GUY???? INVALID TOKEN!!!!")
	}
}

export const userExists = async (username: string, password: string): Promise<boolean> => {

	try {
		const res = await pool.query("SELECT * FROM users WHERE username=$1", [username])

		if (res.rows.length === 0) {
			console.log("we have 0 rows")
			return false;
		}

		const user = res.rows[0]

		const storedHash = user.password_hash

		const isMatch = await bcrypt.compare(String(password), storedHash)

		return isMatch;
	}
	catch (error) {
		console.log("ERROR this was the username: ", username)
		console.log("ERROR this was the password: ", password)
		console.error("ERROR auth error: ", error)
		return false
	}

}

export const agentAuthMiddleware = () => {

}


