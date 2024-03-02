package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
)

type ProjectCategoryUseCase interface {
	ListProjectCategories() ([]entity.ProjectCategoryDto, error)
}

type projectCategoryUseCase struct {
	projectCategoryRepository repository.ProjectCategoryRepository
}

type ProjectCategoryUseCaseOptions struct {
	repository.ProjectCategoryRepository
}

func NewProjectCategoryUseCase(options *ProjectCategoryUseCaseOptions) ProjectCategoryUseCase {
	return &projectCategoryUseCase{
		projectCategoryRepository: options.ProjectCategoryRepository,
	}
}

func (uc *projectCategoryUseCase) ListProjectCategories() ([]entity.ProjectCategoryDto, error) {
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
