package infrastructure

import (
	"product-app/domain"
	"product-app/persistence"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_AddAndGetById(t *testing.T) {
	TruncateTestData(ctx, dbPool)

	repo := persistence.NewUserRepository(dbPool)
	now := time.Now()

	user := domain.User{
		Username:  "john",
		Email:     "john@test.com",
		Password:  "123456",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := repo.AddUser(user)
	assert.NoError(t, err)

	savedUser, err := repo.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, savedUser.Username)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	TruncateTestData(ctx, dbPool)

	repo := persistence.NewUserRepository(dbPool)
	now := time.Now()

	err := repo.AddUser(domain.User{
		Username:  "admin",
		Email:     "admin@test.com",
		Password:  "admin123",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.NoError(t, err)

	user, err := repo.GetByUsername("admin")
	assert.NoError(t, err)
	assert.Equal(t, "admin@test.com", user.Email)
}

func TestUserRepository_Update(t *testing.T) {
	TruncateTestData(ctx, dbPool)

	repo := persistence.NewUserRepository(dbPool)
	now := time.Now()

	_ = repo.AddUser(domain.User{
		Username:  "oldname",
		Email:     "user@test.com",
		Password:  "123",
		CreatedAt: now,
		UpdatedAt: now,
	})

	user, err := repo.GetById(1)
	assert.NoError(t, err)

	user.Username = "newname"
	user.UpdatedAt = time.Now()

	err = repo.UpdateUser(user)
	assert.NoError(t, err)

	updated, err := repo.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, "newname", updated.Username)
}

func TestUserRepository_DeleteById(t *testing.T) {
	TruncateTestData(ctx, dbPool)

	repo := persistence.NewUserRepository(dbPool)
	now := time.Now()

	_ = repo.AddUser(domain.User{
		Username:  "todelete",
		Email:     "delete@test.com",
		Password:  "123",
		CreatedAt: now,
		UpdatedAt: now,
	})

	err := repo.DeleteById(1)
	assert.NoError(t, err)

	_, err = repo.GetById(1)
	assert.Error(t, err)
}
