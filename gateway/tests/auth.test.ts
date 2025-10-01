import 'dotenv/config'
import express from 'express'
import request from 'supertest'
import app from '../src/app'
import fs from 'fs'
import https from 'https'


describe("GET /dashboard-stream", () => {

	it("should give NO data since unauthenticated", async () => {
		await request(app).get("/dashboard-stream").expect(401)
	})


	it("should NOT give stream data", async () => {

		const csrfResp = await request(app).get("/csrf")

		const csrfToken = csrfResp.body.csrfToken


		const loginRes = await request(app).post("/login").send({ "username": "notarealuser", "password": "notarealpassword" }
		).set("Content-Type", "application/json").set("X-CSRF-Token", csrfToken).expect(401)

		const cookies = loginRes.headers['set-cookie']

		await request(app).get("/dashboard-stream").set("Accept", "text/event-stream").set("Cookie", cookies).expect(401)
	})



	it("should give stream data", async () => {


		const csrfResp = await request(app).get("/csrf")

		const csrfToken = csrfResp.body.csrfToken

		const loginRes = await request(app).post("/login").send({ "username": "admin", "password": "awnoidroppeditinthewater" }
		).set("Content-Type", "application/json").set("X-CSRF-Token", csrfToken).expect(200)

		const cookies = loginRes.headers['set-cookie']

		await request(app).get("/dashboard-stream").set("Accept", "text/event-stream").set("Cookie", cookies).expect(200)

	})
})
