package main

import (
	"SATD/server/auth"
	"SATD/types"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {

	jwt_received := ""

	c := http.DefaultClient

	auth.AuthToDash(c, 1, "https://localhost:3000/login", types.DashCreds{Username: "admin", Password: "awnoidroppeditinthewater"}, &jwt_received)

	println(jwt_received)
	if jwt_received == "" {
		t.Errorf("no jwt was generated")
	}
}
