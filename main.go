package main

import (
	"go-auth-with-chi/api"
	"go-auth-with-chi/cmd"
)

func main() {
	cmd.SetBaseConfig()
	api.StartAPIServer()
}
