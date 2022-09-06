package ioc

import "go-auth-with-chi/repository"

type iocRepo struct {
	Users repository.UserRepository
}

var Repo iocRepo
