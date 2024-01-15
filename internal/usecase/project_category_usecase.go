package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
)

type ProjectCategoryUsecase interface {
	ListProjectCategories() ([]entity.ProjectCategoryDto, error)
}

type projectCategoryUsecase struct {
	projectCategoryRepository repository.ProjectCategoryRepository
}

type ProjectCategoryUsecaseOptions struct {
	repository.ProjectCategoryRepository
}

func NewProjectCategoryUsecase(options *ProjectCategoryUsecaseOptions) ProjectCategoryUsecase {
	return &projectCategoryUsecase{
		projectCategoryRepository: options.ProjectCategoryRepository,
	}
}

func (uc *projectCategoryUsecase) ListProjectCategories() ([]entity.ProjectCategoryDto, error) {
	categories, err := uc.projectCategoryRepository.FindAll()
	if err != nil {
		return nil, err
	}

	categoryDtos := make([]entity.ProjectCategoryDto, 0, len(categories))
	for _, category := range categories {
		categoryDtos = append(categoryDtos, *category.ToProjectCategoryDto())
	}

	return categoryDtos, nil
}
