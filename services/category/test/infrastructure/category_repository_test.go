package infrastructure

import (
	"testing"

	"product-app/services/category/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository_GetAll(t *testing.T) {
	setupFullTestData()

	expected := []domain.Category{
		{Id: 1, Name: "Elektronik", Description: "Elektronik ürünler"},
		{Id: 2, Name: "Beyaz Eşya", Description: "Beyaz eşya ürünleri"},
		{Id: 3, Name: "Dekorasyon", Description: "Ev dekorasyonu"},
	}

	actual := categoryRepository.GetAllCategories()

	assert.Len(t, actual, 3)
	assert.Equal(t, expected, actual)
}

func TestCategoryRepository_GetById(t *testing.T) {
	setupFullTestData()

	category, err := categoryRepository.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), category.Id)

	_, err = categoryRepository.GetById(33)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category not found")
}

func TestCategoryRepository_Add(t *testing.T) {
	clearTestData()

	newCategory := domain.Category{
		Name:        "Teknoloji",
		Description: "Teknoloji ürünleri",
	}

	err := categoryRepository.AddCategory(newCategory)
	assert.NoError(t, err)

	categories := categoryRepository.GetAllCategories()
	assert.Len(t, categories, 1)
	assert.Equal(t, "Teknoloji", categories[0].Name)
}

func TestCategoryRepository_Update(t *testing.T) {
	setupFullTestData()

	before, _ := categoryRepository.GetById(1)
	assert.Equal(t, "Elektronik", before.Name)

	updatedCategory := domain.Category{
		Id:          1,
		Name:        "Güncel Elektronik",
		Description: "Güncel elektronik ürünler",
	}

	err := categoryRepository.UpdateCategory(updatedCategory)
	assert.NoError(t, err)

	after, _ := categoryRepository.GetById(1)
	assert.Equal(t, "Güncel Elektronik", after.Name)
}

func TestCategoryRepository_DeleteById(t *testing.T) {
	setupFullTestData()

	err := categoryRepository.DeleteById(1)
	assert.NoError(t, err)

	_, err = categoryRepository.GetById(1)
	assert.Error(t, err)
}
