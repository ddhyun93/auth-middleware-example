package repository

import "go-auth-with-chi/domain"

type UserRepository interface {
	Get(ID string) (*domain.UserDAO, error)
	GetByEmail(email string) (*domain.UserDAO, error)
	Upsert(dao *domain.UserDAO) (*domain.UserDAO, error)

	Destroy() bool
}