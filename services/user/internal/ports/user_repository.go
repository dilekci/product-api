package ports

import "product-app/services/user/internal/domain"

type UserRepository interface {
	GetById(userId int64) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	AddUser(user domain.User) error
	UpdateUser(user domain.User) error
	DeleteById(userId int64) error
}
