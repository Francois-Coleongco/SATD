import bcrypt from 'bcrypt'
import { Client } from 'pg'
import pool from './db'

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

		// const user = res.rows[0]

		// console.log("did user exist?? ", user)
		// const storedHash = user.password_hash

		// console.log("okay was i undefined? ", storedHash)
		//
		// const isMatch = await bcrypt.compare(String(password), storedHash)

		return true;
	}
	catch (error) {
		console.log("this was the username: ", username)
		console.log("this was the password: ", password)
		console.error("auth error: ", error)
		return false
	}



}
