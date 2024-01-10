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
	GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error)
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

	newProject.Owner = *user

	return newProject.ToProjectDto(), nil
}

func (uc *projectUsecase) GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error) {
	ownerID := uuid.MustParse(requestOwnerID)

	projects, err := uc.projectRepository.FindAllByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}

	projectDtos := make([]entity.ProjectDto, 0, len(projects))
	for _, project := range projects {
		projectDtos = append(projectDtos, *project.ToProjectDto())
	}

	return projectDtos, nil
}
