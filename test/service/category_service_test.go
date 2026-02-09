package service

import (
	"testing"

	"product-app/internal/domain"
	"product-app/internal/usecase"

	"github.com/stretchr/testify/assert"
)

func setupCategoryService() usecase.ICategoryService {
	initialCategories := []domain.Category{
		{Id: 1, Name: "Electronics", Description: "Electronic items"},
		{Id: 2, Name: "Books", Description: "Books and magazines"},
	}

	fakeRepository := NewFakeCategoryRepository(initialCategories)
	return usecase.NewCategoryService(fakeRepository)
}

func Test_ShouldGetAllCategories(t *testing.T) {
	categoryService := setupCategoryService()
	category := categoryService.GetAllCategories()
	assert.Len(t, category, 2)
}

func Test_ShouldCategoryGetById(t *testing.T) {
	categoryService := setupCategoryService()
	category, err := categoryService.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), category.Id)
	assert.Equal(t, "Electronics", category.Name)
	assert.Equal(t, "Electronic items", category.Description)
}

func Test_ShouldCategoryDeleteById(t *testing.T) {
	categoryService := setupCategoryService()
	err := categoryService.DeleteById(1)
	assert.NoError(t, err)
	_, err = categoryService.GetById(1)
	assert.Error(t, err)
}

func Test_ShouldUpdateCategory(t *testing.T) {
	categoryService := setupCategoryService()
	before, err := categoryService.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, "Electronics", before.Name)
	assert.Equal(t, "Electronic items", before.Description)

	updatedCategory := domain.Category{
		Id:          1,
		Name:        "Electronics and Gadgets",
		Description: "Electronic items and smart gadgets",
	}
	err = categoryService.UpdateCategory(updatedCategory)
	assert.NoError(t, err)

	after, err := categoryService.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, "Electronics and Gadgets", after.Name)
	assert.Equal(t, "Electronic items and smart gadgets", after.Description)
	assert.Equal(t, int64(1), after.Id)
}

func Test_ShouldAddCategory_Success(t *testing.T) {
	categoryService := setupCategoryService()
	beforeCategories := categoryService.GetAllCategories()
	assert.Equal(t, 2, len(beforeCategories))

	newCategory := domain.Category{
		Name:        "Clothing",
		Description: "Apparel and accessories",
	}
	err := categoryService.AddCategory(newCategory)
	assert.NoError(t, err)

	afterCategories := categoryService.GetAllCategories()
	assert.Equal(t, 3, len(afterCategories))

	addedCategory := afterCategories[2]
	assert.Equal(t, int64(3), addedCategory.Id)
	assert.Equal(t, "Clothing", addedCategory.Name)
	assert.Equal(t, "Apparel and accessories", addedCategory.Description)
}
