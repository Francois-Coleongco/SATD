import bcrypt from 'bcrypt'
import { Client } from 'pg'
import pool from './db'
import { NextFunction, Request, Response } from 'express'
import jsonwebtoken from 'jsonwebtoken'
import { JwtPayload } from './types'


export const authMiddleware = (req: Request, res: Response, next: NextFunction) => {
	var token = req.cookies.jwt

	try {
		const decoded = jsonwebtoken.verify(token, String(process.env.SECRET_JWT_KEY)) as JwtPayload
		console.log("this was decoded", decoded)
		if (decoded.role !== "user") {
			return res.status(401).send("INVALID TOKEN!!!!")
		}
		res.locals.user = { username: decoded.username }
		next();
	} catch {
		return res.status(401).send("INVALID TOKEN!!!!")
	}
}

export const userExists = async (username: String, password: String): Promise<boolean> => {

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
		console.error("ERROR auth error: ", error)
		return false
	}

}

export const analysisServerAuthMiddleware = (req: Request, res: Response, next: NextFunction) => {

	const authHeader = req.headers.authorization

	if (authHeader === undefined) {
		return res.status(401).send("INVALID TOKEN!!!!")
	}

	const parts = authHeader.split(' ')

	var token = ""

	if (parts.length === 2 && parts[0] == 'Bearer') {
		token = parts[1]
	}


	try {
		const decoded = jsonwebtoken.verify(token, String(process.env.SECRET_JWT_KEY)) as JwtPayload

		if (decoded.role !== "server") {
			return res.status(401).send("INVALID TOKEN!!!!")
		}

		res.locals.user = { username: decoded.username }
		next();
	} catch {
		return res.status(401).send("INVALID TOKEN!!!!")
	}

}


