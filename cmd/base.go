package cmd

import (
	"go-auth-with-chi/mongo"
	"os"
)

func SetBaseConfig() {
	env := "dev"
	if os.Getenv("ENV") != "" {
		env = os.Getenv("ENV")
	}

	conn := mongo.NewMongoDB(env)
	conn.RegisterRepos()
}