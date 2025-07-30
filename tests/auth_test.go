package main

import (
	"SATD/server/auth"
	"SATD/types"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	c := http.DefaultClient

	auth.AuthToDash(c, 1, "https://localhost:3000/login", types.DashCreds{Username: "admin", Password: "awnoidroppeditinthewater"})
}
