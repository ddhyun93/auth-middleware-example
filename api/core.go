package api

import (
	"fmt"
	"github.com/go-chi/chi"
	"go-auth-with-chi/middleware"
	"log"
	"net/http"
	"os"
)

func StartAPIServer() {
	port := ":8080"
	env := "dev"
	if os.Getenv("ENV") != "" {
		env = os.Getenv("ENV")
	}
	mux := chi.NewRouter()

	mux.Use(middleware.AuthMiddleware())
	AddUserHandler(mux)

	info := fmt.Sprintf("start api server .... port: %s env: %s", port, env)
	log.Println(info)
	_ = http.ListenAndServe(port, mux)
}
