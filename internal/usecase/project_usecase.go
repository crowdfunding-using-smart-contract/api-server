package usecase

import (
	"fmt"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProjectUsecase interface {
	CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error)
}

type projectUsecase struct {
	projectRepository repository.ProjectRepository
	userRepository    repository.UserRepository
}

type ProjectUsecaseOptions struct {
	repository.ProjectRepository
	repository.UserRepository
}

func NewProjectUsecase(options *ProjectUsecaseOptions) ProjectUsecase {
	return &projectUsecase{
		projectRepository: options.ProjectRepository,
		userRepository:    options.UserRepository,
	}
}

func (uc *projectUsecase) CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error) {
	ownerID := uuid.MustParse(project.OwnerID)

	user, err := uc.userRepository.FindById(ownerID)
	if err != nil {
		return nil, err
	}

	if user.Role != entity.Creator {
		err = fmt.Errorf("user is not a creator")
		return nil, err
	}

	payload := &entity.Project{
		Title:         project.Title,
		Description:   project.Description,
		Image:         project.Image,
		TargetAmount:  project.TargetAmount,
		CurrentAmount: decimal.NewFromInt(0),
		OwnerID:       ownerID,
	}
	newProject, err := uc.projectRepository.Create(payload)
	if err != nil {
		return nil, err
	}

	return newProject.ToProjectDto(), nil

}
