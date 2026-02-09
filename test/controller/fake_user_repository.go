package controller

import (
	"errors"
	"fmt"
	"product-app/internal/domain"
	"product-app/internal/ports"
)

type FakeUserRepository struct {
	users []domain.User
}

func NewFakeUserRepository(initial []domain.User) ports.UserRepository {
	return &FakeUserRepository{users: initial}
}

func (repo *FakeUserRepository) GetById(userId int64) (domain.User, error) {
	for _, user := range repo.users {
		if user.Id == userId {
			return user, nil
		}
	}
	return domain.User{}, errors.New(fmt.Sprintf("user not found with id %d", userId))
}

func (repo *FakeUserRepository) GetByUsername(username string) (domain.User, error) {
	for _, user := range repo.users {
		if user.Username == username {
			return user, nil
		}
	}
	return domain.User{}, errors.New(fmt.Sprintf("user not found with username %s", username))
}

func (repo *FakeUserRepository) GetByEmail(email string) (domain.User, error) {
	for _, user := range repo.users {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, errors.New(fmt.Sprintf("user not found with email %s", email))
}

func (repo *FakeUserRepository) AddUser(user domain.User) error {
	user.Id = int64(len(repo.users)) + 1
	repo.users = append(repo.users, user)
	return nil
}

func (repo *FakeUserRepository) UpdateUser(user domain.User) error {
	for i, u := range repo.users {
		if u.Id == user.Id {
			repo.users[i] = user
			return nil
		}
	}
	return errors.New(fmt.Sprintf("user not found with id %d", user.Id))
}

func (repo *FakeUserRepository) DeleteById(userId int64) error {
	foundIndex := -1
	for i, user := range repo.users {
		if user.Id == userId {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return errors.New(fmt.Sprintf("user not found with id %d", userId))
	}
	repo.users = append(repo.users[:foundIndex], repo.users[foundIndex+1:]...)
	return nil
}
