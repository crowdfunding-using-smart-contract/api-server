package usecase

import (
	"fmt"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProjectUseCase interface {
	CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error)
	GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error)
}

type projectUseCase struct {
	projectRepository repository.ProjectRepository
}

type ProjectUseCaseOptions struct {
	repository.ProjectRepository
}

func NewProjectUseCase(options *ProjectUseCaseOptions) ProjectUseCase {
	return &projectUseCase{
		projectRepository: options.ProjectRepository,
	}
}

func (uc *projectUseCase) CreateProject(project *entity.ProjectCreatePayload) (*entity.ProjectDto, error) {
	ownerID := uuid.MustParse(project.OwnerID)

	endDate, err := time.Parse(time.RFC3339, project.EndDate)
	if err != nil {
		return nil, err
	}

	var launchDate time.Time = time.Now()
	if project.LaunchDate != "" {
		if launchDate, err = time.Parse(time.RFC3339, project.LaunchDate); err != nil {
			fmt.Println("Error while parsing launch date: ", err)
			return nil, err
		}
	}

	payload := &entity.Project{
		Title:          project.Title,
		SubTitle:       project.SubTitle,
		CategoryID:     uuid.MustParse(project.CategoryID),
		SubCategoryID:  uuid.MustParse(project.SubCategoryID),
		Image:          project.Image,
		Description:    project.Description,
		TargetFunding:  project.TargetFunding,
		MonetaryUnit:   project.MonetaryUnit,
		StartDate:      time.Now(),
		EndDate:        endDate,
		LaunchDate:     launchDate,
		CurrentFunding: decimal.NewFromInt(0),
		OwnerID:        ownerID,
	}
	newProject, err := uc.projectRepository.Create(payload)
	if err != nil {
		return nil, err
	}

	return newProject.ToProjectDto(), nil
}

func (uc *projectUseCase) GetProjectsByOwnerID(requestOwnerID string) ([]entity.ProjectDto, error) {
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
