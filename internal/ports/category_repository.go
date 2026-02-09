package ports

import "product-app/internal/domain"

type CategoryRepository interface {
	GetAllCategories() []domain.Category
	GetById(categoryId int64) (domain.Category, error)
	AddCategory(category domain.Category) error
	UpdateCategory(category domain.Category) error
	DeleteById(categoryId int64) error
}
