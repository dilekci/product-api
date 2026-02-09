package controller

import (
	"errors"
	"fmt"
	"product-app/domain"
	"product-app/persistence"
)

type FakeCategoryRepository struct {
	categories []domain.Category
}

func NewFakeCategoryRepository(initial []domain.Category) persistence.ICategoryRepository {
	return &FakeCategoryRepository{categories: initial}
}

func (repo *FakeCategoryRepository) GetAllCategories() []domain.Category {
	return repo.categories
}

func (repo *FakeCategoryRepository) GetById(categoryId int64) (domain.Category, error) {
	for _, category := range repo.categories {
		if category.Id == categoryId {
			return category, nil
		}
	}
	return domain.Category{}, errors.New(fmt.Sprintf("Category not found with id %d", categoryId))
}

func (repo *FakeCategoryRepository) AddCategory(category domain.Category) error {
	repo.categories = append(repo.categories, domain.Category{
		Id:          int64(len(repo.categories)) + 1,
		Name:        category.Name,
		Description: category.Description,
	})
	return nil
}

func (repo *FakeCategoryRepository) UpdateCategory(category domain.Category) error {
	for i, cat := range repo.categories {
		if cat.Id == category.Id {
			repo.categories[i] = category
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Category not found with id %d", category.Id))
}

func (repo *FakeCategoryRepository) DeleteById(categoryId int64) error {
	foundIndex := -1
	for i, category := range repo.categories {
		if category.Id == categoryId {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return errors.New(fmt.Sprintf("Category not found with id %d", categoryId))
	}
	repo.categories = append(repo.categories[:foundIndex], repo.categories[foundIndex+1:]...)
	return nil
}
